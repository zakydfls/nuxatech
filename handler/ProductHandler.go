package handler

import (
	"fmt"
	"net/http"
	"nuxatech-nextmedis/dto/request"
	"nuxatech-nextmedis/dto/response"
	"nuxatech-nextmedis/service"
	"nuxatech-nextmedis/utils"

	"github.com/gin-gonic/gin"
)

type ProductHandler interface {
	GetAllProducts(c *gin.Context)
	CreateProduct(c *gin.Context)
	GetProduct(c *gin.Context)
}

type productHandler struct {
	productService service.ProductService
}

// @Summary Get product by ID
// @Description Get detailed product information
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID" format(uuid)
// @Success 200 {object} response.APIResponse{data=response.ProductResponse} "Product retrieved"
// @Failure 404 {object} response.APIResponse "Product not found"
// @Router /product/{id} [get]
func (p *productHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	fmt.Println(id)
	product, err := p.productService.GetProduct(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.APIResponse{
			Success: false,
			Message: "Product not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Data:    product,
		Message: "Success to get product",
	})

}

// @Summary Create product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body request.CreateProductRequest true "Product details"
// @Success 201 {object} response.APIResponse{data=response.ProductResponse} "Product created"
// @Failure 400 {object} response.APIResponse "Invalid request"
// @Failure 401 {object} response.APIResponse "Unauthorized"
// @Router /products [post]
// @Security BearerAuth
func (p *productHandler) CreateProduct(c *gin.Context) {
	var req request.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	product, err := p.productService.CreateProduct(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to create product",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Success: true,
		Data:    product,
		Message: "Success to create product",
	})

}

// @Summary Get products
// @Description Get paginated list of products
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Param search query string false "Search term"
// @Success 200 {object} response.APIResponse{data=response.ProductPagingResponse} "Products retrieved"
// @Failure 400 {object} response.APIResponse "Invalid parameters"
// @Router /products [get]
func (p *productHandler) GetAllProducts(c *gin.Context) {
	params := service.ProductQueryParams{
		Page:   utils.ParseIntWithDefault(c.Query("page"), 1),
		Limit:  utils.ParseIntWithDefault(c.Query("limit"), 10),
		Search: c.Query("search"),
	}

	products, err := p.productService.GetAllProducts(c, params)
	if len(products.Result) == 0 {
		c.JSON(http.StatusNotFound, response.APIResponse{
			Success: false,
			Message: "No product found",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Message: "Failed to get product data",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Data:    products,
		Message: "Success to get product data",
	})

}

func NewProductHandler(ps service.ProductService) ProductHandler {
	return &productHandler{
		productService: ps,
	}
}
