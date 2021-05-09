package test

import (
	"bestsell/common"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"math"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

type TableStruct struct{
	Name string
	Key string
	Ctrl string
	Desc string
}
var ExcelDatas map[string][]map[string]interface{}
var ExcelStructs map[string][]*TableStruct

var _excelPath  string
var _jsonPath  string


func init() {
	initColumKeys()
	_excelPath = common.ExcelPath
	_jsonPath = common.JsonPath

	ExcelDatas = map[string][]map[string]interface{}{}
	ExcelStructs = map[string][]*TableStruct{}
	//ConvertDatas()

	
	//jsonText, err := json.Marshal(ExcelDatas)
	//if err != nil {
	//	log.Fatalf("JSON marshaling failed: %s", err)
	//}
	//fmt.Printf("%s\n", jsonText)
}

func ConvertDatas()  {
	dirPath := _excelPath
	//files, _ := ioutil.ReadDir(dirPath)
	//inPrefix := ".xlsx"
	//exPrefix := "~$"
	//var fileNames []string
	//for _, file := range files {
	//	filename := file.Name()
	//	if strings.HasSuffix(filename, inPrefix) && !strings.HasPrefix(filename, exPrefix) {
	//		fmt.Println("filename:"+filename)
	//		fileNames = append(fileNames, filename)
	//	}
	//}

	// fileNames := []string {
	// 	"address.xlsx",
	// 	"uibanner.xlsx",
	// }
	// loadExcelOneSheet(&fileNames, &dirPath)

	fileNames := []string {
		"goods.xlsx",
	}
	for _, fileName := range fileNames {
		fullPath := path.Join(dirPath, fileName)
		loadExcelMutilSheet(fileName, fullPath)
	}
	//backupPath := "../assert/data/backup/"
	//saveExcel(&fileNames, &backupPath)
	fmt.Println("ConvertDatas finish====>>")
}

func loadExcelMutilSheet(fileName string, fullPath string) {
	xlsx, err := excelize.OpenFile(fullPath)
	if err != nil {
		panic(err)
	}
	fileName = fileName[:len(fileName)-5]
	xlsName := "cfg_"+fileName
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
		keyRow := rows[1]
		ctrlRow := rows[2]
		// nameRow := rows[0]
		// descRow := rows[3]

		// var structs = []*TableStruct{}
		// for c :=0; c<len(ctrlRow); c++ {
		// 	ctrlCell := ctrlRow[c]
		// 	if len(ctrlCell) > 0 {
		// 		key := keyRow[c]
		// 		if len(key) == 0 {
		// 			continue
		// 		}
		// 		strc := &TableStruct{
		// 			Name:nameRow[c],
		// 			Key:keyRow[c],
		// 			Ctrl:ctrlRow[c],
		// 			Desc:descRow[c],
		// 		}
		// 		structs = append(structs, strc)
		// 	}
		// }
		// ExcelStructs[xlsName] = structs
		for i := 4;i<len(rows);i++ {
			row := rows[i]
			if len(row[0]) <= 0 {
				continue
			}
			data := make(map[string]interface{})
			for c :=0; c<len(ctrlRow); c++ {
				ctrlCell := ctrlRow[c]
				if len(ctrlCell) > 0 {
					key := keyRow[c]
					if len(key) == 0 {
						continue
					}
					cell := row[c]
					if c == 0 && len(cell) == 0{
						data = nil
						break
					}
					if ctrlCell[0] == 'i' {
						val, ok := strconv.Atoi(cell)
						if ok != nil {
							val = 0
						}
						data[key] = val
					}else if ctrlCell[0] == 'f' {
						val, ok := strconv.ParseFloat(cell, 64)
						if ok != nil {
							val = 0.0
						}
						data[key] = val
					}else{
						data[key] = cell
					}
				}
			}
			if data != nil {
				datas = append(datas, data)
			}
		}
	}
	ExcelDatas[xlsName] = datas

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

func loadExcelOneSheet(fileNames *[]string, dirPath *string) {
	for _,fileName := range(*fileNames) {
		fullPath := path.Join(*dirPath,fileName)
		xlsx, err := excelize.OpenFile(fullPath)
		if err != nil {
			panic(err)
			return
		}
		fileName = fileName[:len(fileName)-5]
		sheetMap := xlsx.GetSheetMap()
		for _,sheet := range sheetMap {
			if len(sheet) == 0 {
				continue
			}
			if strings.HasPrefix(sheet, "Sheet") {
				continue
			}
			xlsName := "cfg_"+fileName
			rows := xlsx.GetRows(sheet)
			if len(rows) == 0 {
				continue
			}
			keyRow := rows[1]
			ctrlRow := rows[2]
			nameRow := rows[0]
			descRow := rows[3]

			var structs = []*TableStruct{}
			for c :=0; c<len(ctrlRow); c++ {
				ctrlCell := ctrlRow[c]
				if len(ctrlCell) > 0 {
					key := keyRow[c]
					if len(key) == 0 {
						continue
					}
					strc := &TableStruct{
						Name:nameRow[c],
						Key:keyRow[c],
						Ctrl:ctrlRow[c],
						Desc:descRow[c],
					}
					structs = append(structs, strc)
				}
			}
			ExcelStructs[xlsName] = structs

			var datas = []map[string]interface{}{}
			for i := 4;i<len(rows);i++ {
				row := rows[i]
				if len(row[0]) <= 0 {
					continue
				}
				data := make(map[string]interface{})
				for c :=0; c<len(ctrlRow); c++ {
					ctrlCell := ctrlRow[c]
					if len(ctrlCell) > 0 {
						key := keyRow[c]
						if len(key) == 0 {
							continue
						}
						cell := row[c]
						if c == 0 && len(cell) == 0{
							data = nil
							break
						}
						if ctrlCell[0] == 'i' {
							val, ok := strconv.Atoi(cell)
							if ok != nil {
								val = 0
							}
							data[key] = val
						}else if ctrlCell[0] == 'f' {
							val, ok := strconv.ParseFloat(cell, 64)
							if ok != nil {
								val = 0.0
							}
							data[key] = val
						}else{
							data[key] = cell
						}
					}
				}
				if data != nil {
					datas = append(datas, data)
				}
			}
			ExcelDatas[xlsName] = datas

			str, err := json.Marshal(datas)
			if err != nil {
				fmt.Println(err)
			}
			jsonPath := path.Join(_jsonPath, xlsName+".json")

			if common.Exists(jsonPath) {
				os.Remove(jsonPath)
			}
			writeFile(jsonPath, []byte(str))
			break
		}
	}
}


var columKeys []string
func initColumKeys()  {
	size := int('Z'-'A')+1
	for i := 0; i < 100; i++ {
		idx := i%size
		count := math.Floor(float64(i/size))
		buffer := bytes.NewBuffer(nil)
		writer := bufio.NewWriter(buffer)
		if count == 0 {
			b := uint8(idx)+'A'
			writer.WriteByte(b)
		}else{
			b := uint8(count-1)+'A'
			writer.WriteByte(b)
			b = uint8(idx)+'A'
			writer.WriteByte(b)
		}
		_ = writer.Flush()
		key := string(buffer.Bytes())
		columKeys = append(columKeys, key)
	}
}

func SaveExcel(exceFile string) {
	fileNames := []string{exceFile+".xlsx"}
	dirPath := _excelPath
	saveExcel(&fileNames, &dirPath)
}

func LoadExcel(cfgName string)*[]map[string]interface{} {
	cfgs, ok := ExcelDatas[cfgName]
	if ok {
		return &cfgs
	}
	items := strings.Split(cfgName,"_")
	exceFile := items[1]
	fileNames := []string{exceFile+".xlsx"}
	dirPath := _excelPath
	loadExcelOneSheet(&fileNames, &dirPath)
	cfgs, ok = ExcelDatas[cfgName]
	if ok {
		return &cfgs
	}
	return nil
}

func saveExcel(fileNames *[]string, backupPath *string)  {
	//ExcelDatas = map[string][]map[string]interface{}{}
	for _,fileName := range(*fileNames) {
		if strings.HasSuffix(fileName, ".xlsx") {
			fileName = fileName[0:len(fileName)-len(".xlsx")]
		}
		eDatas := map[string]*[]map[string]interface{}{}
		for key,val := range ExcelDatas {
			if strings.Contains(key, fileName) {
				eDatas[key] = &val
			}
		}
		fmt.Println("fileName = "+fileName)
		var sheetName string
		if len(eDatas) == 0 {
			continue
		}
		xlsx := excelize.NewFile()

		for key,eData := range eDatas {
			tmps := strings.Split(key, "_")
			if len(tmps) == 2 {
				sheetName = tmps[1]
			}else if len(tmps) == 3{
				sheetName = tmps[2]
			}else{
				fmt.Println("sheetname format error ="+key)
				continue
			}
			index := xlsx.NewSheet(sheetName)
			xlsx.SetActiveSheet(index)

			excelStruct,ok := ExcelStructs[key]
			if !ok || len(excelStruct) == 0 {
				excelStruct = []*TableStruct{}
				for _,data := range *eData {
					for k,v := range data {
						stru := &TableStruct{}
						stru.Key = k
						switch v.(type) {
							case int:
								stru.Ctrl = "i"
						case float64:
								stru.Ctrl = "f"
						default:
								stru.Ctrl = "s"
						}
						excelStruct = append(excelStruct, stru)
					}
					break
				}
			}

			for j := 0; j < len(excelStruct); j++ {
				stru := (excelStruct)[j]
				axis := columKeys[j] + strconv.Itoa(1)
				xlsx.SetCellValue(sheetName, axis, stru.Name)

				axis = columKeys[j] + strconv.Itoa(2)
				xlsx.SetCellValue(sheetName, axis, stru.Key)

				axis = columKeys[j] + strconv.Itoa(3)
				xlsx.SetCellValue(sheetName, axis, stru.Ctrl)

				axis = columKeys[j] + strconv.Itoa(4)
				xlsx.SetCellValue(sheetName, axis, stru.Desc)
			}

			for i,data := range *eData {
				for j := 0; j < len(excelStruct); j++ {
					stru := (excelStruct)[j]
					v := data[stru.Key]
					axis := columKeys[j] + strconv.Itoa(i+5)
					switch v.(type) {
					case string:
						xlsx.SetCellValue(sheetName, axis, v)
					case int:
						xlsx.SetCellValue(sheetName, axis, strconv.Itoa(v.(int)))
					case float64:
						xlsx.SetCellValue(sheetName, axis, strconv.FormatFloat(v.(float64), 'f', -1, 64))
					default:
						panic("no support type:" + reflect.TypeOf(v).String())
					}
				}
			}
		}
		xlsxPath := *backupPath+fileName+".xlsx"
		if err := xlsx.SaveAs(xlsxPath); err != nil {
			fmt.Println(err)
		}
	}
}

func GetCfgList(cfgName string)[]map[string]interface{} {
	return ExcelDatas[cfgName]
}

func GetCfgById(cfgName string, id int)*map[string]interface{} {
	cfgs, ok := ExcelDatas[cfgName]
	if ok && len(cfgs) > 0 {
		val := cfgs[0]["id"]
		switch val.(type) {
			case int:
			default:
				return nil
		}
		for _,cfg := range(cfgs) {
			val = cfg["id"]
			if val == id {
				return &cfg
			}
		}
	}
	return nil
}

func GetCfgByKey(cfgName string, _key string)*map[string]interface{} {
	return GetCfgWith(cfgName, "key", _key)
}

func GetCfgFieldI(cfgName string, pkey string, pval interface{}, fkey string) (int,bool){
	cfg := GetCfgWith(cfgName,pkey, pval)
	if cfg == nil {
		return 0,false
	}
	field := (*cfg)[fkey]
	if field == nil {
		return 0,false
	}
	return field.(int),true
}

func GetCfgFieldS(cfgName string, pkey string, pval interface{}, fkey string) (string,bool){
	cfg := GetCfgWith(cfgName, pkey, pval)
	if cfg == nil {
		return "",false
	}
	field := (*cfg)[fkey]
	if field == nil {
		return "",false
	}
	return field.(string),true
}

func GetCfgWith(cfgName string, _key string, _val interface{})*map[string]interface{} {
	cfgs, ok := ExcelDatas[cfgName]
	if ok && len(cfgs) > 0 {
		tmp := cfgs[0]
		val,ok1 := tmp[_key]
		if !ok1 {
			return nil
		}
		switch val.(type) {
			case int:
				switch _val.(type) {
					case int:
						for _,cfg := range(cfgs) {
							val = cfg[_key]
							if val == _val {
								return &cfg
							}
						}
				}
			case float64:
				switch _val.(type) {
				case float64:
					for _,cfg := range(cfgs) {
						val = cfg[_key]
						if val == _val {
							return &cfg
						}
					}
				}
			case string:
				switch _val.(type) {
					case string:
						for _,cfg := range(cfgs) {
							val = cfg[_key]
							if strings.Compare(val.(string), _val.(string)) == 0 {
								return &cfg
							}
						}
				}
		}
	}
	return nil
}

func writeFile(filePath string, data []byte)  {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	n, err := file.Write(data)
	if err != nil && n != len(data) {
		panic(err)
	}
}

func LoadExcel2Map(fullPath string)*map[string][]*map[string]interface{} {
	retDatas := map[string][]*map[string]interface{}{}
	xlsx, err := excelize.OpenFile(fullPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
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
		keyRow := rows[1]
		ctrlRow := rows[2]

		var datas = []*map[string]interface{}{}
		for i := 4;i<len(rows);i++ {
			row := rows[i]
			data := map[string]interface{}{}
			for c :=0; c<len(ctrlRow); c++ {
				ctrlCell := ctrlRow[c]
				if len(ctrlCell) > 0 {
					key := keyRow[c]
					if len(key) == 0 {
						continue
					}
					cell := row[c]
					if c == 0 && len(cell) == 0{
						data = nil
						break
					}
					if ctrlCell[0] == 'i' {
						val, ok := strconv.Atoi(cell)
						if ok != nil {
							val = 0
						}
						data[key] = val
					}else if ctrlCell[0] == 'f' {
						val, ok := strconv.ParseFloat(cell, 64)
						if ok != nil {
							val = 0.0
						}
						data[key] = val
					}else{
						data[key] = cell
					}
				}
			}
			if data != nil {
				datas = append(datas, &data)
			}
		}
		retDatas[sheet] = datas
		break
	}
	return &retDatas
}

//func init() {
//	testFilePath := "/Users/linyou/svn/plan1/wechat/server/asset/gm/goods.xlsx"
//	sheets := LoadExcel2Map(testFilePath)
//	fmt.Println("sheets =", sheets)
//}
//func testExcel1(){
//	cfgs := GetCfgList("cfg_attr")
//
//	fmt.Println(cfgs)
//}
//
//func testExcel2(){
//	cfg := GetCfg("cfg_attr", "def_w")
//	fmt.Println("name =", cfg["name"])
//	fmt.Println("type =", cfg["type"])
//	fmt.Println("magic_type =", cfg["magic_type"])
//}

func Test(){
	GetCfgList("cfg_global")
	{
		cfg :=GetCfgByKey("cfg_global", "hero_max")
		fmt.Println("cfg =", cfg)
	}
	{
		cfg :=GetCfgById("cfg_global_power", 12)
		fmt.Println("cfg =", cfg)
	}
	{
		cfg :=GetCfgWith("cfg_global_holiday1", "date", "2020-1-25")
		fmt.Println("cfg =", cfg)
	}
}