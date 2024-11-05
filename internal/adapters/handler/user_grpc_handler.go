package handler

import (
	"context"

	userPb "github.com/LordMoMA/Hexagonal-Architecture/internal/core/proto/user"
	"github.com/LordMoMA/Hexagonal-Architecture/internal/core/services"
)

type (
	userGrpcHandler struct {
		userPb.UnimplementedUserGrpcServiceServer
		userUsecase services.UserUsecaseService
	}
)

func NewUserGrpcHandler(userUsecase services.UserUsecaseService) *userGrpcHandler {
	return &userGrpcHandler{
		UnimplementedUserGrpcServiceServer: userPb.UnimplementedUserGrpcServiceServer{},
		userUsecase:                        userUsecase,
	}
}

func (g *userGrpcHandler) LoginGrpc(ctx context.Context, req *userPb.LoginGrpcReq) (*userPb.LoginGrpcRes, error) {
	return g.userUsecase.LoginGrpc(req)
}
func (g *userGrpcHandler) CreateUserGrpc(ctx context.Context, req *userPb.CreateUserGrpcReq) (*userPb.CreateUserGrpcRes, error) {
	return g.userUsecase.CreateUserGrpc(req)
}
func (g *userGrpcHandler) ForgetPassword(ctx context.Context, req *userPb.ForgetPasswordReq) (*userPb.ForgetPasswordRes, error) {
	return g.userUsecase.ForgetPasswordGrpc(req)
}
func (g *userGrpcHandler) ResetPassword(ctx context.Context, req *userPb.ResetPasswordReq) (*userPb.ResetPasswordRes, error) {
	return g.userUsecase.ResetPasswordGrpc(req)

}
