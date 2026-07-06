package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bt-scu/Loft/backend/internal/spotify"
	"github.com/joho/godotenv"
)

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
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
