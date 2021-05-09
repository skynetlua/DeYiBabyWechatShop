package generate

import (
	"bestsell/common"
	// "bestsell/config"
	"bestsell/router"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

var configDir string
var handleDir string
var templateFile string

var mysqls = []string{
	"player",
	"playerInfo",
	"order",
	"cart",
	"favorite",
	"address",
	"goods",
	"reputation",
	"refund",
	"cashLog",
	"team",
	"teamInfo",
	"myTeam",
	"record",
	"operation",
	"goodsInfo",
	"category",
	"logistics",

	"record",
	"operation",

	"notice",
	"coupon",
	"myCoupon",
}

var playerMysqls = []string{
	"member",
	"commission",
	"teamLog",
	"teamNotice",
	"recommend",
	"shop",
	//"operation",
	//"withdrawLog",
}

//自动生成数据库数据结构
func generateMysqls() {
	handleDir = path.Join(common.SrcPath, "/bestsell/mysqld")
	templateFile = path.Join(common.AssetPath, "/template/tempmysql.txt")
	for _,item := range mysqls{
		fileName := strings.ToLower(item)
		goFile := path.Join(handleDir, fileName+".go")
		if isExists(goFile) {
			continue
		}
		bytes := readFile(templateFile)
		content := string(bytes)
		item = strings.ToUpper(item[0:1])+item[1:]
		content = strings.ReplaceAll(content, "{ClassName}", item)
		writeFile(goFile, []byte(content))
	}

	templateFile = path.Join(common.AssetPath, "/template/tempplayermysql.txt")
	for _,item := range playerMysqls{
		fileName := strings.ToLower(item)
		goFile := path.Join(handleDir, fileName+".go")
		if isExists(goFile) {
			continue
		}
		bytes := readFile(templateFile)
		content := string(bytes)
		item = strings.ToUpper(item[0:1])+item[1:]
		content = strings.ReplaceAll(content, `{ClassName}`, item)
		writeFile(goFile, []byte(content))
	}
}

//自动生成excel数据结构
func generateCfg(cfgName string)  {
	idx := strings.Index(cfgName, "cfg_")
	if idx != 0 {
		panic("generateCfg cfgName ="+cfgName)
	}
	name := cfgName[4:]
	goFileName := name+".go"
	goFileName = path.Join(configDir, goFileName)
	var names []string
	tmps := strings.Split(name, "_")
	for _,tmp := range tmps{
		tmp = strings.ToUpper(tmp[0:1])+tmp[1:]
		names = append(names, tmp)
	}
	name = strings.Join(names, "")

	fmt.Println("generateCfg goFileName:", goFileName)
	fmt.Println("generateCfg name:", name)

	var content string
	isNew := true
	if !isExists(goFileName) {
		bytes := readFile(templateFile)
		content = string(bytes)
		content = strings.ReplaceAll(content, "{ConfigName}", cfgName)
		content = strings.ReplaceAll(content, "{ClassName}", name)
	}else{
		content = string(readFile(goFileName))
		isNew = false
	}

	sidx1 := strings.Index(content, `//{`+`{1`)
	eidx1 := strings.Index(content, `//}}1`)
	sidx2 := strings.Index(content, `//{`+`{2`)
	eidx2 := strings.Index(content, `//}}2`)

	originCt1 := content[sidx1:eidx1]
	originCt2 := content[sidx2:eidx2]

	theLen := len(content)
	otherCt := content[eidx1:theLen]
	theLen = len(otherCt)
	eidx := strings.Index(otherCt[6:theLen], "}")
	endxIdx := eidx+6
	if endxIdx > theLen {
		endxIdx = theLen
	}
	otherCt = otherCt[:endxIdx]
	excelStruct := common.ExcelStructs[cfgName]
	var lines1 []string
	var lines2 []string
	for _,obj := range(excelStruct) {
		lkeyName := obj.Key
		ukeyName := obj.Key[0:1]
		ukeyName = strings.ToUpper(ukeyName)+obj.Key[1:]
		var typeName string
		switch obj.Ctrl {
		case "i":
			typeName = "int"
		case "f":
			typeName = "float64"
		default:
			typeName = "string"
		}
		if !strings.Contains(originCt1, ukeyName+" ") {
			if strings.Contains(otherCt, ukeyName+" ") {
				continue
			}
		}
		line1 := "\t"+ukeyName+" "+typeName+" `json:\""+lkeyName+"\"`"
		lines1 = append(lines1, line1)

		line2 := "\t\titem."+ukeyName+" = cfg[\""+lkeyName+"\"].("+typeName+")"
		lines2 = append(lines2, line2)
	}

	newCt1 := strings.Join([]string{`//{`+"{1\n",strings.Join(lines1,"\n"), "\n"}, "")+"\t"
	newCt2 := strings.Join([]string{`//{`+"{2\n",strings.Join(lines2,"\n"), "\n"}, "")+"\t\t"
	if !isNew {
		if strings.Compare(originCt1, newCt1) == 0 && strings.Compare(originCt2, newCt2) == 0 {
			return
		}
	}
	contents := []string{
		content[:sidx1],
		newCt1,
		content[eidx1:sidx2],
		newCt2,
		content[eidx2:len(content)],
	}
	writeFile(goFileName, []byte(strings.Join(contents,"")))
}

//自动生成路由方法
func generateRouteMethod(routeParam []string) {
	route := routeParam[0]
	items := strings.Split(route, "/")
	items = items[1:]
	hdFilePath := path.Join(handleDir, items[0]+".go")
	items[0] = strings.ToUpper(items[0][0:1])+items[0][1:]
	mpath := strings.Join(items, "_")

	newCt :="//=>"+strings.Join(routeParam," ")+"\nfunc "+mpath+"(ctx iris.Context, sess *common.BSSession"
	var content string
	isNew := true
	if !isExists(hdFilePath) {
		content = `package handle
import (
	"github.com/kataras/iris/v12"
	"bestsell/common"
)
`
	}else{
		content = string(readFile(hdFilePath))
		isNew = false
	}
	var originCt string
	var sidx int
	var eidx int
	if !isNew {
		var idx = 0
		for idx < len(content) {
			ptr := content[idx:]
			sidx = strings.Index(ptr, "//=>")
			if sidx < 0 {
				break
			}
			ptr = ptr[sidx:]
			eidx = strings.Index(ptr, ")")
			if eidx < 0 {
				break
			}
			tmp := ptr[:eidx]
			if strings.Contains(tmp, mpath) {
				originCt = tmp
				sidx += idx
				eidx += sidx
				break
			}
			idx += sidx+eidx
		}
	}
	if len(originCt) > 10 {
		if strings.Compare(originCt, newCt) == 0 {
			return
		}
		contents := []string{
			content[:sidx],
			newCt,
			content[eidx:],
		}
		content = strings.Join(contents,"")
	}else{
		content = content+"\n"+newCt+") {\n\tempty(\""+route+"\")\n}\n"
	}
	writeFile(hdFilePath, []byte(content))
	fmt.Println("generateRouteMethod:", hdFilePath, route)
}

//自动生成路由文件
func generateRouteMap()  {
	routeMapFile := path.Join(common.SrcPath, "/bestsell/router/map.go")
	content := string(readFile(routeMapFile))
	sidx := strings.Index(content, `//{`+"{")
	eidx := strings.Index(content, `//}`+"}")
	originCt := content[sidx:eidx]
	var lines []string
	for _,data := range router.RoutesParams{
		route := data[0]
		items := strings.Split(route, "/")
		items = items[1:]
		items[0] = strings.ToUpper(items[0][0:1])+items[0][1:]
		mpath := strings.Join(items, "_")
		line := "\tROUTE(\""+route+"\", handle."+mpath+")"
		lines = append(lines, line)
		generateRouteMethod(data[:])
	}
	newCt := strings.Join([]string{`//{`+"{\n",strings.Join(lines,"\n"), "\n"}, "")+"\t"
	if strings.Compare(originCt, newCt) == 0 {
		return
	}
	header := content[:sidx]
	tail := content[eidx:]
	contents := []string{
		header,
		newCt,
		tail,
	}
	writeFile(routeMapFile, []byte(strings.Join(contents,"")))
	fmt.Println("generateRouteMap:", routeMapFile)
}

func Generate() {
	if strings.Compare(common.Config.Mode, "tool") != 0 {
		return
	}
	fmt.Println("Generate======================>")
	generateMysqls()

	configDir = path.Join(common.SrcPath, "/bestsell/config")
	handleDir = path.Join(common.SrcPath, "/bestsell/handle")
	templateFile = path.Join(common.AssetPath, "/template/template.txt")

	// for cfgName := range common.ExcelDatas {
	// 	ret := excludeFiles[cfgName]
	// 	if !ret {
	// 		generateCfg(cfgName)
	// 	}
	// }

	// ch := make(chan bool)
	// config.StartServer(&ch)
	// close(ch)

	generateRouteMap()
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func isFile(path string) bool {
	return !isDir(path)
}

var excludeFiles = map[string]bool{
}

func readFile(filePath string)[]byte  {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0777)
	var data []byte
	if err != nil {
		panic(err)
		return data
	}
	defer file.Close()
	buffer := make([]byte, 4096)
	for{
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		data = append(data, buffer[:n]...)
	}
	return data
}

func writeFile(filePath string, data []byte)  {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
		return
	}
	defer file.Close()
	n, err := file.Write(data)
	if err != nil && n != len(data) {
		panic(err)
	}
}
