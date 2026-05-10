package store

import (
	"github.com/masamodelkin/nudge-server/internal/model"
)

func (s *Store) CreateUser(user *model.User) error {
	_, err := s.db.NamedExec(
		`INSERT INTO users (id, username, password_hash, email, auth_provider, provider_id)
         VALUES (:id, :username, :password_hash, :email, :auth_provider, :provider_id)`,
		user,
	)
	return err
}

func (s *Store) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := s.db.Get(&user,
		`SELECT id, username, password_hash, email, auth_provider, created_at
         FROM users WHERE username = ?`,
		username,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUserByID(id string) (*model.User, error) {
	var user model.User
	err := s.db.Get(&user,
		`SELECT id, username, password_hash, email, auth_provider, created_at
         FROM users WHERE id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) CreateRefreshToken(token *model.RefreshToken) error {
	_, err := s.db.NamedExec(
		`INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
         VALUES (:id, :user_id, :token_hash, :expires_at)`,
		token,
	)
	return err
}

func (s *Store) GetRefreshToken(id string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := s.db.Get(&token,
		`SELECT id, user_id, token_hash, expires_at, created_at
         FROM refresh_tokens WHERE id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (s *Store) GetRefreshTokenByHash(hash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := s.db.Get(&token,
		`SELECT id, user_id, token_hash, expires_at, created_at
         FROM refresh_tokens WHERE token_hash = ?`,
		hash,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (s *Store) DeleteRefreshToken(id string) error {
	_, err := s.db.Exec("DELETE FROM refresh_tokens WHERE id = ?", id)
	return err
}

func (s *Store) DeleteRefreshTokenByHash(hash string) error {
	_, err := s.db.Exec("DELETE FROM refresh_tokens WHERE token_hash = ?", hash)
	return err
}
