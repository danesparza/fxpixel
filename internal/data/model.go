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
	R int `json:"R,omitempty"` // Red brightness level
	G int `json:"G,omitempty"` // Green brightness level
	B int `json:"B,omitempty"` // Blue brightness level
	W int `json:"W,omitempty"` // White brightness level
}

type SolidMeta struct {
	Color MetaColor `json:"color"` // Color indicates what MetaColor to display
}

type FadeMeta struct {
	Color MetaColor `json:"color"` // Color indicates what color to fade to
}

type GradientMeta struct {
	StartColor MetaColor `json:"start-color"` // StartColor indicates the first color in the gradient
	EndColor   MetaColor `json:"end-color"`   // EndColor indicates the second color in the gradient
}

type SequenceMeta struct {
	Sequence []MetaColor `json:"sequence"` // Sequence defines a repeating array of colors
}

type ZipMeta struct {
	Color MetaColor `json:"color"` // Color indicates what color to 'zip'
}

type LightningMeta struct {
	Bursts          int    `json:"bursts,omitempty"`           // Bursts indicates the number of bursts to fire in a single lightning effect
	BurstType       string `json:"burst-type"`                 // BurstType can be 'fixed' or 'random'.  Defaults to random
	BurstSpacing    int    `json:"burst-spacing,omitempty"`    // BurstSpacing indicates how much time (in ms) should exist between bursts
	BurstLength     int    `json:"burst-length,omitempty"`     // BurstLength indicates how long each flash should show
	BurstBrightness int    `json:"burst-brightness,omitempty"` // BurstBrightness indicates how bright each flash is
}

type TriggerMeta struct {
	Verb    string   `json:"verb,omitempty"`    // Verb indicates the HTTP verb to use.  Defaults to 'POST'
	URL     string   `json:"url"`               // URL indicates what url should be used
	Headers []string `json:"headers,omitempty"` // Headers indicates what HTTP headers should be passed
	Body    []byte   `json:"body,omitempty"`    // Body indicates what HTTP body should be passed.  Defaults to empty
}
