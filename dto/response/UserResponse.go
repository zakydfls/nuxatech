package response

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}
