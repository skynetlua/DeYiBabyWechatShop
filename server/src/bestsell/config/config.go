package config

import (
	"bestsell/common"
	"fmt"
	"strings"
	"io/ioutil"
	"path"
	"encoding/json"
)

func StartServer(ch *chan bool) {
	go func() {
		fmt.Println("StartServer config")
		InitAddress()
		InitGoodsData()
		// InitGoodsLogistics()
		// initGoodsCategory()
		// InitGoodsDetail()
		// initGoods()
		InitUiBanner()
		(*ch) <- true
	}()
	<-(*ch)
}

func converUrl(url string)string {
	if !strings.Contains(url, "http") {
		if strings.HasPrefix(url, "/") {
			url = common.StaticUrl+url
		}else{
			url = common.StaticUrl+"/"+url
		}
	}
	return url
}

func LoadJson(cfgName string)*[]map[string]interface{} {
	_jsonPath := common.JsonPath
	jsonPath := path.Join(_jsonPath, cfgName+".json")
	bytes, _ := ioutil.ReadFile(jsonPath)
	if bytes == nil {
		fmt.Println("LoadJson Failed jsonPath =", jsonPath)
		return nil
	}
	fmt.Println("LoadJson Success jsonPath =", jsonPath)
	datas := []map[string]interface{}{}
	err := json.Unmarshal(bytes, &datas)
	if err != nil {
		panic(err)
	}
	return &datas
}