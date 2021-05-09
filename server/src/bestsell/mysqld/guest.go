package mysqld

import (
	//"fmt"
	//"bestsell/common"
)

type DBGuest struct {
    MysqlModel
	PlayerId  	int
    GuestId  	int
    GuestName  	string
}

func startDBGuest()  {
	if !db.HasTable(&DBGuest{}) {
		db.CreateTable(&DBGuest{})
	}
}

//DBGuestBox
type DBGuestBox struct {
	MysqlModelBox
	PlayerId int
	items []*DBGuest
}
func (p *DBGuestBox)AddGuest(item *DBGuest) {
	item.PlayerId = p.PlayerId
	item.Insert()
	p.BeginWrite()
	p.items = append(p.items, item)
	p.EndWrite()
}
func (p *DBGuestBox)GetGuest(id int)*DBGuest {
	for _,item := range p.items {
		if item.ID == id {
			return item
		}
	}
	return nil
}
func (p *DBGuestBox)GetGuests()*[]*DBGuest {
	return &p.items
}
func (p *DBGuestBox)RemoveGuest(id int) {
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
func GetDBGuestBoxFromDB(playerId int)*DBGuestBox  {
	box := DBGuestBox{
		PlayerId:playerId,
	}
	db.Where("player_id = ?", playerId).Find(&box.items)
	return &box
}


//DBGuest
func (p *DBGuest)Insert()  {
	if p.ID < 0 {
		panic("(p *DBGuest)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBGuest)Load(){
	if p.ID < 0 {
		panic("(p *DBGuest)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBGuest)Save(){
	if p.ID < 0 {
		panic("(p *DBGuest)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBGuest)Remove(){
	db.Delete(p)
}
