package utility

import (
	"log"

	"webprobe/metrics"
	"webprobe/scanner"
)

func PushToPrometheus(pushGatewayURL string, urlDataWithReachability []scanner.URLDataWithReachability) {

	metrics.InitMetric()
	log.Printf("开始进行结果推送...")
	metrics.SaveMetrics(urlDataWithReachability)
	metrics.PushMetrics(pushGatewayURL)
	log.Printf("结果推送已完成，请访问 pushgateway 页面 %s 查看.", pushGatewayURL)

}
