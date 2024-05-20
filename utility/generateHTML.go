package utility

import (
	"log"

	"webprobe/config"
	"webprobe/handler"

	"webprobe/internal"
)

func generateHTML() {
	cfg := internal.GetConfig(config.ConfigPath)
	urlDataWithReachability, _ := internal.GeturlDataWithReachability(cfg)
	log.Printf("扫描已完成.")
	log.Printf("开始生成 HTML 文件...")
	handler.GenerateHTMLFile(urlDataWithReachability, cfg.Output, cfg.ScannerConfig.Crawl.Depth)
	log.Printf("HTML 文件已生成，请访问 %s 查看.", cfg.Output)
}
