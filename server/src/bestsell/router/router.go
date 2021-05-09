package router

import (
	"bestsell/common"
	"bestsell/handle"
	"bestsell/module"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"net/http"
	"path"
	"strings"
	"time"
)

type Handler func(iris.Context, *common.BSSession)


var appIris = newApp()
//var appIrisSess *sessions.Sessions
var controlCloseServer = false

func init() {
}

func newApp() *iris.Application {
	//dbPath := path.Join(common.ProjectPath,"sessions.db")
	//db, err := boltdb.New(dbPath, os.ModePerm)
	//if err != nil {
	//	panic(err)
	//}
	//db := redis.New(redis.Config{
	//	Network:   "tcp",
	//	Addr:      "127.0.0.1:6379",
	//	Timeout:   time.Duration(30) * time.Second,
	//	MaxActive: 10,
	//	Password:  "",
	//	Database:  "",
	//	Prefix:    "",
	//	Delim:     "-",
	//	Driver:    redis.Redigo(), // redis.Radix() can be used instead.
	//})
	//
	//iris.RegisterOnInterrupt(func() {
	//	db.Close()
	//})
	//defer db.Close()
	//
	//sess := sessions.New(sessions.Config{
	//	Cookie: "sell",
	//	Expires: 24 * time.Hour,
	//	AllowReclaim: true,
	//})
	//sess.UseDatabase(db)
	//appIrisSess = sess

	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	const refreshEvery = 1 * time.Second
	app.Use(iris.Cache304(refreshEvery))

	//app.Get("/", before, mainHandler, after)
	//app.Use(before)
	//app.Done(after)

	//app.Handle("GET", "/", func(ctx iris.Context) {
	//	ctx.HTML("<h1>Welcome to the Good Island</h1>")
	//})

	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString("pong")
	})

	if common.Config.Enable["debug"] == 1 {
		app.Logger().SetLevel("debug")
	}

	if common.Config.Enable["web"] != 1 {
		return app
	}

	filePath := path.Join(common.StaticPath, "/favicon.ico")
	app.Favicon(filePath)

	//filePath = "/Users/linyou/svn/plan1/vant/client/dist"
	//app.HandleDir("/", filePath, iris.DirOptions{Gzip: false,ShowList: false})
	//templateView := iris.HTML(filePath, ".html")
	//app.RegisterView(templateView)

	filePath =  path.Join(common.AssetPath, "/view")
	templateView := iris.HTML(filePath, ".html")
	templateView.Layout("layout.html")
	templateView.Reload(true)
	app.RegisterView(templateView)

	filePath = common.StaticPath
	app.HandleDir("/static", filePath, iris.DirOptions{Gzip: false,ShowList: false})
	return app
}

