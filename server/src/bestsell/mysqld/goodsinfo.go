package mysqld

import (
	"bestsell/common"
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"
)

type DBGoodsSku struct {
	Label string `json:"label"`
	Id int `json:"id"`
	Name string `json:"name"`
	Price int `json:"price"`
	Icon string `json:"icon"`
	publicIcon string
}

var Module_GetFilePublicPath func(string)string

type DBGoodsInfo struct {
    MysqlModel
	GoodsId  		int `gorm:"not null;unique"`
	PlayerId  		int
	CategoryId  	int
	Name   			string `gorm:"type:varchar(256);not null;unique"`
    EnterPrice   	int
    OriginPrice  	int
	SellPrice  	 	int
	MinPrice  		int
	VipPrice  	 	int
	BarCode  		string `gorm:"type:varchar(64);"`
	EnterDate  		string `gorm:"type:varchar(32);"`
	ProductDate  	string `gorm:"type:varchar(32);"`
	Desc   			string `gorm:"type:varchar(256);"`
	Supplier   		string `gorm:"type:varchar(64);"`
	Brand  			string `gorm:"type:varchar(64);"`
	Promote   		string `gorm:"type:text;"`
	ResId  			string `gorm:"type:varchar(128);"`
	Order 			int
	Express 		int
	Icon  			string `gorm:"type:varchar(256);"`
	Pics   			string `gorm:"type:text;"`
	Contents  		string `gorm:"type:text;"`

	MainType  		string `gorm:"type:varchar(32);"`
	SubType  		string `gorm:"type:varchar(32);"`
	SubType2  		string `gorm:"type:varchar(32);"`
	NumberSell  	int

	SkuStruct  		string `gorm:"type:text;"`
	goodsSkus  		[]*DBGoodsSku
    skuLevel 		int

	icon  			string
	pics   			string
	contents  		string

	publicIcon  	string
	publicPics   	[]string
	publicContents  []string

	skuPics   	map[int]string
}

var dbGoodsInfoSafeMap common.SafeMap

func GetDBGoodsInfo(goodsId int)*DBGoodsInfo {
	ret := dbGoodsInfoSafeMap.Get(goodsId)
	if ret == nil {
		return nil
	}
	return ret.(*DBGoodsInfo)
}

func AddDBGoodsInfo(item *DBGoodsInfo)  {
	if item.ID < 0 {
		panic("AddDBGoodsInfo item.ID < 0")
	}
	if item.GoodsId <= 0 {
		panic("AddDBGoodsInfo item.GoodsId <= 0")
	}
	old := GetDBGoodsInfo(item.GoodsId)
	if old != nil {
		panic("AddDBGoodsInfo repeat")
		//return
	}
	dbGoodsInfoSafeMap.Set(item.GoodsId, item)
}

func RemoveDBGoodsInfo(goodsId int)  {
	dbGoodsInfoSafeMap.Remove(goodsId)
}

func LoadDBGoodsInfos()  {
	var _dbItemsSlice []*DBGoodsInfo
	db.Find(&_dbItemsSlice)
	dbGoodsInfoSafeMap = *common.NewSafeMap()
	for _,item := range _dbItemsSlice {
		AddDBGoodsInfo(item)
		if item.ID > 0 {
			item.TryLoadImage()
			item.GetGoodsSkus()
		}
	}
	println("LoadDBGoodsInfos success")
}

func GetDBGoodsInfoOrFromDB(goodsId int)*DBGoodsInfo {
	item := GetDBGoodsInfo(goodsId)
	if item != nil {
		return item
	}
	item = &DBGoodsInfo{
		GoodsId:goodsId,
	}
	item.LoadWithGoodsId()
	if item.ID == 0 {
		return nil
	}
	AddDBGoodsInfo(item)
	return item
}

func RemoveDBGoodsInfoAndDB(goodsId int) {
	dbGoodsInfo := GetDBGoodsInfoOrFromDB(goodsId)
	if dbGoodsInfo == nil {
		return
	}
	RemoveDBGoodsInfo(goodsId)
	dbGoodsInfo.Remove()
}

func startDBGoodsInfo() {
	if !db.HasTable(&DBGoodsInfo{}) {
		db.CreateTable(&DBGoodsInfo{})
	}
	//orderSkuId := 202020
	//splitOrderSKuId(orderSkuId)
	LoadDBGoodsInfos()
}

