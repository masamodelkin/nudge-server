package model

import "github.com/jmoiron/sqlx/types"

type Trigger struct {
	ID          string         `db:"id"             json:"id"`
	Name        string         `db:"name"           json:"name"`
	Type        string         `db:"type"           json:"type"`
	Config      types.JSONText `db:"config"         json:"config"`
	IsExclusive bool           `db:"is_exclusive"   json:"is_exclusive"`
	UserID      string         `db:"user_id"        json:"-"`
}
