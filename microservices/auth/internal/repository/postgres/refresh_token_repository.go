package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/domain"
)

type refreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *refreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, revoked, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx, query,
		token.ID, token.UserID, token.Token, token.ExpiresAt, token.Revoked, token.CreatedAt,
	)

	return err
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, revoked, created_at
		FROM refresh_tokens
		WHERE token = $1 AND revoked = false
	`

	refreshToken := &domain.RefreshToken{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token,
		&refreshToken.ExpiresAt, &refreshToken.Revoked, &refreshToken.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("refresh token not found")
	}

	return refreshToken, err
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, token string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *refreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
