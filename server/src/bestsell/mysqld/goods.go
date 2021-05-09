package mysqld

import (
	"bestsell/common"
	"bestsell/config"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
)


type EGoodsStatus int
const (
	EGoodsStatusAll    	EGoodsStatus = -1
	EGoodsStatusEditing EGoodsStatus = 0
	EGoodsStatusEdited  EGoodsStatus = 1
	EGoodsStatusUp  	EGoodsStatus = 2
	EGoodsStatusDown  	EGoodsStatus = 3
	//GoodsStatusIdle  	GoodsStatus = 4
)

type EGoodsMark int
const (
	EGoodsMarkNone    	EGoodsMark = -1
	EGoodsMarkNormal  	EGoodsMark = 0
	EGoodsMarkRecomm  	EGoodsMark = 1
	EGoodsMarkSeckill	EGoodsMark = 2
	EGoodsMarkTeam  	EGoodsMark = 3
)

type DBGoods struct {
    MysqlModel
	GoodsId  		int `gorm:"not null;unique"`
	ShopId  		int
	Status  		int  //0开始编辑，1完成编辑，2上架，3下架，4废弃
	Mark 		 	int  //
	NumberFav  		int
	NumberRepute  	int
	NumberOrder  	int
	NumberSell  	int
	NumberView  	int
	NumberStore  	int
	BuyLimit 		int
	StartTime       int64
	EndTime       	int64
	EditLock 		int
	goodsInfo 		*DBGoodsInfo
}

var dbGoodsSafeMap common.SafeMap

var allGoodsSlice []*DBGoods
var goodsSliceMap map[int][]*DBGoods
var subTypesGoodsMap map[int]map[string]*[]*DBGoods

var lockMutex sync.RWMutex

func GetAllGoodsSlice() []*DBGoods {
	lockMutex.Lock()
	defer lockMutex.Unlock()

	return allGoodsSlice[:]
}

func GetGoodsSliceByByCategoryId(categoryId int) []*DBGoods {
	lockMutex.Lock()
	defer lockMutex.Unlock()

	ret, ok := goodsSliceMap[categoryId]
	if !ok {
		return []*DBGoods{}
	}
	return ret[:]
}

func GetSubGoodsGroupByByCategoryId(categoryId int) *map[string]*[]*DBGoods {
	lockMutex.Lock()
	defer lockMutex.Unlock()

	ret, ok := subTypesGoodsMap[categoryId]
	if !ok {
		return nil
	}
	typeGoodsGroup := map[string]*[]*DBGoods{}
	for key, value := range ret {
		goodsList := (*value)[:]
		typeGoodsGroup[key] = &goodsList
	}
	return &typeGoodsGroup
}

type CustomSort struct {
	t []*DBGoods
	less func(x *DBGoods, y *DBGoods) bool
}

func (x CustomSort) Len() int {
	return len(x.t)
}
func (x CustomSort) Less(i, j int) bool {
	return x.less(x.t[i], x.t[j])
}

func (x CustomSort) Swap(i, j int) {
	x.t[i], x.t[j] = x.t[j], x.t[i]
}
func SortGoodsSlice(goodsSlice *[]*DBGoods)  {
	sort.Sort(CustomSort{
		*goodsSlice,
		func(a *DBGoods, b *DBGoods) bool {
			return a.GetGoodsInfo().Order < b.GetGoodsInfo().Order
		},
	})
}

