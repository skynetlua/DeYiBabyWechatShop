const CONFIG = require('../../config.js')
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
        saveTime: Date.now()+1000*60*1
    }
    APP.caches[key] = cache
}

Page({
    /**
     * 页面的初始数据
     */
    data: {
        categories: [],
        selectCategoryId: null,
        currentGoods: [],
        scrolltop: 0,
        skuCurGoods: null,
        curGoods: null,
        loadingMoreHidden: true,
        groups: {},
        caches: {},
        subTypesList: {},
        curCategory: {},
        nextTime: 0,
        isInit:false
    },
    onLoad(options) {
        wx.showShareMenu({
            withShareTicket: true
        })
        this.loadData()
    },
    async categories() {
        const res = await WXAPI.goodsCategory()
        let categories = []
        let selectCategoryId = this.data.selectCategoryId
        if (res.code == 0) {
            let datas = res.data || []
            datas.sort(function(a, b){
                return a.order - b.order
            })
            if (datas.length > 0 && selectCategoryId) {
                const curCategory = datas.find(ele => {
                    return ele.id == selectCategoryId
                })
                selectCategoryId = curCategory.id
            }
            for (let i = 0; i < datas.length; i++) {
                let item = res.data[i]
                categories.push(item)
                if (i == 0 && !selectCategoryId) {
                    selectCategoryId = item.id
                }
                // groups[item.id] = {
                //     categoryId: item.id,
                //     goods: [],
                //     curPage: 0,
                //     pageSize: 20
                // }
            }
            this.setData({
                // groups: groups,
                categories: categories,
                selectCategoryId: selectCategoryId,
                subTypesList: {}
            })
            this.getSubTypesList()

            let curTime = Date.now()
            wx.setStorageSync('category_category_time', curTime)
        }
    },
    // async getGoodsList() {
    //     let categoryId = this.data.selectCategoryId
    //     let group = this.data.groups[categoryId]
    //     let params = {
    //         page: group.curPage,
    //         pageSize: group.pageSize,
    //         categoryId: categoryId
    //     }
    //     let reqKey = [categoryId, params.page, params.pageSize].join(':')
    //     let res = this.getCache(reqKey)
    //     if(res){
    //         // console.log("getGoodsList=====>> 无需下载")
    //         // wx.hideLoading()
    //     } else {
    //         // console.log("getGoodsList=====>> 下载")
    //         wx.showLoading({title: '加载中'})
    //         res = await WXAPI.goods(params)
    //         this.setCache(reqKey, res)
    //         wx.hideLoading()
    //     }
    //     var datas = res.data || []
    //     let goodsMap = {}
    //     let goods = group.goods
    //     for (let i = 0; i < goods.length; i++) {
    //         let item = goods[i]
    //         goodsMap[item.id] = item
    //     }
    //     for (let i = 0; i < datas.length; i++) {
    //         let item = datas[i]
    //         if (!goodsMap[item.id]) {
    //             goods.push(item);
    //         }
    //     }
    //     let loadingMoreHidden = es.code != 10000
    //     this.setData({
    //         currentGoods:goods,
    //         loadingMoreHidden: loadingMoreHidden
    //     })
    // },
    async getSubTypesList() {
        let categoryId = this.data.selectCategoryId
        if (!categoryId) {
            categoryId = 1
        }
        let subTypeList = this.data.subTypesList[categoryId]
        if (!subTypeList) {
            let reqKey = "category"+categoryId
            let res = APP.getCache(reqKey)
            if(!res){
                wx.showLoading({title: '加载中'})
                res = await WXAPI.goodsSubtypes(categoryId)
                APP.setCache(reqKey, res)
                wx.hideLoading()
            }
            // wx.showLoading({title: '加载中'})
            // let res = await WXAPI.goodsSubtypes(categoryId)
            // wx.hideLoading()
            subTypeList = res.data || []
            this.data.subTypesList[categoryId] = subTypeList
        }
        let curCategory
        for (var i = 0; i < this.data.categories.length; i++) {
            let category = this.data.categories[i]
            if (category.id == categoryId) {
                curCategory = category
                break
            }
        }

        let subTypeGroup = []
        let groupMap = {}
        for (var i = 0; i < subTypeList.length; i++) {
            let item = subTypeList[i];
            if (!item.mainType || item.mainType.length == 0) {
                item.mainType = curCategory.name;
            }
            let group = groupMap[item.mainType]
            if (!group) {
                group = {
                    name:item.mainType == '-' ? curCategory.name:item.mainType,
                    mainType:item.mainType,
                    list:[]
                }
                groupMap[item.mainType] = group
                subTypeGroup.push(group)
            }
            group.list.push(item)
        }

        this.setData({
            curCategory:curCategory || {},
            // subTypeList:subTypeList || [],
            subTypeGroup:subTypeGroup
        })
        this.data.isInit = true
    },
    // onReachBottom() {
    //     if (!this.data.loadingMoreHidden) {
    //         return
    //     }
    //     let categoryId = this.data.selectCategoryId
    //     let group = this.data.groups[categoryId]
    //     group.curPage = group.curPage + 1
    //     this.setData({
    //         groups: this.data.groups
    //     })
    //     this.getGoodsList(true)
    // },
    // onPullDownRefresh() {
    //     // this.data.loadingMoreHidden = false
    //     // let categoryId = this.data.selectCategoryId
    //     // let group = this.data.groups[categoryId]
    //     // group.curPage = 0
    //     // group.goods = []
    //     // this.setData({
    //     //     groups: this.data.groups
    //     // })
    //     // this.getGoodsList(true)
    //     this.data.subTypesList = {}
    //     this.getSubTypesList()
    //     wx.stopPullDownRefresh()
    // },
    // toDetailsTap(e) {
    //     wx.navigateTo({url: "/pages/goods-details/index?id=" + e.currentTarget.dataset.id})
    // },
    onCategoryClick(e) {
        this.setSelectCategoryId(e.target.dataset.id)
    },
    bindinput(e) {
        this.setData({inputVal: e.detail})
    },
    bindconfirm(e) {
        let value = e.detail;
        if (!value || value.length == 0) {
            return;
        }
        this.setData({inputVal: value})
        wx.navigateTo({url: '/pages/goods/list?name=' + value + "&id="+this.data.selectCategoryId})
    },
    searchTap() {
        let value = this.data.inputVal;
        if (!value || value.length == 0) {
            return;
        }
        wx.navigateTo({url: '/pages/goods/list?name=' + value})
    },
    onShareAppMessage() {
        return {
            title: '"' + wx.getStorageSync('mallName') + '" ' + CONFIG.shareProfile,
            path: '/pages/start/loading?inviter_id=' + wx.getStorageSync('uid') + '&route=/pages/category/category'
        }
    },
    setSelectCategoryId(categoryId) {
        if (categoryId === this.data.selectCategoryId) {
            this.setData({scrolltop: 0})
        } else {
            var curCategoryId
            var categories = this.data.categories
            for (var i = 0; i < categories.length; i++) {
                let item = categories[i]
                if (item.id == categoryId) {
                    curCategoryId = categoryId
                    break
                }
            }
            this.setData({
                selectCategoryId: curCategoryId,
                scrolltop: 0
            });
            wx.setStorageSync('_categoryId', curCategoryId)
            // this.getGoodsList()
            this.getSubTypesList()
        }
    },
    loadData() {
        let curTime = Date.now()
        if (curTime < this.data.nextTime) {
            return
        }
        this.data.nextTime = curTime+1000*2

        let tabBarBadge_time = wx.getStorageSync('tabBarBadge_time') || 0
        tabBarBadge_time = Number(tabBarBadge_time)
        if (tabBarBadge_time+1000*60*3 < curTime) {
            wx.setStorageSync('tabBarBadge_time', curTime)
            AUTH.checkHasLogined().then(isLogined => {
                if (isLogined) {
                    this.setData({wxlogin: isLogined})
                    TOOLS.showTabBarBadge()
                }
            })
        }

        const _categoryId = wx.getStorageSync('_categoryId')
        let category_time = wx.getStorageSync('category_category_time') || 0
        category_time = Number(category_time)
        if (this.data.isInit && category_time+1000*60*1 > curTime) {
            if (_categoryId && _categoryId != this.data.selectCategoryId) {
                this.setSelectCategoryId(_categoryId)
            }
            return
        }
        if (_categoryId) {
            this.data.selectCategoryId = _categoryId
        }
        this.categories()
    },
    onShow() {
        this.loadData()
    },
    async addShopCar(e) {
        const goodsId = e.currentTarget.dataset.id
        const curGoods = this.data.currentGoods.find(ele => {
            return ele.id == goodsId
        })
        if (!curGoods) {
            wx.showToast({
                title: '商品列表无该商品',
                icon: 'none'
            })
            this.getGoodsList();
            return
        }
        if (curGoods.numberStore <= 0) {
            wx.showToast({
                title: '已售罄~',
                icon: 'none'
            })
            return
        }
        this.setData({
            curGoods: curGoods
        })
        this.addShopCarCheck({
            goodsId: goodsId,
            buyNumber: 1,
            skuId: -1
        })
    },
    async addShopCarCheck(options) {
        AUTH.checkHasLogined().then(isLogined => {
            this.setData({
                wxlogin: isLogined
            })
            if (isLogined) {
                this.addShopCarDone(options)
            }else{
                wx.showToast({
                    title: '下单前，需要登陆',
                    icon: 'none'
                })
                AUTH.login()
            }
        })
    },
    async addShopCarDone(options) {
        const res = await WXAPI.cartAdd(options.goodsId, options.buyNumber, options.skuId)
        if (res.code == 30002) {
            const skuRes = await WXAPI.goodsSku(options.goodsId)
            if (skuRes.code != 0) {
                wx.showToast({
                    title: skuRes.msg,
                    icon: 'none'
                })
                return
            }
            wx.hideTabBar();
            const skuCurGoods = skuRes.data;
            skuCurGoods.storesBuy = 1;
            this.setData({
                skuCurGoods
            });
            return;
        }
        if (res.code != 0) {
            wx.showToast({
                title: res.msg,
                icon: 'none'
            })
            return
        }
        wx.showToast({
            title: '加入成功',
            icon: 'success'
        })
        this.setData({
            skuCurGoods: null
        })
        wx.showTabBar()
        TOOLS.showTabBarBadge()
    },
    storesJia() {
        const skuCurGoods = this.data.skuCurGoods
        if (skuCurGoods.storesBuy < skuCurGoods.numberStore) {
            skuCurGoods.storesBuy++
            this.setData({
                skuCurGoods
            })
        }else{
            wx.showToast({
                title: '库存有限',
                icon: 'none'
            })
        }
    },
    storesJian() {
        const skuCurGoods = this.data.skuCurGoods
        if (skuCurGoods.storesBuy > 1) {
            skuCurGoods.storesBuy--
            this.setData({
                skuCurGoods
            })
        }else{
            wx.showToast({
                title: '不能小于1',
                icon: 'none'
            })
        }
    },
    closeSku() {
        this.setData({
            skuCurGoods: null
        })
        wx.showTabBar()
    },
    skuSelect(e) {
        const skuId = e.currentTarget.dataset.id
        const skuCurGoods = this.data.skuCurGoods
        let skuInfo = null
        skuCurGoods.skus.forEach(sku => {
            sku.active = sku.id == skuId
            if (sku.active) {
                skuInfo = sku
            }
        })
        if (skuInfo) {
            skuCurGoods.selectSku = skuInfo
        }
        this.setData({
            skuCurGoods
        })
    },
    addCarSku() {
        const skuCurGoods = this.data.skuCurGoods
        let skuId = 0;
        if (skuCurGoods.skus && skuCurGoods.skus.length > 1) {
            const skuItem = skuCurGoods.skus.find(ele => { return ele.active })
            if (!skuItem) {
                wx.showToast({
                    title: '请选择规格',
                    icon: 'none'
                })
                return;
            }
            skuId = skuItem.id;
        }
        const options = {
            goodsId: skuCurGoods.goodsId,
            buyNumber: skuCurGoods.storesBuy,
            skuId: skuId
        }
        this.addShopCarDone(options)
    }
    // onReachBottom() {
    //     console.log("onReachBottom==========>>")
    // },
    // onPullDownRefresh() {
    //     console.log("onPullDownRefresh==========>>")
    // }
})