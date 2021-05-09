package mysqld

import (
	"fmt"
	"bestsell/common"
)

type DBNotice struct {
    MysqlModel
    State int
	Title string
	Content string `gorm:"type:text;"`
}

var dbNoticeSafeMap common.SafeMap

func GetDBNotice(id int)*DBNotice {
	ret := dbNoticeSafeMap.Get(id)
	if ret == nil {
		return nil
	}
	return ret.(*DBNotice)
}

func AddDBNotice(item *DBNotice)  {
	old := GetDBNotice(item.ID)
	if old != nil {
		fmt.Println("DBNotice repeat")
		return
	}
	dbNoticeSafeMap.Set(item.ID, item)
}

func GetDBNoticeList()[]*DBNotice{
	var items []*DBNotice
	iterFunc := func(key int, v interface{}) bool{
		item := v.(*DBNotice)
		items = append(items, item)
		return true
	}
	dbNoticeSafeMap.RangeSafe(iterFunc)
	return items
}

func loadDBNotices() {
	var _dbItemsSlice []*DBNotice
	db.Find(&_dbItemsSlice)
	dbNoticeSafeMap = *common.NewSafeMap()
	for _,item := range _dbItemsSlice{
		if item.State == 1 {
			dbNoticeSafeMap.Set(item.ID, item)
		}
	}
}

func GetDBNoticeCount()int {
	count := 0
	iterFunc := func(key int, v interface{}) bool{
		count++
		return true
	}
	dbNoticeSafeMap.RangeSafe(iterFunc)
	return count
}

func startDBNotice()  {
	if !db.HasTable(&DBNotice{}) {
		db.CreateTable(&DBNotice{})
	}
	loadDBNotices()
	count := GetDBNoticeCount()
	fmt.Println("startDBNotice count =", count)
	if count == 0 {
		dbNotice := &DBNotice{}
		dbNotice.State = 1
		dbNotice.Title = "商城新开张，优惠多多，戳戳戳我看详情。"
		dbNotice.Content = `<html>
<head><meta charset="utf-8"></head>
<body>
    <p>尊敬的客户：</p>
    <p>商场新开张，点击产品看优惠！新商城暂不太完善，具体产品和优惠也可咨询客服</p>
</body>
</html>`
		dbNotice.Insert()
		dbNotice.Load()
	}
}

//DBNotice
func (p *DBNotice)Insert()  {
	if p.ID < 0 {
		panic("(p *DBNotice)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBNotice)Load(){
	if p.ID < 0 {
		panic("(p *DBNotice)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBNotice)Save(){
	if p.ID < 0 {
		panic("(p *DBNotice)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBNotice)Remove(){
	db.Delete(p)
}
