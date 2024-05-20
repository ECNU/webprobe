package internal

import (
	"fmt"
	"log"
	"webprobe/models"

	"gorm.io/gorm/schema"
)

var ConfigPath string

func ScrapyWeb() {
	Cfg = GetConfig(ConfigPath)
	DB, _ = InitDB(Cfg)
	DB.NamingStrategy = schema.NamingStrategy{
		TablePrefix:   "",   // 表前缀，如果有的话
		SingularTable: true, // 禁用表名的复数形式
	}

	// 执行自动迁移
	err := DB.AutoMigrate(&models.Reachability{}, &models.URLStatus{})
	if err != nil {
		log.Fatalf("Error performing auto migration: %v", err)
	}
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}
	UrlDataWithReachability, _ = GeturlDataWithReachability(Cfg)
}
