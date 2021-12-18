package tplgo
 
var ControllerModelTmpl = `package {{.PackageController}}
{{- if .ControllerImportList}}
{{/* 注释是这么写 */}} 
import (
	{{- range $v:= .ControllerImportList}}
	"{{$v}}"
	{{- end}}
)
{{- end }}
{{range $table:= .TableList}}

func Add{{.PoName}}(c *gin.Context){
	req := struct {
		model.{{.PoName}}
	}{}
	c.ShouldBindJSON(&req)
	{{- range .ColumnList}}{{/*不能为空的情况*/}}
	{{- if eq .AppNotnull "notnull" }}
	if strings.TrimSpace(req.{{ .PropName}}) == "" {
		c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusErrorTip, Msg: "分区名称不可为空。", EnglishMsg: "AreaName can't be blank."})
		return
	}
	{{- end}}
	{{- end}}
}
{{- end}}
`

