package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"tiktok-app-microservice/common/err/apiErr"

	"tiktok-app-microservice/service/http/internal/config"
	"tiktok-app-microservice/service/http/internal/handler"
	"tiktok-app-microservice/service/http/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/app.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case apiErr.ApiErr:
			return http.StatusOK, e.Response()
		case apiErr.ErrInternal:
			return http.StatusOK, e.Response(c.RestConf)
		default:
			return http.StatusInternalServerError, err
		}
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
