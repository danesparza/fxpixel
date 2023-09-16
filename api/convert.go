package api

import (
	"github.com/danesparza/fxpixel/internal/data"
	"github.com/danesparza/fxpixel/internal/data/const/effect"
	"github.com/danesparza/fxpixel/internal/data/const/step"
	"time"
)

// Convert internal data model to api format
func TimelineToApi(tl data.Timeline) Timeline {

	unixTimeUTC := time.Unix(tl.Created, 0) //gives unix time stamp in utc

	//	Convert the base timeline information
	retval := Timeline{
		ID:      tl.ID,
		Enabled: tl.Enabled,
		Created: unixTimeUTC.Format(time.RFC3339),
		Name:    tl.Name,
		GPIO:    int(tl.GPIO.Int32),
		Tags:    tl.Tags,
	}

	//	For each step ...
	for _, item := range tl.Steps {
		newStep := TimelineStep{
			ID:     item.ID,
			Type:   item.Type.String(),
			Effect: item.Effect.String(),
			Leds:   item.Leds.String,
			Time:   int(item.Time.Int32),
			Number: item.Number,
		}

		//	... determine the step type
		switch item.Type {
		case step.Effect:
			/* If it's an effect, load effect meta */
			switch item.Effect {
			case effect.Unknown:
				//	We don't know what to do here
			case effect.Solid:
				md := item.MetaInfo.(data.SolidMeta)
				newStep.MetaInfo = SolidMeta{
					Color: MetaColor{
						R: md.Color.R,
						G: md.Color.G,
						B: md.Color.B,
						W: md.Color.W,
					},
				}
			case effect.Fade:
				md := item.MetaInfo.(data.FadeMeta)
				newStep.MetaInfo = FadeMeta{
					Color: MetaColor{
						R: md.Color.R,
						G: md.Color.G,
						B: md.Color.B,
						W: md.Color.W,
					},
				}
			case effect.Gradient:
				md := item.MetaInfo.(data.GradientMeta)
				newStep.MetaInfo = GradientMeta{
					StartColor: MetaColor{
						R: md.StartColor.R,
						G: md.StartColor.G,
						B: md.StartColor.B,
						W: md.StartColor.W,
					},
					EndColor: MetaColor{
						R: md.EndColor.R,
						G: md.EndColor.G,
						B: md.EndColor.B,
						W: md.EndColor.W,
					},
				}
			case effect.Sequence:
				md := item.MetaInfo.(data.SequenceMeta)
				//	Copy the sequence
				sequenceSlice := []MetaColor{}
				for _, item := range md.Sequence {
					sequenceItem := MetaColor{
						R: item.R,
						G: item.G,
						B: item.B,
						W: item.W,
					}
					sequenceSlice = append(sequenceSlice, sequenceItem)
				}
				newStep.MetaInfo = SequenceMeta{Sequence: sequenceSlice}
			case effect.Rainbow:
				//	Don't need to do anything
			case effect.Zip:
				md := item.MetaInfo.(data.ZipMeta)
				newStep.MetaInfo = ZipMeta{
					Color: MetaColor{
						R: md.Color.R,
						G: md.Color.G,
						B: md.Color.B,
						W: md.Color.W,
					},
				}
			case effect.KnightRider:
				//	Don't need to do anything
			case effect.Lightning:
				md := item.MetaInfo.(data.LightningMeta)
				newStep.MetaInfo = LightningMeta{
					Bursts:          md.Bursts,
					BurstType:       md.BurstType,
					BurstSpacing:    md.BurstSpacing,
					BurstLength:     md.BurstLength,
					BurstBrightness: md.BurstBrightness,
				}
			}
		case step.Sleep:
		case step.RandomSleep:
		case step.Loop:
		default:
		}

		//	Then add the step to the list of steps:
		retval.Steps = append(retval.Steps, newStep)
	}

	//	Return the timeline
	return retval
}
