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

var po_target string           //新生成的model存放的路径
var po_source string           //model的源定义文件，yml文件等
var controller_target string   //新生成的controller存放的路径
// var controller_source string   //controller的源定义文件,yml文件等
var po_same_name_as_table bool //PO是否同名于表名字段名
// var showModel bool // 是否显示生成的model代码
// var sourceFileFormat string // 配置模型的文件类型，例如json，yaml，yml,参考 SourceFormat
// var SourceFormat consts.SourceFormat
var Orm string // 数据库层使用的持久化框架
var OrmType consts.OrmType
var JsonCase string

var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "生成实体类对象",
	Long:  "暂时只支持生成go语言实体类对象",
	// Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// zap.L().Info("公用配置文件路径", zap.String("commonFile", commonFile))
		zap.L().Debug("语言名称", zap.String("language", language))
		zap.L().Debug("生成的 model 文件", zap.String("po-target", po_target))
		zap.L().Debug("生成的 controller 文件", zap.String("controller-target", controller_target))
		zap.L().Debug("PO定义源文件", zap.String("po-source", po_source))
		zap.L().Debug("PO是否同名于表名字段名", zap.Bool("po-same-name-as-table", po_same_name_as_table))
		zap.L().Debug("PO持久化框架", zap.String("Orm", Orm))
		// zap.L().Info("收到", zap.Any("cmd", cmd), zap.Any("args", args))
		zap.L().Debug("args", zap.Strings("model-args", args))
		SourceFormat := getSourceFileType()
		getOrm()
		zap.L().Info("", zap.Bool("sameName", po_same_name_as_table))
		allInfo := parser.GetAllInfo(po_source, SourceFormat)
		// zap.L().Info("allInfo", zap.Reflect("allInfo", allInfo))
		allInfo.InferenceColumnDefaultTime()
		allInfo.CompatibleGoType()
		allInfo.InferenceColumnType()
		allInfo.InferencePropTypeIsKey()
		allInfo.InferenceColumnTypeRange()
		allInfo.InferenceOmitempty()
		allInfo.InferenceXormNotnull()
		allInfo.InferenceUnique()
		allInfo.InferenceJsonName(JsonCase)
		allInfo.InferenceXormDefault()
		allInfo.CollectImport()
		allInfo.CollectControllerImport()
		if po_same_name_as_table {
			allInfo.SetTableName()
			allInfo.SetColumnName()
		}
		switch language {
		case "go":
			allInfo.CreatePoModel(po_target)
			allInfo.CreateControllerModel(controller_target)
		default:
			zap.L().Error("暂不支持生成的语言源文件。", zap.String("language", language))
		}
		// TODO 生成完毕后 用代码 对文件再执行一次 go fmt
		format_po := "go fmt " + po_target
		format_controller := "go fmt " + controller_target
		commond_format_po := kexec.CommandString(format_po)
		// commond_format_controller := kexec.CommandString(format_controller)
		commond_format_po.Run()
		// commond_format_controller.Run()
		zap.L().Debug("格式化执行命令完毕", zap.String("format_po", format_po),zap.String("format_controller", format_controller))
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
func getSourceFileType() consts.SourceFormat {
	// switch sourceFileFormat {
	// case "json":
	// 	SourceFormat = consts.Json
	// case "yaml", "yml":
	// 	SourceFormat = consts.Yaml
	// case "":
	var SourceFormat consts.SourceFormat
	ext := path.Ext(po_source)
	if strings.EqualFold(".yaml", ext) || strings.EqualFold(".yml", ext) {
		SourceFormat = consts.Yaml
	} else if strings.EqualFold(".json", ext) {
		SourceFormat = consts.Json
	} else {
		zap.L().Error("无法判断源文件类型")
		os.Exit(1)
	}
	// }
	zap.L().Info("SourceFormat", zap.Reflect("SourceFormat", SourceFormat))
	return SourceFormat
}

func init() {
	// 第4个参数为默认值
	rootCmd.AddCommand(modelCmd)
	modelCmd.Flags().StringVarP(&po_target, "po-target", "t", "", "请输入po类存储文件")
	modelCmd.MarkFlagRequired("po-target") // 必填
	modelCmd.Flags().StringVarP(&po_source, "po-source", "s", "", "请输入po定义源文件")
	modelCmd.MarkFlagRequired("po-source") // 必填
	modelCmd.Flags().BoolVarP(&po_same_name_as_table, "po-same-name-as-table", "n", false, "PO是否同名于表名字段名")
	// modelCmd.Flags().StringVarP(&sourceFileFormat, "sourceFileFormat", "f", "", "配置模型的文件类型,无值时根据文件后缀判断，例如json，yaml，yml")
	modelCmd.Flags().StringVarP(&Orm, "orm", "o", "xorm", "数据库持久化框架，默认xorm,例如 xorm,gorm,mybatis")
	modelCmd.Flags().StringVarP(&JsonCase, "jsoncase", "c", "origin", "生成的po中的json首字母使用大写或小写，默认使用origin(与字段相同),例如 origin,lower,upper")
	modelCmd.Flags().StringVarP(&controller_target, "controller-target", "", "", "请输入controller类存储文件")
	// modelCmd.MarkFlagRequired("controller-target") // 必填
	// modelCmd.Flags().StringVarP(&controller_source, "controller-source", "", "", "请输入controller定义源文件")
	// modelCmd.MarkFlagRequired("controller-source") // 必填

}
