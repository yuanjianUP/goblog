package main

import (
	"database/sql"
	"fmt"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db = database.DB
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

//声明存储数据库数据
type Article struct {
	Title, Body string
	ID          int64
}

//中间件
func forceHTMLMiddleware(next http.Handler) http.Handler {
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

func removeTrailingSlash(next http.Handler) http.Handler { //移除路由后的/符号
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		//请求传递下去
		next.ServeHTTP(w, r)
	})
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	id := vars["id"]
	//读取对应的文章数据
	article, err := getArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.LogError(err)
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 未找到数据")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		rowsAffected, err := article.Delete()
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		} else {
			//未发生错误
			if rowsAffected > 0 {
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
}
func (a Article) Delete() (RowsAffected int64, err error) {
	rs, err := db.Exec("DELETE FROM articles WHERE id =" + strconv.FormatInt(a.ID, 10))
	if err != nil {
		return 0, err
	}
	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}
	return 0, nil
}
func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}
func main() {
	database.Initialize()
	bootstrap.SetupDB()

	router = bootstrap.SetupRoute()
	route.SetRoute(router)
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	//使用中间件：强制内容类型为HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":8090", removeTrailingSlash(router))
}
