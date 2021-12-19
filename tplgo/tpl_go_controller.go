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
			zap.L().Error("根据 {{.PropName}} 查询 {{$table.PoName}} 时异常", zap.Error(e))
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
	{{- range .ColumnList}}        {{/*单独设置 Id CreatedBy */}}
		{{- if eq .PropName "Id" }}
			po.Id = id.CreateTimeId15()
		{{- end }}
		{{- if eq .PropName "CreatedBy" }}
			po.CreatedBy = GetCurrentStaffName(c)
		{{- end }}
		{{- if eq .PropName "CanDel" }}
			po.CanDel = true
		{{- end }}
	{{- end }}
	if _, err := db.Engine.Insert(&po); err != nil {
		zap.L().Error("添加 {{$table.PoName}} 时异常", zap.Error(err))
		c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusErrorTip, Msg: "添加失败", EnglishMsg: "Add failed"})
		return
	}
	{{ if $table.MarkCannotDel -}}
	{{$table.MarkCannotDel}}
	{{ end -}}
	c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusOkTip, Msg: "添加成功", EnglishMsg: "Add success", Data: gin.H{"Id": po.Id}})
}

func Del{{.PoName}}(c *gin.Context) {
	req := struct {
		Ids []int64  ` + "`" + `json:"Ids"` + "`" + `
	}{}
	c.ShouldBindJSON(&req)

	if len(req.Ids) < 1 {
		c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusErrorTip, Msg: "Ids不可为空", EnglishMsg: "Ids can't be blank"})
		return
	}

	if _, err := db.Engine.NewSession().In(model.{{$table.PoName}}_Id_DB, req.Ids).And(model.{{$table.PoName}}_CanDel_DB + "=true").Delete(model.{{$table.PoName}}Ptr); err != nil {
		zap.L().Error("删除 {{$table.PoName}} 时异常", zap.Error(err))
		c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusErrorTip, Msg: "删除失败。", EnglishMsg: "Failed."})
		return
	}
	c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusOkTip, Msg: "删除成功。", EnglishMsg: "Success."})
}

func Edit{{.PoName}}(c *gin.Context) {
	req := struct {
		model.{{.PoName}}
	}{}
	body := GetBody(c, &req)

	if req.Id < 1 {
		c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusErrorTip, Msg: "Id不可为空", EnglishMsg: "Id can't be empty"})
		return
	}

	po := model.{{.PoName}}{Id: req.Id}
	session := db.Engine.NewSession()

	{{- range .ColumnList}}        {{/* 判断前端是否传值过来*/}}

		{{ if or .IsKey (eq .PropName "CreatedBy") (eq .PropName "ModifiedBy") (eq .PropName "CreateTime") (eq .PropName  "ModifyTime") (eq .PropName "CanDel") }}
			{{continue}}
		{{ end -}}
		{{- if and .ForeignKey  (eq .PropType "int64") }}
		if strings.Contains(body, strings.ToUpper(model.{{$table.PoName}}_{{.PropName}}_GO)) && req.{{.PropName}} >0 {
		{{ else }}
		if strings.Contains(body, strings.ToUpper(model.{{$table.PoName}}_{{.PropName}}_GO))  {
		{{ end -}}
			session.Cols(model.{{$table.PoName}}_{{.PropName}}_DB)
			po.{{.PropName}} = req.{{.PropName}}
			{{- if eq .AppNotnull "notnull" }}     {{/*判断必传参数必须存在*/}}
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
			{{- if eq .AppNotRepeat "notrepeat" }} {{/*判断数据不重复*/}}
			{
				{{ if  $table.ZoneKey -}}
				exist, e := db.Engine.In(model.{{$table.PoName}}_{{$table.ZoneKey}}_DB, req.JobId).Exist(&model.{{$table.PoName}}{ {{.PropName}}: req.{{.PropName}} })
				{{ else }}
				exist, e := db.Engine.Exist(&model.{{$table.PoName}}{ {{.PropName}}: req.{{.PropName}} })
				{{ end -}}
				if e != nil {
					zap.L().Error("根据 {{.PropName}} 查询 {{$table.PoName}} 时异常", zap.Error(e))
					c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusErrorTip, Msg: "添加失败。", EnglishMsg: "Add failed"})
					return
				}
				if exist {
					c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusErrorTip, Msg: "添加失败,{{.PropComment}}重复。", EnglishMsg: "Add failed,duplicate {{.PropName}}."})
					return
				}
			}{{ end -}}{{/*判断数据不重复*/}}
		}
	 
	{{ end -}}
	{{- range .ColumnList}}        {{/*单独设置 ModifiedBy */}}
		{{- if eq .PropName "ModifiedBy" }}
			po.ModifiedBy = GetCurrentStaffName(c)
			session.Cols(model.{{$table.PoName}}_{{.PropName}}_DB)
		{{- end }}
	{{- end }}
	_, err := session.ID(req.Id).Update(&po){{/* 保存修改的内容 */}}
	if err != nil {
		zap.L().Error("修改 {{$table.PoName}} 失败", zap.Error(err))
		c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusErrorTip, Msg: "修改失败", EnglishMsg: "Failed."})
		return
	}
	{{- if .MarkCannotDel}}
	{{.MarkCannotDel}}
	{{- end }}
	c.JSON(http.StatusOK, &ctx.Resp{Status: enum.StatusOkTip, Msg: "修改成功。", EnglishMsg: "Success."})
}

func Get{{.PoName}}Page(c *gin.Context) {
	req := struct {
		model.{{.PoName}}
		ctx.Req
		Search string ` + "`" + `json:"search"` + "`" + `
	}{}
	c.ShouldBindJSON(&req)
	 
	{{ if $table.ZoneKey }}
	if req.{{$table.ZoneKey}} < 1 {
		c.JSON(http.StatusOK, ctx.Resp{Status: enum.StatusErrorTip, Msg: "请填写正确的项目ID。", EnglishMsg: "{{$table.ZoneKey}} is incorrect."})
		return
	}
	{{ end }}
	{{- range .SearchSli}}
		//{{.}}
	{{- end }}
}
{{- end}}
`
