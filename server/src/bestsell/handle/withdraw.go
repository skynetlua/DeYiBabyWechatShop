package handle
import (
	"bestsell/module"
	"bestsell/mysqld"
	"github.com/kataras/iris/v12"
	"bestsell/common"
)

//=>/withdraw/apply true post {money} 
func Withdraw_apply(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg": "未登陆"})
		return
	}
	money := common.AtoI(ctx.FormValue("money"))
	if player.Balance < float64(money) {
		ctx.JSON(iris.Map{"code": -1, "msg": "账户余额不足"})
		return
	}
	box := player.GetWithdrawLogBox()


	item := &mysqld.DBWithdrawLog{
		PlayerId: player.ID,
		Status: int(mysqld.WithdrawLogCheck),
		Money: money,
	}
	box.AddItem(item)
	ctx.JSON(iris.Map{"code": 0, "msg": "提现申请成功，等审核"})
}

//=>/withdraw/detail true get {id} 
func Withdraw_detail(ctx iris.Context, sess *common.BSSession) {



}

//=>/withdraw/list true post {} 
func Withdraw_list(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	box := player.GetWithdrawLogBox()
	dbItems := box.GetItems()
	var datas []*map[string]interface{}
	for idx,item := range *dbItems {
		data := map[string]interface{}{
			"index":  idx,
			"money": item.Money,
			"status":item.Status,
			"statusStr":item.GetWithdrawTypeName(),
			"dateAdd":item.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		datas = append(datas, &data)
	}
	ctx.JSON(iris.Map{"code": 0, "data": datas})
}
