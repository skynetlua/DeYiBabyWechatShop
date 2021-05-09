package mysqld

import "sync"

type EGoodsAction int
const (
	EGoodsActionVisit  	 EGoodsAction = 0
	EGoodsActionCart  	 EGoodsAction = 1
	EGoodsActionPrepare  EGoodsAction = 2
	EGoodsActionOrder  	 EGoodsAction = 3
	EGoodsActionPay  	 EGoodsAction = 4
)

type DBGoodsStat struct {
    MysqlModel
	GoodsId  int
	SkuId int
	PlayerId int
	OrderId int
	Action int
    Number int

    playerName string
	playerIcon string
	goodsName string
}

var _lastGoodsStatSlice []*DBGoodsStat
var _lastGoodsStatsCache []map[string]interface{}
var _goodsStatLockMutex sync.RWMutex


func loadLastGoodsStatSlice() {
	var goodsStatSlice []*DBGoodsStat
	db.Where("action = ?", EGoodsActionPay).Order("id desc").Limit(12).Find(&goodsStatSlice)
	_lastGoodsStatSlice = []*DBGoodsStat{}
	for _,item := range goodsStatSlice {
		dbPlayer := GetDBPlayerByPlayerIdOrFromDB(item.PlayerId)
		if dbPlayer != nil {
			dbPlayerInfo := dbPlayer.GetPlayerInfo()
			dbGoods := GetDBGoods(item.GoodsId)
			if dbGoods != nil {
				dbGoodsInfo := dbGoods.GetGoodsInfo()
				item.playerName = dbPlayerInfo.NickName
				item.goodsName = dbGoodsInfo.Name
				item.playerIcon = dbPlayerInfo.AvatarUrl

				_lastGoodsStatSlice = append(_lastGoodsStatSlice, item)
				if len(_lastGoodsStatSlice) >= 10 {
					break
				}
			}
		}
	}
	loadLastGoodsStatsCache()
}

func loadLastGoodsStatsCache() {
	_goodsStatLockMutex.Lock()
	defer _goodsStatLockMutex.Unlock()

	var tmps []map[string]interface{}
	for _,item := range _lastGoodsStatSlice {
		if item.playerName == "量子出车" {
			item.playerName = "小天使"
		}
		data := map[string]interface{} {
			"playerName":  item.playerName,
			"goodsId":  item.GoodsId,
			"goodsName":  item.goodsName,
			"playerIcon":  item.playerIcon,
		}
		tmps = append(tmps, data)
	}
	_lastGoodsStatsCache = tmps
}

func addLastGoodsStats(goodsStat *DBGoodsStat) {
	_goodsStatLockMutex.Lock()
	defer _goodsStatLockMutex.Unlock()

	if len(_lastGoodsStatSlice) < 10 {
		_lastGoodsStatSlice = append(_lastGoodsStatSlice, goodsStat)
	} else {
		_lastGoodsStatSlice = append([]*DBGoodsStat{goodsStat}, _lastGoodsStatSlice...)
		_lastGoodsStatSlice = _lastGoodsStatSlice[:10]
	}
	loadLastGoodsStatsCache()
}

func GetLastGoodsStatList() *[]map[string]interface{} {
	_goodsStatLockMutex.Lock()
	defer _goodsStatLockMutex.Unlock()

	var tmps []map[string]interface{}
	tmps = _lastGoodsStatsCache
	return &tmps
}

func GetGoodsStatSlice() []*DBGoodsStat {
	var goodsStatSlice []*DBGoodsStat
	db.Order("id desc").Limit(200).Find(&goodsStatSlice)
	_lastGoodsStatSlice = []*DBGoodsStat{}
	for _,item := range goodsStatSlice {
		dbPlayer := GetDBPlayerByPlayerIdOrFromDB(item.PlayerId)
		if dbPlayer != nil {
			dbPlayerInfo := dbPlayer.GetPlayerInfo()
			dbGoods := GetDBGoods(item.GoodsId)
			if dbGoods != nil {
				dbGoodsInfo := dbGoods.GetGoodsInfo()
				if dbGoodsInfo != nil {
					item.playerName = dbPlayerInfo.NickName
					item.goodsName = dbGoodsInfo.Name
					item.playerIcon = dbPlayerInfo.AvatarUrl
				}
			}
		}
	}
	return goodsStatSlice
}


func startDBGoodsStat() {
	if !db.HasTable(&DBGoodsStat{}) {
		db.CreateTable(&DBGoodsStat{})
	}
	loadLastGoodsStatSlice()
}

func MakeVisitGoodsStat(action EGoodsAction, player *DBPlayer, goods *DBGoods, skuId int, orderId int, number int) {
	goodsStat := DBGoodsStat {
		PlayerId: player.ID,
		GoodsId: goods.GoodsId,
		SkuId: skuId,
		OrderId: orderId,
		Action: int(action),
		Number: number,
	}
	goodsStat.Insert()
	if action == EGoodsActionPay {
		goodsStatBox := player.GetGoodsStatBox()
		goodsStatBox.AddGoodsStat(&goodsStat)

		goodsStat.playerName = player.GetPlayerInfo().NickName
		goodsStat.goodsName = goods.GetGoodsInfo().Name
		goodsStat.playerIcon = player.GetPlayerInfo().AvatarUrl
		addLastGoodsStats(&goodsStat)
	}
}

//DBGoodsStatBox
type DBGoodsStatBox struct {
	MysqlModelBox
	PlayerId int
	items []*DBGoodsStat
}

func (p *DBGoodsStatBox)AddGoodsStat(item *DBGoodsStat) {
	p.BeginWrite()
	p.items = append(p.items, item)
	p.EndWrite()
}

func (p *DBGoodsStatBox)BuyCount(goodsId int) int {
	count := 0
	for _,item := range p.items {
		if item.GoodsId == goodsId && item.Action == int(EGoodsActionPay) {
			count++
		}
	}
	return count
}

func GetDBGoodsStatBoxFromDB(playerId int) *DBGoodsStatBox {
	box := DBGoodsStatBox {
		PlayerId:playerId,
	}
	db.Where("player_id = ? and action = ?", playerId, EGoodsActionPay).Find(&box.items)
	return &box
}

//DBGoodsStat
func (p *DBGoodsStat)Insert() {
	if p.ID < 0 {
		panic("(p *DBLogistics)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBGoodsStat)Load() {
	if p.ID < 0 {
		panic("(p *DBLogistics)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBGoodsStat)Save() {
	if p.ID < 0 {
		panic("(p *DBLogistics)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBGoodsStat)Remove() {
	//db.Delete(p)
}
