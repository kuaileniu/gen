package consts

type SourceFormat string

const (
	Json SourceFormat = "json"
	Yaml SourceFormat = "yaml"
)

type OrmType string

const (
	Xorm    OrmType = "xorm"
	Gorm    OrmType = "gorm"
	MyBatis OrmType = "mybatis"
)
