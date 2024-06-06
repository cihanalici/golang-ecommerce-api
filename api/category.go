package api

import (
	"strconv"
	"time"

	db "github.com/cihanalici/api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type categoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type categoryResponse struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func categoryNotation(category db.Category) categoryResponse {
	return categoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func categoriesNotation(categories []db.Category) []categoryResponse {
	result := make([]categoryResponse, len(categories))

	for i, category := range categories {
		result[i] = categoryNotation(category)
	}

	return result
}

// CreateCategory godoc
// @Summary Create a category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param request body categoryRequest true "Category request"
// @Success 200 {object} categoryResponse
// @Router /categories [post]
func (server *Server) createCategory(ctx *gin.Context) {
	var req categoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	category, err := server.store.CreateCategory(ctx, db.CreateCategoryParams{
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, categoryNotation(category))
}

// GetCategory godoc
// @Summary Get a category
// @Description Get a category by id
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} categoryResponse
// @Router /categories/{id} [get]

type getCategoryRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getCategory(ctx *gin.Context) {
	var req getCategoryRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	category, err := server.store.GetCategoryById(ctx, req.ID)

	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, categoryNotation(category))
}

// GetCategories godoc
// @Summary Get all categories
// @Description Get all categories
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} categoryResponse
// @Router /categories [get]

type getCategoriesRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) getCategories(ctx *gin.Context) {
	var req getCategoriesRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.ListCategoriesParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	categories, err := server.store.ListCategories(ctx, arg)

	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, categoriesNotation(categories))
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update a category by id
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body categoryRequest true "Category request"
// @Success 200 {object} categoryResponse
// @Router /categories/{id} [put]

type updateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (server *Server) updateCategory(ctx *gin.Context) {
	catId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	category, err := server.store.GetCategoryById(ctx, int32(catId))

	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	var req updateCategoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.UpdateCategoryParams{
		ID:          category.ID,
		Name:        req.Name,
		Description: req.Description,
	}

	category, err = server.store.UpdateCategory(ctx, arg)

	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, categoryNotation(category))
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by id
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200
// @Router /categories/{id} [delete]
func (server *Server) deleteCategory(ctx *gin.Context) {
	catId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	err = server.store.DeleteCategory(ctx, int32(catId))

	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, gin.H{"status": "ok"})
}
