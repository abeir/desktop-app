package config

import (
	"errors"
	"fmt"
	"github.com/abeir/desktop-app/core"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// NewApiConfig 创建ApiConfig
func NewApiConfig() *ApiConfig{
	return &ApiConfig{load:Unload}
}

type ApiConfig struct {
	Urls map[string]string  `json:"url" yaml:"url"`
	Apis []Api 	`json:"api" yaml:"api"`
	//加载配置状态
	load LoadState
}

type Api struct {
	Id string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
	Url string `json:"url" yaml:"url"`
	Get string `json:"get" yaml:"get"`
}

func (a *Api) IsEmpty() bool {
	return a.Id==""
}


func (a *ApiConfig) parseUrl() {
	if a.Apis==nil || len(a.Apis)==0{
		return
	}
	for _, api := range a.Apis {
		api.Url = core.NewTemplate().Parse(api.Url, a.Urls)
	}
}

func (a *ApiConfig) findYml() (string, error){
	configPath := os.Getenv(apiEnvVar)
	if configPath!="" {
		if core.IsExists(configPath) {
			return configPath, nil
		}
	}
	currentPath, err := core.CurrentPath()
	if err!=nil {
		return "", fmt.Errorf("获取当前路径失败：%w", err)
	}
	return filepath.Join(currentPath, "config", "api.yml"), nil
}

func (a *ApiConfig) loadFromYml(path string) error{
	fmt.Println("加载api.yml：", path)
	data, err := ioutil.ReadFile(path)
	if err!=nil {
		return fmt.Errorf("读取api.yml失败：%w", err)
	}
	err = yaml.Unmarshal(data, a)
	if err!=nil {
		return fmt.Errorf("解析api.yml失败：%w", err)
	}
	return nil
}

func (a *ApiConfig) Load() error{
	var err  error
	if a.load == Loading {
		return errors.New("正在加载api配置，请不要重复调用Load")
	}
	if a.load == Loaded {
		return errors.New("api配置已加载，请不要重复调用Load")
	}
	a.load = Loading
	defer func() {
		if err==nil {
			a.load = Loaded
		}else{
			a.load = Unload
		}
	}()
	file, err := a.findYml()
	if err!=nil {
		return err
	}
	if err=a.loadFromYml(file); err!=nil {
		return err
	}
	a.parseUrl()
	return nil
}

func (a *ApiConfig) IsEmpty() bool{
	return a.Apis==nil || len(a.Apis)==0
}
