package controller

import (
	"encoding/json"
	ctlr "github.com/abeir/desktop-app/restful/controller"
	"github.com/abeir/desktop-app/restful/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestTestController(t *testing.T) {
	testController := ctlr.NewTestController()
	w := NewBaseTest("/test", testController.Index).DoRequest("GET", "/test", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	body := w.Body
	assert.NotEmpty(t, body.String())

	rs := &model.ResultMessage{}
	assert.NoError(t, json.Unmarshal(body.Bytes(), rs), "解析body json格式错误：" + body.String())
	assert.Equal(t, model.SuccessCode, rs.Code)
}
