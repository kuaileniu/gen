package models

import (
	"bytes"
	"log"
	"strings"
	"text/template"

	"github.com/kuaileniu/gen/tplgo"
	"github.com/kuaileniu/sliceutil"
	"github.com/kuaileniu/zfile"
)

//数据库对应的类型，例如：varchar
// type ColumnType string

// go，java开发语言对应的属性类型，例如int，uint32
type PropType string

func (info *ModelInfo) CreatePoModel(pathFile string) error {
	tmpl, err := template.New("model").Parse(tplgo.PoModelTmpl)
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, info)

	err = zfile.ReWriteFile(pathFile, []byte(buf.Bytes()))
	return err
}

func (info *ModelInfo) CreateControllerModel(pathFile string) error {
	tmpl, err := template.New("controller").Parse(tplgo.ControllerModelTmpl)
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, info)

	err = zfile.ReWriteFile(pathFile, []byte(buf.Bytes()))
	return err
}

// 兼容go类型
func (info *ModelInfo) CompatibleGoType() {

	for table_index, table := range info.TableList {
		for col_index, col := range table.ColumnList {
			if strings.EqualFold("string", col.PropType) {
				col.PropType = "string"
			}
			table.ColumnList[col_index] = col
		}
		info.TableList[table_index] = table
	}
}

// 设置自动生成默认时间
func (info *ModelInfo) InferenceColumnDefaultTime() {
	for table_index, table := range info.TableList {
		for col_index, col := range table.ColumnList {
			defalutTime := strings.TrimSpace(col.DefaultTime)
			if defalutTime != "" {
				// created updated
				col.DefaultTime = " " + defalutTime
				table.ColumnList[col_index] = col
			}
		}
		info.TableList[table_index] = table
	}
}

// 推理出字段类型
func (info *ModelInfo) InferenceColumnType() {
	for table_index, table := range info.TableList {
		for col_index, col := range table.ColumnList {
			if col.ColumnType != "" {
				//这里不改变设置进来的ColumnType
				continue
			}
			if strings.EqualFold("string", col.PropType) && col.ColumnType == "" {
				col.ColumnType = "varchar"
			} else if strings.EqualFold("int64", col.PropType) {
				col.ColumnType = "bigint"
			} else if strings.EqualFold("int32", col.PropType) {
				col.ColumnType = "int"
			} else if strings.EqualFold("int16", col.PropType) {
				col.ColumnType = "smallint"
			} else if strings.EqualFold("int8", col.PropType) {
				col.ColumnType = "tinyint"
			} else if strings.HasPrefix("ztype.Money", col.PropType) {
				col.ColumnType = "decimal"
			} else if strings.EqualFold("ztype.Time", col.PropType) {
				col.ColumnType = "timestamp"
			} else if strings.EqualFold("bool", col.PropType) {
			} else if strings.EqualFold("float64", col.PropType) {
			} else if strings.EqualFold("zconst.KeYongStatus", col.PropType) {
				col.ColumnType = "varchar"
			} else if strings.EqualFold("zconst.TaskLockId", col.PropType) {
				col.ColumnType = "varchar"
			} else if strings.EqualFold("zconst.KeyType", col.PropType) {
				col.ColumnType = "varchar"
			} else if strings.EqualFold("zconst.EncriptType", col.PropType) {
				col.ColumnType = "varchar"
			} else if strings.EqualFold("zconst.CanShow", col.PropType) {
				col.ColumnType = "bool"
			} else if strings.EqualFold("zconst.Lang", col.PropType) {
				col.ColumnType = "varchar"
			} else {
				log.Printf("未将结构体属性类型[%v]映射到go数据库字段类型,%#v", col.PropType, col)
			}
			table.ColumnList[col_index] = col
		}
		info.TableList[table_index] = table
	}
}

// 推理出字段是否为主键
func (info *ModelInfo) InferencePropTypeIsKey() {
	for table_index, table := range info.TableList {
		for col_index, col := range table.ColumnList {
			if col.IsKey {
				col.PK = " pk" // pk前面的空格不可去掉
			}
			table.ColumnList[col_index] = col
		}
		info.TableList[table_index] = table
	}
}

// 推理出数据库字段类型的范围varchar(25) 或 decimal(10,2) 中的 （25） （10,2)
func (info *ModelInfo) InferenceColumnTypeRange() {
	for table_index, table := range info.TableList {
		for col_index, col := range table.ColumnList {
			if strings.EqualFold("varchar", col.ColumnType) && col.Length != "" {
				col.ColumnTypeRange = "(" + col.Length + ")"
			} else if strings.EqualFold("varchar", col.ColumnType) && col.Length == "" {
				col.ColumnTypeRange = "(256)"
			} else if strings.EqualFold("decimal", col.ColumnType) {
				if col.Precision != "" && col.Precision != "0" {
					col.ColumnTypeRange = "(" + col.Length + "," + col.Precision + ")"
				} else {
					col.ColumnTypeRange = "(" + col.Length + ")"
				}
			}
			table.ColumnList[col_index] = col
		}
		info.TableList[table_index] = table
	}
}

