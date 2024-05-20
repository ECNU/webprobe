package internal

import (
	"log"

	"webprobe/config"
	"webprobe/handler"
	"webprobe/scanner"
)

var UrlDataWithReachability []scanner.URLDataWithReachability

// geturlDataWithReachability
func GeturlDataWithReachability(cfg config.Config) ([]scanner.URLDataWithReachability, error) {

	// scanner结果的接收数组

	log.Printf("开始进行扫描...")
	urlDataWithReachability := make([]scanner.URLDataWithReachability, 0)
	urlDataWithReachability = handler.Scanner(cfg.URLs, cfg.ScannerConfig, urlDataWithReachability)
	return urlDataWithReachability, nil
}
