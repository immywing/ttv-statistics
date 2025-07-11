package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"ttv-statistics/handlers"
	"ttv-statistics/helixclient"
	"ttv-statistics/testutil"
)

func TestGetStreamerVideoStatistics(t *testing.T) {

	stubServer := httptest.NewServer(testutil.StubServerMux())
	defer stubServer.Close()

	helixclient.HelixHost = stubServer.URL
	helixclient.ClientID = "stub-client-id"

	type testCase struct {
		name         string
		userName     string
		queryParams  map[string]string
		expectedBody string
		expectedCode int
	}

	testCases := []testCase{
		{
			name:         "Valid request",
			userName:     "123",
			queryParams:  map[string]string{"N": "3"},
			expectedBody: `{"video_lengths_sum":3600000000000,"most_viewed_video_title":"Title: Sample Video 1. View Count: 150","view_count_sum":300,"view_count_avg":100,"avg_view_per_minute":5}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Missing N param",
			userName:     "123",
			queryParams:  map[string]string{},
			expectedBody: `missing required URL param "N"`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid N param",
			userName:     "123",
			queryParams:  map[string]string{"N": "abc"},
			expectedBody: `invalid URL param "N" must be a valid integer`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Missing username",
			userName:     "",
			queryParams:  map[string]string{"N": "3"},
			expectedBody: "missing required path param \"username\"", // this will depend on your router
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Build URL
			urlPath := fmt.Sprintf("/streamer/%s/statistics", tc.userName)
			query := url.Values{}
			for k, v := range tc.queryParams {
				query.Set(k, v)
			}

			req := httptest.NewRequest(http.MethodGet, urlPath+"?"+query.Encode(), nil)
			req.SetPathValue(handlers.UserNamePathParam, tc.userName)

			rec := httptest.NewRecorder()
			handlers.GetStreamerVideoStatistics(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			bodyBytes, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != tc.expectedCode {
				t.Errorf("expected status %d, got %d", tc.expectedCode, resp.StatusCode)
			}

			bodyStr := strings.TrimSpace(string(bodyBytes))
			if !strings.Contains(bodyStr, tc.expectedBody) {
				t.Errorf("want%q\n got %q", tc.expectedBody, bodyStr)
			}
		})
	}
}