func LoadGoodsCaches() {
	lockMutex.Lock()
	defer lockMutex.Unlock()

	allGoodsSlice = GetDBGoodsList(EGoodsStatusUp)
	SortGoodsSlice(&allGoodsSlice)
	goodsSliceMap = make(map[int][]*DBGoods)
	subTypesGoodsMap = make(map[int]map[string]*[]*DBGoods)

	dbCategoryList := GetDBCategoryList()
	for _, dbCategory := range dbCategoryList {
		goodsList := GetDBGoodsListByCategoryId(dbCategory.ID, EGoodsStatusUp)
		SortGoodsSlice(&goodsList)
		goodsSliceMap[dbCategory.ID] = goodsList
		typeGoodsGroup := map[string]*[]*DBGoods{}
		for _, dbGoods := range goodsList {
			goodsInfo := dbGoods.GetGoodsInfo()
			typeGoodsList, ok := typeGoodsGroup[goodsInfo.SubType]
			if !ok {
				typeGoodsGroup[goodsInfo.SubType] = &[]*DBGoods{}
				typeGoodsList, ok = typeGoodsGroup[goodsInfo.SubType]
				if !ok {
					panic("LoadGoodsCaches typeGoodsList, ok = typeGoodsGroup[goodsInfo.subType] goodsInfo.SubType ="+goodsInfo.SubType)
				}
			}
			*typeGoodsList = append(*typeGoodsList, dbGoods)
		}
		for _, typeGoodsList := range typeGoodsGroup {
			SortGoodsSlice(typeGoodsList)
		}
		subTypesGoodsMap[dbCategory.ID] = typeGoodsGroup
	}
}

func init() {
}

func GetDBGoods(goodsId int)*DBGoods {
	ret := dbGoodsSafeMap.Get(goodsId)
	if ret == nil {
		return nil
	}
	return ret.(*DBGoods)
}

func AddDBGoods(item *DBGoods) {
	// fmt.Println("AddDBGoods  goodsId =", item.GoodsId, "ID =", item.ID)
	old := GetDBGoods(item.GoodsId)
	if old != nil {
		panic("AddDBGoods DBGoods repeat")
	}
	dbGoodsSafeMap.Set(item.GoodsId, item)
}

func RemoveDBGoods(goodsId int) {
	dbGoodsSafeMap.Remove(goodsId)
}

func LoadDBGoodss() {
	var _dbItemsSlice []*DBGoods
	db.Find(&_dbItemsSlice)
	for _,item := range _dbItemsSlice{
		if item.GoodsId > 0 {
			item.LoadGoodsInfo()

			AddDBGoods(item)
		}
	}
}

func startDBGoods() {
	if !db.HasTable(&DBGoods{}) {
		db.CreateTable(&DBGoods{})
	}
	dbGoodsSafeMap = *common.NewSafeMap()
	LoadDBGoodss()
	LoadGoodsCaches()
}

func GetDBGoodsList(status EGoodsStatus)[]*DBGoods {
	var items []*DBGoods
	iterFunc := func(key int, v interface{}) bool {
		item := v.(*DBGoods)
		if status == -1 || item.GetStatus() == status {
			items = append(items, item)
		}
		return true
	}
	dbGoodsSafeMap.RangeSafe(iterFunc)
	return items
}

//func MakeNewGoodsId() int {
//	newGoodsId := 10000000
//	iterFunc := func(key int, v interface{}) bool{
//		item := v.(*DBGoods)
//		if item.GoodsId >= 10000000 {
//			if item.GoodsId >= newGoodsId {
//				newGoodsId = item.GoodsId+1
//			}
//		}
//		return true
//	}
//	dbGoodsSafeMap.RangeSafe(iterFunc)
//	return newGoodsId
//}

func GetDBGoodsListByCategoryId(categoryId int, status EGoodsStatus)[]*DBGoods {
	var items []*DBGoods
	iterFunc := func(key int, v interface{}) bool{
		item := v.(*DBGoods)
		if status == -1 || item.GetStatus() == status {
			dbGoodsInfo := item.GetGoodsInfo()
			if dbGoodsInfo.CategoryId == categoryId {
				items = append(items, item)
			}
		}
		return true
	}
	dbGoodsSafeMap.RangeSafe(iterFunc)
	return items
}

func GetDBGoodsListByStatus(status EGoodsStatus)[]*DBGoods{
	var items []*DBGoods
	iterFunc := func(key int, v interface{}) bool{
		item := v.(*DBGoods)
		if item.GetStatus() == status {
			items = append(items, item)
		}
		return true
	}
	dbGoodsSafeMap.RangeSafe(iterFunc)
	return items
}

