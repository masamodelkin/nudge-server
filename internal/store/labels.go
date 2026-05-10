package store

import (
	"github.com/masamodelkin/nudge-server/internal/model"
)

func (s *Store) CreateLabel(label *model.Label) error {
	_, err := s.db.NamedExec(
		`INSERT INTO labels (id, name, color, user_id)
         VALUES (:id, :name, :color, :user_id)`,
		label,
	)
	return err
}

func (s *Store) ListLabels(userID string) ([]model.Label, error) {
	var labels []model.Label
	err := s.db.Select(&labels,
		"SELECT id, name, color, user_id FROM labels WHERE user_id = ?",
		userID,
	)
	return labels, err
}

func (s *Store) GetLabel(id string, userID string) (*model.Label, error) {
	var label model.Label
	err := s.db.Get(&label,
		"SELECT id, name, color, user_id FROM labels WHERE id = ? AND user_id = ?",
		id, userID,
	)
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (s *Store) DeleteLabel(id string, userID string) error {
	_, err := s.db.Exec(
		"DELETE FROM labels WHERE id = ? AND user_id = ?",
		id, userID,
	)
	return err
}
