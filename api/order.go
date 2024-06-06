package api

import (
	"database/sql"
	"strconv"
	"time"

	db "github.com/cihanalici/api/db/sqlc"
	"github.com/cihanalici/api/util"
	"github.com/gin-gonic/gin"
)

type orderRequest struct {
	UserID      *int32 `json:"user_id" binding:"required"`
	TotalAmount string `json:"total_amount" binding:"required"`
	Status      string `json:"status" binding:"required"`
}

type orderResponse struct {
	ID          int32     `json:"id"`
	UserID      *int32    `json:"user_id"`
	TotalAmount string    `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func orderNotation(order db.Order) orderResponse {
	var userID *int32
	if order.UserID.Valid {
		userID = &order.UserID.Int32
	}

	return orderResponse{
		ID:          order.ID,
		UserID:      userID,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}

func ordersNotation(orders []db.Order) []orderResponse {
	result := make([]orderResponse, len(orders))

	for i, order := range orders {
		result[i] = orderNotation(order)
	}

	return result
}

// CreateOrder godoc
// @Summary Create a new order
// @Tags orders
// @Description create a new order
// @Accept  json
// @Produce  json
// @Param input body orderRequest true "Order Request"
// @Success 200 {object} orderResponse
// @Router /orders [post]

func (server *Server) createOrder(ctx *gin.Context) {
	var req orderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.CreateOrderParams{
		UserID:      sql.NullInt32{Int32: *req.UserID, Valid: true},
		TotalAmount: req.TotalAmount,
		Status:      req.Status,
	}

	order, err := server.store.CreateOrder(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, orderNotation(order))
}

// GetOrder godoc
// @Summary Get an order by ID
// @Tags orders
// @Description get an order by ID
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Success 200 {object} orderResponse

type getOrderRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getOrder(ctx *gin.Context) {
	var req getOrderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	order, err := server.store.GetOrderById(ctx, req.ID)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, orderNotation(order))
}

// ListOrders godoc
// @Summary List all orders
// @Tags orders
// @Description list all orders
// @Accept  json
// @Produce  json
// @Param user_id query int false "User ID"
// @Param status query string false "Status"
// @Success 200 {array} orderResponse

type listOrdersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) ListOrders(ctx *gin.Context) {
	var req listOrdersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.ListOrdersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	orders, err := server.store.ListOrders(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, ordersNotation(orders))
}

// UpdateOrder godoc
// @Summary Update an order
// @Tags orders
// @Description update an order
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Param input body orderRequest true "Order Request"
// @Success 200 {object} orderResponse

type updateOrderRequest struct {
	ID          int32  `json:"id"`
	UserID      *int32 `json:"user_id"`
	TotalAmount string `json:"total_amount"`
	Status      string `json:"status"`
}

func (server *Server) updateOrder(ctx *gin.Context) {
	orderId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	order, err := server.store.GetOrderById(ctx, int32(orderId))
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	if order.ID != int32(orderId) {
		ctx.JSON(400, errorResponse(err))
		return
	}

	var req updateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	userId := sql.NullInt32{Valid: false}
	if req.UserID != nil {
		userId = sql.NullInt32{Int32: *req.UserID, Valid: true}
	}

	arg := db.UpdateOrderParams{
		ID:          int32(orderId),
		TotalAmount: req.TotalAmount,
		Status:      req.Status,
		UserID:      userId,
	}

	order, err = server.store.UpdateOrder(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, orderNotation(order))
}

// DeleteOrder godoc
// @Summary Delete an order
// @Tags orders
// @Description delete an order
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Success 200

func (server *Server) deleteOrder(ctx *gin.Context) {
	orderId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	err = server.store.DeleteOrder(ctx, int32(orderId))
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, gin.H{"status": "ok"})
}

// GetOrdersByUserId godoc
// @Summary Get orders by user ID
// @Tags orders
// @Description get orders by user ID
// @Accept  json
// @Produce  json
// @Param user_id path int true "User ID"
// @Success 200 {array} orderResponse

type getOrdersByUserIdRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) getOrdersByUserId(ctx *gin.Context) {
	var req getOrdersByUserIdRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	// bu alınan tokeni int32'ye çevirip userId'e atamamız gerekiyor
	userIdByToken, _ := ctx.Get("userId")
	userId, _ := userIdByToken.(int32)

	arg := db.GetOrdersByUserIdParams{
		UserID: util.ToInt32ToNullInt32(userId),
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	orders, err := server.store.GetOrdersByUserId(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, ordersNotation(orders))
}
