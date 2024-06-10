package logic

import (
	"context"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"tiktok-app-microservice/common/err/rpcErr"
	"tiktok-app-microservice/common/model"

	"tiktok-app-microservice/service/rpc/user/internal/svc"
	"tiktok-app-microservice/service/rpc/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByNameLogic {
	return &GetUserByNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByNameLogic) GetUserByName(in *user.GetUserByNameRequest) (*user.GetUserReply, error) {
	var u model.User
	err := l.svcCtx.DBList.Mysql.Where("username = ?", in.Name).First(&u).Error

	if err == gorm.ErrRecordNotFound {
		return nil, status.Error(rpcErr.UserNotExist.Code, rpcErr.UserNotExist.Message)
	} else if err != nil {
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	return &user.GetUserReply{
		Id:          int64(u.ID),
		Name:        u.Username,
		Password:    u.Password,
		FollowCount: u.FollowCount,
		FanCount:    u.FanCount,
	}, nil
}