func AddNewGoodsInfo(item *DBGoodsInfo)bool {
	if item.GoodsId <= 0 {
		panic("AddNewGoodsInfo p.GoodsId < 0")
	}
	item.LoadWithGoodsId()
	if item.ID > 0 {
		return true
	}
	item.InsertWithGoodsId()
	if item.ID < 0 {
		return false
	}
	return true
}

func (p *DBGoodsInfo)SetSellPrice(sellPrice int) {
	p.BeginWrite()
	p.SellPrice = sellPrice
	p.EndWrite()
}

func (p *DBGoodsInfo)GetPublicPics()[]string {
	if len(p.pics) == 0 {
		return []string{}
	}
	if len(p.publicPics) > 0 {
		return p.publicPics
	}
	publicPics := addPublicUrlHosts(p.pics)
	p.BeginWrite()
	p.publicPics = publicPics
	p.EndWrite()
	return p.publicPics
}

func (p *DBGoodsInfo)GetPics()[]string {
	if len(p.pics) == 0 {
		return []string{}
	}
	return addUrlHosts(p.pics)
}

func (p *DBGoodsInfo)GetPublicContents()[]string {
	if len(p.contents) == 0 {
		return []string{}
	}
	if len(p.publicContents) > 0 {
		return p.publicContents
	}
	publicContents := addPublicUrlHosts(p.contents)
	p.BeginWrite()
	p.publicContents = publicContents
	p.EndWrite()
	return p.publicContents
}

func (p *DBGoodsInfo)GetContents()[]string {
	if len(p.contents) == 0 {
		return []string{}
	}
	return addUrlHosts(p.contents)
}

func (p *DBGoodsInfo)GetPublicIcon()string {
	publicPics := p.GetPublicPics()
	if len(publicPics) > 0 {
		return publicPics[0]
	}
	if len(p.icon) == 0 {
		return ""
	}
	if len(p.publicIcon) > 0 {
		return p.publicIcon
	}
	publicIcon := addPublicUrlHost(p.icon)
	p.BeginWrite()
	p.publicIcon = publicIcon
	p.EndWrite()
	return p.publicIcon
}

func (p *DBGoodsInfo)GetIcon()string {
	if len(p.icon) == 0 {
		return ""
	}
	return addUrlHost(p.icon)
}

func (p *DBGoodsInfo)GetOrderPublicIcon(orderSkuId int) string {
	if orderSkuId == 0 {
		return p.GetPublicIcon()
	}
	skuId1 := orderSkuId % 100
	skuIcon, ok := p.skuPics[skuId1]
	if !ok || len(skuIcon) == 0 {
		return p.GetPublicIcon()
	}
	return skuIcon
}

func (p *DBGoodsInfo)GetSkuPics() map[int]string {
	return p.skuPics
}

func (p *DBGoodsInfo)GetSkuLevel()int {
	if p.skuLevel > 0 {
		return p.skuLevel
	}
	if len(p.SkuStruct) > 0 {
		p.GetGoodsSkus()
	}
	return p.skuLevel
}

func (p *DBGoodsInfo)SetSkuStruct(skuJson string) bool {
	goodsSkus := []*DBGoodsSku{}
	err := json.Unmarshal([]byte(skuJson), &goodsSkus)
	if err != nil {
		fmt.Println("(p *DBGoodsInfo)SetGoodsSkus() err:", err)
		return false
	}
	p.BeginWrite()
	p.SkuStruct = skuJson
	p.skuLevel = 0
	p.goodsSkus = []*DBGoodsSku{}
	p.EndWrite()

	p.GetGoodsSkus()
	return true
}

