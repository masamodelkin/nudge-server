package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/masamodelkin/nudge-server/internal/model"
	"github.com/masamodelkin/nudge-server/internal/store"
)

var ErrStatusNotFound = errors.New("status not found")

type StatusService struct {
	store *store.Store
}

func NewStatusService(s *store.Store) *StatusService {
	return &StatusService{store: s}
}

func (s *StatusService) Create(userID string, name string) (*model.Status, error) {
	status := &model.Status{
		ID:     uuid.New().String(),
		Name:   name,
		UserID: userID,
	}
	if err := s.store.CreateStatus(status); err != nil {
		return nil, err
	}
	return status, nil
}

func (s *StatusService) List(userID string) ([]model.Status, error) {
	return s.store.ListStatuses(userID)
}

func (s *StatusService) Delete(id string, userID string) error {
	_, err := s.store.GetStatus(id, userID)
	if err != nil {
		return ErrStatusNotFound
	}
	return s.store.DeleteStatus(id, userID)
}
