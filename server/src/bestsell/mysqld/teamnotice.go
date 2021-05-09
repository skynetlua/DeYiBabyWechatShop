package mysqld

import (
	//"fmt"
	//"bestsell/common"
)

type DBTeamNotice struct {
    MysqlModel
	PlayerId int
}

func startDBTeamNotice()  {
	if !db.HasTable(&DBTeamNotice{}) {
		db.CreateTable(&DBTeamNotice{})
	}
}

//DBTeamNoticeBox
type DBTeamNoticeBox struct {
	MysqlModelBox
	PlayerId int
	items []*DBTeamNotice
}
func (p *DBTeamNoticeBox)AddTeamNotice(item *DBTeamNotice) {
	item.PlayerId = p.PlayerId
	item.Insert()
	p.BeginWrite()
	p.items = append(p.items, item)
	p.EndWrite()
}
func (p *DBTeamNoticeBox)GetTeamNotice(id int)*DBTeamNotice {
	for _,item := range p.items {
		if item.ID == id {
			return item
		}
	}
	return nil
}
func (p *DBTeamNoticeBox)GetTeamNotices()*[]*DBTeamNotice {
	return &p.items
}
func (p *DBTeamNoticeBox)RemoveTeamNotice(id int) {
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
func GetDBTeamNoticeBoxFromDB(playerId int)*DBTeamNoticeBox  {
	box := DBTeamNoticeBox{
		PlayerId:playerId,
	}
	db.Where("player_id = ?", playerId).Find(&box.items)
	return &box
}


//DBTeamNotice
func (p *DBTeamNotice)Insert()  {
	if p.ID < 0 {
		panic("(p *DBTeamNotice)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBTeamNotice)Load(){
	if p.ID < 0 {
		panic("(p *DBTeamNotice)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBTeamNotice)Save(){
	if p.ID < 0 {
		panic("(p *DBTeamNotice)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBTeamNotice)Remove(){
	db.Delete(p)
}
