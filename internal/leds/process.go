package leds

import "github.com/danesparza/fxpixel/internal/data"

type PlayTimelineRequest struct {
	ProcessID         string
	RequestedTimeline data.Timeline
}
