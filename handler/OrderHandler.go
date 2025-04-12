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

type OrderHandler interface {
	CreateOrder(c *gin.Context)
	GetOrder(c *gin.Context)
	UpdateOrderStatus(c *gin.Context)
	GetUserOrders(c *gin.Context)
}

type orderHandler struct {
	orderService service.OrderService
}

// @Summary Create new order
// @Description Create a new order from cart items
// @Tags orders
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body request.CreateOrderRequest true "Order creation request"
// @Success 201 {object} response.APIResponse{data=response.OrderResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /orders [post]
// @Security BearerAuth
func (h *orderHandler) CreateOrder(c *gin.Context) {
	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
		return
	}

	userID := utils.GetUserID(c)
	order, err := h.orderService.CreateOrder(c, userID, &req)
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "insufficient stock") {
			status = http.StatusConflict
		}

		c.JSON(status, response.APIResponse{
			Success: false,
			Message: "Failed to create order",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Order created successfully",
		Data:    order,
	})
}

// @Summary Get order by ID
// @Description Get detailed information about a specific order
// @Tags orders
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Order ID" format(uuid)
// @Success 200 {object} response.APIResponse{data=response.OrderResponse} "Order retrieved successfully"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 404 {object} response.APIResponse "Order not found"
// @Router /orders/{id} [get]
// @Security BearerAuth
func (h *orderHandler) GetOrder(c *gin.Context) {
	userID := utils.GetUserID(c)
	orderID := c.Param("id")

	order, err := h.orderService.GetOrder(c, userID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.APIResponse{
			Success: false,
			Message: "Order not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Order retrieved successfully",
		Data:    order,
	})
}

// @Summary Update order status
// @Description Update the status of an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Order ID" format(uuid)
// @Param request body request.UpdateOrderStatusRequest true "Status update request"
// @Success 200 {object} response.APIResponse{data=response.OrderResponse} "Order status updated successfully"
// @Failure 400 {object} response.APIResponse "Invalid request"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 404 {object} response.APIResponse "Order not found"
// @Router /orders/{id}/status [put]
// @Security BearerAuth
func (h *orderHandler) UpdateOrderStatus(c *gin.Context) {
	var req request.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	userID := utils.GetUserID(c)
	orderID := c.Param("id")

	order, err := h.orderService.UpdateOrderStatus(c, userID, orderID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Failed to update order status",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Order status updated successfully",
		Data:    order,
	})
}

// @Summary Get user orders
// @Description Get paginated list of user orders
// @Tags orders
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Param search query string false "Search term"
// @Success 200 {object} response.APIResponse{data=response.OrderPagingResponse} "Orders retrieved successfully"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Failure 404 {object} response.APIResponse "No orders found"
// @Router /orders [get]
// @Security BearerAuth
func (h *orderHandler) GetUserOrders(c *gin.Context) {
	userID := utils.GetUserID(c)

	params := service.ProductQueryParams{
		Page:   utils.ParseIntWithDefault(c.Query("page"), 1),
		Limit:  utils.ParseIntWithDefault(c.Query("limit"), 10),
		Search: c.Query("search"),
	}

	orders, err := h.orderService.GetUserOrders(c, userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get orders",
			Error:   err.Error(),
		})
		return
	}

	if len(orders.Result) == 0 {
		c.JSON(http.StatusNotFound, response.APIResponse{
			Success: false,
			Message: "No orders found",
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Orders retrieved successfully",
		Data:    orders,
	})
}

func NewOrderHandler(orderService service.OrderService) OrderHandler {
	return &orderHandler{
		orderService: orderService,
	}
}
