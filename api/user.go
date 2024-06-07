package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "github.com/cihanalici/api/db/sqlc"
	"github.com/cihanalici/api/token"
	"github.com/cihanalici/api/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Name      string          `json:"name" binding:"required"`
	Email     string          `json:"email" binding:"required,email"`
	Password  string          `json:"password" binding:"required,min=6"`
	Addresses json.RawMessage `json:"addresses"`
}
type Address struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	State    string `json:"state"`
	Zip      string `json:"zip"`
	Country  string `json:"country"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

type createUserResponse struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
	// addresses json data type
	Addresses json.RawMessage `json:"addresses"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func userResponse(user db.User) createUserResponse {
	addressesBytes, err := json.Marshal(user.Addresses)
	if err != nil {
		return createUserResponse{}
	}

	return createUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Addresses: addressesBytes,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func usersResponse(users []db.User) []createUserResponse {
	var result []createUserResponse
	for _, user := range users {
		result = append(result, userResponse(user))
	}
	return result
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      "user",
		Addresses: req.Addresses,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

		fmt.Println(err)

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println(user)

	rsp := userResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	Token string             `json:"token"`
	User  createUserResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	token, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		Token: token,
		User:  userResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}

type getUserRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

// users/:id
func (server *Server) getUserByID(ctx *gin.Context) {
	var req getUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, userResponse(user))
}

type getUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) getUsers(ctx *gin.Context) {
	var req getUsersRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	users, err := server.store.ListUsers(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, usersResponse(users))
}

type updateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`

	Role      string          `json:"role"`
	Addresses json.RawMessage `json:"addresses"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, int32(userId))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if user.ID != int32(userId) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	var req updateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams{
		Name:      req.Name,
		Email:     req.Email,
		Role:      req.Role,
		Addresses: req.Addresses,
		ID:        user.ID,
	}

	user, err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, userResponse(user))
}

func (server *Server) deleteUser(ctx *gin.Context) {
	userId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, int32(userId))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if user.ID != int32(userId) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteUser(ctx, int32(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// reset password işlemleri

type requestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type resetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (server *Server) requestPasswordReset(ctx *gin.Context) {
	var req requestPasswordResetRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	token, err := token.GenerateResetToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.storeResetToken(ctx, user.ID, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.SendResetEmail(req.Email, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "password reset email sent"})
}

func (server *Server) storeResetToken(ctx *gin.Context, userID int32, token string) error {
	arg := db.CreatePasswordResetParams{
		UserID:     userID,
		ResetToken: token,
		ExpiresAt:  time.Now().Add(time.Hour * 1),
	}

	_, err := server.store.CreatePasswordReset(ctx, arg)
	if err != nil {
		return err
	}

	return nil
}

func (server *Server) resetPassword(ctx *gin.Context) {
	var req resetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// bu alınan tokeni int32'ye çevirip userId'e atamamız gerekiyor
	userIdByToken, _ := ctx.Get("userId")
	userId, _ := userIdByToken.(int32)

	fmt.Println(userId)
	fmt.Println(req.Token)

	err := server.verifyToken(ctx, userId, req.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:       userId,
		Password: hashedPassword,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "password reset successfully"})
}

func (server *Server) verifyToken(ctx *gin.Context, userID int32, token string) error {
	row, err := server.store.GetPasswordResetByToken(ctx, token)
	if err != nil {
		return err
	}

	if row.UserID != userID {
		return util.ErrInvalidToken
	}

	if time.Now().After(row.ExpiresAt) {
		return util.ErrExpiredToken
	}

	return nil
}
