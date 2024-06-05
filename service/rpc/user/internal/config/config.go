package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	// 数据库配置
	DBList DBListConf
}

type DBListConf struct {
	Mysql MysqlConf `yaml:"mysql" json:"mysql"`
	Redis RedisConf `yaml:"redis" json:"redis"`
}

type MysqlConf struct {
	Address     string `yaml:"address" json:"address"`
	Username    string `yaml:"username" json:"username"`
	Password    string `yaml:"password" json:"password"`
	DBName      string `yaml:"dbname" json:"dbname"`
	TablePrefix string `yaml:"tablePrefix" json:"tablePrefix"`
}

type RedisConf struct {
	Address string `yaml:"address" json:"address"`
	DB      int    `yaml:"db" json:"db"`
}
