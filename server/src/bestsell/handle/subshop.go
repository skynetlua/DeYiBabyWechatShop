package handle
import (
	"bestsell/common"
	"fmt"
	"github.com/kataras/iris/v12"
)

//=>/subshop/list true post {} 
func Subshop_list(ctx iris.Context, sess *common.BSSession) {
	empty("/subshop/list")
}

//=>/subshop/my true get {token} 
func Subshop_my(ctx iris.Context, sess *common.BSSession) {
	empty("/subshop/my")
}

//=>/subshop/detail true get {id} 
func Subshop_detail(ctx iris.Context, sess *common.BSSession) {
	_id := ctx.FormValue("id")
	fmt.Println("_id =", _id)
	info := map[string]interface{}{
		"pic":common.StaticUrl+"/static/api/qrcode/20200430/53b76d26893ec68d8d9b44586d3edbce.jpg",
		"name":"Q-Baby母婴生活馆（敏捷店）",
		"address":"",
	}
	data := map[string]interface{}{
		"info":info,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/subshop/apply true post {} 
func Subshop_apply(ctx iris.Context, sess *common.BSSession) {
	empty("/subshop/apply")
}