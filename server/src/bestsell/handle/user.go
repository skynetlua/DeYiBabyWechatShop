package handle

import (
	"bestsell/common"
	"bestsell/module"
	"fmt"
	"github.com/kataras/iris/v12"
)


//=>/user/check/token true get {token} 
func User_check_token(ctx iris.Context, sess *common.BSSession) {
	token := ctx.FormValue("token")
	fmt.Println("User_check_token token =", token)
	if len(token) == 0 {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	userLogin := &module.UserLogin{}
	userLogin.Token = token
	userLogin.Code = ""
	userLogin.Ip = ctx.RemoteAddr()
	ret := module.OnLogin(userLogin)
	if ret < 0 {
		ctx.JSON(iris.Map{"code": -1})
		return
	}else if ret != 0 {
		ctx.JSON(iris.Map{"code": ret})
		return
	}
	player := userLogin.Player
	data := map[string]interface{}{
		"token": player.Token,
		"uid": player.ID,
		"gm": player.GM,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/user/check/referrer true get {referrer} 
func User_check_referrer(ctx iris.Context, sess *common.BSSession) {
	empty("/user/check/referrer")
}

//=>/user/detail true get {token} 
func User_detail(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	playerInfo := player.GetPlayerInfo()
	player.GM = 1
	userInfo := map[string]interface{}{
		"id" 			:player.ID,
		"refererId" 	:player.RefererId,
		"refererName" 	:player.RefererName,
		"avatarUrl" 	:playerInfo.AvatarUrl,
		"nickName" 		:playerInfo.NickName,
		"mobile" 		:playerInfo.Mobile,
		"balance" 		:player.Balance,
		"amountCost" 	:player.AmountCost,
		"growth" 		:player.Growth,
		// "isSeller" 		:player.IsSeller(),
	}
	data := map[string]interface{}{
		"userInfo":userInfo,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/user/wxinfo true get {token} 
func User_wxinfo(ctx iris.Context, sess *common.BSSession) {
	empty("/user/wxinfo")
}

//=>/user/amount true get {token} 
func User_amount(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	data := map[string]interface{}{
		"balance": player.Balance,
		"amountCost": player.AmountCost,
		"growth": player.Growth,
		"score": player.Score,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}
//=>/user/cashLog true post {} 
func User_cashLog(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	dbCashLogBox := player.GetCashLogBox()
	dbCashLogs := dbCashLogBox.GetCashLogs()
	var items []*map[string]interface{}
	for idx,cashLog := range *dbCashLogs {
		item := map[string]interface{}{
			"index":  idx,
			"amount": cashLog.Amount,
			"behavior":cashLog.Behavior,
			"cashType":cashLog.CashType,
			"typeStr":cashLog.GetCashTypeName(),
			"dateAdd":cashLog.GetUpdateDateStr(),
		}
		items = append(items, &item)
	}
	ctx.JSON(iris.Map{"code": 0, "data": items})
}

//=>/user/payLog true post {} 
func User_payLog(ctx iris.Context, sess *common.BSSession) {
	empty("/user/payLog")
}
