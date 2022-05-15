package models

import "strings"

type Column struct {
	// 是主键
	IsKey bool   `json:"is_key" yaml:"is_key"`
	PK    string `json:"-" yaml:"-"` //在go代码中根据 IsKey 推导出 是"pk"或空
	// Omitempty string `json:"omit"`// json中omitempty是关键词，不能用于读取数据进来
	Omitempty string `json:"omitempty" yaml:"omitempty"` // json中omitempty是关键词，不能用于读取数据进来
	// 唯一索引名字
	Unique string `json:"unique" yaml:"unique"` //   unique: unique;unique: unique(uniquename)
	// 索引名字
	Index string `json:"index" yaml:"index"` //   index: index; index: index(indexname)
	// go，java对应的属性名字
	PropName string `json:"prop_name" yaml:"prop_name"`
	// 是外键
	ForeignKey bool `json:"foreign_key" yaml:"foreign_key"`
	// 配置出来，反推 columnType
	PropType string `json:"prop_type" yaml:"prop_type"`
	// 导出和导入json时的名字
	JsonName string `json:"json_name" yaml:"json_name"`
	//数据库库中字段不可为空
	Notnull string `json:"notnull" yaml:"notnull"` // "notnull": "notnull",
	//代码中判断不可为空
	AppNotnull string `json:"app_notnull" yaml:"app_notnull"` // "app_notnull": notnull,
	//代码中判断不可重复
	AppNotRepeat string `json:"app_notrepeat" yaml:"app_notrepeat"` // "app_notrepeat": app_notrepeat,
	// 数据库对应的字段名字
	ColumnName string `json:"column_name" yaml:"column_name"` //在数据库中不需要建立字段时 column_name: '-'
	//列类型 column_type: 'DECIMAL(10,2)'
	ColumnType string `json:"column_type" yaml:"column_type"` // 如果没传入值则在go代码中（非模板中）根据PropType推导出来
	// 设置默认生成时间created 或 updated， xorm:"CreateTime timestamp created"` " xorm:"ModifyTime timestamp updated"`
	DefaultTime string `json:"defalut_time" yaml:"defalut_time"`
	//长度
	Length string `json:"length" yaml:"length"`
	//精度 可用 column_type: 'DECIMAL(10,2)' 替代
	Precision string `json:"precision" yaml:"precision"`
	// 可用 column_type: 'DECIMAL(10,2)' 替代
	//数据库字段类型的范围varchar(25) 或 decimal(10,2) 中的 （25） （10,2)
	ColumnTypeRange string // 在go中根据 ColumnType、Length、Precision 三个条件推断出来
	//默认值
	Default string `json:"default" yaml:"default"`
	//备注
	PropComment string `json:"prop_comment" yaml:"prop_comment"`
	// 完整的备注
	PropAllComment string `json:"prop_comment" yaml:"prop_all_comment"`
	// 在数据库中生成comment
	CommentShowInDB bool `json:"comment_in_db" yaml:"comment_in_db"`
}

type Table struct {
	TableName  string   `json:"table_name" yaml:"table_name"`
	PoName     string   `json:"po_name" yaml:"po_name"`
	ColumnList []Column `json:"column_list" yaml:"column_list"`
	//备注
	PoComment string `json:"po_comment" yaml:"po_comment"`
	// 区域主键，比例如JobId， 可以替代 JtblArea_JobId_DB中的JobId (db.Engine.In(model.JtblArea_JobId_DB, req.JobId).Exist(&model.JtblArea{AreaName: req.AreaName}) )
	ZoneKey string `json:"zone_key" yaml:"zone_key"`
	// # 范围Key的注释 ，例如 项目ID
	ZoneKeyComment string `json:"zone_key" yaml:"zone_key_comment"`
	// 在数据库中生成comment
	CommentInDB    bool `json:"comment_in_db" yaml:"comment_in_db"`
	CannotDelModel string

	//# 修改标记为不能删除的go语句
	//mark_cannot_del : 'MarkCannotDel(&CannotDelModel{JtblJobId: po.JobId})'
	MarkCannotDel string `json:"mark_cannot_del" yaml:"mark_cannot_del"`
	// 查询字段 []{AreaName,AreaDescription}
	SearchSli []string `json:"search_sli" yaml:"search_sli"`
	// Deprecated
	// vo 属性数组,适合 根据外键只取同一个关联表中的一列场景，不适合取同一个关联表中的多列场景
	VoPropSli []VoProp `json:"vo_prop_sli" yaml:"vo_prop_sli"`

	// multiple
	// vo 属性数组,适合 根据外键取同一个关联表中的多列场景
	VoMultiPropSli []VoMultiProp `yaml:"vo_multi_prop_sli"`
}

