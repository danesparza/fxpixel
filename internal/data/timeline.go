package data

import (
	"context"
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

	//query := `select id, filepath, description, created, tags
	//			from media;`
	//
	//stmt, err := a.DB.PreparexContext(ctx, query)
	//if err != nil {
	//	return retval, err
	//}
	//
	//rows, err := stmt.QueryxContext(ctx)
	//if err != nil {
	//	return retval, err
	//}
	//
	//defer func() {
	//	if closeErr := rows.Close(); closeErr != nil {
	//		log.Err(closeErr).Msg("unable to close rows")
	//	}
	//}()
	//
	//for rows.Next() {
	//	item := File{}
	//	tags := []byte{}
	//	if err := rows.Scan(&item.ID, &item.FilePath, &item.Description, &item.Created, &tags); err != nil {
	//		return retval, fmt.Errorf("problem reading into struct: %v", err)
	//	}
	//
	//	//	If we have data in tags ...
	//	if tags != nil {
	//		//	Unmarshal the JSON tag array
	//		if err := json.Unmarshal(tags, &item.Tags); err != nil {
	//			return retval, fmt.Errorf("problem decoding tags for %v: %v", item.ID, err)
	//		}
	//	}
	//
	//	retval = append(retval, item)
	//}

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
