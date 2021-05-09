package config

import (
	// "bestsell/common"
	"fmt"
	"sort"
)

type CfgUiBanner struct {
    //{{1
	Id int `json:"id"`
	Type string `json:"type"`
	Order int `json:"order"`
	PicUrl string `json:"picUrl"`
	LinkUrl string `json:"linkUrl"`
	Status int `json:"status"`
	//}}1
}

type CfgUiBannerSlice []*CfgUiBanner
var cfgUiBannerSlice CfgUiBannerSlice

func (p *CfgUiBannerSlice) Len() int { return len(*p) }
func (p *CfgUiBannerSlice) Swap(i, j int) { (*p)[i], (*p)[j] = (*p)[j], (*p)[i] }
func (p *CfgUiBannerSlice) Less(i, j int) bool { return (*p)[i].Order < (*p)[j].Order }

var cfgUiBannerMap map[string]*CfgUiBannerSlice

func GetCfgUiBannerSlice()*CfgUiBannerSlice {
	return &cfgUiBannerSlice
}

func GetCfgUiBannerSliceByMap(key string)*CfgUiBannerSlice {
	return cfgUiBannerMap[key]
}

func InitUiBanner()  {
	cfgs := *LoadJson("cfg_uibanner")
	if cfgs == nil {
		fmt.Println("json cfg_uibanner is lost")
		return
	}
	var _cfgUiBannerSlice CfgUiBannerSlice
	_cfgUiBannerMap := make(map[string]*CfgUiBannerSlice)
	for _, cfg := range cfgs {
		item := &CfgUiBanner{}
        //{{2
		item.Id = int(cfg["id"].(float64))
		item.Type = cfg["type"].(string)
		item.Order = int(cfg["order"].(float64))
		item.PicUrl = cfg["picUrl"].(string)
		item.LinkUrl = cfg["linkUrl"].(string)
		item.Status = int(cfg["status"].(float64))
		//}}2

		if item.Status != 1 {
			continue
		}

		item.PicUrl = converUrl(item.PicUrl)
		_cfgUiBannerSlice = append(_cfgUiBannerSlice, item)

		arr,ok := _cfgUiBannerMap[item.Type]
		if !ok {
			arr = &CfgUiBannerSlice{}
		}
		*arr = append(*arr, item)
		_cfgUiBannerMap[item.Type] = arr
		
		// fmt.Println("json cfg_uibanner item =", item)
	}
	cfgUiBannerSlice = _cfgUiBannerSlice
	cfgUiBannerMap = _cfgUiBannerMap
	for _, bans := range cfgUiBannerMap {
		sort.Sort(bans)
	}
}

//func init() {
//	InitUiBanner()
//}
