package tplgo

var PoModelTmpl = `package {{.PackageName}}
{{- if .ImportList}}

import (
	{{- range $v:= .ImportList}}
	"{{$v}}"
	{{- end}}
)
{{- end }}
{{range $table:= .TableList}}
// {{$table.PoComment}}
var {{$table.PoName}}Ptr = &{{$table.PoName}}{}

{{- end}}
{{range $table:= .TableList}}
func (po *{{.PoName}}) TableName() string {
	return "{{.TableName}}"
}

// {{$table.PoComment}}
type {{$table.PoName}} struct {
	{{- range .ColumnList}}
	{{- if .PropComment}}
	// {{ .PropComment}}
	{{- end}}
	{{.PropName}} {{.PropType}} ` + "`" + `json:"{{.JsonName}}{{.Omitempty}}" xorm:"{{.ColumnName}}{{.PK}} {{.ColumnType}}{{.ColumnTypeRange}}{{.DefaultTime}}{{.Notnull}}{{.Default}}{{.Unique}}"` + "`" + `
	{{- end}}
}

{{- if .ColumnList}}
const (
{{- range .ColumnList}}
	{{$table.PoName}}_{{.PropName}} = "{{.ColumnName}}"
{{- end}}
)
{{- end}}
{{end}}`
