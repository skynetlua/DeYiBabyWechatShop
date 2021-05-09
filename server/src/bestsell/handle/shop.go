package handle

import (
	"bestsell/common"
	// "bestsell/config"
	// "fmt"
	"github.com/kataras/iris/v12"
	// "strconv"
)


//=>/shop/goods/category/info true get {id} 
func Shop_goods_category_info(ctx iris.Context, sess *common.BSSession) {
	empty("/shop/goods/category/info")
}

//=>/shop/goods/price true post {goodsId,propertyChildIds} 
func Shop_goods_price(ctx iris.Context, sess *common.BSSession) {
	// _goodsId := ctx.FormValue("goodsId")
	// _propertyChildIds := ctx.FormValue("propertyChildIds")
	// fmt.Println("_propertyChildIds:", _propertyChildIds)
	// goodsId,err := strconv.Atoi(_goodsId)
	// if err != nil {
	// 	fmt.Println("Shop_goods_price err =", err)
	// 	ctx.JSON(iris.Map{"code": -1})
	// 	return
	// }
	// cfgGoods := (*config.GetCfgGoodsMap())[goodsId]
	// if cfgGoods == nil {
	// 	ctx.JSON(iris.Map{"code": -1})
	// 	return
	// }
	// data := map[string]interface{}{
	// 	"price":cfgGoods.MinPrice,
	// 	"originalPrice":cfgGoods.OriginalPrice,
	// 	//"score":cfgGoods.MinScore,
	// 	"stores": cfgGoods.Stores,
	// }
	// ctx.JSON(iris.Map{"code": 0, "data": data})
}