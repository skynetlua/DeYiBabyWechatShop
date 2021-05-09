package mysqld

import "bestsell/common"
import "fmt"

type DBRecommend struct {
    MysqlModel
	GoodsId int `gorm:"not null;unique"`
    Style   int
}

var dbRecommendSafeMap common.SafeMap

func GetDBRecommend(goodsId int)*DBRecommend {
	ret := dbRecommendSafeMap.Get(goodsId)
	if ret == nil {
		return nil
	}
	return ret.(*DBRecommend)
}

func AddDBRecommend(item *DBRecommend)  {
	old := GetDBRecommend(item.GoodsId)
	if old != nil {
		fmt.Println("DBRecommend repeat")
		return
	}
	dbRecommendSafeMap.Set(item.GoodsId, item)
}

func LoadDBRecommends()  {
	var _dbItemsSlice []*DBRecommend
	db.Find(&_dbItemsSlice)
	for _,item := range _dbItemsSlice{
		dbRecommendSafeMap.Set(item.GoodsId, item)
	}
}

func startDBRecommend()  {
	if !db.HasTable(&DBRecommend{}) {
		db.CreateTable(&DBRecommend{})
	}
	LoadDBRecommends()
}

//DBRecommend
func (p *DBRecommend)Insert()  {
	if p.ID < 0 {
		panic("(p *DBRecommend)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBRecommend)Load(){
	if p.ID < 0 {
		panic("(p *DBRecommend)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBRecommend)Save(){
	if p.ID < 0 {
		panic("(p *DBRecommend)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBRecommend)Remove(){
	db.Delete(p)
}