func GetDBGoodsOrFromDB(goodsId int)*DBGoods {
	item := GetDBGoods(goodsId)
	if item != nil {
		return item
	}
	item = &DBGoods{
		GoodsId: goodsId,
	}
	item.LoadWithGoodsId()
	if item.ID <= 0 {
		return nil
	}
	AddDBGoods(item)
	return item
}

func RemoveDBGoodsAndDB(goodsId int)  {
	dbGoods := GetDBGoodsOrFromDB(goodsId)
	if dbGoods == nil {
		return
	}
	RemoveDBGoods(goodsId)
	dbGoods.Remove()

	dbGoods.LoadGoodsInfo()
	if dbGoods.goodsInfo == nil {
		return
	}
	dbGoods.goodsInfo.Remove()
	dbGoods.goodsInfo = nil
}

func AddNewDBGoodsByInfo(goodsId int, name string)*DBGoods {
	goodsInfo := GetDBGoodsInfo(goodsId)
	if goodsInfo == nil {
		goodsInfo = &DBGoodsInfo{
			GoodsId: goodsId,
			Name: name,
		}
		AddNewGoodsInfo(goodsInfo)
	}
	dbGoods := GetDBGoodsOrFromDB(goodsInfo.GoodsId)
	if dbGoods != nil {
		if dbGoods.goodsInfo != nil && dbGoods.goodsInfo.GoodsId != goodsInfo.GoodsId {
			panic("AddNewDBGoodsByInfo  goodsInfo no equal")
		}
		dbGoods.goodsInfo = goodsInfo
		return dbGoods
	}
	dbGoods = &DBGoods {
		GoodsId: goodsInfo.GoodsId,
	}
	dbGoods.LoadWithGoodsId()
	if dbGoods.ID <= 0 {
		dbGoods.Insert()
	}
	AddDBGoods(dbGoods)
	dbGoods.goodsInfo = goodsInfo
	return dbGoods
}

func GetDBGoodsListGroup() *map[int]*[]*DBGoods {
	goodsListGroup := map[int]*[]*DBGoods{}
	teamDBGoodsList := []*DBGoods{}
	recomDBGoodsList := []*DBGoods{}
	seckillDBGoodsList := []*DBGoods{}
	iterFunc := func(key int, v interface{}) bool {
		item := v.(*DBGoods)
		if item.GetStatus() == EGoodsStatusUp {
			switch EGoodsMark(item.Mark) {
			case EGoodsMarkRecomm:
				recomDBGoodsList = append(recomDBGoodsList, item)
			case EGoodsMarkSeckill:
				seckillDBGoodsList = append(seckillDBGoodsList, item)
			case EGoodsMarkTeam:
				teamDBGoodsList = append(teamDBGoodsList, item)
			}
		}
		return true
	}
	dbGoodsSafeMap.RangeSafe(iterFunc)

	goodsListGroup[int(EGoodsMarkRecomm)] = &recomDBGoodsList
	goodsListGroup[int(EGoodsMarkSeckill)] = &seckillDBGoodsList
	goodsListGroup[int(EGoodsMarkTeam)] = &teamDBGoodsList
	return &goodsListGroup
}

func (p *DBGoods)GetTag() int {
	return p.Mark
}

func (p *DBGoods)GetStatus() EGoodsStatus {
	return EGoodsStatus(p.Status)
}

func (p *DBGoods)SetStatus(status EGoodsStatus) {
	p.BeginWrite()
	defer p.EndWrite()
	p.Status = int(status)
}

func (p *DBGoods)GetShowNumberSell() int {
	dbGoodsInfo := p.GetGoodsInfo()
	return dbGoodsInfo.NumberSell + p.NumberSell
}

func (p *DBGoods)AddSellNum(numberSell int) {
	p.BeginWrite()
	defer p.EndWrite()

	p.NumberSell = p.NumberSell + numberSell
}

func (p *DBGoods)AddOrderNum() {
	p.BeginWrite()
	defer p.EndWrite()

	p.NumberOrder++
}

