package handle

import (
	"github.com/kataras/iris/v12"
	"bestsell/common"
)

type CfgCoupon struct {
	Id int `json:"id"`
	Name string `json:"name"`
	DateEndDays int `json:"dateEndDays"`
	DateEndType int `json:"dateEndType"`
	DateStartType int `json:"dateStartType"`
	MoneyHreshold int `json:"moneyHreshold"`
	MoneyMax int `json:"moneyMax"`
	MoneyMin int `json:"moneyMin"`
	NeedScore int `json:"needScore"`
	NeedSignedContinuous int `json:"needSignedContinuous"`
	NumberGit int `json:"numberGit"`
	NumberGitNumber int `json:"numberGitNumber"`
	NumberLeft int `json:"numberLeft"`
	NumberPersonMax int `json:"numberPersonMax"`
	NumberTotle int `json:"numberTotle"`
	NumberUsed int `json:"numberUsed"`
	Status int `json:"status"`
	StatusStr string `json:"statusStr"`
}
type CfgCouponSlice []CfgCoupon

var cfgCouponSlice = []CfgCoupon{
	CfgCoupon{
		Id: 5453,
		Name:"一人一份",
		DateEndDays : 7,
		DateEndType : 1,
		DateStartType :1,
		MoneyHreshold : 0.00,
		MoneyMax : 1.00,
		MoneyMin : 1.00,
		NeedScore : 0,
		NeedSignedContinuous : 0,
		NumberGit : 505,
		NumberGitNumber : 505,
		NumberLeft : 999999494,
		NumberPersonMax : 1,
		NumberTotle : 999999999,
		NumberUsed : 1,
		Status : 0,
		StatusStr : "正常",
	},
	CfgCoupon{
		Id: 223,
		Name: "新店优惠",
		DateEndDays: 15,
		DateEndType: 1,
		DateStartType: 1,
		MoneyHreshold: 3000.00,
		MoneyMax: 40.00,
		MoneyMin: 40.00,
		NeedScore: 0,
		NeedSignedContinuous: 0,
		NumberGit: 12689,
		NumberGitNumber: 6825,
		NumberLeft: 993322,
		NumberPersonMax: 999999,
		NumberTotle: 995699,
		NumberUsed: 64,
		Status: 0,
		StatusStr: "正常",
	},
	CfgCoupon{
		Id: 222,
		Name: "新店优惠",
		DateEndDays: 15,
		DateEndType: 1,
		DateStartType: 1,
		MoneyHreshold: 2000.00,
		MoneyMax: 25.00,
		MoneyMin: 25.00,
		NeedScore: 0,
		NeedSignedContinuous: 0,
		NumberGit: 10588,
		NumberGitNumber: 5484,
		NumberLeft: 998464,
		NumberPersonMax: 999999,
		NumberTotle: 999999,
		NumberUsed: 64,
		Status: 0,
		StatusStr: "正常",
	},
}

//=>/discount/coupon true get {} 
func Discount_coupon(ctx iris.Context, sess *common.BSSession) {
	ctx.JSON(iris.Map{"code": 0, "data": cfgCouponSlice})
}
//=>/discount/detail true get {id} 
func Discount_detail(ctx iris.Context, sess *common.BSSession) {
	empty("/discount/detail")
}

//=>/discount/my true get {} 
func Discount_my(ctx iris.Context, sess *common.BSSession) {
	empty("/discount/my")
}

//=>/discount/fetch true post {} 
func Discount_fetch(ctx iris.Context, sess *common.BSSession) {
	empty("/discount/fetch")
}

//=>/discount/send true post {} 
func Discount_send(ctx iris.Context, sess *common.BSSession) {
	empty("/discount/send")
}

//=>/discount/exchange true post {token,number,pwd} 
func Discount_exchange(ctx iris.Context, sess *common.BSSession) {
	empty("/discount/exchange")
}