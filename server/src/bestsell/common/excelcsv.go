package common

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"io/ioutil"
	"math"
	"path"
	"strings"
	"github.com/axgle/mahonia"
)

type TableStruct struct{
	Name string
	Key string
	Ctrl string
	Desc string
}
var ExcelDatas map[string][]map[string]interface{}
var ExcelStructs map[string][]*TableStruct
var columKeys []string

func init() {
	// initColumKeys()

	// fmt.Println(columKeys)

	// ExcelDatas = map[string][]map[string]interface{}{}
	// ExcelStructs = map[string][]*TableStruct{}

	// initDatas()

	//jsonText, err := json.Marshal(ExcelDatas)
	//if err != nil {
	//	log.Fatalf("JSON marshaling failed: %s", err)
	//}
	//fmt.Printf("%s\n", jsonText)
}

func initDatas()  {
	dirPath := ExcelPath
	files, _ := ioutil.ReadDir(dirPath)
	inPrefix := ".csv"
	exPrefix := "~$"
	var fileNames []string
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, inPrefix) && !strings.HasPrefix(filename, exPrefix) {
			fileNames = append(fileNames, filename)
		}
	}
	loadExcels(&fileNames, &dirPath)
	//backupPath := "../assert/data/backup/"
	//saveExcel(&fileNames, &backupPath)
}

func LoadExcel(cfgName string)*[]map[string]interface{} {
	fmt.Println("exclcsv.LoadExcel cfgName=", cfgName)
	cfgs, ok := ExcelDatas[cfgName]
	if ok {
		return &cfgs
	}
	// items := strings.Split(cfgName, "_")
	name := cfgName[4:]
	exceFile := name+".csv"
	fmt.Println("exclcsv.LoadExcel exceFile=", exceFile)
	loadExcel(exceFile)
	cfgs, ok = ExcelDatas[cfgName]
	if ok {
		fmt.Println("exclcsv.LoadExcel Get obj cfgName=", cfgName)
		return &cfgs
	}
	fmt.Println("exclcsv.LoadExcel Get nil cfgName=", cfgName)
	return nil
}

func ConvertToString(src string, srcCode string, tagCode string) string {
    srcCoder := mahonia.NewDecoder(srcCode)
    srcResult := srcCoder.ConvertString(src)
    tagCoder := mahonia.NewDecoder(tagCode)
    _, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
    result := string(cdata)
    return result
}

