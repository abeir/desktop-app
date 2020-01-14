package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewIndexController() *IndexController {
	return &IndexController{}
}

type IndexController struct {

}

func (i *IndexController) Index(ct *gin.Context){
	ct.HTML(http.StatusOK, "index.html", nil)
}
