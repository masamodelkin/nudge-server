package store

import (
	"github.com/masamodelkin/nudge-server/internal/model"
)

func (s *Store) CreateTrigger(trigger *model.Trigger) error {
	_, err := s.db.NamedExec(
		`INSERT INTO triggers (id, name, type, config, is_exclusive, user_id)
         VALUES (:id, :name, :type, :config, :is_exclusive, :user_id)`,
		trigger,
	)
	return err
}

func (s *Store) ListTriggers(userID string) ([]model.Trigger, error) {
	var triggers []model.Trigger
	err := s.db.Select(&triggers,
		"SELECT id, name, type, config, is_exclusive, user_id FROM triggers WHERE user_id = ?",
		userID,
	)
	return triggers, err
}

func (s *Store) GetTrigger(id string, userID string) (*model.Trigger, error) {
	var trigger model.Trigger
	err := s.db.Get(&trigger,
		"SELECT id, name, type, config, is_exclusive, user_id FROM triggers WHERE id = ? AND user_id = ?",
		id, userID,
	)
	if err != nil {
		return nil, err
	}
	return &trigger, nil
}

func (s *Store) UpdateTrigger(trigger *model.Trigger) error {
	_, err := s.db.NamedExec(
		"UPDATE triggers SET name = :name, type = :type, config = :config, is_exclusive = :is_exclusive WHERE id = :id AND user_id = :user_id",
		trigger,
	)
	return err
}

func (s *Store) DeleteTrigger(id string, userID string) error {
	_, err := s.db.Exec(
		"DELETE FROM triggers WHERE id = ? AND user_id = ?",
		id, userID,
	)
	return err
}
