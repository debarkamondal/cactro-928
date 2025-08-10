package main

import (
	"encoding/json"
	"net/http"
)

func fetchTopSongs(w http.ResponseWriter, accessToken string) (*TopTracksResponse, error) {

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/tracks?time_range=long_term", nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	// auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error making request", http.StatusInternalServerError)
		return nil, err
	}
	defer resp.Body.Close()

	var response TopTracksResponse
	json.NewDecoder(resp.Body).Decode(&response)
	return &response, nil
}
func fetchArtists(w http.ResponseWriter, accessToken string) (*FollowedArtistsResponse, error) {

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/following?type=artist", nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	// auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error making request", http.StatusInternalServerError)
		return nil, err
	}
	defer resp.Body.Close()
	var response FollowedArtistsResponse
	json.NewDecoder(resp.Body).Decode(&response)
	return &response, nil
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	accessToken, err := r.Cookie("access_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tracks, err := fetchTopSongs(w, accessToken.Value)
	artists, err := fetchArtists(w, accessToken.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		body := map[string]any{"message": "Couldn't fetch Artists."}
		json.NewEncoder(w).Encode(body)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	var resp struct {
		Tracks  TopTracksResponse       `json:"tracks"`
		Artists FollowedArtistsResponse `json:"artists"`
	}
	resp.Tracks = *tracks
	resp.Artists = *artists
	json.NewEncoder(w).Encode(resp)
}
