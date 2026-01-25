package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/repository"
)

type customerUserRepository struct {
	db *sql.DB
}

// NewCustomerUserRepository creates a new PostgreSQL customer user repository
func NewCustomerUserRepository(db *sql.DB) repository.CustomerUserRepository {
	return &customerUserRepository{db: db}
}

func (r *customerUserRepository) Create(ctx context.Context, user *domain.CustomerUser) error {
	query := `
		INSERT INTO customer_users (
			id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.EmailVerified,
		user.EmailVerificationToken,
		user.EmailVerificationExpiresAt,
		user.PasswordResetToken,
		user.PasswordResetExpiresAt,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create customer user: %w", err)
	}

	return nil
}

func (r *customerUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.CustomerUser, error) {
	query := `
		SELECT id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		FROM customer_users
		WHERE id = $1
	`

	user := &domain.CustomerUser{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.EmailVerified,
		&user.EmailVerificationToken,
		&user.EmailVerificationExpiresAt,
		&user.PasswordResetToken,
		&user.PasswordResetExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("customer user not found")
		}
		return nil, fmt.Errorf("failed to find customer user: %w", err)
	}

	return user, nil
}

func (r *customerUserRepository) FindByEmail(ctx context.Context, email string) (*domain.CustomerUser, error) {
	query := `
		SELECT id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		FROM customer_users
		WHERE email = $1
	`

	user := &domain.CustomerUser{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.EmailVerified,
		&user.EmailVerificationToken,
		&user.EmailVerificationExpiresAt,
		&user.PasswordResetToken,
		&user.PasswordResetExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found is not an error for FindByEmail
		}
		return nil, fmt.Errorf("failed to find customer user by email: %w", err)
	}

	return user, nil
}

func (r *customerUserRepository) FindByVerificationToken(ctx context.Context, token string) (*domain.CustomerUser, error) {
	query := `
		SELECT id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		FROM customer_users
		WHERE email_verification_token = $1
	`

	user := &domain.CustomerUser{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.EmailVerified,
		&user.EmailVerificationToken,
		&user.EmailVerificationExpiresAt,
		&user.PasswordResetToken,
		&user.PasswordResetExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("verification token not found")
		}
		return nil, fmt.Errorf("failed to find customer user by verification token: %w", err)
	}

	return user, nil
}

func (r *customerUserRepository) FindByResetToken(ctx context.Context, token string) (*domain.CustomerUser, error) {
	query := `
		SELECT id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		FROM customer_users
		WHERE password_reset_token = $1
	`

	user := &domain.CustomerUser{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.EmailVerified,
		&user.EmailVerificationToken,
		&user.EmailVerificationExpiresAt,
		&user.PasswordResetToken,
		&user.PasswordResetExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("reset token not found")
		}
		return nil, fmt.Errorf("failed to find customer user by reset token: %w", err)
	}

	return user, nil
}

func (r *customerUserRepository) Update(ctx context.Context, user *domain.CustomerUser) error {
	query := `
		UPDATE customer_users
		SET email = $2,
			password_hash = $3,
			email_verified = $4,
			email_verification_token = $5,
			email_verification_expires_at = $6,
			password_reset_token = $7,
			password_reset_expires_at = $8,
			updated_at = $9
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.EmailVerified,
		user.EmailVerificationToken,
		user.EmailVerificationExpiresAt,
		user.PasswordResetToken,
		user.PasswordResetExpiresAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update customer user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("customer user not found")
	}

	return nil
}

func (r *customerUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM customer_users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("customer user not found")
	}

	return nil
}
