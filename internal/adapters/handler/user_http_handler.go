package handler

import (
	"net/http"

	"github.com/LordMoMA/Hexagonal-Architecture/internal/core/domain"
	"github.com/LordMoMA/Hexagonal-Architecture/internal/core/services"
	"github.com/gin-gonic/gin"
)

type (
	UserHttpHandlerService interface {
		LoginUser(ctx *gin.Context)
		CreateUser(ctx *gin.Context)
		ForgetPassword(ctx *gin.Context)
		ResetPassword(ctx *gin.Context)
	}

	userHttpHandler struct {
		userUsecase services.UserUsecaseService
	}
)

func NewUserHttpHandler(userUsecase services.UserUsecaseService) UserHttpHandlerService {
	return &userHttpHandler{userUsecase}
}

func (h *userHttpHandler) LoginUser(ctx *gin.Context) {
	req := new(domain.LoginUserReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := h.userUsecase.LoginUser(req)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *userHttpHandler) CreateUser(ctx *gin.Context) {
	req := new(domain.CreateUserReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	err := h.userUsecase.CreateUser(req)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "New user created successfully",
	})
}

func (h *userHttpHandler) ForgetPassword(ctx *gin.Context) {
	req := new(domain.ForgetPasswordReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	response, err := h.userUsecase.ForgetPassword(req)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (h *userHttpHandler) ResetPassword(ctx *gin.Context) {
	req := new(domain.ResetPasswordReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	err := h.userUsecase.ResetPassword(req)
	if err != nil {
		HandleError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}
