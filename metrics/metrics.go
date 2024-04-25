package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"log"
	"strconv"
	"webprobe/scanner"
)

var (
	Depth = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sample_depth",
		Help: "The depth of the URL",
	}, []string{"url", "ip_version", "depth"})
	Up = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sample_up",
		Help: "Is the URL up",
	}, []string{"url", "ip_version", "depth"})

	StatusCode = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sample_status_code",
		Help: "The status code for the URL",
	}, []string{"url", "ip_version", "depth"})

	Latency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sample_latency",
		Help: "The latency for the URL",
	}, []string{"url", "ip_version", "depth"})

	CertExpire = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sample_cert_expire",
		Help: "The cert expire time for the URL",
	}, []string{"url", "ip_version", "depth"})

	FirstLevelReach = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "first_level_reach",
		Help: "First level reachability",
	}, []string{"url", "ip_version", "depth"})
	SecondLevelReach = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "second_level_reach",
		Help: "Second level reachability",
	}, []string{"url", "ip_version", "depth"})
)

func InitMetric() {
	prometheus.MustRegister(Depth, Up, StatusCode, Latency, CertExpire, FirstLevelReach, SecondLevelReach)
}
func SaveMetrics(urlDataWithReachability []scanner.URLDataWithReachability) {
	for _, data := range urlDataWithReachability {
		status := data.URLStatus // 从URLDataWithReachability中获取URLStatus
		labelValues := []string{status.URL, status.IPVersion, strconv.Itoa(status.Depth)}
		Depth.WithLabelValues(labelValues...).Set(float64(status.Depth))
		// 如果 status 为假(false)，则设置为 0.0，如果为真(true)，则设置为 1.0
		Up.WithLabelValues(labelValues...).Set(0.0)
		if status.Up {
			Up.WithLabelValues(labelValues...).Set(1.0)
		}
		StatusCode.WithLabelValues(labelValues...).Set(float64(status.StatusCode))
		Latency.WithLabelValues(labelValues...).Set(status.Latency.Seconds())
		CertExpire.WithLabelValues(labelValues...).Set(float64(status.CertExpire.Unix()))

		if status.IPVersion == "ipv4" {
			FirstLevelReach.WithLabelValues(labelValues...).Set(data.ReachabilityIPv4.FirstLevelReach)
			SecondLevelReach.WithLabelValues(labelValues...).Set(data.ReachabilityIPv4.SecondLevelReach)
		}
		if status.IPVersion == "ipv6" {
			FirstLevelReach.WithLabelValues(labelValues...).Set(data.ReachabilityIPv6.FirstLevelReach)
			SecondLevelReach.WithLabelValues(labelValues...).Set(data.ReachabilityIPv6.SecondLevelReach)
		}
	}
}
func PushMetrics(pushgatewayURL string) {
	if err := push.New(pushgatewayURL, "http_response"). // push.New("pushgateway地址", "job名称")
								Collector(Depth). // Collector(completionTime) 给指标赋值
								Grouping("type", "depth").
								Push(); err != nil {
		log.Fatalf("Could not push to Pushgateway: %v", err)
	}
	if err := push.New(pushgatewayURL, "http_response"). // push.New("pushgateway地址", "job名称")
								Collector(StatusCode). // Collector(completionTime) 给指标赋值
								Grouping("type", "statusCode").
								Push(); err != nil {
		log.Fatalf("Could not push to Pushgateway: %v", err)
	}
	if err := push.New(pushgatewayURL, "http_response"). // push.New("pushgateway地址", "job名称")
								Collector(Up). // Collector(completionTime) 给指标赋值
								Grouping("type", "up").
								Push(); err != nil {
		log.Fatalf("Could not push to Pushgateway: %v", err)
	}
	if err := push.New(pushgatewayURL, "http_response"). // push.New("pushgateway地址", "job名称")
								Collector(Latency). // Collector(completionTime) 给指标赋值
								Grouping("type", "latency").
								Push(); err != nil {
		log.Fatalf("Could not push to Pushgateway: %v", err)
	}
	if err := push.New(pushgatewayURL, "http_response"). // push.New("pushgateway地址", "job名称")
								Collector(CertExpire). // Collector(completionTime) 给指标赋值
								Grouping("type", "certExpire").
								Push(); err != nil {
		log.Fatalf("Could not push to Pushgateway: %v", err)
	}
	if err := push.New(pushgatewayURL, "http_response"). // push.New("pushgateway地址", "job名称")
								Collector(FirstLevelReach). // Collector(completionTime) 给指标赋值
								Grouping("type", "firstLevelReach").
								Push(); err != nil {
		log.Fatalf("Could not push to Pushgateway: %v", err)
	}
	if err := push.New(pushgatewayURL, "http_response"). // push.New("pushgateway地址", "job名称")
								Collector(SecondLevelReach). // Collector(completionTime) 给指标赋值
								Grouping("type", "secondLevelReach").
								Push(); err != nil {
		log.Fatalf("Could not push to Pushgateway: %v", err)
	}
}
