package article

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"strconv"
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
func GetAll() ([]Article, error) {
	var articles []Article
	if err := model.DB.Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}
func (article *Article) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatUint(article.ID, 10))
}

//create创建文章，通过article.ID来判断是否创建成功
func (article *Article) Create() (err error) {
	if err = model.DB.Create(&article).Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}
