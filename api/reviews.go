package api

import (
	"strconv"
	"time"

	db "github.com/cihanalici/api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type reviewRequest struct {
	ProductID int32  `json:"product_id" binding:"required"`
	Rating    int32  `json:"rating" binding:"required"`
	Comment   string `json:"comment" binding:"required"`
}

type reviewResponse struct {
	ID        int32     `json:"id"`
	ProductID int32     `json:"product_id"`
	UserID    int32     `json:"user_id"`
	Rating    int32     `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func reviewNotation(review db.Review) reviewResponse {
	return reviewResponse{
		ID:        review.ID,
		ProductID: review.ProductID,
		UserID:    review.UserID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}
}

func reviewsNotation(reviews []db.Review) []reviewResponse {
	result := make([]reviewResponse, len(reviews))

	for i, review := range reviews {
		result[i] = reviewNotation(review)
	}

	return result
}

// CreateReview godoc
// @Summary Create a new review
// @Tags reviews
// @Description create a new review
// @Accept  json
// @Produce  json
// @Param input body reviewRequest true "Review Request"
// @Success 200 {object} reviewResponse
// @Router /reviews [post]

func (server *Server) createReview(ctx *gin.Context) {
	var req reviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	// bu alınan tokeni int32'ye çevirip userId'e atamamız gerekiyor
	userIdByToken, _ := ctx.Get("userId")
	userId, _ := userIdByToken.(int32)

	review, err := server.store.CreateReview(ctx, db.CreateReviewParams{
		ProductID: req.ProductID,
		UserID:    userId,
		Rating:    req.Rating,
		Comment:   req.Comment,
	})
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, reviewNotation(review))
}

// getReview godoc
// @Summary Get a review
// @Description Get a review
// @ID get-review
// @Accept  json
// @Produce  json
// @Param id path int true "Review ID"

type getReviewRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getReview(ctx *gin.Context) {
	var req getReviewRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	review, err := server.store.GetReviewById(ctx, req.ID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, reviewNotation(review))
}

// deleteReview godoc
// @Summary Delete a review
// @Tags reviews
// @Description delete a review
// @ID delete-review
// @Accept  json
// @Produce  json
// @Param id path int true "Review ID"
// @Success 200

func (server *Server) deleteReview(ctx *gin.Context) {
	reviewId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	err = server.store.DeleteReview(ctx, int32(reviewId))
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, gin.H{"status": "ok"})
}

// listReviews godoc
// @Summary List reviews
// @Description List reviews
// @ID list-reviews
// @Accept  json
// @Produce  json
// @Param page_id query int true "Page ID"
// @Param page_size query int true "Page Size"
// @Success 200 {array} reviewResponse

type listReviewsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listReviews(ctx *gin.Context) {
	var req listReviewsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.ListReviewsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	reviews, err := server.store.ListReviews(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, reviewsNotation(reviews))
}

// getReviewsByProductId godoc
// @Summary Get reviews by product ID
// @Description Get reviews by product ID
// @ID get-reviews-by-product-id
// @Accept  json
// @Produce  json
// @Param product_id query int true "Product ID"
// @Success 200 {array} reviewResponse

type getReviewsByProductIdRequest struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (server *Server) getReviewsByProductId(ctx *gin.Context) {
	productId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	_, err = server.store.GetProductById(ctx, int32(productId))
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	var req getReviewsByProductIdRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.GetReviewsByProductIdParams{
		ProductID: int32(productId),
		Limit:     req.Limit,
		Offset:    req.Offset,
	}

	reviews, err := server.store.GetReviewsByProductId(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, reviewsNotation(reviews))
}
