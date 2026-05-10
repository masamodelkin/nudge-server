package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/masamodelkin/nudge-server/internal/model"
	"github.com/masamodelkin/nudge-server/internal/store"
)

var ErrLabelNotFound = errors.New("label not found")

type LabelService struct {
	store *store.Store
}

func NewLabelService(s *store.Store) *LabelService {
	return &LabelService{store: s}
}

func (s *LabelService) Create(userID string, name string, color *string) (*model.Label, error) {
	label := &model.Label{
		ID:     uuid.New().String(),
		Name:   name,
		Color:  color,
		UserID: userID,
	}
	if err := s.store.CreateLabel(label); err != nil {
		return nil, err
	}
	return label, nil
}

func (s *LabelService) List(userID string) ([]model.Label, error) {
	return s.store.ListLabels(userID)
}

func (s *LabelService) Delete(id string, userID string) error {
	_, err := s.store.GetLabel(id, userID)
	if err != nil {
		return ErrLabelNotFound
	}
	return s.store.DeleteLabel(id, userID)
}
