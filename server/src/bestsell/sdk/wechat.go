package sdk

import (
	"bestsell/common"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	WeChatAppid     = "XXXXXXXXXXX"
	WeChatSecret    = "5XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	WeChatmchid     = "XXXXXXXXXXX"
)

var wechatAccessToken string
var wechatAccessTokenExpire time.Time
var writeCounter int32


type WechatPayInfo struct {
	MchId string 			`json:"mchid"`
	AppId string 			`json:"appid"`
	OutTradeNo string 		`json:"out_trade_no"`
	TransactionId string 	`json:"transaction_id"`
	TradeType string 		`json:"trade_type"`
	TradeState string 		`json:"trade_state"`
	TradeStateDesc string 	`json:"trade_state_desc"`
	BankType string 		`json:"bank_type"`
	Attach string 			`json:"attach"`
	SuccessTime string 		`json:"success_time"`
	Payer map[string]string `json:"payer"`
	Amount map[string]interface{} `json:"amount"`
}

func init() {
}

func OnWeChatParsePayInfo(payInFoMsg string) (*WechatPayInfo, error) {
	var payInfo WechatPayInfo
	err = json.Unmarshal([]byte(payInFoMsg), &payInfo)
	if err != nil {
		return nil, err
	}
	if payInfo.MchId != WeChatmchid {
		fmt.Println("OnWeChatParsePayInfo payInfo.MchId =", payInfo.MchId, "WeChatmchid =", WeChatmchid)
		return nil, fmt.Errorf("OnWeChatParsePayInfo payInfo.MchId != WeChatmchid")
	}
	if payInfo.AppId != WeChatAppid {
		fmt.Println("OnWeChatParsePayInfo payInfo.AppId =", payInfo.AppId, "WeChatAppid =", WeChatAppid)
		return nil, fmt.Errorf("OnWeChatParsePayInfo payInfo.AppId != WeChatAppid")
	}
	return &payInfo, nil
}

func OnWeChatDecryptMsg(associatedData, nonce, ciphertext string) (*WechatPayInfo, error) {
	msg, err := DecryptWeChatMsg(associatedData, nonce, ciphertext)
	if err != nil {
		return nil, err
	}
	return OnWeChatParsePayInfo(msg)
}

func OnWeChatPayOrder(params *map[string]interface{}, reply *map[string]interface{}) error {
	openId := (*params)["openId"].(string)
	trade_no := (*params)["trade_no"].(string)
	amount := (*params)["amount"].(int)
	description := (*params)["description"].(string)
	goods_tag := (*params)["goods_tag"].(string)
	attach := (*params)["attach"].(string)
	ret, err := OnWeChatCreateOrder(openId, trade_no, amount, description, goods_tag, attach)
	if err != nil {
		return err
	}
	code := fastjson.GetString(ret, "code")
	if len(code) > 0 {
		(*reply)["code"] = code
		(*reply)["message"] = fastjson.GetString(ret, "message")
		return nil
	}
	prepay_id := fastjson.GetString(ret, "prepay_id")
	(*params)["prepay_id"] = prepay_id
	(*params)["appId"] = WeChatAppid
	return SignatureWeChatPay(params, reply)
}

func OnWeChatRefundOrder(params *map[string]interface{}, reply *map[string]interface{}) error {
	transactionId := (*params)["transactionId"].(string)
	outTradeNo := (*params)["outTradeNo"].(string)
	outRefundNo := (*params)["outRefundNo"].(string)
	refund := (*params)["refund"].(int)
	total := (*params)["total"].(int)
	ret, err := _OnWeChatRefundOrder(transactionId, outTradeNo, outRefundNo, refund, total)
	if err != nil {
		return err
	}
	code := fastjson.GetString(ret, "code")
	if len(code) > 0 {
		(*reply)["code"] = code
		(*reply)["message"] = fastjson.GetString(ret, "message")
		return nil
	}
	(*params)["transaction_id"] = fastjson.GetString(ret, "transaction_id")
	(*params)["out_refund_no"] = fastjson.GetString(ret, "out_refund_no")
	return nil
}

func OnWeChatCreateOrder(openId string, trade_no string, amount int, description string, goods_tag string, attach string) ([]byte, error) {
	total := amount
	createOrderUrl := "https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi"
	reqdata := map[string]interface{}{
		"appid": WeChatAppid,
		"mchid": WeChatmchid,
		"description": description,
		"out_trade_no": trade_no,
		"attach": attach,
		"notify_url": "https://www.bestsellmall.com:444/api/wechat/pay_notify",
		"goods_tag": goods_tag,
		"amount" : map[string]interface{}{
			"total": total,
			"currency": "CNY",
		},
		"payer": map[string]interface{}{
			"openid": openId,
		},
	}
	return PostWeChatPay(createOrderUrl, &reqdata)
}

func OnWeChatQueryOrderByMCH(tradeNo string) ([]byte, error) {
	queryOrderUrl := "https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/"
	requestUrl := queryOrderUrl+tradeNo+"?mchid="+WeChatmchid
	return GetWeChatPay(requestUrl)
}

