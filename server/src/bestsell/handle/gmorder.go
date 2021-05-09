package handle

import (
	"bestsell/common"
	"bestsell/module"
	"bestsell/mysqld"
	"bestsell/sdk"
	"fmt"
	"github.com/kataras/iris/v12"
	"strconv"
	"time"
)

//=>/gm/order/list true get {} 
func Gm_order_list(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	var dbOrders *[]*mysqld.DBOrder
	status := common.AtoI(ctx.FormValue("status"))
	switch mysqld.EStatusType(status) {
		case mysqld.EStatusPay:
			dbOrders = mysqld.GetOrdersByStatus(status)
		case mysqld.EStatusSend:
			dbOrders = mysqld.GetOrdersByStatus(status)
		case mysqld.EStatusReceive:
			dbOrders = mysqld.GetOrdersByStatus(status)
		default:
			fmt.Println("Gm_order_list 不支持 status =", status)
			ctx.JSON(iris.Map{"code": -1})
			return
	}
	var orderList []*map[string]interface{}
	goodsMap := make(map[int][]map[string]interface{})
	for _,dbOrder := range *dbOrders {
		dbOrder.CheckStatus()
		if dbOrder.IsOrderDelete() {
			continue
		}
		dbOrder.CalAmountReal()
		timeStamp := time.Unix(int64(dbOrder.TimeStamp), 0)
		item := map[string]interface{}{
			"orderId"		:dbOrder.ID,
			"playerId"		:dbOrder.PlayerId,
			"orderNumber"	:dbOrder.GetOrderNumber(),
			"status"		:dbOrder.Status,
			"refundId"		:dbOrder.RefundId,
			"refundStatus"	:dbOrder.RefundStatus,
			"remark"		:dbOrder.Remark,
			"statusStr"		:dbOrder.GetDBOrderStatusName(),
			"timeStamp"		:dbOrder.TimeStamp,
			"couponId"		:dbOrder.CouponId,
			"amountCoupon"	:dbOrder.AmountCoupon,
			"amountGoods"	:dbOrder.AmountGoods,
			"amountReal"	:dbOrder.AmountReal,
			"dateAdd"		:timeStamp.Format("2006-01-02 15:04:05"),
		}

		orderList = append(orderList, &item)
		goodsNumber := 0
		goodsInfos := dbOrder.GetOrderGoodsInfos()
		var goodsList []map[string]interface{}
		for _, goodsInfo := range goodsInfos {
			goods := map[string]interface{}{
				"goodsId": goodsInfo.GoodsId,
				"name":    goodsInfo.Name,
				"number":  goodsInfo.Number,
				"price":   goodsInfo.Price,
				"pic":     goodsInfo.GetPublicPic(),
				"skuId":   goodsInfo.SkuId,
				"skuName": goodsInfo.SkuName,
				"skuPrice":goodsInfo.SkuPrice,
				"amount":  goodsInfo.Amount,
			}
			goodsList = append(goodsList, goods)
			goodsNumber += goodsInfo.Number
		}
		item["goodsNumber"] = goodsNumber
		goodsMap[dbOrder.ID] = goodsList
	}
	data := map[string]interface{}{
		"orderList": orderList,
		"goodsMap": goodsMap,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/gm/order/do true get {orderId, status}
func Gm_order_do(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	orderId := common.AtoI(ctx.FormValue("orderId"))
	status := common.AtoI(ctx.FormValue("status"))
	playerId := common.AtoI(ctx.FormValue("playerId"))
	player = mysqld.GetDBPlayerByPlayerIdOrFromDB(playerId)
	if player == nil {
		fmt.Println("Gm_order_do player == nil playerId =", playerId)
		ctx.JSON(iris.Map{"code": -1, "msg":"未找到该买家:"+strconv.Itoa(playerId)})
		return
	}

	orderBox :=	player.GetOrderBox()
	dbOrder := orderBox.GetOrder(orderId)
	if dbOrder == nil {
		fmt.Println("Gm_do_order orderId =", ctx.FormValue("orderId"))
		ctx.JSON(iris.Map{"code": -1, "msg": "订单orderId出错"})
		return
	}
	switch mysqld.EStatusType(status) {
	case mysqld.EStatusReceive:
		if dbOrder.GetDBOrderStatus() != int(mysqld.EStatusSend) {
			ctx.JSON(iris.Map{"code": -1, "msg": "订单不是待发货状态"})
			return
		}
		dbOrder.SendOrder()
	default:
		fmt.Println("Gm_do_order 不支持 status =", status)
		ctx.JSON(iris.Map{"code": -1, "msg":"商品不在待发货状态"})
		return
	}
	ctx.JSON(iris.Map{"code": 0, "msg":"发货成功"})
}

//=>/gm/order/detail true get {orderId, playerId}
func Gm_order_detail(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	orderId := common.AtoI(ctx.FormValue("orderId"))
	playerId := common.AtoI(ctx.FormValue("playerId"))
	player = mysqld.GetDBPlayerByPlayerIdOrFromDB(playerId)
	if player == nil {
		fmt.Println("Gm_order_do player == nil playerId =", playerId)
		ctx.JSON(iris.Map{"code": -1, "msg":"未找到该买家:"+strconv.Itoa(playerId)})
		return
	}
	orderBox :=	player.GetOrderBox()
	dbOrder := orderBox.GetOrder(orderId)
	if dbOrder == nil {
		fmt.Println("Gm_do_order orderId =", ctx.FormValue("orderId"))
		ctx.JSON(iris.Map{"code": -1, "msg": "订单orderId出错"})
		return
	}
	data := dbOrder.GetAllDetail()
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/gm/order/coupon true post {orderId, playerId, amount, couponId}
func Gm_order_coupon(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	orderId := common.AtoI(ctx.FormValue("orderId"))
	playerId := common.AtoI(ctx.FormValue("playerId"))
	player = mysqld.GetDBPlayerByPlayerIdOrFromDB(playerId)
	if player == nil {
		fmt.Println("Gm_order_coupon player == nil playerId =", ctx.FormValue("playerId"))
		ctx.JSON(iris.Map{"code": -1, "msg":"未找到买家:"+strconv.Itoa(playerId)})
		return
	}
	orderBox :=	player.GetOrderBox()
	dbOrder := orderBox.GetOrder(orderId)
	if dbOrder == nil {
		fmt.Println("Gm_order_coupon orderId =", ctx.FormValue("orderId"))
		ctx.JSON(iris.Map{"code": -1, "msg": "订单orderId出错"})
		return
	}
	couponId := common.AtoI(ctx.FormValue("couponId"))
	if couponId > 0 && couponId != dbOrder.CouponId {
		fmt.Println("Gm_order_coupon couponId == dbOrder.CouponId  couponId =", couponId, "dbOrder.CouponId =", dbOrder.CouponId)
		ctx.JSON(iris.Map{"code": -1, "msg": "订单优惠券重复创建"})
		return
	}
	amount := common.AtoI(ctx.FormValue("amount"))
	if amount < 0 {
		fmt.Println("Gm_order_coupon amount < 0")
		ctx.JSON(iris.Map{"code": -1, "msg":"订单优惠券金额出错"})
		return
	}

	dbMyCounpon := dbOrder.GetOrCreateMyCounpon()
	if dbMyCounpon == nil {
		ctx.JSON(iris.Map{"code": -1, "msg": "创建优惠券发生错误"})
		return
	}

	dbMyCounpon.BeginWrite()
	dbMyCounpon.Amount = amount
	dbMyCounpon.EndWrite()
	dbMyCounpon.Save()

	dbOrder.SetAmountCoupon(dbMyCounpon.Amount)
	dbOrder.Save()

	ctx.JSON(iris.Map{"code": 0})
}

//=>/gm/refund/confirm true get {orderId}
func Gm_refund_confirm(ctx iris.Context, sess *common.BSSession) {
	player := module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	orderId := common.AtoI(ctx.FormValue("orderId"))
	playerId := common.AtoI(ctx.FormValue("playerId"))
	player = mysqld.GetDBPlayerByPlayerIdOrFromDB(playerId)
	if player == nil {
		fmt.Println("Gm_refund_confirm player == nil playerId =", playerId)
		ctx.JSON(iris.Map{"code": -1, "msg":"未找到该买家:"+strconv.Itoa(playerId)})
		return
	}

	orderBox :=	player.GetOrderBox()
	dbOrder := orderBox.GetOrder(orderId)
	if dbOrder == nil {
		fmt.Println("Gm_refund_confirm orderId =", ctx.FormValue("orderId"))
		ctx.JSON(iris.Map{"code": -1, "msg": "订单orderId出错"})
		return
	}
	dbRefund := dbOrder.GetDBRefund()

	outRefundNo := dbRefund.GetRefundNumber()
	params := map[string]interface{}{
		"transactionId": dbOrder.TransactionId,
		"outTradeNo" : dbOrder.OrderNumber,
		"outRefundNo" : outRefundNo,
		"refund" : dbRefund.AmountRefund,
		"total" : dbRefund.AmountTotal,
	}
	reply := map[string]interface{}{}
	err := sdk.OnWeChatRefundOrder(&params, &reply)
	if err != nil {
		fmt.Println("Gm_refund_confirm orderId =", orderId, "AmountRefund =", dbRefund.AmountRefund, "err:", err)
		ctx.JSON(iris.Map{"code": -1, "msg":"调用微信退款发生错误"})
		return
	}
	dbOrder.SetDBOrderStatus(mysqld.EStatusRefundFinish)
	dbOrder.Save()

	_, ok := reply["code"]
	if ok {
		orderBytes, err := sdk.OnWeChatQueryRefundByMCH(outRefundNo)
		if err != nil {
			fmt.Println("Gm_refund_confirm orderId =", orderId, "amountReal =", dbOrder.AmountReal, "err:", err)
			ctx.JSON(iris.Map{"code": -1, "msg":"查询微信退款发生错误，请联系客服"})
			return
		}
		fmt.Println("Gm_refund_confirm orderBytes =", string(orderBytes))
		ctx.JSON(iris.Map{"code": 1, "msg": reply["message"]})
		return
	} else {
		if outRefundNo == reply["out_refund_no"].(string) {
			dbRefund.TransactionId = reply["transaction_id"].(string)
			dbRefund.Save()
		} else {
			fmt.Println("Gm_refund_confirm outRefundNo =", outRefundNo, "out_refund_no =", reply["out_refund_no"].(string))
		}
	}

	ctx.JSON(iris.Map{"code": 0})
}

//=>/gm/refund/cancel true get {orderId}
func Gm_refund_cancel(ctx iris.Context, sess *common.BSSession) {
	player := module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	orderId := common.AtoI(ctx.FormValue("orderId"))
	playerId := common.AtoI(ctx.FormValue("playerId"))
	player = mysqld.GetDBPlayerByPlayerIdOrFromDB(playerId)
	if player == nil {
		fmt.Println("Gm_refund_cancel player == nil playerId =", playerId)
		ctx.JSON(iris.Map{"code": -1, "msg":"未找到该买家:" + ctx.FormValue("playerId")})
		return
	}

	orderBox :=	player.GetOrderBox()
	dbOrder := orderBox.GetOrder(orderId)
	if dbOrder == nil {
		fmt.Println("Gm_refund_cancel orderId =", ctx.FormValue("orderId"))
		ctx.JSON(iris.Map{"code": -1, "msg": "订单orderId出错"})
		return
	}

	dbOrder.RemoveDBRefund()
	ctx.JSON(iris.Map{"code": 0})
}

