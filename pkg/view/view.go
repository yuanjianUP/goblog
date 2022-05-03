package view

import (
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"io"
	"path/filepath"
	"strings"
	"text/template"
)

type D map[string]interface{}

func Render(w io.Writer, data interface{}, tplFiles ...string) {
	RenderTemplate(w, "app", data, tplFiles...)
}

func RenderSimple(w io.Writer, data interface{}, tplFiles ...string) {
	RenderTemplate(w, "simple", data, tplFiles...)
}

func RenderTemplate(w io.Writer, name string, data interface{}, tplFiles ...string) {
	viewDir := "resources/views/"
	for i, f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}
	//所有布局模版文件slice
	files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
	logger.LogError(err)
	//在slice里新增我们目标文件
	newFiles := append(files, tplFiles...)
	//读取成功
	tmpl, err := template.New("show.gohtml").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).ParseFiles(newFiles...)
	logger.LogError(err)
	//渲染模版，将所有文行数据传输进去
	err = tmpl.ExecuteTemplate(w, name, data)
	logger.LogError(err)
}
