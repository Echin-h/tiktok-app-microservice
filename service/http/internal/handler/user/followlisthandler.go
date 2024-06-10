package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tiktok-app-microservice/service/http/internal/logic/user"
	"tiktok-app-microservice/service/http/internal/svc"
	"tiktok-app-microservice/service/http/internal/types"
)

func FollowListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FollowListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewFollowListLogic(r.Context(), svcCtx)
		resp, err := l.FollowList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}