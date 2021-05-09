package mysqld

import (
	"fmt"
	"bestsell/common"
)

type DBCart struct {
    MysqlModel
    PlayerId int
	GoodsId int
	SkuId int
	NumberBuy int
}

var dbCartSafeMap common.SafeMap
func init() {
}

func GetDBCart(id int)*DBCart {
	ret := dbCartSafeMap.Get(id)
	if ret == nil {
		return nil
	}
	return ret.(*DBCart)
}

func AddDBCart(item *DBCart)  {
	old := GetDBCart(item.ID)
	if old != nil {
		fmt.Println("DBCart repeat")
		return
	}
	dbCartSafeMap.Set(item.ID, item)
}

func LoadDBCarts()  {
	var _dbItemsSlice []*DBCart
	db.Find(&_dbItemsSlice)
	for _,item := range _dbItemsSlice{
		dbCartSafeMap.Set(item.ID, item)
	}
}

func startDBCart()  {
	if !db.HasTable(&DBCart{}) {
		db.CreateTable(&DBCart{})
	}
	//dbCartSafeMap = *common.NewSafeMap()
	//LoadDBCarts()
}

type DBCartBox struct {
	MysqlModelBox
	PlayerId int
	carts []*DBCart
}

func (p *DBCartBox)AddCart(item *DBCart) {
	item.PlayerId = p.PlayerId
	item.Insert()

	p.BeginWrite()
	defer p.EndWrite()
	p.carts = append(p.carts, item)
}

func (p *DBCartBox)GetCart(id int)*DBCart {
	p.BeginWrite()
	defer p.EndWrite()

	for _,item := range p.carts {
		if item.ID == id {
			return item
		}
	}
	return nil
}

func (p *DBCartBox)GetCarts() []*DBCart {
	//p.BeginWrite()
	//defer p.EndWrite()
	return p.carts
}

func (p *DBCartBox)GetCartCount() int {
	p.BeginWrite()
	defer p.EndWrite()
	return len(p.carts)
}

func (p *DBCartBox)RemoveCart(id int) {
	p.BeginWrite()
	defer p.EndWrite()

	for idx,item := range p.carts {
		if item.ID == id {
			p.carts = append(p.carts[:idx], p.carts[idx+1:]...)
			item.Remove()
			return
		}
	}
}

func (p *DBCartBox)Clear() {
	p.BeginWrite()
	defer p.EndWrite()

	carts := p.carts
	p.carts = p.carts[0:0]
	for _,cart := range carts {
		cart.Remove()
	}
}

func GetDBCartBoxFromDB(playerId int)*DBCartBox  {
	box := DBCartBox{
		PlayerId:playerId,
	}
	db.Where("player_id = ?", playerId).Find(&box.carts)
	return &box
}


//DBCart
func (p *DBCart)Insert()  {
	if p.ID < 0 {
		panic("(p *DBCart)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBCart)Load(){
	if p.ID < 0 {
		panic("(p *DBCart)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBCart)Save(){
	if p.ID < 0 {
		panic("(p *DBCart)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBCart)Remove(){
	db.Delete(p)
}
