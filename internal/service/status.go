package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/masamodelkin/nudge-server/internal/model"
	"github.com/masamodelkin/nudge-server/internal/store"
)

var (
	ErrStatusNotFound   = errors.New("status not found")
	ErrStatusValidation = errors.New("validation failed")
)

type StatusService struct {
	store *store.Store
}

func NewStatusService(s *store.Store) *StatusService {
	return &StatusService{store: s}
}

type StatusRequest struct {
	Name         string
	NextStatusID *string
	IsDone       bool
}

func (s *StatusService) validateStatus(status *model.Status) error {
	if status.Name == "" {
		return fmt.Errorf("%w: name is required", ErrStatusValidation)
	}

	if status.NextStatusID != nil {
		if *status.NextStatusID == status.ID {
			return fmt.Errorf("%w: status cannot point to itself", ErrStatusValidation)
		}

		_, err := s.store.GetStatus(*status.NextStatusID, status.UserID)
		if err != nil {
			return fmt.Errorf("%w: next status not found", ErrStatusValidation)
		}
	}

	return nil
}

func (s *StatusService) Create(userID string, req *StatusRequest) (*model.Status, error) {
	status := &model.Status{
		ID:           uuid.New().String(),
		Name:         req.Name,
		NextStatusID: req.NextStatusID,
		IsDone:       req.IsDone,
		UserID:       userID,
	}

	if err := s.validateStatus(status); err != nil {
		return nil, err
	}

	if err := s.store.CreateStatus(status); err != nil {
		return nil, err
	}
	return status, nil
}

func (s *StatusService) List(userID string) ([]model.Status, error) {
	return s.store.ListStatuses(userID)
}

func (s *StatusService) Update(id string, userID string, req *StatusRequest) (*model.Status, error) {
	status, err := s.store.GetStatus(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStatusNotFound
		}
		return nil, err
	}

	status.Name = req.Name
	status.NextStatusID = req.NextStatusID
	status.IsDone = req.IsDone

	if err := s.validateStatus(status); err != nil {
		return nil, err
	}

	if err := s.store.UpdateStatus(status); err != nil {
		return nil, err
	}
	return status, nil
}

func (s *StatusService) Delete(id string, userID string) error {
	_, err := s.store.GetStatus(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrStatusNotFound
		}
		return err
	}
	return s.store.DeleteStatus(id, userID)
}
