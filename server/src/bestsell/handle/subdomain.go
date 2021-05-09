package handle

import (
	"bestsell/common"
	"github.com/kataras/iris/v12"
	"strings"
)

//=>/subdomain/appid/wxapp false get {appid} 
func Subdomain_appid_wxapp(ctx iris.Context, sess *common.BSSession) {
	appid := ctx.FormValue("appid")
	if len(appid) == 0 {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	if strings.Compare(appid, common.Config.Appid) != 0 {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	data := map[string]interface{}{
		"subdomain": "api",
		"host": "http://127.0.0.1:444",
		"vipLevel": 0,
		"config": []map[string]interface{} {
			{"key":"mallName", "value":"Q-Baby母婴生活馆"},
			// {"key":"recharge_amount_min", "value":0},
			// {"key":"WITHDRAW_MIN", "value":10},
			//{"key":"ALLOW_SELF_COLLECTION", "value":1},
			// {"key":"order_hx_uids", "value":""},
			// {"key":"subscribe_ids", "value":""},
		},
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}