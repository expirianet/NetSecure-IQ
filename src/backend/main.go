package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

var db *sql.DB

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Separate struct for registration
type RegisterRequest struct {
	Email string `json:"email"`
}

// Separate struct for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type OrganizationRequest struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	VATNumber    string `json:"vat_number"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone"`
	UserID       string `json:"user_id"`
}

var orgReq OrganizationRequest

func main() {
	var err error

	// Récupérer les infos de connexion depuis les variables d'environnement
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	if host == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("One or more required environment variables are not set: POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB")
	}

	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// Test connection to the database
	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection test failed:", err)
	}
	fmt.Println("✅ Connected to PostgreSQL successfully")

	http.HandleFunc("/api/register", withCORS(handleRegister))
	http.HandleFunc("/api/login", withCORS(handleLogin))
	http.HandleFunc("/api/protected", withCORS(jwtMiddleware(handleProtected)))
	http.HandleFunc("/api/data/routers", withCORS(jwtMiddleware(handleRouters)))
	http.HandleFunc("/api/complete-organization", withCORS(handleCompleteOrganization))

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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

	// Generate a random secure password
	generatedPassword, err := password.Generate(16, 4, 4, false, false)
	if err != nil {
		http.Error(w, "Failed to generate password", http.StatusInternalServerError)
		return
	}

	// Hash the generated password
	hash, err := bcrypt.GenerateFromPassword([]byte(generatedPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Save user to DB with default role_id = 3 (ex: tenant)
	_, err = db.Exec(`INSERT INTO users (email, password_hash, role_id) VALUES ($1, $2, $3)`, req.Email, string(hash), 2)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Database insert failed: " + err.Error(),
		})
		return
	}

	// Send email with temporary password
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

	d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 2525, "f8a85b7da3ad77", "8dfd2cbb9ac6f2")

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
	err := db.QueryRow(`
		SELECT u.password_hash, r.name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.email = $1
	`, req.Email).Scan(&hash, &roleName)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Génération du token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   req.Email,
		"role":    roleName,
		"exp":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Expire dans 24h
		"user_id": orgReq.UserID,
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Réponse JSON avec le token
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   tokenString,
	})
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
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

		// Optionnel : tu peux récupérer les claims ici et les passer dans le contexte

		next(w, r)
	}
}

func handleRouters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	// Exemple statique — à remplacer par une vraie requête en base si besoin
	routers := []map[string]interface{}{
		{"mac": "00:11:22:33:44:55", "value": "online", "time": time.Now().Format(time.RFC3339)},
		{"mac": "66:77:88:99:AA:BB", "value": "offline", "time": time.Now().Add(-10 * time.Minute).Format(time.RFC3339)},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routers)
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

	// Insert new organization
	var orgID string
	err := db.QueryRow(`
		INSERT INTO organizations (name, address, vat_number, contact_email, contact_phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, now(), now())
		RETURNING id
	`, req.Name, req.Address, req.VATNumber, req.ContactEmail, req.ContactPhone).Scan(&orgID)

	if err != nil {
		http.Error(w, "Failed to insert organization: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update user with organization_id
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
