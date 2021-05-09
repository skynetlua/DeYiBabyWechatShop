package mysqld

import (
	"fmt"
	"bestsell/common"
)

type DBCoupon struct {
    MysqlModel
	Name string `json:"name"`
}

var dbCouponSafeMap common.SafeMap

func GetDBCoupon(id int)*DBCoupon {
	ret := dbCouponSafeMap.Get(id)
	if ret == nil {
		return nil
	}
	return ret.(*DBCoupon)
}

func AddDBCoupon(item *DBCoupon)  {
	old := GetDBCoupon(item.ID)
	if old != nil {
		fmt.Println("DBCoupon repeat")
		return
	}
	dbCouponSafeMap.Set(item.ID, item)
}

func LoadDBCoupons()  {
	var _dbItemsSlice []*DBCoupon
	db.Find(&_dbItemsSlice)
	dbCouponSafeMap = *common.NewSafeMap()
	for _,item := range _dbItemsSlice{
		dbCouponSafeMap.Set(item.ID, item)
	}
}

func startDBCoupon()  {
	if !db.HasTable(&DBCoupon{}) {
		db.CreateTable(&DBCoupon{})
	}
	// LoadDBCoupons()
}

//DBCoupon
func (p *DBCoupon)Insert()  {
	if p.ID < 0 {
		panic("(p *DBCoupon)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBCoupon)Load(){
	if p.ID < 0 {
		panic("(p *DBCoupon)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBCoupon)Save(){
	if p.ID < 0 {
		panic("(p *DBCoupon)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBCoupon)Remove(){
	db.Delete(p)
}
