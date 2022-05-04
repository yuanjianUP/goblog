package bootstrap

import (
	"goblog/app/models/user"
	"goblog/pkg/model"
	"goblog/pkg/model/article"
	"time"

	"gorm.io/gorm"
)

//初始化数据库ORM
func SetupDB() {
	//建立连接
	db := model.ConnectDB()
	sqlDB, _ := db.DB()
	//设置最大连接数
	sqlDB.SetMaxOpenConns(100)
	//设置最大空闲连接数
	sqlDB.SetMaxIdleConns(25)
	//设置每个连接的过期时间
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	migration(db)
}

func migration(db *gorm.DB) {
	db.AutoMigrate(
		&user.User{},

		&article.Article{},
	)
}