// 模版中调用go函数
// wg := pool.NewWaitGroup({{$table.MultiPropSize}})
// vo 中
func (table Table) MultiPropSize() int64 {
	return 19 // 此方法未使用，19数据为随意乱造
}

type VoMultiProp struct {
	// 	vo_multi_prop_sli:
	// 	# PO 名称
	//   - target_po: JtblReportGrouping
	// 	# po 的主键对应的名称
	// 	target_po_key: Id
	// 	multi_prop:
	// 	# vo json 显示名称
	// 	- vo_show: ReportGroupingName
	// 	  # json 属性的类型
	// 	  vo_show_type: string
	// 	  # vo 属性的注释
	// 	  vo_show_comment: 报告组名称
	// 	  # 查 vo_show 的数据源时所依赖vo的属性名称
	// 	  ref_prop_in_vo: ReportGroupingId
	// 	  # vo_show 的数据来源,po 中对应 vo_show 数据的字段
	// 	  the_po_prop: ReportGroupingName

	// PO 名称
	TargetPo    string `yaml:"target_po"`
	TargetPoKey string `yaml:"target_po_key"`
	// 查 vo_show 的数据源时所依赖vo的属性名称
	RefPropInVo string `yaml:"ref_prop_in_vo"`

	MultiPropSli []MultiProp `yaml:"multi_prop"`
}

func (p VoMultiProp) SelectColumn() string {
	// if err := db.Engine.In(model.{{.TargetPo}}_{{.TargetPoKey}}_DB, {{.RefPropInVo}}Sli).Cols(model.{{.TargetPo}}_{{.TargetPoKey}}_DB, model.{{.TargetPo}}_{{.PoDependent}}_DB).Find(&poMap); err != nil {
	columnSli := make([]string, 0)
	TargetPo := p.TargetPo

	columnSli = append(columnSli, "model."+TargetPo+"_"+p.TargetPoKey+"_DB")
	for _, prop := range p.MultiPropSli {
		columnSli = append(columnSli, "model."+TargetPo+"_"+prop.ThePoProp+"_DB")
	}
	return strings.Join(columnSli, ", ")
}

