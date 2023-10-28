package leds

import (
	"github.com/Jon-Bright/ledctl/pixarray"
	"github.com/danesparza/fxpixel/internal/data"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
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

// ProcessLightningEffect processes the passed lightning effect meta
func (sp StepProcessor) ProcessLightningEffect(step data.TimelineStep) error {

	//	Convert the gradient meta information:
	meta := step.MetaInfo.(data.LightningMeta)

	//"type": "effect",
	//"effect": "lightning",
	//"meta-info": {
	//	"bursts": 3, /* Optional: For random: up to this number of bursts.  For fixed:  This many bursts */
	//	"burst-type": "fixed", /* Optional: fixed/random - defaults to random */
	//	"burst-spacing": 100, /* Optional: Maximum time (in ms) to space the bursts */
	//	"burst-length": 300, /* Optional: Maximum time (in ms) to show each burst */
	//	"burst-brightness": 128 /* Optional: Brightness of the white light for each burst */
	//}

	//	Set our defaults:
	rand.Seed(time.Now().UnixNano())
	defaultNumberOfBursts := rand.Intn(5)    //	Default number of bursts (if not specified)
	defaultBurstSpacing := rand.Intn(5)      // Default time (in ms) to space the bursts (if not specified)
	defaultBurstLength := rand.Intn(100)     // Default time (in ms) to show the bursts (if not specified)
	defaultBurstBrightness := rand.Intn(128) // Default time (in ms) to show the bursts (if not specified)

	if meta.Bursts == 0 {
		meta.Bursts = defaultNumberOfBursts
	}

	if meta.BurstSpacing == 0 {
		meta.BurstSpacing = defaultBurstSpacing
	}

	if meta.BurstLength == 0 {
		meta.BurstLength = defaultBurstLength
	}

	if meta.BurstBrightness == 0 {
		meta.BurstBrightness = defaultBurstBrightness
	}

	//	Create a lightning 'on' pixel based on the burst brightness
	ln := pixarray.Pixel{
		R: 0,
		G: 0,
		B: (meta.BurstBrightness + 1) / 2,
		W: meta.BurstBrightness,
	}

	//	Create an 'off' pixel (since we need to flash)
	loff := pixarray.Pixel{}

	//	For now, just log the meta information we have:
	log.Debug().
		Str("stepid", step.ID).
		Int32("time", step.Time.Int32).
		Any("bursts", meta.Bursts).
		Any("bursttype", meta.BurstType).
		Any("bursttype", meta.BurstType).
		Any("burstspacing", meta.BurstSpacing).
		Any("burstlength", meta.BurstLength).
		Any("burstbrightness", meta.BurstBrightness).
		Msg("Processing effect: lightning")

	//	Cycle through our bursts
	for b := 0; b < meta.Bursts; b++ {

		//	Lightning flash
		sp.PixArray.SetAll(ln)
		sp.PixArray.Write()
		time.Sleep(time.Duration(meta.BurstLength) * time.Millisecond)

		//	Flash over
		sp.PixArray.SetAll(loff)
		sp.PixArray.Write()

		//	Add burst spacing
		time.Sleep(time.Duration(meta.BurstSpacing) * time.Millisecond)
	}

	return nil
}
