package mysqld

type DBWithdrawLog struct {
    MysqlModel
	PlayerId int
	Status int
	Money int
}

//var dbWithdrawLogBoxSafeMap common.SafeMap
//
//func GetDBWithdrawLogBox(playerId int)*DBWithdrawLogBox {
//	ret := dbWithdrawLogBoxSafeMap.Get(playerId)
//	if ret == nil {
//		return nil
//	}
//	return ret.(*DBWithdrawLogBox)
//}
//
//func AddDBWithdrawLogBox(item *DBWithdrawLogBox)  {
//	old := GetDBWithdrawLogBox(item.PlayerId)
//	if old != nil {
//		fmt.Println("DBWithdrawLogBox repeat")
//		return
//	}
//	dbWithdrawLogBoxSafeMap.Set(item.PlayerId, item)
//}

//func GetDBWithdrawLogBoxOrFromDB(playerId int)*DBWithdrawLogBox {
//	dbWithdrawLogBox := GetDBWithdrawLogBox(playerId)
//	if dbWithdrawLogBox != nil  {
//		return dbWithdrawLogBox
//	}
//	dbWithdrawLogBox = GetDBWithdrawLogBoxFromDB(playerId)
//	return dbWithdrawLogBox
//}
//
//func LoadDBWithdrawLogs()  {
//	var _dbItemsSlice []*DBWithdrawLog
//	db.Find(&_dbItemsSlice)
//	dbWithdrawLogBoxSafeMap = *common.NewSafeMap()
//	for _,item := range _dbItemsSlice{
//		dbWithdrawLogBoxSafeMap.Set(item.ID, item)
//	}
//}

func startDBWithdrawLog()  {
	if !db.HasTable(&DBWithdrawLog{}) {
		db.CreateTable(&DBWithdrawLog{})
	}
}

//func AddNewWithdrawLog(item *DBWithdrawLog)  {
//	box := GetDBWithdrawLogBox(item.PlayerId)
//	if box != nil {
//		box.AddWithdrawLog(item)
//		return
//	}
//	item.Insert()
//}


//DBWithdrawLogBox
type DBWithdrawLogBox struct {
	MysqlModelBox
	PlayerId int
	items []*DBWithdrawLog
}
func (p *DBWithdrawLogBox)AddItem(item *DBWithdrawLog) {
	item.PlayerId = p.PlayerId
	item.Insert()
	p.BeginWrite()
	p.items = append(p.items, item)
	p.EndWrite()
}
func (p *DBWithdrawLogBox)GetItem(id int)*DBWithdrawLog {
	for _,item := range p.items {
		if item.ID == id {
			return item
		}
	}
	return nil
}
func (p *DBWithdrawLogBox)GetItems()*[]*DBWithdrawLog {
	return &p.items
}
func (p *DBWithdrawLogBox)RemoveItem(id int) {
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
func GetDBWithdrawLogBoxFromDB(playerId int)*DBWithdrawLogBox  {
	box := DBWithdrawLogBox{
		PlayerId:playerId,
	}
	db.Where("player_id = ?", playerId).Find(&box.items)
	return &box
}

type WithdrawLogStatus int
const (
	//购买商铺消费
	WithdrawLogCheck      	WithdrawLogStatus = 0
	WithdrawLogFinish  		WithdrawLogStatus = 1
)

//DBWithdrawLog
func (p *DBWithdrawLog)GetWithdrawTypeName()string{
	switch WithdrawLogStatus(p.Status) {
	case WithdrawLogCheck:
		return "审核中"
	case WithdrawLogFinish:
		return "已到账"
	}
	return "未知"
}

func (p *DBWithdrawLog)Insert()  {
	if p.ID < 0 {
		panic("(p *DBWithdrawLog)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBWithdrawLog)Load(){
	if p.ID < 0 {
		panic("(p *DBWithdrawLog)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBWithdrawLog)Save(){
	if p.ID < 0 {
		panic("(p *DBWithdrawLog)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBWithdrawLog)Remove(){
	db.Delete(p)
}
