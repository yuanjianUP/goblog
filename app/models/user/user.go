package user

import (
	"goblog/app/models"
)

type User struct {
	models.BaseModel

	Name     string `gorm:"type:varchar(255);not null;unique" valid:"name"`
	Email    string `gorm:"type:varchar(255);not null;default'';unique" valid:"email"`
	Password string `gorm:"type:varchar(255);not null;default '';" valid:"password"`
	//gorm:"-" -- 设置gorm在读写时掠过次字段
	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}
