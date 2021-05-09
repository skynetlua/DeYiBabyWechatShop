package module

import (
	"bestsell/common"
	"bestsell/mysqld"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type FileConfig struct {
	Name  	string
	DstMd5 	string
	PublicPath string
}

// api/goods/1/pic/0.png 30a2f8e549736ed4bb90abd7b13ec76b
func (p *FileConfig)init(){
	dstPath := p.getDstPath(true)
	dirLen := len(common.StaticPath)
	if len(dstPath) >= dirLen {
		p.PublicPath = dstPath[dirLen:]
	} else {
		p.PublicPath = ""
	}
	// fmt.Println("[file]loadConfig Name =", p.Name)
	// fmt.Println("[file]loadConfig DstMd5 =", p.DstMd5)
	// fmt.Println("[file]loadConfig srcPath =", p.getSrcPath())
	// fmt.Println("[file]loadConfig dstPath =", dstPath)
	// fmt.Println("[file]loadConfig PublicPath =", p.PublicPath)
}

func (p *FileConfig)getSrcPath()string {
	// fmt.Println("[file]getSrcPath Name =", p.Name)
	// fmt.Println("[file]getSrcPath subDomain =", subDomain)
	idx := strings.Index(p.Name, subDomain)
	fileName := p.Name[idx+len(subDomain)+1:]
	srcPath := path.Join(common.ApiPath, fileName)
	// fmt.Println("[file]getSrcPath fileName =", fileName)
	// fmt.Println("[file]getSrcPath srcPath =", srcPath)
	return srcPath
}

func (p *FileConfig)getDstExtName()string {
	extName := path.Ext(p.Name)
	return extName
}

func (p *FileConfig)getDstPath(isCreate bool)string {
	fileName := p.DstMd5
	if len(fileName) == 0 {
		return ""
	}
	folder := fileName[0:2]
	folder = path.Join(common.ApiPath, "public" ,folder)
	if isCreate {
		common.CreateDir(folder)
	}
	extName := p.getDstExtName()
	dstPath := path.Join(folder, fileName+extName)
	return dstPath
}

func (p *FileConfig)removeDstFile() {
	filePath := p.getDstPath(false)
	if len(filePath) > 0 {
		RemoveFile(filePath)
	}
}

func (p *FileConfig)createDstFile() {
	dstPath := p.getDstPath(true)
	filePath := p.getSrcPath()
	copyFile(filePath, dstPath)

	fmt.Println("[file]loadConfig createDstFile =", p.Name)
}

func (p *FileConfig)isExistSrcFile() bool {
	srcPath := p.getSrcPath()
	if common.Exists(srcPath) {
		return true
	}
	return false
}

func (p *FileConfig)isExistDstFile() bool {
	dstPath := p.getDstPath(true)
	if common.Exists(dstPath) {
		return true
	}
	return false
}

func createFileConfig(name string, dstMd5 string)*FileConfig {
	cfg := &FileConfig{
		Name:   name,
		DstMd5: dstMd5,
	}
	cfg.init()
	return cfg
}

var fileConfigMap map[string]*FileConfig
var subDomain string
var isNeedSave bool
var locker sync.RWMutex

//var goodsResIds []string


//func GetGoodsResIds() *[]string {
//	var resIds []string
//	locker.Lock()
//	resIds = goodsResIds[:]
//	locker.Unlock()
//	return &resIds
//}

func loadConfig()  {
	subDomain = common.Config.Subdomain
	configPath := path.Join(common.AssetPath, "config.txt")
	fmt.Println("[file]loadConfig configPath =", configPath)

	if !common.Exists(configPath) {
		fmt.Println("[file]loadConfig configPath no exist")
		return
	}

	bytes, err := ioutil.ReadFile(configPath)
	content := ""
	if err != nil {
		fmt.Println("loadConfig err =", err)
		return
	} else {
		content = string(bytes)
	}
	items := strings.Split(content, "\n")
	locker.Lock()
	fileConfigMap = make(map[string]*FileConfig)
	for _,item := range items {
		if len(item) > 2 {
			parts := strings.Split(item, " ")
			if len(parts) >= 2 {
				fileName := parts[0]
				if len(fileName) <= 1 {
					continue
				}
				cfg := createFileConfig(fileName, parts[1])
				if !cfg.isExistSrcFile() {
					fmt.Println("[file]loadConfig isExistSrcFile not fileName=", fileName)
					continue
				}
				if !cfg.isExistDstFile() {
					cfg.createDstFile()
				}
				fileConfigMap[cfg.Name] = cfg
			}
		}
	}
	locker.Unlock()
	isNeedSave = false
}

func saveConfig() {
	if !isNeedSave {
		return
	}
	configPath := path.Join(common.AssetPath, "config.txt")
	var items []string
	locker.Lock()
	for _,fileCfg :=range fileConfigMap {
		item := fileCfg.Name + " " + fileCfg.DstMd5
		items = append(items, item)
	}
	locker.Unlock()

	content := strings.Join(items, "\n")
	err := ioutil.WriteFile(configPath, []byte(content), os.ModePerm)
	if err != nil {
		 fmt.Println("saveConfig err =", err)
	}
}

func md5SumFile(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}
	value := md5.Sum(data)
	md5str := fmt.Sprintf("%x", value)
	return md5str
}

