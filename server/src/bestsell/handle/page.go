package handle

import (
	"bestsell/common"
	"bestsell/config"
	"bestsell/module"
	"bestsell/mysqld"
	"github.com/kataras/iris/v12"
	// "fmt"
	"time"
)

//=>/page/index true get  
func Page_index(ctx iris.Context, sess *common.BSSession) {
	dbCategorys := mysqld.GetDBCategoryList()
	var categoryList []*map[string]interface{}
	var goodsGroups []*map[string]interface{}
	for _, item := range dbCategorys {
		if item.Status != 1 {
			continue
		}
		itemInfo := &map[string]interface{}{
			"id" :item.ID,
			"name" :item.Name,
			"level":item.Level,
			"order":item.Order,
			"icon":item.GetPublicIcon(),
			"pic":item.GetPublicPic(),
		}
		categoryList = append(categoryList, itemInfo)

		goodsList := mysqld.GetGoodsSliceByByCategoryId(item.ID)
		//if len(goodsList) > 6 {
		//	goodsList = goodsList[0:6]
		//}
		var goodsInfos []*map[string]interface{}
		for _, goods := range goodsList {
			if goods.GetTag() > 0 {
				continue
			}
			itemInfo := goods.GetInfo()
			goodsInfos = append(goodsInfos, itemInfo)
			if len(goodsInfos) >= 6 {
				break
			}
		}
		goodsGroup := &map[string]interface{}{
			"categoryId" :item.ID,
			"goodsList": goodsInfos,
		}
		goodsGroups = append(goodsGroups, goodsGroup)
	}
	banners := config.GetCfgUiBannerSliceByMap("index")

	dbGoodsGroup := mysqld.GetDBGoodsListGroup()
	var recomGoodsDatas []*map[string]interface{}
	for _, item := range *((*dbGoodsGroup)[int(mysqld.EGoodsMarkRecomm)]) {
		itemInfo := item.GetInfo()
		recomGoodsDatas = append(recomGoodsDatas, itemInfo)
	}

	curTime := time.Now().Unix()
	var seckillGoodsDatas []*map[string]interface{}
	for _, item := range *((*dbGoodsGroup)[int(mysqld.EGoodsMarkSeckill)]) {
		if item.EndTime < curTime {
			item.SetStatus(mysqld.EGoodsStatusDown)
			continue
		}
		itemInfo := item.GetInfo()
		(*itemInfo)["startTime"] = item.StartTime
		(*itemInfo)["endTime"]   = item.EndTime
		seckillGoodsDatas = append(seckillGoodsDatas, itemInfo)
	}

	var teamGoodsDatas []*map[string]interface{}
	for _, item := range *((*dbGoodsGroup)[int(mysqld.EGoodsMarkTeam)]) {
		//if item.EndTime < curTime {
		//	item.SetStatus(mysqld.EGoodsStatusDown)
		//	continue
		//}
		itemInfo := item.GetInfo()
		(*itemInfo)["startTime"] = item.StartTime
		(*itemInfo)["endTime"]   = item.EndTime
		teamGoodsDatas = append(teamGoodsDatas, itemInfo)
	}

	dbNoticeList := mysqld.GetDBNoticeList()
	var noticeList []*map[string]interface{}
	for _,dbNotice := range dbNoticeList {
		_notice := map[string]interface{}{
			"id"   :dbNotice.ID,
			"title":dbNotice.Title,
		}
		noticeList = append(noticeList, &_notice)
	}

	goodsDynamic := mysqld.GetLastGoodsStatList()

	data := map[string]interface{}{
		"categories" 	:categoryList,
		"recomGoods" 	:recomGoodsDatas,
		"seckillGoods" 	:seckillGoodsDatas,
		"teamGoods" 	:teamGoodsDatas,
		"banners" 		:banners,
		"noticeList" 	:noticeList,
		"goodsGroups" 	:goodsGroups,
		"goodsDynamic" 	:goodsDynamic,
	}
	player :=  module.GetPlayer(sess)
	if player != nil {
		data["gm"] = player.GM

		cartBox := player.GetCartBox()
		data["cartCount"] = cartBox.GetCartCount()
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/page/goods/detail true get  
func Page_goods_detail(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	goodsId := common.AtoI(ctx.FormValue("goodsId"))
	dbGoods := mysqld.GetDBGoodsOrFromDB(goodsId)
	if dbGoods == nil || dbGoods.GoodsId != goodsId{
		ctx.JSON(iris.Map{"code": -1, "msg":"商品id出错"})
		return
	}
	dbGoodsInfo := dbGoods.GetGoodsInfo()
	if dbGoodsInfo == nil || goodsId != dbGoodsInfo.GoodsId {
		ctx.JSON(iris.Map{"code": -1, "msg":"商品id出错"})
		return
	}
	if dbGoods.GetStatus() != mysqld.EGoodsStatusUp {
		if player.GM == 0 {
			ctx.JSON(iris.Map{"code": -1, "msg": "商品已下架"})
			return
		}
	}

	cartNum := player.GetCartBox().GetCartCount()

	favoriteBox := player.GetFavoriteBox()
	dbFavorite := favoriteBox.GetFavoriteByGoodsId(goodsId)
	isFavorite := 0
	if dbFavorite != nil {
		isFavorite = 1
	}
	logistics := map[string]interface{}{
		//"isFree": cfgLogistics.IsFree,
		//"freeType": cfgLogistics.FreeType,
		//"freeTypeStr": cfgLogistics.FreeTypeStr,
	}

	mysqld.MakeVisitGoodsStat(mysqld.EGoodsActionVisit, player, dbGoods, 0, 0, 0)

	buyCount := player.GoodsBuyCount(goodsId)
	buyLimit := 0
	if dbGoodsInfo.SellPrice == 0 {
		buyLimit = 10
	}

	goods := dbGoods.GetDetailInfo()
	data := map[string]interface{}{
		"goods":goods,
		"playerId": player.ID,
		"shopId": 1,
		"logistics":logistics,
		"cartNum": cartNum,
		"isFavorite": isFavorite,
		"buyCount": buyCount,
		"buyLimit": buyLimit,
	}

	if dbGoods.GetTag() == int(mysqld.EGoodsMarkTeam) {
		teamBuyOrders := mysqld.GetGoodsTeamBuyOrders(goodsId)
		data["teamBuyList"] = teamBuyOrders
	}

	ctx.JSON(iris.Map{"code": 0, "data": data})
}
