package main

type FollowedArtistsResponse struct {
	Artists struct {
		Next  *string `json:"next"`
		Total int     `json:"total"`
		Items []struct {
			Name         string   `json:"name"`
			Popularity   int      `json:"popularity"`
			Genres       []string `json:"genres"`
			URI          string   `json:"uri"`
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

type CurrentlyPlaying struct {
	IsPlaying            bool   `json:"is_playing"`
	Timestamp            int64  `json:"timestamp"`
	ProgressMs           int    `json:"progress_ms"`
	CurrentlyPlayingType string `json:"currently_playing_type"`

	Context struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		Type string `json:"type"`
		Uri  string `json:"uri"`
	} `json:"context"`

	Item struct {
		DiscNumber  int     `json:"disc_number"`
		DurationMs  int     `json:"duration_ms"`
		Explicit    bool    `json:"explicit"`
		Href        string  `json:"href"`
		Id          string  `json:"id"`
		Name        string  `json:"name"`
		Popularity  int     `json:"popularity"`
		PreviewUrl  *string `json:"preview_url"`
		TrackNumber int     `json:"track_number"`
		Type        string  `json:"type"`
		Uri         string  `json:"uri"`

		ExternalIds struct {
			Isrc string `json:"isrc"`
		} `json:"external_ids"`

		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`

		Artists []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href string `json:"href"`
			Id   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
			Uri  string `json:"uri"`
		} `json:"artists"`
	}
}

type TopTracksResponse struct {
	Items []struct {
		Name         string `json:"name"`
		URI          string `json:"uri"`
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
