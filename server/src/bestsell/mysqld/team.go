package mysqld

import (
	"fmt"
	"bestsell/common"
)

type TeamPostType int
const (
	TeamPostChairman  	TeamPostType = 4
	TeamPostSupremo  	TeamPostType = 3
	TeamPostManager  	TeamPostType = 2
	TeamPostStaff     	TeamPostType = 1
	TeamPostJobless     TeamPostType = 0
)

type OptTeam struct {
	Posts [3]int
}

type DBTeam struct {
    MysqlModel

	Score int
	Level int
    teamInfo *DBTeamInfo
	memberBox *DBMemberBox

    optData *OptTeam
}

var dbTeamSafeMap common.SafeMap

func GetDBTeam(id int)*DBTeam {
	ret := dbTeamSafeMap.Get(id)
	if ret == nil {
		return nil
	}
	return ret.(*DBTeam)
}

func AddDBTeam(item *DBTeam)  {
	old := GetDBTeam(item.ID)
	if old != nil {
		fmt.Println("DBTeam repeat")
		return
	}
	dbTeamSafeMap.Set(item.ID, item)
}

func LoadDBTeams()  {
	var _dbItemsSlice []*DBTeam
	db.Find(&_dbItemsSlice)
	dbTeamSafeMap = *common.NewSafeMap()
	for _,item := range _dbItemsSlice{
		dbTeamSafeMap.Set(item.ID, item)
	}
}

func startDBTeam()  {
	if !db.HasTable(&DBTeam{}) {
		db.CreateTable(&DBTeam{})
	}
	LoadDBTeams()
}

func EachTeam(eachFunc func(team *DBTeam))  {
	iterFunc := func(key int, v interface{}) bool{
		eachFunc(v.(*DBTeam))
		return true
	}
	dbTeamSafeMap.RangeSafe(iterFunc)
}

//DBTeam
func (p *DBTeam)GetOptData()*OptTeam {
	if p.optData == nil {
		p.optData = &OptTeam{}
	}
	return p.optData
}

func (p *DBTeam)GetMemberBox()*DBMemberBox {
	if p.memberBox == nil {
		p.memberBox = GetDBMemberBoxFromDB(p.ID)
	}
	return p.memberBox
}

func (p *DBTeam)GetTeamInfo()*DBTeamInfo {
	if p.teamInfo == nil {
		p.teamInfo = &DBTeamInfo{
			TeamId:p.ID,
		}
		p.teamInfo.LoadByTeamId()
		if p.teamInfo.ID == 0 {
			p.teamInfo.Insert()
		}
	}
	return p.teamInfo
}

func (p *DBTeam)Insert()  {
	if p.ID < 0 {
		panic("(p *DBTeam)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBTeam)Load(){
	if p.ID < 0 {
		panic("(p *DBTeam)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBTeam)Save(){
	if p.ID < 0 {
		panic("(p *DBTeam)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBTeam)Remove(){
	db.Delete(p)
}
