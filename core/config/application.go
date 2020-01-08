package config

import (
	"encoding/json"
	"errors"
	"github.com/abeir/desktop-app/core"
	"github.com/gookit/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

func NewApplicationConfig() *ApplicationConfig{
	return &ApplicationConfig{load:Unload}
}

type ApplicationConfig struct {
	Environment string
	EnvironmentConfig
	//配置加载状态
	load LoadState
}


func (c *ApplicationConfig) findEnvironmentConfig(content *ConfigContent, env string) {
	for _, config := range content.Configurations {
		if config.Profile==env {
			c.EnvironmentConfig = config
			color.Printf("<light_green>global config:</> %+v\n", c.EnvironmentConfig)
			return
		}
	}
}

func (c *ApplicationConfig) loadFromJson(path string) error{
	color.Println("<light_green>load config from:</>", path)
	data, err := ioutil.ReadFile(path)
	if err!=nil {
		return err
	}
	content := &ConfigContent{}
	err = json.Unmarshal(data, content)
	if err!=nil {
		return err
	}
	c.Environment = content.Environment
	c.findEnvironmentConfig(content, content.Environment)
	return nil
}

func (c *ApplicationConfig) loadFromYml(path string) error{
	color.Println("<light_green>load config from:</>", path)
	data, err := ioutil.ReadFile(path)
	if err!=nil {
		return err
	}
	content := &ConfigContent{}
	err = yaml.Unmarshal(data, content)
	if err!=nil {
		return err
	}
	c.Environment = content.Environment
	c.findEnvironmentConfig(content, content.Environment)
	return nil
}

func (c *ApplicationConfig) currentPath() (string, error) {
	dir, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(dir), nil
}

//从环境变量中获取配置文件路径
func (c *ApplicationConfig) fileFromEnvVar() string{
	configPath := os.Getenv(applicationEnvVar)
	if configPath!="" {
		if core.IsExists(configPath) {
			return configPath
		}
	}
	return ""
}

func (c *ApplicationConfig) loadFromEnvVar() (bool, error){
	configFile := c.fileFromEnvVar()
	if configFile=="" {
		return false, nil
	}
	color.Println("<light_green>load config from environment variable:</>", applicationEnvVar, configFile)
	extName := filepath.Ext(configFile)
	if extName==".yaml" || extName==".yml" {
		return true, c.loadFromYml(configFile)
	}else if extName==".json" {
		return true, c.loadFromJson(configFile)
	}
	return false, nil
}

//Load 加载配置文件，先尝试从环境变量中的配置文件位置中读取，再尝试从程序所在位置下的config目录中读取
func (c *ApplicationConfig) Load() error{
	var err error
	if c.load == Loading {
		return errors.New("正在加载应用配置，请不要重复调用Load")
	}
	if c.load == Loaded {
		return errors.New("应用配置已加载，请不要重复调用Load")
	}
	c.load = Loading

	//返回前，根据err是否为空修改加载状态
	defer func() {
		if err!=nil {
			c.load = Unload
		}else{
			c.load = Loaded
		}
	}()

	//首先，尝试从环境变量中读取配置文件
	var isLoad bool
	isLoad, err = c.loadFromEnvVar()
	if err!=nil {
		return err
	}
	if isLoad {
		return nil
	}
	//从当前目录下读取
	var currentPath string
	currentPath, err = c.currentPath()
	if err!=nil {
		return err
	}
	configFile := filepath.Join(currentPath, "config", "application.yml")
	if core.IsExists(configFile) {
		if err = c.loadFromYml(configFile); err!=nil {
			return err
		}
	}
	configFile = filepath.Join(currentPath, "config", "application.yaml")
	if core.IsExists(configFile) {
		if err = c.loadFromYml(configFile); err!=nil {
			return err
		}
	}
	configFile = filepath.Join(currentPath, "config", "application.json")
	if core.IsExists(configFile) {
		if err = c.loadFromJson(configFile); err!=nil {
			return err
		}
	}
	err = errors.New("configuration file not found: " + configFile)
	return err
}

func (c *ApplicationConfig) IsDev() bool{
	return c.Environment=="dev"
}

func (c *ApplicationConfig) IsProd() bool{
	return c.Environment=="prod"
}

func (c *ApplicationConfig) IsTest() bool{
	return c.Environment=="test"
}

func (c *ApplicationConfig) Is(env string) bool{
	return c.Environment==env
}


type Database struct {
	Name string 	`json:"name" yaml:"name"`
	Url string 		`json:"url" yaml:"url"`
}

type Server struct {
	Port string 	`json:"port" yaml:"port"`
}

type Logger struct {
	Level string 		`json:"level" yaml:"level"`
	Path string 		`json:"path" yaml:"path"`
	Filename string 	`json:"filename" yaml:"filename"`
	MaxAge string 		`json:"maxAge" yaml:"maxAge"`
	RotationTime string 	`json:"rotationTime" yaml:"rotationTime"`
}

type EnvironmentConfig struct {
	Profile string		`json:"profile" yaml:"profile"`

	Database Database 	`json:"database" yaml:"database"`

	Server Server 		`json:"server" yaml:"server"`

	Logger Logger 		`json:"logger" yaml:"logger"`
}

type ConfigContent struct {
	Environment string 			`json:"environment" yaml:"environment"`
	Configurations []EnvironmentConfig		`json:"configurations" yaml:"configurations"`
}
