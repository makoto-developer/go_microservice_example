package domain


// JWT represents JWT value object
type JWT struct {
	AccessToken string `db:"access_token" json:"-"`
	RefreshToken string `db:"refresh_token" json:"-"`
	ExpiresIn int `db:"expires_in" json:"expires_in"`
}

// NewJWT creates a new JWT instance
func NewJWT() *JWT {
	return &JWT{}
}
