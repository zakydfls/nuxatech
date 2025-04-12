package route

import (
	"nuxatech-nextmedis/handler"
	"nuxatech-nextmedis/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	userHandler handler.UserHandler,
	authHandler handler.AuthHandler,
	productHandler handler.ProductHandler,
	cartHandler handler.CartHandler,
	accountHandler handler.AccountHandler,
	orderHandler handler.OrderHandler,
) *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/api/v1")

	// Submission TASK 1
	user := v1.Group("/user")
	user.POST("/create", userHandler.CreateUser)
	user.GET("/me", middleware.AuthMiddleware(), userHandler.GetUser)
	user.GET("/find", middleware.AuthMiddleware(), userHandler.FindUser)

	auth := v1.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/register", authHandler.Register)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.DELETE("/logout", authHandler.Logout)

	product := v1.Group("/product")
	product.GET("/", productHandler.GetAllProducts)
	product.GET("/:id", productHandler.GetProduct)
	product.POST("/", productHandler.CreateProduct)

	cart := v1.Group("cart")
	cart.POST("/add", middleware.AuthMiddleware(), cartHandler.AddToCart)
	cart.GET("/", middleware.AuthMiddleware(), cartHandler.GetCart)
	cart.PUT("/item/:id", middleware.AuthMiddleware(), cartHandler.UpdateCartItem)
	cart.DELETE("/item/:id", middleware.AuthMiddleware(), cartHandler.RemoveFromCart)

	order := v1.Group("order")
	order.POST("/", middleware.AuthMiddleware(), orderHandler.CreateOrder)
	order.GET("/:id", middleware.AuthMiddleware(), orderHandler.GetOrder)
	order.PUT("/:id/status", middleware.AuthMiddleware(), orderHandler.UpdateOrderStatus)
	order.GET("/", middleware.AuthMiddleware(), orderHandler.GetUserOrders)

	// Submission TASK 3
	user.POST("/wallet", middleware.AuthMiddleware(), accountHandler.CreateAccount)
	user.POST("/wallet/:id/deposit", middleware.AuthMiddleware(), accountHandler.Deposit)
	user.POST("/wallet/:id/withdraw", middleware.AuthMiddleware(), accountHandler.Withdraw)
	user.GET("/wallet/:id", middleware.AuthMiddleware(), accountHandler.GetAccount)

	return router
}
