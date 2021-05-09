package handle

import (
	"bestsell/common"
	"bestsell/mysqld"
	"github.com/kataras/iris/v12"
)


//=>/notice/list true post {} 
func Notice_list(ctx iris.Context, sess *common.BSSession) {
	dbNoticeList := mysqld.GetDBNoticeList()
	var noticeList []*map[string]interface{}
	for _,dbNotice := range dbNoticeList {
		_notice := map[string]interface{}{
			"id":dbNotice.ID,
			"title":dbNotice.Title,
			"content":dbNotice.Content,
		}
		noticeList = append(noticeList, &_notice)
	}
	data := map[string]interface{}{
		"totalRow": 1,
		"totalPage": 1,
		"noticeList": noticeList,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/notice/lastone true get {type} 
func Notice_lastone(ctx iris.Context, sess *common.BSSession) {
	dbNoticeList := mysqld.GetDBNoticeList()
	if len(dbNoticeList) == 0 {
		ctx.JSON(iris.Map{"code": 0})
		return
	}
	dbNotice :=  dbNoticeList[0]
	data := map[string]interface{}{
		"id":dbNotice.ID,
		"title":dbNotice.Title,
		"content":dbNotice.Content,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/notice/detail true get {id} 
func Notice_detail(ctx iris.Context, sess *common.BSSession) {
	id := common.AtoI(ctx.FormValue("id"))
	dbNotice := mysqld.GetDBNotice(id)
	if dbNotice == nil {
		ctx.JSON(iris.Map{"code": 0})
		return
	}
	data := map[string]interface{}{
		"id":dbNotice.ID,
		"title":dbNotice.Title,
		"content":dbNotice.Content,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}