package handle

import (
	"bestsell/common"
	"bestsell/config"
	"bestsell/module"
	"bestsell/mysqld"
	"fmt"
	"github.com/kataras/iris/v12"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
)


//=>/gm/goods/list true get {status, page, pageSize} 
func Gm_goods_list(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	status := common.AtoI(ctx.FormValue("status"))
	page := common.AtoI(ctx.FormValue("page"))
	pageSize := common.AtoI(ctx.FormValue("pageSize"))
	goodsList := mysqld.GetDBGoodsListByStatus(mysqld.EGoodsStatus(status))

	num := len(goodsList)
	startIdx := page*pageSize
	endIdx := (page+1)*pageSize
	if startIdx > num {
		startIdx = num
	}
	if endIdx > num {
		endIdx = num
	}
	goodsList = goodsList[startIdx:endIdx]

	var data []*map[string]interface{}
	for _, item := range goodsList {
		dbGoodsInfo := item.GetGoodsInfo()
		itemInfo := &map[string]interface{}{
			"goodsId" 		:item.GoodsId,
			"status" 		:int(item.GetStatus()),
			"numberStore" 	:item.NumberStore,
			"mark"  		:item.Mark,
			"sellPrice" 	:dbGoodsInfo.SellPrice,
			"minPrice" 		:dbGoodsInfo.MinPrice,
			"categoryId" 	:dbGoodsInfo.CategoryId,
			"startTime" 	:item.StartTime,
			"endTime" 		:item.EndTime,
			"name"  		:dbGoodsInfo.Name,
			"pic" 			:dbGoodsInfo.GetPublicIcon(),
			"order" 		:dbGoodsInfo.Order,
		}
		data = append(data, itemInfo)
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/gm/goods/info true get {goodsId} 
func Gm_goods_info(ctx iris.Context, sess *common.BSSession) {
	player := module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	goodsId := common.AtoI(ctx.FormValue("goodsId"))
	var dbGoods *mysqld.DBGoods
	if goodsId > 0 {
		dbGoods = mysqld.GetDBGoodsOrFromDB(goodsId)
	}
	if dbGoods == nil {
		dbGoods = &mysqld.DBGoods{
			GoodsId:0,
		}
	}
	data := map[string]interface{}{
		"goods":dbGoods.GetEditInfo(),
	}
	ctx.JSON(iris.Map{"code": 0, "data":data})
}

//=>/gm/goods/update true post {} 
func Gm_goods_update(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	goodsId := common.AtoI(ctx.FormValue("goodsId"))

	//skus := ctx.FormValue("skus")
	//skuPrices := ctx.FormValue("skuPrices")
	name := ctx.FormValue("name")
	barCode := ctx.FormValue("barCode")
	categoryId := common.AtoI(ctx.FormValue("categoryId"))
	if categoryId <= 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"没有指定主目录"})
		return
	}
	mainType := ctx.FormValue("mainType")
	subType := ctx.FormValue("subType")

	promote := ctx.FormValue("promote")
	order := common.AtoI(ctx.FormValue("order"))
	status := common.AtoI(ctx.FormValue("status"))
	pics := ctx.FormValue("pics")
	if len(pics) > 0 {
		pics = removeUrlHosts(pics)
	}
	contents := ctx.FormValue("contents")
	if len(contents) > 0 {
		contents = removeUrlHosts(contents)
	}
	skuJson := ctx.FormValue("skuJson")

	originPrice := common.AtoI(ctx.FormValue("originPrice"))
	sellPrice := common.AtoI(ctx.FormValue("sellPrice"))
	if sellPrice > originPrice {
		originPrice = sellPrice
	}
	//minPrice := common.MakeMoneyValue(common.AtoF(ctx.FormValue("minPrice")))
	mark := 0
	if len(ctx.FormValue("mark")) > 0 {
		mark = common.AtoI(ctx.FormValue("mark"))
	}
	numberStore := common.AtoI(ctx.FormValue("numberStore"))
	numberSell := common.AtoI(ctx.FormValue("numberSell"))
	cfgGoodsData := config.GetCfgGoodsDataByBarCode(barCode)
	//fmt.Println("Gm_goods_update  goodsId =", goodsId)
	dbGoods := mysqld.GetDBGoodsOrFromDB(goodsId)
	if dbGoods == nil && goodsId > 100000 {
		//barCode
		if cfgGoodsData == nil {
			ctx.JSON(iris.Map{"code": -1, "msg":"后台无商品数据"})
			return
		}
		if cfgGoodsData.GoodsId != goodsId {
			ctx.JSON(iris.Map{"code": -1, "msg":"商品ID不一致"})
			return
		}
		dbGoods := mysqld.AddNewDBGoodsByInfo(goodsId, name)
		dbGoodsInfo := dbGoods.GetGoodsInfo()

		dbGoodsInfo.PlayerId = player.ID
		dbGoodsInfo.CategoryId = categoryId
		if len(pics) > 0 {
			dbGoodsInfo.Pics = pics
		}
		if len(contents) > 0 {
			dbGoodsInfo.Contents = contents
		}
		dbGoodsInfo.Order = order
		dbGoodsInfo.CategoryId = categoryId
		if len(name) > 0{
			dbGoodsInfo.Name = name
		} else {
			dbGoodsInfo.Name = cfgGoodsData.Name
		}
		dbGoodsInfo.Promote = promote
		dbGoodsInfo.SellPrice = sellPrice
		dbGoodsInfo.OriginPrice = originPrice
		//dbGoodsInfo.MinPrice = minPrice
		dbGoodsInfo.BarCode = barCode
		if len(mainType) > 0 {
			dbGoodsInfo.MainType = mainType
		} else {
			dbGoodsInfo.MainType = cfgGoodsData.MainType
		}
		if len(subType) > 0 {
			dbGoodsInfo.SubType = subType
		} else {
			dbGoodsInfo.SubType = cfgGoodsData.SubType
		}
		dbGoodsInfo.NumberSell = numberSell
		dbGoodsInfo.TryLoadImage()
		//dbGoodsInfo.Skus = skus
		//dbGoodsInfo.SkuPrices = skuPrices
		dbGoodsInfo.Save()

		dbGoods.NumberStore = numberStore
		dbGoods.SetStatus(mysqld.EGoodsStatus(status))
		dbGoods.Save()

		mysqld.RemoveDBGoods(goodsId)
		dbGoods = mysqld.GetDBGoodsOrFromDB(goodsId)
		dbGoods.GetGoodsInfo()
		mysqld.LoadGoodsCaches()

		ctx.JSON(iris.Map{"code": 0, "msg":"保存成功"})
		return
	}
	if dbGoods == nil || dbGoods.GoodsId != goodsId{
		ctx.JSON(iris.Map{"code": -1, "msg":"商品id出错"})
		return
	}
	dbGoodsInfo := dbGoods.GetGoodsInfo()
	if dbGoodsInfo == nil || goodsId != dbGoodsInfo.GoodsId {
		ctx.JSON(iris.Map{"code": -1, "msg":"商品详情id出错"})
		return
	}
	for  {
		if dbGoodsInfo.Order != order { break }
		if dbGoodsInfo.CategoryId != categoryId { break }
		if dbGoodsInfo.Promote != promote { break }
		if dbGoodsInfo.SellPrice != sellPrice { break }
		if dbGoodsInfo.OriginPrice != originPrice { break }
		//if dbGoodsInfo.MinPrice != minPrice { break }
		//if dbGoodsInfo.Skus != skus { break }
		//if dbGoodsInfo.SkuPrices != skuPrices { break }
		if dbGoodsInfo.Pics != pics { break }
		if dbGoodsInfo.Contents != contents { break }
		if dbGoodsInfo.Name != name { break }
		if dbGoodsInfo.MainType != mainType { break }
		if dbGoodsInfo.SubType != subType { break }
		if dbGoods.NumberStore != numberStore { break }
		if dbGoods.NumberSell != numberSell { break }
		if dbGoodsInfo.SkuStruct != skuJson { break }
		if dbGoods.GetStatus() != mysqld.EGoodsStatus(status) {
			dbGoods.SetStatus(mysqld.EGoodsStatus(status))
			dbGoods.Save()
		}
		mysqld.ReloadGoods(goodsId)
		ctx.JSON(iris.Map{"code": 0, "msg":"保存成功"})
		return
	}
	if !dbGoodsInfo.SetSkuStruct(skuJson){
		ctx.JSON(iris.Map{"code": -1, "msg":"商品规格有错误"})
		return
	}

	dbGoodsInfo.BeginWrite()

	dbGoodsInfo.PlayerId = player.ID
	if len(name) > 0{
		dbGoodsInfo.Name = name
	} else {
		dbGoodsInfo.Name = cfgGoodsData.Name
	}
	if len(mainType) > 0 {
		dbGoodsInfo.MainType = mainType
	} else {
		dbGoodsInfo.MainType = cfgGoodsData.MainType
	}
	if len(subType) > 0 {
		dbGoodsInfo.SubType = subType
	} else {
		dbGoodsInfo.SubType = cfgGoodsData.SubType
	}
	dbGoodsInfo.Promote = promote
	//dbGoodsInfo.Name = name
	//dbGoodsInfo.BarCode = barCode
	//dbGoodsInfo.Desc = ctx.FormValue("desc")
	dbGoodsInfo.CategoryId = categoryId
	if len(pics) > 0 {
		dbGoodsInfo.Pics = pics
	}
	if len(contents) > 0 {
		dbGoodsInfo.Contents = contents
	}
	dbGoodsInfo.Order = order
	dbGoodsInfo.EndWrite()

	for  {
		if dbGoodsInfo.SellPrice != sellPrice { break }
		if dbGoodsInfo.OriginPrice != originPrice { break }
		//if dbGoodsInfo.MinPrice != minPrice { break }
		//if dbGoodsInfo.Skus != skus { break }
		//if dbGoodsInfo.SkuPrices != skuPrices { break }
		if dbGoods.NumberStore != numberStore { break }
		if dbGoods.NumberSell != numberSell { break }
		dbGoodsInfo.Save()
		if dbGoods.GetStatus() != mysqld.EGoodsStatus(status) {
			dbGoods.SetStatus(mysqld.EGoodsStatus(status))
			dbGoods.Save()
		}
		mysqld.ReloadGoods(goodsId)
		ctx.JSON(iris.Map{"code": 0, "msg":"保存成功"})
		return
	}
	//if status == int(mysqld.GoodsStatusUp) && dbGoods.GetStatus() == mysqld.GoodsStatusUp {
	//	dbGoodsInfo.Save()
	//	ctx.JSON(iris.Map{"code": -1, "msg":"需要把商品下架才能编辑价格规格"})
	//	return
	//}

	//dbGoodsInfo.SetSkus(skus)
	//dbGoodsInfo.SetSkuPrices(skuPrices)

	dbGoodsInfo.BeginWrite()
	dbGoodsInfo.PlayerId = player.ID
	//dbGoodsInfo.EnterPrice = enterPrice
	//dbGoodsInfo.MinPrice = minPrice
	dbGoodsInfo.OriginPrice = originPrice
	dbGoodsInfo.SellPrice = sellPrice
	dbGoodsInfo.NumberSell = numberSell

	dbGoodsInfo.EndWrite()
	dbGoodsInfo.Save()

	dbGoods.SetStatus(mysqld.EGoodsStatus(status))

	dbGoods.BeginWrite()
	dbGoods.Mark = mark
	dbGoods.NumberStore = numberStore
	dbGoods.EndWrite()
	dbGoods.Save()

	mysqld.ReloadGoods(goodsId)

	ctx.JSON(iris.Map{"code": 0, "msg":"保存成功"})
}

