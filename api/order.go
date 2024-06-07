package api

import (
	"fmt"
	"strconv"
	"time"

	db "github.com/cihanalici/api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type orderRequest struct {
	UserID      int32               `json:"user_id" binding:"required"`
	TotalAmount string              `json:"total_amount" binding:"required"`
	Status      string              `json:"status" binding:"required"`
	Items       []orderItemsRequest `json:"items" binding:"required"`
}

type orderItemsRequest struct {
	ProductVariantID int32  `json:"product_variant_id"`
	Quantity         int32  `json:"quantity"`
	Price            string `json:"price"`
}

type orderResponse struct {
	ID          int32     `json:"id"`
	UserID      int32     `json:"user_id"`
	TotalAmount string    `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func orderNotation(order db.Order) orderResponse {

	return orderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
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
		UserID:      req.UserID,
		TotalAmount: req.TotalAmount,
		Status:      req.Status,
	}

	order, err := server.store.CreateOrder(ctx, arg)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	fmt.Println(req.Items)

	// create order items
	for _, item := range req.Items {
		itemArg := db.CreateOrderItemParams{
			OrderID:          order.ID,
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
			Price:            item.Price,
		}

		orderItem, err := server.store.CreateOrderItem(ctx, itemArg)
		fmt.Println(orderItem)

		if err != nil {
			ctx.JSON(500, errorResponse(err))
			return
		}
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
	UserID      int32  `json:"user_id"`
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

	arg := db.UpdateOrderParams{
		ID:          int32(orderId),
		TotalAmount: req.TotalAmount,
		Status:      req.Status,
		UserID:      req.UserID,
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
		UserID: userId,
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

// GetMonthlySales godoc
// @Summary Get monthly revenue
// @Tags orders
// @Description get monthly revenue
// @Accept  json
// @Produce  json
// @Success 200 {object} monthlyRevenueResponse

func (server *Server) getMonthlySales(ctx *gin.Context) {
	year := time.Now().Year()
	startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)

	monthlySales, err := server.store.GetMonthlySales(ctx, startOfYear)
	if err != nil {
		ctx.JSON(500, errorResponse(err))
		return
	}

	ctx.JSON(200, monthlySales)
}
