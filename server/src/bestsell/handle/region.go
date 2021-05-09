package handle
import (
	"bestsell/config"
	"bestsell/module"
	"github.com/kataras/iris/v12"
	"bestsell/common"
	"strconv"
	"fmt"
)

//=>/region/province false get  
func Region_province(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	data := config.GetCfgProvinceAddressSlice()
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/region/child false get {pid} 
func Region_child(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	_pid := ctx.FormValue("pid")
	pid,err1 := strconv.Atoi(_pid)
	if err1 != nil {
		fmt.Println("Common_region_child ",err1)
		ctx.JSON(iris.Map{"code": -1, "msg":"pid出错"})
		return
	}
	cfgAddressSlice := config.GetCfgAddressSub(pid)
	data := cfgAddressSlice
	ctx.JSON(iris.Map{"code": 0, "data": data})
}