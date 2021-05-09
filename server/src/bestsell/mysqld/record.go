package mysqld

import (
	"fmt"
	"bestsell/common"
)

type DBRecord struct {
    MysqlModel
	Name string `json:"name"`
}

var dbRecordSafeMap common.SafeMap

func GetDBRecord(id int)*DBRecord {
	ret := dbRecordSafeMap.Get(id)
	if ret == nil {
		return nil
	}
	return ret.(*DBRecord)
}

func AddDBRecord(item *DBRecord)  {
	old := GetDBRecord(item.ID)
	if old != nil {
		fmt.Println("DBRecord repeat")
		return
	}
	dbRecordSafeMap.Set(item.ID, item)
}

func LoadDBRecords()  {
	var _dbItemsSlice []*DBRecord
	db.Find(&_dbItemsSlice)
	dbRecordSafeMap = *common.NewSafeMap()
	for _,item := range _dbItemsSlice{
		dbRecordSafeMap.Set(item.ID, item)
	}
}

func startDBRecord()  {
	if !db.HasTable(&DBRecord{}) {
		db.CreateTable(&DBRecord{})
	}
	LoadDBRecords()
}

//DBRecord
func (p *DBRecord)Insert()  {
	if p.ID < 0 {
		panic("(p *DBRecord)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBRecord)Load(){
	if p.ID < 0 {
		panic("(p *DBRecord)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBRecord)Save(){
	if p.ID < 0 {
		panic("(p *DBRecord)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBRecord)Remove(){
	db.Delete(p)
}