//DBGoods
func (p *DBGoods)GetGoodsInfo() *DBGoodsInfo {
	if p.goodsInfo == nil {
		p.LoadGoodsInfo()
	}
	return p.goodsInfo
}

func (p *DBGoods)LoadGoodsInfo() {
	p.BeginWrite()
	defer p.EndWrite()

	if p.goodsInfo == nil {
		p.goodsInfo = GetDBGoodsInfoOrFromDB(p.GoodsId)
		if p.goodsInfo != nil {
			if p.GoodsId != p.goodsInfo.GoodsId {
				panic("(p *DBGoods)LoadGoodsInfo1 p.GoodsId != p.goodsInfo.GoodsId")
			}
			return
		}
		p.goodsInfo = &DBGoodsInfo{
			GoodsId:p.GoodsId,
		}
	} else {
		if p.GoodsId != p.goodsInfo.GoodsId {
			panic("(p *DBGoods)LoadGoodsInfo2 p.GoodsId != p.goodsInfo.GoodsId")
		}
	}
	p.goodsInfo.LoadWithGoodsId()
	if p.goodsInfo.ID < 0 {
		p.goodsInfo.InsertWithGoodsId()
		p.goodsInfo.LoadWithGoodsId()
	}
}

func (p *DBGoods)Insert() {
	if p.GoodsId <= 0 {
		panic("(p *DBGoods)Insert p.GoodsId <= 0")
	}
	db.Create(p)
}

func (p *DBGoods)LoadWithGoodsId() {
	if p.GoodsId <= 0 {
		panic("(p *DBGoods)LoadWithGoodsId p.GoodsId < 0")
	}
	db.Where("goods_id = ?", p.GoodsId).First(p)
}

func (p *DBGoods)Save() {
	if p.ID < 0 {
		panic("(p *DBGoods)Save p.ID < 0")
	}
	if p.GoodsId <= 0 {
		panic("(p *DBGoods)Save p.GoodsId <= 0")
	}
	// fmt.Println("(p *DBGoods)Save()  p.ID =", p.ID)
	// fmt.Println("(p *DBGoods)Save()  p.GoodsId =", p.GoodsId)
	db.Save(p)
}

func (p *DBGoods)Remove() {
	//if p.ID < 0 {
	//	panic("(p *DBGoods)Remove p.ID < 0")
	//}
	//db.Delete(p)
}

func (p *DBGoods)GetInfo() *map[string]interface{} {
	dbGoodsInfo := p.GetGoodsInfo()
	itemInfo := &map[string]interface{}{
		//"uid": p.ID,
		"id": p.GoodsId,
		"name": dbGoodsInfo.Name,
		"mark": p.Mark,
		"order": dbGoodsInfo.Order,
		"isPromote": len(dbGoodsInfo.Promote)>0,
		"sellPrice": dbGoodsInfo.SellPrice,
		"originPrice": dbGoodsInfo.OriginPrice,
		// "minPrice": dbGoodsInfo.MinPrice,
		"numberStore": p.NumberStore,
		"numberSell": p.GetShowNumberSell(),
		"updateTime": p.CreatedAt.Unix(),
		"subType": dbGoodsInfo.SubType,
		"subType2": dbGoodsInfo.SubType2,
		"pic": dbGoodsInfo.GetPublicIcon(),
	}
	return itemInfo
}

func (p *DBGoods)GetDetailInfo() *map[string]interface{} {
	dbGoodsInfo := p.GetGoodsInfo()
	goodsDetail := &map[string]interface{}{
		"id": dbGoodsInfo.GoodsId,
		"name": dbGoodsInfo.Name,
		"mark": p.Mark,
		"promote": dbGoodsInfo.Promote,
		"numberSell": p.GetShowNumberSell(),
		"numberStore": p.NumberStore,
		"startTime": p.StartTime,
		"endTime": p.EndTime,
		"updateTime": p.CreatedAt.Unix(),
		"sellPrice": dbGoodsInfo.SellPrice,
		"originPrice": dbGoodsInfo.OriginPrice,
		"skuJson" :dbGoodsInfo.SkuStruct,
		"skuPics":dbGoodsInfo.GetSkuPics(),
		// "minPrice": dbGoodsInfo.MinPrice,
		//"skus": dbGoodsInfo.GetSkuInfos(),
		"icon": dbGoodsInfo.GetPublicIcon(),
		"pics": dbGoodsInfo.GetPublicPics(),
		"contents": dbGoodsInfo.GetPublicContents(),
	}
	return goodsDetail
}

