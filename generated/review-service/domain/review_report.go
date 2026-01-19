package domain

import (
	"github.com/google/uuid"
	"time"
)

// ReviewReport represents ReviewReport
type ReviewReport struct {
	Id uuid.UUID `db:"id" json:"id"`
	ReviewId uuid.UUID `db:"review_id" json:"review_id"`
	ReporterId uuid.UUID `db:"reporter_id" json:"reporter_id"`
	Reason ReportReason `db:"reason" json:"reason"`
	Description *text `db:"description" json:"description,omitempty"`
	Status ReportStatus `db:"status" json:"status"`
	ResolvedBy *uuid.UUID `db:"resolved_by" json:"resolved_by,omitempty"`
	ResolvedAt *time.Time `db:"resolved_at" json:"resolved_at,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewReviewReport creates a new ReviewReport instance
func NewReviewReport() *ReviewReport {
	return &ReviewReport{}
}
