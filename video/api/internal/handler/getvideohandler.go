package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tiktok-app-microservice/video/api/internal/logic"
	"tiktok-app-microservice/video/api/internal/svc"
	"tiktok-app-microservice/video/api/internal/types"
)

func getVideoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VideoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetVideoLogic(r.Context(), svcCtx)
		resp, err := l.GetVideo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
