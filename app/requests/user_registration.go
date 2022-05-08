package requests

import (
	"goblog/app/models/user"

	"github.com/thedevsaddam/govalidator"
)

//表单验证
func ValidateRegistrationFrom(data user.User) map[string][]string {
	//定制规则
	rules := govalidator.MapData{
		"name":             []string{"required", "alpha_num", "between:3,20", "not_exists:users,name"},
		"email":            []string{"required", "min:4", "max:30", "email", "not_exists:users,email"},
		"password":         []string{"required", "min:6"},
		"password_confirm": []string{"required"},
	}
	//定制错误消息
	message := govalidator.MapData{
		"name": []string{
			"required:用户名为必填项",
			"alpha_num:格式错误，只允许数字和英文",
			"between:用户名长度需在3-20之间",
		},
		"email": []string{
			"required:Email为必填项",
			"min:Email必须大鱼4",
			"max:Email必须小雨30",
			"email:Email格式不正确",
		},
		"password": []string{
			"required:密码为必填项",
			"min:长度需大于 6",
		},
		"password_confirm": []string{
			"required:确认密码为必填项",
		},
	}
	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid",
		Messages:      message,
	}
	errs := govalidator.New(opts).ValidateStruct()
	if data.Password != data.PasswordConfirm {
		errs["PasswordConfirm"] = append(errs["PasswordConfirm"], "两次密码输入不一致")
	}
	return errs
}