//=>/gm/upload/goods true post {} 
func Gm_upload_goods(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	//maxSize := ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()
	maxSize := int64(1 << 20)
	err := ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		//ctx.WriteString(err.Error())
		fmt.Println("Gm_upload_goods err =", err.Error())
		ctx.JSON(iris.Map{"code": -1, "msg":"文件太大"})
		return
	}
	form := ctx.Request().MultipartForm
	if form.File["upfile"] == nil || len(form.File["upfile"]) == 0 {
		fmt.Println("Gm_upload_goods  no file")
		ctx.JSON(iris.Map{"code": -1, "msg":"上传文件丢失"})
		return
	}
	file := form.File["upfile"][0]
	contentType := file.Header.Get("content-type")
	fmt.Println("contentType =", contentType)
	fileExtension := ".png"
	if strings.Contains(contentType, "jpeg") {
		fileExtension = ".jpeg"
	}
	categoryId := ctx.FormValue("categoryId")
	goodsId := ctx.FormValue("goodsId")
	part := ctx.FormValue("part")
	idx := ctx.FormValue("idx")
	if len(categoryId) == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"参数缺失1"})
		return
	}
	if len(goodsId) == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"参数缺失2"})
		return
	}
	if len(part) == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"参数缺失3"})
		return
	}
	if len(idx) == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"参数缺失4"})
		return
	}
	dirFolder := path.Join(common.UploadGoodsPath, categoryId)
	common.CreateDir(dirFolder)

	saveFolder := path.Join(dirFolder, categoryId+"_"+goodsId)
	common.CreateDir(saveFolder)
	saveFolder = path.Join(saveFolder, part)
	common.CreateDir(saveFolder)
	saveFile := path.Join(saveFolder, idx+fileExtension)
	_, err = saveGoodsUploadedFile(file, saveFile)
	if err != nil {
		fmt.Println("failed to upload:", file.Filename)
		ctx.JSON(iris.Map{"code": -1, "msg":"保存文件失败"})
		return
	}
	urlPath := common.FillPathHeader(saveFile)
	data := map[string]interface{}{
		"categoryId":categoryId,
		"goodsId" 	:common.AtoI(goodsId),
		"url" 		:urlPath,
		"part" 		:part,
		"idx" 		:common.AtoI(idx),
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
	module.RemoveFilePublicPath(urlPath)
}

