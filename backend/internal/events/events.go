package events

import "context"

type ListEventsResponse struct {
	Foo string `json:"foo"`
}

func (s *Service) ListEvents(_ context.Context) (*ListEventsResponse, error) {
	return &ListEventsResponse{Foo: "bar!"}, nil
}
