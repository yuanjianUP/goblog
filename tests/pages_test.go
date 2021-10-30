package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomePage(t *testing.T) { //单元测试
	baseURL := "http://localhost:3000"

	//请求 -- 模拟用户访问浏览器
	var (
		resp *http.Response
		err  error
	)
	resp, err = http.Get(baseURL + "/")
	assert.NoError(t, err, "有错误发生,err不为空")
	assert.Equal(t, 200, resp.StatusCode, "应用返回状态 200")
}
