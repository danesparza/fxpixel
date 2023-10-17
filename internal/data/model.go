package data

import (
	"database/sql"
	"github.com/danesparza/fxpixel/internal/data/const/effect"
	"github.com/danesparza/fxpixel/internal/data/const/step"
	"time"
)

// SystemConfig represents the system configuration information
type SystemConfig struct {
	GPIO           int    `json:"gpio"`
	LEDs           int    `json:"leds"`
	PixelOrder     string `json:"pixel_order"`
	NumberOfColors int    `json:"number_of_colors"`
}

// Timeline represents a series of event frames to be shown in order
type Timeline struct {
	ID      string         `json:"id"`             // Unique Timeline ID
	Enabled bool           `json:"enabled"`        // Timeline enabled or not
	Created time.Time      `json:"created"`        // Timeline create time
	Name    string         `json:"name"`           // Timeline name
	GPIO    sql.NullInt32  `json:"gpio,omitempty"` // The GPIO device to play the timeline on.  Optional.  If not set, uses the default
	Steps   []TimelineStep `json:"steps"`          // Steps for the timeline
	Tags    []string       `json:"tags"`           // List of Tags to associate with this timeline
}

// TimelineStep represents a single step in a timeline
type TimelineStep struct {
	ID       string            `json:"id"`                  // The timeline step id
	Type     step.StepType     `json:"type"`                // Timeline frame type (effect/sleep/trigger/loop)
	Effect   effect.EffectType `json:"effect,omitempty"`    // The Effect type (if Type=effect)
	Leds     sql.NullString    `json:"leds,omitempty"`      // Leds to use for the scene (optional) If not set and is required for the type, defaults to entire strip
	Time     sql.NullInt32     `json:"time,omitempty"`      // Time (in milliseconds).  Some things (like trigger) don't require time
	MetaInfo any               `json:"meta-info,omitempty"` // Additional information required for specific types
	Number   int               `json:"number"`              // The step number (ordinal position in the timeline)
}

type MetaColor struct {
	R int `json:"R,omitempty"`
	G int `json:"G,omitempty"`
	B int `json:"B,omitempty"`
	W int `json:"W,omitempty"`
}

type SolidMeta struct {
	Color MetaColor `json:"color"`
}

type FadeMeta struct {
	Color MetaColor `json:"color"`
}

type GradientMeta struct {
	StartColor MetaColor `json:"start-color"`
	EndColor   MetaColor `json:"end-color"`
}

type SequenceMeta struct {
	Sequence []MetaColor `json:"sequence"`
}

type ZipMeta struct {
	Color MetaColor `json:"color"`
}

type LightningMeta struct {
	Bursts          int    `json:"bursts,omitempty"`
	BurstType       string `json:"burst-type"`
	BurstSpacing    int    `json:"burst-spacing,omitempty"`
	BurstLength     int    `json:"burst-length,omitempty"`
	BurstBrightness int    `json:"burst-brightness,omitempty"`
}

type TriggerMeta struct {
	Verb    string   `json:"verb,omitempty"`
	URL     string   `json:"url"`
	Headers []string `json:"headers,omitempty"`
}
