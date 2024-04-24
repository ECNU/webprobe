package handler

import (
	"log"
	"webprobe/config"
	"webprobe/scanner"
)

func Scanner(URLs []string, config config.ScannerConfig, urlDataWithReachability []scanner.URLDataWithReachability) []scanner.URLDataWithReachability {
	// 用于存储爬虫结果，包括url、深度、父链接
	urlWithDepth := make([]scanner.URLDepth, 0)
	for index, url := range URLs {
		log.Printf("开始进行第 %d/%d 个 链接爬虫...", index+1, len(URLs))
		urlWithDepth = scanner.Crawler(url, config.Crawl, urlWithDepth)
		log.Printf("链接爬虫已完成，共抓取到 %d 条链接. ", len(urlWithDepth))
	}

	// 用于存储检测结果，包括FatherURL Depth URL IPVersion Up StatusCode Latency CertExpire
	urlData := make([]scanner.URLStatus, 0)
	log.Printf("开始进行链接检测...")
	urlData = scanner.Check(urlWithDepth, config.Check, urlData)
	log.Printf("链接检测已完成.")

	log.Printf("开始进行子链可达率计算...")
	urlDataWithReachability = scanner.CalculateReachability(urlData, config.Check.UseIPV6, urlDataWithReachability)
	log.Printf("子链可达率计算已完成.")
	return urlDataWithReachability
}
