package handlers

import (
	"encoding/json"
	"fmt"
	"log"
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
		http.Error(w, fmt.Sprintf("message=%s innermessage=%s", "missing required path param", "username"), http.StatusNotFound)
		return
	}

	n := r.URL.Query().Get(LastN)
	if n == "" {
		http.Error(w, fmt.Sprintf("message=%s innermessage=%s", "missing required URL param", "N"), http.StatusBadRequest)
		return
	}

	intN, err := strconv.Atoi(n)
	if err != nil {
		http.Error(w, fmt.Sprintf("message=%s innermessage=%s", "invalid URL param", "N must be a valid integer"), http.StatusBadRequest)
		return
	}

	userData, err := helixclient.GetUserData(ctx, userName)
	if err != nil {
		http.Error(w, fmt.Sprintf("message=%s innermessage=%v", "error occured obtaining ttv user data", err), http.StatusInternalServerError)
		return
	}

	if len(userData.Data) == 0 {
		http.Error(w, fmt.Sprintf("message=%s", "no user data found"), http.StatusNoContent)
		return
	}

	if len(userData.Data) > 1 {
		log.Printf("Warning: Helix API returned more than 1 result in User Data array")
	}

	videosData, err := helixclient.GetStreamerFirstNVideoStatistics(ctx, userData.Data[0].ID, intN)
	if err != nil {
		http.Error(w, fmt.Sprintf("message=%s innermessage%v", "error occured obtaining ttv video data", err), http.StatusInternalServerError)
		return
	}

	aggregateData, err := statstools.AggregateStreamerVideoStatistics(videosData.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(aggregateData)
	if err != nil {
		http.Error(w, fmt.Sprintf("message=%s innermessage=%v", "failed to marshal response body", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set(constants.ContentTypeHeaderKey, constants.ContentTypeApplicationJson)
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
