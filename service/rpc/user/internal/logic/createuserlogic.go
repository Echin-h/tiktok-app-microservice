package logic

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/status"
	"tiktok-app-microservice/common/err/rpcErr"
	"tiktok-app-microservice/common/model"

	"tiktok-app-microservice/service/rpc/user/internal/svc"
	"tiktok-app-microservice/service/rpc/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateUserLogic) CreateUser(u *user.CreateUserRequest) (*user.CreatUserReply, error) {
	if u.Name == "" || u.Password == "" {
		return nil, status.Error(rpcErr.NameOrPasswordEmpty.Code, rpcErr.NameOrPasswordEmpty.Message)
	}
	tx := l.svcCtx.DBList.Mysql.Begin()

	var count int64
	if err := tx.Model(&model.User{}).Where("username = ?", u.Name).Count(&count).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}
	if count > 0 {
		tx.Rollback()
		return nil, status.Error(rpcErr.UserAlreadyExist.Code, rpcErr.UserAlreadyExist.Message)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return nil, status.Error(rpcErr.PassWordEncryptFailed.Code, err.Error())
	}

	newUser := &model.User{
		Username: u.Name,
		Password: string(password),
	}

	if err := tx.Create(newUser).Error; err != nil {
		tx.Rollback()
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, status.Error(rpcErr.DataBaseError.Code, err.Error())
	}

	return &user.CreatUserReply{
		Id: int64(newUser.ID),
	}, nil
}
