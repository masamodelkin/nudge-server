package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/masamodelkin/nudge-server/internal/model"
	"github.com/masamodelkin/nudge-server/internal/store"
)

var (
	ErrValidation   = errors.New("validation failed")
	ErrTaskNotFound = errors.New("task not found")
)

type TaskRequest struct {
	Name        *string
	Description *string
	IsDraft     bool
	DueDate     *int64
	Priority    *int
	Duration    *int
	StatusID    *string
	LabelIDs    []string
	TriggerIDs  []string
}

type TaskService struct {
	store *store.Store
}

func NewTaskService(s *store.Store) *TaskService {
	return &TaskService{store: s}
}

func (s *TaskService) validateTask(task *model.Task, labelIDs []string, triggerIDs []string, userID string) error {
	if task.IsDraft {
		if (task.Name == nil || *task.Name == "") && (task.Description == nil || *task.Description == "") {
			return fmt.Errorf("%w: draft requires at least a name or description", ErrValidation)
		}
	} else {
		if task.Name == nil || *task.Name == "" {
			return fmt.Errorf("%w: name is required", ErrValidation)
		}

		if task.DueDate != nil && *task.DueDate < 0 {
			return fmt.Errorf("%w: due date must be positive", ErrValidation)
		}

		if task.Duration == nil {
			return fmt.Errorf("%w: duration is required", ErrValidation)
		}

		if *task.Duration < 0 {
			return fmt.Errorf("%w: duration must be non-negative", ErrValidation)
		}

		if task.DueDate == nil && task.Priority == nil {
			return fmt.Errorf("%w: due date or priority is required", ErrValidation)
		}
	}

	if task.StatusID != nil {
		_, err := s.store.GetStatus(*task.StatusID, userID)
		if err != nil {
			return fmt.Errorf("%w: status not found", ErrValidation)
		}
	}

	for _, labelID := range labelIDs {
		_, err := s.store.GetLabel(labelID, userID)
		if err != nil {
			return fmt.Errorf("%w: label not found", ErrValidation)
		}
	}

	for _, triggerID := range triggerIDs {
		_, err := s.store.GetTrigger(triggerID, userID)
		if err != nil {
			return fmt.Errorf("%w: trigger not found", ErrValidation)
		}
	}

	return nil
}

func (s *TaskService) Create(userID string, req *TaskRequest) (*model.Task, error) {
	task := &model.Task{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		IsDraft:     req.IsDraft,
		DueDate:     req.DueDate,
		Priority:    req.Priority,
		Duration:    req.Duration,
		TimeSpent:   0,
		StatusID:    req.StatusID,
		UserID:      userID,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	if err := s.validateTask(task, req.LabelIDs, req.TriggerIDs, userID); err != nil {
		return nil, err
	}

	if err := s.store.CreateTask(task, req.LabelIDs, req.TriggerIDs); err != nil {
		return nil, err
	}

	return s.store.GetTask(task.ID, userID)
}

func (s *TaskService) Get(id string, userID string) (*model.Task, error) {
	task, err := s.store.GetTask(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return task, nil
}

func (s *TaskService) List(userID string) ([]model.Task, error) {
	return s.store.ListTasks(userID)
}

func (s *TaskService) Update(id string, userID string, req *TaskRequest) (*model.Task, error) {
	task, err := s.store.GetTask(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	task.Name = req.Name
	task.Description = req.Description
	task.IsDraft = req.IsDraft
	task.DueDate = req.DueDate
	task.Priority = req.Priority
	task.Duration = req.Duration
	task.StatusID = req.StatusID

	if err := s.validateTask(task, req.LabelIDs, req.TriggerIDs, userID); err != nil {
		return nil, err
	}

	if err := s.store.UpdateTask(task, req.LabelIDs, req.TriggerIDs); err != nil {
		return nil, err
	}

	return s.store.GetTask(id, userID)
}

func (s *TaskService) Delete(id string, userID string) error {
	_, err := s.store.GetTask(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	return s.store.DeleteTask(id, userID)
}

func (s *TaskService) AddTime(id string, userID string, timeAdded int) (*model.Task, error) {
	task, err := s.store.GetTask(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	if task.TimeSpent+timeAdded < 0 {
		return nil, fmt.Errorf("%w: time spent cannot be negative", ErrValidation)
	}

	if err := s.store.AddTaskTime(id, userID, timeAdded); err != nil {
		return nil, err
	}

	return s.store.GetTask(id, userID)
}

func (s *TaskService) SetTime(id string, userID string, timeSpent int) (*model.Task, error) {
	_, err := s.store.GetTask(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	if timeSpent < 0 {
		return nil, fmt.Errorf("%w: time spent cannot be negative", ErrValidation)
	}

	if err := s.store.SetTaskTime(id, userID, timeSpent); err != nil {
		return nil, err
	}

	return s.store.GetTask(id, userID)
}
