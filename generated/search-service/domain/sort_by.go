package domain

// SortBy represents SortBy type
type SortBy string

const (
	SortByRelevance SortBy = "RELEVANCE"
	SortByPriceAsc SortBy = "PRICE_ASC"
	SortByPriceDesc SortBy = "PRICE_DESC"
	SortByRatingDesc SortBy = "RATING_DESC"
	SortByNewest SortBy = "NEWEST"
	SortByReviewCountDesc SortBy = "REVIEW_COUNT_DESC"
)

// SortByValues returns all possible values
func SortByValues() []SortBy {
	return []SortBy{
		SortByRelevance,
		SortByPriceAsc,
		SortByPriceDesc,
		SortByRatingDesc,
		SortByNewest,
		SortByReviewCountDesc,
	}
}

// IsValid checks if the value is valid
func (e SortBy) IsValid() bool {
	switch e {
	case SortByRelevance:
	case SortByPriceAsc:
	case SortByPriceDesc:
	case SortByRatingDesc:
	case SortByNewest:
	case SortByReviewCountDesc:
		return true
	}
	return false
}
