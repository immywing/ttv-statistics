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

	var mockResponse helixclient.UsersResponseBody

	switch userName {
	case "":
		http.Error(w, "no login provided", http.StatusBadRequest)
		return
	case "good_user":
		mockResponse = helixclient.UsersResponseBody{
			Data: []struct {
				ID              string `json:"id"`
				Login           string `json:"login"`
				DisplayName     string `json:"display_name"`
				ProfileImageURL string `json:"profile_image_url"`
			}{
				{
					ID:              "good_user",
					Login:           userName,
					DisplayName:     "Streamer A",
					ProfileImageURL: "https://example.com/streamerA.png",
				},
			},
		}
	case "bad_user":
		http.Error(w, "bad user", http.StatusBadRequest)
		return
	case "no_data_user":
		// use the default var
	case "extra_data_user":
		mockResponse = helixclient.UsersResponseBody{
			Data: []struct {
				ID              string `json:"id"`
				Login           string `json:"login"`
				DisplayName     string `json:"display_name"`
				ProfileImageURL string `json:"profile_image_url"`
			}{
				{
					ID:              "good_user",
					Login:           userName,
					DisplayName:     "Streamer A",
					ProfileImageURL: "https://example.com/streamerA.png",
				},
				{
					ID:              "good_user",
					Login:           userName,
					DisplayName:     "Streamer A",
					ProfileImageURL: "https://example.com/streamerA.png",
				},
			},
		}
	case "good_user_bad_video_request":
		mockResponse = helixclient.UsersResponseBody{
			Data: []struct {
				ID              string `json:"id"`
				Login           string `json:"login"`
				DisplayName     string `json:"display_name"`
				ProfileImageURL string `json:"profile_image_url"`
			}{
				{
					ID:              "00000",
					Login:           userName,
					DisplayName:     "Streamer A",
					ProfileImageURL: "https://example.com/streamerA.png",
				},
			},
		}
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(mockResponse)
}

func mockGetHelixVideosData(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID != "good_user" {
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
