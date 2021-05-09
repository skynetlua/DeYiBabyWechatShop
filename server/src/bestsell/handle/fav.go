package handle

import (
	"bestsell/common"
	"bestsell/module"
	"bestsell/mysqld"
	"github.com/kataras/iris/v12"
	"strconv"
	"fmt"
)

//=>/fav/list true post {} 
func Fav_list(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	items := []map[string]interface{}{}
	favoriteBox := player.GetFavoriteBox()
	dbFavorites := favoriteBox.GetFavorites()
	for _,dbFavorite := range *dbFavorites{
		dbGoods := mysqld.GetDBGoods(dbFavorite.GoodsId)
		if dbGoods == nil {
			continue
		}
		dbGoodsInfo := dbGoods.GetGoodsInfo()
		item := map[string]interface{}{
			"goodsId" : dbGoods.GoodsId,
			"pic": dbGoodsInfo.GetPublicIcon(),
			"name" : dbGoodsInfo.Name,
		}
		items = append(items, item)
	}
	data := items
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/fav/add true post {token,goodsId} 
func Fav_add(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	_goodsId := ctx.FormValue("goodsId")
	goodsId,err := strconv.Atoi(_goodsId)
	if err != nil {
		fmt.Println("Fav_add ",err)
		ctx.JSON(iris.Map{"code": -1, "msg":"商品id出错"})
		return
	}
	dbGoods := mysqld.GetDBGoods(goodsId)
	if dbGoods == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"该商品已存在"})
		return
	}
	favoriteBox := player.GetFavoriteBox()
	dbFavorite := favoriteBox.GetFavoriteByGoodsId(goodsId)
	if dbFavorite != nil {
		ctx.JSON(iris.Map{"code": 0, "msg":"已添加收藏"})
		return
	}
	dbFavorite = &mysqld.DBFavorite{
		GoodsId:goodsId,
	}
	favoriteBox.AddFavorite(dbFavorite)
	ctx.JSON(iris.Map{"code": 0, "msg":"收藏添加成功"})
}

//=>/fav/check true get {token,goodsId} 
func Fav_check(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	_goodsId := ctx.FormValue("goodsId")
	goodsId,err := strconv.Atoi(_goodsId)
	if err != nil {
		fmt.Println("Fav_add ",err)
		ctx.JSON(iris.Map{"code": -1, "msg":"商品id出错"})
		return
	}
	favoriteBox := player.GetFavoriteBox()
	dbFavorite := favoriteBox.GetFavoriteByGoodsId(goodsId)
	if dbFavorite != nil {
		ctx.JSON(iris.Map{"code": 0})
		return
	}
	ctx.JSON(iris.Map{"code": -1})
}

//=>/fav/delete true post {token,goodsId} 
func Fav_delete(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"请先登录"})
		return
	}
	_goodsId := ctx.FormValue("goodsId")
	goodsId,err := strconv.Atoi(_goodsId)
	if err != nil {
		fmt.Println("Fav_add ",err)
		ctx.JSON(iris.Map{"code": -1, "msg":"商品id出错"})
		return
	}
	favoriteBox := player.GetFavoriteBox()
	favoriteBox.RemoveFavoriteByGoodsId(goodsId)
	ctx.JSON(iris.Map{"code": 0, "msg":"收藏删除成功"})
}