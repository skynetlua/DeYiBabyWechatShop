package mysqld

import (
	"fmt"
	"bestsell/common"
)

type DBLogistics struct {
    MysqlModel
	Name string
}

var dbLogisticsSafeMap common.SafeMap

func GetDBLogistics(id int)*DBLogistics {
	ret := dbLogisticsSafeMap.Get(id)
	if ret == nil {
		return nil
	}
	return ret.(*DBLogistics)
}

func AddDBLogistics(item *DBLogistics)  {
	old := GetDBLogistics(item.ID)
	if old != nil {
		fmt.Println("DBLogistics repeat")
		return
	}
	dbLogisticsSafeMap.Set(item.ID, item)
}

func LoadDBLogisticss()  {
	var _dbItemsSlice []*DBLogistics
	db.Find(&_dbItemsSlice)
	dbLogisticsSafeMap = *common.NewSafeMap()
	for _,item := range _dbItemsSlice{
		dbLogisticsSafeMap.Set(item.ID, item)
	}
}

func startDBLogistics()  {
	if !db.HasTable(&DBLogistics{}) {
		db.CreateTable(&DBLogistics{})
	}
	LoadDBLogisticss()
}

//DBLogistics
func (p *DBLogistics)Insert()  {
	if p.ID < 0 {
		panic("(p *DBLogistics)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBLogistics)Load(){
	if p.ID < 0 {
		panic("(p *DBLogistics)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBLogistics)Save(){
	if p.ID < 0 {
		panic("(p *DBLogistics)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBLogistics)Remove(){
	db.Delete(p)
}
