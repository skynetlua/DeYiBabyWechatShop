package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)


var PROGOPATH = "PROGOPATH"
var CONFIGFILE = "config.json"



type ConfigInfo struct {
	ServerAddr 	string
	Domain string

	Mode        string
	MySqlAddr 	string
	StaticUrl   string
	// ProjectPath string
	HttpsCrt 	string
	HttpsKey 	string

	Appid   	string
	Subdomain 	string
	Enable      map[string]int `json:"Enable"`
}

var Config ConfigInfo

var ProjectPath string
var AssetPath string
var StaticPath string
var ExcelPath  string
var JsonPath  string
var ViewPath string
var GoodsPath string
var QRCodePath string
var UploadPath string
var UploadGoodsPath string
var UploadCategoryPath string
var ApiPath string
var PicturePath string
var IconPath string

var SrcPath string
var StaticUrl string

var isInit = false
func init() {
	if isInit {
		return
	}
	isInit = true
	fmt.Println("common:init config")
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	UnmarshalJson()
	// if len(Config.ProjectPath) < 2 {
	// 	ProjectPath = getProjectPath()
	// 	Config.ProjectPath = ProjectPath
	// }else{
	// ProjectPath = Config.ProjectPath
	// }
	fmt.Println("common:init ProjectPath:", ProjectPath)
 
	AssetPath = path.Join(ProjectPath, "/asset")
	StaticPath = path.Join(ProjectPath, "/static")
	SrcPath = path.Join(ProjectPath, "src")

	ExcelPath = path.Join(AssetPath, "/xlsx/")
	ViewPath = path.Join(AssetPath, "/view/")
	GoodsPath = path.Join(AssetPath, "/goods/")
	JsonPath = path.Join(AssetPath, "/json/")

	ApiPath = path.Join(StaticPath, Config.Subdomain)
	QRCodePath = path.Join(ApiPath ,"/qrcode")
	UploadPath = path.Join(ApiPath ,"/upload")
	UploadGoodsPath = path.Join(ApiPath ,"/auto")
	UploadCategoryPath = path.Join(ApiPath ,"/category")

	PicturePath = path.Join(ApiPath ,"/picture")
	IconPath = path.Join(ApiPath ,"/icon")

	CreateDir(UploadGoodsPath)
	CreateDir(UploadPath)

	StaticUrl = Config.StaticUrl

	fmt.Println("AssetPath:", AssetPath)
	fmt.Println("StaticPath:", StaticPath)
	fmt.Println("SrcPath:", SrcPath)
	fmt.Println("ExcelPath:", ExcelPath)
	fmt.Println("ViewPath:", ViewPath)
	fmt.Println("GoodsPath:", GoodsPath)
	fmt.Println("ApiPath:", ApiPath)
	fmt.Println("QRCodePath:", QRCodePath)
	fmt.Println("UploadPath:", UploadPath)
	fmt.Println("UploadGoodsPath:", UploadGoodsPath)
	fmt.Println("UploadCategoryPath:", UploadCategoryPath)
	fmt.Println("StaticUrl:", StaticUrl)
}

func getProjectPath() string {
	var projectPath = os.Getenv(PROGOPATH)
	if len(projectPath) == 0 {
		execPath := os.Args[0]
		projectPath = path.Dir(execPath)
		cfgPath := path.Join(projectPath, CONFIGFILE)
		fmt.Println("getProjectPath cfgPath:", cfgPath)
		if Exists(cfgPath) {
			if len(projectPath) > 2 {
				projectPath = path.Dir(projectPath)
				projectPath = path.Dir(projectPath)
			}
			fmt.Println("getProjectPath1 projectPath:", projectPath)
			return projectPath
		}
		projectPath = ""
		goPath := os.Getenv("GOPATH")
		items := strings.Split(goPath, ":")
		for _,item := range items{
			if !strings.HasSuffix(item, "/go") {
				projectPath = item
				break
			}
		}
		if len(projectPath) == 0 {
			projectPath = "./"
		}
	}
	fmt.Println("getProjectPath2 projectPath:", projectPath)
	return projectPath
}

func UnmarshalJson() {
	var configPath string
	var mode string
	fileName := "./"+CONFIGFILE
	count := len(os.Args)
	for i:=1;i<count;i++ {
		arg := os.Args[i]
		switch arg {
		case "-m":
			i++
			mode = os.Args[i]
		case "-c":
			i++
			configPath = os.Args[i]
		case "-p":
			i++
			ProjectPath = os.Args[i]
		}
	}
	//for i,v := range os.Args {
	//	fmt.Println("i =", i, "v =",v)
	//}
	//flag.StringVar(&configPath, "c", "./"+fileName, "配置文件，默认为./config.json")
	//flag.StringVar(&mode, "m", "", "模式:api,res,tool")
	//flag.Parse()
	if len(ProjectPath) == 0 {
		ProjectPath = getProjectPath()
	}
	bytes, _ := ioutil.ReadFile(configPath)
	if bytes == nil {
		configPath = path.Join(ProjectPath, "src","bestsell", fileName)
		bytes, _ = ioutil.ReadFile(configPath)
	}
	fmt.Println("configPath:", configPath)
	err := json.Unmarshal(bytes, &Config)
	if err != nil {
		panic(err)
	}
	if len(mode) > 0 {
		Config.Mode = mode
	}
}
