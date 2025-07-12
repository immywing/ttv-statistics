package statstools_test

import (
	"testing"
	"time"
	"ttv-statistics/helixclient"
	"ttv-statistics/statstools"
)

func TestAggregateStreamerVideoStatistics(t *testing.T) {

	type testCase struct {
		name          string
		inputs        []helixclient.VideoInfo
		expected      statstools.LastNVideoStatistics
		expectedError string
	}

	testCases := []testCase{
		{
			name:          "No video data",
			inputs:        []helixclient.VideoInfo{},
			expected:      statstools.LastNVideoStatistics{},
			expectedError: `message="no video data provided"`,
		},
		{
			name: "Valid video data",
			inputs: []helixclient.VideoInfo{
				{
					Duration:  "1h30m",
					ViewCount: 1000,
					Title:     "First Video",
				},
				{
					Duration:  "2h15m",
					ViewCount: 2000,
					Title:     "Second Video",
				},
				{
					Duration:  "45m",
					ViewCount: 1500,
					Title:     "Third Video",
				},
			},
			expected: statstools.LastNVideoStatistics{
				VideoLengthsSum:  4*time.Hour + 30*time.Minute,
				ViewCountSum:     4500,
				MostViewedVideo:  statstools.MostViewedVideo{Title: "Second Video", ViewCount: 2000},
				ViewCountAvg:     1500,
				ViewPerMinuteAvg: 16,
			},
		},
		{
			name:          "Valid video data with a duration of 0",
			expectedError: "",
			inputs: []helixclient.VideoInfo{
				{
					Duration:  "0m",
					ViewCount: 1000,
					Title:     "First Video",
				},
				{
					Duration:  "0m",
					ViewCount: 2000,
					Title:     "Second Video",
				},
				{
					Duration:  "0m",
					ViewCount: 1500,
					Title:     "Third Video",
				},
			},
			expected: statstools.LastNVideoStatistics{
				VideoLengthsSum:  0 * time.Minute,
				ViewCountSum:     4500,
				MostViewedVideo:  statstools.MostViewedVideo{Title: "Second Video", ViewCount: 2000},
				ViewCountAvg:     1500,
				ViewPerMinuteAvg: 0,
			},
		},
		{
			name:          "Valid video data with a bad duration format",
			expectedError: `message="failed to parse duration" innermessage=time: invalid duration "invalid"`,
			inputs: []helixclient.VideoInfo{
				{
					Duration:  "invalid",
					ViewCount: 1000,
					Title:     "First Video",
				},
			},
			expected: statstools.LastNVideoStatistics{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			result, err := statstools.AggregateStreamerVideoStatistics(tc.inputs)

			if err != nil && err.Error() != tc.expectedError {
				t.Errorf("unexpected error: want: %v, got: %v", tc.expectedError, err)
			}

			if result != tc.expected {
				t.Errorf("unexpected result: \nwant: %+v, \n got: %+v", tc.expected, result)
			}
		})
	}
}
