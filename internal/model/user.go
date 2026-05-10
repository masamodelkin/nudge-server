package model

type User struct {
	ID           string  `db:"id"            json:"id"`
	Username     string  `db:"username"      json:"username"`
	PasswordHash string  `db:"password_hash" json:"-"`
	Email        *string `db:"email"         json:"email"`
	AuthProvider string  `db:"auth_provider" json:"auth_provider"`
	ProviderID   *string `db:"provider_id"   json:"provider_id,omitempty"`
	CreatedAt    int64   `db:"created_at"    json:"created_at"`
}

type RefreshToken struct {
	ID        string `db:"id"`
	UserID    string `db:"user_id"`
	TokenHash string `db:"token_hash"`
	ExpiresAt int64  `db:"expires_at"`
	CreatedAt int64  `db:"created_at"`
}
