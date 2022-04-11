package bootstrap

import (
	"goblog/pkg/model"
	"time"
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
}
