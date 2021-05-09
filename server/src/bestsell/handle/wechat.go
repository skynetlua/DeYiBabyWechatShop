package handle

import (
	aesCbc "bestsell/aescbc"
	"bestsell/common"
	"bestsell/module"
	"bestsell/mysqld"
	"bestsell/sdk"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"os"
	"path"
	"strconv"
	"strings"
)

//=>/wechat/login true post {code,type} 
func Wechat_login(ctx iris.Context, sess *common.BSSession) {
	code := ctx.FormValue("code")
	_type := ctx.FormValue("type")
	fmt.Println("Wechat_login code =", code, "type =", _type)

	reply := make(map[string]interface{})
	if sdk.OnWeChatLogin(code, &reply) != 0 {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	fmt.Println("reply =", reply)
	openid := reply["openid"].(string)
	session_key := reply["session_key"].(string)
	unionid := reply["unionid"].(string)
	if len(openid) == 0 {
		ctx.JSON(iris.Map{"code": -1})
		return
	}

	userLogin := &module.UserLogin{
		Openid: openid,
		SessionKey: session_key,
		Unionid: unionid,
		Token: "",
		Ip: ctx.RemoteAddr(),
	}
	ret := module.OnLogin(userLogin)
	if ret < 0 {
		ctx.JSON(iris.Map{"code": ret})
		return
	}else if ret != 0 {
		ctx.JSON(iris.Map{"code": ret})
		return
	}
	player := userLogin.Player
	playerInfo := player.GetPlayerInfo()
	data := map[string]interface{}{
		"token" :player.Token,
		"uid" 	:player.ID,
		"nickName": playerInfo.NickName,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/wechat/register/simple true post {} 
func Wechat_register_simple(ctx iris.Context, sess *common.BSSession) {
	code := ctx.FormValue("code")
	referrer := ctx.FormValue("referrer")

	reply := make(map[string]interface{})
	if sdk.OnWeChatLogin(code, &reply) != 0 {
		fmt.Println(" Wechat_register_simple login wx error")
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	fmt.Println("reply =", reply)
	sessionKey := reply["session_key"].(string)
	openId := reply["openid"].(string)
	if len(sessionKey) == 0 {
		fmt.Println(" Wechat_register_simple no sessionKey")
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	var userInfo = module.WxUserInfo{}
	userInfo.Referrer = referrer
	userInfo.OpenId = openId
	if module.OnLoginWx(&userInfo) != 0 {
		fmt.Println("Wechat_register_simple OnLoginWx error")
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	ctx.JSON(iris.Map{"code": 0})
}

//=>/wechat/register/complex true post {} 
func Wechat_register_complex(ctx iris.Context, sess *common.BSSession) {
	code := ctx.FormValue("code")
	encryptedData := ctx.FormValue("encryptedData")
	iv := ctx.FormValue("iv")
	referrer := ctx.FormValue("referrer")

	reply := make(map[string]interface{})
	if sdk.OnWeChatLogin(code, &reply) != 0 {
		fmt.Println(" User_wxapp_register_complex login wx error")
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	fmt.Println("reply =", reply)
	sessionKey := reply["session_key"].(string)
	loginOpenid := reply["openid"].(string)
	if len(sessionKey) == 0 {
		fmt.Println(" User_wxapp_register_complex no sessionKey")
		ctx.JSON(iris.Map{"code": -1})
		return
	}

	_sessionKey, err1 := base64.StdEncoding.DecodeString(sessionKey)
	if err1 != nil {
		fmt.Println("User_wxapp_register_complex DecodeString sessionKey error", err1)
		ctx.JSON(iris.Map{"code": -1})
		return
	}

	_encryptedData, err2 := base64.StdEncoding.DecodeString(encryptedData)
	if err2 != nil {
		fmt.Println("User_wxapp_register_complex DecodeString encryptedData error", err2)
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	_iv, err3 := base64.StdEncoding.DecodeString(iv)
	if err3 != nil {
		fmt.Println("User_wxapp_register_complex DecodeString iv error", err3)
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	aes := aesCbc.NewAesCipher128(_sessionKey, _iv)
	userInfoJson := aes.Decrypt(_encryptedData)

	txtJson := string(userInfoJson)
	retOjb := *common.ForceParseJson(txtJson)

	var userInfo = module.WxUserInfo{}
	userInfo.OpenId = retOjb["openId"]
	userInfo.NickName = retOjb["nickName"]
	userInfo.Gender,_ = strconv.Atoi(retOjb["gender"])
	userInfo.Language = retOjb["language"]
	userInfo.City = retOjb["city"]
	userInfo.Province = retOjb["province"]
	userInfo.Country = retOjb["country"]
	userInfo.AvatarUrl = retOjb["avatarUrl"]
	userInfo.Watermark.Timestamp ,_ = strconv.Atoi(retOjb["timestamp"])
	userInfo.Watermark.Appid = retOjb["appid"]
	userInfo.Referrer = referrer

	if strings.Compare(common.Config.Appid, userInfo.Watermark.Appid) != 0 {
		fmt.Println("User_wxapp_register_complex Appid error native Appid =", common.Config.Appid)
		fmt.Println("User_wxapp_register_complex Appid error weixin Appid =", userInfo.Watermark.Appid)
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	if strings.Compare(loginOpenid, userInfo.OpenId) != 0 {
		fmt.Println("User_wxapp_register_complex OpenId error native OpenId =", loginOpenid)
		fmt.Println("User_wxapp_register_complex OpenId error weixin OpenId =", userInfo.OpenId)
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	if module.OnLoginWx(&userInfo) != 0 {
		fmt.Println("User_wxapp_register_complex OnLoginWx error")
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	ctx.JSON(iris.Map{"code": 0})
}


//Accept:[*/*]
//Cache-Control:[no-cache]
//Connection:[Keep-Alive]
//Content-Length:[971]
//Content-Type:[application/json]
//Pragma:[no-cache]
//User-Agent:[Mozilla/4.0]
//Wechatpay-Nonce:[hTFmm0PX5hr799VyNnL0OPmQa0eLfQP3]
//Wechatpay-Serial:[1FE2D6609C894CA9327E555405746DBD226CFF02]
//Wechatpay-Signature:[C8x/vZ3vE89lxK8Aaby7H3/+xN7I6jya2x5mz8PkVu5Yj1gxNk2IwRiL64K8mzaiewiugZClREOsy5uFd3Q1O9OMvYsjzno51QD81lHMCbKX9BGAZcRQz4+gOFQ0vKXCQpfKLsE0/acbzVmeP3gHMVXMU/L6py5WFfCLzwxDzWP86pJWR+GZzFk/e+Xwo/9qEgB+hhR0E/gmx6R2zCaA5AaQttlZvK31bztrwmhJxPPLHUjAzk9T96Ku5RWMEpM0YyePNLJjjzdYm6uIiHjC73WSqW3zMqAyseEJJQR5y1M+mKjfe2tIbPcbN2/cIuqzjQfvVsaFGlucDZUQwi6AeA==]
//Wechatpay-Timestamp:[1616121263]
//{
//"id":"a6ca741b-b256-526c-b4ec-9dd2ec6c23da",
//"create_time":"2021-03-19T10:34:23+08:00",
//"resource_type":"encrypt-resource",
//"event_type":"TRANSACTION.SUCCESS",
//"summary":"支付成功",
//"resource":{
//	"original_type":"transaction",
//	"algorithm":"AEAD_AES_256_GCM",
//	"ciphertext":"bdJx/Oybyj71EBTrq5Ier0QUacfzGynlDU92Xx4PwvoadvIL8Bfw7c+fdx6DLXpReW0VrI/HBpmUmOvyrIRFLwWB1hmAi8RBig1+oid0t+hzOtnQD1FYPJvcj07B48UCctzQ7uLXicypbvDBtXQtk5S8kxfsPA6/8SivHJnAZRJVJjzmG0TkH4lxUzZu/p4tPUf2b7fI1szbQVGe9eas1JUjckoU3M1zkpdlM+gFysM+kPv9DXgNg4rumWGoVHeM7wcEvakBvx/0ZSve1RflxnImD9kjVOiWhAcBY/vVkc3FDK7Y/12jEGUJF9YwbdkuxinmT7e5uLEGHW6WW9sBnw58dDCDnUKvgfKl3Tn6THCpVoHJPmAv26aGjTtaa14sgh0YzpTMGFVLrNGhuJ8sFPXwNiez+uLf5ElnkYd7WpQ5pMfYG20+w6AdmigbHRs45F7CWuXRPB2XbeFyg8T79bY/Ie8w3DtypzqcLve8fkUBv9hiN1D1YNWbb8EVuXSd85fu/tAL2ioKlbEy2Z4PP6H36EccdLkTksDCQcjYUrS8P8bcWmfDdcJ05T8aeMwTHFDRe1pSQfsBlBCAUZu7/DPlxwFCErmf+OdyXO+10W4v4EsT/eNQHdQMYG5hr9LtIw==",
//	"associated_data":"transaction",
//	"nonce":"478T95cxUjQP"
//	}
//}
type WechatPayNotify struct {
	Id string 					`json:"id"`
	CreateTime string 			`json:"create_time"`
	ResourceType string 		`json:"resource_type"`
	EventType string 			`json:"event_type"`
	Summary string 				`json:"summary"`
	Resource map[string]string 	`json:"resource"`
}

func init() {
//	testTxt := `{
//	"id":"a6ca741b-b256-526c-b4ec-9dd2ec6c23da",
//	"create_time":"2021-03-19T10:34:23+08:00",
//	"resource_type":"encrypt-resource",
//	"event_type":"TRANSACTION.SUCCESS",
//	"summary":"支付成功",
//	"resource":{
//		"original_type":"transaction",
//		"algorithm":"AEAD_AES_256_GCM",
//		"ciphertext":"bdJx/Oybyj71EBTrq5Ier0QUacfzGynlDU92Xx4PwvoadvIL8Bfw7c+fdx6DLXpReW0VrI/HBpmUmOvyrIRFLwWB1hmAi8RBig1+oid0t+hzOtnQD1FYPJvcj07B48UCctzQ7uLXicypbvDBtXQtk5S8kxfsPA6/8SivHJnAZRJVJjzmG0TkH4lxUzZu/p4tPUf2b7fI1szbQVGe9eas1JUjckoU3M1zkpdlM+gFysM+kPv9DXgNg4rumWGoVHeM7wcEvakBvx/0ZSve1RflxnImD9kjVOiWhAcBY/vVkc3FDK7Y/12jEGUJF9YwbdkuxinmT7e5uLEGHW6WW9sBnw58dDCDnUKvgfKl3Tn6THCpVoHJPmAv26aGjTtaa14sgh0YzpTMGFVLrNGhuJ8sFPXwNiez+uLf5ElnkYd7WpQ5pMfYG20+w6AdmigbHRs45F7CWuXRPB2XbeFyg8T79bY/Ie8w3DtypzqcLve8fkUBv9hiN1D1YNWbb8EVuXSd85fu/tAL2ioKlbEy2Z4PP6H36EccdLkTksDCQcjYUrS8P8bcWmfDdcJ05T8aeMwTHFDRe1pSQfsBlBCAUZu7/DPlxwFCErmf+OdyXO+10W4v4EsT/eNQHdQMYG5hr9LtIw==",
//		"associated_data":"transaction",
//		"nonce":"478T95cxUjQP"
//	}
//}`
	//var payNotify WechatPayNotify
	//err := json.Unmarshal([]byte(testTxt), &payNotify)
	//if err != nil {
	//	panic(err)
	//}
}

//=>/wechat/cancel_notify true post {}
func Wechat_cancel_notify(ctx iris.Context, sess *common.BSSession) {

	fmt.Println("Wechat_cancel_notify Header:", ctx.Request().Header)
	body, _ := ctx.GetBody()
	fmt.Println("Wechat_cancel_notify body:", string(body))

	ctx.JSON(iris.Map{"code": "SUCCESS", "message": "成功"})
}

//=>/wechat/pay_notify true post {}
func Wechat_pay_notify(ctx iris.Context, sess *common.BSSession) {
	fmt.Println("Wechat_pay_notify Header:", ctx.Request().Header)
	body, _ := ctx.GetBody()
	fmt.Println("Wechat_pay_notify body:", string(body))
	var payNotify WechatPayNotify
	err := json.Unmarshal(body, &payNotify)
	if err != nil {
		fmt.Println("Wechat_pay_notify pasrse error:", err)
		ctx.JSON(iris.Map{"code": "SUCCESS", "message": "成功"})
		return
	}
	if payNotify.EventType == "TRANSACTION.SUCCESS" {
		associatedData := payNotify.Resource["associated_data"]
		nonce := payNotify.Resource["nonce"]
		ciphertext := payNotify.Resource["ciphertext"]
		fmt.Println("Wechat_pay_notify decrypt==>>")
		payInfo, err := sdk.OnWeChatDecryptMsg(associatedData, nonce, ciphertext)
		if err != nil {
			fmt.Println("Wechat_pay_notify decrypt error:", err)
			ctx.JSON(iris.Map{"code": "SUCCESS", "message": "成功"})
			return
		}
		openId := payInfo.Payer["openid"]
		outTradeNo := payInfo.OutTradeNo
		dbOrder := mysqld.GetDBOrderByOpenIdAndOrderNumber(openId, outTradeNo)
		if dbOrder == nil {
			fmt.Println("Wechat_pay_notify dbOrder == nil")
			ctx.JSON(iris.Map{"code": "SUCCESS", "message": "成功"})
			return
		}
		amountPay := int(payInfo.Amount["total"].(float64))
		amountPayerPay := int(payInfo.Amount["payer_total"].(float64))
		transactionId := payInfo.TransactionId
		dbOrder.FinishPay(amountPay, amountPayerPay, transactionId)
		fmt.Println("Wechat_pay_notify success EventType:", payNotify.EventType)
	} else {
		fmt.Println("Wechat_pay_notify fail EventType:", payNotify.EventType)
	}
	ctx.JSON(iris.Map{"code": "SUCCESS", "message": "成功"})
}

//=>/wechat/pay true post {} 
func Wechat_pay(ctx iris.Context, sess *common.BSSession) {
	player := module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	money := common.AtoI(ctx.FormValue("money"))
	if money < 0 {
		fmt.Println("Wechat_pay money < 0 money =",money)
		ctx.JSON(iris.Map{"code": -1, "msg":"订单金额出错，请联系客服"})
		return
	}
	orderId := common.AtoI(ctx.FormValue("orderId"))
	orderBox :=	player.GetOrderBox()
	dbOrder := orderBox.GetOrder(orderId)
	if dbOrder == nil {
		fmt.Println("Wechat_pay orderId =", ctx.FormValue("orderId"))
		ctx.JSON(iris.Map{"code": -1, "msg":"订单不存在，请联系客服"})
		return
	}
	if dbOrder.GetDBOrderStatus() != int(mysqld.EStatusPay) {
		fmt.Println("Wechat_pay orderId =", orderId, "orderStatus=", dbOrder.GetDBOrderStatus())
		ctx.JSON(iris.Map{"code": -1, "msg":"请求支付失败，订单已关闭或者完成"})
		return
	}
	dbPlayerInfo := player.GetPlayerInfo()
	ret := dbOrder.CalGoodsAmount()
	if ret == -1 {
		ctx.JSON(iris.Map{"code": -1, "msg":"有商品已下架，请联系客服"})
		return
	} else if ret == -2 {
		ctx.JSON(iris.Map{"code": -1, "msg":"商品库存不够，请联系客服"})
		return
	}
	dbOrder.CalAmountReal()
	fmt.Println("Wechat_pay playerId =", player.ID, "token =", player.Token, "name =", dbPlayerInfo.NickName, "mobile =", dbPlayerInfo.Mobile)
	if dbOrder.AmountReal != money {
		ctx.JSON(iris.Map{"code": -1, "msg":"订单金额不一致，请联系客服"})
		return
	}
	if money == 0 {
		dbOrder.FinishPay(0, 0, "free")
		ctx.JSON(iris.Map{"code": 1, "msg": "免费订单，欢迎再次光临!"})
		return
	}
	trade_no := dbOrder.GetPayOrderNumber()
	fmt.Println("Wechat_pay orderId =", orderId, "amountReal =", dbOrder.AmountReal, "money =", money)
	params := map[string]interface{} {
		"openId": player.OpenId,
		"trade_no": trade_no,
		"amount": dbOrder.AmountReal,
		"description": "东城水岸（锦兴防疫站旁）",
		"goods_tag": dbOrder.GetOrderGoodsTag(),
		"attach": dbOrder.GetOrderAttach(),
		"scene_info": map[string]string {
			"payer_client_ip": ctx.RemoteAddr(),
		},
	}
	//params["amount"] = 0.01
	reply := map[string]interface{}{}
	err := sdk.OnWeChatPayOrder(&params, &reply)
	if err != nil {
		fmt.Println("Wechat_pay orderId =", orderId, "amountReal =", dbOrder.AmountReal, "err:", err)
		ctx.JSON(iris.Map{"code": -1, "msg":"调用微信支付发生错误"})
		return
	}
	_, ok := reply["code"]
	if ok {
		orderBytes, err := sdk.OnWeChatQueryOrderByMCH(trade_no)
		if err != nil {
			fmt.Println("Wechat_pay orderId =", orderId, "amountReal =", dbOrder.AmountReal, "err:", err)
			ctx.JSON(iris.Map{"code": -1, "msg":"查询微信订单发生错误，请联系客服"})
			return
		}
		payInfo, err := sdk.OnWeChatParsePayInfo(string(orderBytes))
		if err != nil {
			fmt.Println("Wechat_pay ParsePayInfo error:", err)
			ctx.JSON(iris.Map{"code": "FAIL", "msg": "查询微信订单发生错误，请联系客服"})
			return
		}
		if payInfo.TradeState == "NOTPAY" {
			orderBytes, err = sdk.OnWeChatCloseOrder(trade_no)
			if err != nil {
				fmt.Println("Wechat_pay  OnWeChatCloseOrder orderId =", orderId, "amountReal =", dbOrder.AmountReal, "err:", err)
				ctx.JSON(iris.Map{"code": -1, "msg":"查询微信订单发生错误，请联系客服"})
				return
			}
		}
		openId := payInfo.Payer["openid"]
		if len(openId) == 0 {
			ctx.JSON(iris.Map{"code":  -1, "msg": "微信支付出现错误"})
			return
		}
		outTradeNo := payInfo.OutTradeNo
		dbOrder := mysqld.GetDBOrderByOpenIdAndOrderNumber(openId, outTradeNo)
		if dbOrder == nil {
			fmt.Println("Wechat_pay dbOrder == nil")
			ctx.JSON(iris.Map{"code":  -1, "msg": "查询微信订单，服务器没订单数据，请联系客服"})
			return
		}
		amountPay := int(payInfo.Amount["total"].(float64))
		amountPayerPay := int(payInfo.Amount["payer_total"].(float64))
		transactionId := payInfo.TransactionId
		dbOrder.FinishPay(amountPay, amountPayerPay, transactionId)

		ctx.JSON(iris.Map{"code": 1, "msg": "订单已支付"})
		return
	}
	ctx.JSON(iris.Map{"code": 0, "data": reply})
}

//=>/wechat/qrcode true post {} 
func Wechat_qrcode(ctx iris.Context, sess *common.BSSession) {
	player := module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	scene := ctx.FormValue("scene")
	//page := ctx.FormValue("page")
	//scene = url.QueryEscape(scene)
	//is_hyaline := ctx.FormValue("is_hyaline")
	//auto_color := ctx.FormValue("autoColor")
	//expireHours := ctx.FormValue("expireHours")

	params := make(map[string]interface{})
	params["scene"] = scene
	//params["page"] = page
	//params["is_hyaline"] = is_hyaline
	//params["auto_color"] = auto_color
	idStr := strconv.Itoa(player.ID)
	fileDir := common.CreateTokenDir(common.QRCodePath, idStr)
	fileName := common.MakeMd5(scene) +".jpg"
	fileName = path.Join(fileDir, fileName)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		reply := make(map[string]interface{})
		if sdk.OnWeChatQcode(fileName,&params, &reply) != 0 {
			ctx.JSON(iris.Map{"code": -1, "msg": "生成二维码发生错误"})
			return
		}
	}
	resUrl := common.MakeResUrl(fileName)
	fmt.Println("OnWeChatQcode")
	//escapeUrl := url.QueryEscape(urlStr)
	//enEscapeUrl, _ := url.QueryUnescape(escapeUrl)
	ctx.JSON(iris.Map{"code": 0, "data": resUrl})
}

//=>/wechat/bindMobile true post {encryptedData,iv} 
func Wechat_bindMobile(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": 10002})
		return
	}
	encryptedData := ctx.FormValue("encryptedData")
	iv := ctx.FormValue("iv")

	sessionKey := player.SessionKey
	if len(sessionKey) == 0 {
		fmt.Println(" Wechat_bindMobile no sessionKey")
		ctx.JSON(iris.Map{"code": 10002})
		return
	}
	_sessionKey, err1 := base64.StdEncoding.DecodeString(sessionKey)
	if err1 != nil {
		fmt.Println("Wechat_bindMobile DecodeString sessionKey error", err1)
		ctx.JSON(iris.Map{"code": 10002})
		return
	}
	_encryptedData, err2 := base64.StdEncoding.DecodeString(encryptedData)
	if err2 != nil {
		fmt.Println("Wechat_bindMobile DecodeString encryptedData error", err2)
		ctx.JSON(iris.Map{"code": 10002})
		return
	}
	_iv, err3 := base64.StdEncoding.DecodeString(iv)
	if err3 != nil {
		fmt.Println("Wechat_bindMobile DecodeString iv error", err3)
		ctx.JSON(iris.Map{"code": 10002})
		return
	}
	aes := aesCbc.NewAesCipher128(_sessionKey, _iv)
	byteJson := aes.Decrypt(_encryptedData)
	retOjb := *common.ForceParseJson(string(byteJson))
	if len(retOjb["phoneNumber"]) > 0 {
		playerInfo := player.GetPlayerInfo()
		if strings.Compare(playerInfo.Mobile, retOjb["phoneNumber"]) == 0 {
			fmt.Println("Wechat_bindMobile phoneNumber no change")
		}else{
			playerInfo.Mobile = retOjb["phoneNumber"]
			playerInfo.DelaySave(playerInfo)
		}
	}else{
		fmt.Println("Wechat_bindMobile phoneNumber empty")
	}
	ctx.JSON(iris.Map{"code": 0})
}

