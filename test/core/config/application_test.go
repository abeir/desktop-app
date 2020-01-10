package config

import (
	conf "github.com/abeir/desktop-app/core/config"
	"os"
	"testing"
)

func TestApplicationConfigLoad(t *testing.T) {
	_ = os.Setenv(conf.ApplicationEnvVar, "/home/abeir/workspace/go/desktop-app/config/application.yml")
	app := conf.NewApplicationConfig()
	if err := app.Load(); err!=nil {
		t.Error(err)
	}
	if app.Environment!= "dev" {
		t.Errorf("获取环境类型错误，预期：dev, 实际：%s", app.Environment)
	}
}