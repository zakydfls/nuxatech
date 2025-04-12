package handler

import (
	"net/http"
	"nuxatech-nextmedis/dto/request"
	"nuxatech-nextmedis/dto/response"
	"nuxatech-nextmedis/service"
	"nuxatech-nextmedis/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	CreateUser(c *gin.Context)
	GetUser(c *gin.Context)
	FindUser(c *gin.Context)
	// UpdateUser(c *gin.Context)
	FindUserByEmail(c *gin.Context)
	FindByUsername(c *gin.Context)
	FindUserAfterDate(c *gin.Context)
	// DeleteUser(c *gin.Context)
}

type userHandler struct {
	userService service.UserService
}

// FindUser implements UserHandler.
func (h *userHandler) FindUser(c *gin.Context) {
	username := c.Query("username")
	email := c.Query("email")
	after := c.Query("after")
	if after != "" {
		user, err := h.userService.FindUserAfterDate(c, after)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.APIResponse{
				Success: false,
				Message: "Failed to get user data",
				Error:   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, response.APIResponse{
			Message: "Success get user data",
			Data:    user,
			Success: true,
		})
	} else if email != "" {
		user, err := h.userService.FindUserByEmail(c, email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.APIResponse{
				Success: false,
				Message: "Failed to get user data",
				Error:   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, response.APIResponse{
			Message: "Success get user data",
			Data:    user,
			Success: true,
		})
	} else if username != "" {
		user, err := h.userService.FindByUsername(c, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.APIResponse{
				Success: false,
				Message: "Failed to get user data",
				Error:   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, response.APIResponse{
			Message: "Success get user data",
			Data:    user,
			Success: true,
		})
	} else {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Message: "Missing query params",
			Success: false,
		})
	}
}

// FindByUsername implements UserHandler.
func (h *userHandler) FindByUsername(c *gin.Context) {
	username := c.Query("username")
	user, err := h.userService.FindByUsername(c, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get user data",
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.APIResponse{
		Message: "Success get user data",
		Data:    user,
		Success: true,
	})
}

// FindUserByEmail implements UserHandler.
func (h *userHandler) FindUserByEmail(c *gin.Context) {
	email := c.Query("email")
	user, err := h.userService.FindByUsername(c, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get user data",
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.APIResponse{
		Message: "Success get user data",
		Data:    user,
		Success: true,
	})
}

// FindUserAfterDate implements UserHandler.
func (h *userHandler) FindUserAfterDate(c *gin.Context) {
	after := c.Query("after")
	user, err := h.userService.FindUserAfterDate(c, after)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get user data",
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.APIResponse{
		Message: "Success get user data",
		Data:    user,
		Success: true,
	})
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

// @Summary Get user profile
// @Description Get current user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} response.APIResponse{data=response.UserResponse} "Profile retrieved"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Router /users/profile [get]
// @Security BearerAuth
func (h *userHandler) GetUser(c *gin.Context) {
	userId := utils.GetUserID(c)
	user, err := h.userService.GetUser(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get user data",
			Error:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.APIResponse{
		Message: "Success get user data",
		Data:    user,
		Success: true,
	})
}

func (h *userHandler) CreateUser(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to create user",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "User created successfully",
	})
}
