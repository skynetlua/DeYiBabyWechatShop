package mysqld

import (
	"fmt"
	"bestsell/common"
)

type DBReputation struct {
    MysqlModel
	PlayerId 	int
	OrderId  	int
	GoodsId     int
	SkuId  		int
	Repute 		int
	Remark 		string

    PlayerName string
	AvatarUrl string
}

var dbReputationBoxSafeMap common.SafeMap

func GetDBReputationBox(id int)*DBReputationBox {
	ret := dbReputationBoxSafeMap.Get(id)
	if ret == nil {
		return nil
	}
	return ret.(*DBReputationBox)
}

func AddDBReputationBox(item *DBReputationBox) {
	old := GetDBReputationBox(item.GoodsId)
	if old != nil {
		fmt.Println("AddDBReputationBox repeat")
		return
	}
	dbReputationBoxSafeMap.Set(item.GoodsId, item)
}

//func LoadDBReputationsBox()  {
//	var _dbItemsSlice []*DBReputationBox
//	db.Find(&_dbItemsSlice)
//	dbReputationSafeMap = *common.NewSafeMap()
//	for _,item := range _dbItemsSlice{
//		dbReputationSafeMap.Set(item.ID, item)
//	}
//}

func startDBReputation() {
	if !db.HasTable(&DBReputation{}) {
		db.CreateTable(&DBReputation{})
	}
	dbReputationBoxSafeMap = *common.NewSafeMap()
	//LoadDBReputations()
}

func AddNewReputation(reputation *DBReputation) {
	reputationBox := GetDBReputationBox(reputation.GoodsId)
	if reputationBox != nil {
		reputationBox.AddItem(reputation)
		return
	}
	reputation.Insert()
}

//
//DBReputationBox
type DBReputationBox struct {
	MysqlModelBox
	GoodsId int
	isDirty bool
	items []*DBReputation
}

func GetDBReputationBoxFromDB(goodsId int)*DBReputationBox  {
	box := DBReputationBox{
		GoodsId:goodsId,
	}
	db.Where("goods_id = ?", goodsId).Find(&box.items)
	return &box
}

func GetDBReputationBoxOrFromDB(goodsId int)*DBReputationBox {
	dbReputationBox := GetDBReputationBox(goodsId)
	if dbReputationBox != nil  {
		return dbReputationBox
	}
	dbReputationBox = GetDBReputationBoxFromDB(goodsId)
	return dbReputationBox
}

func (p *DBReputationBox)AddItem(item *DBReputation) {
	item.Insert()
	p.BeginWrite()
	p.items = append(p.items, item)
	p.EndWrite()
}

func (p *DBReputationBox)GetItem(id int)*DBReputation {
	for _,item := range p.items {
		if item.ID == id {
			return item
		}
	}
	return nil
}

func (p *DBReputationBox)GetItems()*[]*DBReputation {
	return &p.items
}

func (p *DBReputationBox)RemoveItem(id int) {
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

//DBReputation
func (p *DBReputation)Insert()  {
	if p.ID < 0 {
		panic("(p *DBReputation)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBReputation)Load(){
	if p.ID < 0 {
		panic("(p *DBReputation)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBReputation)Save(){
	if p.ID < 0 {
		panic("(p *DBReputation)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBReputation)Remove(){
	db.Delete(p)
}
