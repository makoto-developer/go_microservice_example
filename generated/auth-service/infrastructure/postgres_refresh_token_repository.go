package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/auth-service/domain"
)

type PostgresRefreshTokenRepository struct {
	db *sql.DB
}

func NewPostgresRefreshTokenRepository(db *sql.DB) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{db: db}
}

func (r *PostgresRefreshTokenRepository) Create(ctx context.Context, refreshToken *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (
			id, user_id, token, expires_at, revoked, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		refreshToken.Id,
		refreshToken.UserId,
		refreshToken.Token,
		refreshToken.ExpiresAt,
		refreshToken.Revoked,
		refreshToken.CreatedAt,
	)

	return err
}

func (r *PostgresRefreshTokenRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, revoked, created_at
		FROM refresh_tokens
		WHERE id = $1
	`

	var refreshToken domain.RefreshToken

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&refreshToken.Id,
		&refreshToken.UserId,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.Revoked,
		&refreshToken.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("refresh token not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *PostgresRefreshTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, revoked, created_at
		FROM refresh_tokens
		WHERE token = $1
	`

	var refreshToken domain.RefreshToken

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&refreshToken.Id,
		&refreshToken.UserId,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.Revoked,
		&refreshToken.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("refresh token not found with token: %s", token)
	}
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *PostgresRefreshTokenRepository) Update(ctx context.Context, refreshToken *domain.RefreshToken) error {
	query := `
		UPDATE refresh_tokens
		SET user_id = $1,
			token = $2,
			expires_at = $3,
			revoked = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(ctx, query,
		refreshToken.UserId,
		refreshToken.Token,
		refreshToken.ExpiresAt,
		refreshToken.Revoked,
		refreshToken.Id,
	)

	return err
}

func (r *PostgresRefreshTokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresRefreshTokenRepository) List(ctx context.Context, limit, offset int) ([]*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, revoked, created_at
		FROM refresh_tokens
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refreshTokens []*domain.RefreshToken
	for rows.Next() {
		var refreshToken domain.RefreshToken

		err := rows.Scan(
			&refreshToken.Id,
			&refreshToken.UserId,
			&refreshToken.Token,
			&refreshToken.ExpiresAt,
			&refreshToken.Revoked,
			&refreshToken.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		refreshTokens = append(refreshTokens, &refreshToken)
	}

	return refreshTokens, rows.Err()
}
