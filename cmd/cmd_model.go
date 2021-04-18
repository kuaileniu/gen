package cmd

import (
	"os"
	"path"

	// "runtime"
	"strings"

	"github.com/codeskyblue/kexec"
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
var Orm string // 数据库层使用的持久化框架
var OrmType consts.OrmType

var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "生成实体类对象",
	Long:  "暂时只支持生成go语言实体类对象",
	// Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// zap.L().Info("公用配置文件路径", zap.String("commonFile", commonFile))
		zap.L().Debug("语言名称", zap.String("language", language))
		zap.L().Debug("生成的model文件", zap.String("targetModelFile", targetModelFile))
		zap.L().Debug("实体定义源文件", zap.String("sourceModelFile", sourceModelFile))
		zap.L().Debug("PO是否同名于表名字段名", zap.Bool("modelFieldSameNameAsTable", modelFieldSameNameAsTable))
		zap.L().Debug("PO持久化框架", zap.String("Orm", Orm))
		// zap.L().Info("收到", zap.Any("cmd", cmd), zap.Any("args", args))
		zap.L().Debug("args", zap.Strings("model-args", args))
		getSourceFileType()
		getOrm()
		zap.L().Info("", zap.Bool("sameName", modelFieldSameNameAsTable))
		allInfo := parser.GetAllInfo(sourceModelFile, SourceFormat)
		// zap.L().Info("allInfo", zap.Reflect("allInfo", allInfo))
		allInfo.InferenceColumnDefaultTime()
		allInfo.CompatibleGoType()
		allInfo.InferenceColumnType()
		allInfo.InferencePropTypeIsKey()
		allInfo.InferenceColumnTypeRange()
		allInfo.InferenceOmitempty()
		allInfo.InferenceXormNotnull()
		allInfo.InferenceUnique()
		allInfo.InferenceJsonName()
		allInfo.InferenceXormDefault()
		allInfo.CollectImport()
		if modelFieldSameNameAsTable {
			allInfo.SetTableName()
			allInfo.SetColumnName()
		}
		switch language {
		case "go":
			allInfo.CreatePoModel(targetModelFile)
		default:
			zap.L().Error("暂不支持生成的语言源文件。", zap.String("language", language))
		}
		// TODO 生成完毕后 用代码 对文件再执行一次 go fmt
		cmdStr := "go fmt " + targetModelFile
		p := kexec.CommandString(cmdStr)
		p.Run()
		zap.L().Debug("格式化执行命令完毕", zap.String("cmdStr", cmdStr))
	},
}

// 获取orm类型
func getOrm() {
	orm := strings.TrimSpace(strings.ToLower(Orm))
	zap.L().Info("orm", zap.String("orm", orm))
	switch orm {
	case "xorm":
		OrmType = consts.Xorm
	case "gorm":
		OrmType = consts.Gorm
	case "mybatis":
		OrmType = consts.MyBatis
	}
	zap.L().Info("持久化框架", zap.Reflect("OrmType", OrmType))
}

// 判断模型文件的格式是json 或 yaml 或 ...
func getSourceFileType() {
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
	zap.L().Info("SourceFormat", zap.Reflect("SourceFormat", SourceFormat))
}

func init() {
	// 第4个参数为默认值
	rootCmd.AddCommand(modelCmd)
	modelCmd.Flags().StringVarP(&targetModelFile, "target", "t", "", "请输入实体类存储文件")
	modelCmd.MarkFlagRequired("target") // 必填
	modelCmd.Flags().StringVarP(&sourceModelFile, "source", "s", "", "请输入实体定义源文件")
	modelCmd.MarkFlagRequired("source") // 必填
	modelCmd.Flags().BoolVarP(&modelFieldSameNameAsTable, "modelFieldSameNameAsTable", "n", false, "PO是否同名于表名字段名")
	modelCmd.Flags().StringVarP(&sourceFileFormat, "sourceFileFormat", "f", "", "配置模型的文件类型,无值时根据文件后缀判断，例如json，yaml，yml")
	modelCmd.Flags().StringVarP(&Orm, "orm", "o", "xorm", "数据库持久化框架，默认xorm,例如 xorm,gorm,mybatis")
}
