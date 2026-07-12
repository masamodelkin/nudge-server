package store

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/masamodelkin/nudge-server/internal/model"
)

func (s *Store) CreateTask(task *model.Task, labelIDs []string) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExec(
		`INSERT INTO tasks (id, name, description, is_draft, due_date, priority, duration, time_spent, status_id, user_id, created_at, updated_at)
         VALUES (:id, :name, :description, :is_draft, :due_date, :priority, :duration, :time_spent, :status_id, :user_id, :created_at, :updated_at)`,
		task,
	)
	if err != nil {
		return err
	}

	for _, labelID := range labelIDs {
		_, err = tx.Exec(
			"INSERT OR IGNORE INTO task_labels (task_id, label_id) VALUES (?, ?)",
			task.ID, labelID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) GetTask(id string, userID string) (*model.Task, error) {
	var task model.Task
	err := s.db.Get(&task,
		`SELECT id, name, description, is_draft, due_date, priority, duration, time_spent, status_id, user_id, created_at, updated_at
         FROM tasks WHERE id = ? AND user_id = ?`,
		id, userID,
	)
	if err != nil {
		return nil, err
	}

	task.Labels, _ = s.GetTaskLabels(task.ID)
	if task.Labels == nil {
		task.Labels = []model.Label{}
	}

	return &task, nil
}

func (s *Store) ListTasks(userID string) ([]model.Task, error) {
	var tasks []model.Task
	err := s.db.Select(&tasks,
		`SELECT id, name, description, is_draft, due_date, priority, duration, time_spent, status_id, user_id, created_at, updated_at
         FROM tasks WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return []model.Task{}, nil
	}

	taskIDs := make([]string, len(tasks))
	for i, t := range tasks {
		taskIDs[i] = t.ID
	}

	labelsMap, _ := s.GetTasksLabels(taskIDs)

	for i := range tasks {
		tasks[i].Labels = labelsMap[tasks[i].ID]
		if tasks[i].Labels == nil {
			tasks[i].Labels = []model.Label{}
		}
	}

	return tasks, nil
}

func (s *Store) UpdateTask(task *model.Task, labelIDs []string) error {
	task.UpdatedAt = time.Now().Unix()

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExec(
		`UPDATE tasks SET name = :name, description = :description, is_draft = :is_draft,
         due_date = :due_date, priority = :priority, duration = :duration,
         time_spent = :time_spent, status_id = :status_id, updated_at = :updated_at
         WHERE id = :id AND user_id = :user_id`,
		task,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM task_labels WHERE task_id = ?", task.ID)
	if err != nil {
		return err
	}

	for _, labelID := range labelIDs {
		_, err = tx.Exec("INSERT INTO task_labels (task_id, label_id) VALUES (?, ?)", task.ID, labelID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) DeleteTask(id string, userID string) error {
	_, err := s.db.Exec(
		"DELETE FROM tasks WHERE id = ? AND user_id = ?",
		id, userID,
	)
	return err
}

func (s *Store) GetTaskLabels(taskID string) ([]model.Label, error) {
	var labels []model.Label
	err := s.db.Select(&labels,
		`SELECT l.id, l.name, l.color, l.user_id
         FROM labels l
         JOIN task_labels tl ON l.id = tl.label_id
         WHERE tl.task_id = ?`,
		taskID,
	)
	return labels, err
}

func (s *Store) GetTasksLabels(taskIDs []string) (map[string][]model.Label, error) {
	if len(taskIDs) == 0 {
		return make(map[string][]model.Label), nil
	}

	type taskLabel struct {
		TaskID string `db:"task_id"`
		model.Label
	}

	query, args, err := sqlx.In(
		`SELECT tl.task_id, l.id, l.name, l.color, l.user_id
         FROM labels l
         JOIN task_labels tl ON l.id = tl.label_id
         WHERE tl.task_id IN (?)`,
		taskIDs,
	)
	if err != nil {
		return nil, err
	}

	var rows []taskLabel
	err = s.db.Select(&rows, s.db.Rebind(query), args...)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]model.Label)
	for _, row := range rows {
		result[row.TaskID] = append(result[row.TaskID], row.Label)
	}
	return result, nil
}

func (s *Store) AddTaskTime(id string, userID string, seconds int) error {
	_, err := s.db.Exec(
		"UPDATE tasks SET time_spent = time_spent + ?, updated_at = ? WHERE id = ? AND user_id = ?",
		seconds, time.Now().Unix(), id, userID,
	)
	return err
}

func (s *Store) SetTaskTime(id string, userID string, seconds int) error {
	_, err := s.db.Exec(
		"UPDATE tasks SET time_spent = ?, updated_at = ? WHERE id = ? AND user_id = ?",
		seconds, time.Now().Unix(), id, userID,
	)
	return err
}
