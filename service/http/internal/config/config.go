package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	UserRpc zrpc.RpcClientConf
	Redis   RedisConf
}

type RedisConf struct {
	Address string `yaml:"address" json:"address"`
	DB      int    `yaml:"db" json:"db"`
}
