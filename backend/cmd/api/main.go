package main

import (
	"log"
	"net/http"
	"os"
	"strings"
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
	mux.HandleFunc("/api/auth/login", auth.Login)
	mux.HandleFunc("/api/login", auth.Login) // alias pour l'ancien front
	mux.HandleFunc("/api/auth/me", auth.Me)

	// Organization (GET/POST)
	mux.HandleFunc("/api/complete-organization", org.CompleteOrganization)

	// Server
	srv := &http.Server{
		Addr:              addr,
		// CORS GLOBAL + logging (le CORS s'applique à toutes les routes et mêmes aux 404/500)
		Handler:           corsMiddleware(logRequests(mux)),
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

// -------- CORS global, compatible préflight --------
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		// Autorise l'origine appelante (au lieu de '*', pour éviter certains blocages navigateurs)
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		reqHdrs := r.Header.Get("Access-Control-Request-Headers")
		if strings.TrimSpace(reqHdrs) == "" {
			reqHdrs = "Authorization, Content-Type"
		}
		w.Header().Set("Access-Control-Allow-Headers", reqHdrs)

		// Préflight : renvoyer tout de suite
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logging minimal
func logRequests(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}
