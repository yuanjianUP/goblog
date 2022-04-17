package main

import (
	"fmt"
	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"goblog/pkg/route"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter().StrictSlash(true) //strictslash可以使go和go/都能正确访问

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>hello,欢迎来到goblog</h1>")
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客事用以记录编程笔记，如您有反馈或建议，请联系")
}
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+"<p>如有疑问，请联系我们</p>")
}

func main() {
	bootstrap.SetupDB()

	router = bootstrap.SetupRoute()
	route.SetRoute(router)
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	//使用中间件：强制内容类型为HTML
	http.ListenAndServe(":8090", middlewares.RemoveTrailingSlash(router))
}
