package middleware

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"tiktok-app-microservice/common/err/apiErr"
	"tiktok-app-microservice/common/middleware"
	"tiktok-app-microservice/service/http/internal/config"
)

type AuthMiddleware struct {
	Config config.Config
}

func NewAuthMiddleware(c config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		Config: c,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string
		if token = r.URL.Query().Get("token"); token == "" {
			token = r.PostFormValue("token")
		}
		if token == "" {
			httpx.OkJson(w, apiErr.NotLogin)
			return
		}
		isTimeOut, err := middleware.ValidToken(token, m.Config.Auth.AccessSecret)
		if err != nil || isTimeOut {
			httpx.OkJson(w, apiErr.InvalidToken)
			return
		}
		next(w, r)
	}
}
