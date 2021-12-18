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

	PackageController string   `json:"package_controller" yaml:"package_controller"`
	ControllerImportList [] string 
	ModelPath string `yaml:"model_path"`
}
