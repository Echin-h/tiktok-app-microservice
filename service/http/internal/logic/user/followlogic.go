package user

import (
	"context"
	"tiktok-app-microservice/common/err/apiErr"
	"tiktok-app-microservice/common/middleware"
	"tiktok-app-microservice/service/rpc/user/types/user"

	"tiktok-app-microservice/service/http/internal/svc"
	"tiktok-app-microservice/service/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowLogic {
	return &FollowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FollowLogic) Follow(req *types.FollowReq) (resp *types.FollowResp, err error) {
	l.WithContext(l.ctx).Infof("关注用户: %v", req)

	var id int64
	id, err = middleware.GetUserIdFromToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, apiErr.InvalidToken
	}

	if id == req.ToUserId {
		return nil, apiErr.IllegalOperation.WithDetails("不能关注自己")
	}

	if req.ActionType == 1 {
		// 查看是否已经关注
		isFollowReply, _ := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
			UserId:       id,
			FollowUserId: req.ToUserId,
		})
		if isFollowReply.IsFollow {
			logx.WithContext(l.ctx).Errorf("已经关注过了")
			return nil, apiErr.AlreadyFollowed
		}

		// 关注
		_, err := l.svcCtx.UserRpc.FollowUser(l.ctx, &user.FollowUserRequest{
			UserId:       id,
			FollowUserId: req.ToUserId,
		})

		if err != nil {
			logx.WithContext(l.ctx).Errorf("关注失败: %v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

		// TODO: 异步任务添加好友关系
	} else if req.ActionType == 0 {
		// 查看是否关注
		isFollowReply, _ := l.svcCtx.UserRpc.IsFollow(l.ctx, &user.IsFollowRequest{
			UserId:       id,
			FollowUserId: req.ToUserId,
		})

		if !isFollowReply.IsFollow {
			logx.WithContext(l.ctx).Errorf("还未关注")
			return nil, apiErr.NotFollowed
		}

		// 取消关注
		_, err := l.svcCtx.UserRpc.UnFollowUser(l.ctx, &user.UnFollowUserRequest{
			UserId:         id,
			UnFollowUserId: req.ToUserId,
		})

		if err != nil {
			logx.WithContext(l.ctx).Errorf("取消关注失败: %v", err)
			return nil, apiErr.InternalError(l.ctx, err.Error())
		}

		// TODO: 异步任务解除好友关系
	} else {
		logx.WithContext(l.ctx).Errorf("ActionType参数错误")
		return nil, apiErr.InvalidParams.WithDetails("ActionType参数错误")
	}
	return &types.FollowResp{
		BasicReply: types.BasicReply(apiErr.Success),
	}, nil
}
