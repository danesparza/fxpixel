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

// ProcessGradientEffect processes the passed gradient effect meta
func (sp StepProcessor) ProcessGradientEffect(step data.TimelineStep) error {

	//	Convert the gradient meta information:
	meta := step.MetaInfo.(data.GradientMeta)

	//	For now, just log the meta information we have:
	log.Debug().
		Str("stepid", step.ID).
		Int32("time", step.Time.Int32).
		Any("startcolor", meta.StartColor).
		Any("endcolor", meta.EndColor).
		Msg("Processing effect: gradient")

	var a Artist

	a = &Gradient{
		Colors: [][]int{
			{meta.StartColor.R, meta.StartColor.G, meta.StartColor.B, meta.StartColor.W},
			{meta.EndColor.R, meta.EndColor.G, meta.EndColor.B, meta.EndColor.W},
		},
	}

	a.Draw(sp.PixArray)

	////	Write the data
	err := sp.PixArray.Write()
	if err != nil {
		log.Err(err).Msg("Problem writing to strip")
	}

	return nil
}
