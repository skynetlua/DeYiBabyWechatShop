package router

import (
	"bestsell/common"
	"bestsell/handle"
	// "fmt"
	"github.com/kataras/iris/v12"
	"path"
	"strings"
)

//needSubDomain, method, data, token
var RoutesParams = [][5]string{
	{"/subdomain/appid/wxapp","false","get","{appid}"},
	//{"/page/data","true","get", "{ path }"},
	//{"/config/vipLevel","true","get"},
	{"/page/index","true","get"},
	{"/page/goods/detail","true","get"},

	//wechat
	{"/wechat/login","true","post","{code,type}"},
	{"/wechat/register/complex","true","post","{}"},
	{"/wechat/register/simple","true","post","{}"},
	{"/wechat/pay","true","post","{}"},
	{"/wechat/qrcode","true","post","{}"},
	{"/wechat/bindMobile","true","post","{encryptedData,iv}"},
	{"/wechat/pay_notify","true","post","{}"},
	{"/wechat/cancel_notify","true","post","{}"},

	//config
	{"/config/value","true","get","{key}"},
	{"/config/values","true","get","{keys}"},

	//user
	{"/user/check/token","true","get","{token}"},
	{"/user/check/referrer","true","get","{referrer}"},
	{"/user/detail","true","get","{token}"},
	{"/user/wxinfo","true","get","{token}"},
	{"/user/amount","true","get","{token}"},
	{"/user/cashLog","true","post","{}"},
	{"/user/payLog","true","post","{}"},

	//address
	{"/address/add","true","post","{}"},
	{"/address/update","true","post","{}"},
	{"/address/delete","true","post","{id,token}"},
	{"/address/list","true","get","{token}"},
	{"/address/default","true","get","{token}"},
	{"/address/detail","true","get","{id,token}"},

	{"/address/update","true","post","{}"},
	{"/address/update","true","post","{}"},

	//banner
	{"/banner/list","true","get","{}"},

	//goods
	{"/goods/category/all","true","get"},
	{"/goods/category/info","true","get","{id}"},
	{"/goods/list","true","post","{}"},
	{"/goods/detail","true","get","{id}"},
	{"/goods/sku","true","get","{id}"},
	{"/goods/price","true","post","{goodsId,propertyChildIds}"},
	{"/goods/reputation","true","post","{}"},
	{"/goods/category/subtypes","true","get","{categoryId}"},
	{"/goods/category/sublist","true","get","{categoryId, subType}"},

	//fav
	{"/fav/list","true","post","{}"},
	{"/fav/add","true","post","{token,goodsId}"},
	{"/fav/check","true","get","{token,goodsId}"},
	{"/fav/delete","true","post","{token,goodsId}"},

	//subshop
	{"/subshop/list","true","post","{}"},
	{"/subshop/my","true","get","{token}"},
	{"/subshop/detail","true","get","{id}"},
	{"/subshop/apply","true","post","{}"},

	//cart
	{"/cart/info","true","get","{token}"},
	{"/cart/list","true","get","{token}"},
	{"/cart/add","true","post","{token,goodsId,skuId,buyNumber}"},
	{"/cart/modifyNumber","true","post","{token,key,number}"},
	{"/cart/remove","true","post","{token,key}"},
	{"/cart/empty","true","post","{token}"},
	{"/cart/quick","true","get","{token,goodsId, skuId, buyNumber}"},

	//notice
	{"/notice/list","true","post","{}"},
	{"/notice/lastone","true","get","{type}"},
	{"/notice/detail","true","get","{id}"},

	//discount
	{"/discount/coupon","true","get","{}"},
	{"/discount/detail","true","get","{id}"},
	{"/discount/my","true","get","{}"},
	{"/discount/fetch","true","post","{}"},
	{"/discount/send","true","post","{}"},
	{"/discount/exchange","true","post","{token,number,pwd}"},

	//live
	{"/live/rooms","true","get"},
	{"/live/his","true","get","{roomId}"},

	//order
	{"/order/prepare","true","post","{}"},
	{"/order/create","true","post","{}"},
	{"/order/list","true","post","{}"},
	{"/order/detail","true","get","{id,token,hxNumber}"},
	{"/order/delivery","true","post","{orderId,token}"},
	{"/order/reputation","true","post","{}"},
	{"/order/close","true","post","{orderId,token}"},
	{"/order/delete","true","post","{orderId,token}"},
	{"/order/pay","true","post","{orderId,token}"},
	{"/order/hx","true","post","{hxNumber}"},
	{"/order/statistics","true","get","{token}"},
	{"/order/refund","true","get","{token,orderId}"},
	{"/order/refundApply/apply","true","post","{}"},
	{"/order/refundApply/info","true","get","{token,orderId}"},
	{"/order/refundApply/cancel","true","post","{token,orderId}"},

	//distribute
	{"/distribute/info","true","get","{}"},
	{"/distribute/apply","true","post","{name,mobile}"},
	{"/distribute/apply/progress","true","get","{}"},
	{"/distribute/members","true","post","{}"},
	{"/distribute/log","true","post","{}"},

	//common
	{"/region/province","false","get"},
	{"/region/child","false","get","{pid}"},

	//upload
	{"/upload/file","true","post","{}"},

	//withdraw
	{"/withdraw/apply","true","post","{money}"},
	{"/withdraw/detail","true","get","{id}"},
	{"/withdraw/list","true","post","{}"},

	//gm
	{"/gm/order/list", "true", "get", "{}"},
	{"/gm/order/do", "true", "get", "{orderId, status, playerId}"},
	{"/gm/order/detail", "true", "get", "{orderId, playerId}"},
	{"/gm/order/coupon", "true", "post", "{orderId, playerId, amount, couponId}"},
	{"/gm/refund/confirm", "true", "get", "{orderId, playerId}"},
	{"/gm/refund/cancel", "true", "get", "{orderId, playerId}"},
	{"/gm/team/list", "true", "get", "{}"},
	{"/gm/do/team", "true", "get", "{teamId, status}"},
	{"/gm/goods/info", "true", "get", "{goodsId}"},
	{"/gm/goods/update", "true", "post", "{}"},
	{"/gm/goods/list", "true", "get", "{status, page, pageSize}"},
	{"/gm/goods/barcode", "true", "get", "{barCode}"},
	{"/gm/upload/goods","true","post","{}"},
	{"/gm/goods/update/info", "true", "post", "{}"},
	{"/gm/goods/remove", "true", "get", "{goodsId}"},
	{"/gm/goods/category", "true", "get", "{categoryId}"},
	{"/gm/goods/goodsdata", "true", "get", "{barCode}"},
	{"/gm/goods/goodsdatas", "true", "get", "{}"},

	{"/gm/category/list", "true", "get", "{}"},
	{"/gm/category/update", "true", "post", "{}"},
	{"/gm/category/remove", "true", "get", "{}"},
	{"/gm/upload/category", "true","post","{}"},
	//{"/gm/upload/excel","true","post","{}"},
	//{"/gm/goods/load/picture","true","get","{}"},


	//{"/common/mobile-segment/location","false","get","{mobile}"},
	//{"/common/mobile-segment/next","false","post","{}"},

	//{"/score/send/rule","true","post","{}"},
	//{"/score/sign/rules","true","get","{}"},
	//{"/score/sign","true","post","{token}"},
	//{"/score/sign/logs","true","post","{}"},
	//{"/score/today-signed","true","get","{token}"},
	//{"/score/exchange","true","post","{number,token}"},
	//{"/score/exchange/cash","true","post","{deductionScore,token}"},
	//{"/score/logs","true","post","{}"},
	//{"/score/share/wxa/group","true","post","{code,referrer,encryptedData,iv}"},

	//{"/template-msg/wxa/formId","true","post","{token,type,formId}"},
	//{"/template-msg/put","true","post","{}"},
	//{"/pay/tt/microapp","true","post","{}"},
	//{"/pay/query","true","get","{token,outTradeId}"},
	//{"/pay/lcsw/wxapp","true","post","{}"},
	//{"/pay/wepayez/wxapp","true","post","{}"},
	//{"/pay/alipay/semiAutomatic/payurl","true","post","{}"},
	//{"/user/wxapp/login/mobile","true","post","{code,encryptedData,iv}"},
	//{"/user/username/login","true","post","{}"},
	//{"/user/username/bindUsername","true","post","{token,username,pwd}"},
	//{"/user/m/login","true","post","{mobile,pwd,deviceId,deviceName}"},
	//{"/user/m/reset-pwd","true","post","{mobile,pwd,code}"},
	//{"/user/email/reset-pwd","true","post","{email,pwd,code}"},
	//
	//{"/user/username/register","true","post","{}"},
	//{"/user/m/register","true","post","{}"},
	//


	//{"/friendly-partner/list","true","post","{type}"},
	//{"/user/friend/list","true","post","{}"},
	//{"/user/friend/add","true","post","{token,uid}"},
	//{"/user/friend/detail","true","get","{token,uid}"},
	//{"/media/video/detail","true","get","{videoId}"},
	//{"/user/m/bind-mobile","true","post","{token,mobile,code,pwd}"},



	//{"/user/recharge/send/rule","true","get"},
	//{"/payBill/discounts","true","get"},
	//{"/payBill/pay","true","post","{token,money}"},

	//{"/dfs/upload/url","true","post","{remoteFileUrl,ext}"},
	//{"/dfs/upload/list","true","post","{path}"},

	//{"/cms/category/list","true","get","{}"},
	//{"/cms/category/info","true","get","{id}"},
	//{"/cms/news/list","true","post","{}"},
	//{"/cms/news/useful/logs","true","post","{}"},
	//{"/cms/news/detail","true","get","{id}"},
	//{"/cms/news/preNext","true","get","{id}"},
	//{"/cms/news/put","true","post","{}"},
	//{"/cms/news/del","true","post","{token,id}"},
	//{"/cms/news/useful","true","post","{}"},
	//{"/cms/page/info/v2","true","get","{key}"},
	//{"/cms/tags/list","true","get","{}"},
	//{"/invoice/list","true","post","{}"},
	//{"/invoice/apply","true","post","{}"},
	//{"/invoice/info","true","get","{token,id}"},
	//{"/deposit/list","true","post","{}"},
	//{"/deposit/pay","true","post","{}"},
	//{"/deposit/info","true","get","{token,id}"},
	//{"/deposit/back/apply","true","post","{token,id}"},

	//{"/comment/add","true","post","{}"},
	//{"/comment/list","true","post","{}"},
	//{"/user/modify","true","post","{}"},
	//{"/uniqueId/get","true","get","{type}"},
	//{"/barcode/info","true","get","{barcode}"},
	//{"/luckyInfo/info/v2","true","get","{id}"},
	//{"/luckyInfo/join","true","post","{id,token}"},
	//{"/luckyInfo/join/my","true","get","{id,token}"},
	//{"/luckyInfo/join/logs","true","post","{}"},
	//{"/json/list","true","post","{}"},
	//{"/json/set","true","post","{}"},
	//{"/json/delete","true","post","{token,id}"},
	//{"/verification/pic/check","true","post","{key,code}"},
	//{"/common/short-url/shorten","false","post","{url}"},
	//{"/verification/sms/get","true","get","{mobile,key,picCode}"},
	//{"/verification/sms/check","true","post","{mobile,code}"},
	//{"/verification/mail/get","true","get","{mail}"},
	//{"/verification/mail/check","true","post","{mail,code}"},
	//{"/common/map/distance","false","get","{lat1,lng1,lat2,lng2}"},
	//{"/common/map/qq/address","false","get","{location,coord_type}"},
	//{"/common/map/qq/search","false","post","{}"},
	//{"/virtualTrader/list","true","post","{}"},
	//{"/virtualTrader/info","true","get","{token,id}"},
	//{"/virtualTrader/buy","true","post","{token,id}"},
	//{"/virtualTrader/buy/logs","true","post","{}"},
	//{"/queuing/types","true","get","{status}"},
	//{"/queuing/get","true","post","{token,typeId,mobile}"},
	//{"/queuing/my","true","get","{token,typeId,status}"},
	//{"/user/idcard","true","post","{token,name,idCardNo}"},
	//{"/user/loginout","true","get","{token}"},
	//{"/user/level/list","true","post","{}"},
	//{"/user/level/info","true","get","{id}"},
	//{"/user/level/prices","true","get","{levelId}"},
	//{"/user/level/buy","true","post","{token,userLevelPriceId,isAutoRenew,remark}"},
	//{"/user/level/buyLogs","true","post","{}"},
	//{"/user/message/list","true","post","{}"},
	//{"/user/message/read","true","post","{token,id}"},
	//{"/user/message/del","true","post","{token,id}"},
	//{"/user/wxapp/bindOpenid","true","post","{token,code,type}"},
	//{"/user/wxapp/decode/encryptedData","true","post","{code,encryptedData,iv}"},
	//{"/score/deduction/rules","true","get","{type}"},
	//{"/vote/items","true","post","{}"},
	//{"/vote/info","true","get","{id}"},
	//{"/vote/vote","true","post","{token,voteId,items}"},
	//{"/vote/vote/info","true","get","{token,voteId}"},
	//{"/vote/vote/list","true","post","{}"},

	//{"/user/email/register","true","post","{}"},
	//{"/user/email/login","true","post","{}"},
	//{"/user/email/bindUsername","true","post","{token,email,code,pwd}"},
	//{"/site/statistics","true","get"},
	//{"/cms/news/fav/add","true","post","{token,newsId}"},
	//{"/cms/news/fav/check","true","get","{token,newsId}"},
	//{"/cms/news/fav/list","true","post","{}"},
	//{"/cms/news/fav/delete","true","post","{token,id}"},
	//{"/cms/news/fav/delete","true","post","{token,newsId}"},
	//
	//{"/growth/logs","true","post","{}"},
	//{"/growth/exchange","true","post","{token,deductionScore}"},
}

