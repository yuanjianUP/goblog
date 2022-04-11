package article

import (
	"goblog/pkg/model"
	"goblog/pkg/types"
)

//声明存储数据库数据
type Article struct {
	Title, Body string
	ID          uint64
}

func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToUint64(idstr)
	if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}
	return article, nil
}
