package store

import (
	"github.com/masamodelkin/nudge-server/internal/model"
)

func (s *Store) CreateStatus(status *model.Status) error {
	_, err := s.db.NamedExec(
		`INSERT INTO statuses (id, name, user_id)
         VALUES (:id, :name, :user_id)`,
		status,
	)
	return err
}

func (s *Store) ListStatuses(userID string) ([]model.Status, error) {
	var statuses []model.Status
	err := s.db.Select(&statuses,
		"SELECT id, name, user_id FROM statuses WHERE user_id = ?",
		userID,
	)
	return statuses, err
}

func (s *Store) GetStatus(id string, userID string) (*model.Status, error) {
	var status model.Status
	err := s.db.Get(&status,
		"SELECT id, name, user_id FROM statuses WHERE id = ? AND user_id = ?",
		id, userID,
	)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (s *Store) DeleteStatus(id string, userID string) error {
	_, err := s.db.Exec(
		"DELETE FROM statuses WHERE id = ? AND user_id = ?",
		id, userID,
	)
	return err
}