var noTokenList = []string{
	"/subdomain/appid/wxapp",
}

type HandleArg struct {
	handle Handler
	path string
	hasSubdomain bool
	fullPath string
	method string
	args string
	hasToken bool
	session *common.BSSession
}

var handleMap map[string]*HandleArg
var subdomain = "/api"

func initHandle()  {
	handleMap = make(map[string]*HandleArg)
	//{{
	ROUTE("/subdomain/appid/wxapp", handle.Subdomain_appid_wxapp)
	ROUTE("/page/index", handle.Page_index)
	ROUTE("/page/goods/detail", handle.Page_goods_detail)
	ROUTE("/wechat/login", handle.Wechat_login)
	ROUTE("/wechat/register/complex", handle.Wechat_register_complex)
	ROUTE("/wechat/register/simple", handle.Wechat_register_simple)
	ROUTE("/wechat/pay", handle.Wechat_pay)
	ROUTE("/wechat/qrcode", handle.Wechat_qrcode)
	ROUTE("/wechat/bindMobile", handle.Wechat_bindMobile)
	ROUTE("/wechat/pay_notify", handle.Wechat_pay_notify)
	ROUTE("/wechat/cancel_notify", handle.Wechat_cancel_notify)
	ROUTE("/config/value", handle.Config_value)
	ROUTE("/config/values", handle.Config_values)
	ROUTE("/user/check/token", handle.User_check_token)
	ROUTE("/user/check/referrer", handle.User_check_referrer)
	ROUTE("/user/detail", handle.User_detail)
	ROUTE("/user/wxinfo", handle.User_wxinfo)
	ROUTE("/user/amount", handle.User_amount)
	ROUTE("/user/cashLog", handle.User_cashLog)
	ROUTE("/user/payLog", handle.User_payLog)
	ROUTE("/address/add", handle.Address_add)
	ROUTE("/address/update", handle.Address_update)
	ROUTE("/address/delete", handle.Address_delete)
	ROUTE("/address/list", handle.Address_list)
	ROUTE("/address/default", handle.Address_default)
	ROUTE("/address/detail", handle.Address_detail)
	ROUTE("/banner/list", handle.Banner_list)
	ROUTE("/goods/category/all", handle.Goods_category_all)
	ROUTE("/goods/category/info", handle.Goods_category_info)
	ROUTE("/goods/list", handle.Goods_list)
	ROUTE("/goods/detail", handle.Goods_detail)
	ROUTE("/goods/sku", handle.Goods_sku)
	ROUTE("/goods/price", handle.Goods_price)
	ROUTE("/goods/reputation", handle.Goods_reputation)
	ROUTE("/goods/category/subtypes", handle.Goods_category_subtypes)
	ROUTE("/goods/category/sublist", handle.Goods_category_sublist)
	ROUTE("/fav/list", handle.Fav_list)
	ROUTE("/fav/add", handle.Fav_add)
	ROUTE("/fav/check", handle.Fav_check)
	ROUTE("/fav/delete", handle.Fav_delete)
	ROUTE("/subshop/list", handle.Subshop_list)
	ROUTE("/subshop/my", handle.Subshop_my)
	ROUTE("/subshop/detail", handle.Subshop_detail)
	ROUTE("/subshop/apply", handle.Subshop_apply)
	ROUTE("/cart/info", handle.Cart_info)
	ROUTE("/cart/list", handle.Cart_list)
	ROUTE("/cart/add", handle.Cart_add)
	ROUTE("/cart/modifyNumber", handle.Cart_modifyNumber)
	ROUTE("/cart/remove", handle.Cart_remove)
	ROUTE("/cart/empty", handle.Cart_empty)
	ROUTE("/cart/quick", handle.Cart_quick)
	ROUTE("/notice/list", handle.Notice_list)
	ROUTE("/notice/lastone", handle.Notice_lastone)
	ROUTE("/notice/detail", handle.Notice_detail)
	ROUTE("/discount/coupon", handle.Discount_coupon)
	ROUTE("/discount/detail", handle.Discount_detail)
	ROUTE("/discount/my", handle.Discount_my)
	ROUTE("/discount/fetch", handle.Discount_fetch)
	ROUTE("/discount/send", handle.Discount_send)
	ROUTE("/discount/exchange", handle.Discount_exchange)
	ROUTE("/live/rooms", handle.Live_rooms)
	ROUTE("/live/his", handle.Live_his)
	ROUTE("/order/prepare", handle.Order_prepare)
	ROUTE("/order/create", handle.Order_create)
	ROUTE("/order/list", handle.Order_list)
	ROUTE("/order/detail", handle.Order_detail)
	ROUTE("/order/delivery", handle.Order_delivery)
	ROUTE("/order/reputation", handle.Order_reputation)
	ROUTE("/order/close", handle.Order_close)
	ROUTE("/order/delete", handle.Order_delete)
	ROUTE("/order/pay", handle.Order_pay)
	ROUTE("/order/hx", handle.Order_hx)
	ROUTE("/order/statistics", handle.Order_statistics)
	ROUTE("/order/refund", handle.Order_refund)
	ROUTE("/order/refundApply/apply", handle.Order_refundApply_apply)
	ROUTE("/order/refundApply/info", handle.Order_refundApply_info)
	ROUTE("/order/refundApply/cancel", handle.Order_refundApply_cancel)
	ROUTE("/distribute/info", handle.Distribute_info)
	ROUTE("/distribute/apply", handle.Distribute_apply)
	ROUTE("/distribute/apply/progress", handle.Distribute_apply_progress)
	ROUTE("/distribute/members", handle.Distribute_members)
	ROUTE("/distribute/log", handle.Distribute_log)
	ROUTE("/region/province", handle.Region_province)
	ROUTE("/region/child", handle.Region_child)
	ROUTE("/upload/file", handle.Upload_file)
	ROUTE("/withdraw/apply", handle.Withdraw_apply)
	ROUTE("/withdraw/detail", handle.Withdraw_detail)
	ROUTE("/withdraw/list", handle.Withdraw_list)
	ROUTE("/gm/order/list", handle.Gm_order_list)
	ROUTE("/gm/order/do", handle.Gm_order_do)
	ROUTE("/gm/order/detail", handle.Gm_order_detail)
	ROUTE("/gm/order/coupon", handle.Gm_order_coupon)
	ROUTE("/gm/refund/confirm", handle.Gm_refund_confirm)
	ROUTE("/gm/refund/cancel", handle.Gm_refund_cancel)
	ROUTE("/gm/team/list", handle.Gm_team_list)
	ROUTE("/gm/do/team", handle.Gm_do_team)
	ROUTE("/gm/goods/info", handle.Gm_goods_info)
	ROUTE("/gm/goods/update", handle.Gm_goods_update)
	ROUTE("/gm/goods/list", handle.Gm_goods_list)
	ROUTE("/gm/goods/barcode", handle.Gm_goods_barcode)
	ROUTE("/gm/upload/goods", handle.Gm_upload_goods)
	ROUTE("/gm/goods/update/info", handle.Gm_goods_update_info)
	ROUTE("/gm/goods/remove", handle.Gm_goods_remove)
	ROUTE("/gm/goods/category", handle.Gm_goods_category)
	ROUTE("/gm/category/list", handle.Gm_category_list)
	ROUTE("/gm/category/update", handle.Gm_category_update)
	ROUTE("/gm/category/remove", handle.Gm_category_remove)
	ROUTE("/gm/upload/category", handle.Gm_upload_category)
	ROUTE("/gm/goods/goodsdata", handle.Gm_goods_goodsdata)
	ROUTE("/gm/goods/goodsdatas", handle.Gm_goods_goodsdatas)

	//ROUTE("/gm/upload/excel", handle.Gm_upload_excel)
	//ROUTE("/gm/goods/load/picture", handle.Gm_goods_load_picture)
	//}}
}

