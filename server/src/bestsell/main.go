
package main

import (
	"bestsell/common"
	_ "bestsell/common"
	"bestsell/config"
	"bestsell/generate"
	"bestsell/handle"
	"bestsell/module"
	"bestsell/mysqld"
	"bestsell/router"
	"bestsell/test"
	"fmt"
	"log"
	"runtime"
)


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Llongfile)
	fmt.Println("main Config：", common.Config)
	mode := common.Config.Mode
	switch mode {
	case "api":
		startApiServer()
	case "res":
		startResServer()
	case "tool":
		startToolServer()
	case "excel":
		startExcelServer()
	default:
		panic("unknown mode:"+common.Config.Mode)
	}
}

func startExcelServer() {
	fmt.Println("startExcelServer")
	test.ConvertDatas()
	test.ConvertExcel2JsonArrays()

	ch := make(chan bool)
	mysqld.StartServer(&ch)
	fmt.Println("[main]startExcelServer mysqld")

	close(ch)

	mysqld.LoadConfigs()

	fmt.Println("[main]startExcelServer close")
}

func startToolServer()  {
	generate.Generate()
}

func startResServer()  {
	fmt.Println("startResServer")
	router.StartServer(nil)

	fmt.Println("[main]StartServer close")
}

func startApiServer() {
	fmt.Println("startApiServer")
	ch := make(chan bool)
	config.StartServer(&ch)
	fmt.Println("[main]ApiServer config")

	mysqld.StartServer(&ch)
	fmt.Println("[main]ApiServer mysqld")

	handle.StartServer(&ch)
	fmt.Println("[main]ApiServer handle")

	module.StartServer(&ch)
	fmt.Println("[main]ApiServer module")

	close(ch)

	fmt.Println("[main]ApiServer router")
	router.StartServer(nil)

	fmt.Println("[main]ApiServer close")
}


