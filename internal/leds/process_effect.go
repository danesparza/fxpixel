package leds

import (
	"github.com/Jon-Bright/ledctl/pixarray"
	"github.com/danesparza/fxpixel/internal/data"
	"github.com/rs/zerolog/log"
)

// ProcessSolidEffect processes the passed solid effect meta
func (sp StepProcessor) ProcessSolidEffect(step data.TimelineStep) error {

	//	Convert the solid meta information:
	meta := step.MetaInfo.(data.SolidMeta)

	//	For now, just log the meta information we have:
	log.Debug().
		Str("stepid", step.ID).
		Int32("time", step.Time.Int32).
		Any("color", meta.Color).
		Msg("Processing effect: solid")

	//	Create an individual pixel to set the color
	p := pixarray.Pixel{
		R: meta.Color.R,
		G: meta.Color.G,
		B: meta.Color.B,
		W: meta.Color.W,
	}

	//	Set all pixels in the array to this color
	sp.PixArray.SetAll(p)

	//	Write the data
	err := sp.PixArray.Write()
	if err != nil {
		log.Err(err).Msg("Problem writing to strip")
	}

	return nil
}
