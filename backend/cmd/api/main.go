package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"backend/internal/handlers"
	"backend/internal/repository"
)

func main() {
	// Env
	dsn := getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/netsecure?sslmode=disable")
	secret := []byte(getenv("JWT_SECRET", "devsecret"))
	addr := getenv("ADDR", ":8081")

	// DB
	repo, err := repository.NewPostgres(dsn)
	if err != nil {
		log.Fatal("db connect:", err)
	}
	defer repo.Close()

	// Handlers
	auth := handlers.NewAuthHandler(repo, secret)
	org := handlers.NewOrganizationHandler(secret)

	// Router
	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc("/api/auth/login", withCORS(auth.Login))
	mux.HandleFunc("/api/login",       withCORS(auth.Login)) // alias pour l'ancien front
	mux.HandleFunc("/api/auth/me",     withCORS(auth.Me))

	// Organization (GET/POST)
	mux.HandleFunc("/api/complete-organization", withCORS(org.CompleteOrganization))

	srv := &http.Server{
		Addr:              addr,
		Handler:           logRequests(mux),
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Println("API listening on", addr)
	log.Fatal(srv.ListenAndServe())
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// CORS minimal
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	}
}

// Logging minimal
func logRequests(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}
