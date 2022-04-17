package middlewares

import (
	"net/http"
	"strings"
)

//中间件
func ForceHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//1.设置标头
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		//2.继续处理请求
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/") //去掉右边的/
		}
		next.ServeHTTP(w, r)
	})
}
