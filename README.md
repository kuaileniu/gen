# go1.16
go build 和 go test 默认情况下不再修改 go.mod 和 go.sum。可通过 go mod tidy ，go get 或者手动完成；

go mod tidy

# gen
生成Go代码，Java代码，可配置xorm，gorm，dao，service,可使用yaml，json

# yaml 和 json 文件在线转换
- https://www.bejson.com/json/json2yaml

# 安装
go install -ldflags "-s -w" github.com/kuaileniu/gen // install 后的可执行文件名称为 gen; windows下可以，

当 git clone github.com/kuaileniu/gen 后，在gen根目录下执行
go install -ldflags "-s -w" github.com/kuaileniu/gen ，将安装clone下来的代码
go install -ldflags "-s -w" gen.go ，将安装clone下来的代码

# 使用

go run gen.go -l go model -t target/models/po_model.go -s testdata/model_info.json -n true -o xorm

go run gen.go -l go model -t target/models/po_model.go -s testdata/model_info.yml -n true -o xorm


gen -l go model -t target/models/po_model.go -s model_info.yml -n true -o xorm