func ROUTE(_path string, handle func(iris.Context, *common.BSSession)) {
	handleArg := HandleArg{
		handle:handle,
		path:_path,
		hasToken:false,
	}
	//handleMap[_path] = &handleArg
	isNew := false
	for _,route := range RoutesParams{
		key := route[0]+"_"+route[1]+"_"+route[2]
		_,ok := handleMap[key]
		if ok {
			continue
		}
		isNew = true
		if strings.Compare(route[0], _path) == 0 {
			if strings.Compare(route[1], "true") == 0 {
				handleArg.fullPath = path.Join(subdomain, _path)
				handleArg.hasSubdomain = true
			}else{
				handleArg.fullPath = _path
				handleArg.hasSubdomain = false
			}
			handleArg.method = route[2]
			handleArg.args = route[3]
			handleMap[key] = &handleArg

			for _,item := range noTokenList{
				if strings.Compare(item, _path) == 0 {
					handleArg.hasToken = false
					break
				}
			}
			break
		}
	}
	if !isNew {
		return
	}
	if len(handleArg.fullPath) == 0 {
		panic("can't find the path:"+_path)
		return
	}
	HandleHandle(&handleArg)
}

//func Handle_common(ctx iris.Context, sess *sessions.Session) {
//	rPath := ctx.Path()
//	fmt.Println("Handle_common path:=", rPath)
//}

func getHandleArg(_path string)*HandleArg{
	return handleMap[_path]
}

func checkHandle()  {
	for _,route := range RoutesParams {
		rPath := route[0]
		_,ok := handleMap[rPath]
		if !ok {
			// fmt.Println("unimplement path:=", rPath)
		}
	}
}

func startMap() {
	if common.Config.Enable["map"] != 1 {
		return
	}
	initHandle()
	checkHandle()
}