func saveGoodsUploadedFile(fh *multipart.FileHeader, saveFile string) (int64, error) {
	src, err := fh.Open()
	if err != nil {
		return 0, err
	}
	defer src.Close()
	out, err := os.OpenFile(saveFile, os.O_WRONLY|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		return 0, err
	}
	defer out.Close()
	return io.Copy(out, src)
}

//=>/gm/goods/barcode true get {barCode} 
func Gm_goods_barcode(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	barCode := ctx.FormValue("barCode")
	cfgGoodsData := config.GetCfgGoodsDataByBarCode(barCode)
	if cfgGoodsData == nil {
		ctx.JSON(iris.Map{"code": 1})
		return
	}
	data := map[string]interface{}{
		"goodsId": cfgGoodsData.GoodsId,
		"barCode": barCode,
	}
	ctx.JSON(iris.Map{"code": 0, "data":data})
}

//=>/gm/goods/update/info true post {} 
func Gm_goods_update_info(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	goodsId := common.AtoI(ctx.FormValue("goodsId"))
	dbGoods := mysqld.GetDBGoodsOrFromDB(goodsId)
	if dbGoods == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"参数有错"})
		return
	}
	status := mysqld.EGoodsStatus(common.AtoI(ctx.FormValue("status")))
	//dbGoods.SetStatus(status)
	//orginStatus := dbGoods.GetStatus()
	//if orginStatus != status {
	//	if orginStatus == mysqld.GoodsStatusEditing {
	//		ctx.JSON(iris.Map{"code": -1, "msg":"录入商品不能修改状态"})
	//		return
	//	}else if orginStatus == mysqld.GoodsStatusEdited {
	//		if status != mysqld.GoodsStatusUp {
	//			ctx.JSON(iris.Map{"code": -1, "msg":"待上架商品只能改成上架"})
	//			return
	//		}
	//	}else if orginStatus == mysqld.GoodsStatusUp {
	//		if status != mysqld.GoodsStatusDown {
	//			ctx.JSON(iris.Map{"code": -1, "msg":"上架商品只能改成下架"})
	//			return
	//		}
	//	}else if orginStatus == mysqld.GoodsStatusDown {
	//		if status != mysqld.GoodsStatusUp {
	//			ctx.JSON(iris.Map{"code": -1, "msg":"下架商品只能改成上架"})
	//			return
	//		}
	//	}else{
	//		ctx.JSON(iris.Map{"code": -1, "msg":"商品不存在"})
	//		return
	//	}
	//}
	numberStore := common.AtoI(ctx.FormValue("numberStore"))
	numberSell := common.AtoI(ctx.FormValue("numberSell"))
	sellPrice := common.AtoI(ctx.FormValue("sellPrice"))
	originPrice := common.AtoI(ctx.FormValue("originPrice"))
	mark := common.AtoI(ctx.FormValue("mark"))
	categoryId := common.AtoI(ctx.FormValue("categoryId"))
	order := common.AtoI(ctx.FormValue("order"))

	startTime := 0
	if len(ctx.FormValue("startTime")) > 0 {
		startTime = common.AtoI(ctx.FormValue("startTime"))
	}
	endTime := 0
	if len(ctx.FormValue("endTime")) > 0 {
		endTime = common.AtoI(ctx.FormValue("endTime"))
	}

	dbGoodsInfo := dbGoods.GetGoodsInfo()

	//if mark == int(mysqld.GoodsMarkSeckill) {
	//	skus := dbGoodsInfo.GetSkus()
	//	if len(skus) != 0 {
	//		ctx.JSON(iris.Map{"code": -1, "msg":"有规格的产品，不能参与秒杀"})
	//		return
	//	}
	//}
	//if mark == int(mysqld.GoodsMarkTeam) {
	//	skus := dbGoodsInfo.GetSkus()
	//	if len(skus) != 0 {
	//		ctx.JSON(iris.Map{"code": -1, "msg":"有规格的产品，不能参与拼团"})
	//		return
	//	}
	//}

	for {
		if dbGoods.Mark != mark { break }
		if dbGoods.Status != int(status) { break }
		if dbGoods.StartTime != int64(startTime) { break }
		if dbGoods.EndTime != int64(endTime) { break }
		if dbGoods.NumberStore != numberStore { break }
		if dbGoodsInfo.NumberSell != numberSell { break }
		if dbGoodsInfo.CategoryId != categoryId { break }
		if dbGoodsInfo.Order != order { break }
		if dbGoodsInfo.SellPrice != sellPrice { break }
		if dbGoodsInfo.OriginPrice != originPrice { break }
		fmt.Println("Gm_goods_update_info nochange goodsId =", goodsId)
		ctx.JSON(iris.Map{"code": -1, "msg":"无任何修改"})
		return
	}
	dbGoods.SetStatus(status)

	//if orginStatus == mysqld.GoodsStatusUp {
	//	if dbGoodsInfo.SellPrice != sellPrice {
	//		ctx.JSON(iris.Map{"code": -1, "msg":"上架商品，不能修改价格"})
	//		return
	//	}
	//}
	dbGoods.BeginWrite()
	dbGoods.Mark = mark
	dbGoods.StartTime = int64(startTime)
	dbGoods.EndTime = int64(endTime)
	dbGoods.NumberStore = numberStore
	dbGoods.EndWrite()
	dbGoods.Save()

	dbGoodsInfo.BeginWrite()
	dbGoodsInfo.PlayerId = player.ID
	if categoryId > 0 {
		dbGoodsInfo.CategoryId = categoryId
	}
	if order > 0 {
		dbGoodsInfo.Order = order
	}
	dbGoodsInfo.OriginPrice = originPrice
	dbGoodsInfo.EndWrite()
	dbGoodsInfo.SetNumberSell(numberSell)
	dbGoodsInfo.SetSellPrice(sellPrice)
	dbGoodsInfo.Save()

	mysqld.ReloadGoods(goodsId)
	ctx.JSON(iris.Map{"code": 0})
}

