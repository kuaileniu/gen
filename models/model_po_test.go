package models

import (
	"testing"

	"github.com/kuaileniu/zlog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	// flag.Set("alsologtostderr", "true")
	zlog.InitLogger(zlog.LogConfig{Filename: "../logs/gen.log", Console: true})
	zap.L().Info("执行开始\n")
	m.Run()
	zap.L().Info("执行完毕\n")
}

func TestStringFirstToUpper(t *testing.T) {
	s := StringFirstToUpper("abC")
	assert.Equal(t, "AbC", s)
	zap.L().Info("s", zap.String("s", s))
}

func TestStringFirstLower(t *testing.T) {
	s := StringFirstToLower("abC")
	assert.Equal(t, "abC", s)
	zap.L().Info("s", zap.String("s", s))

	s2 := StringFirstToLower("MbC")
	assert.Equal(t, "mbC", s2)
	zap.L().Info("s2", zap.String("s", s2))

	s3 := StringFirstToLower("中bC")
	assert.Equal(t, "中bC", s3)
	zap.L().Info("s3", zap.String("s", s3))
}