func (p *DBGoods)GetEditInfo() *map[string]interface{} {
	if p.GoodsId == 0 {
		info := & map[string]interface{} {
			"goodsId" 		:0,
			"categoryId" 	:0,
			"name" 			:"",
			"barCode" 		:"",

			"originPrice" 	:0,
			"sellPrice" 	:0,
			// "minPrice"		: 0,

			"promote" 		:"",
			"order" 		:99999999,

			"mainType" 		:"",
			"subType" 		:"",
			"numberSell" 	:0,

			"skuJson" 	:"",

			"numberStore" 	:99999,
		}
		return info
	}
	dbGoodsInfo := p.GetGoodsInfo()
	info := &map[string]interface{}{
		"id"			:p.ID,
		"goodsId" 		:dbGoodsInfo.GoodsId,
		"categoryId" 	:dbGoodsInfo.CategoryId,
		"name" 			:dbGoodsInfo.Name,
		"barCode" 		:dbGoodsInfo.BarCode,

		"originPrice" 	:dbGoodsInfo.OriginPrice,
		"sellPrice" 	:dbGoodsInfo.SellPrice,
		// "minPrice"		:dbGoodsInfo.MinPrice,

		"promote" 		:dbGoodsInfo.Promote,
		"order" 		:dbGoodsInfo.Order,

		"mainType" 		:dbGoodsInfo.MainType,
		"subType" 		:dbGoodsInfo.SubType,
		"numberSell" 	:dbGoodsInfo.NumberSell,
		"skuJson" 		:dbGoodsInfo.SkuStruct,

		"status"   		:int(p.GetStatus()),
		"numberStore" 	:p.NumberStore,
		"mark"  		:p.Mark,

		//"icon" 			:dbGoodsInfo.GetIcon(),
		"pics" 			:dbGoodsInfo.GetPics(),
		"contents" 		:dbGoodsInfo.GetContents(),

		//"publicIcon" 	:dbGoodsInfo.GetPublicIcon(),
		"publicPics" 	:dbGoodsInfo.GetPublicPics(),
		"publicContents" :dbGoodsInfo.GetPublicContents(),
	}
	return info
}

func ReloadGoods(goodsId int) {
	RemoveDBGoods(goodsId)
	RemoveDBGoodsInfoAndDB(goodsId)
	dbGoods := GetDBGoodsOrFromDB(goodsId)
	dbGoods.GetGoodsInfo()
	LoadGoodsCaches()
}

