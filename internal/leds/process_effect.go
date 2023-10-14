package leds

import (
	"github.com/danesparza/fxpixel/internal/data"
	"github.com/rs/zerolog/log"
)

// ProcessSolidEffect processes the passed solid effect meta
func ProcessSolidEffect(step data.TimelineStep) error {

	//	Convert the solid meta information:
	meta := step.MetaInfo.(data.SolidMeta)

	//	For now, just log the meta information we have:
	log.Debug().
		Str("stepid", step.ID).
		Int32("time", step.Time.Int32).
		Any("color", meta.Color).
		Msg("Processing effect: solid")

	return nil
}
