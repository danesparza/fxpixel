package api

import (
	"database/sql"
	"encoding/json"
	"github.com/danesparza/fxpixel/internal/data"
	"github.com/danesparza/fxpixel/internal/data/const/effect"
	"github.com/danesparza/fxpixel/internal/data/const/step"
	"github.com/rs/zerolog/log"
	"time"
)

// TimelineToApi converts internal data model to api format
func TimelineToApi(tl data.Timeline) Timeline {

	//	Convert the base timeline information
	retval := Timeline{
		ID:      tl.ID,
		Enabled: tl.Enabled,
		Created: tl.Created.Format(time.RFC3339),
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
		case step.Trigger:
			md := item.MetaInfo.(data.TriggerMeta)
			newStep.MetaInfo = TriggerMeta{
				Verb:    md.Verb,
				URL:     md.URL,
				Headers: md.Headers,
				Body:    md.Body,
			}
		default:
		}

		//	Then add the step to the list of steps:
		retval.Steps = append(retval.Steps, newStep)
	}

	//	Return the timeline
	return retval
}

// ApiToTimeline converts api format to internal data model
func ApiToTimeline(tl Timeline) data.Timeline {

	//	Convert the base timeline information
	retval := data.Timeline{
		ID:      tl.ID,
		Enabled: tl.Enabled,
		Name:    tl.Name,
		GPIO:    sql.NullInt32{Int32: int32(tl.GPIO), Valid: true},
		Tags:    tl.Tags,
	}

	//	For each step ...
	for _, item := range tl.Steps {
		newStep := data.TimelineStep{
			ID:     item.ID,
			Type:   step.FromString(item.Type),
			Effect: effect.FromString(item.Effect),
			Leds:   sql.NullString{String: item.Leds, Valid: true},
			Time:   sql.NullInt32{Int32: int32(item.Time), Valid: true},
			Number: item.Number,
		}

		//	Convert the meta info to a json string:
		metaInfo, _ := json.Marshal(item.MetaInfo)
		jsonString := string(metaInfo)

		//	... determine the step type
		switch newStep.Type {
		case step.Effect:
			/* If it's an effect, load effect meta */
			switch newStep.Effect {
			case effect.Unknown:
				//	We don't know what to do here
			case effect.Solid:
				em := data.SolidMeta{}
				err := json.Unmarshal([]byte(jsonString), &em)
				if err != nil {
					log.Err(err).Msg("Problem unmarshalling SolidMeta")
				}

				newStep.MetaInfo = em
			case effect.Fade:
				em := data.FadeMeta{}
				json.Unmarshal([]byte(jsonString), &em)
				newStep.MetaInfo = em
			case effect.Gradient:
				em := data.GradientMeta{}
				json.Unmarshal([]byte(jsonString), &em)
				newStep.MetaInfo = em
			case effect.Sequence:
				em := data.SequenceMeta{}
				json.Unmarshal([]byte(jsonString), &em)
				newStep.MetaInfo = em
			case effect.Rainbow:
				//	Don't need to do anything
			case effect.Zip:
				em := data.ZipMeta{}
				json.Unmarshal([]byte(jsonString), &em)
				newStep.MetaInfo = em
			case effect.KnightRider:
				//	Don't need to do anything
			case effect.Lightning:
				em := data.LightningMeta{}
				json.Unmarshal([]byte(jsonString), &em)
				newStep.MetaInfo = em
			}
		case step.Sleep:
		case step.RandomSleep:
		case step.Loop:
		case step.Trigger:
			em := data.TriggerMeta{}
			json.Unmarshal([]byte(jsonString), &em)
			newStep.MetaInfo = em
		default:
		}

		//	Then add the step to the list of steps:
		retval.Steps = append(retval.Steps, newStep)
	}

	//	Return the timeline
	return retval
}
