package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"tiktok-app-microservice/user/rpc/types/user"
	"tiktok-app-microservice/video/api/internal/config"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user.UserServiceClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: user.NewUserServiceClient(zrpc.MustNewClient(c.UserRpc)),
	}
}
