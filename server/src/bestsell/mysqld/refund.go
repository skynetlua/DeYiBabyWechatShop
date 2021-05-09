package mysqld

import (
	"strconv"
	"time"
)

type DBRefund struct {
    MysqlModel
	PlayerId  		int `json:"playerId"`
	OrderId 		int `json:"orderId"`
	Status 			int `json:"status"`
	OrderStatus 	int `json:"orderStatus"`
	RefundType 		int `json:"refundType"`
	Logistics 		int `json:"logistics"`
	ReasonId   		int `json:"reasonId"`
	AmountTotal 	int `json:"amountTotal"`
	AmountRefund 	int `json:"amountRefund"`
	TimeStamp  		int `json:"timeStamp"`
	TransactionId   string
	Remark 			string `json:"remark"`
	Pics 			string `json:"pics"`
}

//var dbRefundSafeMap common.SafeMap

//func GetDBRefund(orderId int)*DBRefund {
//	ret := dbRefundSafeMap.Get(orderId)
//	if ret == nil {
//		return nil
//	}
//	return ret.(*DBRefund)
//}

//func AddDBRefund(item *DBRefund)  {
//	old := GetDBRefund(item.OrderId)
//	if old != nil {
//		fmt.Println("DBRefund repeat")
//		return
//	}
//	dbRefundSafeMap.Set(item.OrderId, item)
//}

//func LoadDBRefunds()  {
//	var _dbItemsSlice []*DBRefund
//	db.Find(&_dbItemsSlice)
//	dbRefundSafeMap = *common.NewSafeMap()
//	for _,item := range _dbItemsSlice{
//		dbRefundSafeMap.Set(item.ID, item)
//	}
//}

func startDBRefund()  {
	if !db.HasTable(&DBRefund{}) {
		db.CreateTable(&DBRefund{})
	}
	//dbRefundSafeMap = *common.NewSafeMap()
	//LoadDBRefunds()
}

//func AddNewRefund(refund *DBRefund)  {
//	refund.Insert()
//}
//
//func GetDBRefundOrFromDB(orderId int)*DBRefund {
//	dbRefund := GetDBRefund(orderId)
//	if dbRefund != nil  {
//		return dbRefund
//	}
//	dbRefund = &DBRefund{
//		OrderId:orderId,
//	}
//	dbRefund.LoadWithOrderId()
//	if dbRefund.ID == 0 {
//		return nil
//	}
//	AddDBRefund(dbRefund)
//	return dbRefund
//}

//DBRefund
func (p *DBRefund)Insert()  {
	if p.ID < 0 {
		panic("(p *DBRefund)Insert p.ID < 0")
	}
	p.TimeStamp = int(time.Now().Unix())
	db.Create(p)
}

func (p *DBRefund)Load(){
	if p.ID < 0 {
		panic("(p *DBRefund)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBRefund)LoadWithOrderId(){
	if p.OrderId < 0 {
		panic("(p *DBRefund)LoadWithOrderId p.OrderId < 0")
	}
	db.Where("`order_id` = ?", p.OrderId).First(p)
}

func (p *DBRefund)Save(){
	if p.ID < 0 {
		panic("(p *DBRefund)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBRefund)Remove(){
	db.Delete(p)
}

func (p *DBRefund)GetRefundNumber() string {
	tm := time.Unix(int64(p.TimeStamp),0)
	orderNumber := tm.Format("20060102150405")
	orderNumber = strconv.Itoa(p.PlayerId)+orderNumber[2:]+"*"+strconv.Itoa(p.ID)
	return orderNumber
}