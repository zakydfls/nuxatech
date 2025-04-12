package handler

import (
	"net/http"
	"nuxatech-nextmedis/dto/request"
	"nuxatech-nextmedis/dto/response"
	"nuxatech-nextmedis/service"
	"nuxatech-nextmedis/utils"

	"github.com/gin-gonic/gin"
)

type CartHandler interface {
	AddToCart(c *gin.Context)
	GetCart(c *gin.Context)
	UpdateCartItem(c *gin.Context)
	RemoveFromCart(c *gin.Context)
}

type cartHandler struct {
	cartService service.CartService
}

func (h *cartHandler) AddToCart(c *gin.Context) {
	var req request.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	userID := utils.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   "user not authenticated",
		})
		return
	}

	cart, err := h.cartService.AddToCart(c, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Failed to add to cart",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Successfully added to cart",
		Data:    cart,
	})
}

func (h *cartHandler) GetCart(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   "user not authenticated",
		})
		return
	}

	cart, err := h.cartService.GetCart(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get cart",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Successfully retrieved cart",
		Data:    cart,
	})
}

func (h *cartHandler) UpdateCartItem(c *gin.Context) {
	var req request.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   "user not authenticated",
		})
		return
	}

	cartItemID := c.Param("id")
	if cartItemID == "" {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid cart item ID",
		})
		return
	}

	cart, err := h.cartService.UpdateCartItem(c, userID, cartItemID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Failed to update cart item",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Successfully updated cart item",
		Data:    cart,
	})
}

func (h *cartHandler) RemoveFromCart(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, response.APIResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   "user not authenticated",
		})
		return
	}

	cartItemID := c.Param("id")
	if cartItemID == "" {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid cart item ID",
		})
		return
	}

	err := h.cartService.RemoveFromCart(c, userID, cartItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Failed to remove item from cart",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Successfully removed item from cart",
	})
}

func NewCartHandler(cartService service.CartService) CartHandler {
	return &cartHandler{
		cartService: cartService,
	}
}
