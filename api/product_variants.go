package api

import (
	"strconv"
	"time"

	db "github.com/cihanalici/api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type productVariantRequest struct {
	ProductID int32  `json:"product_id"`
	Color     string `json:"color" binding:"required"`
	Size      string `json:"size" binding:"required"`
	Stock     int32  `json:"stock" binding:"required"`
	Price     string `json:"price" binding:"required"`
}

type productVariantResponse struct {
	ID        int32     `json:"id"`
	ProductID int32     `json:"product_id"`
	Color     string    `json:"color"`
	Size      string    `json:"size"`
	Stock     int32     `json:"stock"`
	Price     string    `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func productVariantNotation(productVariant db.ProductVariant) productVariantResponse {

	return productVariantResponse{
		ID:        productVariant.ID,
		ProductID: productVariant.ProductID,
		Color:     productVariant.Color,
		Size:      productVariant.Size,
		Stock:     productVariant.Stock,
		Price:     productVariant.Price,
		CreatedAt: productVariant.CreatedAt,
		UpdatedAt: productVariant.UpdatedAt,
	}
}

func productVariantsNotation(productVariants []db.ProductVariant) []productVariantResponse {
	result := make([]productVariantResponse, len(productVariants))

	for i, productVariant := range productVariants {
		result[i] = productVariantNotation(productVariant)
	}

	return result
}

// CreateProductVariant godoc
// @Summary Create a product variant
// @Description Create a product variant
// @ID create-product-variant
// @Accept  json
// @Produce  json
// @Tags product_variants
// @Param request body productVariantRequest true "Product variant data"
// @Success 200 {object} productVariantResponse
// @Router /product_variants [post]

func (server *Server) createProductVariant(ctx *gin.Context) {
	var req productVariantRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.CreateProductVariantParams{
		ProductID: req.ProductID,
		Color:     req.Color,
		Size:      req.Size,
		Stock:     req.Stock,
		Price:     req.Price,
	}

	// if req.ProductID != nil {
	// 	arg.ProductID = sql.NullInt32{Int32: *req.ProductID, Valid: true}
	// }

	productVariant, err := server.store.CreateProductVariant(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, productVariantNotation(productVariant))
}

// GetProductVariant godoc
// @Summary Get a product variant
// @Description Get a product variant by ID
// @ID get-product-variant
// @Accept  json
// @Produce  json
// @Tags product_variants
// @Param id path int true "Product Variant ID"
// @Success 200 {object} productVariantResponse

type getProductVariantRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getProductVariant(ctx *gin.Context) {
	var req getProductVariantRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	productVariant, err := server.store.GetProductVariantById(ctx, req.ID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, productVariantNotation(productVariant))
}

// ListProductVariants godoc
// @Summary List product variants
// @Description List all product variants
// @ID list-product-variants
// @Accept  json
// @Produce  json
// @Tags product_variants
// @Success 200 {array} productVariantResponse
// @Router /product_variants [get]

type listProductVariantsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listProductVariants(ctx *gin.Context) {
	var req listProductVariantsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	productVariants, err := server.store.ListProductVariants(ctx, db.ListProductVariantsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, productVariantsNotation(productVariants))
}

// UpdateProductVariant godoc
// @Summary Update a product variant
// @Description Update a product variant by ID
// @ID update-product-variant
// @Accept  json
// @Produce  json
// @Tags product_variants

func (server *Server) updateProductVariant(ctx *gin.Context) {
	variantId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	variant, err := server.store.GetProductVariantById(ctx, int32(variantId))
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	var req productVariantRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.UpdateProductVariantParams{
		ID:        variant.ID,
		Color:     req.Color,
		Size:      req.Size,
		Stock:     req.Stock,
		Price:     req.Price,
		ProductID: req.ProductID,
	}

	variant, err = server.store.UpdateProductVariant(ctx, arg)

	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, productVariantNotation(variant))
}

// DeleteProductVariant godoc
// @Summary Delete a product variant
// @Description Delete a product variant by ID
// @ID delete-product-variant
// @Accept  json
// @Produce  json
// @Tags product_variants
// @Param id path int true "Product Variant ID"

func (server *Server) deleteProductVariant(ctx *gin.Context) {
	variantId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	err = server.store.DeleteProductVariant(ctx, int32(variantId))

	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, gin.H{"status": "ok"})
}
