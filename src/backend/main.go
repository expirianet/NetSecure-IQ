package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// MikroTikData represents the structure of the data sent from MikroTik
type MikroTikData struct {
	MAC    string `json:"mac"`
	Status string `json:"status"` // "online", "offline", etc.
}

// RouterStatus represents a single router record from InfluxDB
type RouterStatus struct {
	MAC   string `json:"mac"`
	Value string `json:"value"`
	Time  string `json:"time"`
}

// RegisterRequest represents a user registration payload
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var db *sql.DB

func main() {
	var err error
	connStr := "host=postgresql user=netsecure_iq password=your_secure_password dbname=netsecure_iq_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	http.HandleFunc("/api/ping", withCORS(handlePing))
	http.HandleFunc("/api/register", withCORS(handleRegister))
	http.HandleFunc("/api/login", withCORS(handleLogin))
	http.HandleFunc("/api/data/routers", withCORS(handleRouterData))
	http.HandleFunc("/api/register-mac", withCORS(handleRegisterMAC)) // ✅ NEW

	fmt.Println("Go Backend is running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	var data MikroTikData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Check if MAC is already known (exists in Influx)
	const (
		token  = "my-token"
		bucket = "netsecure"
		org    = "netsecure-org"
		url    = "http://influxdb:8086"
	)

	client := influxdb2.NewClient(url, token)
	defer client.Close()

	queryAPI := client.QueryAPI(org)
	query := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: -30d)
			|> filter(fn: (r) => r._measurement == "device_status")
			|> filter(fn: (r) => r.mac == "%s")
			|> limit(n:1)
	`, bucket, data.MAC)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil || !result.Next() {
		http.Error(w, "MAC not known. Please register it first.", http.StatusUnauthorized)
		return
	}

	// MAC exists — update with new status
	if err := saveToInflux(data.MAC, data.Status); err != nil {
		log.Printf("InfluxDB write failed: %v\n", err)
		http.Error(w, "Failed to save data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Status updated"))
}

func handleRegisterMAC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	type reqBody struct {
		MAC string `json:"mac"`
	}

	var req reqBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.MAC == "" {
		http.Error(w, "Invalid MAC address", http.StatusBadRequest)
		return
	}

	// Check if MAC already exists in InfluxDB
	const (
		token  = "my-token"
		bucket = "netsecure"
		org    = "netsecure-org"
		url    = "http://influxdb:8086"
	)

	client := influxdb2.NewClient(url, token)
	defer client.Close()

	queryAPI := client.QueryAPI(org)
	query := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: -30d)
			|> filter(fn: (r) => r._measurement == "device_status")
			|> filter(fn: (r) => r.mac == "%s")
			|> limit(n:1)
	`, bucket, req.MAC)

	result, err := queryAPI.Query(context.Background(), query)
	if err == nil && result.Next() {
		http.Error(w, "MAC already exists", http.StatusConflict)
		return
	}

	if result.Err() != nil {
		http.Error(w, "InfluxDB query failed: "+result.Err().Error(), http.StatusInternalServerError)
		return
	}

	// Write only if it doesn't exist
	if err := saveToInflux(req.MAC, "unknown"); err != nil {
		http.Error(w, "Failed to register MAC: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("MAC registered with status 'unknown'"))
}

func saveToInflux(mac, status string) error {
	const (
		token  = "my-token"
		bucket = "netsecure"
		org    = "netsecure-org"
		url    = "http://influxdb:8086"
	)

	client := influxdb2.NewClient(url, token)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(org, bucket)

	point := influxdb2.NewPoint(
		"device_status",
		map[string]string{"mac": mac},
		map[string]interface{}{"status": status},
		time.Now(),
	)

	return writeAPI.WritePoint(context.Background(), point)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Password hash failed", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`INSERT INTO users (email, password_hash) VALUES ($1, $2)`, req.Email, string(hash))
	if err != nil {
		http.Error(w, "Database insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User registered successfully")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	var hash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE email = $1", req.Email).Scan(&hash)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	fmt.Fprintln(w, "Login successful")
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func handleRouterData(w http.ResponseWriter, r *http.Request) {
	const (
		token  = "my-token"
		bucket = "netsecure"
		org    = "netsecure-org"
		url    = "http://influxdb:8086"
	)

	client := influxdb2.NewClient(url, token)
	defer client.Close()

	queryAPI := client.QueryAPI(org)
	query := `
		from(bucket: "netsecure")
			|> range(start: -7d)
			|> filter(fn: (r) => r._measurement == "device_status")
			|> filter(fn: (r) => r._field == "status")
	`

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "InfluxDB query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var output []RouterStatus

	for result.Next() {
		record := result.Record()
		entry := RouterStatus{
			MAC:   record.ValueByKey("mac").(string),
			Value: fmt.Sprintf("%v", record.Value()),
			Time:  record.Time().Format(time.RFC3339),
		}
		output = append(output, entry)
	}

	if result.Err() != nil {
		http.Error(w, "Query parsing failed: "+result.Err().Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
