package mysqld

type DBShop struct {
    MysqlModel
	PlayerId int
}

func startDBShop()  {
	if !db.HasTable(&DBShop{}) {
		db.CreateTable(&DBShop{})
	}
}

//DBShopBox
type DBShopBox struct {
	MysqlModelBox
	PlayerId int
	items []*DBShop
}
func (p *DBShopBox)AddShop(item *DBShop) {
	item.PlayerId = p.PlayerId
	item.Insert()
	p.BeginWrite()
	p.items = append(p.items, item)
	p.EndWrite()
}
func (p *DBShopBox)GetShop(id int)*DBShop {
	for _,item := range p.items {
		if item.ID == id {
			return item
		}
	}
	return nil
}
func (p *DBShopBox)GetShops()*[]*DBShop {
	return &p.items
}
func (p *DBShopBox)RemoveShop(id int) {
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
func GetDBShopBoxFromDB(playerId int)*DBShopBox  {
	box := DBShopBox{
		PlayerId:playerId,
	}
	db.Where("player_id = ?", playerId).Find(&box.items)
	return &box
}


//DBShop
func (p *DBShop)Insert()  {
	if p.ID < 0 {
		panic("(p *DBShop)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBShop)Load(){
	if p.ID < 0 {
		panic("(p *DBShop)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBShop)Save(){
	if p.ID < 0 {
		panic("(p *DBShop)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBShop)Remove(){
	db.Delete(p)
}
