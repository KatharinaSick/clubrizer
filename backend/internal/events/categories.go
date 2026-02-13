package events

import "context"

func (s *Service) ListCategories(ctx context.Context) ([]*Category, error) {
	categories, err := s.store.getAllCategories(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
