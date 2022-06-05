package controllers

import (
	"encoding/json"
	"fmt"
	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/view"
	"net/http"
)

type AuthController struct{}

func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	_user := user.User{
		Name:            r.PostFormValue("name"),
		Email:           r.PostFormValue("email"),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirmation"),
	}
	errs := requests.ValidateRegistrationFrom(_user)
	if len(errs) > 0 {
		data, _ := json.MarshalIndent(errs, "", "")
		fmt.Fprint(w, string(data))
	} else {
		_user.Create()
		if _user.ID > 0 {
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建用户失败，请联系管理员")
		}
	}
	//验证通过，入库并跳转
	//表单不通过
}
func (*AuthController) Login(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.login")
}
func (*AuthController) Dologin(w http.ResponseWriter, r *http.Request) {

}
