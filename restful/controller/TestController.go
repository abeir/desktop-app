package controller

import (
	"github.com/abeir/desktop-app/restful/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewTestController() *TestController {
	return &TestController{}
}

type TestController struct {

}

func (t *TestController) Index(ct *gin.Context){
	rs := model.SuccessResultMessage("success")
	ct.JSON(http.StatusOK, rs)
}


