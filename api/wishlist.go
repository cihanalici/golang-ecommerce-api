package api

import (
	"strconv"
	"time"

	db "github.com/cihanalici/api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type WishlistRequest struct {
	UserID    int32 `json:"user_id"`
	ProductID int32 `json:"product_id"`
}

type WishlistResponse struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	ProductID int32     `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func WishlistNotation(wishlist db.Wishlist) WishlistResponse {

	return WishlistResponse{
		ID:        wishlist.ID,
		UserID:    wishlist.UserID,
		ProductID: wishlist.ProductID,
		CreatedAt: wishlist.CreatedAt,
		UpdatedAt: wishlist.UpdatedAt,
	}
}

func WishlistsNotation(wishlists []db.Wishlist) []WishlistResponse {
	result := make([]WishlistResponse, len(wishlists))

	for i, wishlist := range wishlists {
		result[i] = WishlistNotation(wishlist)
	}

	return result
}

// CreateWishlist godoc
// @Summary Create a wishlist item
// @Description create a wishlist item
// @ID create-wishlist-item
// @Accept  json
// @Produce  json
// @Param input body WishlistRequest true "Wishlist item to create"
// @Success 200 {object} WishlistResponse
// @Router /wishlist [post]

func (server *Server) createWishlist(ctx *gin.Context) {
	var req WishlistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	userId, _ := ctx.Get("userId")
	userIDInt32, _ := userId.(int32)

	arg := db.CreateWishlistItemParams{
		UserID:    userIDInt32,
		ProductID: req.ProductID,
	}

	wishlist, err := server.store.CreateWishlistItem(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, WishlistNotation(wishlist))
}

// getWishlist godoc
// @Summary Get a wishlist item
// @Description Get a wishlist item
// @ID get-wishlist-item
// @Accept  json
// @Produce  json
// @Param id path int true "Wishlist item ID"
// @Success 200 {object} WishlistResponse

type getWishlistRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getWishlist(ctx *gin.Context) {
	var req getWishlistRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	wishlist, err := server.store.GetWishlistItemById(ctx, req.ID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, WishlistNotation(wishlist))
}

// listWishlist godoc
// @Summary List wishlist items
// @Description List wishlist items
// @ID list-wishlist-items
// @Accept  json
// @Produce  json
// @Success 200 {array} WishlistResponse

type listWishlistRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listWishlist(ctx *gin.Context) {
	var req listWishlistRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.ListWishlistItemsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	wishlists, err := server.store.ListWishlistItems(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, WishlistsNotation(wishlists))
}

// getWishlistByUser godoc
// @Summary Get wishlist items by user
// @Description Get wishlist items by user
// @ID get-wishlist-items-by-user
// @Accept  json
// @Produce  json
// @Param user_id query int false "User ID"
// @Success 200 {array} WishlistResponse

type getWishlistByUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) getWishlistByUser(ctx *gin.Context) {
	var req getWishlistByUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	// bu alınan tokeni int32'ye çevirip userId'e atamamız gerekiyor
	userIdByToken, _ := ctx.Get("userId")
	userId, _ := userIdByToken.(int32)

	arg := db.GetWishlistItemsByUserIdParams{
		UserID: userId,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	wishlists, err := server.store.GetWishlistItemsByUserId(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, WishlistsNotation(wishlists))
}

// deleteWishlist godoc
// @Summary Delete a wishlist item
// @Description delete a wishlist item
// @ID delete-wishlist-item
// @Accept  json
// @Produce  json
// @Param id path int true "Wishlist item ID"

func (server *Server) deleteWishlist(ctx *gin.Context) {
	wishlistId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	err = server.store.DeleteWishlistItem(ctx, int32(wishlistId))
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, gin.H{"status": "ok"})
}
