package statstools

import (
	"fmt"
	"time"
	"ttv-statistics/helixclient"
)

type LastNVideoStatistics struct {
	VideoLengthsSum      time.Duration `json:"video_lengths_sum"`
	MostViewedVideoTitle string        `json:"most_viewed_video_title"`
	ViewCountSum         int           `json:"view_count_sum"`
	ViewCountAvg         int           `json:"view_count_avg"`
	ViewPerMinuteAvg     int           `json:"avg_view_per_minute"`
}

func AggregateStreamerVideoStatistics(videosData []helixclient.VideoInfo) (aggregateData LastNVideoStatistics, err error) {

	if len(videosData) == 0 {
		return LastNVideoStatistics{}, fmt.Errorf("message=%q", "no video data provided")
	}

	topVideoViewCount := 0

	for _, videoData := range videosData {

		aggregateData.ViewCountSum += videoData.ViewCount

		if topVideoViewCount < videoData.ViewCount {
			topVideoViewCount = videoData.ViewCount
			aggregateData.MostViewedVideoTitle = videoData.Title
		}

		duration, err := time.ParseDuration(videoData.Duration)
		if err != nil {
			return LastNVideoStatistics{}, fmt.Errorf("message=%q innermessage=%v", "failed to parse duration", err)
		}

		aggregateData.VideoLengthsSum += duration

	}

	aggregateData.ViewCountAvg = aggregateData.ViewCountSum / len(videosData)
	if aggregateData.VideoLengthsSum > 0 {
		aggregateData.ViewPerMinuteAvg = aggregateData.ViewCountSum / int(aggregateData.VideoLengthsSum.Minutes())
	}

	aggregateData.MostViewedVideoTitle = fmt.Sprintf("Title: %s. View Count: %d", aggregateData.MostViewedVideoTitle, topVideoViewCount)

	return aggregateData, err
}
