package handle

import (
	"github.com/kataras/iris/v12"
	"bestsell/common"
	"strings"
)

//=>/config/vipLevel true get  
func Config_vipLevel(ctx iris.Context, sess *common.BSSession) {
	data := map[string]interface{}{
		"vipLevel":1,
		"config":[]map[string]interface{} {
			{"key":"mallName", "value":"大卖场"},
			{"key":"recharge_amount_min", "value":0},
			{"key":"WITHDRAW_MIN", "value":0},
			{"key":"ALLOW_SELF_COLLECTION", "value":0},
			{"key":"order_hx_uids", "value":"orderid123456"},
			{"key":"subscribe_ids", "value":"scrid123456"},
		},
	}
	//config.key, config.value
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/config/value true get {key} 
func Config_value(ctx iris.Context, sess *common.BSSession) {
	key := ctx.FormValue("key")
	ctx.JSON(iris.Map{"code": 0, "data": key})
}

//mallName
//recharge_amount_min
//WITHDRAW_MIN
//ALLOW_SELF_COLLECTION
//order_hx_uids
//subscribe_ids
//=>/config/values true get {keys} 
func Config_values(ctx iris.Context, sess *common.BSSession) {
	_keys := ctx.FormValue("keys")
	keys := strings.Split(_keys, ",")
	ctx.JSON(iris.Map{"code": 0, "data": keys})
}