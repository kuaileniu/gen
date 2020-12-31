package main

import (
	"github.com/kuaileniu/gen/cmd"
	"github.com/kuaileniu/zlog"
	"go.uber.org/zap"
)

func init() {
	zlog.InitLogger(zlog.LogConfig{
		Filename: "./logs/gen.log",
		Level:    "debug",
		// Level:      "info",
		MaxSize:    5,
		MaxBackups: 10,
		MaxAge:     10,
		Console:    true,
	})
}

func main() {
	err := cmd.Execute()
	if err != nil {
		zap.L().Error("cmd.Execute err", zap.Error(err))
	}
	zap.L().Info("gen over.")
	// os.Exit(0)  //添加这句会杀掉同名的进程，导致  同一个方法上面 不能执行多次//go:generate gen
}
