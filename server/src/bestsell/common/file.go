package common

import (
	"os"
	"path/filepath"
	"time"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"sort"
)

func CreateDateDir(dirPath string) string {
	folderName := time.Now().Format("20060102")
	folderPath := filepath.Join(dirPath, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.Mkdir(folderPath, 0777) //0777也可以os.ModePerm
		os.Chmod(folderPath, os.ModePerm)
	}
	return folderPath
}

func CreateTokenDir(dirPath string, token string) string {
	folderPath := filepath.Join(dirPath, token)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		//0777也可以os.ModePerm
		os.Mkdir(folderPath, 0777)
		os.Chmod(folderPath, os.ModePerm)
	}
	return folderPath
}

func CreateDir(folderPath string) string {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.Mkdir(folderPath, 0777)
		os.Chmod(folderPath, os.ModePerm)
	}
	return folderPath
}

func FillPathHeader(filePath string) string {
	urlPath := filePath[len(StaticPath):]
	if urlPath[0] != '/' && urlPath[0] != '\\' {
		urlPath = "/"+urlPath
	}
	return urlPath
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func continueAllFiles(pathName string, filePaths []string)([]string, error){
	return GetAllFiles(pathName, filePaths)
}

func GetAllFiles(pathName string, filePaths []string)([]string, error){
	rd, err := ioutil.ReadDir(pathName)
	if err != nil {
		fmt.Println("GetAllFiles1 read dir fail:", err)
		return filePaths, err
	}
	for _, fi := range rd {
		if fi.Name()[0] == '.' {
			continue
		}
		if fi.IsDir() {
			fullDir := path.Join(pathName, fi.Name())
			filePaths, err = continueAllFiles(fullDir, filePaths)
			if err != nil {
				fmt.Println("GetAllFiles2 read dir fail:", err)
				return filePaths, err
			}
		} else {
			fullName := path.Join(pathName, fi.Name())
			filePaths = append(filePaths, fullName)
		}
	}
	return filePaths, nil
}

func SortFileList(filePaths *[]string) {
	sort.Slice(*filePaths, func(i, j int) bool {
    	filePath1 := (*filePaths)[i]
        filePath2 := (*filePaths)[j]
        if len(filePath1) == len(filePath2) {
        	for i := 0; i < len(filePath1); i++ {
        		if filePath1[i] != filePath2[i] {
        			return filePath1[i] < filePath2[i]
        		}
        	}
        	return false
        }
        return len(filePath1) < len(filePath2)
    })
}

func GetFilesList(pathName string, filePaths []string)([]string, error){
	rd, err := ioutil.ReadDir(pathName)
	if err != nil {
		fmt.Println("GetAllFiles1 read dir fail:", err)
		return filePaths, err
	}
	for _, fi := range rd {
		if fi.Name()[0] == '.' {
			continue
		}
		fullName := path.Join(pathName, fi.Name())
		filePaths = append(filePaths, fullName)
	}
	return filePaths, nil
}

func RemoveUrlHost(url string)string {
	if len(url) == 0 {
		return ""
	}
	hostUrl := StaticPath
	if strings.Contains(url, hostUrl) {
		url = url[len(hostUrl):]
	}
	return url
}

func RemoveUrlHosts(urlPaths []string)string {
	if len(urlPaths) == 0{
		return ""
	}
	var items []string
	for _,tmp := range urlPaths {
		item := RemoveUrlHost(tmp)
		items = append(items, item)
	}
	return strings.Join(items, ";")
}