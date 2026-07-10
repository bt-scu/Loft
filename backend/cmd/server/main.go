package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bt-scu/Loft/backend/internal/spotify"
	"github.com/joho/godotenv"
)

func withCORS(next http.Handler) http.Handler {
	allowedOrigin := os.Getenv("FRONTEND_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Function that runs when you call root url
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprint(w, "Hello, World! Your Go server is working.")
}

func main() {
	// Load SPOTIFY_CLIENT_ID / SPOTIFY_CLIENT_SECRET etc. from .env, if present
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on process environment")
	}

	// Create a router (if you see some extra path route it to some function)
	mux := http.NewServeMux()

	//Create spotify handlers
	spotifyHandlers := spotify.NewHandlers()

	// Register routing handlers
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("GET /spotify/auth", spotifyHandlers.SpotifyAuthHandler)

	//not called by user (function initialized by spotify callback)
	mux.HandleFunc("GET /auth/callback", spotifyHandlers.CompleteAuth)

	// Define the network port
	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)

	// Start the server and log errors if it fails to start
	err := http.ListenAndServe(port, withCORS(mux))
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
