package service

import (
	"github.com/abeir/desktop-app/core/config"
	"github.com/abeir/desktop-app/restful"
)

type BaseService struct {
	api config.Api
}

// FindUrl 从api配置中，根据id获取url
func (b *BaseService) FindUrl(id string) string{
	if id=="" {
		return ""
	}
	if !b.api.IsEmpty() {
		return b.api.Url
	}
	apis := restful.Gobal.Api.Apis
	for _, api := range apis {
		if id == api.Id {
			b.api = api
			return api.Url
		}
	}
	return ""
}

func (b *BaseService) request(){

}
