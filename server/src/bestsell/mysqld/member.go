package mysqld

type DBMember struct {
    MysqlModel
	TeamId  	int
    Level  		int
	Post  		int
	MemberId 	int
	Nickname 	string
}

func startDBMember()  {
	if !db.HasTable(&DBMember{}) {
		db.CreateTable(&DBMember{})
	}
}

//DBMemberBox
type DBMemberBox struct {
	MysqlModelBox
	TeamId int
	items []*DBMember
}
func (p *DBMemberBox)AddMember(item *DBMember) {
	item.TeamId = p.TeamId
	item.Insert()
	p.BeginWrite()
	p.items = append(p.items, item)
	p.EndWrite()
}
func (p *DBMemberBox)GetMember(id int)*DBMember {
	for _,item := range p.items {
		if item.ID == id {
			return item
		}
	}
	return nil
}
func (p *DBMemberBox)GetMembers()*[]*DBMember {
	return &p.items
}
func (p *DBMemberBox)RemoveMember(id int) {
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
func GetDBMemberBoxFromDB(teamId int)*DBMemberBox  {
	box := DBMemberBox{
		TeamId:teamId,
	}
	db.Where("team_id = ?", teamId).Find(&box.items)
	return &box
}

//DBMember
func (p *DBMember)Insert()  {
	if p.ID < 0 {
		panic("(p *DBMember)Insert p.ID < 0")
	}
	db.Create(p)
}

func (p *DBMember)Load(){
	if p.ID < 0 {
		panic("(p *DBMember)Load p.ID < 0")
	}
	db.First(p, p.ID)
}

func (p *DBMember)Save(){
	if p.ID < 0 {
		panic("(p *DBMember)Save p.ID < 0")
	}
	db.Save(p)
}

func (p *DBMember)Remove(){
	db.Delete(p)
}
