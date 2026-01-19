package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/auth/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			id, email, password_hash, role, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	
	_, err := r.db.ExecContext(
		ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Role, user.EmailVerified,
		user.EmailVerificationToken, user.EmailVerificationExpiresAt,
		user.PasswordResetToken, user.PasswordResetExpiresAt,
		user.CreatedAt, user.UpdatedAt,
	)
	
	return err
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT 
			id, email, password_hash, role, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		FROM users
		WHERE id = $1
	`
	
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.EmailVerified,
		&user.EmailVerificationToken, &user.EmailVerificationExpiresAt,
		&user.PasswordResetToken, &user.PasswordResetExpiresAt,
		&user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	
	return user, err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT 
			id, email, password_hash, role, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		FROM users
		WHERE email = $1
	`
	
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.EmailVerified,
		&user.EmailVerificationToken, &user.EmailVerificationExpiresAt,
		&user.PasswordResetToken, &user.PasswordResetExpiresAt,
		&user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	
	return user, err
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET
			email = $2,
			password_hash = $3,
			role = $4,
			email_verified = $5,
			email_verification_token = $6,
			email_verification_expires_at = $7,
			password_reset_token = $8,
			password_reset_expires_at = $9,
			updated_at = $10
		WHERE id = $1
	`
	
	user.UpdatedAt = time.Now()
	
	_, err := r.db.ExecContext(
		ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Role, user.EmailVerified,
		user.EmailVerificationToken, user.EmailVerificationExpiresAt,
		user.PasswordResetToken, user.PasswordResetExpiresAt,
		user.UpdatedAt,
	)
	
	return err
}

func (r *userRepository) UpdateEmailVerification(ctx context.Context, userID uuid.UUID, verified bool) error {
	query := `UPDATE users SET email_verified = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID, verified, time.Now())
	return err
}

func (r *userRepository) UpdatePasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt *time.Time) error {
	query := `
		UPDATE users SET 
			password_reset_token = $2,
			password_reset_expires_at = $3,
			updated_at = $4
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt, time.Now())
	return err
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID, passwordHash, time.Now())
	return err
}
