package main

import (
	"log"
	"webprobe/config"
	"webprobe/handler"
	"webprobe/metrics"
	"webprobe/scanner"
)

func main() {
	log.Printf("开始导入配置文件...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("导入配置文件失败: %v", err)
	}
	log.Printf("配置文件导入成功.")
	// scanner结果的接收数组
	metrics.InitMetric()
	log.Printf("开始进行扫描...")
	urlDataWithReachability := make([]scanner.URLDataWithReachability, 0)
	urlDataWithReachability = handler.Scanner(cfg.URLs, cfg.ScannerConfig, urlDataWithReachability)
	log.Printf("扫描已完成.")
	log.Printf("开始生成 HTML 文件...")
	handler.GenerateHTMLFile(urlDataWithReachability, cfg.Output, cfg.ScannerConfig.Crawl.Depth)
	log.Printf("HTML 文件已生成，请访问 %s 查看.", cfg.Output)

	if cfg.DBConfig.UseDB {
		log.Printf("正在写入数据库...")
		db := handler.InitDB(cfg.DBConfig.DBPath)
		handler.SaveData(db, urlDataWithReachability, cfg.ScannerConfig.Crawl.Depth)
		defer db.Close()
		log.Printf("写入数据库完成.")
	}

	if cfg.Push.PushToPrometheus {
		log.Printf("开始进行结果推送...")
		metrics.SaveMetrics(urlDataWithReachability)
		metrics.PushMetrics(cfg.Push.PushGatewayURL)
		log.Printf("结果推送已完成，请访问 pushgateway 页面 %s 查看.", cfg.Push.PushGatewayURL)
	}
}
