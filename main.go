package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//articlesFormData创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

var router = mux.NewRouter().StrictSlash(true) //strictslash可以使go和go/都能正确访问
var db *sql.DB

func initDB() {
	var err error
	config := mysql.Config{
		User:                 "homestead",
		Passwd:               "secret",
		Addr:                 "127.0.0.1:33060",
		Net:                  "tcp",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}
	//准备数据库连减池
	db, err = sql.Open("mysql", config.FormatDSN())
	checkError(err)

	//设置最大连接数
	db.SetMaxOpenConns(25)
	//设置最大空闲连接数
	db.SetMaxIdleConns(25)
	//设置每个链接过期时间
	db.SetConnMaxLifetime(5 * time.Minute)

	//尝试连接，失败会报错
	err = db.Ping()
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>hello,欢迎来到goblog</h1>")
}
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客事用以记录编程笔记，如您有反馈或建议，请联系")
}
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+"<p>如有疑问，请联系我们</p>")
}

//声明存储数据库数据
type Article struct {
	Title, Body string
	ID          int64
}

func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

//表单验证
func validateArticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需要介于3-40"
	}

	//验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需要大于或等于10个字节"
	}
	return errors
}
func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	//2.读取对应文章列表
	article, err := getArticleByID(id)
	//如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			//3.1数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404文章未找到")
		} else {
			//数据库错误
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500服务器内部错误")
		}
	} else {
		//读取成功
		tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{
			"RouteName2URL": RouteName2URL,
			"Int64ToString": Int64ToString,
		}).ParseFiles("resources/views/articles/show.gohtml")
		checkError(err)
		tmpl.Execute(w, article)
	}
}

//通过路由名称来获取URL
func RouteName2URL(routeName string, pairs ...string) string {
	url, err := router.Get(routeName).URL(pairs...)
	if err != nil {
		checkError(err)
		return ""
	}
	return url.String()
}

//将int64转换为string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}
func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	//执行查询语句
	rows, err := db.Query("SELECT * FROM articles")
	checkError(err)
	defer rows.Close()

	var articles []Article
	//循环读取结果
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Title, &article.Body)
		checkError(err)
		//将article追加到articles的这个数据中
		articles = append(articles, article)
	}

	//检测遍历时是否发生错误
	err = rows.Err()
	checkError(err)

	//加载模板
	tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
	checkError(err)

	//渲染模板，将所有文章的数据传输进去
	tmpl.Execute(w, articles)
}
func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	errors := validateArticleFormData(title, body)
	//检查是否有错误
	if len(errors) == 0 {
		lastInsertID, err := saveArticleToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprint(w, "插入成功，ID为"+strconv.FormatInt(lastInsertID, 10))
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		storeURL, _ := router.Get("articles.store").URL()
		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}
		tmpl.Execute(w, data)
	}
}

//插入文章
func saveArticleToDB(title string, body string) (int64, error) {
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)
	//1.获取一个prepare声明语句
	stmt, err = db.Prepare("INSERT INTO articles (title,body) VALUES(?,?)")
	//例行的错误检测
	if err != nil {
		return 0, err
	}
	//2.在此函数运行结束后关闭此语句，防止占用sql连接
	defer stmt.Close() //defer延迟语句
	//3.执行请求，传参进入绑定的内容
	rs, err = stmt.Exec(title, body)
	if err != nil {
		return 0, err
	}
	//4.插入成功的话，会返回自增ID
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}
	return 0, err
}
func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	html := `
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <title>创建文章 —— 我的技术博客</title>
        </head>
        <body>
            <form action="%s" method="post">
                <p><input type="text" name="title"></p>
                <p><textarea name="body" cols="30" rows="10"></textarea></p>
                <p><button type="submit">提交</button></p>
            </form>
        </body>
        </html>
    `
	storeURL, _ := router.Get("articles.store").URL()
	fmt.Fprintf(w, html, storeURL)
}

//中间件
func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//1.设置标头
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		//2.继续处理请求
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/") //去掉右边的/
		}
		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler { //移除路由后的/符号
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		//请求传递下去
		next.ServeHTTP(w, r)
	})
}
func createTable() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
		id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
		title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
		body longtext COLLATE utf8mb4_unicode_ci
	);`
	_, err := db.Exec(createArticlesSQL)
	checkError(err)
}

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
	//获取参数
	vars := mux.Vars(r)
	id := vars["id"]
	//读取对应的数据
	articles := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&articles.ID, &articles.Title, &articles.Body)
	if err != nil {
		if err == sql.ErrNoRows {
			//未找到数据
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404未找到数据")
		} else {
			//数据库错误
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		//读取成功，显示表单
		updateURL, _ := router.Get("articles.update").URL("id", id)
		data := ArticlesFormData{
			Title:  articles.Title,
			Body:   articles.Body,
			URL:    updateURL,
			Errors: nil,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		checkError(err)
		tmpl.Execute(w, data)
	}
}
func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			checkError(err)
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 未找到数据")
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 数据库错误")
		}
	} else {
		//表单验证
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")
		errors := validateArticleFormData(title, body)

		if len(errors) == 0 {
			//验证通过
			query := "UPDATE articles SET title = ?,body = ? WHERE id = ?"
			rs, err := db.Exec(query, title, body, id)

			if err != nil {
				checkError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500服务器内部错误")
			}

			//更新成功，跳转到文章详情页
			if n, _ := rs.RowsAffected(); n > 0 {
				showURL, _ := router.Get("articles.show").URL("id", id)
				http.Redirect(w, r, showURL.String(), http.StatusFound)
			} else {
				fmt.Fprint(w, "您没有做任何更改")
			}
		} else {
			updateURL, _ := router.Get("articles.update").URL("id", id)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: errors,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			checkError(err)
			tmpl.Execute(w, data)
		}
	}
}
func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	id := vars["id"]
	//读取对应的文章数据
	article, err := getArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			checkError(err)
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 未找到数据")
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		}
	} else {
		rowsAffected, err := article.Delete()
		if err != nil {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器错误")
		} else {
			//未发生错误
			if rowsAffected > 0 {
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}
}
func (a Article) Link() string { //语言函数方法
	showURL, err := router.Get("articles.show").URL("id", strconv.FormatInt(a.ID, 10))
	if err != nil {
		checkError(err)
		return ""
	}
	return showURL.String()
}
func (a Article) Delete() (RowsAffected int64, err error) {
	rs, err := db.Exec("DELETE FROM articles WHERE id =" + strconv.FormatInt(a.ID, 10))
	if err != nil {
		return 0, err
	}
	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}
	return 0, nil
}
func main() {
	initDB()
	createTable()
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")

	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}/update", articlesUpdateHandler).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	//使用中间件：强制内容类型为HTML
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
