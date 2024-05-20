/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"webprobe/internal"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "webprobe",
	Short: "支持IPv4、IPv6的网站探测工具，	支持探测多级子链并计算支持度",
	Long: `支持IPv4、IPv6的网站探测工具，能够探测多级子链并计算支持度，同时提供网络性能优化、报告和可视化、灵活配置选项和安全性考虑。
	
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		internal.ScrapyWeb()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.webprobe.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&internal.ConfigPath, "config", "c", ".", "config file directory")
	//rootCmd.MarkPersistentFlagRequired("config") // 标记为必需
}
