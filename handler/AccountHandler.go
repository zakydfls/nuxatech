package handler

import (
	"net/http"
	"nuxatech-nextmedis/dto/request"
	"nuxatech-nextmedis/dto/response"
	"nuxatech-nextmedis/service"

	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	CreateAccount(c *gin.Context)
	GetAccount(c *gin.Context)
	Deposit(c *gin.Context)
	Withdraw(c *gin.Context)
}

type accountHandler struct {
	accountService service.AccountService
}

func (h *accountHandler) CreateAccount(c *gin.Context) {
	var req request.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	account, err := h.accountService.CreateAccount(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to create account",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "Account created successfully",
		Data:    account,
	})
}

func (h *accountHandler) GetAccount(c *gin.Context) {
	id := c.Param("id")
	account, err := h.accountService.GetAccount(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.APIResponse{
			Success: false,
			Message: "Account not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Account retrieved successfully",
		Data:    account,
	})
}

func (h *accountHandler) Deposit(c *gin.Context) {
	id := c.Param("id")
	var req request.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	transaction, err := h.accountService.Deposit(c, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to process deposit",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Deposit processed successfully",
		Data:    transaction,
	})
}

func (h *accountHandler) Withdraw(c *gin.Context) {
	id := c.Param("id")
	var req request.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	transaction, err := h.accountService.Withdraw(c, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to process withdrawal",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Withdrawal processed successfully",
		Data:    transaction,
	})
}

func NewAccountHandler(accountService service.AccountService) AccountHandler {
	return &accountHandler{
		accountService: accountService,
	}
}
