package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"tiktok-app-microservice/service/http/internal/config"
	"tiktok-app-microservice/service/http/internal/middleware"
	"tiktok-app-microservice/service/rpc/user/userclient"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc userclient.User
	Auth    rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Auth:    middleware.NewAuthMiddleware(c).Handle,
	}
}
