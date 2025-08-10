package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func PauseHandler(w http.ResponseWriter, r *http.Request) {

	accessToken, err := r.Cookie("access_token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Make request
	req, err := http.NewRequest("PUT", "https://api.spotify.com/v1/me/player/pause", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken.Value)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}
func PlayHandler(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	contextUri := r.URL.Query().Get("context_uri")

	accessToken, err := r.Cookie("access_token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Make request
	var body struct {
		Uris       []string `json:"uris,omitempty"`
		ContextUri string   `json:"context_uri,omitempty"`
	}
	body.Uris = append(body.Uris, uri)
	body.ContextUri = contextUri
	marshalleldBody, err := json.Marshal(body)
	bodyReader := bytes.NewBuffer(marshalleldBody)
	req, err := http.NewRequest("PUT", "https://api.spotify.com/v1/me/player/play", bodyReader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken.Value)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}
