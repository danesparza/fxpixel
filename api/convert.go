package api

import (
	"encoding/json"
	"github.com/danesparza/fxpixel/internal/data"
	"github.com/danesparza/fxpixel/internal/data/const/effect"
	"github.com/danesparza/fxpixel/internal/data/const/step"
	"github.com/rs/zerolog/log"
)

// Convert internal data model to api format
func TimelineToApi(tl data.Timeline) Timeline {

	//	Convert the base timeline information
	retval := Timeline{
		ID:      tl.ID,
		Enabled: tl.Enabled,
		Created: tl.Created,
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

		//	Convert the meta info to a json string:
		metaInfo, _ := json.Marshal(item.MetaInfo)
		var jsonString string
		json.Unmarshal(metaInfo, &jsonString)

		//	... determine the step type
		switch item.Type {
		case step.Effect:
			/* If it's an effect, load effect meta */
			switch item.Effect {
			case effect.Unknown:
				//	We don't know what to do here
			case effect.Solid:
				em := SolidMeta{}
				err := json.Unmarshal([]byte(jsonString), &em)
				if err != nil {
					log.Err(err).Msg("Problem unmarshalling SolidMeta")
				}

				newStep.MetaInfo = em
			case effect.Fade:
				em := FadeMeta{}
				json.Unmarshal([]byte(jsonString), &em)
				newStep.MetaInfo = em
			case effect.Gradient:
				em := GradientMeta{}
				json.Unmarshal([]byte(jsonString), &em)
				newStep.MetaInfo = em
			case effect.Sequence:
				em := SequenceMeta{}
				json.Unmarshal([]byte(jsonString), &em)
				newStep.MetaInfo = em
			case effect.Rainbow:
				//	Don't need to do anything
			case effect.Zip:
				em := ZipMeta{}
				json.Unmarshal([]byte(jsonString), &em)
				newStep.MetaInfo = em
			case effect.KnightRider:
				//	Don't need to do anything
			case effect.Lightning:
				em := LightningMeta{}
				json.Unmarshal([]byte(jsonString), &em)
				newStep.MetaInfo = em
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
