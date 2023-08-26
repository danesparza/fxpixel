package data

import "time"

// Timeline represents a series of event frames to be shown in order
type Timeline struct {
	ID      string          `json:"id"`             // Unique Timeline ID
	Enabled bool            `json:"enabled"`        // Timeline enabled or not
	Created time.Time       `json:"created"`        // Timeline create time
	Name    string          `json:"name"`           // Timeline name
	GPIO    int             `json:"gpio,omitempty"` // The GPIO device to play the timeline on.  Optional.  If not set, uses the default
	Frames  []TimelineFrame `json:"frames"`         // Frames for the timeline
}

// TimelineFrame represents a single event frame in a timeline
type TimelineFrame struct {
	Type      string         `json:"type"`               // Timeline frame type (scene/sleep/fade) Fade 'fades' between the previous channel state and this frame
	Channels  []ChannelValue `json:"channels,omitempty"` // Channel information to set for the scene (optional) Required if type = scene or fade
	SleepTime int            `json:"sleeptime"`          // Sleep type in seconds (optional) Required if type = sleep
}

// ChannelValue needs to be renamed and reworked to fit into LED strips
type ChannelValue struct {
	Channel int  `json:"channel"` // Unique Fixture ID
	Value   byte `json:"value"`   // Optional fixture name
}
