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
	helixclient.HelixHost = stubServer.URL
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
			userName:     "good_user",
			queryParams:  map[string]string{"N": "3"},
			expectedBody: `{"video_lengths_sum":3600000000000,"view_count_sum":300,"view_count_avg":100,"view_per_minute_avg":5,"most_viewed_video":{"title":"Sample Video 1","view_count":150}}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Missing N param",
			userName:     "good_user",
			queryParams:  map[string]string{},
			expectedBody: `message=missing required URL param innermessage=N`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid N param",
			userName:     "good_user",
			queryParams:  map[string]string{"N": "abc"},
			expectedBody: `message=invalid URL param innermessage=N must be a valid integer`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Missing username",
			userName:     "",
			queryParams:  map[string]string{"N": "3"},
			expectedBody: `message=missing required path param innermessage=username`,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "helix client fails to get user data",
			userName:     "bad_user",
			queryParams:  map[string]string{"N": "3"},
			expectedBody: fmt.Sprintf(`message=error occured obtaining ttv user data innermessage=message=received unexpected status code url=%s/users?login=bad_user status_code=400`, stubServer.URL),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "helix client fails to get user data",
			userName:     "no_data_user",
			queryParams:  map[string]string{"N": "3"},
			expectedBody: `message=no user data found`,
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "helix client fails to get user data",
			userName:     "extra_data_user",
			queryParams:  map[string]string{"N": "3"},
			expectedBody: `{"video_lengths_sum":3600000000000,"view_count_sum":300,"view_count_avg":100,"view_per_minute_avg":5,"most_viewed_video":{"title":"Sample Video 1","view_count":150}}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "helix client fails to get user data",
			userName:     "good_user_bad_video_request",
			queryParams:  map[string]string{"N": "3"},
			expectedBody: fmt.Sprintf(`message=error occured obtaining ttv video data innermessagemessage=received unexpected status code url=%s/videos?first=3&user_id=00000 status_code=400`, stubServer.URL),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

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

			bodyStr := strings.Trim(string(bodyBytes), "\n")

			if bodyStr != tc.expectedBody {
				t.Errorf("\nwant %q\n got %q", tc.expectedBody, bodyStr)
			}
		})
	}
}
