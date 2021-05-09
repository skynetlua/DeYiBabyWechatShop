package config

import (
	// "bestsell/common"
	"fmt"
)

type CfgUi struct {
    //{{1
	Id int `json:"id"`
	Type string `json:"type"`
	Order int `json:"order"`
	PicUrl string `json:"picUrl"`
	LinkUrl string `json:"linkUrl"`
	Status int `json:"status"`
	//}}1
}

var cfgUiSlice []*CfgUi
var cfgUiMap map[int]*CfgUi

func GetCfgUiSlice()*[]*CfgUi  {
	return &cfgUiSlice
}

func GetCfgUi(id int)*CfgUi {
    return cfgUiMap[id]
	//for _,item := range cfgUiSlice{
	//	if item.Id == id {
	//		return item
	//	}
	//}
	//return nil
}

func InitUi()  {
	cfgs := *LoadJson("cfg_ui")
	if cfgs == nil {
		fmt.Println("json cfg_ui is lost")
		return
	}
	var _cfgUiSlice []*CfgUi
	var _cfgUiMap = make(map[int]*CfgUi)
	for _, cfg := range cfgs {
		item := &CfgUi{}
        //{{2
		item.Id = cfg["id"].(int)
		item.Type = cfg["type"].(string)
		item.Order = cfg["order"].(int)
		item.PicUrl = cfg["picUrl"].(string)
		item.LinkUrl = cfg["linkUrl"].(string)
		item.Status = cfg["status"].(int)
		//}}2
		_cfgUiSlice = append(_cfgUiSlice, item)
		_cfgUiMap[item.Id] = item
	}
	cfgUiSlice = _cfgUiSlice
	cfgUiMap = _cfgUiMap
}

func init() {
	// InitUi()
}

