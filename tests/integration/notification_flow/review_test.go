package notification_flow

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/makoto-developer/go_microservice_example/tests/integration/notification_flow/clients"
)

// TestCreateReview tests review creation flow
func TestCreateReview(t *testing.T) {
	// Connect to review database
	reviewClient, err := clients.NewReviewClient(
		"postgresql://postgres:postgres_password@localhost:22018/review_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer reviewClient.Close()

	// Test data
	customerID := uuid.New()
	productID := uuid.New()
	orderID := uuid.New()
	rating := 5
	reviewText := "Excellent product! Highly recommended."

	// 1. Create review
	review, err := reviewClient.CreateReview(customerID, productID, orderID, rating, reviewText)
	require.NoError(t, err)
	assert.NotNil(t, review)
	assert.Equal(t, customerID, review.CustomerID)
	assert.Equal(t, productID, review.ProductID)
	assert.Equal(t, orderID, review.OrderID)
	assert.Equal(t, rating, review.Rating)
	assert.Equal(t, reviewText, review.ReviewText)

	// 2. Verify review is created
	retrievedReview, err := reviewClient.GetReview(review.ID)
	require.NoError(t, err)
	assert.Equal(t, review.ID, retrievedReview.ID)
	assert.Equal(t, rating, retrievedReview.Rating)

	// Cleanup
	err = reviewClient.DeleteReview(review.ID)
	require.NoError(t, err)
}

// TestReviewRatingValidation tests rating validation (1-5)
func TestReviewRatingValidation(t *testing.T) {
	reviewClient, err := clients.NewReviewClient(
		"postgresql://postgres:postgres_password@localhost:22018/review_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer reviewClient.Close()

	customerID := uuid.New()
	productID := uuid.New()
	orderID := uuid.New()

	// Test invalid rating (too low)
	_, err = reviewClient.CreateReview(customerID, productID, orderID, 0, "Test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rating must be between 1 and 5")

	// Test invalid rating (too high)
	_, err = reviewClient.CreateReview(customerID, productID, orderID, 6, "Test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rating must be between 1 and 5")

	// Test valid ratings
	for rating := 1; rating <= 5; rating++ {
		review, err := reviewClient.CreateReview(
			customerID,
			productID,
			uuid.New(), // Different order ID for each review
			rating,
			"Test review",
		)
		require.NoError(t, err)
		assert.Equal(t, rating, review.Rating)

		// Cleanup
		err = reviewClient.DeleteReview(review.ID)
		require.NoError(t, err)
	}
}

// TestGetReviewsByProduct tests fetching reviews by product ID
func TestGetReviewsByProduct(t *testing.T) {
	reviewClient, err := clients.NewReviewClient(
		"postgresql://postgres:postgres_password@localhost:22018/review_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer reviewClient.Close()

	productID := uuid.New()
	customerID := uuid.New()

	// Create multiple reviews for the same product
	reviewIDs := []uuid.UUID{}
	for i := 1; i <= 3; i++ {
		review, err := reviewClient.CreateReview(
			customerID,
			productID,
			uuid.New(),
			i+2, // ratings 3, 4, 5
			"Test review",
		)
		require.NoError(t, err)
		reviewIDs = append(reviewIDs, review.ID)
	}

	// Get all reviews for the product
	reviews, err := reviewClient.GetReviewsByProduct(productID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(reviews), 3)

	// Verify reviews are ordered by created_at DESC
	if len(reviews) >= 2 {
		assert.True(t, reviews[0].CreatedAt.After(reviews[1].CreatedAt) ||
			reviews[0].CreatedAt.Equal(reviews[1].CreatedAt))
	}

	// Cleanup
	for _, reviewID := range reviewIDs {
		err = reviewClient.DeleteReview(reviewID)
		require.NoError(t, err)
	}
}

// TestAverageRating tests average rating calculation
func TestAverageRating(t *testing.T) {
	reviewClient, err := clients.NewReviewClient(
		"postgresql://postgres:postgres_password@localhost:22018/review_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer reviewClient.Close()

	productID := uuid.New()
	customerID := uuid.New()

	// Create reviews with ratings: 3, 4, 5 (average = 4.0)
	reviewIDs := []uuid.UUID{}
	ratings := []int{3, 4, 5}
	for _, rating := range ratings {
		review, err := reviewClient.CreateReview(
			customerID,
			productID,
			uuid.New(),
			rating,
			"Test review",
		)
		require.NoError(t, err)
		reviewIDs = append(reviewIDs, review.ID)
	}

	// Calculate average rating
	avgRating, err := reviewClient.GetAverageRating(productID)
	require.NoError(t, err)
	assert.InDelta(t, 4.0, avgRating, 0.1) // Allow small floating point error

	// Cleanup
	for _, reviewID := range reviewIDs {
		err = reviewClient.DeleteReview(reviewID)
		require.NoError(t, err)
	}
}

// TestUpdateReview tests review update functionality
func TestUpdateReview(t *testing.T) {
	reviewClient, err := clients.NewReviewClient(
		"postgresql://postgres:postgres_password@localhost:22018/review_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer reviewClient.Close()

	// 1. Create initial review
	customerID := uuid.New()
	productID := uuid.New()
	orderID := uuid.New()
	review, err := reviewClient.CreateReview(customerID, productID, orderID, 3, "Initial review")
	require.NoError(t, err)

	// 2. Update review
	newRating := 5
	newText := "Updated review - much better experience!"
	err = reviewClient.UpdateReview(review.ID, newRating, newText)
	require.NoError(t, err)

	// 3. Verify update
	updatedReview, err := reviewClient.GetReview(review.ID)
	require.NoError(t, err)
	assert.Equal(t, newRating, updatedReview.Rating)
	assert.Equal(t, newText, updatedReview.ReviewText)
	assert.True(t, updatedReview.UpdatedAt.After(updatedReview.CreatedAt))

	// Cleanup
	err = reviewClient.DeleteReview(review.ID)
	require.NoError(t, err)
}

// TestReviewEditableUntil tests that reviews have an editable_until timestamp
func TestReviewEditableUntil(t *testing.T) {
	reviewClient, err := clients.NewReviewClient(
		"postgresql://postgres:postgres_password@localhost:22018/review_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer reviewClient.Close()

	// Create review
	review, err := reviewClient.CreateReview(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		4,
		"Test review",
	)
	require.NoError(t, err)

	// Verify editable_until is set (should be 30 days from now)
	assert.NotZero(t, review.EditableUntil)
	assert.True(t, review.EditableUntil.After(review.CreatedAt))

	// Cleanup
	err = reviewClient.DeleteReview(review.ID)
	require.NoError(t, err)
}

// TestDuplicateReviewPrevention tests that duplicate order+product reviews are prevented
func TestDuplicateReviewPrevention(t *testing.T) {
	reviewClient, err := clients.NewReviewClient(
		"postgresql://postgres:postgres_password@localhost:22018/review_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer reviewClient.Close()

	customerID := uuid.New()
	productID := uuid.New()
	orderID := uuid.New()

	// Create first review
	review1, err := reviewClient.CreateReview(customerID, productID, orderID, 4, "First review")
	require.NoError(t, err)

	// Try to create duplicate review (same order + product)
	_, err = reviewClient.CreateReview(customerID, productID, orderID, 5, "Duplicate review")
	assert.Error(t, err) // Should fail due to unique constraint

	// Cleanup
	err = reviewClient.DeleteReview(review1.ID)
	require.NoError(t, err)
}
