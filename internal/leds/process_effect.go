package leds

import (
	"context"
	"github.com/Jon-Bright/ledctl/effects"
	"github.com/Jon-Bright/ledctl/pixarray"
	"github.com/danesparza/fxpixel/internal/data"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
)

// ProcessSolidEffect processes the passed solid effect meta
func (sp StepProcessor) ProcessSolidEffect(step data.TimelineStep) error {

	//	Convert the meta information:
	meta := step.MetaInfo.(data.SolidMeta)

	//	Log the meta information we have:
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

	//	Convert the meta information:
	meta := step.MetaInfo.(data.GradientMeta)

	//	Log the meta information we have:
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
func (sp StepProcessor) ProcessLightningEffect(ctx context.Context, step data.TimelineStep) error {

	//	Convert the meta information:
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

	//	Log the meta information we have:
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

		select {
		default:

			//	Lightning flash
			sp.PixArray.SetAll(ln)
			sp.PixArray.Write()

			select {
			case <-time.After(time.Duration(meta.BurstLength) * time.Millisecond):
				//	Flash over
				sp.PixArray.SetAll(loff)
				sp.PixArray.Write()

				//	Add burst spacing
				select {
				case <-time.After(time.Duration(meta.BurstSpacing) * time.Millisecond):
					continue
				case <-ctx.Done():
					return nil
				}

			case <-ctx.Done():
				return nil
			}

		case <-ctx.Done():
			//	Reset all pixels:
			sp.PixArray.SetAll(pixarray.Pixel{})
			sp.PixArray.Write()

			return nil
		}

	}

	return nil
}

// ProcessSequenceEffect processes the passed sequence effect meta
func (sp StepProcessor) ProcessSequenceEffect(step data.TimelineStep) error {

	//	Convert the meta information:
	meta := step.MetaInfo.(data.SequenceMeta)

	//	Log the meta information we have:
	log.Debug().
		Str("stepid", step.ID).
		Int32("time", step.Time.Int32).
		Any("sequence", meta.Sequence).
		Msg("Processing effect: sequence")

	var a Artist

	//	Build the color sequence slice
	colorSeq := [][]int{}
	for _, color := range meta.Sequence {
		colorSeq = append(colorSeq, []int{
			color.R, color.G, color.B, color.W,
		})
	}

	//	Create the sequence
	a = &Sequence{
		Colors: colorSeq,
	}

	//	Draw
	a.Draw(sp.PixArray)

	////	Write the data
	err := sp.PixArray.Write()
	if err != nil {
		log.Err(err).Msg("Problem writing to strip")
	}

	return nil
}

// ProcessFadeEffect processes the passed fade effect meta
func (sp StepProcessor) ProcessFadeEffect(ctx context.Context, step data.TimelineStep) error {

	//	Convert the meta information:
	meta := step.MetaInfo.(data.FadeMeta)

	//	Log the meta information we have:
	log.Debug().
		Str("stepid", step.ID).
		Int32("time", step.Time.Int32).
		Any("color", meta.Color).
		Msg("Processing effect: fade")

	var d time.Duration

	fade := effects.NewFade(time.Duration(step.Time.Int32)*time.Millisecond, pixarray.Pixel{
		R: meta.Color.R,
		G: meta.Color.G,
		B: meta.Color.B,
		W: meta.Color.W,
	})

	fade.Start(sp.PixArray, time.Now())

	//	Create a ticker to process work:
	ticker := time.NewTicker(1 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			d = fade.NextStep(sp.PixArray, time.Now())
			err := sp.PixArray.Write()
			if err != nil {
				log.Err(err).Msg("Problem writing to strip")
			}

			//	This is a weird way to signal this,
			//	but a duration of 0 means the fade is 'done'
			if d == 0 {
				break
			}
		case <-ctx.Done():
			//	Reset all pixels:
			sp.PixArray.SetAll(pixarray.Pixel{})
			sp.PixArray.Write()

			return nil
		}
	}

	return nil
}

// ProcessKnightRiderEffect processes the knight rider effect
func (sp StepProcessor) ProcessKnightRiderEffect(ctx context.Context, step data.TimelineStep) error {

	//	Log the meta information we have:
	log.Debug().
		Str("stepid", step.ID).
		Int32("time", step.Time.Int32).
		Msg("Processing effect: knightrider")

	var d time.Duration

	kr := effects.NewKnightRider(1*time.Second, 5)
	kr.Start(sp.PixArray, time.Now())

	//	Create a ticker to process work:
	ticker := time.NewTicker(1 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			d = kr.NextStep(sp.PixArray, time.Now())
			err := sp.PixArray.Write()
			if err != nil {
				log.Err(err).Msg("Problem writing to strip")
			}

			//	This is a weird way to signal this,
			//	but a duration of 0 means the fade is 'done'
			if d == 0 {
				break
			}
		case <-ctx.Done():
			//	Reset all pixels:
			sp.PixArray.SetAll(pixarray.Pixel{})
			sp.PixArray.Write()

			return nil
		}
	}
}

// ProcessRainbowEffect processes the rainbow effect
func (sp StepProcessor) ProcessRainbowEffect(ctx context.Context, step data.TimelineStep) error {

	//	Log the meta information we have:
	log.Debug().
		Str("stepid", step.ID).
		Int32("time", step.Time.Int32).
		Msg("Processing effect: rainbow")

	var d time.Duration

	rainbow := effects.NewRainbow(20 * time.Second)
	rainbow.Start(sp.PixArray, time.Now())

	//	Create a ticker to process work:
	ticker := time.NewTicker(1 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			d = rainbow.NextStep(sp.PixArray, time.Now())
			sp.PixArray.Write()

			//	This is a weird way to signal this,
			//	but a duration of 0 means the fade is 'done'
			if d == 0 {
				break
			}

		case <-ctx.Done():
			//	Reset all pixels:
			sp.PixArray.SetAll(pixarray.Pixel{})
			sp.PixArray.Write()

			return nil
		}
	}
}

// ProcessZipEffect processes the rainbow effect
func (sp StepProcessor) ProcessZipEffect(ctx context.Context, step data.TimelineStep) error {

	//	Convert the meta information:
	meta := step.MetaInfo.(data.ZipMeta)

	//	Log the meta information we have:
	log.Debug().
		Str("stepid", step.ID).
		Int32("time", step.Time.Int32).
		Any("color", meta.Color).
		Msg("Processing effect: zip")

	var d time.Duration

	//	Use the time from the step, but default to 2 seconds if it's not set
	zipDuration := int(step.Time.Int32)
	if zipDuration == 0 {
		zipDuration = 2000
	}

	zip := effects.NewZip(time.Duration(zipDuration)*time.Millisecond, pixarray.Pixel{
		R: meta.Color.R,
		G: meta.Color.G,
		B: meta.Color.B,
		W: meta.Color.W,
	})

	zip.Start(sp.PixArray, time.Now())

	//	Create a ticker to process work:
	ticker := time.NewTicker(1 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			d = zip.NextStep(sp.PixArray, time.Now())
			sp.PixArray.Write()

			//	This is a weird way to signal this,
			//	but a duration of 0 means the fade is 'done'
			if d == 0 {
				break
			}

		case <-ctx.Done():
			//	Reset all pixels:
			sp.PixArray.SetAll(pixarray.Pixel{})
			sp.PixArray.Write()

			return nil
		}
	}
}
