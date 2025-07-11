package testutil

import (
	"encoding/json"
	"net/http"
	"ttv-statistics/helixclient"
)

func StubServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(helixclient.HelixUsersEndpoint, mockGetHelixUserData)
	mux.HandleFunc(helixclient.HelixVideosEndpoint, mockGetHelixVideosData)
	return mux
}

func mockGetHelixUserData(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("login")
	if userName == "" {
		http.Error(w, "no login provided", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	mockResponse := helixclient.UsersResponseBody{
		Data: []struct {
			ID              string `json:"id"`
			Login           string `json:"login"`
			DisplayName     string `json:"display_name"`
			ProfileImageURL string `json:"profile_image_url"`
		}{
			{
				ID:              "123",
				Login:           userName,
				DisplayName:     "Streamer A",
				ProfileImageURL: "https://example.com/streamerA.png",
			},
		},
	}

	_ = json.NewEncoder(w).Encode(mockResponse)
}

func mockGetHelixVideosData(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID != "123" {
		http.Error(w, "invalid or missing user_id", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	resp := helixclient.VideosResponseBody{
		Data: []helixclient.VideoInfo{
			{
				ID:        "v1",
				Title:     "Sample Video 1",
				Duration:  "30m",
				ViewCount: 150,
			},
			{
				ID:        "v2",
				Title:     "Sample Video 2",
				Duration:  "20m",
				ViewCount: 100,
			},
			{
				ID:        "v3",
				Title:     "Sample Video 3",
				Duration:  "10m",
				ViewCount: 50,
			},
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}
