package cmd

import (
	"os"
	"path"
	"strings"

	"github.com/kuaileniu/gen/consts"
	"github.com/kuaileniu/gen/parser"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var targetModelFile string         //新生成的model存放的路径
var sourceModelFile string         //model的源定义文件
var modelFieldSameNameAsTable bool //PO是否同名于表名字段名
// var showModel bool // 是否显示生成的model代码
var sourceFileFormat string // 配置模型的文件类型，例如json，yaml，yml,参考 SourceFormat
var SourceFormat consts.SourceFormat
var orm string // 数据库层使用的持久化框架
var OrmType consts.OrmType

var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "生成实体类对象",
	Long:  "暂时只支持生成go语言实体类对象",
	// Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// zap.L().Info("公用配置文件路径", zap.String("commonFile", commonFile))
		zap.L().Info("语言名称", zap.String("language", language))
		zap.L().Info("生成的model文件", zap.String("targetModelFile", targetModelFile))
		zap.L().Info("实体定义源文件", zap.String("sourceModelFile", sourceModelFile))
		zap.L().Info("PO是否同名于表名字段名", zap.Bool("modelFieldSameNameAsTable", modelFieldSameNameAsTable))
		// zap.L().Info("收到", zap.Any("cmd", cmd), zap.Any("args", args))
		zap.L().Info("args", zap.Strings("model-args", args))
		checkSourceFileType()
		allInfo := parser.GetAllInfo(sourceModelFile, SourceFormat)
		zap.L().Info("allInfo", zap.Reflect("allInfo", allInfo))
	},
}

// 判断模型文件的格式是json 或 yaml 或 ...
func checkSourceFileType() {
	switch sourceFileFormat {
	case "json":
		SourceFormat = consts.Json
	case "yaml", "yml":
		SourceFormat = consts.Yaml
	case "":
		ext := path.Ext(sourceModelFile)
		if strings.EqualFold(".yaml", ext) || strings.EqualFold(".yml", ext) {
			SourceFormat = consts.Yaml
		} else if strings.EqualFold(".json", ext) {
			SourceFormat = consts.Json
		} else {
			zap.L().Error("无法判断源文件类型")
			os.Exit(1)
		}
	}
	zap.L().Info(sourceModelFile+"的格式", zap.Reflect("SourceFormat", SourceFormat))
}

func init() {
	// 第4个参数为默认值
	rootCmd.AddCommand(modelCmd)
	modelCmd.Flags().StringVarP(&targetModelFile, "target", "t", "", "请输入实体类存储文件")
	modelCmd.MarkFlagRequired("target") // 必填
	modelCmd.Flags().StringVarP(&sourceModelFile, "source", "s", "", "请输入实体定义源文件")
	modelCmd.MarkFlagRequired("source") // 必填
	modelCmd.Flags().BoolVarP(&modelFieldSameNameAsTable, "modelFieldSameNameAsTable", "n", false, "PO是否同名于表名字段名")
	modelCmd.Flags().StringVarP(&sourceFileFormat, "sourceFileFormat", "f", "", "配置模型的文件类型，例如json，yaml，yml")
	modelCmd.Flags().StringVarP(&orm, "orm", "o", "xorm", "")
}
