package logic

import (
	"context"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"tiktok-app-microservice/common/err/rpcErr"
	"tiktok-app-microservice/common/model"
	"tiktok-app-microservice/common/utils"

	"tiktok-app-microservice/service/rpc/user/internal/svc"
	"tiktok-app-microservice/service/rpc/user/types/user"

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

// 先从缓存中获取用户信息，如果缓存中没有，再从数据库中获取
func (l *GetUserByIdLogic) GetUserById(in *user.GetUserByIdRequest) (*user.GetUserReply, error) {
	result, err := l.svcCtx.DBList.Redis.HGetAll(l.ctx, utils.GenUserInfoCacheKey(in.Id)).Result()
	if err == nil && len(result) != 0 {
		l.WithContext(l.ctx).Info("Get user info from cache")
		return &user.GetUserReply{
			Id:          utils.Str2Int64(result["Id"]),
			Name:        result["Name"],
			Password:    result["Password"],
			FollowCount: utils.Str2Int64(result["FollowCount"]),
			FanCount:    utils.Str2Int64(result["FanCount"]),
		}, nil
	} else if err != nil && err != redis.Nil {
		l.WithContext(l.ctx).Error(rpcErr.CacheError.Code, err.Error())
	}

	var u model.User
	err = l.svcCtx.DBList.Mysql.Where("id = ?", in.Id).First(&u).Error

	if err == gorm.ErrRecordNotFound {
		return nil, status.Error(rpcErr.UserNotExist.Code, rpcErr.UserNotExist.Message)
	} else if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	if model.IsPopularUser(u.FanCount) {
		// 缓存个人信息
		err = l.svcCtx.DBList.Redis.HMSet(l.ctx, utils.GenUserInfoCacheKey(in.Id), map[string]interface{}{
			"Id":          u.ID,
			"Name":        u.Username,
			"Password":    u.Password,
			"FollowCount": u.FollowCount,
			"FanCount":    u.FanCount,
		}).Err()
		if err != nil {
			return nil, status.Error(rpcErr.CacheError.Code, err.Error())
		}

		err = l.svcCtx.DBList.Redis.LPush(l.ctx, utils.GenPopUserListCacheKey(), u.ID).Err()
		if err != nil {
			return nil, status.Error(rpcErr.CacheError.Code, err.Error())
		}
	}

	return &user.GetUserReply{
		Id:          int64(u.ID),
		Name:        u.Username,
		Password:    u.Password,
		FollowCount: u.FollowCount,
		FanCount:    u.FanCount,
	}, nil
}
