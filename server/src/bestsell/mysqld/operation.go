package mysqld

import (
	"fmt"
	"bestsell/common"
)

type DBOperation struct {
    MysqlModel
	PlayerId  	int
}

var dbOperationBoxSafeMap common.SafeMap

func GetDBOperationBox(playerId int)*DBOperationBox {
	ret := dbOperationBoxSafeMap.Get(playerId)
	if ret == nil {
		return nil
	}
	return ret.(*DBOperationBox)
}

func AddDBOperationBox(item *DBOperationBox)  {
	old := GetDBOperationBox(item.PlayerId)
	if old != nil {
		fmt.Println("DBOperationBox repeat")
		return
	}
	dbOperationBoxSafeMap.Set(item.PlayerId, item)
}

func GetDBOperationBoxOrFromDB(playerId int)*DBOperationBox {
	dbOperationBox := GetDBOperationBox(playerId)
	if dbOperationBox != nil  {
		return dbOperationBox
	}
	dbOperationBox = GetDBOperationBoxFromDB(playerId)
	return dbOperationBox
}

func LoadDBOperations()  {
	var _dbItemsSlice []*DBOperation
	db.Find(&_dbItemsSlice)
	dbOperationBoxSafeMap = *common.NewSafeMap()
	for _,item := range _dbItemsSlice{
		dbOperationBoxSafeMap.Set(item.ID, item)
	}
}

func startDBOperation()  {
	if !db.HasTable(&DBOperation{}) {
		db.CreateTable(&DBOperation{})
	}
	dbOperationBoxSafeMap = *common.NewSafeMap()
}

func AddNewOperation(item *DBOperation)  {
	box := GetDBOperationBox(item.PlayerId)
	if box != nil {
		box.AddOperation(item)
		return
	}
	item.Insert()
}


//DBOperationBox
type DBOperationBox struct {
	MysqlModelBox
	PlayerId int
	items []*DBOperation
}
func (p *DBOperationBox)AddOperation(item *DBOperation) {
	item.PlayerId = p.PlayerId
	item.Insert()
	p.BeginWrite()
	p.items = append(p.items, item)
	p.EndWrite()
}
func (p *DBOperationBox)GetOperation(id int)*DBOperation {
	for _,item := range p.items {
		if item.ID == id {
			return item
		}
	}
	return nil
}
func (p *DBOperationBox)GetOperations()*[]*DBOperation {
	return &p.items
}
func (p *DBOperationBox)RemoveOperation(id int) {
	for idx,item := range p.items {
		if item.ID == id {
			p.BeginWrite()
			p.items = append(p.items[:idx], p.items[idx+1:]...)
			p.EndWrite()
			item.Remove()
			return
		}
	}
}
func GetDBOperationBoxFromDB(playerId int)*DBOperationBox  {
	box := DBOperationBox{
		PlayerId:playerId,
	}
	db.Where("player_id = ?", playerId).Find(&box.items)
	return &box
}


//DBOperation
func (p *DBOperation)Insert()  {
	if p.ID < 0 {
		panic("(p *DBOperation)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBOperation)Load(){
	if p.ID < 0 {
		panic("(p *DBOperation)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBOperation)Save(){
	if p.ID < 0 {
		panic("(p *DBOperation)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBOperation)Remove(){
	db.Delete(p)
}