type MultiProp struct {
	// vo_multi_prop_sli:
	//     # PO 名称
	//   - target_po: JtblPressureGuage
	//     # po 的主键对应的名称
	//     target_po_key: Id
	//     multi_prop:
	//     # vo json 显示名称
	//     - vo_show: PressureGuageSerNo1
	//       # vo tag,为可选项，外侧是单引号，内部是左上角的撇符号
	//       vo_show_tag: '`json:"PressureGuageSerNo1,omitempty"`'
	//       # json 属性的类型
	//       vo_show_type: string
	//       # vo 属性的注释
	//       vo_show_comment: 压力表序列号1
	//       # 查 vo_show 的数据源时所依赖vo的属性名称
	//       ref_prop_in_vo: PressureGuageId1
	//       # vo_show 的数据来源,po 中对应 vo_show 数据的字段
	//       the_po_prop: SerialNumber
	//     - vo_show: CalibrationCertNo1
	//       # vo tag,为可选项，外侧是单引号，内部是左上角的撇符号
	//       vo_show_tag: '`json:"CalibrationCertNo1,omitempty"`'
	//       # json 属性的类型
	//       vo_show_type: string
	//       # vo 属性的注释
	//       vo_show_comment: 校验证书号1
	//       # 查 vo_show 的数据源时所依赖vo的属性名称
	//       ref_prop_in_vo: PressureGuageId1
	//       # vo_show 的数据来源,po 中对应 vo_show 数据的字段
	//       the_po_prop: CalibrationCertNo
	//   - target_po: JtblPressureGuage
	//     # po 的主键对应的名称
	//     target_po_key: Id
	//     multi_prop:
	//     # vo json 显示名称
	//     - vo_show: PressureGuageSerNo2
	//       # vo tag,为可选项，外侧是单引号，内部是左上角的撇符号
	//       vo_show_tag: '`json:"PressureGuageSerNo2,omitempty"`'
	//       # json 属性的类型
	//       vo_show_type: string
	//       # vo 属性的注释
	//       vo_show_comment: 压力表序列号2
	//       # 查 vo_show 的数据源时所依赖vo的属性名称
	//       ref_prop_in_vo: PressureGuageId2
	//       # vo_show 的数据来源,po 中对应 vo_show 数据的字段
	//       the_po_prop: SerialNumber
	//     - vo_show: CalibrationCertNo2
	//       # vo tag,为可选项，外侧是单引号，内部是左上角的撇符号
	//       vo_show_tag: '`json:"CalibrationCertNo2,omitempty"`'
	//       # json 属性的类型
	//       vo_show_type: string
	//       # vo 属性的注释
	//       vo_show_comment: 校验证书号2
	//       # 查 vo_show 的数据源时所依赖vo的属性名称
	//       ref_prop_in_vo: PressureGuageId2
	//       # vo_show 的数据来源,po 中对应 vo_show 数据的字段
	//       the_po_prop: CalibrationCertNo
	// vo 中显示给前端的属性名称
	VoShow string `yaml:"vo_show"`
	// vo属性的类型
	VoShowType string `yaml:"vo_show_type"`
	// vo tag,外侧是单引号，内部是左上角的撇符号
	VoShowTag string `yaml:"vo_show_tag"`
	// vo属性的注释
	VoShowComment string `yaml:"vo_show_comment"`
	// vo_show 的数据来源,po 中对应 vo_show 数据的字段
	ThePoProp string `yaml:"the_po_prop"`
}

type VoProp struct {
	// vo_prop_sli:
	// ## vo 对应的属性
	// ## vo json 显示名称
	// - vo_show: ReportGroupingName
	//   # json 属性的类型
	//   vo_show_type: string
	//   # 查 vo_prop_show 依赖的vo中的属性
	//   vo_dependent: ReportGroupingId
	//   # PO 名称
	//   po: JtblReportGrouping
	//   # po 的主键对应的名称
	//   po_key_name: Id
	//   # po 中对应 vo_show 的字段
	//   po_dependent: ReportGroupingName

	// vo 中显示给前端的属性名称
	VoShow string `yaml:"vo_show"`
	// vo属性的类型
	VoShowType string `yaml:"vo_show_type"`
	// vo属性的注释
	VoShowComment string `yaml:"vo_show_comment"`
	// 查 vo_prop_show 依赖的vo中的属性
	VoDependent string `yaml:"vo_dependent"`
	// PO 名称
	Po string `yaml:"po"`
	// po 的主键对应的名称
	PoKeyName string `yaml:"po_key_name"`
	// po 中对应 vo_show 的字段
	PoDependent string `yaml:"po_dependent"`
}

type ModelInfo struct {
	PackageName string `json:"package_name" yaml:"package_name"`

	TableList            []Table  `json:"table_list" yaml:"table_list"`
	ImportList           []string `json:"-"`
	Engine_name          string   `json:"engine_name" yaml:"engine_name"`
	PackageController    string   `json:"package_controller" yaml:"package_controller"`
	ControllerImportList []string
	ModelPath            string `yaml:"model_path"`
}
