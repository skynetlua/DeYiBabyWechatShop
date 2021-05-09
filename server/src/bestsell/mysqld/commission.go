package mysqld

import (
	//"fmt"
	//"bestsell/common"
)

type DBCommission struct {
    MysqlModel
	PlayerId  	int
	Level  		int
	Money  		float64
	Ratio  		float64
	SellerId  	int
	SellerName  string
	BuyerId 	int
	BuyerName 	string
}

func startDBCommission()  {
	if !db.HasTable(&DBCommission{}) {
		db.CreateTable(&DBCommission{})
	}
}

//DBCommissionBox
type DBCommissionBox struct {
	MysqlModelBox
	PlayerId int
	items []*DBCommission
}
func (p *DBCommissionBox)AddCommission(item *DBCommission) {
	item.PlayerId = p.PlayerId
	item.Insert()
	p.BeginWrite()
	p.items = append(p.items, item)
	p.EndWrite()
}
func (p *DBCommissionBox)GetCommission(id int)*DBCommission {
	for _,item := range p.items {
		if item.ID == id {
			return item
		}
	}
	return nil
}
func (p *DBCommissionBox)GetCommissions()*[]*DBCommission {
	return &p.items
}
func (p *DBCommissionBox)RemoveCommission(id int) {
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
func GetDBCommissionBoxFromDB(playerId int)*DBCommissionBox  {
	box := DBCommissionBox{
		PlayerId:playerId,
	}
	db.Where("player_id = ?", playerId).Find(&box.items)
	return &box
}


//DBCommission
func (p *DBCommission)Insert() {
	if p.ID < 0 {
		panic("(p *DBCommission)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBCommission)Load(){
	if p.ID < 0 {
		panic("(p *DBCommission)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBCommission)Save(){
	if p.ID < 0 {
		panic("(p *DBCommission)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBCommission)Remove(){
	db.Delete(p)
}
