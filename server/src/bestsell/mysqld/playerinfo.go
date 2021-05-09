package mysqld

type DBPlayerInfo struct {
    MysqlModel
	PlayerId 		int  `gorm:"not null;unique"`
	NickName   		string
	Gender   		int
	Language   		string
	City   			string
	Province   		string
	Country   		string
	AvatarUrl   	string
	Timestamp   	int
	Referrer   		string
	Mobile     		string
}

//var dbPlayerInfoSafeMap common.SafeMap
//
//func GetDBPlayerInfo(id int)*DBPlayerInfo {
//	ret := dbPlayerInfoSafeMap.Get(id)
//	if ret == nil {
//		return nil
//	}
//	return ret.(*DBPlayerInfo)
//}
//
//func AddDBPlayerInfo(item *DBPlayerInfo)  {
//	old := GetDBPlayerInfo(item.ID)
//	if old != nil {
//		fmt.Println("DBPlayerInfo repeat")
//		return
//	}
//	dbPlayerInfoSafeMap.Set(item.ID, item)
//}
//
//func LoadDBPlayerInfos()  {
//	var _dbItemsSlice []*DBPlayerInfo
//	db.Find(&_dbItemsSlice)
//	dbPlayerInfoSafeMap = *common.NewSafeMap()
//	for _,item := range _dbItemsSlice{
//		dbPlayerInfoSafeMap.Set(item.ID, item)
//	}
//}
func GetDBPlayerInfosByPlayerIdsFromDB(playerId *[]int) *[]*DBPlayerInfo{
	var playerInfos []*DBPlayerInfo
	db.Where("player_id in (?)", *playerId).Find(&playerInfos)
	return &playerInfos
}

func startDBPlayerInfo()  {
	if !db.HasTable(&DBPlayerInfo{}) {
		db.CreateTable(&DBPlayerInfo{})
	}
	//LoadDBPlayerInfos()
}

//DBPlayerInfo
func (p *DBPlayerInfo)Insert()  {
	if p.ID < 0 {
		panic("(p *DBPlayerInfo)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBPlayerInfo)Load(){
	if p.ID < 0 {
		panic("(p *DBPlayerInfo)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBPlayerInfo)LoadByPlayerId(){
	if p.PlayerId <= 0 {
		panic("(p *DBPlayer)Load p.Token < 0")
	}
	p.ID = 0
	db.Where("player_id = ?", p.PlayerId).First(p)
	if p.ID <= 0 {
		//panic("(p *DBPlayerInfo)LoadWithPlayerId p.ID < 0")
		return
	}
}

func (p *DBPlayerInfo)Save(){
	if p.ID < 0 {
		panic("(p *DBPlayerInfo)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBPlayerInfo)Remove(){
	db.Delete(p)
}

func GetDBPlayerInfoFromDB(playerId int, needCreate bool)*DBPlayerInfo  {
	item := &DBPlayerInfo{
		PlayerId: playerId,
	}
	item.LoadByPlayerId()
	if needCreate && item.ID == 0 {
		item.Insert()
	}
	return item
}
