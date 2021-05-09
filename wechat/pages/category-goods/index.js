const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')
const TOOLS = require('../../utils/tools.js')

const APP = getApp()

APP.caches = APP.caches || {}
APP.getCache = APP.getCache || function(key) {
    let cache = APP.caches[key]
    if (!cache) {
        return
    }
    if (cache.saveTime > Date.now()) {
        return cache.data
    }
}
APP.setCache = APP.setCache || function(key, val) {
    let cache = {
        data: val,
        saveTime: Date.now()+1000*60*3
    }
    APP.caches[key] = cache
}

Page({
    data: {
        listType: 1, // 1为1个商品一行，2为2个商品一行    
        name: '', // 搜索关键词
        // orderBy: '', // 排序规则
        loadingMoreHidden: true,
        groups: [
            {
                id: 1,
                value: 0,
                options: [
                    {text: '全部品牌', value: 0}
                ]
            },
            {
                id: 2,
                value: 0,
                options: [
                    {text: '默认排序', value: 0},
                    {text: '新品排序', value: 1},
                    {text: '销量排序', value: 2},
                    {text: '价格排序', value: 3}
                ]
            }
        ],
        searchValue: '',
        showSkuPanel: false,
        selectSkuIds: []
    },
    onLoad(options) {
        let data = {}
        if (options.name) {
            wx.setNavigationBarTitle({
                title: options.name
            })
            data.title = options.name
        }
        let listType = wx.getStorageSync("listType")
        if (!listType) {
            listType = 2
        }

        data.listType = listType
        data.categoryId = options.id
        if (options.subtype) {
            data.subtype = options.subtype
        }
        if (options.index) {
            data.index = options.index
        }
        this.setData(data)
        this.requestGoodsList()
    },
    onShareAppMessage() {
        const groups = this.data.groups
        let index = groups[0] && groups[0].value || 0
        let categoryId = this.data.categoryId
        let uid = wx.getStorageSync('uid');
        let path = '/pages/start/loading?inviter_id='+uid+'&route=/pages/category-goods/index%3fid%3d'+categoryId+'%26index%3d'+index+'%26name%3d'+encodeURI(this.data.title);
        return {
            title: this.data.title,
            path: path
        };
    },
    onReady() {
    },
    onShow() {
    },
    getNormalName(categoryId, index) {
        if (categoryId == 1) {
            if (index == 1) {
                return '全部品牌'
            } else {
                return '全部规格'
            }
        }
        if (index != 1) {
            return '二级规格'
        }
        return '全部商品'
    },
    getGroupSelectText(group) {
        for (var i = 0; i < group.options.length; i++) {
            let option = group.options[i]
            if (option.value == group.value) {
                return option.text
            }
        }
        console.log("getGroupSelectText error group =", group)
    },
    updateChange() {
        const groups = this.data.groups
        const originGoodsList = this.data.originGoodsList
        let group = groups[0]
        let goodsList
        if (group.value == 0) {
            goodsList = originGoodsList
        } else {
            goodsList = []
            let selectType = this.getGroupSelectText(group)
            for (let i = 0; i < originGoodsList.length; i++) {
                let goods = originGoodsList[i]
                if (goods.subType == selectType) {
                    goodsList.push(goods)
                }
            }
        }
        if (groups.length > 2) {
            group = groups[1]
            if (group.value > 0) {
                let selectType = this.getGroupSelectText(group)
                let tmps = goodsList
                goodsList = []
                for (let i = 0; i < tmps.length; i++) {
                    let goods = tmps[i]
                    if (goods.subType2 == selectType) {
                        goodsList.push(goods)
                    }
                }
            }
        }
        group = groups[groups.length-1]
        if (group.value > 0) {
            let orderBy = group.value
            goodsList.sort(function(a, b) {
                //新品排序
                if (orderBy == 1) {
                    return b.updateTime-a.updateTime
                //价格排序
                } else if (orderBy == 3) {
                    return a.sellPrice-b.sellPrice
                //销量排序
                } else {
                    return b.numberSell-a.numberSell
                }
            })
        }
        this.setData({
            goodsList: goodsList
        });
    },
    onChange1(event) {
        let group = this.data.groups[0]
        group.value = event.detail
        this.updateChange()
    },
    onChange2(event) {
        let group = this.data.groups[1]
        group.value = event.detail
        this.updateChange()
    },
    onChange3(event) {
        let group = this.data.groups[2]
        group.value = event.detail
        this.updateChange()
    },
    async requestGoodsList() {
        // wx.showLoading({title: '加载中'});
        let categoryId = this.data.categoryId
        let data = {
            page: 0,
            pageSize: 1000,
            categoryId: categoryId
        }

        let reqKey = "categorygoods"+categoryId
        let res = APP.getCache(reqKey)
        if(!res){
            wx.showLoading({title: '加载中'})
            res = await WXAPI.goods(data)
            APP.setCache(reqKey, res)
            wx.hideLoading()
        }
        let loadingMoreHidden = res.code != 10000
        if (res.code == 0 || res.code == 10000) {
            var goodsList = res.data || []
            let originGoodsList = goodsList
            var subTypes = {}
            var subType2s = {}
            for (var i = 0; i < goodsList.length; i++) {
                var goods = goodsList[i]
                if (!goods.subType || goods.subType.length == 0) {
                    goods.subType = "其他"
                }
                if (goods.subType2 && goods.subType2.length > 0) {
                    subType2s[goods.subType2] = true
                }
                subTypes[goods.subType] = true
            }

            let groups = []
            //typeGroup===>>
            let typeName = this.getNormalName(categoryId, 1)
            let typeGroup = {
                value: 0,
                options: [
                    {text: typeName, value: 0}
                ]
            }
            var idx = 1
            var isHasOther = false
            for (var key in subTypes) {
                if (key != '其他') {
                    var item = {
                        text: key, 
                        value: idx
                    };
                    idx++;
                    typeGroup.options.push(item)
                } else {
                    isHasOther = true
                }
            }
            if (isHasOther) {
                var item = {
                    text: "其他", 
                    value: idx
                };
                typeGroup.options.push(item)
            }
            groups.push(typeGroup)

            //type2Group===>>
            typeName = this.getNormalName(categoryId, 2)
            let options = [{text: typeName, value: 0}]
            idx = 1
            for (var key in subType2s) {
                var item = {
                    text: key, 
                    value: idx
                };
                idx++;
                options.push(item)
            }
            if (options.length > 1) {
                let type2Group = {
                    value: 0,
                    options: options,
                }
                groups.push(type2Group)
            }

            //sortGroup===>>
            let sortGroup = {
                value: 0,
                options: [
                    {text: '默认排序', value: 0},
                    {text: '新品排序', value: 1},
                    {text: '销量排序', value: 2},
                    {text: '价格排序', value: 3}
                ]
            }
            groups.push(sortGroup)
            for (var i = 0; i < groups.length; i++) {
                var group = groups[i]
                group.id = i+1
                group.changeFunc = 'onChange'+group.id
            }
            let subType = this.data.subtype
            if (subType) {
                this.data.subtype = null
                let group = groups[0]
                for (var i = 0; i < group.options.length; i++) {
                    let option = group.options[i]
                    if (option.text == subType) {
                        group.value = option.value
                        break
                    }
                }
                goodsList = []
                let selectType = this.getGroupSelectText(group)
                for (let i = 0; i < originGoodsList.length; i++) {
                    let goods = originGoodsList[i]
                    if (goods.subType == selectType) {
                        goodsList.push(goods)
                    }
                }
            }
            let index = this.data.index
            if (index) {
                this.data.index = null
                let group = groups[0]
                for (var i = 0; i < group.options.length; i++) {
                    let option = group.options[i]
                    if (option.value == index) {
                        group.value = option.value
                        break
                    }
                }
                goodsList = []
                let selectType = this.getGroupSelectText(group)
                for (let i = 0; i < originGoodsList.length; i++) {
                    let goods = originGoodsList[i]
                    if (goods.subType == selectType) {
                        goodsList.push(goods)
                    }
                }
            }
            this.setData({
                groups: groups,
                goodsList: goodsList,
                originGoodsList: originGoodsList,
                loadingMoreHidden: loadingMoreHidden
            });
        } else {
            wx.showToast({title: res.msg, icon: 'none'});
        }
    },
    onHide() {
    },
    onUnload() {
    },
    // onPullDownRefresh() {
    //     console.log("onPullDownRefresh====>>")
    // },
    // onReachBottom() {
    //     console.log("onReachBottom====>>")
    // },
    changeShowType() {
        let listType = 0
        if (this.data.listType == 1) {
            listType = 2
        } else {
            listType = 1
        }
        this.setData({
            listType: listType
        })
        wx.setStorageSync("listType", listType)
    },
    onSearch(e) {
        let value = e.detail
        if (!value || value.length == 0) {
            return
        }
        wx.navigateTo({url: '/pages/goods/list?name=' + value})
    },
    onClear(e) {
        this.updateChange()
    },

    async onAddShopCar(e) {
        const goodsId = e.currentTarget.dataset.id;
        let curGoods = null;
        for (var i = 0; i < this.data.goodsList.length; i++) {
            let goods = this.data.goodsList[i];
            if (goods.id == goodsId) {
                curGoods = goods;
                break;
            }
        }
        if (!curGoods) {
            wx.showToast({
                title: '商品列表无该商品',
                icon: 'none'
            })
            return;
        }
        if (curGoods.numberStore <= 0) {
            wx.showToast({
                title: '已售罄~',
                icon: 'none'
            })
            return
        }
        AUTH.checkHasLogined().then(isLogined => {
            this.setData({
                wxlogin: isLogined
            })
            if (isLogined) {
                // 处理加入购物车的业务逻辑
                this.addShopCarDone(curGoods)
            } else {
                wx.showToast({
                    title: '下单前，需要登陆',
                    icon: 'none'
                })
                AUTH.login()
            }
        })
    },
    async addShopCarDone(curGoods) {
        if (!curGoods.skuGroups) {
            let res = await WXAPI.goodsSku(curGoods.id);
            if (res.code != 0) {
                wx.showToast({ title: res.msg, icon: 'none' });
                return;
            }
            let data = res.data;
            if (curGoods.id != data.goodsId) {
                wx.showToast({ title: "物品ID不匹配", icon: 'none' });
                return
            }
            curGoods.numberStore = data.numberStore;
            curGoods.buyCount = data.buyCount;
            curGoods.buyLimit = data.buyLimit;
            curGoods.skuPics = data.skuPics;
            if (!TOOLS.isCanBuy(curGoods)) {
                return;
            }
            TOOLS.parseSkus(curGoods, data.skuJson);
        }
        this.data.selectSkuIds = [];
        this.data.curGoods = curGoods;
        curGoods.buyNumber = 1;

        this.updateSkuInfos();

        this.setData({
            curGoods: curGoods,
            showSkuPanel: true
        });
    },
    updateSkuInfos: function() {
        let curGoods = this.data.curGoods;
        let skuGroups = curGoods.skuGroups;

        let selectSkuIds = this.data.selectSkuIds;
        let selectSkuPrice = curGoods.sellPrice;
        let selectSkuIcon = curGoods.pic;
        for (var i = 0; i < skuGroups.length; i++) {
            let selectSkuId = selectSkuIds[i];
            let skuGroup = skuGroups[i];
            skuGroup.skuList.forEach(sku => {
                if (sku.id == selectSkuId) {
                    sku.active = true;
                    if (selectSkuId < 100) {
                        if (curGoods.skuPics && curGoods.skuPics[selectSkuId]) {
                            selectSkuIcon = curGoods.skuPics[selectSkuId]
                        }
                    }
                    if (sku.price && sku.price > 0) {
                        selectSkuPrice = sku.price
                    }
                }else{
                    sku.active = false
                }
            })
        }
        this.setData({
            skuGroups: skuGroups,
            selectSkuIcon: selectSkuIcon,
            selectSkuPrice: selectSkuPrice
        });
    },
    numJiaTap() {
        let curGoods = this.data.curGoods;
        if (curGoods.buyNumber < curGoods.numberStore) {
            curGoods.buyNumber++
            this.setData({
                curGoods
            })
        } else {
            wx.showToast({
                title: '库存有限',
                icon: 'none'
            })
        }
    },
    numJianTap() {
        const curGoods = this.data.curGoods;
        if (curGoods.buyNumber > 1) {
            curGoods.buyNumber--
            this.setData({
                curGoods
            })
        } else {
            wx.showToast({
                title: '不能小于1',
                icon: 'none'
            })
        }
    },
    onCloseSkuPanel() {
        this.setData({
            showSkuPanel: false
        })
    },
    async labelItemTap(e) {
        const skuId = e.currentTarget.dataset.id;
        let level = TOOLS.getLevel(skuId)-1;
        let skuGroups = this.data.skuGroups;
        let skuGroup = skuGroups[level];
        if (!skuGroup) {
            return
        }
        let selectSkuIds = this.data.selectSkuIds;
        selectSkuIds[level] = skuId;
        this.updateSkuInfos();
    },
    async addCarSku() {
        let data = this.data;
        let orderSkuId = TOOLS.makeOrderSkuId(data.curGoods, data.skuGroups, data.selectSkuIds);
        if (orderSkuId == null) {
            return;
        }
        const curGoods = this.data.curGoods;
        let goodsId = curGoods.id;
        let buyNumber = curGoods.buyNumber;
        const res = await WXAPI.cartAdd(goodsId, buyNumber, orderSkuId);
        if (res.code != 0) {
            wx.showToast({
                title: res.msg,
                icon: 'none'
            });
            return;
        }
        this.onCloseSkuPanel();
        wx.showToast({
            title: '加入购物车',
            icon: 'success'
        });
    },
    previewImage(e) {
        const url = e.currentTarget.dataset.url;
        if(typeof url == "string"){
            wx.previewImage({current: url, urls: [url]});
        }else{
            wx.previewImage({current: url[0], urls: url});
        }
    }
})