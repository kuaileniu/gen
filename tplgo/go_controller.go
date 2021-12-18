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

func Add{{.PoName}}(c *gin.Context) {
	req := struct {
		model.{{.PoName}}
	}{}
	c.ShouldBindJSON(&req)
	{{- range .ColumnList}}{{/*不能为空的情况*/}}
	{{- if eq .AppNotnull "notnull" }}
	{{- if eq .PropType "string"}}
	if strings.TrimSpace(req.{{ .PropName}}) == "" {
		c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusErrorTip, Msg: "{{.PropComment}}不可为空。", EnglishMsg: "{{ .PropName}} can't be blank."})
		return
	}
	{{- end}}
	{{- if eq .PropType "int64"}}
	if req.{{ .PropName}} < 1 {
		c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusErrorTip, Msg: "{{.PropComment}}不可为空。", EnglishMsg: "{{ .PropName}} can't be blank."})
		return
	}
	{{- end}}
	{{- end}}
	{{- end}}

	{{- range .ColumnList}}{{/*不能为空的情况*/}}
	{{- if eq .AppNotRepeat "notrepeat" }}
	{
		{{ if  $table.ZoneKey -}}
		exist, e := db.Engine.In(model.{{$table.PoName}}_{{$table.ZoneKey}}_DB, req.JobId).Exist(&model.{{$table.PoName}}{ {{.PropName}}: req.{{.PropName}} })
		{{ else }}
		exist, e := db.Engine.Exist(&model.{{$table.PoName}}{ {{.PropName}}: req.{{.PropName}} })
		{{ end -}}
		if e != nil {
			zap.L().Error("根据 AreaName 查询 JtblArea 时异常", zap.Error(e))
			c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusErrorTip, Msg: "添加失败。", EnglishMsg: "Add failed"})
			return
		}
		if exist {
			c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusErrorTip, Msg: "添加失败,{{.PropComment}}重复。", EnglishMsg: "Add failed,duplicate {{.PropName}}."})
			return
		}
	}
	
	{{- end}}
	{{- end}}
}

{{- end}}
`
