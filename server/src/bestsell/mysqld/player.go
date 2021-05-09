package mysqld

import (
	"bestsell/common"
	"fmt"
)

type DBPlayer struct {
    MysqlModel
	Token       	string 	`gorm:"not null;unique"`
	OpenId     		string	`gorm:"not null;unique"`
	SessionKey 		string
	UnionId    		string

	RefererId       int
	RefererName     string
    TeamId     		int
    TeamPost  		int

	Balance     	float64
	AmountCost     	float64
	Score     		int
	Growth     		int
    GM  			int
    VIP     		int

	LoginIP 		string

    playerInfo      *DBPlayerInfo
    // myTeam       	*DBMyTeam
    addressBox 		*DBAddressBox
    cartBox 		*DBCartBox
    favoriteBox 	*DBFavoriteBox
    orderBox 		*DBOrderBox

	cashLogBox 		*DBCashLogBox
	withdrawLogBox 	*DBWithdrawLogBox
	commissionBox 	*DBCommissionBox
	goodsStatBox 	*DBGoodsStatBox
}

var dbPlayerSafeMap common.SafeMapS
var dbPlayerSafeMapI common.SafeMap


func init() {
}

func GetDBPlayerByOpenIdOrFromDB(openId string) *DBPlayer {
	token := openId
	player := GetDBPlayer(token)
	if player != nil {
		return player
	}
	player = &DBPlayer{
		Token: token,
	}
	player.LoadWithToken()
	if player.ID <= 0 {
		return nil
	}
	AddDBPlayer(player)
	return player
}

func GetDBPlayerByPlayerId(playerId int)*DBPlayer {
	ret := dbPlayerSafeMapI.Get(playerId)
	if ret == nil {
		return nil
	}
	return ret.(*DBPlayer)
}

func GetDBPlayerByPlayerIdOrFromDB(playerId int) *DBPlayer {
	ret := dbPlayerSafeMapI.Get(playerId)
	if ret != nil {
		return ret.(*DBPlayer)
	}
	var player DBPlayer
	db.First(&player, playerId)
	if player.ID <= 0 {
		return nil
	}
	return &player
}

func GetDBPlayersByPlayerIdsFromDB(playerId *[]int) *[]*DBPlayer{
	var players []*DBPlayer
	db.Where("id in (?)", *playerId).Find(&players)
	return &players
}

func GetDBPlayer(token string)*DBPlayer {
	ret := dbPlayerSafeMap.Get(token)
	if ret == nil {
		return nil
	}
	return ret.(*DBPlayer)
}

func AddDBPlayer(item *DBPlayer) {
	old := GetDBPlayer(item.Token)
	if old != nil {
		fmt.Println("DBPlayer repeat")
		return
	}
	dbPlayerSafeMap.Set(item.Token, item)
	dbPlayerSafeMapI.Set(item.ID, item)
}

//func LoadDBPlayers() {
//	var _dbItemsSlice []*DBPlayer
//	db.Find(&_dbItemsSlice)
//	for _,item := range _dbItemsSlice{
//		dbPlayerSafeMap.Set(item.Token, item)
//	}
//}

func startDBPlayer() {
	if !db.HasTable(&DBPlayer{}) {
		db.CreateTable(&DBPlayer{})
	}
	dbPlayerSafeMap = *common.NewSafeMapS()
	dbPlayerSafeMapI = *common.NewSafeMap()
	//LoadDBPlayers()
}

//func (p *DBPlayer)SetWorkGoodsId(goodsId int) {
//	p.BeginWrite()
//	p.workGoodsId = goodsId
//	p.EndWrite()
//}

