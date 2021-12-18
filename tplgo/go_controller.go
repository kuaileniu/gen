package tplgo

var ControllerModelTmpl = `package {{.PackageController}}
{{- if .ControllerImportList}}
{{/* 注释是这么写 ,双花括号中的横线放最前面时，此段多出一个空行，放最后时无空行*/}} 
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

	po := req.{{$table.PoName}}
	po.Id = id.CreateTimeId15()
	po.CreatedBy = GetCurrentStaffName(c)
	po.CanDel = true
	if _, err := db.Engine.Insert(&po); err != nil {
		zap.L().Error("添加 {{$table.PoName}} 时异常", zap.Error(err))
		c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusErrorTip, Msg: "添加失败", EnglishMsg: "Add failed"})
		return
	}
	MarkCannotDel(&CannotDelModel{JtblJobId: po.JobId})
	c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusOkTip, Msg: "添加成功", EnglishMsg: "Add success", Data: gin.H{"Id": po.Id}})
}

{{- end}}
`
