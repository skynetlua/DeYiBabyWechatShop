package config

import (
	"bestsell/common"
	// "bestsell/common"
	"fmt"
	"strconv"
)

type CfgGoodsData struct {
    //{{1
	BarCode string `json:"barCode"`
	SkuId string `json:"skuId"`
	Name string `json:"name"`
	MainType string `json:"mainType"`
	SubType string `json:"subType"`
	Unit string `json:"unit"`
	EnterPrice string `json:"enterPrice"`
	SellPrice string `json:"sellPrice"`
	GoodsId int
	//}}1

	enterPrice int
	sellPrice int
}

func (p *CfgGoodsData)GetSellPrice() int {
	return p.sellPrice
}


var cfgGoodsDatasMap map[string]*CfgGoodsData
//var goodsIdcfgGoodsDatasdMap map[int]*CfgGoodsData

func GetCfgGoodsDataByBarCode(barCode string)*CfgGoodsData {
	if len(barCode) == 0 {
		return nil
	}
    return cfgGoodsDatasMap[barCode]
}

func GetCfgGoodsDatas() map[string]*CfgGoodsData {
    return cfgGoodsDatasMap
}
//func GetCfgGoodsDataByGoodsId(goodsId int)*CfgGoodsData {
//	return goodsIdcfgGoodsDatasdMap[goodsId]
//}

func InitGoodsData()  {
	cfgs := *LoadJson("cfgn_goodsdatas")
	if cfgs == nil {
		fmt.Println("json cfgn_goodsdatas is lost")
		return
	}
	var cfgGoodsDatasSkuIdMap = map[string]*CfgGoodsData{}
	var _cfgGoodsDatasMap = make(map[string]*CfgGoodsData)
	//var _goodsIdcfgGoodsDatasdMap = make(map[int]*CfgGoodsData)
	for _, cfg := range cfgs {
		item := &CfgGoodsData{}
        //{{2
		item.BarCode = cfg["barCode"].(string)
		item.SkuId = cfg["skuId"].(string)
		item.Name = cfg["name"].(string)
		item.MainType = cfg["mainType"].(string)
		item.SubType = cfg["subType"].(string)
		item.Unit = cfg["unit"].(string)
		item.EnterPrice = cfg["enterPrice"].(string)
		item.SellPrice = cfg["sellPrice"].(string)
		//}}2

		item.enterPrice = int(common.AtoF(item.EnterPrice)*10)*10
		item.sellPrice = int(common.AtoF(item.SellPrice)*10)*10

		_cfgGoodsDatasMap[item.BarCode] = item

		goodsIdStr := item.SkuId[3:]
		goodsId, err := strconv.Atoi(goodsIdStr)
		if err != nil {
			panic("goodsId errer err")
		}
		item.GoodsId = goodsId+1000000
		tmp := cfgGoodsDatasSkuIdMap[item.SkuId]
		if tmp != nil || len(item.SkuId) == 0 {
			panic("repeat item")
		}
		cfgGoodsDatasSkuIdMap[item.SkuId] = item
		//_goodsIdcfgGoodsDatasdMap[item.GoodsId] = item
	}
	cfgGoodsDatasMap = _cfgGoodsDatasMap
	//goodsIdcfgGoodsDatasdMap = _goodsIdcfgGoodsDatasdMap
}


