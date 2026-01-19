package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/auth-service/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			id, email, password_hash, role, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.Id,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.EmailVerified,
		user.EmailVerificationToken,
		user.EmailVerificationExpiresAt,
		user.PasswordResetToken,
		user.PasswordResetExpiresAt,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	var emailVerificationToken sql.NullString
	var emailVerificationExpiresAt sql.NullTime
	var passwordResetToken sql.NullString
	var passwordResetExpiresAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.EmailVerified,
		&emailVerificationToken,
		&emailVerificationExpiresAt,
		&passwordResetToken,
		&passwordResetExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	if emailVerificationToken.Valid {
		user.EmailVerificationToken = &emailVerificationToken.String
	}
	if emailVerificationExpiresAt.Valid {
		user.EmailVerificationExpiresAt = &emailVerificationExpiresAt.Time
	}
	if passwordResetToken.Valid {
		user.PasswordResetToken = &passwordResetToken.String
	}
	if passwordResetExpiresAt.Valid {
		user.PasswordResetExpiresAt = &passwordResetExpiresAt.Time
	}

	return &user, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	var emailVerificationToken sql.NullString
	var emailVerificationExpiresAt sql.NullTime
	var passwordResetToken sql.NullString
	var passwordResetExpiresAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.EmailVerified,
		&emailVerificationToken,
		&emailVerificationExpiresAt,
		&passwordResetToken,
		&passwordResetExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found with email: %s", email)
	}
	if err != nil {
		return nil, err
	}

	if emailVerificationToken.Valid {
		user.EmailVerificationToken = &emailVerificationToken.String
	}
	if emailVerificationExpiresAt.Valid {
		user.EmailVerificationExpiresAt = &emailVerificationExpiresAt.Time
	}
	if passwordResetToken.Valid {
		user.PasswordResetToken = &passwordResetToken.String
	}
	if passwordResetExpiresAt.Valid {
		user.PasswordResetExpiresAt = &passwordResetExpiresAt.Time
	}

	return &user, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1,
			password_hash = $2,
			role = $3,
			email_verified = $4,
			email_verification_token = $5,
			email_verification_expires_at = $6,
			password_reset_token = $7,
			password_reset_expires_at = $8,
			updated_at = $9
		WHERE id = $10
	`

	_, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.EmailVerified,
		user.EmailVerificationToken,
		user.EmailVerificationExpiresAt,
		user.PasswordResetToken,
		user.PasswordResetExpiresAt,
		time.Now(),
		user.Id,
	)

	return err
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, email_verified,
			email_verification_token, email_verification_expires_at,
			password_reset_token, password_reset_expires_at,
			created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		var emailVerificationToken sql.NullString
		var emailVerificationExpiresAt sql.NullTime
		var passwordResetToken sql.NullString
		var passwordResetExpiresAt sql.NullTime

		err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
			&user.EmailVerified,
			&emailVerificationToken,
			&emailVerificationExpiresAt,
			&passwordResetToken,
			&passwordResetExpiresAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if emailVerificationToken.Valid {
			user.EmailVerificationToken = &emailVerificationToken.String
		}
		if emailVerificationExpiresAt.Valid {
			user.EmailVerificationExpiresAt = &emailVerificationExpiresAt.Time
		}
		if passwordResetToken.Valid {
			user.PasswordResetToken = &passwordResetToken.String
		}
		if passwordResetExpiresAt.Valid {
			user.PasswordResetExpiresAt = &passwordResetExpiresAt.Time
		}

		users = append(users, &user)
	}

	return users, rows.Err()
}

// BeginTx starts a transaction
func (r *PostgresUserRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
