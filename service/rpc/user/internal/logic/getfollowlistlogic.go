package logic

import (
	"context"

	"tiktok-app-microservice/service/rpc/user/internal/svc"
	"tiktok-app-microservice/service/rpc/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowListLogic {
	return &GetFollowListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowListLogic) GetFollowList(in *user.GetFollowListRequest) (*user.GetFollowListReply, error) {
	// todo: add your logic here and delete this line

	return &user.GetFollowListReply{}, nil
}