func (p *DBGoodsInfo)GetGoodsSkus()[]*DBGoodsSku {
	p.BeginWrite()
	defer p.EndWrite()

	if len(p.SkuStruct) == 0 {
		if len(p.goodsSkus) > 0 {
			p.skuLevel = 0
			p.goodsSkus = []*DBGoodsSku{}
		}
		return p.goodsSkus
	}
	if len(p.goodsSkus) > 0 {
		if len(p.SkuStruct) > 0 {
			return p.goodsSkus
		}
		p.skuLevel = 0
		p.goodsSkus = []*DBGoodsSku{}
	}
	if len(p.SkuStruct) == 0 {
		return p.goodsSkus
	}
	goodsSkus := []*DBGoodsSku{}
	err := json.Unmarshal([]byte(p.SkuStruct), &goodsSkus)
	if err != nil {
		fmt.Println("(p *DBGoodsInfo)GetGoodsSkus() err:", err)
		return []*DBGoodsSku{}
	}
	p.skuLevel = 0
	for _,goodsSku := range goodsSkus {
		if goodsSku.Id < 100 {
			if p.skuLevel == 0 {
				p.skuLevel = 1
			}
		} else if goodsSku.Id < 10000 {
			if p.skuLevel < 2 {
				p.skuLevel = 2
			}
		} else if goodsSku.Id < 1000000 {
			if p.skuLevel < 3 {
				p.skuLevel = 3
			}
		} else {
			panic("(p *DBGoodsInfo)GetGoodsSkus too much sku level")
		}
	}
	p.goodsSkus = goodsSkus
	return p.goodsSkus
}

func splitOrderSKuId(orderSkuId int) (int, int, int) {
	if orderSkuId <= 0 {
		return 0, 0, 0
	}
	skuId1 := orderSkuId % 100
	orderSkuId = orderSkuId/100
	if orderSkuId <= 0 {
		return skuId1, 0, 0
	}
	skuId2 := (orderSkuId % 100)*100
	orderSkuId = orderSkuId/100
	if orderSkuId <= 0 {
		return skuId1, skuId2, 0
	}
	skuId3 := (orderSkuId % 100)*10000
	return skuId1, skuId2, skuId3
}

func (p *DBGoodsInfo)IsValidOrderSkuId(orderSkuId int) bool {
	if orderSkuId == 0 {
		if p.skuLevel == 0 {
			return true
		}
		return false
	}
	goodsSkus := p.GetGoodsSkus()
	if p.skuLevel == 1 {
		skuId1, skuId2, skuId3 := splitOrderSKuId(orderSkuId)
		if skuId2 > 0 || skuId3 > 0 {
			return false
		}
		for _, goodsSku := range goodsSkus {
			if goodsSku.Id == skuId1 {
				return true
			}
		}
		return false
	}
	if p.skuLevel == 2 {
		skuId1, skuId2, skuId3 := splitOrderSKuId(orderSkuId)
		if skuId3 > 0 || skuId1 == 0 || skuId2 == 0 {
			return false
		}
		for _, goodsSku := range goodsSkus {
			if goodsSku.Id == skuId1 {
				skuId1 = 0
			} else if goodsSku.Id == skuId2 {
				skuId2 = 0
			}
		}
		if skuId1 == 0 && skuId2 == 0 {
			return true
		}
		return false
	}
	skuId1, skuId2, skuId3 := splitOrderSKuId(orderSkuId)
	if skuId3 == 0 || skuId1 == 0 || skuId2 == 0 {
		return false
	}
	for _, goodsSku := range goodsSkus {
		if goodsSku.Id == skuId1 {
			skuId1 = 0
		} else if goodsSku.Id == skuId2 {
			skuId2 = 0
		} else if goodsSku.Id == skuId3 {
			skuId3 = 0
		}
	}
	if skuId1 == 0 && skuId2 == 0 && skuId3 == 0 {
		return true
	}
	return false
}

func (p *DBGoodsInfo)GetSkuPrice(orderSkuId int) int {
	goodsSkus := p.GetGoodsSkus()
	if len(goodsSkus) == 0{
		return 0
	}
	skuId1, skuId2, skuId3 := splitOrderSKuId(orderSkuId)
	for _,goodsSku := range goodsSkus{
		skuId := goodsSku.Id
		if skuId < 100 {
			if skuId1 == skuId {
				if goodsSku.Price > 0 {
					return goodsSku.Price
				}
			}
		} else if skuId2 > 0 && skuId < 10000 {
			if skuId2 == skuId {
				if goodsSku.Price > 0 {
					return goodsSku.Price
				}
			}
		} else {
			if skuId3 > 0 && skuId3 == skuId {
				if goodsSku.Price > 0 {
					return goodsSku.Price
				}
			}
		}
	}
	return 0
}

