package handle

import (
	"bestsell/common"
	"bestsell/module"
	"bestsell/mysqld"
	"fmt"
	"time"
	"github.com/kataras/iris/v12"
)

//=>/cart/info true get {token} 
func Cart_info(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	cartBox := player.GetCartBox()
	data := map[string]interface{}{
		"number": cartBox.GetCartCount(),
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/cart/list true get {token} 
func Cart_list(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"请先登陆"})
		return
	}
	cartBox := player.GetCartBox()
	carts := cartBox.GetCarts()
	var items []map[string]interface{}
	amountGoods := 0
	for _, cart := range carts {
		dbGoods := mysqld.GetDBGoods(cart.GoodsId)
		if dbGoods == nil {
			fmt.Println("Cart_list GoodsId =", cart.GoodsId)
			cartBox.RemoveCart(cart.ID)
			continue
		}
		dbGoodsInfo := dbGoods.GetGoodsInfo()
		realPrice := dbGoodsInfo.GetRealPrice(cart.SkuId)
		item := map[string]interface{} {
			"id": cart.ID,
			"goodsId": cart.GoodsId,
			"skuId": cart.SkuId,
			"name": dbGoodsInfo.Name,
			"icon": dbGoodsInfo.GetOrderPublicIcon(cart.SkuId),
			"price":  realPrice,
			"numberBuy": cart.NumberBuy,
			"numberStore": dbGoods.NumberStore,
			"skuNames": dbGoodsInfo.GetSkuNames(cart.SkuId),
		}
		items = append(items, item)
		amountGoods += realPrice*cart.NumberBuy
	}
	data := map[string]interface{}{
		"items": items,
		"amountGoods": amountGoods,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/cart/quick true get {token, goodsId, skuId, buyNumber}
func Cart_quick(ctx iris.Context, sess *common.BSSession) {
}

//=>/cart/add true post {token,goodsId,number,sku} 
func Cart_add(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"需要登陆才能操作"})
		return
	}

	_skuId := ctx.FormValue("skuId")
	skuId := common.AtoIDefault(_skuId, -1)
	if skuId == -1 {
		ctx.JSON(iris.Map{"code": 30002})
		return
	}
	goodsId := common.AtoI(ctx.FormValue("goodsId"))
	buyNumber := common.AtoI(ctx.FormValue("buyNumber"))
	dbGoods := mysqld.GetDBGoods(goodsId)
	if dbGoods == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"该商品不存在，请咨询客服"})
		return
	}
	if buyNumber > dbGoods.NumberStore {
		ctx.JSON(iris.Map{"code": -1, "msg":"商品库存不够，请咨询客服"})
		return
	}
	//if player.IsHadBuy(goodsId) {
	//	ctx.JSON(iris.Map{"code": -1, "msg":"很抱歉，您已经买过该免费产品"})
	//	return
	//}
	dbGoodsInfo := dbGoods.GetGoodsInfo()
	if dbGoodsInfo.SellPrice == 0 && buyNumber > 1 {
		ctx.JSON(iris.Map{"code": -1, "msg":"很抱歉，免费产品只能购买1件"})
		return
	}
	if dbGoods.Mark == int(mysqld.EGoodsMarkSeckill) || dbGoods.Mark == int(mysqld.EGoodsMarkTeam) {
		curTime := time.Now().Unix()
        endInterval := dbGoods.EndTime-curTime
        if endInterval <= 0 {
        	ctx.JSON(iris.Map{"code": -1, "msg":"很抱歉，活动已结束！"})
        	return
        }
        startInterval := dbGoods.StartTime-curTime
        if startInterval > 0 {
        	ctx.JSON(iris.Map{"code": -1, "msg":"很抱歉，活动未开始！"})
        	return
        }
    }
	if !dbGoodsInfo.IsValidOrderSkuId(skuId) {
		ctx.JSON(iris.Map{"code": -1, "msg":"选择的商品规格匹配不上！"})
		return
	}
	mysqld.MakeVisitGoodsStat(mysqld.EGoodsActionCart, player, dbGoods, skuId, 0, buyNumber)

	cartBox := player.GetCartBox()
	carts := cartBox.GetCarts()
	var theCart *mysqld.DBCart
	for _, cart := range carts{
		if cart.GoodsId == goodsId && cart.SkuId == skuId {
			if dbGoodsInfo.SellPrice == 0 && cart.NumberBuy > 0 {
				ctx.JSON(iris.Map{"code": -1, "msg":"很抱歉，免费产品只能购买1件"})
				return
			}
			cart.BeginWrite()
			cart.NumberBuy += buyNumber
			cart.EndWrite()
			cart.DelaySave(cart)

			theCart = cart
			break
		}
	}
	if theCart == nil {
		cart := &mysqld.DBCart{
			PlayerId: cartBox.PlayerId,
			GoodsId: goodsId,
			SkuId: skuId,
			NumberBuy: buyNumber,
		}
		cartBox.AddCart(cart)
		theCart = cart
	}
	fmt.Println("Cart_add cartId =", theCart.ID, "goodsId =", theCart.GoodsId, "skuId =", theCart.SkuId, "numberBuy =", theCart.NumberBuy)
	ctx.JSON(iris.Map{"code": 0, "msg":"添加成功"})
}

//=>/cart/modifyNumber true post {token,key,number} 
func Cart_modifyNumber(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	id := common.AtoI(ctx.FormValue("id"))
	number := common.AtoI(ctx.FormValue("number"))
	cartBox := player.GetCartBox()
	cart := cartBox.GetCart(id)
	if cart == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"该订单不存在，请咨询客服"})
		return
	}
	dbGoods := mysqld.GetDBGoods(cart.GoodsId)
	if number > dbGoods.NumberStore {
		ctx.JSON(iris.Map{"code": -1, "msg":"商品库存不够，请咨询店员"})
		return
	}
	dbGoodsInfo := dbGoods.GetGoodsInfo()
	if dbGoodsInfo.SellPrice == 0 && cart.NumberBuy + number > 1 {
		ctx.JSON(iris.Map{"code": -1, "msg":"很抱歉，免费产品只能购买1件"})
		return
	}
	cart.BeginWrite()
	cart.NumberBuy = number
	cart.EndWrite()
	cart.DelaySave(cart)
	fmt.Println("Cart_modifyNumber cartId =", cart.ID, "goodsId =", cart.GoodsId, "skuId =", cart.SkuId, "numberBuy =", cart.NumberBuy)
	data := map[string]interface{}{
		"id": cart.ID,
		"numberBuy": cart.NumberBuy,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data, "msg": "数量修改成功"})
}

//=>/cart/remove true post {token,key} 
func Cart_remove(ctx iris.Context, sess *common.BSSession) {
	player := module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	id := common.AtoI(ctx.FormValue("id"))
	cartBox := player.GetCartBox()
	cart := cartBox.GetCart(id)
	if cart == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"该订单不存在，请咨询客服"})
		return
	}
	fmt.Println("Cart_remove cartId =", cart.ID, "goodsId =", cart.GoodsId, "skuId =", cart.SkuId, "numberBuy =", cart.NumberBuy)
	cartBox.RemoveCart(id)
	ctx.JSON(iris.Map{"code": 0, "msg":"删除成功"})
}

//=>/cart/empty true post {token} 
func Cart_empty(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	cartBox := player.GetCartBox()
	cartBox.Clear()
	ctx.JSON(iris.Map{"code": 0, "msg":"清空成功"})
}