# gen
生成Go代码，Java代码，可配置xorm，gorm，dao，service,可使用yaml，json


# 安装
go install -ldflags "-s -w" github.com/kuaileniu/gen // install 后的可执行文件名称为 gen


# 使用

 go run gen.go -l go model -t target/models/po_model.go -s testdata/model_info.json