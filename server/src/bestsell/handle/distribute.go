package handle

import (
	"bestsell/common"
	"bestsell/module"
	// "bestsell/mysqld"
	"github.com/kataras/iris/v12"
)

//=>/distribute/apply true post {name,mobile} 
func Distribute_apply(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg": "未登陆"})
		return
	}
	// nickname := ctx.FormValue("name")
	// mobile := ctx.FormValue("mobile")
	// myTeam := player.GetMyTeam()

	// myTeam.BeginWrite()
	// myTeam.Nickname = nickname
	// myTeam.Mobile = mobile
	// myTeam.Status = int(mysqld.TeamStatusCheck)
	// myTeam.EndWrite()

	// if myTeam.ID <= 0 {
	// 	myTeam.Insert()
	// }else {
	// 	myTeam.DelaySave(myTeam)
	// }
	ctx.JSON(iris.Map{"code": 0, "msg": "申请成功等审核"})
}

//=>/distribute/apply/progress true get {} 
func Distribute_apply_progress(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": 2000, "msg": "未登陆"})
		return
	}
	// myTeam := player.GetMyTeam()
	// data := map[string]interface{}{
	// 	"status":myTeam.GetTeamStatus(),
	// 	"remark":myTeam.GetTeamStatusName(),
	// }
	// ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/distribute/members true post {} 
func Distribute_members(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg": "未登陆"})
		return
	}
	//memberBox := player.GetMemberBox()
	//members := memberBox.GetMembers()
	//if len(*members) == 0 {
	//	ctx.JSON(iris.Map{"code": 700})
	//	return
	//}
	//var datas []interface{}
	//for _,item :=range *members {
	//	data := map[string]interface{}{
	//		"memberId" 	:item.MemberId,
	//		"level"    	:item.Level,
	//		"nickname" 	:item.Nickname,
	//		"mobile" 	:item.Mobile,
	//	}
	//	datas = append(datas, data)
	//}
	//ctx.JSON(iris.Map{"code": 0, "data": datas})
}

//=>/distribute/info true get {} 
func Distribute_info(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg": "未登陆"})
		return
	}
	// myTeam := player.GetMyTeam()
	// mobile := myTeam.Mobile
	// if len(mobile) == 0 {
	// 	mobile = player.GetPlayerInfo().Mobile
	// }
	// data := map[string]interface{}{
	// 	"refererId": player.RefererId,
	// 	"refererName": player.RefererName,
	// 	"nickname":myTeam.Nickname,
	// 	"mobile": mobile,
	// }
	// ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/distribute/log true post {} 
func Distribute_log(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg": "未登陆"})
		return
	}
	commissionBox := player.GetCommissionBox()
	commissions := commissionBox.GetCommissions()
	var datas []interface{}
	for _,item :=range *commissions {
		data := map[string]interface{}{
			"level"    	:item.Level,
			"money" 	:item.Money,
			"ratio" 	:item.Ratio,
			"sellerId" 	:item.SellerId,
			"sellerName" :item.SellerName,
			"buyerId" 	:item.BuyerId,
			"buyerName" 	:item.BuyerName,
		}
		datas = append(datas, data)
	}
	ctx.JSON(iris.Map{"code": 0, "data": datas})
}
