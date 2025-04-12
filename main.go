package main

import (
	"net/http"
	"nuxatech-nextmedis/config"
	"nuxatech-nextmedis/handler"
	"nuxatech-nextmedis/middleware"
	"nuxatech-nextmedis/repository"
	"nuxatech-nextmedis/route"
	"nuxatech-nextmedis/service"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/basic/docs"
)

// @title           Nextmedis API
// @version         1.0
// @description     API Server for Nextmedis Application
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9000
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	docs.SwaggerInfo.Title = "Nextmedis API"
	docs.SwaggerInfo.Description = "API Documentation for Nextmedis Application"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:9000"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	config.DBInit()
	userRepository := repository.NewUserRepository()
	tokenRepository := repository.NewPersonalTokenRepository()
	productRepository := repository.NewProductRepository()
	cartRepository := repository.NewCartRepository()
	accountRepository := repository.NewAccountRepository()
	orderRepository := repository.NewOrderRepository()
	transactionRepository := repository.NewTransactionRepository()

	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService(userRepository, tokenRepository)
	productService := service.NewProductService(productRepository)
	cartService := service.NewCartService(cartRepository, productRepository)
	accountService := service.NewAccountService(accountRepository, transactionRepository)
	orderService := service.NewOrderService(orderRepository, cartRepository, productRepository, accountRepository)

	userHandler := handler.NewUserHandler(userService)
	authHadler := handler.NewAuthHandler(authService)
	productHandler := handler.NewProductHandler(productService)
	cartHandler := handler.NewCartHandler(cartService)
	accountHandler := handler.NewAccountHandler(accountService)
	orderHandler := handler.NewOrderHandler(orderService)

	middleware.SetAuthService(authService)

	server := route.SetupRoutes(
		userHandler,
		authHadler,
		productHandler,
		cartHandler,
		accountHandler,
		orderHandler,
	)
	server.LoadHTMLGlob("./public/html/*")
	server.Static("/public", "./public")
	server.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"apiVersion": "v1",
			"goVersion":  runtime.Version(),
		})
	})
	server.Run(":" + config.Envs.Port)
}
