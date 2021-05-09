package config

import (
	// "bestsell/common"
	"fmt"
)

type CfgCategory struct {
    //{{1
	Id int `json:"id"`
	Name string `json:"name"`
	Level int `json:"level"`
	Type string `json:"type"`
	Order int `json:"order"`
	IsUse int `json:"isUse"`
	Icon string `json:"icon"`
	//}}1
}

var cfgCategorySlice []*CfgCategory
var cfgCategoryMap map[int]*CfgCategory

func GetCfgCategorySlice()*[]*CfgCategory  {
	return &cfgCategorySlice
}

func GetCfgCategory(id int)*CfgCategory {
    return cfgCategoryMap[id]
	//for _,item := range cfgCategorySlice{
	//	if item.Id == id {
	//		return item
	//	}
	//}
	//return nil
}

func InitCategory()  {
	cfgs := *LoadJson("cfg_category")
	if cfgs == nil {
		fmt.Println("json cfg_category is lost")
		return
	}
	var _cfgCategorySlice []*CfgCategory
	var _cfgCategoryMap = make(map[int]*CfgCategory)
	for _, cfg := range cfgs {
		item := &CfgCategory{}
        //{{2
		item.Id = cfg["id"].(int)
		item.Name = cfg["name"].(string)
		item.Level = cfg["level"].(int)
		item.Type = cfg["type"].(string)
		item.Order = cfg["order"].(int)
		item.IsUse = cfg["isUse"].(int)
		item.Icon = cfg["icon"].(string)
		//}}2
		_cfgCategorySlice = append(_cfgCategorySlice, item)
		_cfgCategoryMap[item.Id] = item
	}
	cfgCategorySlice = _cfgCategorySlice
	cfgCategoryMap = _cfgCategoryMap
}

func init() {
	// InitCategory()
}

