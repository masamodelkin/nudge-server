package model

type Label struct {
	ID     string  `db:"id"      json:"id"`
	Name   string  `db:"name"    json:"name"`
	Color  *string `db:"color"   json:"color"`
	UserID string  `db:"user_id" json:"-"`
}
