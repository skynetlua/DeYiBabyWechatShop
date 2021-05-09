const WXAPI = require('apifm-wxapi')

// 显示购物车tabBar的Badge
function showTabBarBadge() {
    const token = wx.getStorageSync('token')
    if (!token) {
        return
    }
    WXAPI.cartInfo().then(res => {
        if (res.code == 700) {
            wx.removeTabBarBadge({index: 2})
        }
        if (res.code == 0) {
            if (res.data.number == 0) {
                wx.removeTabBarBadge({index: 2})
            } else {
                wx.setTabBarBadge({index: 2, text: `${res.data.number}`})
            }
        }
    })
}

function datestr2time(datestr) {
    if (!datestr || datestr.length <= 1) {
        return 0
    }
    let items = datestr.split(' ')
    if (items.length < 2) {
        return 0
    }
    var theDate = new Date()
    let dateStr = items[0]
    let timeStr = items[1]
    items = dateStr.split('/')
    theDate.setFullYear(Number(items[0]))
    theDate.setMonth(Number(items[1])-1)
    theDate.setDate(Number(items[2]))
    items = timeStr.split(':')
    theDate.setHours(Number(items[0]))
    if (items[1]) {
        theDate.setMinutes(Number(items[1])-1)
    } else {
        theDate.setMinutes(0)
    }
    if (items[2]) {
        theDate.setSeconds(Number(items[2]))
    } else {
        theDate.setSeconds(0)
    }
    return Math.floor(theDate.getTime()/1000)
}

function time2datestr(time) {
    var theDate = new Date()
    if (time && time != 0) {
        theDate.setTime(time*1000)
    }
    let dateStr = theDate.getFullYear()+'/'+(theDate.getMonth()+1)+'/'+theDate.getDate()+' '+theDate.getHours()+':'+theDate.getMinutes()+':'+theDate.getSeconds()
    return dateStr
}

function time2dateformat(time) {
    var theDate = new Date()
    if (time && time != 0) {
        theDate.setTime(time*1000)
    }
    var monthStr = (theDate.getMonth()+1)+'';
    if (monthStr.length == 1) {
        monthStr = '0'+monthStr;
    }
    var dayStr = theDate.getDate()+'';
    if (dayStr.length == 1) {
        dayStr = '0'+dayStr;
    }
    let dateStr = theDate.getFullYear()+'-'+monthStr+'-'+dayStr
    return dateStr
}

function time2timeformat(time) {
    var theDate = new Date()
    if (time && time != 0) {
        theDate.setTime(time*1000)
    }
    var hourStr = theDate.getHours()+'';
    if (hourStr.length == 1) {
        hourStr = '0'+hourStr;
    }
    var minuteStr = theDate.getMinutes()+'';
    if (minuteStr.length == 1) {
        minuteStr = '0'+minuteStr;
    }
    let dateStr = hourStr+':'+minuteStr
    return dateStr
}

function getLevel(skuId) {
    let level = 1;
    while(skuId >= 100) {
        level++;
        skuId = skuId/100;
    }
    return level;
}

function isCanBuy(goods) {
    if (goods.buyLimit && goods.buyLimit > 0) {
        if (goods.buyCount && goods.buyCount > 0) {
            if (goods.buyCount >= goods.buyLimit) {
                wx.showToast({title: '超过限购，限购'+goods.buyLimit+'次数',icon: 'none'})
                return false
            }
        }
    }
    return true
}

function parseSkus(curGoods, skuJson) {
    var skuGroups = [];
    var skuList = [];
    if (skuJson && skuJson.length > 0) {
        skuList = JSON.parse(skuJson);

        // skuList.push({"label":"规格2","id":100,"name":"2NB码96片","price":0});
        // skuList.push({"label":"规格2","id":200,"name":"2NB码96片","price":0});
        // skuList.push({"label":"规格2","id":300,"name":"3NB码96片","price":0});
        // skuList.push({"label":"规格2","id":400,"name":"3NB码96片","price":0});

        let skuGroupMap = {};
        for (var i = 0; i < skuList.length; i++) {
            let sku = skuList[i];
            let skuGroup = skuGroupMap[sku.label];
            let level = getLevel(sku.id)
            if (!skuGroup) {
                skuGroup = {
                    label:sku.label,
                    level:level,
                    skuList:[]
                };
                skuGroupMap[sku.label] = skuGroup;
                skuGroups.push(skuGroup);
            } else {
                if (skuGroup.level != level) {
                    console.log("parseSkus skuGroup.level =", skuGroup.level, "level =", level);
                }
            }
            skuGroup.skuList.push(sku)
        }
    }
    curGoods.skuGroups = skuGroups;
    curGoods.skuList = skuList;
}

function makeOrderSkuId(curGoods, skuGroups, selectSkuIds) {
    if (curGoods.buyNumber < 1) {
        wx.showModal({
            title: '提示',
            content: '亲，购买数量不能为0喔',
            showCancel: false
        })
        return null;
    }
    let count = skuGroups.length;
    if (count == 0) {
        return 0;
    }
    for (var i = 0; i < count; i++) {
        let skuGroup = skuGroups[i];
        let selectSkuId = selectSkuIds[i];
        let isFind = false;
        if (selectSkuId && selectSkuId > 0) {
            for (var j = 0; j < skuGroup.skuList.length; j++) {
                let sku = skuGroup.skuList[j];
                if (sku.id == selectSkuId) {
                    if (sku.active) {
                        isFind = true;
                    }
                    break
                }
            }
        }
        if (!isFind) {
            wx.showModal({
                title: '提示',
                content: '请选择商品规格！',
                showCancel: false
            });
            return null;
        }
    }

    let orderSkuId = selectSkuIds[0];
    if (count > 1) {
        orderSkuId = orderSkuId+selectSkuIds[1];
    }
    if (count > 2) {
        orderSkuId = orderSkuId+selectSkuIds[2];
    }
    return orderSkuId
}

module.exports = {
    getLevel: getLevel,
    isCanBuy: isCanBuy,
    parseSkus: parseSkus,
    makeOrderSkuId: makeOrderSkuId,
    showTabBarBadge: showTabBarBadge,
    time2timeformat: time2timeformat,
    time2dateformat: time2dateformat,
    time2datestr: time2datestr
}