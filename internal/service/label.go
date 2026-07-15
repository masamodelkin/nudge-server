package service

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/masamodelkin/nudge-server/internal/model"
	"github.com/masamodelkin/nudge-server/internal/store"
)

var (
	ErrLabelNotFound   = errors.New("label not found")
	ErrLabelValidation = errors.New("validation failed")
	hexColorRegex      = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)
)

type LabelService struct {
	store *store.Store
}

func NewLabelService(s *store.Store) *LabelService {
	return &LabelService{store: s}
}

type LabelRequest struct {
	Name  string
	Color *string
}

func validateLabel(label *model.Label) error {
	if label.Name == "" {
		return fmt.Errorf("%w: name is required", ErrLabelValidation)
	}

	if label.Color != nil && !hexColorRegex.MatchString(*label.Color) {
		return fmt.Errorf("%w: color must be a valid hex color (e.g. #FF0000)", ErrLabelValidation)
	}

	return nil
}

func (s *LabelService) Create(userID string, req *LabelRequest) (*model.Label, error) {
	label := &model.Label{
		ID:     uuid.NewString(),
		Name:   req.Name,
		Color:  req.Color,
		UserID: userID,
	}

	if err := validateLabel(label); err != nil {
		return nil, err
	}

	if err := s.store.CreateLabel(label); err != nil {
		return nil, err
	}
	return label, nil
}

func (s *LabelService) List(userID string) ([]model.Label, error) {
	return s.store.ListLabels(userID)
}

func (s *LabelService) Update(id string, userID string, req *LabelRequest) (*model.Label, error) {
	label, err := s.store.GetLabel(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLabelNotFound
		}
		return nil, err
	}

	label.Name = req.Name
	label.Color = req.Color

	if err := validateLabel(label); err != nil {
		return nil, err
	}

	if err := s.store.UpdateLabel(label); err != nil {
		return nil, err
	}
	return label, nil
}

func (s *LabelService) Delete(id string, userID string) error {
	_, err := s.store.GetLabel(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrLabelNotFound
		}
		return err
	}
	return s.store.DeleteLabel(id, userID)
}
