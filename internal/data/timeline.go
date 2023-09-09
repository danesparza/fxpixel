package data

import (
	"context"
	"fmt"
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
		ts.step_time, ts.step_number
	from
		timeline tl
		join timeline_step ts
			on ts.timeline_id = tl.id;`

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
		step := TimelineStep{}

		if err := rows.Scan(&item.ID, &item.Enabled, &item.Created, &item.Name, &item.GPIO,
			&step.ID, &step.Type, &step.Effect, &step.Leds,
			&step.Time, &step.Number); err != nil {
			return retval, fmt.Errorf("problem reading into struct: %v", err)
		}

		//	If the tracked timeline doesn't exist yet, add it:
		_, found := timelines[item.ID]
		if !found {
			timelines[item.ID] = item
		}

		//	If we have data in the step.ID,
		//	add the step information to the referenced timeline
		if step.ID != "" {

			//	Get a copy
			if entry, ok := timelines[item.ID]; ok {

				// Then we modify the copy
				entry.Steps = append(timelines[item.ID].Steps, step)

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