func loadConfig(_cfg *map[string]interface{}) {
	cfg := *_cfg
	goodsId := int(cfg["id"].(float64))
	if goodsId <= 0 {
		panic("loadConfig goodsId is illegal")
	}

	order := int(cfg["order"].(float64))
	if order <= 0 {
		order = 99999999
	}

	name := cfg["name"].(string)
	if len(name) == 0 {
		panic("loadConfig name is illegal")
	}

	desc := ""
	if _, ok := cfg["desc"]; ok {
		desc = cfg["desc"].(string)
	}

	resId := ""
	if _, ok := cfg["image"]; ok {
		resId = cfg["image"].(string)
	}

	promote := ""
	if _, ok := cfg["promote"]; ok {
		promote = cfg["promote"].(string)
	}

	mainType := ""
	if _, ok := cfg["mainType"]; ok {
		mainType = cfg["mainType"].(string)
	}
	subType := ""
	if _, ok := cfg["subType"]; ok {
		subType = cfg["subType"].(string)
	}
	if len(subType) == 0 {
		subType = "其他"
	}

	subType2 := ""
	if _, ok := cfg["subType2"]; ok {
		subType2 = cfg["subType2"].(string)
	}

	numberSell := 0
	if _, ok := cfg["numberSell"]; ok {
		numberSell = int(cfg["numberSell"].(float64))
	}

	store := int(cfg["store"].(float64))
	category := int(cfg["category"].(float64))

	status := 1
	if _, ok := cfg["status"]; ok {
		switch cfg["status"].(type) {
			case float64:
				status = int(cfg["status"].(float64))
			case string:
				statusStr := cfg["status"].(string)
				if statusStr == "停售" {
					status = 1
				} else if statusStr == "上架" {
					status = 2
				}
		}
	}
	
	enterPrice := 0
	if _, ok := cfg["enterPrice"]; ok {
		enterPrice = int(cfg["enterPrice"].(float64)*10)*10
	}
	// minPrice := 0.0
	// if _, ok := cfg["minPrice"]; ok {
	// 	minPrice = cfg["minPrice"].(float64)
	// }
	sellPrice := 0
	if _, ok := cfg["sellPrice"]; ok {
		sellPrice = int(cfg["sellPrice"].(float64)*10)*10
	}
	originPrice := 0
	if _, ok := cfg["originPrice"]; ok {
		originPrice = int(cfg["originPrice"].(float64)*10)*10
	}
	if originPrice == 0 {
		originPrice = sellPrice
	}
	if originPrice < enterPrice {
		panic("loadConfig originPrice < enterPrice resId"+resId)
	}

	//sekillPrice := 0.0
	//if _, ok := cfg["sekillPrice"]; ok {
	//	sekillPrice = cfg["sekillPrice"].(float64)
	//}
	//startTime := int64(0)
	//if _, ok := cfg["startTime"]; ok {
	//	startTimeStr := cfg["startTime"].(string)
	//	stamp, _ := time.Parse("2006/01/02 15:04:05", startTimeStr)
	//	startTime = stamp.Unix()
	//}
	//
	//endTime := int64(0)
	//if _, ok := cfg["endTime"]; ok {
	//	endTimeStr := cfg["endTime"].(string)
	//	stamp, _ := time.Parse("2006/01/02 15:04:05", endTimeStr)
	//	endTime = stamp.Unix()
	//}
	//if endTime > 0 {
	//	if startTime == 0 {
	//		startTime = time.Now().Unix()
	//	}
	//	if startTime > endTime {
	//		panic("loadConfig startTime > endTime name:"+name)
	//	}
	//}

	goodsSkus := []*DBGoodsSku{}

	skuLabel1 := "规格"
	if _, ok := cfg["skuLabel1"]; ok {
		skuLabel1 = cfg["skuLabel1"].(string)
	}
	for idx:=1; idx < 10; idx++ {
		//sku1, sku2, sku3, sku4
		key := "sku"+strconv.Itoa(idx)
		tmp, ok := cfg[key]
		if !ok {
			key = "sku1"+strconv.Itoa(idx)
			tmp, ok = cfg[key]
			if !ok {
				break
			}
		}
		skuName := tmp.(string)

		if len(skuName) == 0 {
			continue
		}

		skuPrice := 0
		key2 := "skuPrice"+strconv.Itoa(idx)
		tmp2, ok2 := cfg[key2]
		if ok2 {
			skuPrice = int(tmp2.(float64)*10)*10
		} else {
			key2 = "skuPrice1"+strconv.Itoa(idx)
			tmp2, ok2 = cfg[key2]
			if ok2 {
				skuPrice = int(tmp2.(float64)*10)*10
			}
		}

		goodsSku := &DBGoodsSku{
			Label: skuLabel1,
			Id: idx,
			Name: skuName,
			Price: skuPrice,
		}
		goodsSkus = append(goodsSkus, goodsSku)
	}

	for pidx := 2; pidx < 3; pidx++ {
		skuLabel := "规格"
		labelKey := "skuLabel"+strconv.Itoa(pidx)
		if _, ok := cfg[labelKey]; ok {
			skuLabel = cfg[labelKey].(string)
		}

		pkey := "sku"+strconv.Itoa(pidx)
		key := pkey+"1"
		if _, ok := cfg[key]; !ok {
			break
		}
		for idx := 1; idx < 10; idx++ {
			//sku21, sku22, sku23, sku24
			key := pkey+strconv.Itoa(idx)
			tmp, ok := cfg[key]
			if !ok {
				break
			}
			skuName := tmp.(string)

			if len(skuName) == 0 {
				continue
			}

			skuPrice := 0
			key2 := "skuPrice"+strconv.Itoa(pidx)+strconv.Itoa(idx)
			tmp2, ok2 := cfg[key2]
			if ok2 {
				skuPrice = int(tmp2.(float64)*10)*10
			}
			goodsSku := &DBGoodsSku{
				Label: skuLabel,
				Id: idx*100,
				Name: skuName,
				Price: skuPrice,
			}
			goodsSkus = append(goodsSkus, goodsSku)
		}
	}
	skuStruct := ""
	if len(goodsSkus) > 0 {
		jsonStr, err := json.Marshal(goodsSkus)
		if err != nil {
			fmt.Println(err)
			panic("error")
		}
		skuStruct = string(jsonStr)
	}

	goods := AddNewDBGoodsByInfo(goodsId, name)
	if goods.EditLock > 0 {
		fmt.Println("EditLock goodsId =", goodsId)
		return
	}

	goodsInfo := goods.GetGoodsInfo()
	//enterPrice = common.MakeMoneyValue(enterPrice)
	//sellPrice = common.MakeMoneyValue(sellPrice)
	//originPrice = common.MakeMoneyValue(originPrice)
	// minPrice = common.MakeMoneyValue(minPrice)
	for  {
		if goodsInfo.Order != order { break }
		if goodsInfo.CategoryId != category { break }
		if goodsInfo.Name != name { break }
		if goodsInfo.Desc != desc { break }
		if goodsInfo.ResId != resId { break }
		if goodsInfo.Promote != promote { break }
		if goodsInfo.EnterPrice != enterPrice { break }
		if goodsInfo.SellPrice != sellPrice { break }
		if goodsInfo.OriginPrice != originPrice { break }
		// if goodsInfo.MinPrice != minPrice { break }
		if goodsInfo.MainType != mainType { break }
		if goodsInfo.SubType != subType { break }
		if goodsInfo.SubType2 != subType2 { break }
		if goodsInfo.NumberSell != numberSell { break }
		if goodsInfo.SkuStruct != skuStruct { break }

		if goods.NumberStore != store { break }
		if goods.Status != status { break }

		fmt.Println("no change goodsId =", goodsId)
		return
	}

	goodsInfo.Order = order
	goodsInfo.CategoryId = category
	goodsInfo.Name = name
	goodsInfo.Desc = desc
	goodsInfo.ResId = resId
	goodsInfo.Promote = promote
	goodsInfo.EnterPrice = enterPrice
	goodsInfo.SellPrice = sellPrice
	goodsInfo.OriginPrice = originPrice
	// goodsInfo.MinPrice = minPrice
	goodsInfo.MainType = mainType
	goodsInfo.SubType = subType
	goodsInfo.SubType2 = subType2
	if numberSell > 0 {
		goodsInfo.NumberSell = numberSell
	}
	goodsInfo.SkuStruct = skuStruct
	//goods.StartTime = startTime
	//goods.EndTime = endTime
	//if endTime > 0 && sekillPrice > 0.0{
		//goods.Mark = int(GoodsMarkSeckill)
		//goodsInfo.SellPrice = sekillPrice
	//} else {
		//goods.Mark = int(GoodsMarkNormal)
	//}
	goodsInfo.TryLoadImage()
	goodsInfo.Save()

	goods.NumberStore = store
	goods.SetStatus(EGoodsStatus(status))
	goods.Save()
}

func LoadConfigs() {
	cfgs := *config.LoadJson("cfg_goods")
	for _,cfg := range cfgs {
		loadConfig(&cfg)
	}
}