func OnWeChatQueryRefundByMCH(tradeNo string) ([]byte, error) {
	queryOrderUrl := "https://api.mch.weixin.qq.com/v3/refund/domestic/refunds/"
	requestUrl := queryOrderUrl+tradeNo
	return GetWeChatPay(requestUrl)
}

func _OnWeChatRefundOrder(transactionId string, outTradeNo string, outRefundNo string, refund int, total int) ([]byte, error) {
	cancelOrderUrl := "https://api.mch.weixin.qq.com/v3/refund/domestic/refunds"
	reqdata := map[string]interface{} {
		"transaction_id": transactionId,
		"out_trade_no": outTradeNo,
		"out_refund_no": outRefundNo,
		"reason": "商品已售完",
		"notify_url": "https://www.bestsellmall.com:444/api/wechat/cancel_notify",
		"amount" : map[string]interface{} {
			"refund": refund,
			"total": total,
			"currency": "CNY",
		},
	}
	return PostWeChatPay(cancelOrderUrl, &reqdata)
}

func OnWeChatCloseOrder(tradeNo string) ([]byte, error) {
	queryOrderUrl := "https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/"
	requestUrl := queryOrderUrl+tradeNo+"/close"
	reqdata := map[string]interface{}{
		"mchid": WeChatmchid,
	}
	return PostWeChatPay(requestUrl, &reqdata)
}


func getWeChatAccessToken() string {
	common.BeginWrite(&writeCounter)
	accessToken := wechatAccessToken
	common.EndWrite(&writeCounter)
	if len(accessToken) > 0 {
		curTime := time.Now()
		outlineTime := wechatAccessTokenExpire
		curUnix := curTime.Unix()
		outlineUnix := outlineTime.Unix()
		if curUnix < outlineUnix {
			return accessToken
		}
	}
	reply := make(map[string]interface{})
	if OnWeChatAccessToken(&reply) != 0 {
		panic("OnWeChatAccessToken error")
		return ""
	}
	access_token := reply["access_token"].(string)
	expires_in := reply["expires_in"].(int)-100
	outlineTime := time.Now().Add(time.Duration(expires_in)*time.Second)
	fmt.Println("outlineTime =", outlineTime)

	common.BeginWrite(&writeCounter)
	wechatAccessToken = access_token
	wechatAccessTokenExpire = outlineTime
	common.EndWrite(&writeCounter)

	return wechatAccessToken
}

func OnWeChatLogin(code string, reply *map[string]interface{}) int {
	_url := "https://api.weixin.qq.com/sns/jscode2session?appid="+WeChatAppid+"&secret="+WeChatSecret+"&js_code="+code+"&grant_type=authorization_code"
	req, err := http.NewRequest("GET", _url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("status", resp.Status)
	fmt.Println("response:", resp.Header)
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		panic(err)
	}
	fmt.Println("OnWeChatLogin response =", string(response))

	errcode := fastjson.GetInt(response, "errcode")
	if errcode != 0 {
		errmsg := fastjson.GetString(response, "errmsg")
		fmt.Println("OnWeChatLogin error:", errcode, errmsg)
		return errcode
	}
	openid := fastjson.GetString(response, "openid")
	//OnWeChatCreateOrder(openid)
	session_key := fastjson.GetString(response, "session_key")
	unionid := fastjson.GetString(response, "unionid")
	(*reply)["openid"] = openid
	(*reply)["session_key"] = session_key
	(*reply)["unionid"] = unionid
	return 0
}

//access_token
func OnWeChatAccessToken(reply *map[string]interface{}) int {
	_url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid="+WeChatAppid+"&secret="+WeChatSecret
	req, err := http.NewRequest("GET", _url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("status", resp.Status)
	fmt.Println("response:", resp.Header)
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		panic(err)
	}
	fmt.Println("response =", string(response))
	errcode := fastjson.GetInt(response, "errcode")
	if errcode != 0 {
		errmsg := fastjson.GetString(response, "errmsg")
		fmt.Println("OnWeChatAccessToken error:", errcode, errmsg)
		return errcode
	}
	access_token := fastjson.GetString(response, "access_token")
	expires_in := fastjson.GetInt(response, "expires_in")
	(*reply)["access_token"] = access_token
	(*reply)["expires_in"] = expires_in
	return 0
}

func OnWeChatQcode(fileName string,params *map[string]interface{}, reply *map[string]interface{}) int {
	body, err := json.Marshal(&params)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	accessToken := getWeChatAccessToken()
	if len(accessToken) == 0 {
		return -1
	}
	_url := "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token="+accessToken
	fmt.Println("_url =", _url)
	fmt.Println("boyd =", string(body))

	req, err := http.NewRequest("POST", _url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("status", resp.Status)
	fmt.Println("response:", resp.Header)
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		panic(err)
	}
	ctype := resp.Header.Get("Content-Type")
	if strings.Compare(ctype, "image/jpeg") == 0 {
		ioutil.WriteFile(fileName, response, os.ModePerm)
		return 0
	}
	errcode := fastjson.GetInt(response, "errcode")
	if errcode != 0 {
		errmsg := fastjson.GetString(response, "errmsg")
		fmt.Println("OnWeChatQcode error:", errcode, errmsg)
		return errcode
	}
	return -1
}