// 是否忽略为空
func (info *ModelInfo) InferenceOmitempty() {
	for table_index, table := range info.TableList {
		for col_index, col := range table.ColumnList {
			if strings.EqualFold("omitempty", col.Omitempty) {
				col.Omitempty = ",omitempty"
			}
			table.ColumnList[col_index] = col
		}
		info.TableList[table_index] = table
	}
}

// 不为空
func (info *ModelInfo) InferenceXormNotnull() {
	for table_index, table := range info.TableList {
		for col_index, col := range table.ColumnList {
			if strings.EqualFold("notnull", col.Notnull) {
				col.Notnull = " notnull"
			}
			table.ColumnList[col_index] = col
		}
		info.TableList[table_index] = table
	}
}

func (info *ModelInfo) InferenceJsonName(jsonCase string) {
	// JsonName
	for table_index, table := range info.TableList {
		for col_index, col := range table.ColumnList {
			// col.JsonName = "string"
			if "" != strings.TrimSpace(col.JsonName) {
				col.JsonName = strings.TrimSpace(col.JsonName)
			} else {
				JsonName := strings.TrimSpace(col.PropName)
				// origin,lower,upper
				if strings.EqualFold(jsonCase, "lower") {
					col.JsonName = StringFirstToLower(JsonName)
				} else if strings.EqualFold(jsonCase, "upper") {
					col.JsonName = StringFirstToUpper(JsonName)
				} else {
					col.JsonName = JsonName
				}
			}
			table.ColumnList[col_index] = col
		}
		info.TableList[table_index] = table
	}
}

// 推理出导出json时是否含此字段
func (info *ModelInfo) InferenceUnique() {
	for table_index, table := range info.TableList {
		for col_index, col := range table.ColumnList {
			colUnique := strings.ToLower(col.Unique)
			if strings.HasPrefix(colUnique, "unique") {
				col.Unique = " " + col.Unique
			}
			table.ColumnList[col_index] = col
		}
		info.TableList[table_index] = table
	}
}

// 推理出default
func (info *ModelInfo) InferenceXormDefault() {
	for table_index, table := range info.TableList {
		for col_index, col := range table.ColumnList {
			if col.Default != "" {
				col.Default = " " + col.Default
			}
			table.ColumnList[col_index] = col
		}
		info.TableList[table_index] = table
	}
}

// set Table name with PO name
func (info *ModelInfo) SetTableName() {
	for table_index, table := range info.TableList {
		table.TableName = table.PoName
		info.TableList[table_index] = table
	}
}

// set column name with property name
func (info *ModelInfo) SetColumnName() {
	for table_index, table := range info.TableList {
		for col_index, col := range table.ColumnList {
			if strings.EqualFold("-", col.ColumnName) { // 指明不生成表字段
				continue
			}
			col.ColumnName = col.PropName
			table.ColumnList[col_index] = col
		}
		info.TableList[table_index] = table
	}
}

// 汇集 import项
func (info *ModelInfo) CollectImport() {
	for _, table := range info.TableList {
		for _, col := range table.ColumnList {
			lower := strings.ToLower(col.PropType)
			if strings.HasPrefix(lower, "ztype.") {
				sliceutil.AddNoRepeatStr(&info.ImportList, "github.com/kuaileniu/ztype")
			} else if strings.HasPrefix(lower, "zconst.") {
				sliceutil.AddNoRepeatStr(&info.ImportList, "github.com/kuaileniu/zconst")
			}
		}
	}
}

// 汇集 controller.import项
func (info *ModelInfo) CollectControllerImport() {
	// for _, table := range info.TableList {
	// 	for _, col := range table.ColumnList {
	// 		lower := strings.ToLower(col.PropType)
	// 		if strings.HasPrefix(lower, "ztype.") {
	// 			sliceutil.AddNoRepeatStr(&info.ControllerImportList, "github.com/kuaileniu/ztype")
	// 		} else if strings.HasPrefix(lower, "zconst.") {
	// 			sliceutil.AddNoRepeatStr(&info.ControllerImportList, "github.com/kuaileniu/zconst")
	// 		}
	// 	}
	// }
	sliceutil.AddNoRepeatStr(&info.ControllerImportList, "github.com/gin-gonic/gin")
	sliceutil.AddNoRepeatStr(&info.ControllerImportList, "github.com/kuaileniu/db")
	sliceutil.AddNoRepeatStr(&info.ControllerImportList, info.ModelPath)
	sliceutil.AddNoRepeatStr(&info.ControllerImportList, "github.com/kuaileniu/zweb/ctx")
	sliceutil.AddNoRepeatStr(&info.ControllerImportList, "github.com/kuaileniu/zweb/enum")
	sliceutil.AddNoRepeatStr(&info.ControllerImportList, "go.uber.org/zap")
	sliceutil.AddNoRepeatStr(&info.ControllerImportList, "net/http")
	sliceutil.AddNoRepeatStr(&info.ControllerImportList, "strings")
}

//字符串首字母转化为大写
//AbC -> abC
func StringFirstToUpper(str string) string {
	rstr := []rune(str)
	first := rstr[0]
	firstUpper := strings.ToUpper(string(first))
	return firstUpper + string(rstr[1:])
}

//字符串首字母转化为小写
//abC -> AbC
func StringFirstToLower(str string) string {
	rstr := []rune(str)
	first := rstr[0]
	firstUpper := strings.ToLower(string(first))
	return firstUpper + string(rstr[1:])
}
