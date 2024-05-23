package transport

// UserAuth is the struct that represents the user authentication request
type UserAuth struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// DeleteUserRequest is the struct that represents the delete user request
type DeleteUserRequest struct {
	Username string `json:"username,omitempty"`
}

// InvalidateToken is the struct that represents the invalidate token request
type InvalidateToken struct {
	Username string `json:"username,omitempty"`
}

// ValidateToken is the struct that represents the validate token response
type ValidateToken struct {
	Valid bool `json:"valid,omitempty"`
}

type GetTokenResponse struct {
	JwtToken string `json:"jwt_token"`
}
