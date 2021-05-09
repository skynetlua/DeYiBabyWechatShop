package module

import (
	"bestsell/common"
	"bestsell/mysqld"
	"fmt"
	"strings"
	"time"
)

type WxWatermark struct {
	Timestamp int  `json:"timestamp"`
	Appid string `json:"appid"`
}

type WxUserInfo struct {
	OpenId string `json:"openId"`
	NickName string `json:"nickName"`
	Gender int `json:"gender"`
	Language string `json:"language"`
	City string `json:"city"`
	Province string `json:"province"`
	Country string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
	Watermark WxWatermark `json:"watermark"`
	Referrer string `json:"referrer"`
}

func init() {
}

func GetPlayer(p *common.BSSession)*mysqld.DBPlayer {
	if p.Data == nil {
		return nil
	}
	ret := p.Data["player"]
	if ret == nil {
		return nil
	}
	return ret.(*mysqld.DBPlayer)
}

func OpenId2Token(openId string)string  {
	return openId
	//data := []byte(openId)
	//has := md5.Sum(data)
	//md5str1 := fmt.Sprintf("%x", has)
	//return md5str1
}

//func GetPlayerByOpenId(openId string) *mysqld.DBPlayer {
//	token := OpenId2Token(openId)
//	player := mysqld.GetDBPlayer(token)
//	if player != nil {
//		return player
//	}
//	player = &mysqld.DBPlayer{
//		Token: token,
//	}
//	player.LoadWithToken()
//	if player.ID <= 0 {
//		return nil
//	}
//	mysqld.AddDBPlayer(player)
//	return player
//}

func OnLogin(login *UserLogin) int{
	login.LoginTime = time.Now()
	if len(login.Openid) > 0 {
		login.Token = OpenId2Token(login.Openid)
	}
	if len(login.Token) > 0 {
		ret := login.LoginPlayer()
		if len(login.Openid) > 0 {
			player := login.Player
			if player != nil {
				player.UnionId = login.Unionid
				player.OpenId = login.Openid
				player.SessionKey = login.SessionKey
				player.DelaySave(player)
				return 0
			}
			player = &mysqld.DBPlayer{
				Token: login.Token,
				OpenId: login.Openid,
				SessionKey: login.SessionKey,
				UnionId: login.Unionid,
				LoginIP: login.Ip,
				Balance: 1000,
			}
			player.Insert()
			ret = login.LoginPlayer()
		}
		return ret
	}else{
		panic("OnLogin login.Token")
	}
	return -1
}

func OnLoginWx(login *WxUserInfo) int {
	var token = OpenId2Token(login.OpenId)
	userLogin := &UserLogin{}
	userLogin.Token = token
	userLogin.LoginPlayer()
	if userLogin.Player == nil {
		return -1
	}
	player := userLogin.Player
	if strings.Compare(token, player.Token) != 0 {
		fmt.Println("OnLoginWx error token =", token)
		fmt.Println("OnLoginWx error Account =", player.Token)
		return -1
	}
	playerInfo := player.GetPlayerInfo()
	playerInfo.NickName = login.NickName
	playerInfo.Gender = login.Gender
	playerInfo.Language = login.Language
	playerInfo.City = login.City
	playerInfo.Province = login.Province
	playerInfo.Country = login.Country
	playerInfo.AvatarUrl = login.AvatarUrl
	playerInfo.Timestamp = login.Watermark.Timestamp
	playerInfo.Referrer = login.Referrer
	playerInfo.DelaySave(playerInfo)
	return 0
}

////UserLogin
type UserLogin struct {
	Code string
	Token string
	Player *mysqld.DBPlayer
	LoginTime time.Time
	Ip string
	Openid string
	SessionKey string
	Unionid string
}

func (p *UserLogin)LoginPlayer() int {
	player := mysqld.GetDBPlayer(p.Token)
	if player != nil {
		p.Player = player
		player.LoginIP = p.Ip
		return 0
	}
	player = &mysqld.DBPlayer{
		Token: p.Token,
	}
	player.LoadWithToken()
	if player.ID <= 0 {
		return 1
	}
	mysqld.AddDBPlayer(player)
	p.Player = player
	player.LoginIP = p.Ip
	return 0
}
