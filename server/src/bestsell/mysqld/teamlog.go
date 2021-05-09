package mysqld

import (
	//"fmt"
	//"bestsell/common"
)

type DBTeamLog struct {
    MysqlModel
	PlayerId int
}

func startDBTeamLog()  {
	if !db.HasTable(&DBTeamLog{}) {
		db.CreateTable(&DBTeamLog{})
	}
}

//DBTeamLogBox
type DBTeamLogBox struct {
	MysqlModelBox
	PlayerId int
	items []*DBTeamLog
}
func (p *DBTeamLogBox)AddTeamLog(item *DBTeamLog) {
	item.PlayerId = p.PlayerId
	item.Insert()
	p.BeginWrite()
	p.items = append(p.items, item)
	p.EndWrite()
}
func (p *DBTeamLogBox)GetTeamLog(id int)*DBTeamLog {
	for _,item := range p.items {
		if item.ID == id {
			return item
		}
	}
	return nil
}
func (p *DBTeamLogBox)GetTeamLogs()*[]*DBTeamLog {
	return &p.items
}
func (p *DBTeamLogBox)RemoveTeamLog(id int) {
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
func GetDBTeamLogBoxFromDB(playerId int)*DBTeamLogBox  {
	box := DBTeamLogBox{
		PlayerId:playerId,
	}
	db.Where("player_id = ?", playerId).Find(&box.items)
	return &box
}


//DBTeamLog
func (p *DBTeamLog)Insert()  {
	if p.ID < 0 {
		panic("(p *DBTeamLog)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBTeamLog)Load(){
	if p.ID < 0 {
		panic("(p *DBTeamLog)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBTeamLog)Save(){
	if p.ID < 0 {
		panic("(p *DBTeamLog)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBTeamLog)Remove(){
	db.Delete(p)
}
