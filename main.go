package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

)

var clientID = os.Getenv("CLIENT_ID")
var clientSecret = os.Getenv("CLIENT_SECRET")
var redirectURI = "https://prod.dkmondal.in/spotify/auth/callback"
var stateKey = "spotify_auth_state"

// TokenResponse represents the response from Spotify's token endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// generateRandomString creates a random string of specified length
func generateRandomString(length int) string {
	bytes := make([]byte, 30) // Generate more bytes than needed
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// loginHandler handles the /login endpoint
func loginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateRandomString(16)

	// Set state cookie
	cookie := &http.Cookie{
		Name:  stateKey,
		Value: state,
		Path:  "/",
	}
	http.SetCookie(w, cookie)

	// Build authorization URL
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", clientID)
	params.Add("scope", "user-read-private user-modify-playback-state user-read-currently-playing user-top-read user-follow-read user-read-email")
	params.Add("redirect_uri", redirectURI)
	params.Add("state", state)

	authURL := "https://accounts.spotify.com/authorize?" + params.Encode()
	http.Redirect(w, r, authURL, http.StatusFound)
}

// callbackHandler handles the /callback endpoint
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	// Get stored state from cookie
	cookie, err := r.Cookie(stateKey)
	var storedState string
	if err == nil {
		storedState = cookie.Value
	}

	// Verify state parameter
	if state == "" || state != storedState {
		params := url.Values{}
		params.Add("error", "state_mismatch")
		http.Redirect(w, r, "/#"+params.Encode(), http.StatusFound)
		return
	}

	// Clear state cookie
	clearCookie := &http.Cookie{
		Name:   stateKey,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, clearCookie)

	// Exchange code for tokens
	tokenData := url.Values{}
	tokenData.Set("code", code)
	tokenData.Set("redirect_uri", redirectURI)
	tokenData.Set("grant_type", "authorization_code")

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(tokenData.Encode()))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error making request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var tokenResp TokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
			http.Error(w, "Error decoding response", http.StatusInternalServerError)
			return
		}

		// Use the access token to get user info
		userReq, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
		if err == nil {
			userReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
			userResp, err := client.Do(userReq)
			if err == nil {
				defer userResp.Body.Close()
				var userInfo map[string]any
				json.NewDecoder(userResp.Body).Decode(&userInfo)
				fmt.Printf("User info: %+v\n", userInfo)
			}
		}

		// Redirect with tokens
		http.SetCookie(w, &http.Cookie{
			Name:  "access_token",
			Value: tokenResp.AccessToken,
			// Domain: "prod.dkmondal.in",
			Path:     "/spotify",
			MaxAge:   3600,
			Secure:   true,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/spotify/dashboard", http.StatusFound)
	} else {
		params := url.Values{}
		params.Add("error", "invalid_token")
		http.Redirect(w, r, "/#"+params.Encode(), http.StatusFound)
	}
}

// refreshTokenHandler handles the /refresh_token endpoint
func refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.URL.Query().Get("refresh_token")
	if refreshToken == "" {
		http.Error(w, "Missing refresh_token parameter", http.StatusBadRequest)
		return
	}

	// Prepare form data
	tokenData := url.Values{}
	tokenData.Set("grant_type", "refresh_token")
	tokenData.Set("refresh_token", refreshToken)

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(tokenData.Encode()))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error making request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var tokenResp TokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
			http.Error(w, "Error decoding response", http.StatusInternalServerError)
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"access_token":  tokenResp.AccessToken,
			"refresh_token": tokenResp.RefreshToken,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Error refreshing token", http.StatusBadRequest)
	}
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/public/index.html")
}

func main() {
	mux := http.NewServeMux()
	root := http.NewServeMux()

	root.Handle("/spotify", http.RedirectHandler("/spotify/", http.StatusMovedPermanently))
	root.Handle("/spotify/", http.StripPrefix("/spotify", mux))

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/dashboard", dashboardHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/auth/callback", callbackHandler)
	mux.HandleFunc("/auth/refresh_token", refreshTokenHandler)

	mux.HandleFunc("/play", PlayHandler)
	mux.HandleFunc("/pause", PauseHandler)

	fmt.Println("Listening on port 8081 (proxied)")
	if err := http.ListenAndServe(":8081", root); err != nil {
		fmt.Println(err)
		fmt.Println("Couldn't initiate server on port 8081")
	}

}
