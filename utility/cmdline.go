package utility

import (
	"webprobe/config"
	"webprobe/internal"

	"github.com/gookit/goutil/dump"
)

func Cmdline() {
	cfg := internal.GetConfig(config.ConfigPath)
	urlDataWithReachability, _ := internal.GeturlDataWithReachability(cfg)
	/*options := dump.Options{
		Output:     os.Stdout,
		NoType:     true,
		NoColor:    false,
		IndentLen:  2,
		IndentChar: ' ',
		MaxDepth:   10,
		ShowFlag:   0,
		CallerSkip: 1,
	}
	*/
	// 使用这些选项创建一个新的 dumper
	d := dump.NewWithOptions(dump.WithoutType())

	// 使用新的 dumper 打印变量
	d.Print(urlDataWithReachability)
	/*for _, data := range urlDataWithReachability {
		//currentTime := time.Now().Format(layout) // 格式化当前时间
		status := data.URLStatus // 从URLDataWithReachability中获取URLStatus
		fmt.Printf("ipv4: %s,ipv6:%s, 状态: %s \n", data.ReachabilityIPv4, data.ReachabilityIPv6, status)

	}*/
}
