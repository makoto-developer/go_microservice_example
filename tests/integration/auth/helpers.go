package auth

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

const (
	// Test configuration
	AuthServiceAddr = "localhost:22100"
	DBHost          = "localhost"
	DBPort          = "22010"
	DBUser          = "postgres"
	DBPassword      = "postgres"
	DBName          = "auth_service"

	// JWT secret (should match auth service config)
	JWTSecret = "your-secret-key-here-change-in-production"
)

// DBHelper provides database utilities for testing
type DBHelper struct {
	db *sql.DB
}

// NewDBHelper creates a new database helper
func NewDBHelper(t *testing.T) *DBHelper {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPassword, DBName)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err, "Failed to connect to database")

	err = db.Ping()
	require.NoError(t, err, "Failed to ping database")

	return &DBHelper{db: db}
}

// Close closes the database connection
func (h *DBHelper) Close() error {
	return h.db.Close()
}

// CleanupTestData removes test data from the database
func (h *DBHelper) CleanupTestData(t *testing.T, email string) {
	// Delete from customer_users table
	_, err := h.db.Exec("DELETE FROM customer_users WHERE email = $1", email)
	if err != nil {
		t.Logf("Warning: failed to cleanup customer_users: %v", err)
	}

	// Delete from owner_users table
	_, err = h.db.Exec("DELETE FROM owner_users WHERE email = $1", email)
	if err != nil {
		t.Logf("Warning: failed to cleanup owner_users: %v", err)
	}
}

// CleanupAllTestData removes all test data (emails starting with "test_")
func (h *DBHelper) CleanupAllTestData(t *testing.T) {
	// Delete from customer_users table
	_, err := h.db.Exec("DELETE FROM customer_users WHERE email LIKE 'test_%'")
	if err != nil {
		t.Logf("Warning: failed to cleanup customer_users: %v", err)
	}

	// Delete from owner_users table
	_, err = h.db.Exec("DELETE FROM owner_users WHERE email LIKE 'test_%'")
	if err != nil {
		t.Logf("Warning: failed to cleanup owner_users: %v", err)
	}
}

// UserExistsInCustomer checks if a customer user exists in the database
func (h *DBHelper) UserExistsInCustomer(email string) (bool, error) {
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM customer_users WHERE email = $1)", email).Scan(&exists)
	return exists, err
}

// UserExistsInOwner checks if an owner user exists in the database
func (h *DBHelper) UserExistsInOwner(email string) (bool, error) {
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM owner_users WHERE email = $1)", email).Scan(&exists)
	return exists, err
}

// GetCustomerUserID gets the user ID for a customer
func (h *DBHelper) GetCustomerUserID(email string) (string, error) {
	var userID string
	err := h.db.QueryRow("SELECT user_id FROM customer_users WHERE email = $1", email).Scan(&userID)
	return userID, err
}

// GetOwnerUserID gets the user ID for an owner
func (h *DBHelper) GetOwnerUserID(email string) (string, error) {
	var userID string
	err := h.db.QueryRow("SELECT user_id FROM owner_users WHERE email = $1", email).Scan(&userID)
	return userID, err
}

// TokenHelper provides JWT token utilities
type TokenHelper struct {
	secret string
}

// NewTokenHelper creates a new token helper
func NewTokenHelper() *TokenHelper {
	return &TokenHelper{secret: JWTSecret}
}

// ValidateToken validates a JWT token and returns the claims
func (h *TokenHelper) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(h.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GetUserIDFromToken extracts user_id from token
func (h *TokenHelper) GetUserIDFromToken(tokenString string) (string, error) {
	claims, err := h.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in token")
	}

	return userID, nil
}

// IsTokenExpired checks if token is expired
func (h *TokenHelper) IsTokenExpired(tokenString string) (bool, error) {
	claims, err := h.ValidateToken(tokenString)
	if err != nil {
		return true, err
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return true, fmt.Errorf("exp not found in token")
	}

	expirationTime := time.Unix(int64(exp), 0)
	return time.Now().After(expirationTime), nil
}

// TestDataGenerator generates test data
type TestDataGenerator struct{}

// NewTestDataGenerator creates a new test data generator
func NewTestDataGenerator() *TestDataGenerator {
	return &TestDataGenerator{}
}

// GenerateTestEmail generates a unique test email
func (g *TestDataGenerator) GenerateTestEmail() string {
	return fmt.Sprintf("test_%s@example.com", uuid.New().String()[:8])
}

// GenerateTestPassword generates a test password
func (g *TestDataGenerator) GenerateTestPassword() string {
	return "Test@Password123"
}
