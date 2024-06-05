package user

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"tiktok-app-microservice/common/err/apiErr"
	"tiktok-app-microservice/common/err/rpcErr"
	"tiktok-app-microservice/common/middleware"
	"tiktok-app-microservice/service/http/internal/svc"
	"tiktok-app-microservice/service/http/internal/types"
	"tiktok-app-microservice/service/rpc/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	GetUserByNameReply, err := l.svcCtx.UserRpc.GetUserByName(l.ctx, &user.GetUserByNameRequest{
		Name: req.Username,
	})
	if rpcErr.Is(err, rpcErr.UserNotExist) {
		return nil, apiErr.UserNotFound
	} else if err != nil {
		logx.WithContext(l.ctx).Errorf("LoginLogic.Login GetUserByName err: %v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(GetUserByNameReply.Password), []byte(req.Password))
	if err != nil {
		return nil, apiErr.PasswordIncorrect
	}

	AccessToken, err := middleware.CreatToken(
		"",
		GetUserByNameReply.Id,
		l.svcCtx.Config.Auth.AccessSecret,
		l.svcCtx.Config.Auth.AccessExpire,
	)

	if err != nil {
		logx.WithContext(l.ctx).Errorf("LoginLogic.Login CreateToken err: %v", err)
		return nil, apiErr.GenerateTokenFailed
	}

	return &types.LoginResp{
		BasicReply: types.BasicReply(apiErr.Success),
		UserId:     int(GetUserByNameReply.GetId()),
		Token:      AccessToken,
	}, nil
}