func (p *DBGoodsInfo)GetSkuNames(orderSkuId int) string {
	if orderSkuId == 0 {
		return ""
	}
	skuList := p.GetSkuList(orderSkuId)
	nameList := []string{}
	for _, goodsSku := range *skuList {
		nameList = append(nameList, goodsSku.Name)
	}
	skuNames := strings.Join(nameList, ";")
	return skuNames
}

func (p *DBGoodsInfo)GetRealPrice(orderSkuId int) int {
	if orderSkuId == 0 {
		return p.SellPrice
	}
	skuPrice := p.GetSkuPrice(orderSkuId)
	if skuPrice > 0 {
		return skuPrice
	}
	return p.SellPrice
}

func (p *DBGoodsInfo)GetSkuList(orderSkuId int) *[]*DBGoodsSku {
	skuList := []*DBGoodsSku{}
	if orderSkuId == 0 {
		return &skuList
	}

	skuId1, skuId2, skuId3 := splitOrderSKuId(orderSkuId)
	goodsSkus := p.GetGoodsSkus()
	for _,goodsSku := range goodsSkus {
		if skuId1 == goodsSku.Id || (skuId2 > 0 && skuId2 == goodsSku.Id) || (skuId3 > 0 && skuId3 == goodsSku.Id) {
			skuList = append(skuList, goodsSku)
		}
	}
	return &skuList
}

func (p *DBGoodsInfo)GetOrderSkuInfos(orderSkuId int) *[]map[string]interface{} {
	skuInfos := []map[string]interface{}{}
	if orderSkuId == 0 {
		return &skuInfos
	}

	goodsSkus := p.GetGoodsSkus()
	if len(goodsSkus) == 0 {
		return &skuInfos
	}

	skuId1, skuId2, skuId3 := splitOrderSKuId(orderSkuId)
	for _,goodsSku := range goodsSkus{
		if skuId1 == goodsSku.Id || (skuId2 > 0 && skuId2 == goodsSku.Id) || (skuId3 > 0 && skuId3 == goodsSku.Id) {
			skuInfo := map[string]interface{}{
				"id":goodsSku.Id,
				"name":goodsSku.Name,
				"price":goodsSku.Price,
			}
			skuInfos = append(skuInfos, skuInfo)
		}
	}
	return &skuInfos
}

func (p *DBGoodsInfo)SetNumberSell(numberSell int) {
	p.BeginWrite()
	p.NumberSell = numberSell
	p.EndWrite()
}

func (p *DBGoodsInfo)InsertWithGoodsId() {
	if p.GoodsId <= 0 {
		panic("(p *DBGoodsInfo)InsertWithGoodsId p.GoodsId <= 0")
	}
	if len(p.Name) == 0 {
		p.Name = "商品"+strconv.Itoa(p.GoodsId)
	}
	db.Create(p)
}

func (p *DBGoodsInfo)LoadWithGoodsId() {
	if p.GoodsId <= 0 {
		panic("(p *DBGoodsInfo)LoadWithGoodsId p.GoodsId <= 0")
	}
	db.Where("goods_id = ?", p.GoodsId).First(p)
	if p.ID > 0 {
		p.TryLoadImage()
		p.GetGoodsSkus()
	}
}

func (p *DBGoodsInfo)SaveWithGoodsId() {
	if p.ID < 0 {
		panic("(p *DBGoodsInfo)SaveWithGoodsId p.ID < 0")
	}
	if p.GoodsId <= 0 {
		panic("(p *DBGoodsInfo)SaveWithGoodsId p.GoodsId <= 0")
	}
	db.Save(p)
}

// func (p *DBGoodsInfo)Insert() {
// 	if p.ID < 0 {
// 		panic("(p *DBGoodsInfo)Insert p.ID < 0")
// 	}
// 	db.Create(p)
// }

// func (p *DBGoodsInfo)Load() {
// 	if p.ID < 0 {
// 		panic("(p *DBGoodsInfo)Load p.ID < 0")
// 	}
// 	db.First(p, p.ID)
// 	p.TryLoadImage()
// }

func (p *DBGoodsInfo)Save() {
	if p.ID < 0 {
		panic("(p *DBGoodsInfo)Save p.ID < 0")
	}
	if p.GoodsId <= 0 {
		panic("(p *DBGoodsInfo)Save p.GoodsId <= 0")
	}
	fmt.Println("(p *DBGoodsInfo)Save()  p.ID =", p.ID)
	fmt.Println("(p *DBGoodsInfo)Save()  p.GoodsId =", p.GoodsId)
	db.Save(p)
}

