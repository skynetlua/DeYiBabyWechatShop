package handle

import (
	"bestsell/common"
	"bestsell/module"
	"bestsell/mysqld"
	"fmt"
	"github.com/kataras/iris/v12"
	"path"
	"strconv"
	"strings"
)


//=>/gm/category/list true get {} 
func Gm_category_list(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	dbCategorys := mysqld.GetDBCategoryList()
	var categorys []*map[string]interface{}
	for _, item := range dbCategorys {
		itemInfo := &map[string]interface{}{
			"id" 	 :item.ID,
			"name"   :item.Name,
			"status" :item.Status,
			"icon" 	 :item.GetIcon(),
			"level"  :item.Level,
			"order"  :item.Order,
			"publicIcon":item.GetPublicIcon(),
		}
		categorys = append(categorys, itemInfo)
	}
	resIds := mysqld.GetCategoryResIds()
	data := &map[string]interface{}{
		"resIds"   :resIds,
		"categorys":categorys,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/gm/category/update true post {} 
func Gm_category_update(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	id := common.AtoI(ctx.FormValue("id"))
	dbCategory := mysqld.GetDBCategoryOrFromDB(id)
	status := common.AtoI(ctx.FormValue("status"))
	level := common.AtoI(ctx.FormValue("level"))
	order := common.AtoI(ctx.FormValue("order"))
	name := ctx.FormValue("name")
	icon := ctx.FormValue("icon")
	if dbCategory == nil {
		dbCategorys := mysqld.GetDBCategoryList()
		maxId := 0
		for _, item := range dbCategorys {
			if item.ID > maxId {
				maxId = item.ID
			}
		}
		maxId++
		dbCategory = &mysqld.DBCategory{
			Name:    name+strconv.Itoa(maxId),
			Level:  level,
			Order:  order,
			Status: status,
		}
		dbCategory.SetIcon(icon)
		mysqld.AddNewCategory(dbCategory)
		dbCategory = mysqld.GetDBCategoryOrFromDB(dbCategory.ID)
	} else {
		dbCategory.BeginWrite()
		dbCategory.Name = name
		dbCategory.Level = level
		dbCategory.Order = order
		dbCategory.Status = status
		dbCategory.SetIcon(icon)
		dbCategory.EndWrite()
		dbCategory.Save()
	}
	data := &map[string]interface{}{
		"id" 	 :dbCategory.ID,
		"name"   :dbCategory.Name,
		"status" :dbCategory.Status,
		"icon" 	 :dbCategory.GetIcon(),
		"level"  :dbCategory.Level,
		"order"  :dbCategory.Order,
		"publicIcon":dbCategory.GetPublicIcon(),
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/gm/category/remove true get {} 
func Gm_category_remove(ctx iris.Context, sess *common.BSSession) {
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
	dbCategory := mysqld.GetDBCategoryOrFromDB(categoryId)
	if dbCategory == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"目录已删除"})
		return
	}
	mysqld.RemoveDBCategoryAndDB(categoryId)
	ctx.JSON(iris.Map{"code": 0})
}

//=>/gm/upload/category true post {} 
func Gm_upload_category(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
	if player.GM == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
		return
	}
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
	// fmt.Println("contentType =", contentType)
	fileExtension := ".png"
	if strings.Contains(contentType, "jpeg") {
		fileExtension = ".jpeg"
	}
	categoryId := ctx.FormValue("categoryId")
	if len(categoryId) == 0 {
		ctx.JSON(iris.Map{"code": -1, "msg":"参数缺失1"})
		return
	}
	saveFolder := common.UploadCategoryPath
	common.CreateDir(saveFolder)
	saveFile := path.Join(saveFolder, categoryId+fileExtension)
	_, err = saveGoodsUploadedFile(file, saveFile)
	if err != nil {
		fmt.Println("failed to upload:", file.Filename)
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	urlPath := common.FillPathHeader(saveFile)
	data := map[string]interface{}{
		"url" 		:urlPath,
		"categoryId":common.AtoI(categoryId),
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
	module.RemoveFilePublicPath(urlPath)
}

//=>/gm/goods/resids true get {} 
//func Gm_goods_resids(ctx iris.Context, sess *common.BSSession) {
//	player :=  module.GetPlayer(sess)
//	if player == nil {
//		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
//		return
//	}
//	if player.GM == 0 {
//		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
//		return
//	}
//	resIds := module.GetGoodsResIds()
//	data := map[string]interface{}{
//		"resIds":resIds,
//	}
//	ctx.JSON(iris.Map{"code": 0, "data": data})
//}

//=>/gm/upload/excel true post {} 
//func Gm_upload_excel(ctx iris.Context, sess *common.BSSession) {
//	player :=  module.GetPlayer(sess)
//	if player == nil {
//		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
//		return
//	}
//	if player.GM == 0 {
//		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
//		return
//	}
//	maxSize := int64(1 << 20)
//	err := ctx.Request().ParseMultipartForm(maxSize)
//	if err != nil {
//		ctx.StatusCode(iris.StatusInternalServerError)
//		//ctx.WriteString(err.Error())
//		fmt.Println("Gm_upload_excel err =", err.Error())
//		ctx.JSON(iris.Map{"code": -1, "msg":"文件太大"})
//		return
//	}
//	form := ctx.Request().MultipartForm
//	upfiles := form.File["upfile"]
//	if upfiles == nil || len(upfiles) == 0 {
//		fmt.Println("Gm_upload_excel  no file")
//		ctx.JSON(iris.Map{"code": -1, "msg":"上传文件丢失"})
//		return
//	}
//	file := upfiles[0]
//	fileName := ctx.FormValue("fileName")
//	if len(fileName) == 0 {
//		tm := time.Now()
//		fileName = tm.Format("20060102150405")+".xlsx"
//	}
//	saveFolder := path.Join(common.AssetPath, "gm")
//	common.CreateDir(saveFolder)
//	saveFile := path.Join(saveFolder, fileName)
//	_, err = saveGoodsUploadedFile(file, saveFile)
//	if err != nil {
//		fmt.Println("failed to upload:", file.Filename)
//		ctx.JSON(iris.Map{"code": -1, "msg":"上传文件出错"})
//		return
//	}
//	ret, onlineGoods := importGoodsDatas(saveFile, player)
//	if !ret {
//		ctx.JSON(iris.Map{"code": 0, "msg": "excel表格数据格式有错"})
//		return
//	}
//	ctx.JSON(iris.Map{"code": 0, "msg": "上传成功", "data":onlineGoods})
//}

//func importGoodsDatas(goodsPath string, player *mysqld.DBPlayer)(bool, []string){
//	sheets := common.LoadExcel2Map(goodsPath)
//	if len(*sheets) == 0 {
//		return false, nil
//	}
//	var onlineGoods []string
//	for _,sheet := range *sheets {
//		for _, _item := range sheet {
//			item := *_item
//			goodsInfo := &mysqld.DBGoodsInfo{
//				PlayerId  	:player.ID,
//				Name 		:common.ConvertString(item["name"]),
//				BarCode 	:common.ConvertString(item["barCode"]),
//				Desc 		:common.ConvertString(item["desc"]),
//				Skus  		:common.ConvertString(item["skus"]),
//
//				CategoryId  :common.ConvertInt(item["catogoryId"]),
//				EnterPrice 	:common.ConvertFloat64(item["enterPrice"]),
//				MinPrice 	:common.ConvertFloat64(item["minPrice"]),
//				OriginPrice :common.ConvertFloat64(item["originPrice"]),
//				SellPrice 	:common.ConvertFloat64(item["sellPrice"]),
//
//				EnterDate  	:common.ConvertString(item["enterDate"]),
//				ProductDate :common.ConvertString(item["productDate"]),
//
//				ResId 		:common.ConvertString(item["resId"]),
//			}
//			numberStore := common.ConvertInt(item["numberStore"])
//			if goodsInfo.MinPrice > goodsInfo.SellPrice {
//				goodsInfo.SellPrice = goodsInfo.MinPrice
//			}
//			if goodsInfo.SellPrice > goodsInfo.OriginPrice {
//				goodsInfo.OriginPrice = goodsInfo.SellPrice
//			}
//			if len(goodsInfo.Name) == 0 {
//				continue
//			}
//			dbGoods := mysqld.GetDBGoodsByGoodsName(goodsInfo.Name)
//			if dbGoods != nil {
//				if dbGoods.Status == int(mysqld.GoodsStatusUp) {
//					onlineGoods = append(onlineGoods, goodsInfo.Name)
//					continue
//				}
//				dbGoods.BeginWrite()
//				dbGoods.NumberStore = numberStore
//				dbGoods.Status = int(mysqld.GoodsStatusEdited)
//				dbGoods.EndWrite()
//				dbGoods.Save()
//
//				dbGoodsInfo := dbGoods.GetGoodsInfo()
//				dbGoodsInfo.BeginWrite()
//				dbGoodsInfo.Name = goodsInfo.Name
//				dbGoodsInfo.BarCode = goodsInfo.BarCode
//				dbGoodsInfo.Desc = goodsInfo.Desc
//				dbGoodsInfo.Skus = goodsInfo.Skus
//
//				dbGoodsInfo.CategoryId = goodsInfo.CategoryId
//				dbGoodsInfo.EnterPrice = goodsInfo.EnterPrice
//				dbGoodsInfo.MinPrice = goodsInfo.MinPrice
//				dbGoodsInfo.OriginPrice = goodsInfo.OriginPrice
//
//				dbGoodsInfo.SellPrice = goodsInfo.SellPrice
//				dbGoodsInfo.EnterDate = goodsInfo.EnterDate
//				dbGoodsInfo.ProductDate = goodsInfo.ProductDate
//				dbGoodsInfo.ResId = goodsInfo.ResId
//				dbGoodsInfo.EndWrite()
//			}else{
//				mysqld.AddNewDBGoodsByInfo(goodsInfo)
//				dbGoods = mysqld.GetDBGoodsOrFromDB(goodsInfo.ID)
//
//				dbGoods.BeginWrite()
//				dbGoods.NumberStore = numberStore
//				dbGoods.Status = int(mysqld.GoodsStatusEdited)
//				dbGoods.EndWrite()
//				dbGoods.Save()
//			}
//		}
//		break
//	}
//	return true, onlineGoods
//}

//=>/gm/goods/load/picture true get {} 
//func Gm_goods_load_picture(ctx iris.Context, sess *common.BSSession) {
//	player := module.GetPlayer(sess)
//	if player == nil {
//		ctx.JSON(iris.Map{"code": -1, "msg":"未登陆"})
//		return
//	}
//	if player.GM == 0 {
//		ctx.JSON(iris.Map{"code": -1, "msg":"无权限"})
//		return
//	}
//	dbGoodsList := mysqld.GetDBGoodsList(-1)
//	var onlineGoods []string
//	for _,dbGoods := range dbGoodsList {
//		if dbGoods.Status == int(mysqld.GoodsStatusUp) {
//			continue
//		}
//	 	dbGoodsInfo := dbGoods.GetGoodsInfo()
//		if len(dbGoodsInfo.ResId) > 0 {
//			if module.MakeGoodsFiles(dbGoodsInfo) {
//				dbGoodsInfo.BeginWrite()
//				dbGoodsInfo.ResId = ""
//				dbGoodsInfo.EndWrite()
//				dbGoodsInfo.Save()
//			}else{
//				onlineGoods = append(onlineGoods, dbGoodsInfo.Name)
//			}
//		}
//	}
//	ctx.JSON(iris.Map{"code": 0, "data":onlineGoods})
//}
