package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var language string
// 回头参考下 https://github.com/gohade/hade/blob/main/app/console/kernel.go
var rootCmd = &cobra.Command{
	Use:   "gen",
	Short: "gen is code generate for go or java language",
	Long:  "gen 是一个生成go或java语言代码的工具",
	Run: func(cmd *cobra.Command, args []string) {
		// zap.L().Info("公用配置文件路径", zap.String("commonFile", commonFile))
		zap.L().Info("语言名称", zap.String("language", language))
		zap.L().Debug("收到", zap.Any("root-args", args))
	},
	// 不需要出现cobra默认的completion子命令
	// CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	initFlag()
	cobra.OnInitialize(initConfig)
	// rootCmd.AddCommand(serviceCmd)

}

func initConfig() {
}

func initFlag() {
	// 第4个参数为默认值
	rootCmd.PersistentFlags().StringVarP(&language, "lang", "l", "go", "请输入目标计算机语言名称")
	// rootCmd.MarkPersistentFlagRequired("lang") //必填
	// rootCmd.PersistentFlags().StringVarP(&commonFile, "commonFile", "c", "", "请输入源文件路径")
	// viper.SetDefault("author", "冯江涛 <hggfjt@163.com>")
}
