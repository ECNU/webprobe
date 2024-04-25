package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	URLs          []string      `mapstructure:"url"`
	ScannerConfig ScannerConfig `mapstructure:"scanner"`
	Push          Push          `mapstructure:"push"`
	Output        string        `mapstructure:"output"` // default value: "output.html"
	DBConfig      DBConfig      `mapstructure:"db"`
}

type ScannerConfig struct {
	Crawl CrawlConfig `mapstructure:"crawl"`
	Check CheckConfig `mapstructure:"check"`
}

type CrawlConfig struct {
	Depth   int           `mapstructure:"depth"`   // default value: 2
	Timeout time.Duration `mapstructure:"timeout"` // default value: 3 (seconds)
}

type CheckConfig struct {
	UseIPV6           bool          `mapstructure:"ipv6"` // default value: false
	Retry             ReTryConfig   `mapstructure:"retry"`
	Concurrency       int           `mapstructure:"concurrency"`         // default value: 1000
	DialerTimeout     time.Duration `mapstructure:"dialer_timeout"`      // default value: 10 (seconds)
	HttpClientTimeout time.Duration `mapstructure:"http_client_timeout"` // default value: 15
	ContextTimeout    time.Duration `mapstructure:"context_timeout"`     // default value: 1000
}
type ReTryConfig struct {
	Time     int           `mapstructure:"time"`
	Interval time.Duration `mapstructure:"interval"`
}
type Push struct {
	PushToPrometheus bool   `mapstructure:"push_to_prometheus"` // default value: false
	PushGatewayURL   string `mapstructure:"push_gateway_url"`
}
type DBConfig struct {
	UseDB  bool   `mapstructure:"use_db"`  // default value: true
	DBPath string `mapstructure:"db_path"` // default value: "url_status.db"
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv() // read from environment variables

	// 设置默认值
	viper.SetDefault("output", "output.html")
	viper.SetDefault("scanner.crawl.depth", 2)
	viper.SetDefault("scanner.crawl.timeout", 3*time.Second)
	viper.SetDefault("scanner.check.ipv6", false)
	viper.SetDefault("scanner.check.retry.time", 1)
	viper.SetDefault("scanner.check.retry.interval", 1*time.Second)
	viper.SetDefault("scanner.check.concurrency", 1000)
	viper.SetDefault("scanner.check.dialer_timeout", 10*time.Second)
	viper.SetDefault("scanner.check.dialer_alive_timeout", 5*time.Second)
	viper.SetDefault("scanner.check.http_client_timeout", 15*time.Second)
	viper.SetDefault("scanner.check.context_timeout", 1000*time.Second)
	viper.SetDefault("push.push_to_prometheus", false)
	viper.SetDefault("db.use_db", false)
	viper.SetDefault("db.db_path", "LinkPulse.db")

	var configuration Config

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	err := viper.Unmarshal(&configuration)
	return &configuration, err
}
