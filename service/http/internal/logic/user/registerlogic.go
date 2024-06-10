package user

import (
	"context"
	"tiktok-app-microservice/common/err/apiErr"
	"tiktok-app-microservice/common/err/rpcErr"
	"tiktok-app-microservice/common/middleware"
	"tiktok-app-microservice/service/rpc/user/types/user"

	"tiktok-app-microservice/service/http/internal/svc"
	"tiktok-app-microservice/service/http/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	logx.WithContext(l.ctx).Infof("register: %v", req)
	if len(req.Username) > 32 {
		return nil, apiErr.InvalidParams.WithDetails("username up to 32 characters")
	} else if len(req.Password) > 32 {
		return nil, apiErr.InvalidParams.WithDetails("password up to 32 characters")
	}

	createUser, err := l.svcCtx.UserRpc.CreateUser(l.ctx, &user.CreateUserRequest{
		Name:     req.Username,
		Password: req.Password,
	})

	if rpcErr.Is(err, rpcErr.UserAlreadyExist) {
		return nil, apiErr.UserAlreadyExist
	} else if err != nil {
		logx.WithContext(l.ctx).Errorf("rpc CreateUser failed: %v", err)
		return nil, apiErr.InternalError(l.ctx, err.Error())
	}

	token, err := middleware.CreatToken(
		"",
		createUser.Id,
		l.svcCtx.Config.Auth.AccessSecret,
		l.svcCtx.Config.Auth.AccessExpire,
	)

	if err != nil {
		return nil, apiErr.GenerateTokenFailed
	}

	return &types.RegisterResp{
		BasicReply: types.BasicReply(apiErr.Success),
		UserId:     createUser.GetId(),
		Token:      token,
	}, nil
}
