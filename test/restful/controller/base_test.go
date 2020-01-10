package controller

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http/httptest"
)

type BaseTest struct {
	engine *gin.Engine
	relativePath string
	handler gin.HandlerFunc
}

func NewBaseTest(relativePath string, handler gin.HandlerFunc) *BaseTest {
	return &BaseTest{
		engine: gin.Default(),
		relativePath: relativePath,
		handler: handler,
	}
}

func (b *BaseTest) DoRequest(method string, target string, body io.Reader) *httptest.ResponseRecorder{
	b.engine.GET(b.relativePath, b.handler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, target, body)
	b.engine.ServeHTTP(w, req)
	return w
}