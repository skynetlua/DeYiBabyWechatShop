package mysqld

import (
	"bestsell/common"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"reflect"
	"sync/atomic"
	"time"
)

type MysqlModel struct {
	isSaving     	bool
	writeCounter 	int32
	ID           	int `gorm:"primary_key"`
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
	DeletedAt 		*time.Time `sql:"index"`
}

func (p *MysqlModel)GetUpdateDateStr()string {
	return p.UpdatedAt.Format("2006-01-02 15:04:05")
}

func (p *MysqlModel)GetCreateDateStr()string {
	return p.CreatedAt.Format("2006-01-02 15:04:05")
}

func (p *MysqlModel)BeginWrite() {
	if p.writeCounter == 0 {
		atomic.AddInt32(&p.writeCounter, 1)
	}else{
		if p.writeCounter != 0 {
			waitTime := time.Now().Unix()+3
			waitCount := 0
			for p.writeCounter != 0 {
				waitCount++
				if waitCount%10000 == 0 {
					if time.Now().Unix() > waitTime {
						panic("MysqlModel BeginWrite is lock")
					}
				}
			}
		}
		atomic.AddInt32(&p.writeCounter, 1)
	}
}

func (p *MysqlModel)EndWrite() {
	if p.writeCounter == 0 {
		panic("MysqlModel EndWrite is finish")
	}else{
		atomic.AddInt32(&p.writeCounter, -1)
	}
}

func (p *MysqlModel)Save(){
	panic("(p *MysqlModel)Save empty")
}

func (p *MysqlModel) DelaySave(child interface{}) {
	if p.isSaving {
		return
	}
	p.isSaving = true
	go func() {
		select {
		case <-time.After(time.Second *10):
			//p.Save()
			ref := reflect.ValueOf(child)
			method := ref.MethodByName("Save")
			if method.IsValid() {
				method.Call(make([]reflect.Value, 0))
			} else {
				panic("DelaySave")
			}
			p.isSaving = false
		}
	}()
}



type MysqlModelBox struct {
	isSaving     	bool
	writeCounter 	int32
}
func (p *MysqlModelBox)BeginWrite() {
	if p.writeCounter == 0 {
		atomic.AddInt32(&p.writeCounter, 1)
	}else{
		if p.writeCounter != 0 {
			waitTime := time.Now().Unix()+3
			waitCount := 0
			for p.writeCounter != 0 {
				waitCount++
				if waitCount%10000 == 0 {
					if time.Now().Unix() > waitTime {
						panic("MysqlModelBox BeginWrite is lock")
					}
				}
			}
		}
		atomic.AddInt32(&p.writeCounter, 1)
	}
}

func (p *MysqlModelBox)EndWrite() {
	if p.writeCounter == 0 {
		panic("MysqlModelBox EndWrite is finish")
	}else{
		atomic.AddInt32(&p.writeCounter, -1)
	}
}



//////////////////
var db *gorm.DB

func start() {
	config := common.Config
	if config.MySqlAddr == "" {
		return
	}
	fmt.Println("init mysql", config.MySqlAddr)
	var err error
	db, err = gorm.Open("mysql", config.MySqlAddr+"?charset=utf8&parseTime=true&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	err = db.DB().Ping()
	if err != nil {
		log.Fatal(err)
	}
	db = db.LogMode(config.Enable["debug"] == 1)
}

func StartServer(ch *chan bool) {
	if common.Config.Enable["mysql"] != 1 {
		return
	}
	go func() {
		fmt.Println("StartServer mysqld")
		start()
		startDBPlayerInfo()
		startDBPlayer()
		startDBCashLog()
		startDBRefund()
		startDBCategory()
		startDBGoodsInfo()
		startDBGoods()
		startDBMyCoupon()
		startDBOrder()
		startDBFavorite()
		startDBReputation()
		startDBCart()
		startDBAddress()
		startDBWithdrawLog()
		// startDBMyTeam()
		//startDBTeam()
		//startDBTeamInfo()
		//startDBMember()
		startDBCommission()
		startDBNotice()
		startDBGoodsStat()
		(*ch) <- true
	}()
	<-(*ch)
}