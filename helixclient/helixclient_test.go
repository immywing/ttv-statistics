package helixclient_test

import (
	"context"
	"net/http/httptest"
	"testing"
	"ttv-statistics/helixclient"
	"ttv-statistics/testutil"
)

func TestGetUserData(t *testing.T) {
	server := httptest.NewServer(testutil.StubServerMux())
	defer server.Close()
	helixclient.HelixHost = server.URL

	type testCase struct {
		name            string
		userName        string
		expectError     bool
		expectedDataLen int
	}

	testCases := []testCase{
		{
			name:            "Good user returns one entry",
			userName:        "good_user",
			expectError:     false,
			expectedDataLen: 1,
		},
		{
			name:            "Bad user returns error",
			userName:        "bad_user",
			expectError:     true,
			expectedDataLen: 0,
		},
		{
			name:            "No data user returns empty data",
			userName:        "no_data_user",
			expectError:     false,
			expectedDataLen: 0,
		},
		{
			name:            "Extra data user returns multiple entries",
			userName:        "extra_data_user",
			expectError:     false,
			expectedDataLen: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := helixclient.GetUserData(context.Background(), tc.userName)
			if tc.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.expectError && len(resp.Data) != tc.expectedDataLen {
				t.Errorf("expected %d user data entries, got %d", tc.expectedDataLen, len(resp.Data))
			}
		})
	}
}

func TestGetStreamerFirstNVideoStatistics(t *testing.T) {
	server := httptest.NewServer(testutil.StubServerMux())
	defer server.Close()
	helixclient.HelixHost = server.URL

	type testCase struct {
		name        string
		userID      string
		n           int
		expectError bool
		expectedLen int
	}

	testCases := []testCase{
		{
			name:        "Valid userID returns videos",
			userID:      "good_user",
			n:           3,
			expectError: false,
			expectedLen: 3,
		},
		{
			name:        "Invalid userID returns error",
			userID:      "invalid_user",
			n:           3,
			expectError: true,
			expectedLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := helixclient.GetStreamerFirstNVideoStatistics(context.Background(), tc.userID, tc.n)
			if tc.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tc.expectError && len(resp.Data) != tc.expectedLen {
				t.Errorf("expected %d video data entries, got %d", tc.expectedLen, len(resp.Data))
			}
		})
	}
}