func loadExcel(fileName string) {
	fullPath := path.Join(ExcelPath, fileName)
	fmt.Println("exclcsv.loadExcel:", fullPath)
	csvfile, err := os.Open(fullPath)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
		return
	}
	fileName = fileName[:len(fileName)-4]
	xlsName := "cfg_"+fileName
	fmt.Printf("exclcsv.loadExcel:xlsName %s\n", xlsName)

	defer csvfile.Close()
	reader := csv.NewReader(csvfile)
	var rows [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		rows = append(rows, record)
	}

	nameRow := rows[0]
	keyRow := rows[1]
	ctrlRow := rows[2]
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
			// fmt.Printf("strc.Name=%s ", strc.Name)
			// fmt.Printf("strc.Key=%s ", strc.Key)
			// fmt.Printf("strc.Ctrl=%s ", strc.Ctrl)
			// fmt.Printf("strc.Desc=%s ", strc.Desc)
		}
	}
	fmt.Printf("exclcsv.loadExcel:xlsName %s\n", xlsName)
	ExcelStructs[xlsName] = structs
	var datas = []map[string]interface{}{}
	for i := 4;i<len(rows);i++ {
		row := rows[i]
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

	fmt.Printf("exclcsv.loadExcel:xlsName xlsName=%s len=%d\n", xlsName, len(datas)) 
	ExcelDatas[xlsName] = datas

	str, err := json.Marshal(datas)
    if err != nil {
        fmt.Println(err)
    }
	jsonPath := path.Join(ExcelPath, xlsName+".json")
	writeFile(jsonPath, []byte(str))
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

func loadExcels(fileNames *[]string, dirPath *string) {
	for _,fileName := range(*fileNames) {
		loadExcel(fileName)
	}
}

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

func SaveExcels(exceFile string) {
	fileNames := []string{exceFile+".xlsx"}
	dirPath := ExcelPath
	saveExcel(&fileNames, &dirPath)
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
		//var sheetName string
		if len(eDatas) == 0 {
			continue
		}
		//xlsx := excelize.NewFile()
		//
		//for key,eData := range eDatas {
		//	tmps := strings.Split(key, "_")
		//	if len(tmps) == 2 {
		//		sheetName = tmps[1]
		//	}else if len(tmps) == 3{
		//		sheetName = tmps[2]
		//	}else{
		//		fmt.Println("sheetname format error ="+key)
		//		continue
		//	}
		//	index := xlsx.NewSheet(sheetName)
		//	xlsx.SetActiveSheet(index)
		//
		//	excelStruct,ok := ExcelStructs[key]
		//	if !ok || len(excelStruct) == 0 {
		//		excelStruct = []*TableStruct{}
		//		for _,data := range *eData {
		//			for k,v := range data {
		//				stru := &TableStruct{}
		//				stru.Key = k
		//				switch v.(type) {
		//					case int:
		//						stru.Ctrl = "i"
		//				case float64:
		//						stru.Ctrl = "f"
		//				default:
		//						stru.Ctrl = "s"
		//				}
		//				excelStruct = append(excelStruct, stru)
		//			}
		//			break
		//		}
		//	}
		//
		//	for j := 0; j < len(excelStruct); j++ {
		//		stru := (excelStruct)[j]
		//		axis := columKeys[j] + strconv.Itoa(1)
		//		xlsx.SetCellValue(sheetName, axis, stru.Name)
		//
		//		axis = columKeys[j] + strconv.Itoa(2)
		//		xlsx.SetCellValue(sheetName, axis, stru.Key)
		//
		//		axis = columKeys[j] + strconv.Itoa(3)
		//		xlsx.SetCellValue(sheetName, axis, stru.Ctrl)
		//
		//		axis = columKeys[j] + strconv.Itoa(4)
		//		xlsx.SetCellValue(sheetName, axis, stru.Desc)
		//	}
		//
		//	for i,data := range *eData {
		//		for j := 0; j < len(excelStruct); j++ {
		//			stru := (excelStruct)[j]
		//			v := data[stru.Key]
		//			axis := columKeys[j] + strconv.Itoa(i+5)
		//			switch v.(type) {
		//			case string:
		//				xlsx.SetCellValue(sheetName, axis, v)
		//			case int:
		//				xlsx.SetCellValue(sheetName, axis, strconv.Itoa(v.(int)))
		//			case float64:
		//				xlsx.SetCellValue(sheetName, axis, strconv.FormatFloat(v.(float64), 'f', -1, 64))
		//			default:
		//				panic("no support type:" + reflect.TypeOf(v).String())
		//			}
		//		}
		//	}
		//}
		//xlsxPath := *backupPath+fileName+".xlsx"
		//if err := xlsx.SaveAs(xlsxPath); err != nil {
		//	fmt.Println(err)
		//}
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


func LoadExcel2Map(fullPath string)*map[string][]*map[string]interface{} {
	retDatas := map[string][]*map[string]interface{}{}
	//xlsx, err := excelize.OpenFile(fullPath)
	//if err != nil {
	//	fmt.Println(err)
	//	return nil
	//}
	//sheetMap := xlsx.GetSheetMap()
	//for _,sheet := range sheetMap {
	//	if len(sheet) == 0 {
	//		continue
	//	}
	//	if strings.HasPrefix(sheet, "Sheet") {
	//		continue
	//	}
	//	rows := xlsx.GetRows(sheet)
	//	if len(rows) == 0 {
	//		continue
	//	}
	//	keyRow := rows[1]
	//	ctrlRow := rows[2]
	//
	//	var datas = []*map[string]interface{}{}
	//	for i := 4;i<len(rows);i++ {
	//		row := rows[i]
	//		data := map[string]interface{}{}
	//		for c :=0; c<len(ctrlRow); c++ {
	//			ctrlCell := ctrlRow[c]
	//			if len(ctrlCell) > 0 {
	//				key := keyRow[c]
	//				if len(key) == 0 {
	//					continue
	//				}
	//				cell := row[c]
	//				if c == 0 && len(cell) == 0{
	//					data = nil
	//					break
	//				}
	//				if ctrlCell[0] == 'i' {
	//					val, ok := strconv.Atoi(cell)
	//					if ok != nil {
	//						val = 0
	//					}
	//					data[key] = val
	//				}else if ctrlCell[0] == 'f' {
	//					val, ok := strconv.ParseFloat(cell, 64)
	//					if ok != nil {
	//						val = 0.0
	//					}
	//					data[key] = val
	//				}else{
	//					data[key] = cell
	//				}
	//			}
	//		}
	//		if data != nil {
	//			datas = append(datas, &data)
	//		}
	//	}
	//	retDatas[sheet] = datas
	//	break
	//}
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
	//GetCfgList("cfg_global")
	//{
	//	cfg :=GetCfgByKey("cfg_global", "hero_max")
	//	fmt.Println("cfg =", cfg)
	//}
	//{
	//	cfg :=GetCfgById("cfg_global_power", 12)
	//	fmt.Println("cfg =", cfg)
	//}
	//{
	//	cfg :=GetCfgWith("cfg_global_holiday1", "date", "2020-1-25")
	//	fmt.Println("cfg =", cfg)
	//}
}