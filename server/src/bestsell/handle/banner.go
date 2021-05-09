package handle

import (
	"bestsell/common"
	"bestsell/config"
	"github.com/kataras/iris/v12"
	"strings"
)

//
func init() {
}

//=>/banner/list true get {} 
func Banner_list(ctx iris.Context, sess *common.BSSession) {
	_types := ctx.FormValue("types")
	types := strings.Split(_types, ",")
	data := make(map[string]interface{})
	for _,_type := range types {
		banners := config.GetCfgUiBannerSliceByMap(_type)
		var list config.CfgUiBannerSlice
		for _, banner := range *banners {
			if banner.Status == 1 {
				list = append(list, banner)
			}
		}
		data[_type] = list
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}