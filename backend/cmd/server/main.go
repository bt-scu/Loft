package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bt-scu/Loft/backend/internal/db"
	"github.com/bt-scu/Loft/backend/internal/spotify"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

//A middleware step that runs on every request
func withCORS(next http.Handler) http.Handler {
	allowedOrigin := os.Getenv("FRONTEND_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173"
	}
	
	//NOTE: W and R are go provided objects that you can use to configure your own HTTP objects
	//w: you can write headers, status codes, and body bytes
	//r: a struct that describes everything about the incoming request (method. path, headers, body, etc.)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//frontend is allowed to read responses
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)

		//frontend is allowed to send cookies via these HTTP methods and a specific Content-Type Header
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		//Handles preflight logic (Frontend throws an initial check and backend responds with a Yes with no body)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		//Calls next handler
		next.ServeHTTP(w, r)
	})
}

//Function that runs when you call root url
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
	router := chi.NewRouter()
	router.Use(withCORS)

	//Create spotify handlers
	spotifyHandlers := spotify.InitHandlers()

	// Register routing handlers
	router.Get("/", homeHandler)
	router.Get("/spotify/auth", spotifyHandlers.SpotifyAuthHandler)

	//not called by user (function initialized by spotify callback)
	router.Get("/auth/callback", spotifyHandlers.CompleteAuth)

	// Define the network port
	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)

	pool, err := db.Connect(context.Background())
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer pool.Close()

	// Start the server and log errors if it fails to start
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
