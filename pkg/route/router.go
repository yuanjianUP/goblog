package route

import (
	"net/http"

	"github.com/gorilla/mux"
)

var router *mux.Router

//通过路由名称来获取URL
func Name2URL(routeName string, pairs ...string) string {
	url, err := router.Get(routeName).URL(pairs...)
	if err != nil {
		//checkError(err)
		return ""
	}
	return url.String()
}

func GetRouteVariable(paramName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[paramName]
}
