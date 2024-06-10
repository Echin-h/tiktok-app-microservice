package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"tiktok-app-microservice/common/err/rpcErr"
	"tiktok-app-microservice/common/model"
	"tiktok-app-microservice/common/utils"

	"tiktok-app-microservice/service/rpc/user/internal/svc"
	"tiktok-app-microservice/service/rpc/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FollowUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFollowUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FollowUserLogic {
	return &FollowUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FollowUserLogic) FollowUser(in *user.FollowUserRequest) (*user.Empty, error) {
	l.WithContext(l.ctx).Infof("关注用户: %v", in)

	err := l.svcCtx.DBList.Mysql.Transaction(func(tx *gorm.DB) error {
		var u *model.User
		var followUser *model.User
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", in.UserId).First(u)
		tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", in.FollowUserId).First(followUser)

		// 1. 处理关注用户
		result, err := l.svcCtx.DBList.Redis.Exists(l.ctx, utils.GenUserInfoCacheKey(in.UserId)).Result()
		if err != nil {
			l.Logger.Error(rpcErr.CacheError.Code, err.Error())
		}
		if result == 1 {
			// TODO: 如果是大V (在redis中有缓存), 就只更新缓存, 交给定时任务更新数据库
		} else {
			// 如果是普通用户, 就直接更新数据库
			u.FollowCount++
			err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Update("follow_count", u.FollowCount).Error
			if err != nil {
				l.WithContext(l.ctx).Errorf("更新普通用户关注数失败: %v", err)
				return status.Error(rpcErr.DataBaseError.Code, err.Error())
			}
		}

		// 2. 处理被关注用户
		result, err = l.svcCtx.DBList.Redis.Exists(l.ctx, utils.GenUserInfoCacheKey(in.FollowUserId)).Result()
		if err != nil {
			l.Logger.Error(rpcErr.CacheError.Code, err.Error())
		}
		if result == 1 {
			// TODO: 如果是大V (在redis中有缓存), 就只更新缓存, 交给定时任务更新数据库
		} else {
			followUser.FanCount++
			err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Update("fan_count", followUser.FanCount).Error
			if err != nil {
				return status.Error(rpcErr.DataBaseError.Code, err.Error())
			}
		}

		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(u).Association("Follows").Append(followUser)
		if err != nil {
			return status.Error(rpcErr.DataBaseError.Code, err.Error())
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}
	return &user.Empty{}, nil
}
