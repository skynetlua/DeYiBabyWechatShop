package handle
import (
	"github.com/kataras/iris/v12"
	"bestsell/common"
)

//=>/live/rooms true get  
func Live_rooms(ctx iris.Context, sess *common.BSSession) {
	empty("/live/rooms")
}

//=>/live/his true get {roomId} 
func Live_his(ctx iris.Context, sess *common.BSSession) {
	empty("/live/his")
}