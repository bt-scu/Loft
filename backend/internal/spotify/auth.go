package spotify

import (
	"log"
	"net/http"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

// redirects to the Spotify Developer Dashboard (alongside authcoe and state. THese values are nested inside of the callbackurl [auth/callback])
const redirectURI = "http://127.0.0.1:8080/auth/callback"

// state should be randomly generated per login attempt and tied to the session that started it.
// An attacker's code always comes paired with the attacker's own state (from their own login).
// A victim's session will have a different (or no) expected state stored for that state value,
// so even if a victim's browser is tricked into requesting the callback with the attacker's
// code + state, the server's state-comparison check fails and the callback is rejected —
// preventing the victim's session from getting linked to the attacker's account.
const state = "abc123"

// Handlers holds the Spotify authenticator so the HTTP handlers below can be
// registered as methods instead of relying on a package-level global.
type Handlers struct {
	auth           *spotifyauth.Authenticator
	frontendOrigin string
}

// Handlers are a collection of dependencies(external things to make code works, usually has behaviors). We use pointer syntax to ensure we are using the same handler instance
func NewHandlers() *Handlers {
	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:5173"
	}

	return &Handlers{
		auth: spotifyauth.New(
			spotifyauth.WithRedirectURL(redirectURI),
			spotifyauth.WithScopes(
				spotifyauth.ScopeUserReadCurrentlyPlaying,
				spotifyauth.ScopeUserReadPlaybackState,
				spotifyauth.ScopeUserModifyPlaybackState,
			),
			spotifyauth.WithClientID(os.Getenv("SPOTIFY_CLIENT_ID")),
			spotifyauth.WithClientSecret(os.Getenv("SPOTIFY_CLIENT_SECRET")),
		),
		frontendOrigin: frontendOrigin,
	}
}

// SpotifyAuthHandler sends the user to Spotify's login/consent page.
// w will set  the location of the outgoing response to the redirect URL and write the status code
// r really just being checked to see if the URL is a relative URL (path only) or for the r.Method (to offer fallback logic if user does not access from browser)
func (h *Handlers) SpotifyAuthHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, h.auth.AuthURL(state), http.StatusFound)
}

// CompleteAuth handles Spotify's redirect back after the user logs in:
// it exchanges the code for a token and confirms the token works.
func (h *Handlers) CompleteAuth(w http.ResponseWriter, r *http.Request) {

	//perform token exchange (passes context(metadata on client state, if browser dies than request is cancelled), state, r(where the access code is stored))
	tok, err := h.auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "couldn't get token", http.StatusForbidden)
		log.Println("token exchange failed:", err)
		return
	}

	//takes (token + Context) -> wraps http client around it (which has the bearer token in the header) -> wraps this client and wraps it into a go object with nice go methods
	client := spotify.New(h.auth.Client(r.Context(), tok))

	//Get request to Spotify's v1/m endpoint, sees the BEarer token and returns the user's profile. Fetches for user profiele
	user, err := client.CurrentUser(r.Context())

	if err != nil {
		http.Error(w, "couldn't fetch user profile", http.StatusInternalServerError)
		log.Println("fetching current user failed:", err)
		return
	}

	//affirm access
	log.Printf("logged in as %s (access token: %s)\n", user.ID, tok.AccessToken)

	// No session/token persistence yet (that's the DB layer, still to come) —
	// for now just send the browser back to the frontend so the login loop
	// is click-through-able end to end.
	http.Redirect(w, r, h.frontendOrigin, http.StatusFound)
}
