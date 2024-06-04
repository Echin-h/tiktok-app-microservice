package svc

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"tiktok-app-microservice/common/model"
	"tiktok-app-microservice/service/rpc/user/internal/config"
	"time"
)

type ServiceContext struct {
	Config config.Config
	DBList *DBList
}

type DBList struct {
	Mysql *gorm.DB
	Redis *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		DBList: initDB(c),
	}
}

func initDB(c config.Config) *DBList {
	dbList := new(DBList)
	dbList.Mysql = initMysql(c)
	dbList.Redis = initRedis(c)
	return dbList
}

func initMysql(c config.Config) *gorm.DB {
	fmt.Println("connect Mysql ...")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBList.Mysql.Username,
		c.DBList.Mysql.Password,
		c.DBList.Mysql.Address,
		c.DBList.Mysql.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.DBList.Mysql.TablePrefix, // 表名前缀
			SingularTable: true,                       // 使用单数表名
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	// 自动建表
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}

	fmt.Println("connect Mysql success")
	return db
}

func initRedis(c config.Config) *redis.Client {
	fmt.Println("connect Redis ...")
	client := redis.NewClient(&redis.Options{
		Addr: c.DBList.Redis.Address,
		DB:   c.DBList.Redis.DB,

		// timeout:
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolTimeout:  3 * time.Second,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("connect Redis success")
	return client
}
