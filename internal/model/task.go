package model

type Task struct {
	ID          string  `db:"id"          json:"id"`
	Name        *string `db:"name"        json:"name"`
	Description *string `db:"description" json:"description"`
	IsDraft     bool    `db:"is_draft"    json:"is_draft"`
	DueDate     *int64  `db:"due_date"    json:"due_date"`
	Priority    *int    `db:"priority"    json:"priority"`
	Duration    *int    `db:"duration"    json:"duration"`
	TimeSpent   int     `db:"time_spent"  json:"time_spent"`
	StatusID    *string `db:"status_id"   json:"status_id"`
	UserID      string  `db:"user_id"     json:"-"`
	CreatedAt   int64   `db:"created_at"  json:"created_at"`
	UpdatedAt   int64   `db:"updated_at"  json:"updated_at"`
	Labels      []Label `db:"-"           json:"labels"`
}
