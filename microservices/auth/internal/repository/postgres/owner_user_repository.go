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

type ownerUserRepository struct {
	db *sql.DB
}

// NewOwnerUserRepository creates a new PostgreSQL owner user repository
func NewOwnerUserRepository(db *sql.DB) repository.OwnerUserRepository {
	return &ownerUserRepository{db: db}
}

func (r *ownerUserRepository) Create(ctx context.Context, user *domain.OwnerUser) error {
	query := `
		INSERT INTO owner_users (
			id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			business_verified, business_verification_status,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
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
		user.BusinessVerified,
		user.BusinessVerificationStatus,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create owner user: %w", err)
	}

	return nil
}

func (r *ownerUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.OwnerUser, error) {
	query := `
		SELECT id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			business_verified, business_verification_status,
			created_at, updated_at
		FROM owner_users
		WHERE id = $1
	`

	user := &domain.OwnerUser{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.EmailVerified,
		&user.EmailVerificationToken,
		&user.EmailVerificationExpiresAt,
		&user.PasswordResetToken,
		&user.PasswordResetExpiresAt,
		&user.BusinessVerified,
		&user.BusinessVerificationStatus,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("owner user not found")
		}
		return nil, fmt.Errorf("failed to find owner user: %w", err)
	}

	return user, nil
}

func (r *ownerUserRepository) FindByEmail(ctx context.Context, email string) (*domain.OwnerUser, error) {
	query := `
		SELECT id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			business_verified, business_verification_status,
			created_at, updated_at
		FROM owner_users
		WHERE email = $1
	`

	user := &domain.OwnerUser{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.EmailVerified,
		&user.EmailVerificationToken,
		&user.EmailVerificationExpiresAt,
		&user.PasswordResetToken,
		&user.PasswordResetExpiresAt,
		&user.BusinessVerified,
		&user.BusinessVerificationStatus,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found is not an error for FindByEmail
		}
		return nil, fmt.Errorf("failed to find owner user by email: %w", err)
	}

	return user, nil
}

func (r *ownerUserRepository) FindByVerificationToken(ctx context.Context, token string) (*domain.OwnerUser, error) {
	query := `
		SELECT id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			business_verified, business_verification_status,
			created_at, updated_at
		FROM owner_users
		WHERE email_verification_token = $1
	`

	user := &domain.OwnerUser{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.EmailVerified,
		&user.EmailVerificationToken,
		&user.EmailVerificationExpiresAt,
		&user.PasswordResetToken,
		&user.PasswordResetExpiresAt,
		&user.BusinessVerified,
		&user.BusinessVerificationStatus,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("verification token not found")
		}
		return nil, fmt.Errorf("failed to find owner user by verification token: %w", err)
	}

	return user, nil
}

func (r *ownerUserRepository) FindByResetToken(ctx context.Context, token string) (*domain.OwnerUser, error) {
	query := `
		SELECT id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			business_verified, business_verification_status,
			created_at, updated_at
		FROM owner_users
		WHERE password_reset_token = $1
	`

	user := &domain.OwnerUser{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.EmailVerified,
		&user.EmailVerificationToken,
		&user.EmailVerificationExpiresAt,
		&user.PasswordResetToken,
		&user.PasswordResetExpiresAt,
		&user.BusinessVerified,
		&user.BusinessVerificationStatus,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("reset token not found")
		}
		return nil, fmt.Errorf("failed to find owner user by reset token: %w", err)
	}

	return user, nil
}

func (r *ownerUserRepository) Update(ctx context.Context, user *domain.OwnerUser) error {
	query := `
		UPDATE owner_users
		SET email = $2,
			password_hash = $3,
			email_verified = $4,
			email_verification_token = $5,
			email_verification_expires_at = $6,
			password_reset_token = $7,
			password_reset_expires_at = $8,
			business_verified = $9,
			business_verification_status = $10,
			updated_at = $11
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
		user.BusinessVerified,
		user.BusinessVerificationStatus,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update owner user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("owner user not found")
	}

	return nil
}

func (r *ownerUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM owner_users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete owner user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("owner user not found")
	}

	return nil
}

func (r *ownerUserRepository) ListPendingVerifications(ctx context.Context, limit, offset int) ([]*domain.OwnerUser, error) {
	query := `
		SELECT id, email, password_hash, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			business_verified, business_verification_status,
			created_at, updated_at
		FROM owner_users
		WHERE business_verification_status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list pending verifications: %w", err)
	}
	defer rows.Close()

	var users []*domain.OwnerUser
	for rows.Next() {
		user := &domain.OwnerUser{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.EmailVerified,
			&user.EmailVerificationToken,
			&user.EmailVerificationExpiresAt,
			&user.PasswordResetToken,
			&user.PasswordResetExpiresAt,
			&user.BusinessVerified,
			&user.BusinessVerificationStatus,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan owner user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}
