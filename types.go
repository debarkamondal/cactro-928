package main

type FollowedArtistsResponse struct {
	Artists struct {
		Next  *string `json:"next"`
		Total int     `json:"total"`
		Items []struct {
			Name         string   `json:"name"`
			Popularity   int      `json:"popularity"`
			Genres       []string `json:"genres"`
			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Followers struct {
				Total int `json:"total"`
			} `json:"followers"`
			Images []struct {
				URL string `json:"url"`
			} `json:"images"`
		} `json:"items"`
	} `json:"artists"`
}
type TopTracksResponse struct {
	Items []struct {
		Name         string `json:"name"`
		ExternalURLs struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Album struct {
			Name   string `json:"name"`
			Images []struct {
				URL string `json:"url"`
			} `json:"images"`
			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
		} `json:"album"`
		Artists []struct {
			Name         string `json:"name"`
			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
		} `json:"artists"`
	} `json:"items"`
}
