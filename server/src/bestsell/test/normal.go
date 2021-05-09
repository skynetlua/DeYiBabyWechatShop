package test

import (
	"bestsell/common"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"os"
	"path"
	"strings"
)



func init() {
	_excelPath = common.ExcelPath
	_jsonPath = common.JsonPath
}

func ConvertExcel2JsonArrays() {
	dirPath := _excelPath
	fileNames := []string {
		"goodsdatas.xlsx",
	}
	for _, fileName := range fileNames {
		fullPath := path.Join(dirPath, fileName)
		convertExcel2JsonArray(fileName, fullPath)
	}
}

func convertExcel2JsonArray(fileName string, fullPath string) {
	xlsx, err := excelize.OpenFile(fullPath)
	if err != nil {
		panic(err)
	}
	fileName = fileName[:len(fileName)-5]
	xlsName := "cfgn_"+fileName
	var datas = []map[string]interface{}{}
	fileName = fileName[:len(fileName)-5]
	sheetMap := xlsx.GetSheetMap()
	for _,sheet := range sheetMap {
		if len(sheet) == 0 {
			continue
		}
		if strings.HasPrefix(sheet, "Sheet") {
			continue
		}
		rows := xlsx.GetRows(sheet)
		if len(rows) == 0 {
			continue
		}
		for i := 1;i<len(rows);i++ {
			row := rows[i]
			name := row[0]
			if len(name) <= 0 {
				continue
			}
			skuName := row[1]
			if len(skuName) > 0 {
				name = name+"("+skuName+")"
			}
			skuId := row[2]
			barCode := row[3]
			unit := row[5]
			mainType := row[12]
			subType := row[13]
			enterPrice := row[14]
			sellPrice := row[15]
			data := map[string]interface{} {
				"name": name,
				"skuId": skuId,
				"barCode": barCode,
				"unit": unit,
				"mainType": mainType,
				"subType": subType,
				"enterPrice": enterPrice,
				"sellPrice": sellPrice,
			}
			datas = append(datas, data)
		}
	}
	str, err := json.Marshal(datas)
	if err != nil {
		fmt.Println(err)
	}
	jsonPath := path.Join(_jsonPath, xlsName+".json")
	if common.Exists(jsonPath) {
		os.Remove(jsonPath)
	}
	writeFile(jsonPath, []byte(str))
}
