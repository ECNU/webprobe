package utility

import (
	"log"
	"time"
	"webprobe/scanner"

	"webprobe/models"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

// SaveData 保存 URL 状态数据到数据库
func SaveData(db *gorm.DB, urlDataWithReachability []scanner.URLDataWithReachability, maxDepth int) {
	currentTime := time.Now()

	for _, data := range urlDataWithReachability {
		status := data.URLStatus

		// 保存 URLStatus 数据
		if err := db.Create(&status).Error; err != nil {
			log.Printf("Error inserting URLStatus data: %v", err)
		}

		// 如果深度小于最大深度，则保存 Reachability 数据
		if status.Depth < maxDepth {
			reachability := models.Reachability{
				URL:                  data.URL,
				IPv4FirstLevelReach:  data.ReachabilityIPv4.FirstLevelReach,
				IPv4SecondLevelReach: data.ReachabilityIPv4.SecondLevelReach,
				IPv6FirstLevelReach:  data.ReachabilityIPv6.FirstLevelReach,
				IPv6SecondLevelReach: data.ReachabilityIPv6.SecondLevelReach,
				CreateTime:           currentTime,
			}

			if err := db.Create(&reachability).Error; err != nil {
				log.Printf("Error inserting Reachability data: %v", err)
			}
		}
	}
}