func (p *DBGoodsInfo)Remove() {
	//if p.ID < 0 {
	//	panic("(p *DBGoodsInfo)Remove p.ID < 0")
	//}
	//db.Delete(p)
}

func (p *DBGoodsInfo)TryLoadImage() bool {
	resId := p.ResId
	if len(resId) == 0 {
		if len(p.Icon) > 0 {
			p.icon = p.Icon
		}
		if len(p.Pics) > 0 {
			p.pics = p.Pics
		}
		if len(p.Contents) > 0 {
			p.contents = p.Contents
		}
		//if len(p.Icon) != 0 && len(p.Pics) != 0 && len(p.Contents) != 0 {
		return false
		//}
		//resId = "default"
	}
	return p.MakeGoodsFiles(resId)
}

func (p *DBGoodsInfo)ClearImage() {
	p.BeginWrite()
	p.Icon = ""
	p.Pics = ""
	p.Contents = ""

	p.icon = ""
	p.pics = ""
	p.contents = ""

	p.publicIcon = ""
	p.publicPics = p.publicPics[0:0]
	p.publicContents = p.publicContents[0:0]
	p.EndWrite()
}

func (p *DBGoodsInfo)MakeGoodsFiles(resId string) bool {
	// fmt.Println("MakeGoodsFiles GoodsId:", p.GoodsId, "resId:", resId)
	if len(resId) == 0 {
		p.ClearImage()
		return false
	}
	goodsPath := path.Join(common.ApiPath, "goods", resId)
	if !common.Exists(goodsPath) {
		idx := strings.Index(resId, "_")
		cateIdStr := ""
		if idx > 0 {
			cateIdStr = resId[0:idx]
			goodsPath = path.Join(common.ApiPath, "goods", cateIdStr, resId)
			if !common.Exists(goodsPath) {
				idx = 0
			}
		}
		if idx == 0 {
			picPath := path.Join(common.PicturePath, cateIdStr, resId+".png")
			if !common.Exists(picPath) {
				picPath = path.Join(common.PicturePath, cateIdStr, resId+".jpg")
				if !common.Exists(picPath) {
					p.ClearImage()
					return false
				}
			}
			picPath = removeDir2Url(picPath)

			iconPath := path.Join(common.IconPath, cateIdStr, resId+".png")
			if !common.Exists(iconPath) {
				iconPath = ""
			} else {
				iconPath = removeDir2Url(iconPath)
			}

			// fmt.Println("MakeGoodsFiles2 icon:", iconPath)
			// fmt.Println("MakeGoodsFiles2 pics:", picPath)

			p.BeginWrite()
			p.icon = iconPath
			p.pics = picPath
			p.contents = picPath

			p.Icon = iconPath
			p.Pics = picPath
			p.Contents = picPath
			p.EndWrite()

			p.GetPublicIcon()
			p.GetPublicPics()
			return true
		}
	}
	var allFiles []string
	allFiles,_ = common.GetFilesList(goodsPath, allFiles)
	if len(allFiles) == 0 {
		p.ClearImage()
		return false
	}
	gobalFile := ""
	for _, filePath := range allFiles {
		if common.IsFile(filePath) {
			gobalFile = filePath
		}
	}
	iconPath := path.Join(goodsPath, "icon")
	picPath := path.Join(goodsPath, "pic")
	contentPath := path.Join(goodsPath, "content")

	pics := ""
	var picFiles []string
	if common.Exists(picPath) {
		picFiles, _ = common.GetAllFiles(picPath, picFiles)
		common.SortFileList(&picFiles)
		pics = removeDir2Urls(picFiles)
	}
	if len(pics) == 0 {
		pics = gobalFile
		picFiles = append(picFiles, gobalFile)
	}

	icon := ""
	if common.Exists(iconPath) {
		var iconFiles []string
		iconFiles,_ = common.GetAllFiles(iconPath, iconFiles)
		if len(iconFiles) > 0 {
			icon = iconFiles[0]
			icon = removeDir2Url(icon)
		}
	}
	if len(icon) == 0 {
		icon = picFiles[0]
		icon = removeDir2Url(icon)
	}

	contents := ""
	if common.Exists(contentPath) {
		var contentFiles []string
		contentFiles,_ = common.GetAllFiles(contentPath, contentFiles)
		common.SortFileList(&contentFiles)
		contents = removeDir2Urls(contentFiles)
	}
	if len(contents) == 0 {
		contents = pics
	}

	// fmt.Println("MakeGoodsFiles2 icon:", icon)
	// fmt.Println("MakeGoodsFiles2 pics:", pics)
	// fmt.Println("MakeGoodsFiles2 contents:", contents)
	skuPath := path.Join(goodsPath, "sku")
	if common.Exists(skuPath) {
		goodsSkus := p.GetGoodsSkus()
		skuPics := map[int]string{}
		for _,goodsSku := range goodsSkus {
			if goodsSku.Id < 100 {
				resId := strconv.Itoa(goodsSku.Id)
				picPath := path.Join(skuPath, resId+".png")
				if !common.Exists(picPath) {
					picPath = path.Join(skuPath, resId+".jpg")
					if !common.Exists(picPath) {
						continue
					}
				}
				picPath = removeDir2Url(picPath)
				picPath = addPublicUrlHost(picPath)
				if len(picPath) > 0 {
					skuPics[goodsSku.Id] = picPath
				}
			}
		}
		p.skuPics = skuPics
	}

	p.BeginWrite()
	p.icon = icon
	p.pics = pics
	p.contents = contents
	//p.skuPics = skuPics

	p.Icon = icon
	p.Pics = pics
	p.Contents = contents
	p.EndWrite()

	p.GetPublicIcon()
	p.GetPublicPics()
	p.GetPublicContents()

	return true
}

