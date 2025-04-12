package handler

import (
	"net/http"
	"nuxatech-nextmedis/dto/request"
	"nuxatech-nextmedis/dto/response"
	"nuxatech-nextmedis/service"
	"nuxatech-nextmedis/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
	Logout(c *gin.Context)
}

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

// @Summary User registration
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "Registration details"
// @Success 201 {object} response.APIResponse{data=response.UserResponse} "Registration successful"
// @Failure 400 {object} response.APIResponse "Invalid request format"
// @Failure 409 {object} response.APIResponse "Email already exists"
// @Router /auth/register [post]
func (h *authHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	result, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Registration failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Registration successful",
		Data:    result,
	})
}

// @Summary User login
// @Description Authenticate user and get access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login credentials"
// @Success 200 {object} response.APIResponse{data=response.LoginResponse} "Login successful"
// @Failure 400 {object} response.APIResponse "Invalid request format"
// @Failure 401 {object} response.APIResponse "Invalid credentials"
// @Router /auth/login [post]
func (h *authHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Login failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    result,
	})
}

// @Summary Refresh token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer refresh_token"
// @Success 200 {object} response.APIResponse{data=response.TokenResponse} "Token refreshed"
// @Failure 401 {object} response.APIResponse "Invalid refresh token"
// @Router /auth/refresh [post]
// @Security BearerAuth
func (h *authHandler) RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Refresh token is required",
			Error:   "no token provided",
		})
		return
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Invalid token format",
			Error:   "invalid token format",
		})
		return
	}

	refreshToken := splitToken[1]
	result, err := h.authService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Token refresh failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data:    result,
	})
}

func (h *authHandler) Logout(c *gin.Context) {
	userID := utils.GetUserID(c)
	err := h.authService.Logout(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Logout failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}
