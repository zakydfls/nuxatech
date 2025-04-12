package request

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20" validate:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=6" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20" validate:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=6" validate:"required,min=6"`
}
