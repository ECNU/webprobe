package scanner

import "time"

// URLDepth结构体用来存储URL和它的爬取深度
type URLDepth struct {
	URL       string
	Depth     int
	FatherURL string //父链接
}

// URLStatus存储单个URL的检查结果
type URLStatus struct {
	FatherURL  string
	Depth      int
	URL        string
	IPVersion  string //添加IP版本字段
	Up         bool
	StatusCode int
	Latency    time.Duration
	CertExpire time.Time
}

type URLReachabilityIPv4 struct {
	FirstLevelReach  float64
	SecondLevelReach float64
}

type URLReachabilityIPv6 struct {
	FirstLevelReach  float64
	SecondLevelReach float64
}
type URLDataWithReachability struct {
	URLStatus
	ReachabilityIPv4 URLReachabilityIPv4
	ReachabilityIPv6 URLReachabilityIPv6
}
