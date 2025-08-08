package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

var db *sql.DB
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

var (
	influxToken  = os.Getenv("INFLUXDB_TOKEN")
	influxBucket = os.Getenv("INFLUXDB_BUCKET")
	influxOrg    = os.Getenv("INFLUXDB_ORG")
	influxURL    = os.Getenv("INFLUXDB_URL")
	influxClient influxdb2.Client
)

type RegisterRequest struct {
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type OrganizationRequest struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	VATNumber     string `json:"vat_number"`
	State         string `json:"state"`
	City          string `json:"city"`
	ZipCode       string `json:"zip_code"`
	ContactEmail  string `json:"contact_email"`
	PecEmail      string `json:"pec_email"`
	SdiCode       string `json:"sdi_code"`
	ContactPhone  string `json:"contact_phone"`
	PersonnelInfo string `json:"personnel_info"`
	UserID        string `json:"user_id"`
}

type RouterStatus struct {
	MAC   string `json:"mac"`
	Value string `json:"value"`
	Time  string `json:"time"`
}

type MikroTikData struct {
	MAC    string `json:"mac"`
	Status string `json:"status"` // "online", "offline", etc.
}

type CreateUserRequest struct {
	Email          string  `json:"email"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	OrganizationID *string `json:"organization_id"`
	Role           string  `json:"role"` // <-- Add this field
}

type MikroTikRegisterRequest struct {
	MAC    string  `json:"mac"`
	SiteID *string `json:"site_id"` // Optional for now
}

//var org OrganizationRequest

func main() {
	if err := godotenv.Load("../../config/backend_env/backend.env"); err != nil {
		log.Println("⚠️ .env file not found at custom path — using system env vars")
	}

	var err error

	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	if host == "" || user == "" || password == "" || dbname == "" {
		log.Println("⚠️ One or more required environment variables are not set: POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB")
	} else if influxToken == "" || influxOrg == "" || influxBucket == "" || influxURL == "" {
		log.Println("⚠️ InfluxDB environment variables are missing or incomplete.")
	} else {
		connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Println("❌ Failed to open DB connection:", err)
		} else if err = db.Ping(); err != nil {
			log.Println("❌ Database connection test failed:", err)
		} else {
			fmt.Println("✅ Connected to PostgreSQL successfully")
		}

		influxClient = influxdb2.NewClient(influxURL, influxToken)
		fmt.Println("✅ Connected to InfluxDB")
		defer influxClient.Close()
	}

	if len(jwtSecret) == 0 {
		log.Println("⚠️ JWT_SECRET environment variable is not set — JWT will not work properly")
	}

	http.HandleFunc("/api/register", withCORS(handleRegister))
	http.HandleFunc("/api/login", withCORS(handleLogin))
	http.HandleFunc("/api/ping", withCORS(handlePing))
	http.HandleFunc("/api/protected", withCORS(jwtMiddleware(handleProtected)))
	http.HandleFunc("/api/data/routers", withCORS(jwtMiddleware(handleRouters)))
	http.HandleFunc("/api/complete-organization", withCORS(handleCompleteOrganization))
	http.HandleFunc("/api/users", withCORS(handleCreateUser))
	http.HandleFunc("/api/mikrotik/resource", withCORS(handleMikrotikResource))
	http.HandleFunc("/api/mikrotik/test", withCORS(handleMikrotikTest))
	http.HandleFunc("/api/mikrotik/preregister", withCORS(handleMikrotikPreRegister))

	fmt.Println("🔐 JWT Secret:", jwtSecret)
	fmt.Println("📦 Influx URL:", influxURL)
	fmt.Println("🚀 Server started at http://localhost:8080 (even if DB is down)")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		fmt.Println("📦 Origin:", origin) // 👈 Add this line
		allowedOrigins := []string{
			"http://localhost:8080",
			"http://localhost:8081",
		}

		for _, o := range allowedOrigins {
			if origin == o {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	generatedPassword, err := password.Generate(16, 4, 4, false, false)
	if err != nil {
		http.Error(w, "Failed to generate password", http.StatusInternalServerError)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(generatedPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`INSERT INTO users (email, password_hash, role_id) VALUES ($1, $2, $3)`, req.Email, string(hash), 2)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Database insert failed: " + err.Error(),
		})
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "NetSecure IQ <noreply@netsecure.test>")
	m.SetHeader("To", req.Email)
	m.SetHeader("Subject", "Your NetSecure IQ Access")
	m.SetBody("text/plain", fmt.Sprintf(`Hello,

Welcome to NetSecure IQ!
Your temporary password is:

%s

Please log in and complete your profile.

– NetSecure IQ Team`, generatedPassword))

	d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 2525, "a7579402169dd8", "6644803192d28b")
	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send email:", err)
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully. Check your email.",
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var hash, roleName string
	var userID, orgID sql.NullString
	err := db.QueryRow(`
		SELECT u.id, u.password_hash, r.name, u.organization_id
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.email = $1
	`, req.Email).Scan(&userID, &hash, &roleName, &orgID)

	if err == sql.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   req.Email,
		"role":    roleName,
		"user_id": userID.String,
		"exp":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Final response
	resp := map[string]interface{}{
		"message": "Login successful",
		"token":   tokenString,
		"role":    roleName,
		"user_id": userID.String,
	}
	if orgID.Valid {
		resp["organization_id"] = orgID.String
	} else {
		resp["organization_id"] = nil
	}

	json.NewEncoder(w).Encode(resp)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var data MikroTikData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	fmt.Printf("📡 Ping received: %+v\n", data)

	writeAPI := influxClient.WriteAPIBlocking(influxOrg, influxBucket)

	p := influxdb2.NewPointWithMeasurement("device_status").
		AddTag("mac", data.MAC).
		AddField("status", data.Status).
		SetTime(time.Now())

	if err := writeAPI.WritePoint(context.Background(), p); err != nil {
		http.Error(w, "Failed to write to InfluxDB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleProtected(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(claims)
}

func jwtMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		const prefix = "Bearer "
		if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len(prefix):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func handleRouters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	queryAPI := influxClient.QueryAPI(influxOrg)
	query := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: -7d)
			|> filter(fn: (r) => r._measurement == "device_status")
			|> filter(fn: (r) => r._field == "status")
			|> last()
	`, influxBucket)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "InfluxDB query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var output []RouterStatus

	for result.Next() {
		record := result.Record()
		mac := record.ValueByKey("mac")
		if mac == nil {
			continue
		}
		output = append(output, RouterStatus{
			MAC:   fmt.Sprintf("%v", mac),
			Value: fmt.Sprintf("%v", record.Value()),
			Time:  record.Time().Format(time.RFC3339),
		})
	}

	if result.Err() != nil {
		http.Error(w, "Influx parse error: "+result.Err().Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func handleCompleteOrganization(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req OrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	fmt.Printf("📥 Received OrganizationRequest:\n%+v\n", req)

	var orgID string
	err := db.QueryRow(`
		INSERT INTO organizations (
			name, address, vat_number, state, city, zip_code,
			contact_email, pec_email, sdi_code, contact_phone,
			personnel_info, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6,
				$7, $8, $9, $10,
				$11, now(), now())
		RETURNING id
	`, req.Name, req.Address, req.VATNumber, req.State, req.City, req.ZipCode,
		req.ContactEmail, req.PecEmail, req.SdiCode, req.ContactPhone,
		req.PersonnelInfo).Scan(&orgID)

	if err != nil {
		http.Error(w, "Failed to insert organization: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`UPDATE users SET organization_id = $1, updated_at = now() WHERE id = $2`, orgID, req.UserID)

	if err != nil {
		http.Error(w, "Failed to update user with organization: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message":         "Organization created and linked successfully",
		"organization_id": orgID,
	})
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.FirstName == "" || req.LastName == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Determine role ID from string
	var roleID int
	switch req.Role {
	case "admin":
		roleID = 1
	case "operator":
		roleID = 2
	default:
		roleID = 3 // default to User
	}

	// Generate and hash password
	generatedPassword, err := password.Generate(16, 4, 4, false, false)
	if err != nil {
		http.Error(w, "Failed to generate password", http.StatusInternalServerError)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(generatedPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert user
	_, err = db.Exec(`
    INSERT INTO users (email, password_hash, role_id, first_name, last_name, organization_id, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, now(), now())
  `, req.Email, string(hash), roleID, req.FirstName, req.LastName, req.OrganizationID)

	if err != nil {
		http.Error(w, "Database insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send email
	m := gomail.NewMessage()
	m.SetHeader("From", "NetSecure IQ <noreply@netsecure.test>")
	m.SetHeader("To", req.Email)
	m.SetHeader("Subject", "Your NetSecure IQ Account")
	m.SetBody("text/plain", fmt.Sprintf(`Hello %s,

Your account has been created in NetSecure IQ.
Login Email: %s
Your role is: %s
Temporary Password: %s

Please log in and change your password as soon as possible.

– NetSecure IQ Team`, req.FirstName, req.Email, req.Role, generatedPassword))

	d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 2525, "a7579402169dd8", "6644803192d28b")
	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send email:", err)
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created and password emailed",
	})
}

// Connects to MikroTik RouterOS API and fetches system resource info
func getMikrotikSystemResource() (map[string]interface{}, error) {
	host := os.Getenv("MIKROTIK_HOST")
	port := os.Getenv("MIKROTIK_PORT")
	user := os.Getenv("MIKROTIK_USER")
	pass := os.Getenv("MIKROTIK_PASSWORD")

	if host == "" || port == "" || user == "" || pass == "" {
		return nil, fmt.Errorf("missing MikroTik REST API credentials in environment variables")
	}

	url := fmt.Sprintf("http://%s:%s/rest/system/resource", host, port)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.SetBasicAuth(user, pass)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach MikroTik REST API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return data, nil
}

func handleMikrotikResource(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := getMikrotikSystemResource()
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ Print the response to console for debugging
	fmt.Println("📦 MikroTik REST API Response:")
	prettyJSON, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(prettyJSON))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func getAllMikrotikMetricsAsText() (string, error) {
	host := os.Getenv("MIKROTIK_HOST")
	port := os.Getenv("MIKROTIK_PORT")
	user := os.Getenv("MIKROTIK_USER")
	pass := os.Getenv("MIKROTIK_PASSWORD")

	if host == "" || port == "" || user == "" || pass == "" {
		return "", fmt.Errorf("missing MikroTik REST API credentials")
	}

	baseURL := fmt.Sprintf("http://%s:%s/rest", host, port)
	endpoints := map[string]string{
		"System Resource":             "/system/resource",
		"Interfaces":                  "/interface",
		"ARP Table":                   "/ip/arp",
		"Firewall Service Ports":      "/ip/firewall/service-port",
		"Firewall Connections":        "/ip/firewall/connection",
		"Wireless Registration Table": "/interface/wireless/registration-table",
	}

	var builder strings.Builder
	client := &http.Client{Timeout: 5 * time.Second}

	for label, ep := range endpoints {
		url := baseURL + ep
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			builder.WriteString(fmt.Sprintf("Section: %s\nError: failed to build request: %v\n\n", label, err))
			continue
		}
		req.SetBasicAuth(user, pass)

		resp, err := client.Do(req)
		if err != nil {
			builder.WriteString(fmt.Sprintf("Section: %s\nError: %v\n\n", label, err))
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		builder.WriteString(fmt.Sprintf("===== %s =====\n", label))

		// Attempt pretty JSON print
		var prettyJSON bytes.Buffer
		if json.Indent(&prettyJSON, body, "", "  ") == nil {
			builder.Write(prettyJSON.Bytes())
		} else {
			builder.Write(body) // fallback
		}
		builder.WriteString("\n\n")
	}

	return builder.String(), nil
}

func handleMikrotikTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	result, err := getAllMikrotikMetricsAsText()
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}

//-------------------------------------------------
// WireGuard pre-registration for MikroTik routers
//-------------------------------------------------

// WireGuard IP Pool
var (
	wireguardSubnet = "10.100.99.0/24"
	startOffset     = 10
)

// Generates real WireGuard keys using wg tool
func generateRealWGKeys() (string, string, error) {
	// Run "wg genkey"
	privateCmd := exec.Command("wg", "genkey")
	privateOut := &bytes.Buffer{}
	privateCmd.Stdout = privateOut
	if err := privateCmd.Run(); err != nil {
		return "", "", fmt.Errorf("wg genkey failed: %v", err)
	}
	privateKey := bytes.TrimSpace(privateOut.Bytes())

	// Run "wg pubkey"
	pubCmd := exec.Command("wg", "pubkey")
	pubCmd.Stdin = bytes.NewReader(privateKey)
	pubOut := &bytes.Buffer{}
	pubCmd.Stdout = pubOut
	if err := pubCmd.Run(); err != nil {
		return "", "", fmt.Errorf("wg pubkey failed: %v", err)
	}
	publicKey := bytes.TrimSpace(pubOut.Bytes())

	return string(privateKey), string(publicKey), nil
}

// Assigns next available IP from pool
func getNextAvailableIP() (string, error) {
	_, ipnet, err := net.ParseCIDR(wireguardSubnet)
	if err != nil {
		return "", fmt.Errorf("invalid subnet: %w", err)
	}

	baseIP := ipnet.IP.To4()
	if baseIP == nil {
		return "", fmt.Errorf("not a valid IPv4 subnet")
	}

	// Step 1: Load all currently assigned IPs from DB
	rows, err := db.Query(`SELECT vpn_internal_ip FROM mikrotik_routers`)
	if err != nil {
		return "", fmt.Errorf("failed to query IPs from database: %w", err)
	}
	defer rows.Close()

	usedIPs := make(map[string]bool)
	for rows.Next() {
		var ipStr string
		if err := rows.Scan(&ipStr); err == nil {
			usedIPs[ipStr] = true
		}
	}

	// Step 2: Iterate over IP pool starting from .10
	for i := startOffset; i < 255; i++ {
		candidateIP := net.IPv4(baseIP[0], baseIP[1], baseIP[2], byte(i)).String()
		if !usedIPs[candidateIP] {
			return candidateIP, nil
		}
	}

	return "", fmt.Errorf("no available IPs in subnet")
}

// MikroTik registration logic
func MikroTikPreRegister(mac string) (map[string]string, error) {
	fmt.Println("🔧 Starting MikroTik registration for MAC:", mac)

	// Step 1: Generate WireGuard keys
	privateKey, publicKey, err := generateRealWGKeys()
	if err != nil {
		return nil, fmt.Errorf("key generation failed: %w", err)
	}
	fmt.Println("✅ WireGuard Keys Generated")
	fmt.Println("🔑 Private Key:", privateKey)
	fmt.Println("🔐 Public Key:", publicKey)

	// Step 2: Assign IP
	ip, err := getNextAvailableIP()
	if err != nil {
		return nil, fmt.Errorf("IP allocation failed: %w", err)
	}
	fmt.Println("✅ Assigned IP:", ip)

	// Step 3: Generate MikroTik config
	mikrotikConfig := generateMikroTikConfig(privateKey, ip)
	fmt.Println("📄 MikroTik Configuration Script:\n", mikrotikConfig)

	// Step 4: Generate server-side peer config
	serverConf := fmt.Sprintf(`[Peer]
PublicKey = %s
AllowedIPs = %s/32`, publicKey, ip)
	fmt.Println("📄 WireGuard Peer (Server) Config:\n", serverConf)

	// Step 5: Store into DB
	_, err = db.Exec(`
		INSERT INTO mikrotik_routers (
			id, site_id, mac_address,
			vpn_private_key, vpn_public_key, vpn_internal_ip,
			firmware_version, model, serial_number,
			provisioning_status, created_at, updated_at
		)
		VALUES (
			gen_random_uuid(), NULL, $1,
			$2, $3, $4,
			NULL, 'unknown', 'unknown',
			'PENDING', now(), now()
		)
	`, mac, privateKey, publicKey, ip)
	if err != nil {
		return nil, fmt.Errorf("database insert failed: %w", err)
	}

	// Step 6: Build return object
	result := map[string]string{
		"internal_ip":     ip,
		"private_key":     privateKey,
		"public_key":      publicKey,
		"mikrotik_config": mikrotikConfig,
		"server_peer":     serverConf,
	}

	if err := addWireGuardPeer(publicKey, ip); err != nil {
		log.Println("⚠️ Could not add peer to WireGuard:", err)
		// You can still return success here — the router will connect once wg0 is restarted or manually updated
	}

	return result, nil
}

func generateMikroTikConfig(privateKey, ip string) string {
	return fmt.Sprintf(`/interface wireguard add name=wg1 private-key="%s"
/ip address add address=%s/32 interface=wg1
/interface wireguard peers add allowed-address=0.0.0.0/0 endpoint-address=NETSECURE_PUBLIC_IP endpoint-port=51820 interface=wg1 public-key=NETSECURE_PUBLIC_KEY`,
		privateKey, ip)
}

func handleMikrotikPreRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MikroTikRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 🛑 Check for duplicates
	var exists bool
	err := db.QueryRow(`SELECT EXISTS (SELECT 1 FROM mikrotik_routers WHERE mac_address = $1)`, req.MAC).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "MAC already registered", http.StatusConflict)
		return
	}

	// ✅ Call MikroTikPreRegister and get config data
	data, err := MikroTikPreRegister(req.MAC)
	if err != nil {
		http.Error(w, "Pre-registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ Return JSON with config
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func addWireGuardPeer(publicKey, ip string) error {
	cmd := exec.Command("docker", "exec", "netsecure-iq-wireguard-1",
		"wg", "set", "wg0",
		"peer", publicKey,
		"allowed-ips", ip+"/32")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add WireGuard peer: %v\nOutput: %s", err, string(output))
	}

	fmt.Println("✅ WireGuard peer added via docker exec")
	return nil
}