func removeDir2Url(url string)string {
	if len(url) == 0 {
		return ""
	}
	hostUrl := common.StaticPath
	if strings.Contains(url, hostUrl) {
		url = url[len(hostUrl):]
	}
	return url
}

func removeDir2Urls(urlPaths []string)string {
	if len(urlPaths) == 0{
		return ""
	}
	var items []string
	for _,tmp := range urlPaths {
		item := removeDir2Url(tmp)
		items = append(items, item)
	}
	return strings.Join(items, ";")
}

func addPublicUrlHost(url string)string {
	if len(url) == 0 {
		return ""
	}
	if !strings.Contains(url, "http") {
		url = Module_GetFilePublicPath(url)
		if len(url) == 0 {
			return url
		}
		if strings.HasPrefix(url, "/") {
			url = common.StaticUrl+"/static"+url
		}else{
			url = common.StaticUrl+"/static/"+url
		}
	}
	return url
}

func addPublicUrlHosts(urlPath string)[]string {
	var items []string
	if len(urlPath) == 0 {
		return items
	}
	tmps := strings.Split(urlPath, ";")
	for _,tmp := range tmps {
		item := addPublicUrlHost(tmp)
		items = append(items, item)
	}
	return items
}

func splitItems(txt string)[]string {
	var items []string
	if len(txt) == 0{
		return items
	}
	items = strings.Split(txt, ";")
	return items
}

func removeUrlHost(url string)string {
	if len(url) == 0 {
		return ""
	}
	hostUrl := common.StaticUrl+"/static/"
	if strings.Contains(url, hostUrl) {
		url = url[len(hostUrl):]
	}
	return url
}

func removeUrlHosts(urlPath string)string {
	if len(urlPath) == 0{
		return ""
	}
	tmps := strings.Split(urlPath, ";")
	var items []string
	for _,tmp := range tmps {
		item := removeUrlHost(tmp)
		items = append(items, item)
	}
	return strings.Join(items, ";")
}

func addUrlHost(url string)string {
	if len(url) == 0 {
		return ""
	}
	if !strings.Contains(url, "http") {
		if strings.HasPrefix(url, "/") {
			url = common.StaticUrl+"/static"+url
		}else{
			url = common.StaticUrl+"/static/"+url
		}
	}
	return url
}

func addUrlHosts(urlPath string)[]string {
	var items []string
	if len(urlPath) == 0 {
		return items
	}
	tmps := strings.Split(urlPath, ";")
	for _,tmp := range tmps {
		item := addUrlHost(tmp)
		items = append(items, item)
	}
	return items
}

// func init() {
// 	goodsInfo := &DBGoodsInfo{}
// 	goodsInfo.TryLoadImage() 
// }