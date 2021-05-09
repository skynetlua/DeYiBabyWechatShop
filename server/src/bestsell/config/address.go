package config

import (
	// "bestsell/common"
	"fmt"
)

type CfgAddress struct {
    //{{1
	Id int `json:"id"`
	Name string `json:"name"`
	Level int `json:"level"`
	Pid int `json:"pid"`
	Pinyin string `json:"pinyin"`
	Jianpin string `json:"jianpin"`
	FirstLetter string `json:"firstLetter"`
	//}}1
	subAddress []*CfgAddress
}

var cfgAddressSlice []*CfgAddress
var cfgAddressMap map[int]*CfgAddress
func GetCfgAddressSlice()*[]*CfgAddress  {
	return &cfgAddressSlice
}

func GetCfgAddress(id int)*CfgAddress {
	return cfgAddressMap[id]
}

func GetCfgAddressName(id int)string {
	cfgAddress := cfgAddressMap[id]
	if cfgAddress == nil {
		return ""
	}
	return cfgAddress.Name
}

func GetCfgAddressSub(id int)*[]*CfgAddress {
	cfgAddress := cfgAddressMap[id]
	if cfgAddress != nil {
		return &cfgAddress.subAddress
	}
	return nil
}

var cfgProvinceAddressSlice []*CfgAddress
func GetCfgProvinceAddressSlice()*[]*CfgAddress  {
	return &cfgProvinceAddressSlice
}

func InitAddress()  {
	cfgs := *LoadJson("cfg_address")
	if cfgs == nil {
		fmt.Println("json cfg_address is lost")
		return
	}
	var _cfgAddressMap = make(map[int]*CfgAddress)
	var _cfgAddressSlice []*CfgAddress
	for _, cfg := range cfgs {
		item := &CfgAddress{}
        //{{2
		item.Id = int(cfg["id"].(float64))
		item.Name = cfg["name"].(string)
		item.Level = int(cfg["level"].(float64))
		item.Pid = int(cfg["pid"].(float64))
		item.Pinyin = cfg["pinyin"].(string)
		item.Jianpin = cfg["jianpin"].(string)
		item.FirstLetter = cfg["firstLetter"].(string)
		//}}2
		_cfgAddressSlice = append(_cfgAddressSlice, item)

		_cfgAddressMap[item.Id] = item
	}
	cfgAddressSlice = _cfgAddressSlice
	cfgAddressMap = _cfgAddressMap

	var _cfgProvinceAddressSlice []*CfgAddress
	for _,item := range _cfgAddressMap {
		if item.Pid > 0 {
			pItem := _cfgAddressMap[item.Pid]
			if pItem != nil {
				pItem.subAddress = append(pItem.subAddress, item)
			}
		}else{
			_cfgProvinceAddressSlice = append(_cfgProvinceAddressSlice, item)
		}
	}

	cfgProvinceAddressSlice = _cfgProvinceAddressSlice
}

//func init() {
//	InitAddress()
//}
//}
