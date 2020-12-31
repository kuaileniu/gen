package parser

import (
	"io/ioutil"
	"os"

	"github.com/kuaileniu/beanutil"
	"github.com/kuaileniu/gen/consts"
	"github.com/kuaileniu/gen/models"
	"go.uber.org/zap"
)

func GetAllInfo(pathFile string, format consts.SourceFormat) *models.ModelInfo {
	// jsonBytes, _ := ioutil.ReadFile("./all_info.json")
	switch format {
	case consts.Json:
		return fromJson(pathFile)
	}
	return nil
}

func fromJson(pathFile string) *models.ModelInfo {
	jsonBytes, _ := ioutil.ReadFile(pathFile)
	modelInfo := models.ModelInfo{}
	if err := beanutil.BytesToStruct(jsonBytes, &modelInfo); err != nil {
		zap.L().Error("解析json数据时异常", zap.Error(err))
		os.Exit(1)
	}
	return &modelInfo
}
