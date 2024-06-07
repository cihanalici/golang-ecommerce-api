package api

import (
	"fmt"
	"strconv"
	"time"

	db "github.com/cihanalici/api/db/sqlc"
	"github.com/gin-gonic/gin"
)

type OrderItemRequest struct {
	OrderID          int32  `json:"order_id"`
	ProductVariantID int32  `json:"product_variant_id"`
	Quantity         int32  `json:"quantity"`
	Price            string `json:"price"`
}

type OrderItemResponse struct {
	ID               int32     `json:"id"`
	OrderID          int32     `json:"order_id"`
	ProductVariantID int32     `json:"product_variant_id"`
	Quantity         int32     `json:"quantity"`
	Price            string    `json:"price"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func OrderItemNotation(orderItem db.OrderItem) OrderItemResponse {

	return OrderItemResponse{
		ID:               orderItem.ID,
		OrderID:          orderItem.OrderID,
		ProductVariantID: orderItem.ProductVariantID,
		Quantity:         orderItem.Quantity,
		Price:            orderItem.Price,
		CreatedAt:        orderItem.CreatedAt,
		UpdatedAt:        orderItem.UpdatedAt,
	}
}

func OrderItemsNotation(orderItems []db.OrderItem) []OrderItemResponse {
	result := make([]OrderItemResponse, len(orderItems))

	for i, orderItem := range orderItems {
		result[i] = OrderItemNotation(orderItem)
	}

	return result
}

// getOrderItems godoc
// @Summary Get order items
// @Description Get all order items
// @Tags order_items
// @Accept json
// @Produce json
// @Success 200 {array} OrderItemResponse
// @Router /order_items [get]

type listOrderItemsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listOrderItems(ctx *gin.Context) {
	var req listOrderItemsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.ListOrderItemsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	orderItems, err := server.store.ListOrderItems(ctx, arg)

	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, OrderItemsNotation(orderItems))
}

// getOrderItem godoc
// @Summary Get an order item by ID
// @Description Get an order item by ID
// @Tags order_items
// @Accept json
// @Produce json
// @Param id path int true "Order Item ID"
// @Success 200 {object} OrderItemResponse
// @Router /order_items/{id} [get]

type getOrderItemRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getOrderItem(ctx *gin.Context) {
	var req getOrderItemRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	orderItem, err := server.store.GetOrderItemById(ctx, req.ID)

	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, OrderItemNotation(orderItem))
}

// getOrderItemsByOrderId godoc
// @Summary Get order items by order ID
// @Description Get order items by order ID
// @Tags order_items
// @Accept json
// @Produce json
// @Param order_id path int true "Order ID"
// @Success 200 {array} OrderItemResponse

type getOrderItemsByOrderIdRequest struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (server *Server) getOrderItemsByOrderId(ctx *gin.Context) {
	orderId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)

	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	var req getOrderItemsByOrderIdRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	arg := db.GetOrderItemsByOrderIdParams{
		OrderID: int32(orderId),
		Limit:   req.Limit,
		Offset:  req.Offset,
	}

	fmt.Println(arg)

	orderItems, err := server.store.GetOrderItemsByOrderId(ctx, arg)

	if err != nil {
		ctx.JSON(400, errorResponse(err))
		return
	}

	ctx.JSON(200, OrderItemsNotation(orderItems))
}
