package controllers

import (
	"database/sql"
	"fmt"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"net/http"
	"text/template"
)

type ArticlesController struct {
}

//声明存储数据库数据
type Article struct {
	Title, Body string
	ID          int64
}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	//2.读取对应文章列表
	article, err := getArticleByID(id)
	//如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			//3.1数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404文章未找到")
		} else {
			//数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500服务器内部错误")
		}
	} else {
		//读取成功
		tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
			"Int64ToString": types.Int64ToString,
		}).ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)
		tmpl.Execute(w, article)
	}
}
func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}
