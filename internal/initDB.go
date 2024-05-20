package internal

import (
	"fmt"
	"log"
	"webprobe/config"
	"webprobe/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库并返回一个 *gorm.DB 实例
func InitDB(config config.Config) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	switch config.DBConfig.Dialect {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBConfig.Username, config.DBConfig.Password, config.DBConfig.Host, config.DBConfig.Port, config.DBConfig.DBName)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", config.DBConfig.Host, config.DBConfig.Username, config.DBConfig.Password, config.DBConfig.DBName, config.DBConfig.Port, config.DBConfig.SSLMode)
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN: dsn,
		}), &gorm.Config{})
	// 可以添加其他数据库类型的处理逻辑
	default:
		log.Fatalf("Unsupported database dialect: %s", config.DBConfig.Dialect)
		return nil, fmt.Errorf("unsupported database dialect: %s", config.DBConfig.Dialect)
	}

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
		return nil, err
	}

	// 自动迁移表结构
	if err := db.AutoMigrate(&models.URLStatus{}, &models.Reachability{}); err != nil {
		log.Fatalf("Error auto migrating tables: %v", err)
		return nil, err
	}

	return db, nil
}
