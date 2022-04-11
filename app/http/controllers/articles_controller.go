package controllers

import (
	"database/sql"
	"fmt"
	"goblog/pkg/logger"
	"goblog/pkg/model/article"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"html/template"
	"net/http"
)

type ArticlesController struct {
}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	//2.读取对应文章列表
	article, err := article.Get(id)
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
			"RouteName2URL":  route.Name2URL,
			"Uint64ToString": types.Uint64ToString,
		}).ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)

		tmpl.Execute(w, article)
	}
}