func copyFile(src, des string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
	}
	defer srcFile.Close()

	desFile, err := os.Create(des)
	if err != nil {
		fmt.Println(err)
	}
	defer desFile.Close()
	return io.Copy(desFile, srcFile)
}

func addFileCfg(fileCfg *FileConfig) {
	locker.Lock()
	fileConfigMap[fileCfg.Name] = fileCfg
	locker.Unlock()
}

func getFileCfg(fileName string)*FileConfig {
	locker.Lock()
	fileCfg := fileConfigMap[fileName]
	locker.Unlock()
	return fileCfg
}

func removeFileCfg(fileName string) {
	locker.Lock()
	delete(fileConfigMap, fileName)
	locker.Unlock()
}

func addFileConfigByFile(filePath string)bool {
	if !common.Exists(filePath) {
		return false
	}
	idx := strings.Index(filePath, common.ApiPath)
	if idx >= 0 {
		fileMd5 := md5SumFile(filePath)
		fileName := path.Join(subDomain, filePath[idx+len(common.ApiPath):])
		fileCfg := getFileCfg(fileName)
		if fileCfg != nil {
			if strings.Compare(fileCfg.DstMd5, fileMd5) == 0 {
				if !fileCfg.isExistDstFile() {
					fileCfg.createDstFile()
				}
			}else{
				fileCfg.removeDstFile()
				fileCfg.DstMd5 = fileMd5
				fileCfg.init()
				fileCfg.createDstFile()
				isNeedSave = true
			}
		}else{
			fileCfg = createFileConfig(fileName, fileMd5)
			fileCfg.createDstFile()
			addFileCfg(fileCfg)
			isNeedSave = true
		}
	}
	return true
}

func removeUnusedFiles() {
	for key,fileCfg := range fileConfigMap {
		if !fileCfg.isExistSrcFile() {
			fileCfg.removeDstFile()
			removeFileCfg(key)
			isNeedSave = true
			fmt.Println("removeUnusedFiles:", key)
		}
	}
}

func makeGoodsFile2PublicFiles(){
	goodsPath := path.Join(common.ApiPath, "goods")
	files, err := ioutil.ReadDir(goodsPath)
	if err != nil {
		fmt.Println("makeGoodsFile2PublicFiles read dir fail:", err)
		return
	}

	locker.Lock()
	var goodsFolders []string
	for _, file := range files {
		fileName := file.Name()
		if fileName[0] == '.' {
			continue
		}
		if file.IsDir() {
			fullDir := path.Join(goodsPath, fileName)
			goodsFolders = append(goodsFolders, fullDir)
			//goodsResIds = append(goodsResIds, fileName)
		}
	}
	locker.Unlock()
	for _,folder := range goodsFolders {
		// fmt.Println("makeAllFile2PublicFiles folder =", folder)
		spiderFiles(folder)
	}
}

func getFilePublicPath(fileName string)string {
	if fileName[0] == '/' || fileName[0] == '\\' {
		fileName = fileName[1:]
	}
	fileCfg := getFileCfg(fileName)
	if fileCfg == nil {
		cfg := &FileConfig {
			Name: fileName,
		}
		filePath := cfg.getSrcPath()
		ret := addFileConfigByFile(filePath)
		if ret {
			delaySaveFileConfigs()
			fileCfg = getFileCfg(fileName)
		}
		if fileCfg == nil {
			fmt.Println("getFilePublicPath 文件不存在:", fileName)
			return ""
		}
	}
	return fileCfg.PublicPath
}

func RemoveFilePublicPath(fileName string) {
	if fileName[0] == '/' || fileName[0] == '\\' {
		fileName = fileName[1:]
	}
	removeFileCfg(fileName)
	delaySaveFileConfigs()
}

var isSaving bool
func delaySaveFileConfigs()  {
	if isSaving {
		return
	}
	isSaving = true
	go func() {
		select {
		case <-time.After(time.Second *10):
			saveConfig()
			isSaving = false
		}
	}()
}

// func FileMake()  {
// 	loadConfig()
// 	makeGoodsFile2PublicFiles()
// 	delaySaveFileConfigs()
// }

