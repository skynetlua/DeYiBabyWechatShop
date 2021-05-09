package config

import (
	// "bestsell/common"
	"fmt"
)

type CfgChannel struct {
    //{{1
	Id int `json:"id"`
	Appid string `json:"appid"`
	Subdomain string `json:"subdomain"`
	//}}1
}

type CfgChannelSlice []*CfgChannel
var cfgChannelSlice CfgChannelSlice

func GetCfgChannelSlice()*CfgChannelSlice  {
	return &cfgChannelSlice
}

func InitChannel()  {
	cfgs := *LoadJson("cfg_channel")
	if cfgs == nil {
		fmt.Println("json cfg_channel is lost")
		return
	}
	var _cfgChannelSlice CfgChannelSlice
	for _, cfg := range cfgs {
		item := &CfgChannel{}
        //{{2
		item.Id = cfg["id"].(int)
		item.Appid = cfg["appid"].(string)
		item.Subdomain = cfg["subdomain"].(string)
		//}}2

		_cfgChannelSlice = append(_cfgChannelSlice, item)
	}
	cfgChannelSlice = _cfgChannelSlice
}

//func init() {
//	InitChannel()
//}
