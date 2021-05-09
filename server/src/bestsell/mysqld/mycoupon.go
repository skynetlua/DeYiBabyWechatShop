package mysqld

import (
	"fmt"
	"bestsell/common"
)

type DBMyCoupon struct {
    MysqlModel
	Name 		string
	CouponType  int
	PlayerId  	int
	OrderId  	int
	Amount  	int //0.01元
}

var dbMyCouponSafeMap common.SafeMap

func GetDBMyCoupon(id int)*DBMyCoupon {
	ret := dbMyCouponSafeMap.Get(id)
	if ret == nil {
		return nil
	}
	return ret.(*DBMyCoupon)
}

func AddDBMyCoupon(item *DBMyCoupon)  {
	old := GetDBMyCoupon(item.ID)
	if old != nil {
		fmt.Println("DBMyCoupon repeat")
		return
	}
	dbMyCouponSafeMap.Set(item.ID, item)
}

func LoadDBMyCoupons()  {
	var _dbItemsSlice []*DBMyCoupon
	db.Find(&_dbItemsSlice)
	dbMyCouponSafeMap = *common.NewSafeMap()
	for _,item := range _dbItemsSlice{
		dbMyCouponSafeMap.Set(item.ID, item)
	}
}

func startDBMyCoupon()  {
	if !db.HasTable(&DBMyCoupon{}) {
		db.CreateTable(&DBMyCoupon{})
	}
	// LoadDBMyCoupons()
}

//DBMyCoupon
func (p *DBMyCoupon)Insert()  {
	if p.ID < 0 {
		panic("(p *DBMyCoupon)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBMyCoupon)Load(){
	if p.ID < 0 {
		panic("(p *DBMyCoupon)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBMyCoupon)LoadWithOrderId(){
	if p.OrderId <= 0 {
		panic("(p *DBMyCoupon)LoadWithOrderId p.OrderId < 0")
	}
	db.Where("order_id = ?", p.OrderId).First(p)
}

func (p *DBMyCoupon)Save(){
	if p.ID < 0 {
		panic("(p *DBMyCoupon)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBMyCoupon)Remove(){
	db.Delete(p)
}
