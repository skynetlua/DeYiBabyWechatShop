package mysqld



type CashLogType int
const (
	//购买商铺消费
	CashLogPay      	CashLogType = 0
	CashLogCommission  	CashLogType = 1
	CashLogWithDraw  	CashLogType = 2
)

type DBCashLog struct {
    MysqlModel
	PlayerId int
	CashType int
	Behavior int
	Amount float64
}

//var dbCashLogBoxSafeMap common.SafeMap
//
//func GetDBCashLogBox(playerId int)*DBCashLogBox {
//	ret := dbCashLogBoxSafeMap.Get(playerId)
//	if ret == nil {
//		return nil
//	}
//	return ret.(*DBCashLogBox)
//}
//
//func AddDBCashLogBox(item *DBCashLogBox)  {
//	old := GetDBCashLogBox(item.PlayerId)
//	if old != nil {
//		fmt.Println("DBCashLogBox repeat")
//		return
//	}
//	dbCashLogBoxSafeMap.Set(item.PlayerId, item)
//}

//func GetDBCashLogBoxOrFromDB(playerId int)*DBCashLogBox {
//	dbCashLogBox := GetDBCashLogBox(playerId)
//	if dbCashLogBox != nil  {
//		return dbCashLogBox
//	}
//	dbCashLogBox = GetDBCashLogBoxFromDB(playerId)
//	return dbCashLogBox
//}
//
//func LoadDBCashLogs()  {
//	var _dbItemsSlice []*DBCashLog
//	db.Find(&_dbItemsSlice)
//	dbCashLogBoxSafeMap = *common.NewSafeMap()
//	for _,item := range _dbItemsSlice{
//		dbCashLogBoxSafeMap.Set(item.ID, item)
//	}
//}

//func AddNewCashLog(item *DBCashLog)  {
//	box := GetDBCashLogBox(item.PlayerId)
//	if box != nil {
//		box.AddCashLog(item)
//		return
//	}
//	item.Insert()
//}

func startDBCashLog()  {
	if !db.HasTable(&DBCashLog{}) {
		db.CreateTable(&DBCashLog{})
	}
	//dbCashLogBoxSafeMap = *common.NewSafeMap()
}


//DBCashLogBox
type DBCashLogBox struct {
	MysqlModelBox
	PlayerId int
	items []*DBCashLog
}
func (p *DBCashLogBox)AddCashLog(item *DBCashLog) {
	item.PlayerId = p.PlayerId
	item.Insert()
	p.BeginWrite()
	p.items = append(p.items, item)
	p.EndWrite()
}
func (p *DBCashLogBox)GetCashLog(id int)*DBCashLog {
	for _,item := range p.items {
		if item.ID == id {
			return item
		}
	}
	return nil
}
func (p *DBCashLogBox)GetCashLogs()*[]*DBCashLog {
	return &p.items
}
func (p *DBCashLogBox)RemoveCashLog(id int) {
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
func GetDBCashLogBoxFromDB(playerId int)*DBCashLogBox  {
	box := DBCashLogBox{
		PlayerId:playerId,
	}
	db.Where("player_id = ?", playerId).Find(&box.items)
	return &box
}


//DBCashlog
func (p *DBCashLog)GetCashTypeName()string{
	switch CashLogType(p.CashType) {
		case CashLogPay:
			return "购物消费"
		case CashLogCommission:
			return "佣金获得"
		case CashLogWithDraw:
			return "现金提现"
	}
	return "未知"
}


func (p *DBCashLog)Insert()  {
	if p.ID < 0 {
		panic("(p *DBCashLog)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBCashLog)Load(){
	if p.ID < 0 {
		panic("(p *DBCashLog)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBCashLog)Save(){
	if p.ID < 0 {
		panic("(p *DBCashLog)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBCashLog)Remove(){
	db.Delete(p)
}
