package handler

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
	"webprobe/scanner"
)

func InitDB(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// 创建第一个表：url_status
	createURLStatusTableSQL := `
    CREATE TABLE IF NOT EXISTS url_status (
        FatherURL TEXT,
        Depth INTEGER,
        URL TEXT,
        IPVersion TEXT,
        Up BOOLEAN,
        StatusCode INTEGER,
        Latency INTEGER,
        CertExpire DATETIME,
        CreateTime DATETIME,
        PRIMARY KEY (URL, IPVersion, CreateTime)
	);`
	_, err = db.Exec(createURLStatusTableSQL)
	if err != nil {
		log.Fatalf("Error creating url_status table: %v", err)
	}

	// 创建第二个表：reachability
	createReachabilityTableSQL := `
    CREATE TABLE IF NOT EXISTS reachability (
        URL TEXT,
        IPv4FirstLevelReach FLOAT,
        IPv4SecondLevelReach FLOAT,
        IPv6FirstLevelReach FLOAT,
        IPv6SecondLevelReach FLOAT,
        CreateTime DATETIME,
        PRIMARY KEY (URL, CreateTime)
    );`
	_, err = db.Exec(createReachabilityTableSQL)
	if err != nil {
		log.Fatalf("Error creating reachability table: %v", err)
	}

	return db
}

// 将URL状态数据保存到数据库
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