func init() {
	loadConfig()
	makeGoodsFile2PublicFiles()
	// resIds := GetGoodsResIds()
	// fmt.Println("resIds:", resIds)
	saveConfig()

	mysqld.Module_GetFilePublicPath = getFilePublicPath
}

func continueSpiderFiles(pathName string) error{
	return spiderFiles(pathName)
}

func spiderFiles(pathName string) error{
	files, err := ioutil.ReadDir(pathName)
	if err != nil {
		fmt.Println("spiderFiles read dir fail:", err)
		return err
	}
	for _, file := range files {
		fileName := file.Name()
		if fileName[0] == '.' {
			continue
		}
		fullName := path.Join(pathName, fileName)
		if file.IsDir() {
			err = continueSpiderFiles(fullName)
			if err != nil {
				fmt.Println("spiderFiles read dir fail:", err)
				return err
			}
		} else {
			addFileConfigByFile(fullName)
		}
	}
	return nil
}

func RemoveFile(filePath string)  {
	if !common.Exists(filePath) {
		return
	}
	err := os.Remove(filePath)
	if err != nil{
		fmt.Println("RemoveFile filePath =", filePath)
	}
}

// func continueAllFiles(pathName string, filePaths []string)([]string, error){
// 	return GetAllFiles(pathName, filePaths)
// }

// func GetAllFiles(pathName string, filePaths []string)([]string, error){
// 	rd, err := ioutil.ReadDir(pathName)
// 	if err != nil {
// 		fmt.Println("read dir fail:", err)
// 		return filePaths, err
// 	}
// 	for _, fi := range rd {
// 		if fi.Name()[0] == '.' {
// 			continue
// 		}
// 		if fi.IsDir() {
// 			fullDir := path.Join(pathName, fi.Name())
// 			filePaths, err = continueAllFiles(fullDir, filePaths)
// 			if err != nil {
// 				fmt.Println("read dir fail:", err)
// 				return filePaths, err
// 			}
// 		} else {
// 			fullName := path.Join(pathName, fi.Name())
// 			filePaths = append(filePaths, fullName)
// 		}
// 	}
// 	return filePaths, nil
// }

// func removeUrlHost(url string)string {
// 	if len(url) == 0 {
// 		return ""
// 	}
// 	hostUrl := common.StaticPath
// 	if strings.Contains(url, hostUrl) {
// 		url = url[len(hostUrl):]
// 	}
// 	return url
// }

// func removeUrlHosts(urlPaths []string)string {
// 	if len(urlPaths) == 0{
// 		return ""
// 	}
// 	var items []string
// 	for _,tmp := range urlPaths {
// 		item := removeUrlHost(tmp)
// 		items = append(items, item)
// 	}
// 	return strings.Join(items, ";")
// }

// func MakeGoodsFiles(dbGoodsInfo *mysqld.DBGoodsInfo) bool {
// 	if len(dbGoodsInfo.ResId) == 0 {
// 		return false
// 	}
// 	goodsPath := path.Join(common.ApiPath, "goods", dbGoodsInfo.ResId)
// 	if !common.Exists(goodsPath) {
// 		return false
// 	}
// 	iconPath := path.Join(goodsPath, "icon")
// 	picPath := path.Join(goodsPath, "pic")
// 	contentPath := path.Join(goodsPath, "content")

// 	var icon string
// 	if common.Exists(iconPath) {
// 		var iconFiles []string
// 		iconFiles,_ = common.GetAllFiles(iconPath, iconFiles)
// 		if len(iconFiles) > 0 {
// 			icon = iconFiles[0]
// 			icon = removeUrlHost(icon)
// 		}
// 	}
// 	var pics string
// 	if common.Exists(picPath) {
// 		var picFiles []string
// 		picFiles, _ = common.GetAllFiles(picPath, picFiles)
// 		pics = removeUrlHosts(picFiles)
// 	}
// 	var contents string
// 	if common.Exists(contentPath) {
// 		var contentFiles []string
// 		contentFiles,_ = common.GetAllFiles(contentPath, contentFiles)
// 		contents = removeUrlHosts(contentFiles)
// 	}
// 	fmt.Println("MakeGoodsFiles icon:", icon)
// 	fmt.Println("MakeGoodsFiles pics:", pics)
// 	fmt.Println("MakeGoodsFiles contents:", contents)

// 	dbGoodsInfo.BeginWrite()
// 	dbGoodsInfo.Icon = icon
// 	dbGoodsInfo.Pics = pics
// 	dbGoodsInfo.Contents = contents
// 	dbGoodsInfo.ResId = ""
// 	dbGoodsInfo.EndWrite()
// 	dbGoodsInfo.Save()
// 	return true
// }

//func init() {
//	dbGoodsInfo := &mysqld.DBGoodsInfo{
//		ResId:"5",
//	}
//	MakeGoodsFiles(dbGoodsInfo)
//}