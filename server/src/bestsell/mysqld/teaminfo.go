package mysqld

type DBTeamInfo struct {
    MysqlModel
    TeamId  	int  `gorm:"not null;unique"`
	Name  		string
	Slogan  	string
	LeaderId  	int  `gorm:"not null;unique"`
}

//var dbTeamInfoSafeMap common.SafeMap
//
//func GetDBTeamInfo(id int)*DBTeamInfo {
//	ret := dbTeamInfoSafeMap.Get(id)
//	if ret == nil {
//		return nil
//	}
//	return ret.(*DBTeamInfo)
//}
//
//func AddDBTeamInfo(item *DBTeamInfo)  {
//	old := GetDBTeamInfo(item.ID)
//	if old != nil {
//		fmt.Println("DBTeamInfo repeat")
//		return
//	}
//	dbTeamInfoSafeMap.Set(item.ID, item)
//}
//
//func LoadDBTeamInfos()  {
//	var _dbItemsSlice []*DBTeamInfo
//	db.Find(&_dbItemsSlice)
//	dbTeamInfoSafeMap = *common.NewSafeMap()
//	for _,item := range _dbItemsSlice{
//		dbTeamInfoSafeMap.Set(item.ID, item)
//	}
//}

func startDBTeamInfo()  {
	if !db.HasTable(&DBTeamInfo{}) {
		db.CreateTable(&DBTeamInfo{})
	}
	//LoadDBTeamInfos()
}

//DBTeamInfo
func (p *DBTeamInfo)Insert()  {
	if p.ID < 0 {
		panic("(p *DBTeamInfo)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBTeamInfo)Load(){
	if p.ID < 0 {
		panic("(p *DBTeamInfo)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBTeamInfo)LoadByTeamId(){
	if p.TeamId < 0 {
		panic("(p *DBTeamInfo)Load p.TeamId < 0")
	}
	db.Where("`team_id` = ?", p.TeamId).First(p)
}

func (p *DBTeamInfo)Save(){
	if p.ID < 0 {
		panic("(p *DBTeamInfo)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBTeamInfo)Remove(){
	db.Delete(p)
}
