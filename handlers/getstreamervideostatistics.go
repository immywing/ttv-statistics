package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"ttv-statistics/constants"
	"ttv-statistics/helixclient"
	"ttv-statistics/statstools"
)

const (
	UserNamePathParam = "username"
	LastN             = "N"
)

func GetStreamerVideoStatistics(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userName := r.PathValue(UserNamePathParam)
	if userName == "" {
		http.Error(w, "missing required path param \"username\"", http.StatusNotFound)
		return
	}

	n := r.URL.Query().Get(LastN)
	if n == "" {
		http.Error(w, fmt.Sprintf("missing required URL param %q", "N"), http.StatusBadRequest)
		return
	}

	intN, err := strconv.Atoi(n)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid URL param %q %s", "N", "must be a valid integer"), http.StatusBadRequest)
		return
	}

	userData, err := helixclient.GetUserData(ctx, userName)
	if err != nil {
		http.Error(w, "error occured obtaining ttv user data", http.StatusInternalServerError)
		return
	}

	if len(userData.Data) != 1 {
		http.Error(w, "unexpected behaviour from twitch helix API.", http.StatusInternalServerError)
		return
	}

	videosData, err := helixclient.GetStreamerFirstNVideoStatistics(ctx, userData.Data[0].ID, intN)
	if err != nil {
		http.Error(w, "error occured obtaining ttv video data", http.StatusInternalServerError)
		return
	}

	aggregateData, err := statstools.AggregateStreamerVideoStatistics(videosData.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(aggregateData)
	if err != nil {
		http.Error(w, "failed to marshal response body", http.StatusInternalServerError)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJson)
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
