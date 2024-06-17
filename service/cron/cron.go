package main

import (
	"context"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"tiktok-app-microservice/service/cron/internal/config"
	"tiktok-app-microservice/service/cron/internal/scheduler"
	"tiktok-app-microservice/service/cron/internal/svc"
)

var configFile = flag.String("f", "etc/cron.yaml", "Specify the config file")

func main() {
	flag.Parse()
	// get config
	var c config.Config
	conf.MustLoad(*configFile, &c)
	// get rely
	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()
	// new service
	svcGroup := service.NewServiceGroup()
	defer svcGroup.Stop()
	// add service and start
	svcGroup.Add(scheduler.NewAsynqServer(ctx, svcContext))
	svcGroup.Start()

	logx.DisableStat()
}
