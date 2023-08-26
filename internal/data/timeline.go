package data

import (
	"context"
)

func (a appDataService) AddTimeline(ctx context.Context, name, description string) (Timeline, error) {
	//TODO implement me
	panic("implement me")
}

func (a appDataService) GetTimeline(ctx context.Context, id string) (Timeline, error) {
	//TODO implement me
	panic("implement me")
}

func (a appDataService) GetAllTimelines(ctx context.Context) ([]Timeline, error) {
	//TODO implement me
	panic("implement me")
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
