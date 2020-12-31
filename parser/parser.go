package parser

import (
	"github.com/kuaileniu/gen/consts"
	"github.com/kuaileniu/gen/models"
	"github.com/kuaileniu/beanutil"
	"io/ioutil"
	"log"
)

func GetAllInfo(pathFile string,format consts.SourceFormat) *models.ModelInfo {
	// jsonBytes, _ := ioutil.ReadFile("./all_info.json")
	jsonBytes, _ := ioutil.ReadFile(pathFile)
	modelInfo := models.ModelInfo{}
	if err := beanutil.BytesToStruct(jsonBytes, &modelInfo); err != nil {
		log.Fatal("解析json数据时异常", err)
	}
	return &modelInfo
}
