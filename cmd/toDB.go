/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"webprobe/handler"
	"webprobe/internal"
	"webprobe/utility"

	"github.com/spf13/cobra"
)

// toDBCmd represents the toDB command
var toDBCmd = &cobra.Command{
	Use:   "toDB",
	Short: "Specify the output database type (relational or time-series). Default is time-series.",
	Long: `Specify the type of database to save the scan results. The available options are:

	-- RDB: Save the scan results to a relational database.
	-- TSDB: Save the scan results to a time-series database (default).`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.ScrapyWeb()

		RDBItems, _ := cmd.Flags().GetStringArray("RDB")
		TSDBItems, _ := cmd.Flags().GetStringArray("TSDB")
		//for _, item := range RDBItems {
		if len(RDBItems) > 0 {

			if internal.Cfg.DBConfig.UseDB {
				log.Printf("正在写入数据库...")
				handler.SaveData(internal.DB, internal.UrlDataWithReachability, internal.Cfg.ScannerConfig.Crawl.Depth)
				log.Printf("写入关系数据库" + RDBItems[0] + "完成.")
			}
		}

		//for _, item := range TSDBItems {
		if len(TSDBItems) > 0 {

			utility.PushToPrometheus(internal.Cfg.Push.PushGatewayURL, internal.UrlDataWithReachability)
			log.Printf("写入时序数据库" + TSDBItems[0] + "完成.")
		}

	},
}

func init() {

	rootCmd.AddCommand(toDBCmd)
	// 在这里定义您的标志和配置设置。

	// Cobra支持持久性标志，这些标志将适用于此命令
	// 和所有子命令，例如：
	// toDBCmd.PersistentFlags().String("foo", "", "A help for foo")
	// 创建toDBCmd命令的Flag，作用为指定关系性数据库的名称 (默认为tsdb)
	toDBCmd.PersistentFlags().StringArrayP("RDB", "r", []string{""}, "result to RDB")
	// 创建toDBCmd命令的Flag，作用为指定时序数据库的名称 (默认为tsdb)
	toDBCmd.PersistentFlags().StringArrayP("TSDB", "t", []string{""}, "result to toTSDB")

	// Cobra支持仅在直接调用此命令时运行的本地标志，例如：
	// toDBCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

//
//
