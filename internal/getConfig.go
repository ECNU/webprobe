package internal

import (
	"log"
	"webprobe/config"
)

var Cfg config.Config

func GetConfig(configPath string) config.Config {
	log.Printf("开始导入配置文件...")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("导入配置文件失败: %v", err)
	}
	log.Printf("配置文件导入成功.")
	return *cfg
}
