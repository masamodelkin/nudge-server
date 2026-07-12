package model

type Status struct {
	ID           string  `db:"id"             json:"id"`
	Name         string  `db:"name"           json:"name"`
	NextStatusID *string `db:"next_status_id" json:"next_status_id"`
	IsDone       bool    `db:"is_done"        json:"is_done"`
	UserID       string  `db:"user_id"        json:"-"`
}
