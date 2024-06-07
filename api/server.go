package api

import (
	"fmt"

	"github.com/cihanalici/api/db/sqlc"
	"github.com/cihanalici/api/token"
	"github.com/cihanalici/api/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	store      sqlc.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store sqlc.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/users/:id", server.getUserByID)
	authRoutes.GET("/users", server.getUsers)
	authRoutes.PUT("/users/:id", server.updateUser)
	authRoutes.DELETE("/users/:id", server.deleteUser)

	authRoutes.POST("/categories", server.createCategory)
	router.GET("/categories/:id", server.getCategory)
	router.GET("/categories", server.getCategories)
	authRoutes.PUT("/categories/:id", server.updateCategory)
	authRoutes.DELETE("/categories/:id", server.deleteCategory)

	authRoutes.POST("/products", server.createProduct)
	router.GET("/products/:id", server.getProduct)
	router.GET("/products", server.getProducts)
	authRoutes.PUT("/products/:id", server.updateProduct)
	authRoutes.DELETE("/products/:id", server.deleteProduct)

	authRoutes.POST("/orders", server.createOrder)
	router.GET("/orders/:id", server.getOrder)
	router.GET("/orders", server.ListOrders)
	authRoutes.PUT("/orders/:id", server.updateOrder)
	authRoutes.DELETE("/orders/:id", server.deleteOrder)
	authRoutes.GET("/orders/user", server.getOrdersByUserId)

	//product variants
	authRoutes.POST("/product_variants", server.createProductVariant)
	router.GET("/product_variants/:id", server.getProductVariant)
	router.GET("/product_variants", server.listProductVariants)
	authRoutes.PUT("/product_variants/:id", server.updateProductVariant)
	authRoutes.DELETE("/product_variants/:id", server.deleteProductVariant)

	//wishlist
	authRoutes.POST("/wishlists", server.createWishlist)
	router.GET("/wishlists/:id", server.getWishlist)
	router.GET("/wishlists", server.listWishlist)
	authRoutes.DELETE("/wishlists/:id", server.deleteWishlist)
	authRoutes.GET("/wishlists/user", server.getWishlistByUser)

	//order items
	authRoutes.GET("/order_items", server.listOrderItems)
	router.GET("/order_items/:id", server.getOrderItem)
	authRoutes.GET("/order_items/order/:id", server.getOrderItemsByOrderId)

	//reviews
	authRoutes.POST("/reviews", server.createReview)
	router.GET("/reviews/:id", server.getReview)
	router.GET("/reviews", server.listReviews)
	router.GET("/reviews/product/:id", server.getReviewsByProductId)
	authRoutes.DELETE("/reviews/:id", server.deleteReview)

	// reset password
	authRoutes.POST("/users/request-password-reset", server.requestPasswordReset)
	authRoutes.POST("/users/reset-password", server.resetPassword)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
