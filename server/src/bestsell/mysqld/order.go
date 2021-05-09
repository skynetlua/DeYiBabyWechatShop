package mysqld

import (
	"bestsell/common"
	"bestsell/config"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type EStatusType int
const (
	EStatusPay      		EStatusType = 0
	EStatusSend  			EStatusType = 1
	EStatusReceive  		EStatusType = 2
	EStatusRepute   		EStatusType = 3
	EStatusFinish   		EStatusType = 4
	EStatusRefundApply   	EStatusType = 5
	EStatusRefundRefuse  	EStatusType = 6
	EStatusClose   			EStatusType = -1
	EStatusRefundFinish  	EStatusType = -2
	EStatusHide   			EStatusType = -3
	EStatusCancel   		EStatusType = -4
)

type DBOrderGoodsInfo struct {
	GoodsId 	int
	Name   		string
	Number 		int
	Price       int
	Pic  		string
	Tag 		int

	SkuId 		int
	SkuName     string
	SkuPrice    int

	Amount      int
	CartId      int
}

func (p *DBOrderGoodsInfo)GetPublicPic() string {
	return addPublicUrlHost(p.Pic)
}

func (p *DBOrderGoodsInfo)GetRealPrice() int {
	if p.SkuPrice > 0 {
		return p.SkuPrice
	}
	return p.Price
}

type DBOrder struct {
    MysqlModel
	OrderNumber   	string `gorm:"not null;unique"`
	PlayerId  		int
	Status  		int
	IsTips   		int
	TimeStamp  		int
	InviterId   	int
	CouponId  		int

	AmountGoods  	int
	AmountLogistics int
	AmountCoupon  	int
	AmountReal  	int
	AmountPay  	    int
	AmountPayerPay  int

	GoodsInfos  	string `gorm:"type:text"`
	GoodsId  		int

	SendType  		int
	QuickBuy 		int
	TeamBuy 		int
	Tag 			int

	AreaId  		int64
	Address  		string
    LinkMan  		string
	Mobile   		string
	Remark   		string
	RefundId   		int
	RefundStatus  	int
	TransactionId   string

	goodsInfos  	*[]*DBOrderGoodsInfo
	refund          *DBRefund
	myCoupon 		*DBMyCoupon
    hasTimer  		bool
}

type TeamBuyOrder struct {
	OrderId 	int
	PlayerId  	int
	PlayerName  string
	PlayerIcon  string
	TeamBuy 	int
	Status 		int
	EndTime 	int64
}

type TeamBuyOrderGroup struct {
	GoodsId int
	list []*TeamBuyOrder
	lastTime int64
}

var dbOrderSafeMap common.SafeMap

var teamBuyOrderSliceMap map[int]*TeamBuyOrderGroup
var orderLockMutex sync.RWMutex

func GetGoodsTeamBuyOrders(goodsId int) []*TeamBuyOrder {
	curTime := time.Now().Unix()

	orderLockMutex.Lock()
	tbOrderGroup, ok := teamBuyOrderSliceMap[goodsId]
	if ok {
		//if curTime < tbOrderGroup.lastTime+30 {
		list := tbOrderGroup.list
		orderLockMutex.Unlock()
		return list
		//}
	}
	orderLockMutex.Unlock()

	dbOrders := []*DBOrder{}
	db.Where("goods_id = ? and Status >= 1 and Tag = 3", goodsId).Find(&dbOrders)

	list := []*TeamBuyOrder{}
	for _, dbOrder := range dbOrders {
		dbPlayer := GetDBPlayerByPlayerIdOrFromDB(dbOrder.PlayerId)
		if dbPlayer != nil {
			dbPlayerInfo := dbPlayer.GetPlayerInfo()
			tbOrder := TeamBuyOrder {
				OrderId:	dbOrder.ID,
				PlayerId: 	dbOrder.PlayerId,
				PlayerName: dbPlayerInfo.NickName,
				PlayerIcon: dbPlayerInfo.AvatarUrl,
				EndTime: 	dbOrder.GetTagEndTime(),
				TeamBuy: 	dbOrder.TeamBuy,
				Status: 	dbOrder.Status,
			}
			list = append(list, &tbOrder)
		}
	}
	tbOrderGroup = &TeamBuyOrderGroup{
		GoodsId: goodsId,
		lastTime: curTime,
		list: list,
	}
	orderLockMutex.Lock()
	teamBuyOrderSliceMap[goodsId] = tbOrderGroup
	orderLockMutex.Unlock()

	return list
}

func AddGoodsTeamBuyOrder(goodsId int, teamBuyOrder *TeamBuyOrder) {
	orderLockMutex.Lock()
	tbOrderGroup, ok := teamBuyOrderSliceMap[goodsId]
	orderLockMutex.Unlock()
	if !ok {
		GetGoodsTeamBuyOrders(goodsId)

		orderLockMutex.Lock()
		tbOrderGroup, ok = teamBuyOrderSliceMap[goodsId]
		orderLockMutex.Unlock()
		if !ok {
			fmt.Println("AddGoodsTeamBuyOrder big error")
			return
		}
	}

	if tbOrderGroup.GoodsId != goodsId {
		panic("AddGoodsTeamBuyOrder tbOrderGroup.GoodsId != goodsId")
	}

	orderLockMutex.Lock()
	for idx, item := range tbOrderGroup.list {
		if item.OrderId == teamBuyOrder.OrderId {
			tbOrderGroup.list = append(tbOrderGroup.list[0:idx], tbOrderGroup.list[idx+1:]...)
			break
		}
	}
	tbOrderGroup.list = append(tbOrderGroup.list, teamBuyOrder)
	orderLockMutex.Unlock()
}

func UpdateGoodsTeamBuyOrderByOrder(dbOrder *DBOrder) {
	if dbOrder.Tag != int(EGoodsMarkTeam) {
		return
	}
	if dbOrder.TeamBuy < 2 {
		return
	}
	//dbGoods := GetDBGoods(dbOrder.GoodsId)
	//if dbGoods == nil || dbGoods.GetTag() != int(EGoodsMarkTeam) {
	//	return
	//}
	teamBuyOrders := GetGoodsTeamBuyOrders(dbOrder.GoodsId)
	dbPlayer := GetDBPlayerByPlayerIdOrFromDB(dbOrder.PlayerId)
	if dbPlayer == nil {
		return
	}
	dbPlayerInfo := dbPlayer.GetPlayerInfo()
	if dbOrder.TeamBuy == 2 {
		var mainTeamOrder *TeamBuyOrder = nil
		for _, item := range teamBuyOrders {
			if item.OrderId == dbOrder.ID {
				mainTeamOrder = item
				break
			}
		}
		if mainTeamOrder == nil {
			mainTeamOrder = &TeamBuyOrder {
				OrderId:	dbOrder.ID,
				PlayerId: 	dbOrder.PlayerId,
				PlayerName: dbPlayerInfo.NickName,
				PlayerIcon: dbPlayerInfo.AvatarUrl,
				EndTime: 	dbOrder.GetTagEndTime(),
				TeamBuy: 	dbOrder.TeamBuy,
				Status: 	dbOrder.Status,
			}
			AddGoodsTeamBuyOrder(dbOrder.GoodsId, mainTeamOrder)
		} else {
			orderLockMutex.Lock()

			mainTeamOrder.EndTime = dbOrder.GetTagEndTime()
			mainTeamOrder.TeamBuy = dbOrder.TeamBuy
			mainTeamOrder.Status = dbOrder.Status

			orderLockMutex.Unlock()
		}
	} else {
		var subTeamOrder *TeamBuyOrder = nil
		var mainTeamOrder *TeamBuyOrder = nil
		for _, item := range teamBuyOrders {
			if item.OrderId == dbOrder.ID {
				subTeamOrder = item
			} else if item.OrderId == dbOrder.TeamBuy {
				mainTeamOrder = item
			}
		}

		dbOrder.SetDBOrderStatus(EStatusReceive)
		dbOrder.Save()

		if mainTeamOrder == nil {
			fmt.Println("AddGoodsTeamBuyOrder lost mainTeamOrder")
			dbOrder.TeamBuy = 2
			UpdateGoodsTeamBuyOrderByOrder(dbOrder)
			return
		}

		orderLockMutex.Lock()
		mainTeamOrder.Status = dbOrder.Status
		orderLockMutex.Unlock()

		mainOrder := &DBOrder{}
		mainOrder.ID = dbOrder.TeamBuy
		mainOrder.Load()
		if mainOrder.PlayerId > 0 {
			dbMainPlayer := GetDBPlayerByPlayerIdOrFromDB(mainOrder.PlayerId)
			if dbMainPlayer != nil {
				mainOrder = dbMainPlayer.GetOrderBox().GetOrder(dbOrder.TeamBuy)
				if mainOrder != nil {
					mainOrder.SetDBOrderStatus(EStatusReceive)
					mainOrder.Save()
				}
			}
		}

		if subTeamOrder == nil {
			subTeamOrder = &TeamBuyOrder {
				OrderId:	dbOrder.ID,
				PlayerId: 	dbOrder.PlayerId,
				PlayerName: dbPlayerInfo.NickName,
				PlayerIcon: dbPlayerInfo.AvatarUrl,
				EndTime: 	dbOrder.GetTagEndTime(),
				TeamBuy: 	dbOrder.TeamBuy,
				Status: 	dbOrder.Status,
			}
			AddGoodsTeamBuyOrder(dbOrder.GoodsId, subTeamOrder)
		} else {
			orderLockMutex.Lock()

			subTeamOrder.EndTime = dbOrder.GetTagEndTime()
			subTeamOrder.TeamBuy = dbOrder.TeamBuy
			subTeamOrder.Status = dbOrder.Status

			orderLockMutex.Unlock()
		}
	}
}

func init() {
	teamBuyOrderSliceMap = make(map[int]*TeamBuyOrderGroup)
}

func GetDBOrder(id int)*DBOrder {
	ret := dbOrderSafeMap.Get(id)
	if ret == nil {
		return nil
	}
	return ret.(*DBOrder)
}
func AddDBOrder(item *DBOrder) {
	old := GetDBOrder(item.ID)
	if old != nil {
		fmt.Println("DBOrder repeat")
		return
	}
	dbOrderSafeMap.Set(item.ID, item)
}
func LoadDBOrders() {
	var _dbItemsSlice []*DBOrder
	db.Find(&_dbItemsSlice)
	for _,item := range _dbItemsSlice {
		dbOrderSafeMap.Set(item.ID, item)
	}
}
func startDBOrder() {
	if !db.HasTable(&DBOrder{}) {
		db.CreateTable(&DBOrder{})
	}
	dbOrderSafeMap = *common.NewSafeMap()
	//LoadDBOrders()
}
func GetOrdersByStatus(status int)*[]*DBOrder {
	var orders []*DBOrder
	db.Where("status = ?", status).Find(&orders)
	return &orders
}

//DBOrderBox
type DBOrderBox struct {
	MysqlModelBox
	PlayerId int
	orders []*DBOrder
}
func (p *DBOrderBox)AddOrder(item *DBOrder) {
	item.PlayerId = p.PlayerId
	item.Insert()

	p.BeginWrite()
	defer p.EndWrite()

	p.orders = append(p.orders, item)
}
func (p *DBOrderBox)GetOrder(id int)*DBOrder {
	p.BeginWrite()
	defer p.EndWrite()

	for _,item := range p.orders {
		if item.ID == id {
			return item
		}
	}
	return nil
}
func (p *DBOrderBox)GetOrderByOrderNumber(orderNumber string)*DBOrder {
	if len(orderNumber) < 13 {
		return nil
	}
	idx := strings.Index(orderNumber, "-")
	orderNumber = orderNumber[idx+1:]
	idx = strings.Index(orderNumber, "-")
	if idx > 0 {
		orderNumber = orderNumber[:idx]
	}
	orderId, _ := strconv.Atoi(orderNumber)
	return p.GetOrder(orderId)
}

func (p *DBOrderBox)GetOrders()[]*DBOrder {
	return p.orders
}

func (p *DBOrderBox)RemoveOrder(id int) {
	p.BeginWrite()
	defer p.EndWrite()

	for idx,item := range p.orders {
		if item.ID == id {
			p.orders = append(p.orders[:idx], p.orders[idx+1:]...)
			//item.Remove()
			return
		}
	}
}
func GetDBOrderBoxFromDB(playerId int) *DBOrderBox {
	box := DBOrderBox{
		PlayerId:playerId,
	}
	dbOrders := []*DBOrder{}
	db.Where("player_id = ?", playerId).Find(&dbOrders)

	for _, dbOrder := range dbOrders {
		dbOrder.CheckStatus()
		if dbOrder.IsOrderDelete() {
			continue
		}
		box.orders = append(box.orders, dbOrder)
	}
	return &box
}

func GetDBOrderByOpenIdAndOrderNumber(openId string, tradeNo string) *DBOrder {
	player := GetDBPlayerByOpenIdOrFromDB(openId)
	if player == nil {
		fmt.Println("GetDBOrderByOpenIdAndOrderNumber player == nil, openId =", openId)
		return nil
	}
	dbOrderBox := player.GetOrderBox()
	dbOrder := dbOrderBox.GetOrderByOrderNumber(tradeNo)
	if dbOrder == nil {
		fmt.Println("GetDBOrderByOpenIdAndOrderNumber dbOrder == nil openId =", openId, "tradeNo =", tradeNo)
		return nil
	}
	return dbOrder
}

// custom
func (p *DBOrder)GetOrCreateMyCounpon() *DBMyCoupon {
	if p.myCoupon == nil || p.CouponId != p.myCoupon.ID {
		dbMyCoupon := &DBMyCoupon {
			PlayerId: p.PlayerId,
			OrderId: p.ID,
		}
		dbMyCoupon.LoadWithOrderId()
		fmt.Println("GetOrCreateCounpon dbMyCoupon1 =", dbMyCoupon)
		if dbMyCoupon.ID == 0 {
			dbMyCoupon.Insert()
			dbMyCoupon.LoadWithOrderId()
			fmt.Println("GetOrCreateCounpon dbMyCoupon2 =", dbMyCoupon)
			if dbMyCoupon.ID == 0 {
				return nil
			}
		}
		fmt.Println("GetOrCreateCounpon dbMyCoupon3 =", dbMyCoupon)
		p.myCoupon = dbMyCoupon
		p.CouponId = dbMyCoupon.ID
		p.AmountCoupon = dbMyCoupon.Amount
	}
	fmt.Println("GetOrCreateCounpon dbMyCoupon =", p.myCoupon)
	fmt.Println("GetOrCreateCounpon p =", p)
	return p.myCoupon
}

//DBOrder
func (p *DBOrder)GetOrderNumber() string {
	if p.ID == 0 {
		return ""
	}
	if len(p.OrderNumber) == 0 {
		fmt.Println("(p *DBOrder)GetOrderNumber len(p.OrderNumber) == 0, ID =", p.ID)
		panic("(p *DBOrder)GetOrderNumber len(p.OrderNumber) == 0")
	}
	return p.OrderNumber
}

func (p *DBOrder)GetPayOrderNumber() string {
	orderNumber := p.GetOrderNumber()+"-"+strconv.Itoa(p.AmountReal)
	return orderNumber
}

//20 06 0102 150405 344
func (p *DBOrder)GetOrderGoodsInfos() []*DBOrderGoodsInfo {
	if p.goodsInfos != nil {
		return *p.goodsInfos
	}
	var goodsInfos []*DBOrderGoodsInfo
	err := json.Unmarshal([]byte(p.GoodsInfos), &goodsInfos)
	if err != nil {
		return nil
	}
	p.BeginWrite()
	defer p.EndWrite()

	p.goodsInfos = &goodsInfos
	return *p.goodsInfos
}

func (p *DBOrder)SetGoodsInfos(goodsInfos *[]*DBOrderGoodsInfo) {
	p.goodsInfos = goodsInfos
}

func (p *DBOrder)SetOrderGoodsInfos(goodsInfos *[]*DBOrderGoodsInfo) bool {
	if goodsInfos == nil || len(*goodsInfos) == 0 {
		return true
	}
	infoItems, err := json.Marshal(goodsInfos)
	if err != nil {
		fmt.Println("(p *DBOrder)SetOrderGoodsInfos error:", err)
		return false
	}
	p.BeginWrite()
	defer p.EndWrite()

	p.GoodsInfos = string(infoItems)
	p.goodsInfos = goodsInfos
	return true
}

func (p *DBOrder)GetOrderGoodsInfo(goodsId int, skuId int) *DBOrderGoodsInfo {
	goodsInfos := p.GetOrderGoodsInfos()
	for _,item := range goodsInfos {
		if item.GoodsId == goodsId && (skuId == -1 || item.SkuId == skuId) {
			return item
		}
	}
	return nil
}

func (p *DBOrder)GetDBOrderStatus() int {
	if p.RefundStatus != 0 {
		if p.Status < 0 {
			return p.Status
		}
		return p.RefundStatus
	}
	return p.Status
}

func (p *DBOrder)SetDBOrderStatus(status EStatusType) {
	p.BeginWrite()
	defer p.EndWrite()

	p.Status = int(status)
}

func (p *DBOrder)GetDBOrderStatusName() string {
	switch EStatusType(p.GetDBOrderStatus()) {
	case EStatusPay:
		return "待支付"
	case EStatusSend:
		return "待发货"
	case EStatusReceive:
		return "待收货"
	case EStatusRepute:
		return "待评价"
	case EStatusFinish:
		return "订单完成"
	case EStatusRefundApply:
		return "退款审核"
	case EStatusRefundRefuse:
		return "退款拒绝"
	case EStatusClose:
		return "订单关闭"
	case EStatusRefundFinish:
		return "退款完成"
	}
	return "未知状态"
}

func (p *DBOrder)GetDBRefund() *DBRefund {
	if p.refund == nil {
		p.BeginWrite()
		p.refund = &DBRefund{
			PlayerId: p.PlayerId,
			OrderId: p.ID,
		}
		p.EndWrite()
	}
	refund := p.refund
	if refund.ID <= 0 {
		refund.LoadWithOrderId()
		if refund.ID == 0 {
			p.BeginWrite()
			refund.ID = -1
			p.RefundId = 0
			p.RefundStatus = 0
			p.EndWrite()
		} else {
			if refund.ID > 0 {
				if p.RefundId != refund.ID {
					p.BeginWrite()
					p.RefundId = refund.ID
					p.EndWrite()
				}
			}
		}
	}
	if refund.ID > 0 {
		return refund
	}
	return nil
}

func (p *DBOrder)AddDBRefund(refund *DBRefund) {
	refund.OrderId = p.ID
	refund.Insert()

	p.BeginWrite()
	p.RefundId = refund.ID
	p.RefundStatus = int(EStatusRefundApply)
	p.EndWrite()

	p.Save()
}

func (p *DBOrder)RemoveDBRefund() {
	dbRefund := p.GetDBRefund()
	if dbRefund != nil && dbRefund.ID > 0 {
		dbRefund.Remove()
		p.BeginWrite()
		dbRefund.ID = -1
		p.EndWrite()
	}

	p.BeginWrite()
	p.RefundId = 0
	p.RefundStatus = 0
	p.EndWrite()

	p.Save()
}

func (p *DBOrder)Insert() {
	if p.ID != 0 {
		panic("(p *DBOrder)Insert p.ID < 0")
	}
	dbPlayer := GetDBPlayerByPlayerIdOrFromDB(p.PlayerId)
	if dbPlayer == nil {
		fmt.Println("(p *DBOrder)Insert dbPlayer == nil p.PlayerId =", p.PlayerId)
		panic("(p *DBOrder)Insert dbPlayer == nil")
	}
	var lastOrder DBOrder
	lastOrder.ID = 0
	db.Last(&lastOrder)
	orderId := lastOrder.ID+1
	p.ID = orderId

	tm := time.Unix(int64(p.TimeStamp),0)
	orderNumber := tm.Format("20060102150405")
	orderNumber = strconv.Itoa(p.PlayerId)+orderNumber[2:]+"-"+strconv.Itoa(p.ID)
	p.OrderNumber = orderNumber
	db.Create(p)
}

func (p *DBOrder)Load() {
	if p.ID < 0 {
		panic("(p *DBOrder)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBOrder)Save() {
	if p.ID < 0 {
		panic("(p *DBOrder)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBOrder) DelaySave(child interface{}) {
	fmt.Println("(p *DBOrder) DelaySave")
	panic("(p *DBOrder) DelaySave")
}

func (p *DBOrder)Update(fields []interface{}) {
	if p.ID < 0 {
		panic("(p *DBOrder)Save p.ID < 0")
	}
	db.Model(p).Where("id = ?", p.ID).Update(fields...)
}

func (p *DBOrder)Remove() {
	db.Delete(p)
}

func (p *DBOrder)GetGoodsNumber() int {
	goodsInfos := p.GetOrderGoodsInfos()
	goodsNumber := 0
	for _,goodsInfo := range goodsInfos {
		goodsNumber += goodsInfo.Number
	}
	return goodsNumber
}

func (p *DBOrder)CalGoodsAmount() int {
	p.AmountGoods = 0
	goodsInfos := p.GetOrderGoodsInfos()
	for _,goodsInfo := range goodsInfos {
		dbGoods := GetDBGoods(goodsInfo.GoodsId)
		if dbGoods == nil {
			return -1
		}
		if goodsInfo.Number > dbGoods.NumberStore {
			return -2
		}
		amount := goodsInfo.GetRealPrice()*goodsInfo.Number
		goodsInfo.Amount = amount
		p.AmountGoods += amount
	}
	return 0
}

func (p *DBOrder)CalAmountReal() {
	p.BeginWrite()
	p.AmountReal = p.AmountGoods+p.AmountLogistics-p.AmountCoupon
	p.EndWrite()
}

func (p *DBOrder)SetAmountCoupon(amountCoupon int) {
	p.BeginWrite()
	p.AmountCoupon = amountCoupon
	p.EndWrite()
	p.CalAmountReal()
}

func (p *DBOrder)UpdateTips() {
	if p.IsTips != 0 {
		//go func() {
		p.BeginWrite()
		p.IsTips = 0
		p.EndWrite()
		p.Save()
		//}()
	}
}

func (p *DBOrder)GetTagEndTime() int64 {
	if p.Tag == int(EGoodsMarkSeckill) {
		for {
			goodsInfos := p.GetOrderGoodsInfos()
			if len(goodsInfos) == 0 {
				break
			}
			goodsInfo := goodsInfos[0]
			dbGoods := GetDBGoods(goodsInfo.GoodsId)
			if dbGoods == nil {
				break
			}
			if dbGoods.GetTag() != p.Tag {
				break
			}
			return dbGoods.EndTime
		}
	} else if p.Tag == int(EGoodsMarkTeam) {
		return int64(p.TimeStamp+60*60*24)
	}
	p.hasTimer = false
	return 0
}
//func (p *DBOrder)IsOrderWaitPay() bool {
//	if p.hasTimer {
//		if p.GetDBOrderStatus() != int(EStatusPay) {
//			p.SetDBOrderStatus(EStatusCancel)
//			p.Save()
//			return false
//		}
//		return true
//	}
//	return false
//}
func (p *DBOrder)CheckStatus() {
	p.hasTimer = false
	status := p.GetDBOrderStatus()
	if (status == int(EStatusPay) && p.Tag == int(EGoodsMarkSeckill)) || (status == int(EStatusSend) && p.Tag == int(EGoodsMarkTeam)) {
		for {
			goodsInfos := p.GetOrderGoodsInfos()
			if len(goodsInfos) == 0 {
				fmt.Println("(p *DBOrder)ChechStatus order len(goodsInfos) == 0 ")
				break
			}
			goodsInfo := goodsInfos[0]
			dbGoods := GetDBGoods(goodsInfo.GoodsId)
			if dbGoods == nil {
				fmt.Println("(p *DBOrder)ChechStatus order dbGoods == nil ")
				break
			}
			if dbGoods.GetTag() != p.Tag {
				fmt.Println("(p *DBOrder)ChechStatus order dbGoods.GetTag() != p.Tag")
				break
			}

			tagEndTime := p.GetTagEndTime()
			//if p.Tag == int(EGoodsMarkTeam) {
			//	if curTime > int64(p.TimeStamp+60*60*24) {
			//		break
			//	}
			//	p.hasTimer = true
			//	return
			//}
			if time.Now().Unix() > tagEndTime {
				fmt.Println("(p *DBOrder)ChechStatus order timeout")
				break
			}
			p.hasTimer = true
			return
		}
		p.SetDBOrderStatus(EStatusCancel)
		p.Save()
		return
	}
	if status != int(EStatusPay) {
		return
	}
	if time.Now().Unix() > int64(p.TimeStamp+60*60*3) {
		p.SetDBOrderStatus(EStatusClose)
		p.Save()
	}
}

func (p *DBOrder)FinishPay(amountPay int, amountPayerPay int, transactionId string) {
	p.BeginWrite()
	p.AmountPay = amountPay
	p.AmountPayerPay = amountPayerPay
	p.TransactionId = transactionId
	p.IsTips = 1
	p.Status = int(EStatusSend)
	p.EndWrite()

	dbPlayer := GetDBPlayerByPlayerId(p.PlayerId)
	goodsInfos := p.GetOrderGoodsInfos()
	for _,goodsInfo := range goodsInfos {
		dbGoods := GetDBGoods(goodsInfo.GoodsId)
		if dbGoods != nil {
			dbGoods.AddOrderNum()
			dbGoods.AddSellNum(goodsInfo.Number)
			dbGoods.DelaySave(dbGoods)
			if dbPlayer != nil {
				MakeVisitGoodsStat(EGoodsActionPay, dbPlayer, dbGoods, goodsInfo.SkuId, p.ID, goodsInfo.Number)
			}
		}
	}
	p.Save()

	UpdateGoodsTeamBuyOrderByOrder(p)
	//p.DoTeamBuy()
}

func (p *DBOrder)SendOrder() {
	if p.GetDBOrderStatus() != int(EStatusSend) {
		return
	}
	p.SetDBOrderStatus(EStatusReceive)
	p.Save()

	UpdateGoodsTeamBuyOrderByOrder(p)
}

func (p *DBOrder)CloseOrder() bool {
	status := p.GetDBOrderStatus()
	if status == int(EStatusPay) {
		p.SetDBOrderStatus(EStatusClose)
		p.Save()
	} else if status == int(EStatusRepute) || status == int(EStatusFinish) {
		p.SetDBOrderStatus(EStatusHide)
		p.Save()
	} else {
		return false
	}
	return true
}

func (p *DBOrder)DeleteOrder() bool {
	status := p.GetDBOrderStatus()
	if status == int(EStatusClose) {
		p.SetDBOrderStatus(EStatusCancel)
		p.Save()
	} else if status == int(EStatusRepute) || status == int(EStatusFinish) {
		p.SetDBOrderStatus(EStatusHide)
		p.Save()
	} else {
		return false
	}
	return true
}

func (p *DBOrder)IsOrderDelete() bool {
	status := p.GetDBOrderStatus()
	if status <= int(EStatusHide) || status <= int(EStatusCancel) {
		return true
	}
	return false
}

func (p *DBOrder) GetOrderAttach() string {
	amountGoodsStr := strconv.Itoa(p.AmountGoods)
	amountLogisticsStr := strconv.Itoa(p.AmountLogistics)
	amountCouponStr := strconv.Itoa(p.AmountCoupon)
	amountRealStr := strconv.Itoa(p.AmountReal)
	desc := p.GetOrderNumber()+"-"+strconv.Itoa(p.PlayerId)+"-"+strconv.Itoa(p.TimeStamp)+"-"+amountGoodsStr+"-"+amountLogisticsStr+"-"+amountCouponStr+"-"+amountRealStr+"\n"
	goodsInfos := p.GetOrderGoodsInfos()
	for _,goodsInfo := range goodsInfos {
		sellPrice := goodsInfo.GetRealPrice()
		if len(goodsInfo.SkuName) > 0 {
			desc = desc+strconv.Itoa(goodsInfo.GoodsId)+"["+strconv.Itoa(goodsInfo.SkuId)+"]"+"-"+strconv.Itoa(goodsInfo.Number)+"-"+strconv.Itoa(sellPrice)+"\n"
		} else {
			desc = desc+strconv.Itoa(goodsInfo.GoodsId)+"-"+strconv.Itoa(goodsInfo.Number)+"-"+strconv.Itoa(sellPrice)+"\n"
		}
	}
	if len(desc) > 120 {
		return desc[0:120]
	}
	return desc
}

func (p *DBOrder) GetOrderGoodsTag() string {
	return strconv.Itoa(p.ID)+"-"+strconv.Itoa(p.PlayerId)
}

func (p *DBOrder)GetOrderInfo() *map[string]interface{} {
	p.CalAmountReal()
	goodsDatas := p.GetGoodsInfoList()
	data := map[string]interface{}{
		"orderId":         	p.ID,
		"playerId": 	   	p.PlayerId,
		"orderNumber":     	p.GetOrderNumber(),
		"status":			p.Status,
		"sendType": 	   	p.SendType,
		"amountGoods":     	p.AmountGoods,
		"amountLogistics": 	p.AmountLogistics,
		"amountCoupon":    	p.AmountCoupon,
		"amountReal":      	p.AmountReal,
		"teamBuy":  		p.TeamBuy,
		"tag": 				p.Tag,
		"goodsList":       	goodsDatas,
	}
	return &data
}

func (p *DBOrder)GetItemInfo() *map[string]interface{} {
	p.CalAmountReal()
	tm := time.Unix(int64(p.TimeStamp),0)
	item := &map[string]interface{}{
		"id":			p.ID,
		"playerId": 	p.PlayerId,
		"sendType": 	p.SendType,
		"orderNumber":	p.GetOrderNumber(),
		"status":		p.Status,
		"amountReal":	p.AmountReal,
		"amountGoods":	p.AmountGoods,
		"goodsNumber":  p.GetGoodsNumber(),
		"refundId":		p.RefundId,
		"refundStatus": p.RefundStatus,
		"remark":		p.Remark,
		"statusStr":	p.GetDBOrderStatusName(),
		"timeStamp":	p.TimeStamp,
		"isTips":		p.IsTips,
		"teamBuy":  	p.TeamBuy,
		"tag": 			p.Tag,
		"dateAdd":		tm.Format("2006-01-02 15:04:05"),
	}
	if p.hasTimer {
		(*item)["endTime"] = p.GetTagEndTime()
	}
	return item
}

func (p *DBOrder)GetGoodsInfoList() *[]map[string]interface{} {
	var goodsList []map[string]interface{}
	goodsInfos := p.GetOrderGoodsInfos()
	for _,goodsInfo := range goodsInfos {
		//dbGoods := GetDBGoods(goodsInfo.GoodsId)
		//if dbGoods == nil {
		//	continue
		//}
		//dbGoodsInfo := dbGoods.GetGoodsInfo()
		goods := map[string]interface{} {
			"goodsId":goodsInfo.GoodsId,
			"name":goodsInfo.Name,
			"price":goodsInfo.Price,
			"number":goodsInfo.Number,
			"amount":goodsInfo.Amount,
			"realPrice":goodsInfo.GetRealPrice(),
			"pic":goodsInfo.GetPublicPic(),
			"skuId":goodsInfo.SkuId,
			"skuName":goodsInfo.SkuName,
			"skuPrice":goodsInfo.SkuPrice,
			"tag":goodsInfo.Tag,
		}
		goodsList = append(goodsList, goods)
	}
	return &goodsList
}

func (p *DBOrder)GetAddressInfo() *map[string]interface{} {
	logisticsData := map[string]interface{} {
		"linkMan":	p.LinkMan,
		"mobile":	p.Mobile,
		"address":	p.Address,
		"sendType": p.SendType,
	}
	areaId := int(p.AreaId)
	if areaId > 0 {
		var cfgAreaAddress *config.CfgAddress
		var cfgCityAddress *config.CfgAddress
		var cfgProvinceAddress *config.CfgAddress
		cfgAddress := config.GetCfgAddress(int(p.AreaId))
		if cfgAddress.Level == 3 {
			cfgAreaAddress = cfgAddress
			cfgCityAddress = config.GetCfgAddress(cfgAreaAddress.Pid)
			cfgProvinceAddress = config.GetCfgAddress(cfgCityAddress.Pid)
		}else{
			cfgAreaAddress = cfgAddress
			cfgCityAddress = cfgAddress
			cfgProvinceAddress = config.GetCfgAddress(cfgCityAddress.Pid)
		}
		logisticsData["provinceStr"] = cfgProvinceAddress.Name
		logisticsData["cityStr"] = cfgCityAddress.Name
		logisticsData["areaStr"] = cfgAreaAddress.Name
	}
	return &logisticsData
}

func (p *DBOrder)GetAllDetail() *map[string]interface{} {
	p.CalAmountReal()
	orderData := map[string]interface{}{
		"orderId":			p.ID,
		"playerId": 		p.PlayerId,
		"orderNumber":		p.GetOrderNumber(),
		"sendType": 	    p.SendType,
		"teamBuy":  		p.TeamBuy,
		"tag": 				p.Tag,
		"status":			p.Status,
		"statusStr":		p.GetDBOrderStatusName(),
		"amountGoods":		p.AmountGoods,
		"amountLogistics":	p.AmountLogistics,
		"amountCoupon":     p.AmountCoupon,
		"amountReal":		p.AmountReal,
		"refundId": 		p.RefundId,
		"refundStatus": 	p.RefundStatus,
	}
	if p.hasTimer {
		orderData["endTime"] = p.GetTagEndTime()
	}
	logisticsData := p.GetAddressInfo()
	goodsDatas := p.GetGoodsInfoList()
	data := map[string]interface{}{
		"orderInfo": orderData,
		"logistics": logisticsData,
		"goodsList": goodsDatas,
	}
	return &data
}