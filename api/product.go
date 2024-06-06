package api

import (
	"database/sql"
	"time"

	db "github.com/cihanalici/api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type productRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"` // Pointer type to allow nil value
	Price       string  `json:"price" binding:"required"`
	Stock       int32   `json:"stock" binding:"required"`
	CategoryID  *int32  `json:"category_id"` // Pointer type to allow nil value
}

type productResponse struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"` // Pointer type to allow nil value
	Price       string    `json:"price"`
	Stock       int32     `json:"stock"`
	CategoryID  *int32    `json:"category_id"` // Pointer type to allow nil value
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func productNotation(product db.Product) productResponse {
	var description *string
	if product.Description.Valid {
		description = &product.Description.String
	}

	var categoryID *int32
	if product.CategoryID.Valid {
		categoryID = &product.CategoryID.Int32
	}

	return productResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: description,
		Price:       product.Price,
		Stock:       product.Stock,
		CategoryID:  categoryID,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func productsNotation(products []db.Product) []productResponse {
	result := make([]productResponse, len(products))

	for i, product := range products {
		result[i] = productNotation(product)
	}

	return result
}

// CreateProduct godoc
// @Summary Create a product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param request body productRequest true "Product request"
// @Success 200 {object} productResponse
// @Router /products [post]

func (server *Server) createProduct(ctx *gin.Context) {
	var req productRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	// Handling null values for Description and CategoryID
	description := sql.NullString{Valid: false}
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}

	categoryID := sql.NullInt32{Valid: false}
	if req.CategoryID != nil {
		categoryID = sql.NullInt32{Int32: *req.CategoryID, Valid: true}
	}

	arg := db.CreateProductParams{
		Name:        req.Name,
		Description: description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  categoryID,
	}

	product, err := server.store.CreateProduct(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, productNotation(product))
}

// GetProduct godoc
// @Summary Get a product
// @Description Get a product by id
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} productResponse
// @Router /products/{id} [get]

type getProductRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getProduct(ctx *gin.Context) {
	var req getProductRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	product, err := server.store.GetProductById(ctx, req.ID)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, productNotation(product))
}

// GetProducts godoc
// @Summary Get all products
// @Description Get all products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} productResponse
// @Router /products [get]

type getProductsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) getProducts(ctx *gin.Context) {
	var req getProductsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	products, err := server.store.ListProducts(ctx, db.ListProductsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, productsNotation(products))
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update a product by id
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body productRequest true "Product request"
// @Success 200 {object} productResponse
// @Router /products/{id} [put]

type updateProductRequest struct {
	ID          int32          `uri:"id" binding:"required,min=1"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	Price       string         `json:"price"`
	Stock       int32          `json:"stock"`
	CategoryID  sql.NullInt32  `json:"category_id"`
}

func (server *Server) updateProduct(ctx *gin.Context) {
	var req updateProductRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	description := sql.NullString{Valid: false}
	if req.Description.Valid {
		description = sql.NullString{String: req.Description.String, Valid: true}
	}

	categoryID := sql.NullInt32{Valid: false}
	if req.CategoryID.Valid {
		categoryID = sql.NullInt32{Int32: req.CategoryID.Int32, Valid: true}
	}

	arg := db.UpdateProductParams{
		ID:          req.ID,
		Name:        req.Name,
		Description: description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  categoryID,
	}

	product, err := server.store.UpdateProduct(ctx, arg)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, productNotation(product))
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by id
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
func (server *Server) deleteProduct(ctx *gin.Context) {
	var req getProductRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	err := server.store.DeleteProduct(ctx, req.ID)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, gin.H{"status": "ok"})
}
