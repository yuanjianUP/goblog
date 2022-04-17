package middlewares

import (
	"net/http"
	"strings"
)

//移除路由后的/符号
func RemoveTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		//请求传递下去
		next.ServeHTTP(w, r)
	})
}
