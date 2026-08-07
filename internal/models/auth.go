package models

// LoginRequest defines the payload required for user authentication.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserInfo represents public details of the authenticated user.
type UserInfo struct {
	Username string `json:"username"`
}

// LoginResponse defines the successful response returned upon authentication.
type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}