//=>/gm/goods/remove true get {goodsId} 
func Gm_goods_remove(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	goodsId := common.AtoI(ctx.FormValue("goodsId"))
	dbGoods := mysqld.GetDBGoodsOrFromDB(goodsId)
	if dbGoods == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"商品已删除"})
		return
	}
	if dbGoods.GetStatus() == mysqld.EGoodsStatusUp {
		ctx.JSON(iris.Map{"code": -1, "msg":"商品需要下架才能删除"})
		return
	}
	mysqld.RemoveDBGoodsAndDB(goodsId)
	ctx.JSON(iris.Map{"code": 0})
}

//=>/gm/goods/category true get {categoryId} 
func Gm_goods_category(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	categoryId := common.AtoI(ctx.FormValue("categoryId"))
	goodsList := mysqld.GetDBGoodsListByCategoryId(categoryId, mysqld.EGoodsStatusAll)
	var data []*map[string]interface{}
	for _, item := range goodsList {
		dbGoodsInfo := item.GetGoodsInfo()
		itemInfo := &map[string]interface{}{
			"goodsId" 		:item.GoodsId,
			"status" 		:int(item.GetStatus()),
			"numberStore" 	:item.NumberStore,
			"mark" 			:item.Mark,
			"startTime" 	:item.StartTime,
			"endTime" 		:item.EndTime,
			"name"  		:dbGoodsInfo.Name,
			"pic" 			:dbGoodsInfo.GetPublicIcon(),
			"sellPrice" 	:dbGoodsInfo.SellPrice,
			"originPrice" 	:dbGoodsInfo.OriginPrice,
			"categoryId" 	:dbGoodsInfo.CategoryId,
			"order" 		:dbGoodsInfo.Order,
			"numberSell" 	:dbGoodsInfo.NumberSell,
			"updateTime"	: item.UpdatedAt.Unix(),
		}
		data = append(data, itemInfo)
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/gm/goods/goodsdata true get {barCode}
func Gm_goods_goodsdata(ctx iris.Context, sess *common.BSSession) {
	player := module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg": "未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg": "无权限"})
		return
	}
	barCode := ctx.FormValue("barCode")
	cfgGoodsData := config.GetCfgGoodsDataByBarCode(barCode)
	data := map[string]interface{}{}
	if cfgGoodsData != nil {
		data["goodsId"] 	= cfgGoodsData.GoodsId
		data["barCode"] 	= cfgGoodsData.BarCode
		data["name"] 		= cfgGoodsData.Name
		data["mainType"] 	= cfgGoodsData.MainType
		data["subType"] 	= cfgGoodsData.SubType
		data["unit"] 		= cfgGoodsData.Unit
		//data["enterPrice"] 	= cfgGoodsData.EnterPrice
		data["sellPrice"] 	= cfgGoodsData.GetSellPrice()
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/gm/goods/goodsdatas true get {}
func Gm_goods_goodsdatas(ctx iris.Context, sess *common.BSSession) {
	player := module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg": "未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg": "无权限"})
		return
	}
	cfgGoodsDatas := config.GetCfgGoodsDatas()
	var datas []*map[string]interface{}
	for _, item := range cfgGoodsDatas {
		data := &map[string]interface{} {
			"goodsId":	item.GoodsId,
			"barCode":	item.BarCode,
			"name": 	item.Name,
			"mainType": item.MainType,
			"subType":	item.SubType,
			"unit":		item.Unit,
			"sellPrice":item.GetSellPrice(),
		}
		dbGoods := mysqld.GetDBGoods(item.GoodsId)
		if dbGoods != nil {
			dbGoodsInfo := dbGoods.GetGoodsInfo()
			(*data)["icon"] =  dbGoodsInfo.GetPublicIcon()
		}
		datas = append(datas, data)
	}
	ctx.JSON(iris.Map{"code": 0, "data": datas})
}

