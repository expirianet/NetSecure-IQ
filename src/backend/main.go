package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

var db *sql.DB

// Separate struct for registration
type RegisterRequest struct {
	Email string `json:"email"`
}

// Separate struct for login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	var err error
	connStr := "host=postgresql user=netsecure_iq password=your_secure_password dbname=netsecure_iq_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// ✅ Test connection to the database
	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection test failed:", err)
	}
	fmt.Println("✅ Connected to PostgreSQL successfully")

	http.HandleFunc("/api/register", withCORS(handleRegister))
	http.HandleFunc("/api/login", withCORS(handleLogin))

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

	// Save user to DB with default role
	//defaultRole := "tenant"
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
		// Optional: return error to client if needed
		// http.Error(w, "User created, but failed to send email", http.StatusInternalServerError)
		// return
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

	// respond with login result
	resp := map[string]interface{}{
		"message": "Login successful",
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
