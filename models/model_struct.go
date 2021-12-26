package models

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
	ColumnName string `json:"column_name" yaml:"column_name"`
	//列类型
	ColumnType string `json:"column_type" yaml:"column_type"` // 如果没传入值则在go代码中（非模板中）根据PropType推导出来
	// 设置默认生成时间created 或 updated， xorm:"CreateTime timestamp created"` " xorm:"ModifyTime timestamp updated"`
	DefaultTime string `json:"defalut_time" yaml:"defalut_time"`
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
	// vo 属性数组
	VoPropSli []VoProp `json:"vo_prop_sli" yaml:"vo_prop_sli"`
}

type VoProp struct{
	// vo_prop_sli:
	// ## vo 对应的属性
	// ## vo json 显示名称
	// - vo_show: ReportGroupingName
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

	TableList  []Table  `json:"table_list" yaml:"table_list"`
	ImportList []string `json:"-"`

	PackageController    string `json:"package_controller" yaml:"package_controller"`
	ControllerImportList []string
	ModelPath            string `yaml:"model_path"`
}
