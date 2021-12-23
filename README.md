# go1.18
go build 和 go test 默认情况下不再修改 go.mod 和 go.sum。可通过 go mod tidy ，go get 或者手动完成；

go mod tidy

- go 1.17 版本不支持 模版中有这个语句 {{- continue -}}  
# gen
生成Go代码，Java代码，可配置xorm，gorm，dao，service,可使用yaml，json

# yaml 和 json 文件在线转换
- https://www.bejson.com/json/json2yaml

# 安装
go install -ldflags "-s -w" github.com/kuaileniu/gen@v0.0.7 // install 后的可执行文件名称为 gen; windows下go1.15,go1.16时可以，

当 git clone github.com/kuaileniu/gen 后，在gen根目录下执行
go install -ldflags "-s -w" github.com/kuaileniu/gen ，将安装clone下来的代码
go install -ldflags "-s -w" gen.go ，将安装clone下来的代码

# 其它项目应用方式
//go:generate gen -l go model -t model/config_po.go -s yamlsrc/e_config.yml -n true -o xorm -c lower
// TODO 校验若是在docker容器内，退出系统
func main() {}

# 参数说明
-n true //生成的数据库字段同go字段名称； 去掉此参数设置则 生成的orm映射字段采用 a1,a2 这种非明确含义字段。

# 调试

go run gen.go -l go model -t target/models/po_model.go -s testdata/model_info.json -n true -o xorm

go run gen.go -l go model -t target/models/po_model.go -s testdata/model_info.yml -n true -o xorm -c lower


//gen -l go model -t target/models/po_model.go -s model_info.yml -n true -o xorm


go run gen.go -l go model -t target/model/po_model.go -s testdata/model_info.json -o xorm
go run gen.go --lang go model --po-target target/model/po_model.go --po-source testdata/model_info.json --orm xorm --controller-target target/controllers/controller_model.go
 


//go:generate gen --lang go model --po-target model/jtbl_po.go --po-source yamlsrc/a_jtbl.yml --orm xorm --controller-target target/controllers/controller_model.go
 