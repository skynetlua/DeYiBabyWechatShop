package handle

import (
	"bestsell/common"
	"bestsell/module"
	"bestsell/mysqld"
	"bestsell/sdk"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"time"
)

var refundReasons = map[int]string{
	0:"不喜欢/不想要",
	1:"空包裹",
	2:"未按约定时间发货",
	3:"快递/物流一直未送达",
	4:"货物破损已拒签",
	5:"退运费",
	6:"规格尺寸与商品页面描述不符",
	7:"功能/效果不符",
	8:"质量问题",
	9:"少件/漏发",
	10:"包装/商品破损",
	11:"发票问题",
}


//=>/order/prepare true post {}
func Order_prepare(ctx iris.Context, sess *common.BSSession) {
	player := module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	sendType := common.AtoI(ctx.FormValue("sendType"))
	quickBuy := common.AtoI(ctx.FormValue("quickBuy"))
	teamBuy := common.AtoI(ctx.FormValue("teamBuy"))
	carts := []*mysqld.DBCart{}
	if quickBuy == 1 {
		_skuId := ctx.FormValue("orderSkuId")
		orderSkuId := common.AtoIDefault(_skuId, -1)
		if orderSkuId == -1 {
			ctx.JSON(iris.Map{"code": 30002, "msg":"请选择规格"})
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
		dbGoodsInfo := dbGoods.GetGoodsInfo()
		if dbGoodsInfo.SellPrice == 0 && buyNumber > 1 {
			ctx.JSON(iris.Map{"code": -1, "msg":"很抱歉，免费产品只能购买1件"})
			return
		}
		//if dbGoods.Mark == int(mysqld.EGoodsMarkSeckill) || dbGoods.Mark == int(mysqld.EGoodsMarkTeam) {
		if dbGoods.Mark == int(mysqld.EGoodsMarkSeckill) {
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
		if !dbGoodsInfo.IsValidOrderSkuId(orderSkuId) {
			ctx.JSON(iris.Map{"code": -1, "msg":"商品规格匹配不上！"})
			return
		}
		cart := &mysqld.DBCart{
			PlayerId: player.ID,
			GoodsId: goodsId,
			SkuId: orderSkuId,
			NumberBuy: buyNumber,
		}
		carts = append(carts, cart)
	} else {
		cartBox := player.GetCartBox()
		carts = cartBox.GetCarts()
	}

	//if len(carts) == 0 {
	//	ctx.JSON(iris.Map{"code": -1, "msg":"有商品已下架，请联系客服"})
	//	return
	//}
	var items []*mysqld.DBOrderGoodsInfo
	orderTag := 0
	for _,cart := range carts {
		dbGoods := mysqld.GetDBGoods(cart.GoodsId)
		if dbGoods == nil {
			ctx.JSON(iris.Map{"code": -1, "msg":"有商品已下架，请联系客服"})
			return
		}
		dbGoodsInfo := dbGoods.GetGoodsInfo()
		if cart.NumberBuy > dbGoods.NumberStore {
			ctx.JSON(iris.Map{"code": -1, "msg":"商品["+dbGoodsInfo.Name+"]库存不够，请联系客服"})
			return
		}
		//if dbGoodsInfo.SellPrice == 0 && player.IsHadBuy(cart.GoodsId) {
		//	ctx.JSON(iris.Map{"code": -1, "msg":"很抱歉，您已经买过免费商品["+dbGoodsInfo.Name+"]"})
		//	return
		//}
		orderSKuId := cart.SkuId
		item := &mysqld.DBOrderGoodsInfo {
			GoodsId:  cart.GoodsId,
			Number:   cart.NumberBuy,
			SkuId:    orderSKuId,
			Tag:  	  dbGoods.GetTag(),
			Name:     dbGoodsInfo.Name,
			Pic:      dbGoodsInfo.GetOrderPublicIcon(orderSKuId),
			Price:    dbGoodsInfo.SellPrice,
			SkuName:  dbGoodsInfo.GetSkuNames(orderSKuId),
			SkuPrice: dbGoodsInfo.GetSkuPrice(orderSKuId),
			CartId:   cart.ID,
		}

		if quickBuy == 1 {
			if teamBuy > 0 {
				if teamBuy == 1 {
					item.Price = dbGoodsInfo.OriginPrice
					orderTag = 0
					item.Tag = 0
				} else {
					orderTag = item.Tag
				}
			} else {
				if item.Tag > 0 {
					orderTag = item.Tag
				}
			}
		}

		items = append(items, item)
		mysqld.MakeVisitGoodsStat(mysqld.EGoodsActionPrepare, player, dbGoods, orderSKuId, 0, cart.NumberBuy)
	}

	if len(items) == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"有商品已下架，请联系客服"})
		return
	}

	amountCoupon := 0
	amountLogistics := 0
	dbOrder := &mysqld.DBOrder {
		Status:          int(mysqld.EStatusPay),
		IsTips:      	 1,
		AmountCoupon:    amountCoupon,
		AmountLogistics: amountLogistics,
		QuickBuy: 		 quickBuy,
		SendType:  		 sendType,
		TeamBuy: 		 teamBuy,
		Tag: 			 orderTag,
	}
	dbOrder.SetGoodsInfos(&items)
	ok := dbOrder.CalGoodsAmount()
	if ok == -1 {
		ctx.JSON(iris.Map{"code": -1, "msg":"有商品已下架，请联系客服"})
		return
	} else if ok == -2 {
		ctx.JSON(iris.Map{"code": -1, "msg":"商品库存不够，请联系客服"})
		return
	}
	dbOrder.CalAmountReal()

	data := dbOrder.GetOrderInfo()
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/order/create true post {} 
func Order_create(ctx iris.Context, sess *common.BSSession) {
	player := module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	sendType := common.AtoI(ctx.FormValue("sendType"))
	quickBuy := common.AtoI(ctx.FormValue("quickBuy"))
	teamBuy := common.AtoI(ctx.FormValue("teamBuy"))
	carts := []*mysqld.DBCart{}
	cartIds := []int{}
	orderTag := 0
	if quickBuy == 1 {
		_skuId := ctx.FormValue("orderSkuId")
		orderSkuId := common.AtoIDefault(_skuId, -1)
		if orderSkuId == -1 {
			ctx.JSON(iris.Map{"code": 30002, "msg":"请选择规格"})
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
		dbGoodsInfo := dbGoods.GetGoodsInfo()
		if dbGoodsInfo.SellPrice == 0 && buyNumber > 1 {
			ctx.JSON(iris.Map{"code": -1, "msg":"很抱歉，免费产品只能购买1件"})
			return
		}
		if dbGoods.Mark == int(mysqld.EGoodsMarkSeckill) {
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
		if !dbGoodsInfo.IsValidOrderSkuId(orderSkuId) {
			ctx.JSON(iris.Map{"code": -1, "msg":"商品规格匹配不上！"})
			return
		}
		cart := &mysqld.DBCart{
			PlayerId: player.ID,
			GoodsId: goodsId,
			SkuId: orderSkuId,
			NumberBuy: buyNumber,
		}
		carts = append(carts, cart)
	} else {
		cartBox := player.GetCartBox()
		carts = cartBox.GetCarts()
		for _,cart := range carts {
			cartIds = append(cartIds, cart.ID)
		}
	}

	var items []*mysqld.DBOrderGoodsInfo
	for _,cart := range carts {
		dbGoods := mysqld.GetDBGoods(cart.GoodsId)
		if dbGoods == nil {
			ctx.JSON(iris.Map{"code": -1, "msg":"有商品已下架，请联系客服"})
			return
		}
		dbGoodsInfo := dbGoods.GetGoodsInfo()
		if cart.NumberBuy > dbGoods.NumberStore {
			ctx.JSON(iris.Map{"code": -1, "msg":"商品["+dbGoodsInfo.Name+"]库存不够，请联系客服"})
			return
		}
		//if dbGoodsInfo.SellPrice == 0 && player.IsHadBuy(cart.GoodsId) {
		//	ctx.JSON(iris.Map{"code": -1, "msg":"很抱歉，您已经买过免费商品["+dbGoodsInfo.Name+"]"})
		//	return
		//}
		orderSKuId := cart.SkuId
		item := &mysqld.DBOrderGoodsInfo {
			GoodsId:  cart.GoodsId,
			Number:   cart.NumberBuy,
			SkuId:    orderSKuId,
			Tag:  	  dbGoods.GetTag(),

			Name:     dbGoodsInfo.Name,
			Pic:      dbGoodsInfo.GetOrderPublicIcon(orderSKuId),
			Price:    dbGoodsInfo.SellPrice,
			SkuName:  dbGoodsInfo.GetSkuNames(orderSKuId),
			SkuPrice: dbGoodsInfo.GetSkuPrice(orderSKuId),
			CartId:   cart.ID,
		}

		if quickBuy == 1 {
			if teamBuy > 0 {
				if teamBuy == 1 {
					item.Price = dbGoodsInfo.OriginPrice
					orderTag = 0
					item.Tag = 0
				} else {
					orderTag = item.Tag
				}
			} else {
				if item.Tag > 0 {
					orderTag = item.Tag
				}
			}
		}

		items = append(items, item)
		mysqld.MakeVisitGoodsStat(mysqld.EGoodsActionPrepare, player, dbGoods, orderSKuId, 0, cart.NumberBuy)
	}

	if len(items) == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"有商品已下架，请联系客服"})
		return
	}

	timeStamp := int(time.Now().Unix())
	inviterId := common.AtoI(ctx.FormValue("inviterId"))
	couponId := common.AtoI(ctx.FormValue("couponId"))
	remark := ctx.FormValue("remark")

	amountLogistics := 0
	amountCoupon := 0

	//var dbOrder *mysqld.DBOrder = nil
	//if quickBuy == 1 {
	//	if orderTag == int(mysqld.EGoodsMarkSeckill) || orderTag == int(mysqld.EGoodsMarkTeam) {
	//		//dbOrder
	//	}
	//}

	firstGoodsId := items[0].GoodsId
	dbOrder := &mysqld.DBOrder {
		Status:          int(mysqld.EStatusPay),
		GoodsId: 		 firstGoodsId,
		TimeStamp: 		 timeStamp,
		IsTips:      	 1,
		InviterId:       inviterId,
		CouponId:        couponId,
		AmountCoupon:    amountCoupon,
		AmountLogistics: amountLogistics,
		Remark:          remark,
		QuickBuy:	 	 quickBuy,
		SendType:  	 	 sendType,
		TeamBuy: 		 teamBuy,
		Tag: 			 orderTag,
		RefundId:		 0,
		RefundStatus:	 0,
	}
	if !dbOrder.SetOrderGoodsInfos(&items) {
		ctx.JSON(iris.Map{"code": -1, "msg":"创建订单的时候出错"})
		return
	}
	ok := dbOrder.CalGoodsAmount()
	if ok == -1 {
		ctx.JSON(iris.Map{"code": -1, "msg":"有商品已下架，请联系客服"})
		return
	} else if ok == -2 {
		ctx.JSON(iris.Map{"code": -1, "msg":"商品库存不够，请联系客服"})
		return
	}
	dbOrder.CalAmountReal()

	addressId := common.AtoI(ctx.FormValue("addressId"))
	addressBox := player.GetAddressBox()
	dbAddress := addressBox.GetAddress(addressId)
	if dbAddress != nil {
		areaId := dbAddress.AreaId
		if areaId == 0 {
			areaId = dbAddress.CityId
		}
		dbOrder.AreaId = areaId
		dbOrder.Address = dbAddress.Address
		dbOrder.LinkMan = dbAddress.LinkMan
		dbOrder.Mobile = dbAddress.Mobile
	} else {
		dbOrder.AreaId = 0
	}

	orderBox := player.GetOrderBox()
	orderBox.AddOrder(dbOrder)
	if dbOrder.ID == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"创建订单的时候出错"})
		return
	}

	for _,goodsInfo := range items {
		dbGoods := mysqld.GetDBGoods(goodsInfo.GoodsId)
		if dbGoods != nil {
			mysqld.MakeVisitGoodsStat(mysqld.EGoodsActionOrder, player, dbGoods, goodsInfo.SkuId, dbOrder.ID, goodsInfo.Number)
		}
	}

	if len(cartIds) > 0 {
		cartBox := player.GetCartBox()
		carts = cartBox.GetCarts()
		if len(carts) == len(cartIds) {
			cartBox.Clear()
		} else {
			for _, cartId := range cartIds {
				cartBox.RemoveCart(cartId)
			}
		}
	}

	data := dbOrder.GetOrderInfo()
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/order/list true post {} 
func Order_list(ctx iris.Context, sess *common.BSSession) {
	player := module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	orderBox := player.GetOrderBox()
	dbOrders := orderBox.GetOrders()
	if len(dbOrders) == 0 {
		ctx.JSON(iris.Map{"code": 1})
		return
	}
	var orderList []*map[string]interface{}
	goodsMap := make(map[int][]map[string]interface{})
	//timeStamp := int(time.Now().Unix())
	for _,dbOrder := range dbOrders {
		if dbOrder.IsOrderDelete() {
			continue
		}
		dbOrder.CheckStatus()
		if dbOrder.IsOrderDelete() {
			continue
		}
		//if dbOrder.IsOrderWaitPay() {
		//	continue
		//}
		item := dbOrder.GetItemInfo()
		orderList = append(orderList, item)
		goodsMap[dbOrder.ID] = *(dbOrder.GetGoodsInfoList())
		dbOrder.UpdateTips()
	}
	data := map[string]interface{} {
		"orderList": orderList,
		"goodsMap": goodsMap,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/order/detail true get {id,token,hxNumber} 
func Order_detail(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	orderId := common.AtoI(ctx.FormValue("orderId"))
	orderBox :=	player.GetOrderBox()
	dbOrder := orderBox.GetOrder(orderId)
	if dbOrder == nil {
		fmt.Println("Order_detail orderId =", ctx.FormValue("orderId"))
		ctx.JSON(iris.Map{"code": -1, "msg": "订单orderId出错"})
		return
	}
	dbOrder.CheckStatus()
	data := dbOrder.GetAllDetail()
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/order/pay true post {orderId,token}
func Order_pay(ctx iris.Context, sess *common.BSSession) {
	//player :=  module.GetPlayer(sess)
	//if player == nil {
	//	ctx.JSON(iris.Map{"code": -1})
	//	return
	//}
	//orderId := common.AtoI(ctx.FormValue("orderId"))
	//orderBox :=	player.GetOrderBox()
	//dbOrder := orderBox.GetOrder(orderId)
	//if dbOrder == nil {
	//	fmt.Println("Order_pay orderId =", ctx.FormValue("orderId"))
	//	ctx.JSON(iris.Map{"code": -1, "msg": "订单orderId出错"})
	//	return
	//}
	//if dbOrder.Status > 0 {
	//	fmt.Println("Order_pay orderId =", orderId, "status =", dbOrder.Status)
	//	ctx.JSON(iris.Map{"code": -1, "msg": "订单已支付"})
	//	return
	//}
	//if float64(dbOrder.AmountReal) > player.Balance{
	//	fmt.Println("Order_pay amountReal =", dbOrder.AmountReal, "Balance =", player.Balance)
	//	ctx.JSON(iris.Map{"code": -1, "msg": "账户余额不足"})
	//	return
	//}
	//ctx.JSON(iris.Map{"code": -1})
	//dbOrder.BeginWrite()
	//dbOrder.Status = int(mysqld.StatusSend)
	//dbOrder.IsTips = 1
	//dbOrder.EndWrite()
	//dbOrder.DelaySave(dbOrder)
	//
	//amountReal := common.MakeMoneyValue(dbOrder.AmountReal)
	//player.CostMoney(amountReal, mysqld.CashLogPay)
	//
	//ctx.JSON(iris.Map{"code": 0})
}

//=>/order/delivery true post {orderId,token} 
func Order_delivery(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	orderId := common.AtoI(ctx.FormValue("orderId"))
	orderBox :=	player.GetOrderBox()
	dbOrder := orderBox.GetOrder(orderId)
	if dbOrder == nil {
		fmt.Println("Order_detail orderId =", ctx.FormValue("orderId"))
		ctx.JSON(iris.Map{"code": -1, "msg": "订单orderId出错"})
		return
	}
	dbOrder.SetDBOrderStatus(mysqld.EStatusRepute)
	dbOrder.Save()
	ctx.JSON(iris.Map{"code": 0})
}

type kReputation struct {
	GoodsId  int `json:"goodsId"`
	Repute 	 int `json:"repute"`
	Remark 	 string `json:"remark"`
}

//=>/order/reputation true post {} 
func Order_reputation(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	orderId := common.AtoI(ctx.FormValue("orderId"))
	orderBox :=	player.GetOrderBox()
	dbOrder := orderBox.GetOrder(orderId)
	if dbOrder == nil {
		fmt.Println("Order_reputation orderId =", ctx.FormValue("orderId"))
		ctx.JSON(iris.Map{"code": -1, "msg": "订单orderId出错"})
		return
	}
	if dbOrder.GetDBOrderStatus() != int(mysqld.EStatusRepute) {
		ctx.JSON(iris.Map{"code": -1, "msg": "已评论"})
		return
	}
	_reputations := ctx.FormValue("reputations")
	var reputations []kReputation
	err := json.Unmarshal([]byte(_reputations), &reputations)
	if err != nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	dbOrder.SetDBOrderStatus(mysqld.EStatusFinish)
	dbOrder.Save()

	for _,reputation := range reputations {
		goodsInfo := dbOrder.GetOrderGoodsInfo(reputation.GoodsId, -1)
		skuId := 0
		if goodsInfo != nil {
			skuId = goodsInfo.SkuId
		}
		dbReputation := &mysqld.DBReputation{
			PlayerId 	:player.ID,
			OrderId 	:orderId,
			GoodsId 	:reputation.GoodsId,
			SkuId 		:skuId,
			Repute 		:reputation.Repute,
			Remark  	:reputation.Remark,
			PlayerName 	:player.GetPlayerInfo().NickName,
			AvatarUrl 	:player.GetPlayerInfo().AvatarUrl,
		}
		mysqld.AddNewReputation(dbReputation)
	}
	ctx.JSON(iris.Map{"code": 0})
}

//=>/order/close true post {orderId,token} 
func Order_close(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	orderId := common.AtoI(ctx.FormValue("orderId"))
	orderBox :=	player.GetOrderBox()
	dbOrder := orderBox.GetOrder(orderId)
	if dbOrder == nil {
		fmt.Println("Order_close not find order orderId =", ctx.FormValue("orderId"))
		ctx.JSON(iris.Map{"code": -1, "msg": "很抱歉，订单不存在"})
		return
	}
	if !dbOrder.CloseOrder() {
		ctx.JSON(iris.Map{"code": -1, "msg": "很抱歉，当前订单不允许关闭"})
		return
	}
	//status := dbOrder.GetDBOrderStatus()
	//if status == int(mysqld.StatusPay) {
	//	fmt.Println("Order_close 未支付订单，关闭 orderId =", orderId)
	//	dbOrder.SetDBOrderStatus(mysqld.StatusClose)
	//} else if status == int(mysqld.StatusRepute) || status == int(mysqld.StatusFinish) {
	//	fmt.Println("Order_close 订单完成，隐藏 orderId =", orderId)
	//	dbOrder.SetDBOrderStatus(mysqld.StatusHide)
	//} else {
	//	fmt.Println("Order_close can't close  orderId =", ctx.FormValue("orderId"))
	//	ctx.JSON(iris.Map{"code": -1, "msg": "很抱歉，当前订单不能关闭"})
	//	return
	//}
	ctx.JSON(iris.Map{"code": 0})
}

//=>/order/delete true post {orderId,token} 
func Order_delete(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	orderId := common.AtoI(ctx.FormValue("orderId"))
	orderBox :=	player.GetOrderBox()
	dbOrder := orderBox.GetOrder(orderId)
	if dbOrder == nil {
		fmt.Println("Order_delete orderId =", ctx.FormValue("orderId"))
		ctx.JSON(iris.Map{"code": -1, "msg": "很抱歉，订单不存在"})
		return
	}
	if !dbOrder.DeleteOrder() {
		fmt.Println("Order_delete orderId =", ctx.FormValue("orderId"))
		ctx.JSON(iris.Map{"code": -1, "msg": "很抱歉，订单不允许被删除"})
		return
	}
	ctx.JSON(iris.Map{"code": 0})
}

//=>/order/hx true post {hxNumber} 
func Order_hx(ctx iris.Context, sess *common.BSSession) {
	empty("/order/hx")
}

//=>/order/statistics true get {token} 
func Order_statistics(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	data := map[string]interface{} {
		"countPay":2,
		"countTransfer": 5,
		"countConfirm": 1,
		"countRepute":3,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/order/refund true get {token,orderId} 
func Order_refund(ctx iris.Context, sess *common.BSSession) {
	empty("/order/refund")
}

//=>/order/refundApply/info true get {token,orderId} 
func Order_refundApply_info(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	orderId := common.AtoI(ctx.FormValue("orderId"))
	orderBox :=	player.GetOrderBox()
	dbOrder := orderBox.GetOrder(orderId)
	if dbOrder == nil {
		fmt.Println("Order_refundApply_info orderId =", ctx.FormValue("orderId"))
		ctx.JSON(iris.Map{"code": -1, "msg": "订单orderId出错"})
		return
	}
	refundData := map[string]interface{}{
		"orderId":		dbOrder.ID,
		"orderStatus": 	dbOrder.Status,
		"refundStatus":	dbOrder.RefundStatus,
		"amount": 		dbOrder.AmountReal,
	}
	refund := dbOrder.GetDBRefund()
	if refund != nil {
		refundData["refundId"] = refund.ID
	}
	data := map[string]interface{}{
		"refund": refundData,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/order/refundApply/apply true post {} 
func Order_refundApply_apply(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	orderId := common.AtoI(ctx.FormValue("orderId"))
	orderBox :=	player.GetOrderBox()
	dbOrder := orderBox.GetOrder(orderId)
	if dbOrder == nil {
		fmt.Println("Order_detail orderId =", ctx.FormValue("orderId"))
		ctx.JSON(iris.Map{"code": -1, "msg": "订单orderId出错"})
		return
	}
	refundType := common.AtoI(ctx.FormValue("refundType"))
	logistics := common.AtoI(ctx.FormValue("logistics"))
	reasonId := common.AtoI(ctx.FormValue("reasonId"))
	remark := ctx.FormValue("remark")
	pics := ctx.FormValue("pics")

	dbRefund := &mysqld.DBRefund {
		OrderId:    orderId,
		OrderStatus:dbOrder.Status,
		RefundType: refundType,
		Logistics:  logistics,
		ReasonId:   reasonId,
		AmountTotal: dbOrder.AmountPayerPay,
		Remark:     remark,
		Pics:       pics,
	}
	dbRefund.AmountRefund = dbRefund.AmountTotal
	dbOrder.AddDBRefund(dbRefund)

	if dbOrder.Status == int(mysqld.EStatusSend) {
		outRefundNo := dbRefund.GetRefundNumber()
		trade_no := dbOrder.GetPayOrderNumber()
		params := map[string]interface{}{
			"transactionId": dbOrder.TransactionId,
			"outTradeNo" : trade_no,
			"outRefundNo" : outRefundNo,
			"refund" : dbRefund.AmountRefund,
			"total" : dbRefund.AmountTotal,
		}
		reply := map[string]interface{}{}
		err := sdk.OnWeChatRefundOrder(&params, &reply)
		if err != nil {
			fmt.Println("Order_refundApply_apply orderId =", orderId, "AmountRefund =", dbRefund.AmountRefund, "err:", err)
			ctx.JSON(iris.Map{"code": -1, "msg":"调用微信退款发生错误"})
			return
		}
		_, ok := reply["code"]
		if ok {
			orderBytes, err := sdk.OnWeChatQueryRefundByMCH(outRefundNo)
			if err != nil {
				fmt.Println("Order_refundApply_apply orderId =", orderId, "amountReal =", dbOrder.AmountReal, "err:", err)
				ctx.JSON(iris.Map{"code": -1, "msg":"查询微信退款发生错误，请联系客服"})
				return
			}
			fmt.Println("Order_refundApply_apply orderBytes =", string(orderBytes))
			ctx.JSON(iris.Map{"code": 1, "msg": reply["message"]})
			return
		} else {
			dbOrder.SetDBOrderStatus(mysqld.EStatusRefundFinish)
			dbOrder.Save()
			if outRefundNo == reply["out_refund_no"].(string) {
				dbRefund.TransactionId = reply["transaction_id"].(string)
				dbRefund.Save()
			} else {
				fmt.Println("Order_refundApply_apply outRefundNo =", outRefundNo, "out_refund_no =", reply["out_refund_no"].(string))
			}
		}
	}
	ctx.JSON(iris.Map{"code": 0})
}

//=>/order/refundApply/cancel true post {token,orderId} 
func Order_refundApply_cancel(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	orderId := common.AtoI(ctx.FormValue("orderId"))
	orderBox :=	player.GetOrderBox()
	dbOrder := orderBox.GetOrder(orderId)
	if dbOrder == nil {
		fmt.Println("Order_detail orderId =", ctx.FormValue("orderId"))
		ctx.JSON(iris.Map{"code": -1, "msg": "订单orderId出错"})
		return
	}
	dbOrder.RemoveDBRefund()

	ctx.JSON(iris.Map{"code": 0})
}