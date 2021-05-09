package mysqld

import (
	"fmt"
	"bestsell/common"
)

type DBAddress struct {
    MysqlModel
	PlayerId int
	LinkMan string `json:"linkMan"`
	Address string `json:"address"`
	Mobile string `json:"mobile"`
	Code string `json:"code"`
	IsDefault int `json:"isDefault"`
	ProvinceId int64 `json:"provinceId"`
	CityId int64 `json:"cityId"`
	AreaId int64 `json:"areaId"`
}

var dbAddressSafeMap common.SafeMap
func init() {
}

func GetDBAddress(id int)*DBAddress {
	ret := dbAddressSafeMap.Get(id)
	if ret == nil {
		return nil
	}
	return ret.(*DBAddress)
}

func AddDBAddress(item *DBAddress)  {
	old := GetDBAddress(item.ID)
	if old != nil {
		fmt.Println("DBAddress repeat")
		return
	}
	dbAddressSafeMap.Set(item.ID, item)
}

func LoadDBAddresss()  {
	var _dbItemsSlice []*DBAddress
	db.Find(&_dbItemsSlice)
	for _,item := range _dbItemsSlice{
		dbAddressSafeMap.Set(item.ID, item)
	}
}

func startDBAddress()  {
	if !db.HasTable(&DBAddress{}) {
		db.CreateTable(&DBAddress{})
	}
	//dbAddressSafeMap = *common.NewSafeMap()
	//LoadDBAddresss()
}

//DBAddressBox
type DBAddressBox struct {
	MysqlModelBox
	PlayerId int
	addresss []*DBAddress
}

func (p *DBAddressBox)AddAddress(item *DBAddress) {
	item.PlayerId = p.PlayerId
	item.Insert()
	p.BeginWrite()
	p.addresss = append(p.addresss, item)
	p.EndWrite()
}

func (p *DBAddressBox)GetAddress(id int)*DBAddress {
	for _,item := range p.addresss {
		if item.ID == id {
			return item
		}
	}
	return nil
}

func (p *DBAddressBox)GetDefaultAddress()*DBAddress {
	if len(p.addresss) == 0 {
		return nil
	}
	for _,item := range p.addresss {
		if item.IsDefault == 1 {
			return item
		}
	}
	return p.addresss[0]
}

func (p *DBAddressBox)GetAddresss()*[]*DBAddress {
	return &p.addresss
}

func (p *DBAddressBox)RemoveAddress(id int) {
	for idx,item := range p.addresss {
		if item.ID == id {
			p.BeginWrite()
			p.addresss = append(p.addresss[:idx], p.addresss[idx+1:]...)
			p.EndWrite()
			item.Remove()
			return
		}
	}
}

func GetDBAddressBoxFromDB(playerId int)*DBAddressBox  {
	addressBox := DBAddressBox{
		PlayerId:playerId,
	}
	db.Where("player_id = ?", playerId).Find(&addressBox.addresss)
	return &addressBox
}


//DBAddress
func (p *DBAddress)Insert()  {
	if p.ID < 0 {
		panic("(p *DBAddress)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBAddress)Load(){
	if p.ID < 0 {
		panic("(p *DBAddress)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBAddress)Save(){
	if p.ID < 0 {
		panic("(p *DBAddress)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBAddress)Remove(){
	db.Delete(p)
}
