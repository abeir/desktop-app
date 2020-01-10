package config

import (
	conf "github.com/abeir/desktop-app/core/config"
	"os"
	"testing"
)

func TestApiConfigLoad(t *testing.T) {
	_ = os.Setenv(conf.ApiEnvVar, "/home/abeir/workspace/go/desktop-app/config/api.yml")
	api := conf.NewApiConfig()
	if err := api.Load(); err!=nil {
		t.Error(err)
		return
	}
	if api.IsEmpty() {
		t.Error("加载api.yml内容为空")
		return
	}
	if api.Urls==nil || len(api.Urls)==0 {
		t.Error("读取api.yml中urls节点内容为空")
		return
	}
	url := "https://kyfw.12306.cn"
	if api.Urls["u12306"] != url {
		t.Errorf("读取api.yml中urls节点u12306内容错误，预期：%s, 实际：%s", url, api.Urls["u12306"])
		return
	}
	sid := "station_name"
	if api.Apis[0].Id != sid {
		t.Errorf("读取api.yml中api节点id内容错误，预期：%s, 实际：%s", sid, api.Apis[0].Id)
		return
	}
	surl := url + "/otn/resources/js/framework/station_name.js"
	if api.Apis[0].Url != surl {
		t.Errorf("读取api.yml中api节点url内容错误，预期：%s, 实际：%s", surl, api.Apis[0].Url)
		return
	}
}
