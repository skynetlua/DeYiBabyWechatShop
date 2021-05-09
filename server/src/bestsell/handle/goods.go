package handle

import (
	"bestsell/common"
	"bestsell/module"
	"bestsell/mysqld"
	"github.com/kataras/iris/v12"
	"strings"
)

//=>/goods/category/all true get  
func Goods_category_all(ctx iris.Context, sess *common.BSSession) {
	dbCategorys := mysqld.GetDBCategoryList()
	var data []*map[string]interface{}
	for _, item := range dbCategorys {
		if item.Status != 1 {
			continue
		}
		itemInfo := &map[string]interface{}{
			"id" 	:item.ID,
			"name"  :item.Name,
			"icon" 	:item.GetPublicIcon(),
			"level" :item.Level,
			"order" :item.Order,
		}
		data = append(data, itemInfo)
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/goods/category/subtypes true get {categoryId}
func Goods_category_subtypes(ctx iris.Context, sess *common.BSSession) {
	categoryId := common.AtoIDefault(ctx.FormValue("categoryId"), -1)
	cateTypeGoodsList := mysqld.GetSubGoodsGroupByByCategoryId(categoryId)
	if cateTypeGoodsList == nil {
		ctx.JSON(iris.Map{"code": -1, "msg": "子类别不存在"})
		return
	}
	var data []*map[string]interface{}
	for typeName, itemList := range *cateTypeGoodsList {
		if len(*itemList) == 0 {
			continue
		}
		mainType := ""
		icon := ""
		dbGoods := (*itemList)[0]
		dbGoodsInfo := dbGoods.GetGoodsInfo()
		if dbGoodsInfo != nil {
			icon = dbGoodsInfo.GetPublicIcon()
			mainType = dbGoodsInfo.MainType
		}
		itemInfo := &map[string]interface{}{
			"categoryId": categoryId,
			"mainType": mainType,
			"name" : typeName,
			"icon" : icon,
		}
		data = append(data, itemInfo)
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/goods/category/sublist true get {categoryId, subType}
func Goods_category_sublist(ctx iris.Context, sess *common.BSSession) {
	categoryId := common.AtoIDefault(ctx.FormValue("categoryId"), -1)
	subType := ctx.FormValue("subType")

	cateTypeGoodsList := mysqld.GetSubGoodsGroupByByCategoryId(categoryId)
	if cateTypeGoodsList == nil {
		ctx.JSON(iris.Map{"code": -1, "msg": "子类别不存在"})
		return
	}
	goodsList, ok := (*cateTypeGoodsList)[subType]
	if !ok {
		ctx.JSON(iris.Map{"code": -1, "msg": "子类别不存在"})
		return
	}
	ctx.JSON(iris.Map{"code": 0, "data": goodsList})
}

//=>/goods/category/info true get {id} 
func Goods_category_info(ctx iris.Context, sess *common.BSSession) {
	empty("/goods/category/info")
}

//=>/goods/list true post {} 
func Goods_list(ctx iris.Context, sess *common.BSSession) {
	var retGoodSlice *[]*mysqld.DBGoods
	var retCode = 0
	for {
		page := common.AtoI(ctx.FormValue("page"))
		pageSize := common.AtoIDefault(ctx.FormValue("pageSize"), 20)
		categoryId := common.AtoIDefault(ctx.FormValue("categoryId"), -1)
		var goodsSlice *[]*mysqld.DBGoods
		if categoryId >= 0 {
			//tmp := mysqld.GetDBGoodsListByCategoryId(categoryId, mysqld.GoodsStatusUp)
			tmp := mysqld.GetGoodsSliceByByCategoryId(categoryId)
			goodsSlice = &tmp
		}
		if goodsSlice == nil {
			//tmp := mysqld.GetDBGoodsList(mysqld.GoodsStatusUp)
			tmp := mysqld.GetAllGoodsSlice()
			goodsSlice = &tmp
		}
		nameLike := ctx.FormValue("nameLike")
		if len(nameLike)>0 {
			tmpGoodsSlice := goodsSlice
			goodsSlice = &[]*mysqld.DBGoods{}
			for _,item := range *tmpGoodsSlice{
				dbGoodsInfo := item.GetGoodsInfo()
				if strings.Contains(dbGoodsInfo.Name, nameLike) {
					*goodsSlice = append(*goodsSlice, item)
				}
			}
			mysqld.SortGoodsSlice(goodsSlice)
		}

		num := len(*goodsSlice)
		startIdx := page * pageSize
		endIdx := (page + 1) * pageSize
		if startIdx > num {
			startIdx = num
		}
		if endIdx > num {
			endIdx = num
		}

		_goodsSlice := (*goodsSlice)[startIdx:endIdx]
		retGoodSlice = &_goodsSlice
		if len(_goodsSlice) < pageSize {
			retCode = 10000
		}
		break
	}
	var data []*map[string]interface{}
	for _, dbGoods := range *retGoodSlice {
		if dbGoods.GetTag() > 0 {
			continue
		}
		itemInfo := dbGoods.GetInfo()
		data = append(data, itemInfo)
	}
	ctx.JSON(iris.Map{"code": retCode, "data": data})
}

//=>/goods/detail true get {id} 
func Goods_detail(ctx iris.Context, sess *common.BSSession) {
	goodsId := common.AtoI(ctx.FormValue("goodsId"))
	dbGoods := mysqld.GetDBGoodsOrFromDB(goodsId)
	if dbGoods == nil || dbGoods.GoodsId != goodsId {
		ctx.JSON(iris.Map{"code": -1, "msg":"商品id出错"})
		return
	}
	goodsDetail := dbGoods.GetDetailInfo()
	data := map[string]interface{} {
		"goods":goodsDetail,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/goods/sku true get {id} 
func Goods_sku(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}

	goodsId := common.AtoI(ctx.FormValue("id"))
	dbGoods := mysqld.GetDBGoods(goodsId)
	if dbGoods == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	dbGoodsInfo := dbGoods.GetGoodsInfo()
	if dbGoodsInfo == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}

	if dbGoods.GetStatus() != mysqld.EGoodsStatusUp {
		if player.GM == 0 {
			ctx.JSON(iris.Map{"code": -1, "msg": "商品已下架"})
			return
		}
	}

	buyCount := player.GoodsBuyCount(goodsId)
	buyLimit := 0
	if dbGoodsInfo.SellPrice == 0 {
		buyLimit = 1
	}

	data := map[string]interface{}{
		"goodsId": dbGoodsInfo.GoodsId,
		"skuJson": dbGoodsInfo.SkuStruct,
		"skuPics":dbGoodsInfo.GetSkuPics(),
		"numberStore": dbGoods.NumberStore,
		"buyCount": buyCount,
		"buyLimit": buyLimit,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/goods/price true post {goodsId,propertyChildIds} 
func Goods_price(ctx iris.Context, sess *common.BSSession) {
	empty("/goods/price")
}

//=>/goods/reputation true post {} 
func Goods_reputation(ctx iris.Context, sess *common.BSSession) {
	goodsId := common.AtoI(ctx.FormValue("goodsId"))
	dbGoods := mysqld.GetDBGoods(goodsId)
	if dbGoods == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	dbGoodsInfo := dbGoods.GetGoodsInfo()
	if dbGoodsInfo == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	//cfgGoodsDetail := cfgGoods.GetDetail()
	reputationBox := mysqld.GetDBReputationBoxOrFromDB(goodsId)
	items := reputationBox.GetItems()
	var datas []*map[string]interface{}
	for _,item := range *items {
		data := &map[string]interface{}{
			"repute":item.Repute,
			"remark":item.Remark,
			"avatarUrl":item.AvatarUrl,
			"skuName":dbGoodsInfo.GetSkuNames(item.SkuId),
			"dateStr":item.GetUpdateDateStr(),
		}
		datas = append(datas, data)
	}
	ctx.JSON(iris.Map{"code": 0, "data": datas})
}