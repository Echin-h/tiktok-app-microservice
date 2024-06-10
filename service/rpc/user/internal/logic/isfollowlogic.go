package logic

import (
	"context"
	"fmt"
	"tiktok-app-microservice/common/model"
	"tiktok-app-microservice/common/utils"
	"tiktok-app-microservice/service/rpc/user/internal/svc"
	"tiktok-app-microservice/service/rpc/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFollowLogic {
	return &IsFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFollowLogic) IsFollow(in *user.IsFollowRequest) (*user.IsFollowReply, error) {
	l.WithContext(l.ctx).Infof("判断是否关注: %v", in)

	if l.svcCtx.DBList.Redis.
		Exists(l.ctx, utils.GenFollowUserCacheKey(in.UserId, in.FollowUserId)).
		Val() == 1 {
		fmt.Println("缓存存在")
		return &user.IsFollowReply{
			IsFollow: true,
		}, nil
	}

	var cnt int64
	l.svcCtx.DBList.Mysql.Where("id = ?", in.UserId).
		Preload("Follows", "id = ?", in.FollowUserId).
		First(&model.User{}).Count(&cnt)

	if cnt > 0 {
		fmt.Println("关注了这个人")
		err := l.svcCtx.DBList.Redis.
			Set(l.ctx, utils.GenFollowUserCacheKey(in.UserId, in.FollowUserId), 1, utils.CacheExpire).Err()
		if err != nil {
			return nil, err
		}
		return &user.IsFollowReply{
			IsFollow: true,
		}, nil
	} else {
		fmt.Println("没有关注这个人")
		return &user.IsFollowReply{
			IsFollow: false,
		}, nil
	}
}
