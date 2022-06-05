package aiticle

import "goblog/app/models/user"

type Article struct {
	UserID uint64 `gorm:"not null;index"`
	User   user.User
}
