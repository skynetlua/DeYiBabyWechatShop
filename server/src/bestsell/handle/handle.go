package handle

import (
	"bestsell/common"
	"bestsell/sdk"
	"fmt"
	"github.com/kataras/iris/v12"
	"strings"
	"time"
)


var wechatAccessToken string
var wechatAccessTokenExpire time.Time

func init() {
	//CheckWeChatAccessToken()
}
func empty(name string)  {
	fmt.Println("empty "+name)
}

func redirect(ctx iris.Context, url string) {
	ctx.ViewData("redirect", url)
	ctx.View("sell/redirect.html")
}

func Handle(ctx iris.Context, sess *common.BSSession) {
	_path := ctx.Path()
	sess.Set("lastPath",_path)
	redirect(ctx, "/sell/login")
}

func CheckWeChatAccessToken() int {
	if len(wechatAccessToken) > 0 {
		curTime := time.Now()
		outlineTime := wechatAccessTokenExpire
		curUnix := curTime.Unix()
		outlineUnix := outlineTime.Unix()
		if curUnix < outlineUnix {
			return 0
		}
	}
	reply := make(map[string]interface{})
	if sdk.OnWeChatAccessToken(&reply) != 0 {
		panic("OnWeChatAccessToken error")
		return -1
	}
	access_token := reply["access_token"].(string)
	expires_in := reply["expires_in"].(int)-100
	outlineTime := time.Now().Add(time.Duration(expires_in)*time.Second)
	fmt.Println("outlineTime =", outlineTime)

	wechatAccessToken = access_token
	wechatAccessTokenExpire = outlineTime
	return 0
}

func StartServer(ch *chan bool) {
	go func() {
		fmt.Println("StartServer handle")
		(*ch) <- true
	}()
	<-(*ch)
}

func addUrlHost(url string)string {
	if len(url) == 0 {
		return ""
	}
	if !strings.Contains(url, "http") {
		if strings.HasPrefix(url, "/") {
			url = common.StaticUrl+"/static"+url
		}else{
			url = common.StaticUrl+"/static/"+url
		}
	}
	return url
}

func addUrlHosts(urlPath string)[]string {
	var items []string
	if len(urlPath) == 0{
		return items
	}
	tmps := strings.Split(urlPath, ";")
	for _,tmp := range tmps {
		item := addUrlHost(tmp)
		items = append(items, item)
	}
	return items
}

func removeUrlHost(url string)string {
	if len(url) == 0 {
		return ""
	}
	hostUrl := common.StaticUrl+"/static/"
	if strings.Contains(url, hostUrl) {
		url = url[len(hostUrl):]
	}
	return url
}

func removeUrlHosts(urlPath string)string {
	if len(urlPath) == 0{
		return ""
	}
	tmps := strings.Split(urlPath, ";")
	var items []string
	for _,tmp := range tmps {
		item := removeUrlHost(tmp)
		items = append(items, item)
	}
	return strings.Join(items, ";")
}

func splitItems(txt string)[]string {
	var items []string
	if len(txt) == 0{
		return items
	}
	items = strings.Split(txt, ";")
	return items
}