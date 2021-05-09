package common

import (
	//"bestsell/mysqld"
	"github.com/kataras/iris/v12"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type BSSession struct {
	Time       	time.Time
	Token   	string
	Data   		map[string]interface{}
	Params     	map[string]interface{}
	LockMutex 	sync.RWMutex
}

func (p *BSSession)Get(key string)interface{} {
	p.LockMutex.Lock()
	ret, _ := p.Params[key]
	p.LockMutex.Unlock()
	return ret
}

func (p *BSSession)GetString(key string)string {
	p.LockMutex.Lock()
	ret, ok := p.Params[key]
	p.LockMutex.Unlock()
	if !ok {
		return ""
	}
	switch ret.(type) {
	case string:
		return ret.(string)
	}
	return ""
}

func (p *BSSession)GetInt(key string)int {
	p.LockMutex.Lock()
	ret, ok := p.Params[key]
	p.LockMutex.Unlock()
	if !ok {
		return 0
	}
	switch ret.(type) {
	case int:
		return ret.(int)
	}
	return 0
}

func (p *BSSession)Set(key string, val interface{}) {
	p.LockMutex.Lock()
	p.Params[key] = val
	p.LockMutex.Unlock()
}

var _sessions map[string]*BSSession
var _sessionMutex sync.Mutex
var _tokenChars = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*")

func init() {
	_sessions = make(map[string]*BSSession)
	go clearTimeout()
}

func clearTimeout() {
	ticker := time.Tick(time.Minute)
	for {
		now := time.Now()
		_sessionMutex.Lock()
		for k, s := range _sessions {
			if now.After(s.Time) {
				delete(_sessions, k)
			}
		}
		_sessionMutex.Unlock()
		<-ticker
	}
}

func genToken() string {
	length := len(_tokenChars)
	token := [32]byte{}
	for i := 0; i < len(token); i++ {
		token[i] = _tokenChars[rand.Intn(length)]
	}
	return string(token[:])
}

//func SetSession(ctx iris.Context, userId int32, userType int8, name string) {
//	token := genToken()
//	_sessionMutex.Lock()
//	_sessions[token] = &BSSession{
//		Time:       time.Now().Add(4 * time.Hour),
//		Params:     make(map[string]interface{}),
//	}
//	_sessionMutex.Unlock()
//}


func GetSession(ctx iris.Context) *BSSession {
	token := ctx.FormValue("token")
	if len(token) < 5 {
		return nil
	}
	_sessionMutex.Lock()
	session, ok := _sessions[token]
	if ok && session != nil && strings.Compare(token, session.Token) == 0 {
		session.Time = time.Now().Add(4 * time.Hour)
		_sessionMutex.Unlock()
		return session
	}else{
		_sessions[token] = &BSSession{
			Token: token,
			Time: time.Now().Add(4 * time.Hour),
			Params: make(map[string]interface{}),
			Data: make(map[string]interface{}),
			LockMutex: sync.RWMutex{},
		}
		session = _sessions[token]
	}
	_sessionMutex.Unlock()
	return session
}


//func GetSession(ctx iris.Context) *BSSession {
//	sell := ctx.GetCookie("sell")
//	if len(sell) == 0 {
//		sell = genToken()
//		ctx.SetCookieKV("sell", sell)
//	}
//	_sessionMutex.RLock()
//	session, ok := _sessions[sell]
//	if ok {
//		session.Time = time.Now().Add(4 * time.Hour)
//		_sessionMutex.RUnlock()
//		return session
//	}else{
//		_sessions[sell] = &BSSession{
//			Time: time.Now().Add(4 * time.Hour),
//			Params: make(map[string]interface{}),
//			Data: make(map[string]interface{}),
//		}
//		session = _sessions[sell]
//	}
//	_sessionMutex.RUnlock()
//	return session
//}

func deleteSession(ctx iris.Context) {
	sell := ctx.GetCookie("sell")
	if len(sell) == 0 {
		_sessionMutex.Lock()
		delete(_sessions, sell)
		_sessionMutex.Unlock()
	}
}

