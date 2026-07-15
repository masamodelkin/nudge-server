package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	"github.com/masamodelkin/nudge-server/internal/model"
	"github.com/masamodelkin/nudge-server/internal/store"
)

var (
	ErrTriggerNotFound   = errors.New("trigger not found")
	ErrTriggerValidation = errors.New("validation failed")
)

type TriggerService struct {
	store *store.Store
}

func NewTriggerService(s *store.Store) *TriggerService {
	return &TriggerService{store: s}
}

type TriggerRequest struct {
	Name        string
	Type        string
	Config      types.JSONText
	IsExclusive bool
}

type triggerLocationConfig struct {
	Lat    *float64 `json:"lat"`
	Lng    *float64 `json:"lng"`
	Radius *float64 `json:"radius"`
}

type triggerTimeConfig struct {
	Start *int   `json:"start"`
	End   *int   `json:"end"`
	Days  []bool `json:"days"`
}

type triggerDeviceConfig struct {
	Name *string `json:"name"`
	Type *string `json:"type"`
}

type triggerWifiConfig struct {
	Ssid *string `json:"ssid"`
}

func validateTriggerConfig(configType string, config *types.JSONText) error {
	switch configType {
	case "location":
		var c triggerLocationConfig
		if err := config.Unmarshal(&c); err != nil {
			return fmt.Errorf("%w: invalid location config", ErrTriggerValidation)
		}
		if c.Lat == nil || c.Lng == nil || c.Radius == nil {
			return fmt.Errorf("%w: lat, lng, and radius are all required fields", ErrTriggerValidation)
		}
		if *c.Radius <= 0 {
			return fmt.Errorf("%w: radius must be positive", ErrTriggerValidation)
		}
		return nil
	case "time":
		var c triggerTimeConfig
		if err := config.Unmarshal(&c); err != nil {
			return fmt.Errorf("%w: invalid time config", ErrTriggerValidation)
		}
		if c.Start == nil && c.End == nil {
			return fmt.Errorf("%w: at least one of start and end time is required for the time config", ErrTriggerValidation)
		}
		if c.Days != nil && len(c.Days) != 7 {
			return fmt.Errorf("%w: days should be a bool array of lenght 7", ErrTriggerValidation)
		}
		return nil
	case "device":
		var c triggerDeviceConfig
		if err := config.Unmarshal(&c); err != nil {
			return fmt.Errorf("%w: invalid device config", ErrTriggerValidation)
		}
		if c.Name == nil && c.Type == nil {
			return fmt.Errorf("%w: at least one of name and type is required for the device config", ErrTriggerValidation)
		}
		return nil
	case "wifi":
		var c triggerWifiConfig
		if err := config.Unmarshal(&c); err != nil {
			return fmt.Errorf("%w: invalid wifi config", ErrTriggerValidation)
		}
		if c.Ssid == nil {
			return fmt.Errorf("%w: ssid is required for wifi trigger", ErrTriggerValidation)
		}
		return nil
	default:
		return fmt.Errorf("%w: invalid trigger type. Available types: location, time, device, wifi", ErrTriggerValidation)
	}
}

func (s *TriggerService) validateTrigger(trigger *model.Trigger) error {
	if trigger.Name == "" {
		return fmt.Errorf("%w: trigger name is required", ErrTriggerValidation)
	}
	if trigger.Type == "" {
		return fmt.Errorf("%w: trigger type is required", ErrTriggerValidation)
	}

	if err := validateTriggerConfig(trigger.Type, &trigger.Config); err != nil {
		return err
	}
	return nil
}

func (s *TriggerService) Create(userID string, req *TriggerRequest) (*model.Trigger, error) {
	trigger := &model.Trigger{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Type:        req.Type,
		Config:      req.Config,
		IsExclusive: req.IsExclusive,
		UserID:      userID,
	}

	if err := s.validateTrigger(trigger); err != nil {
		return nil, err
	}

	if err := s.store.CreateTrigger(trigger); err != nil {
		return nil, err
	}

	return trigger, nil
}

func (s *TriggerService) Get(id string, userID string) (*model.Trigger, error) {
	trigger, err := s.store.GetTrigger(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTriggerNotFound
		}
		return nil, err
	}
	return trigger, nil
}

func (s *TriggerService) List(userID string) ([]model.Trigger, error) {
	return s.store.ListTriggers(userID)
}

func (s *TriggerService) Update(id string, userID string, req *TriggerRequest) (*model.Trigger, error) {
	trigger, err := s.store.GetTrigger(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTriggerNotFound
		}
		return nil, err
	}

	trigger.Name = req.Name
	trigger.Type = req.Type
	trigger.Config = req.Config
	trigger.IsExclusive = req.IsExclusive

	if err := s.validateTrigger(trigger); err != nil {
		return nil, err
	}

	if err := s.store.UpdateTrigger(trigger); err != nil {
		return nil, err
	}

	return trigger, nil
}

func (s *TriggerService) Delete(id string, userID string) error {
	_, err := s.store.GetTrigger(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTriggerNotFound
		}
		return err
	}
	return s.store.DeleteTrigger(id, userID)
}