func StartServer(ch *chan bool) {
	//go func() {
		fmt.Println("StartServer router")
		startMap()
		serverAddr := common.Config.ServerAddr
		domain := common.Config.Domain
		//items := strings.Split(serverAddr, ":")
		////ip := items[0]
		//port := items[1]
		//if len(port) == 0 {
		//	if strings.Contains(domain, "https") {
		//		port = "444"
		//	}else{
		//		port = "8080"
		//	}
		//}
		//(*ch) <- true
		fmt.Println("server start "+domain)
		if common.Config.Enable["https"] == 1 {
			//mycert := path.Join(common.AssetPath, "/tls/Nginx/1_www.bestsellmall.com_bundle.crt")
			//mykey := path.Join(common.AssetPath, "/tls/Nginx/2_www.bestsellmall.com.key")
			mycert := path.Join(common.AssetPath, common.Config.HttpsCrt)
			mykey := path.Join(common.AssetPath, common.Config.HttpsKey)

			//target, _ := url.Parse(serverAddr)
			//go host.NewRedirection(ip+":80", target, iris.StatusMovedPermanently).ListenAndServe()
			appIris.Run(iris.TLS(serverAddr, mycert, mykey), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
		}else{
			appIris.Run(iris.Addr(serverAddr), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
		}
		fmt.Println("appIris is exit!!!!")
	//}()
	//<-(*ch)
}

func HandleHandle(arg *HandleArg)  {
	if strings.Compare(arg.method,"post") == 0 {
		POST(arg.fullPath, arg.handle, arg)
	}else{
		GET(arg.fullPath, arg.handle, arg)
	}
}

func permit(ctx iris.Context, path string, arg *HandleArg) int {
	//fmt.Println( "bind path:", path)
	//fmt.Println("permit:",ctx.GetStatusCode(), ctx.Method(), ctx.Path())
	sess := common.GetSession(ctx)
	arg.session = sess
	if common.Config.Enable["token"] != 1 {
		fmt.Println( "not token path =", path ," FormValues =", ctx.FormValues())
		return http.StatusOK
	}
	for sess != nil {
		if sess.Params == nil {
			return http.StatusForbidden
		}
		visits := sess.GetInt("visits")
		sess.Set("visits", visits+1)

		userLogin := &module.UserLogin{
			Token: sess.Token,
		}
		module.OnLogin(userLogin)
		player := userLogin.Player
		if player == nil {
			break
		}
		sess.Data["player"] = player
		fmt.Println( "visits =", visits, "path =", path, "FormValues =", ctx.FormValues())
		return http.StatusOK
	}
	fmt.Println("StatusUnauthorized path =", path, "FormValues =", ctx.FormValues())
	if !arg.hasToken {
		return http.StatusOK
	}
	return http.StatusUnauthorized
}

//func protectRun(entry func()) {
//	// 延迟处理的函数
//	defer func() {
//		// 发生宕机时，获取panic传递的上下文并打印
//		err := recover()
//		switch err.(type) {
//		case runtime.Error: // 运行时错误
//			fmt.Println("runtime error:", err)
//		default: // 非运行时错误
//			fmt.Println("error:", err)
//		}
//	}()
//	entry()
//}

func GET(path string, handler Handler, arg *HandleArg) {
	appIris.Get(path, func(ctx iris.Context) {
		//fmt.Println( "appIris.Get======>> bind path:", path, "ctx.Path():", ctx.Path())
		switch permit(ctx, path, arg) {
		case http.StatusUnauthorized:
			ctx.StatusCode(iris.StatusUnauthorized)
			handle.Handle(ctx, arg.session)
		case http.StatusForbidden:
			ctx.StatusCode(iris.StatusForbidden)
			ctx.WriteString("权限不足:"+path)
			//handle.Handle(ctx, arg.session)
		case http.StatusOK:
			handler(ctx, arg.session)
		}
		//fmt.Println( "appIris.Get======<< bind path:", path, "ctx.Path():", ctx.Path())
	})
}

func POST(path string, handler Handler, arg *HandleArg) {
	appIris.Post(path, func(ctx iris.Context) {
		//fmt.Println( "appIris.Post======>> bind path:", path, "ctx.Path():", ctx.Path())
		//if controlCloseServer {
		//	ctx.StatusCode(iris.StatusNotFound)
		//	ctx.WriteString("服务器正在紧急维护，请等待几分钟")
		//	return
		//}
		switch permit(ctx, path, arg) {
		case http.StatusUnauthorized:
			ctx.StatusCode(iris.StatusUnauthorized)
			handle.Handle(ctx, arg.session)
			//fmt.Println("POST StatusUnauthorized path =", path)
			//ctx.WriteString("StatusUnauthorized:"+path)
		case http.StatusForbidden:
			ctx.StatusCode(iris.StatusForbidden)
			fmt.Println("POST StatusForbidden path =", path)
			ctx.WriteString("StatusForbidden:"+path)
		case http.StatusOK:
			handler(ctx, arg.session)
		}
		//fmt.Println( "appIris.Post======<< bind path:", path, "ctx.Path():", ctx.Path())
	})
}
