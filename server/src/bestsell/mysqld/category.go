package mysqld

import (
	"bestsell/common"
	"strconv"

	// "fmt"
)

type DBCategory struct {
    MysqlModel
	Name string `gorm:"not null;unique"`
    Level int
    Order int
    Status int
    Icon string
	pic string
	publicIcon string
	publicPic string
}

var dbCategorySafeMap common.SafeMap

func GetDBCategory(id int)*DBCategory {
	ret := dbCategorySafeMap.Get(id)
	if ret == nil {
		return nil
	}
	return ret.(*DBCategory)
}

func GetDBCategoryOrFromDB(id int)*DBCategory {
	item := GetDBCategory(id)
	if item != nil  {
		return item
	}
	item = &DBCategory{}
	item.ID = id
	item.Load()
	if len(item.Name) == 0 {
		return nil
	}
	AddDBCategory(item)
	return item
}

func AddDBCategory(item *DBCategory)  {
	old := GetDBCategory(item.ID)
	if old != nil {
		panic("AddDBCategory DBCategory repeat")
		return
	}
	dbCategorySafeMap.Set(item.ID, item)
}

func RemoveDBCategory(id int)  {
	dbCategorySafeMap.Remove(id)
}

func RemoveDBCategoryAndDB(id int)  {
	item := GetDBCategoryOrFromDB(id)
	if item == nil {
		return
	}
	RemoveDBCategory(id)
	item.Remove()
}

func LoadDBCategorys() {
	var _dbItemsSlice []*DBCategory
	db.Find(&_dbItemsSlice)
	dbCategorySafeMap = *common.NewSafeMap()
	for _,item := range _dbItemsSlice{
		dbCategorySafeMap.Set(item.ID, item)
		item.pic = "/api/category/pic"+strconv.Itoa(item.ID)+".png"
		item.GetPublicPic()
		item.GetPublicIcon()
	}
}

func startDBCategory() {
	if !db.HasTable(&DBCategory{}) {
		db.CreateTable(&DBCategory{})
	}
	LoadDBCategorys()
}

func GetDBCategoryList()[]*DBCategory{
	var items []*DBCategory
	iterFunc := func(key int, v interface{}) bool{
		items = append(items, v.(*DBCategory))
		return true
	}
	dbCategorySafeMap.RangeSafe(iterFunc)
	return items
}

func AddNewCategory(item *DBCategory)bool{
	if item.ID == 0 {
		item.Insert()
	}
	if item.ID < 0 {
		return false
	}
	//AddDBCategory(item)
	return true
}

func GetCategoryResIds()*[]string{
	var resIds []string
	folder := common.UploadCategoryPath
	if common.Exists(folder) {
		var filePaths []string
		filePaths,_ = common.GetAllFiles(folder, filePaths)
		for _,tmp := range filePaths {
			item := removeDir2Url(tmp)
			resIds = append(resIds, item)
		}
	}
	return &resIds
}

//DBCategory
func (p *DBCategory)GetPublicPic()string {
	if len(p.publicPic) > 0 {
		return p.publicPic
	}
	p.BeginWrite()
	p.publicPic = addPublicUrlHost(p.pic)
	p.EndWrite()
	return p.publicPic
}

func (p *DBCategory)GetPublicIcon()string {
	if len(p.publicIcon) > 0 {
		return p.publicIcon
	}
	p.BeginWrite()
	p.publicIcon = addPublicUrlHost(p.Icon)
	p.EndWrite()
	return p.publicIcon
}

func (p *DBCategory)SetIcon(icon string) {
	if len(icon) == 0 {
		return
	}
	p.BeginWrite()
	p.Icon = removeUrlHost(icon)
	p.publicIcon = ""
	p.EndWrite()
}

func (p *DBCategory)GetIcon()string {
	return p.Icon
}

func (p *DBCategory)Insert()  {
	if p.ID < 0 {
		panic("(p *DBCategory)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBCategory)Load(){
	if p.ID < 0 {
		panic("(p *DBCategory)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBCategory)Save(){
	if p.ID < 0 {
		panic("(p *DBCategory)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBCategory)Remove(){
	db.Delete(p)
}
