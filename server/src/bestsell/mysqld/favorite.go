package mysqld

import (
	"fmt"
	"bestsell/common"
)

type DBFavorite struct {
    MysqlModel
	PlayerId int
	GoodsId int `json:"goodsId"`
}

var dbFavoriteSafeMap common.SafeMap
func init() {
}

func GetDBFavorite(id int)*DBFavorite {
	ret := dbFavoriteSafeMap.Get(id)
	if ret == nil {
		return nil
	}
	return ret.(*DBFavorite)
}

func AddDBFavorite(item *DBFavorite)  {
	old := GetDBFavorite(item.ID)
	if old != nil {
		fmt.Println("DBFavorite repeat")
		return
	}
	dbFavoriteSafeMap.Set(item.ID, item)
}

func LoadDBFavorites()  {
	var _dbItemsSlice []*DBFavorite
	db.Find(&_dbItemsSlice)
	for _,item := range _dbItemsSlice{
		dbFavoriteSafeMap.Set(item.ID, item)
	}
}

func startDBFavorite()  {
	if !db.HasTable(&DBFavorite{}) {
		db.CreateTable(&DBFavorite{})
	}
	//dbFavoriteSafeMap = *common.NewSafeMap()
	//LoadDBFavorites()
}

//DBFavoriteBox
type DBFavoriteBox struct {
	MysqlModelBox
	PlayerId int
	favorites []*DBFavorite
}

func (p *DBFavoriteBox)AddFavorite(item *DBFavorite) {
	item.PlayerId = p.PlayerId
	item.Insert()
	p.BeginWrite()
	p.favorites = append(p.favorites, item)
	p.EndWrite()
}

func (p *DBFavoriteBox)GetFavorite(id int)*DBFavorite {
	for _,item := range p.favorites {
		if item.ID == id {
			return item
		}
	}
	return nil
}

func (p *DBFavoriteBox)GetFavoriteByGoodsId(goodsId int)*DBFavorite {
	for _,item := range p.favorites {
		if item.GoodsId == goodsId {
			return item
		}
	}
	return nil
}

func (p *DBFavoriteBox)GetFavorites()*[]*DBFavorite {
	return &p.favorites
}

func (p *DBFavoriteBox)RemoveFavorite(id int) {
	for idx,item := range p.favorites {
		if item.ID == id {
			p.BeginWrite()
			p.favorites = append(p.favorites[:idx], p.favorites[idx+1:]...)
			p.EndWrite()
			item.Remove()
			return
		}
	}
}
func (p *DBFavoriteBox)RemoveFavoriteByGoodsId(goodsId int) {
	for idx,item := range p.favorites {
		if item.GoodsId == goodsId {
			p.BeginWrite()
			p.favorites = append(p.favorites[:idx], p.favorites[idx+1:]...)
			p.EndWrite()
			item.Remove()
			return
		}
	}
}
func GetDBFavoriteBoxFromDB(playerId int)*DBFavoriteBox  {
	box := DBFavoriteBox{
		PlayerId:playerId,
	}
	db.Where("player_id = ?", playerId).Find(&box.favorites)
	return &box
}

//DBFavorite
func (p *DBFavorite)Insert()  {
	if p.ID < 0 {
		panic("(p *DBFavorite)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBFavorite)Load(){
	if p.ID < 0 {
		panic("(p *DBFavorite)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBFavorite)Save(){
	if p.ID < 0 {
		panic("(p *DBFavorite)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBFavorite)Remove(){
	db.Delete(p)
}