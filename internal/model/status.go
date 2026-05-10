package model

type Status struct {
	ID     string `db:"id"      json:"id"`
	Name   string `db:"name"    json:"name"`
	UserID string `db:"user_id" json:"-"`
}
