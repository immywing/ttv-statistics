package statstools

import (
	"fmt"
	"time"
	"ttv-statistics/helixclient"
)

type MostViewedVideo struct {
	Title     string `json:"title"`
	ViewCount int    `json:"view_count"`
}

type LastNVideoStatistics struct {
	VideoLengthsSum  time.Duration   `json:"video_lengths_sum"`
	ViewCountSum     int             `json:"view_count_sum"`
	ViewCountAvg     int             `json:"view_count_avg"`
	ViewPerMinuteAvg int             `json:"view_per_minute_avg"`
	MostViewedVideo  MostViewedVideo `json:"most_viewed_video"`
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
			aggregateData.MostViewedVideo = MostViewedVideo{
				Title:     videoData.Title,
				ViewCount: videoData.ViewCount,
			}
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

	return aggregateData, err
}
