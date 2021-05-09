package handle
import (
	"bestsell/common"
	"bestsell/config"
	"bestsell/module"
	"bestsell/mysqld"
	"fmt"
	"github.com/kataras/iris/v12"
	"strconv"
)

//=>/address/add true post {} 
func Address_add(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	addressBox := player.GetAddressBox()
	dbAddresss := addressBox.GetAddresss()
	for _,dbAddr := range *dbAddresss {
		if dbAddr.IsDefault > 0 {
			dbAddr.BeginWrite()
			dbAddr.IsDefault = 0
			dbAddr.EndWrite()
			dbAddr.DelaySave(dbAddr)
		}
	}
	linkMan := ctx.FormValue("linkMan")
	address := ctx.FormValue("address")
	mobile := ctx.FormValue("mobile")
	code := ctx.FormValue("code")
	provinceId := common.AtoI(ctx.FormValue("provinceId"))
	cityId :=  common.AtoI(ctx.FormValue("cityId"))
	areaId :=  common.AtoI(ctx.FormValue("areaId"))

	dbAddr := &mysqld.DBAddress {
		LinkMan:linkMan,
		Address:address,
		Mobile:mobile,
		Code:code,
		ProvinceId:int64(provinceId),
		CityId:int64(cityId),
		AreaId:int64(areaId),
		IsDefault:1,
	}
	addressBox.AddAddress(dbAddr)
	ctx.JSON(iris.Map{"code": 0})
}

//=>/address/update true post {} 
func Address_update(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	_id := ctx.FormValue("id")
	id,err := strconv.Atoi(_id)
	if err != nil {
		fmt.Println("Address_update ",err)
		ctx.JSON(iris.Map{"code": -1, "msg":"id出错"})
		return
	}
	addressBox := player.GetAddressBox()
	dbAddresss := addressBox.GetAddresss()
	for _,dbAddr := range *dbAddresss {
		if dbAddr.IsDefault > 0 {
			dbAddr.BeginWrite()
			dbAddr.IsDefault = 0
			dbAddr.EndWrite()
			dbAddr.DelaySave(dbAddr)
		}
	}
	dbAddr := addressBox.GetAddress(id)
	isDefault := ctx.FormValue("isDefault")
	if len(isDefault) > 0 {
		if dbAddr == nil {
			ctx.JSON(iris.Map{"code": -1, "msg":"地址不存在"})
			return
		}
		dbAddr.BeginWrite()
		dbAddr.IsDefault = 1
		dbAddr.EndWrite()
		dbAddr.DelaySave(dbAddr)
		ctx.JSON(iris.Map{"code": 0, "msg":"默认地址修改成功"})
		return
	}

	linkMan := ctx.FormValue("linkMan")
	address := ctx.FormValue("address")
	mobile := ctx.FormValue("mobile")
	code := ctx.FormValue("code")
	provinceId := common.AtoI(ctx.FormValue("provinceId"))
	cityId :=  common.AtoI(ctx.FormValue("cityId"))
	areaId :=  common.AtoI(ctx.FormValue("areaId"))
	if dbAddr == nil {
		dbAddr = &mysqld.DBAddress{
			LinkMan:linkMan,
			Address:address,
			Mobile:mobile,
			Code:code,
			ProvinceId:int64(provinceId),
			CityId:int64(cityId),
			AreaId:int64(areaId),
			IsDefault:1,
		}
		dbAddr.Insert()
		addressBox.AddAddress(dbAddr)
	}else{
		dbAddr.BeginWrite()
		dbAddr.LinkMan = linkMan
		dbAddr.Address = address
		dbAddr.Mobile = mobile
		dbAddr.Code = code
		dbAddr.ProvinceId = int64(provinceId)
		dbAddr.CityId = int64(cityId)
		dbAddr.AreaId = int64(areaId)
		dbAddr.IsDefault = 1
		dbAddr.EndWrite()
		dbAddr.DelaySave(dbAddr)
	}
	ctx.JSON(iris.Map{"code": 0, "msg":"地址修改成功"})
}

//=>/address/delete true post {id,token} 
func Address_delete(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	_id := ctx.FormValue("id")
	id,err := strconv.Atoi(_id)
	if err != nil {
		fmt.Println("Address_delete ",err)
		ctx.JSON(iris.Map{"code": -1, "msg":"id出错"})
		return
	}
	addressBox := player.GetAddressBox()
	addressBox.RemoveAddress(id)
	ctx.JSON(iris.Map{"code": 0})
}

//=>/address/list true get {token} 
func Address_list(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	addressBox := player.GetAddressBox()
	dbAddresss := addressBox.GetAddresss()
	var addresss []interface{}
	for _,dbAddr := range *dbAddresss {
		item := &map[string]interface{}{
			"id":dbAddr.ID,
			"linkMan":dbAddr.LinkMan,
			"address":dbAddr.Address,
			"mobile":dbAddr.Mobile,
			"isDefault":dbAddr.IsDefault,
			"province":config.GetCfgAddressName(int(dbAddr.ProvinceId)),
			"city":config.GetCfgAddressName(int(dbAddr.CityId)),
			"area":config.GetCfgAddressName(int(dbAddr.AreaId)),
		}
		addresss = append(addresss, item)
	}
	data := map[string]interface{}{
		"addresss":addresss,
	}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/address/default true get {token} 
func Address_default(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	addressBox := player.GetAddressBox()
	dbAddr := addressBox.GetDefaultAddress()
	if dbAddr == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	data := &map[string]interface{}{
		"id" :dbAddr.ID,
		"linkMan" :dbAddr.LinkMan,
		"address" :dbAddr.Address,
		"mobile" :dbAddr.Mobile,
		"isDefault" :dbAddr.IsDefault,
		"provinceId" :dbAddr.ProvinceId,
		"cityId" :dbAddr.CityId,
		"areaId" :dbAddr.AreaId,
		"province" :config.GetCfgAddressName(int(dbAddr.ProvinceId)),
		"city" :config.GetCfgAddressName(int(dbAddr.CityId)),
		"area" :config.GetCfgAddressName(int(dbAddr.AreaId)),
	}
	//data := map[string]interface{}{
	//	"info":info,
	//}
	ctx.JSON(iris.Map{"code": 0, "data": data})
}

//=>/address/detail true get {id,token} 
func Address_detail(ctx iris.Context, sess *common.BSSession) {
	player :=  module.GetPlayer(sess)
	if player == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	_id := ctx.FormValue("id")
	id,err := strconv.Atoi(_id)
	if err != nil {
		fmt.Println("Address_detail ",err)
		ctx.JSON(iris.Map{"code": -1, "msg":"id出错"})
		return
	}
	addressBox := player.GetAddressBox()
	dbAddr := addressBox.GetAddress(id)
	if dbAddr == nil {
		ctx.JSON(iris.Map{"code": -1})
		return
	}
	info := &map[string]interface{}{
		"id":dbAddr.ID,
		"linkMan":dbAddr.LinkMan,
		"address":dbAddr.Address,
		"mobile":dbAddr.Mobile,
		"isDefault":dbAddr.IsDefault,
		"provinceId":dbAddr.ProvinceId,
		"cityId":dbAddr.CityId,
		"areaId":dbAddr.AreaId,
		"province":config.GetCfgAddressName(int(dbAddr.ProvinceId)),
		"city":config.GetCfgAddressName(int(dbAddr.CityId)),
		"area":config.GetCfgAddressName(int(dbAddr.AreaId)),
	}
	data := map[string]interface{}{
		"info": info,
	}
	ctx.JSON(iris.Map{"code": 0, "data":data})
}