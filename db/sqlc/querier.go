// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sqlc

import (
	"context"
	"time"
)

type Querier interface {
	CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error)
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (OrderItem, error)
	CreatePasswordReset(ctx context.Context, arg CreatePasswordResetParams) (CreatePasswordResetRow, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	CreateProductVariant(ctx context.Context, arg CreateProductVariantParams) (ProductVariant, error)
	CreateReview(ctx context.Context, arg CreateReviewParams) (Review, error)
	CreateSale(ctx context.Context, arg CreateSaleParams) (Sale, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateWishlistItem(ctx context.Context, arg CreateWishlistItemParams) (Wishlist, error)
	DeleteCategory(ctx context.Context, id int32) error
	DeleteExpiredPasswordResets(ctx context.Context) error
	DeleteOrder(ctx context.Context, id int32) error
	DeleteOrderItem(ctx context.Context, id int32) error
	DeletePasswordReset(ctx context.Context, resetToken string) error
	DeleteProduct(ctx context.Context, id int32) error
	DeleteProductVariant(ctx context.Context, id int32) error
	DeleteReview(ctx context.Context, id int32) error
	DeleteSale(ctx context.Context, id int32) error
	DeleteUser(ctx context.Context, id int32) error
	DeleteWishlistItem(ctx context.Context, id int32) error
	GetCategoryById(ctx context.Context, id int32) (Category, error)
	GetMonthlySales(ctx context.Context, createdAt time.Time) ([]GetMonthlySalesRow, error)
	GetOrderById(ctx context.Context, id int32) (Order, error)
	GetOrderItemById(ctx context.Context, id int32) (OrderItem, error)
	GetOrderItemsByOrderId(ctx context.Context, arg GetOrderItemsByOrderIdParams) ([]OrderItem, error)
	GetOrdersByUserId(ctx context.Context, arg GetOrdersByUserIdParams) ([]Order, error)
	GetPasswordResetByID(ctx context.Context, id int32) (PasswordReset, error)
	GetPasswordResetByToken(ctx context.Context, resetToken string) (GetPasswordResetByTokenRow, error)
	GetPasswordResetByUserId(ctx context.Context, userID int32) (GetPasswordResetByUserIdRow, error)
	GetPasswordResetByUserIdAndToken(ctx context.Context, arg GetPasswordResetByUserIdAndTokenParams) (GetPasswordResetByUserIdAndTokenRow, error)
	GetProductById(ctx context.Context, id int32) (Product, error)
	GetProductVariantById(ctx context.Context, id int32) (ProductVariant, error)
	GetReviewById(ctx context.Context, id int32) (Review, error)
	GetReviewsByProductId(ctx context.Context, arg GetReviewsByProductIdParams) ([]Review, error)
	GetSaleById(ctx context.Context, id int32) (Sale, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id int32) (User, error)
	GetWishlistItemById(ctx context.Context, id int32) (Wishlist, error)
	GetWishlistItemsByUserId(ctx context.Context, arg GetWishlistItemsByUserIdParams) ([]Wishlist, error)
	ListCategories(ctx context.Context, arg ListCategoriesParams) ([]Category, error)
	ListOrderItems(ctx context.Context, arg ListOrderItemsParams) ([]OrderItem, error)
	ListOrders(ctx context.Context, arg ListOrdersParams) ([]Order, error)
	ListProductVariants(ctx context.Context, arg ListProductVariantsParams) ([]ProductVariant, error)
	ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error)
	ListReviews(ctx context.Context, arg ListReviewsParams) ([]Review, error)
	ListSales(ctx context.Context, arg ListSalesParams) ([]Sale, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	ListWishlistItems(ctx context.Context, arg ListWishlistItemsParams) ([]Wishlist, error)
	UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error)
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error)
	UpdateOrderItem(ctx context.Context, arg UpdateOrderItemParams) (OrderItem, error)
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error)
	UpdateProductVariant(ctx context.Context, arg UpdateProductVariantParams) (ProductVariant, error)
	UpdateReview(ctx context.Context, arg UpdateReviewParams) (Review, error)
	UpdateSale(ctx context.Context, arg UpdateSaleParams) (Sale, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error)
	UpdateWishlistItem(ctx context.Context, arg UpdateWishlistItemParams) (Wishlist, error)
}

var _ Querier = (*Queries)(nil)
