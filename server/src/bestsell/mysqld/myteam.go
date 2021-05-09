package mysqld

import "fmt"

type TeamStatus int
const (
	TeamStatusNo      	TeamStatus = 0
	TeamStatusCheck  	TeamStatus = 1
	TeamStatusRefuse  	TeamStatus = 2
	TeamStatusSeller  	TeamStatus = 3
	TeamStatusCancel  	TeamStatus = 4
)

type DBMyTeam struct {
    MysqlModel
	PlayerId 		int `gorm:"not null;unique"`
	Status 			int
	Nickname 		string
	Mobile          string
}

//var dbMyTeamSafeMap common.SafeMap

//func GetDBMyTeam(id int)*DBMyTeam {
//	ret := dbMyTeamSafeMap.Get(id)
//	if ret == nil {
//		return nil
//	}
//	return ret.(*DBMyTeam)
//}
//
//func AddDBMyTeam(item *DBMyTeam)  {
//	old := GetDBMyTeam(item.ID)
//	if old != nil {
//		fmt.Println("DBMyTeam repeat")
//		return
//	}
//	dbMyTeamSafeMap.Set(item.ID, item)
//}
//
//func LoadDBMyTeams()  {
//	var _dbItemsSlice []*DBMyTeam
//	db.Find(&_dbItemsSlice)
//	dbMyTeamSafeMap = *common.NewSafeMap()
//	for _,item := range _dbItemsSlice{
//		dbMyTeamSafeMap.Set(item.ID, item)
//	}
//}

func startDBMyTeam()  {
	if !db.HasTable(&DBMyTeam{}) {
		db.CreateTable(&DBMyTeam{})
	}
	//LoadDBMyTeams()
}

func GetDBMyTeamFromDB(playerId int, needCreate bool)*DBMyTeam  {
	item := &DBMyTeam{
		PlayerId: playerId,
	}
	item.LoadByPlayerId()
	if needCreate && item.ID == 0 {
		item.Insert()
	}
	return item
}

func GetAllDBMyTeams()*[]*DBMyTeam {
	var _dbItemsSlice []*DBMyTeam
	db.Find(&_dbItemsSlice)
	return &_dbItemsSlice
}

//DBMyTeam
func (p *DBMyTeam)SetTeamStatus(status int)  {
	switch TeamStatus(status) {
		case TeamStatusNo:
		case TeamStatusCheck:
		case TeamStatusRefuse:
		case TeamStatusSeller:
		case TeamStatusCancel:
		default:
			fmt.Println("SetTeamStatus 未知状态 status =", status)
			return
	}
	p.BeginWrite()
	p.Status = status
	p.EndWrite()
	p.DelaySave(p)
}

func (p *DBMyTeam)GetTeamStatus()int {
	if p.ID <= 0 {
		return int(TeamStatusNo)
	}
	return p.Status
}

func (p *DBMyTeam)GetTeamStatusName()string {
	status := p.GetTeamStatus()
	switch TeamStatus(status) {
	case TeamStatusRefuse:
		return "条件不足，请联系客服"
	}
	return ""
}

func (p *DBMyTeam)Insert()  {
	if p.ID < 0 {
		panic("(p *DBMyTeam)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBMyTeam)LoadByPlayerId(){
	if p.PlayerId <= 0 {
		panic("(p *DBMyTeam)LoadByPlayerId p.Token < 0")
	}
	p.ID = 0
	db.Where("player_id = ?", p.PlayerId).First(p)
	if p.ID <= 0 {
		//panic("(p *DBPlayerInfo)LoadWithPlayerId p.ID < 0")
		return
	}
}

func (p *DBMyTeam)Load(){
	if p.ID < 0 {
		panic("(p *DBMyTeam)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBMyTeam)Save(){
	if p.ID < 0 {
		panic("(p *DBMyTeam)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBMyTeam)Remove(){
	db.Delete(p)
}
