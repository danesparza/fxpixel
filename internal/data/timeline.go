package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/danesparza/fxpixel/internal/data/const/effect"
	"github.com/danesparza/fxpixel/internal/data/const/step"
	"github.com/rs/zerolog/log"
)

func (a appDataService) AddTimeline(ctx context.Context, source string) (Timeline, error) {
	//TODO implement me
	panic("implement me")
}

func (a appDataService) GetTimeline(ctx context.Context, id string) (Timeline, error) {
	//TODO implement me
	panic("implement me")
}

func (a appDataService) GetAllTimelines(ctx context.Context) ([]Timeline, error) {
	//	Our return item
	retval := []Timeline{}

	//	Our temporary maps to keep track of Timelines and their step info
	timelines := map[string]Timeline{}

	query := `select
		tl.id, tl.enabled, tl.created, tl.name, tl.gpio,
		ts.id, ts.step_type_id, ts.effect_type_id, ts.led_range,
		ts.step_time, ts.step_meta, ts.step_number
	from
		timeline tl
		join timeline_step ts
			on ts.timeline_id = tl.id
	order by
    	ts.step_number;`

	stmt, err := a.DB.PreparexContext(ctx, query)
	if err != nil {
		return retval, err
	}

	rows, err := stmt.QueryxContext(ctx)
	if err != nil {
		return retval, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Err(closeErr).Msg("unable to close rows")
		}
	}()

	for rows.Next() {
		item := Timeline{
			Steps: []TimelineStep{},
			Tags:  []string{},
		}
		tlStep := TimelineStep{}

		if err := rows.Scan(&item.ID, &item.Enabled, &item.Created, &item.Name, &item.GPIO,
			&tlStep.ID, &tlStep.Type, &tlStep.Effect, &tlStep.Leds,
			&tlStep.Time, &tlStep.MetaInfo, &tlStep.Number); err != nil {
			return retval, fmt.Errorf("problem reading into struct: %v", err)
		}

		//	If the tracked timeline doesn't exist yet, add it:
		_, found := timelines[item.ID]
		if !found {
			timelines[item.ID] = item
		}

		//	If we have data in the tlStep.ID,
		//	add the tlStep information to the referenced timeline
		if tlStep.ID != "" {

			//	Convert the meta info to a json string:
			metaInfo, _ := json.Marshal(tlStep.MetaInfo)
			var jsonString string
			json.Unmarshal(metaInfo, &jsonString)

			//	... determine the step type
			switch tlStep.Type {
			case step.Effect:
				/* If it's an effect, load effect meta */
				switch tlStep.Effect {
				case effect.Unknown:
					//	We don't know what to do here
				case effect.Solid:
					em := SolidMeta{}
					err := json.Unmarshal([]byte(jsonString), &em)
					if err != nil {
						log.Err(err).Msg("Problem unmarshalling SolidMeta")
					}

					tlStep.MetaInfo = em
				case effect.Fade:
					em := FadeMeta{}
					json.Unmarshal([]byte(jsonString), &em)
					tlStep.MetaInfo = em
				case effect.Gradient:
					em := GradientMeta{}
					json.Unmarshal([]byte(jsonString), &em)
					tlStep.MetaInfo = em
				case effect.Sequence:
					em := SequenceMeta{}
					json.Unmarshal([]byte(jsonString), &em)
					tlStep.MetaInfo = em
				case effect.Rainbow:
					//	Don't need to do anything
				case effect.Zip:
					em := ZipMeta{}
					json.Unmarshal([]byte(jsonString), &em)
					tlStep.MetaInfo = em
				case effect.KnightRider:
					//	Don't need to do anything
				case effect.Lightning:
					em := LightningMeta{}
					json.Unmarshal([]byte(jsonString), &em)
					tlStep.MetaInfo = em
				}
			case step.Sleep:
			case step.RandomSleep:
			case step.Loop:
			default:
			}

			//	Get a copy
			if entry, ok := timelines[item.ID]; ok {

				// Then we modify the copy
				entry.Steps = append(timelines[item.ID].Steps, tlStep)

				// Then we reassign map entry
				timelines[item.ID] = entry
			}
		}

		////	If we have data in tags ...
		//if tags != nil {
		//	//	Unmarshal the JSON tag array
		//	if err := json.Unmarshal(tags, &item.Tags); err != nil {
		//		return retval, fmt.Errorf("problem decoding tags for %v: %v", item.ID, err)
		//	}
		//}
	}

	//	For each item in the map, assign to the output slice:
	for _, v := range timelines {
		retval = append(retval, v)
	}

	//	Return our data:
	return retval, nil
}

func (a appDataService) GetAllTimelinesWithTag(ctx context.Context, tag string) ([]Timeline, error) {
	//TODO implement me
	panic("implement me")
}

func (a appDataService) DeleteTimeline(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (a appDataService) UpdateTags(ctx context.Context, id string, tags []string) error {
	//TODO implement me
	panic("implement me")
}
