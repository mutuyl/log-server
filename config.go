package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config 日志文件信息结构体。
type Config struct {
	// 服务监听端口。
	Port int
	// 允许访问的IP地址(用 "," 号分隔)。
	ValidIPs string
}

var cfg *Config

// initConfig 初始化配置信息句柄。
func initConfig() {
	ReloadConfig()
}

// GetConfig 读取日志文件。
func GetConfig() *Config {
	if cfg == nil {
		t := &Config{}

		buf, err := ioutil.ReadFile("config.yml")
		if err != nil {
			log.Fatalf("read config.yml error: %s", err)
		}
		err = yaml.Unmarshal(buf, t)
		if err != nil {
			log.Fatalf("config.yml file error: %s", err)
		}
		cfg = t
	}
	return cfg
}

// ReloadConfig --
func ReloadConfig() {
	t := &Config{}
	buf, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("read config.yml error: %s", err)
	}
	err = yaml.Unmarshal(buf, t)
	if err != nil {
		log.Fatalf("config.yml file error: %s", err)
	}
	cfg = t
}
