package tests

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomePage(t *testing.T) { //单元测试
	baseURL := "http://localhost:8090"

	//请求 -- 模拟用户访问浏览器
	var tests = []struct {
		method   string
		url      string
		expected int
	}{
		{"GET", "/", 200},
		{"GET", "/about", 200},
		{"GET", "/notfound", 404},
		{"GET", "/aritcles", 200},
		{"GET", "/aritcles/create", 200},
		{"GET", "/aritcles/3", 200},
		{"GET", "/aritcles/3/edit", 200},
		{"POST", "/aritcles/3/edit", 200},
		{"POST", "/aritcles", 200},
		{"POST", "/aritcles/1/delete", 200},
	}

	for _, test := range tests {
		t.Logf("当前请求URL:%v \n", test.url)
		var (
			resp *http.Response
			err  error
		)
		switch {
		case test.method == "POST":
			data := make(map[string][]string)
			resp, err = http.PostForm(baseURL+test.url, data)
		default:
			resp, err = http.Get(baseURL + test.url)
		}
		//断言
		assert.NoError(t, err, "请求 "+test.url+" 时报错")
		assert.Equal(t, test.expected, resp.StatusCode, test.url+"应返回状态码"+strconv.Itoa(test.expected))
	}
}
