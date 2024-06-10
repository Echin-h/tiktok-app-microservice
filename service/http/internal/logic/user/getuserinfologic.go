package user

import (
	"context"
	"errors"
	"fmt"
	"tiktok-app-microservice/common/err/apiErr"
	"tiktok-app-microservice/common/err/rpcErr"
	"tiktok-app-microservice/common/middleware"
	"tiktok-app-microservice/service/rpc/user/types/user"

	"tiktok-app-microservice/service/http/internal/svc"
	"tiktok-app-microservice/service/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoRequest) (resp *types.GetUserInfoResp, err error) {
	l.WithContext(l.ctx).Infof("获取用户信息: %v", req)

	id, err := middleware.GetUserIdFromToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.InvalidToken
	}
	userInfo, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &user.GetUserByIdRequest{
		Id: req.UserId,
	})

	if errors.Is(err, rpcErr.UserNotExist) {
		return nil, apiErr.UserNotFound
	} else if err != nil {
		l.WithContext(l.ctx).Errorf("获取用户信息失败: %v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	fmt.Println("这里是token的人是否关注了u.id的人")
	follow, err := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
		UserId:       id,
		FollowUserId: userInfo.GetId(),
	})

	if err != nil {
		l.WithContext(l.ctx).Errorf("获取是否关注失败: %v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	return &types.GetUserInfoResp{
		BasicReply: types.BasicReply(apiErr.Success),
		User: types.User{
			Id:            userInfo.GetId(),
			Name:          userInfo.GetName(),
			FollowCount:   userInfo.GetFollowCount(),
			FollowerCount: userInfo.GetFanCount(),
			IsFollowed:    follow.GetIsFollow(),
		},
	}, nil
}
