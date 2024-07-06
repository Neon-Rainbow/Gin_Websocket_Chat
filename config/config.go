package config

import (
	"encoding/json"
	"log"
	"os"
)

type MySQL struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type Config struct {
	MySQL   MySQL  `json:"mysql"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

// AppConfig 用于存储配置文件的内容
var AppConfig *Config

// LoadConfig 用于加载配置文件
func LoadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("无法打开配置文件: %v", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("关闭文件失败: %v", err)
		}
	}(file)

	AppConfig = &Config{}
	if err := json.NewDecoder(file).Decode(AppConfig); err != nil {
		log.Fatalf("无法解析配置文件: %v", err)
		return err
	}

	return nil
}
