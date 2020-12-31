package main

import(
	"github.com/kuaileniu/gen/cmd"
	"os"
	"github.com/kuaileniu/zlog"
	"go.uber.org/zap"
)

func init() {
	zlog.InitLogger(zlog.LogConfig{
		Filename:   "./logs/gen.log",
		Level:      "debug",
		MaxSize:    5,
		MaxBackups: 10,
		MaxAge:     10,
		Console:    true,
	})
}

func main(){
	err:=cmd.Execute()
	if err !=nil{
		zap.L().Error("cmd.Execute err", zap.Error(err))
	}
	zap.L().Info("gen over.")
	os.Exit(1)
}