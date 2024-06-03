package logic

import (
	"context"

	"tiktok-app-microservice/user/rpc/internal/svc"
	"tiktok-app-microservice/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByIdLogic) GetUserById(in *user.IdRequest) (*user.UserResponse, error) {
	return &user.UserResponse{
		Id:       "2",
		Name:     "test",
		Password: "123456",
	}, nil
}
