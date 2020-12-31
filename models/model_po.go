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

type Column struct {
	// 是主键
	IsKey bool   `json:"is_key" yaml:"is_key"`
	PK    string `json:"-" yaml:"-"` //在go代码中根据 IsKey 推导出 是"pk"或空
	// Omitempty string `json:"omit"`// json中omitempty是关键词，不能用于读取数据进来
	Omitempty string `json:"omitempty" yaml:"omitempty"` // json中omitempty是关键词，不能用于读取数据进来
	// 索引名字
	Unique string `json:"unique" yaml:"unique"`
	// go，java对应的属性名字
	PropName string `json:"prop_name" yaml:"prop_name"`
	// 配置出来，反推 columnType
	PropType string `json:"prop_type" yaml:"prop_type"`
	// 导出和导入json时的名字
	JsonName string `json:"json_name" yaml:"json_name"`
	//不为空
	Notnull string `json:"notnull" yaml:"notnull"` // "notnull": "notnull",
	// 数据库对应的字段名字
	ColumnName string `json:"column_name" yaml:"column_name"`
	//列类型
	ColumnType string `json:"column_type" yaml:"column_type"` // 如果没传入值则在go代码中（非模板中）根据PropType推导出来
	//长度
	Length string `json:"length" yaml:"length"`
	//精度
	Precision string `json:"precision" yaml:"precision"`
	//数据库字段类型的范围varchar(25) 或 decimal(10,2) 中的 （25） （10,2)
	ColumnTypeRange string // 在go中根据 ColumnType、Length、Precision 三个条件推断出来
	//默认值
	Default string `json:"default" yaml:"default"`
	//备注
	PropComment string `json:"prop_comment" yaml:"prop_comment"`
	// 在数据库中生成comment
	CommentShowInDB bool `json:"comment_in_db" yaml:"comment_in_db"`
}

type Table struct {
	TableName  string   `json:"table_name" yaml:"table_name"`
	PoName     string   `json:"po_name" yaml:"po_name"`
	ColumnList []Column `json:"column_list" yaml:"column_list"`
	//备注
	PoComment string `json:"po_comment" yaml:"po_comment"`
	// 在数据库中生成comment
	CommentInDB bool `json:"comment_in_db" yaml:"comment_in_db"`
}

type ModelInfo struct {
	PackageName string   `json:"package_name" yaml:"package_name"`
	TableList   []Table  `json:"table_list" yaml:"table_list"`
	ImportList  []string `json:"-"`
}

func (info *ModelInfo) CreatePoModel(pathFile string) error {
	tmpl, err := template.New("model").Parse(tplgo.PoModelTmpl)
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

func (info *ModelInfo) InferenceJsonName() {
	// JsonName
	for table_index, table := range info.TableList {
		for col_index, col := range table.ColumnList {
			// col.JsonName = "string"
			if "" != strings.TrimSpace(col.JsonName) {
				col.JsonName = strings.TrimSpace(col.JsonName)
			} else {
				col.JsonName = strings.TrimSpace(col.PropName)
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
