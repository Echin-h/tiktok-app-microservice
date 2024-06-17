package scheduler

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"tiktok-app-microservice/common/cron"
	"tiktok-app-microservice/service/cron/internal/svc"
	"time"
)

type AsynqLogic struct {
	stx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAsynqServer(ctx context.Context, svcCtx *svc.ServiceContext) *AsynqLogic {
	return &AsynqLogic{
		stx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AsynqLogic) Start() {
	l.Info("AsynqTask starting...")

	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     l.svcCtx.Config.Redis.Address,
			Password: l.svcCtx.Config.Redis.Password},
		nil,
	)

	u := cron.SyncUserInfoCache()
	uid, err := scheduler.Register("@every 1h", u)
	if err != nil {
		log.Fatal(err)
	}
	l.WithContext(l.stx).Infof("registered an entry: %q", uid)

	v := cron.SyncVideoInfoCache()
	vid, err := scheduler.Register("@every 301s", v)
	if err != nil {
		log.Fatal(err)
	}
	l.WithContext(l.stx).Infof("registered an entry: %q", vid)

	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}
	l.WithContext(l.stx).Info("AsynqTask started")
}

func (l *AsynqLogic) Stop() {
	l.Info("Gracefully shutting down Asynq server...")

	ctx, cancelFunc := context.WithTimeout(l.stx, 5*time.Second)
	defer cancelFunc()

	stopSignal := make(chan struct{})
	close(stopSignal)

	time.Sleep(1 * time.Second)

	l.Info("AsynqTask stopped")
	err := l.svcCtx.Redis.Close()
	if err != nil {
		l.Error("Error closing Redis connection")
	}
	<-ctx.Done()
	l.Info("Redis connection closed successfully...")
}
