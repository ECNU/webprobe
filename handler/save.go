package handler

import (
	"log"
	"time"
	"webprobe/scanner"

	"webprobe/models"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

/* // 将URL状态数据保存到数据库
func SaveData(db *sql.DB, urlDataWithReachability []scanner.URLDataWithReachability, maxDepth int) {
	const layout = "2006-01-02 15:04:05-07:00" // 设置时间格式

	insertStatusSQL := `REPLACE INTO url_status(FatherURL, Depth, URL, IPVersion, Up, StatusCode, Latency, CertExpire, CreateTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	insertReachabilitySQL := `REPLACE INTO reachability(URL, IPv4FirstLevelReach, IPv4SecondLevelReach, IPv6FirstLevelReach, IPv6SecondLevelReach, CreateTime) VALUES (?, ?, ?, ?, ?, ?)`

	for _, data := range urlDataWithReachability {
		currentTime := time.Now().Format(layout) // 格式化当前时间
		status := data.URLStatus                 // 从URLDataWithReachability中获取URLStatus
		_, err := db.Exec(insertStatusSQL, status.FatherURL, status.Depth, status.URL, status.IPVersion, status.Up, status.StatusCode, status.Latency.Milliseconds(), status.CertExpire, currentTime)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
		}
		if data.URLStatus.Depth < maxDepth {
			_, err = db.Exec(insertReachabilitySQL, data.URL, data.ReachabilityIPv4.FirstLevelReach, data.ReachabilityIPv4.SecondLevelReach, data.ReachabilityIPv6.FirstLevelReach, data.ReachabilityIPv6.SecondLevelReach, time.Now())
			if err != nil {
				log.Printf("Error inserting reachability data: %v", err)
			}
		}
	}
}
*/

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
