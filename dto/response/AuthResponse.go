package response

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	User  UserResponse  `json:"user"`
	Token TokenResponse `json:"token"`
}

type RegisterResponse struct {
	User UserResponse `json:"user"`
}
