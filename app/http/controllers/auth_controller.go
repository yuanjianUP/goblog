package controllers

import (
	"fmt"
	"goblog/app/models/user"
	"goblog/pkg/view"
	"net/http"
)

type AuthController struct{}

func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	//表单验证
	_user := user.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	_user.Create()
	if _user.ID > 0 {
		fmt.Fprint(w, "插入成功,ID为"+_user.GetStringId())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "创建用户失败，请联系管理员")
	}

	//验证通过，入库并跳转
	//表单不通过
}