func (p *DBPlayer)Insert()  {
	if p.ID < 0 {
		panic("(p *DBPlayer)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBPlayer)Load(){
	if p.ID <= 0 {
		panic("(p *DBPlayer)Load p.ID < 0")
	}
	p.Token = ""
	db.First(p, p.ID)
	if len(p.Token) == 0 {
		panic("(p *DBPlayer)Load p.Token == ''")
	}
}

func (p *DBPlayer)LoadWithToken(){
	if len(p.Token) <= 0 {
		panic("(p *DBPlayer)Load p.Token < 0")
	}
	p.ID = 0
	db.Where("`token` = ?", p.Token).First(p)
	if p.ID <= 0 {
		//panic("(p *DBPlayer)Load p.ID < 0")
		return
	}
}

// func (p *DBPlayer)IsSeller()int {
// 	if p.GetMyTeam().GetTeamStatus() == int(TeamStatusSeller) {
// 		return 1
// 	}
// 	return 0
// }

func (p *DBPlayer)GetPlayerInfo()*DBPlayerInfo {
	if p.playerInfo == nil {
		box := GetDBPlayerInfoFromDB(p.ID, true)
		p.BeginWrite()
		p.playerInfo = box
		p.EndWrite()
	}
	return p.playerInfo
}

// func (p *DBPlayer)GetMyTeam()*DBMyTeam {
// 	if p.myTeam == nil {
// 		item := GetDBMyTeamFromDB(p.ID, true)
// 		p.BeginWrite()
// 		p.myTeam = item
// 		p.EndWrite()
// 	}
// 	return p.myTeam
// }

// func (p *DBPlayer)ClearMyTeam() {
// 	p.BeginWrite()
// 	p.myTeam = nil
// 	p.EndWrite()
// }

func (p *DBPlayer)GetGoodsStatBox()*DBGoodsStatBox {
	if p.goodsStatBox == nil {
		item := GetDBGoodsStatBoxFromDB(p.ID)
		p.BeginWrite()
		p.goodsStatBox = item
		p.EndWrite()
	}
	return p.goodsStatBox
}

func (p *DBPlayer)GoodsBuyCount(goodsId int) int {
	dbGoodsStatBox := p.GetGoodsStatBox()
	return dbGoodsStatBox.BuyCount(goodsId)
}
//
func (p *DBPlayer)GetCommissionBox()*DBCommissionBox {
	if p.commissionBox == nil {
		box := GetDBCommissionBoxFromDB(p.ID)
		p.BeginWrite()
		p.commissionBox = box
		p.EndWrite()
	}
	return p.commissionBox
}

func (p *DBPlayer)GetAddressBox()*DBAddressBox {
	if p.addressBox == nil {
		box := GetDBAddressBoxFromDB(p.ID)
		p.BeginWrite()
		p.addressBox = box
		p.EndWrite()
	}
	return p.addressBox
}

func (p *DBPlayer)GetCartBox()*DBCartBox {
	if p.cartBox == nil {
		box := GetDBCartBoxFromDB(p.ID)
		p.BeginWrite()
		p.cartBox = box
		p.EndWrite()
	}
	return p.cartBox
}

func (p *DBPlayer)GetFavoriteBox()*DBFavoriteBox {
	if p.favoriteBox == nil {
		box := GetDBFavoriteBoxFromDB(p.ID)
		p.BeginWrite()
		p.favoriteBox = box
		p.EndWrite()
	}
	return p.favoriteBox
}

func (p *DBPlayer)GetOrderBox()*DBOrderBox {
	if p.orderBox == nil {
		box := GetDBOrderBoxFromDB(p.ID)
		p.BeginWrite()
		p.orderBox = box
		p.EndWrite()
	}
	return p.orderBox
}

func (p *DBPlayer)ClearOrderBox() {
	p.BeginWrite()
	p.orderBox = nil
	p.EndWrite()
}

func (p *DBPlayer)GetCashLogBox()*DBCashLogBox {
	if p.cashLogBox == nil {
		box := GetDBCashLogBoxFromDB(p.ID)
		p.BeginWrite()
		p.cashLogBox = box
		p.EndWrite()
	}
	return p.cashLogBox
}

func (p *DBPlayer)GetWithdrawLogBox()*DBWithdrawLogBox {
	if p.withdrawLogBox == nil {
		box := GetDBWithdrawLogBoxFromDB(p.ID)
		p.BeginWrite()
		p.withdrawLogBox = box
		p.EndWrite()
	}
	return p.withdrawLogBox
}

//
func (p *DBPlayer)CostMoney(amount float64, cashType CashLogType)  {
	cashLog := &DBCashLog{
		PlayerId: p.ID,
		CashType: int(cashType),
		Behavior: 1,
		Amount:   amount,
	}
	cashLogBox := p.GetCashLogBox()
	cashLogBox.AddCashLog(cashLog)

	p.BeginWrite()
	p.AmountCost += amount
	p.Balance -= amount
	p.EndWrite()
	p.Save()
}

func (p *DBPlayer)EarnMoney(amount float64)  {
	p.BeginWrite()
	p.Balance += amount
	p.EndWrite()
	p.Save()
}

func (p *DBPlayer)Save(){
	if p.ID < 0 {
		panic("(p *DBPlayer)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBPlayer)Remove(){
	db.Delete(p)
}
