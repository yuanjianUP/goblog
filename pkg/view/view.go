package view

import (
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"io"
	"path/filepath"
	"strings"
	"text/template"
)

func Render(w io.Writer, name string, data interface{}) {
	viewDir := "resources/views/"
	name = strings.Replace(name, ".", "/", -1)
	//所有布局模版文件slice
	files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
	logger.LogError(err)
	//在slice里新增我们目标文件
	newFiles := append(files, viewDir+name+".gohtml")
	//读取成功
	tmpl, err := template.New("show.gohtml").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).ParseFiles(newFiles...)
	logger.LogError(err)
	//渲染模版，将所有文行数据传输进去
	err = tmpl.ExecuteTemplate(w, "app", data)
	logger.LogError(err)
}
