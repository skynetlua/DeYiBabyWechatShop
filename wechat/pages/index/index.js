const WXAPI = require('apifm-wxapi')
const CONFIG = require('../../config.js')
const TOOLS = require('../../utils/tools.js')
const AUTH = require('../../utils/auth')

const APP = getApp()
// fixed首次打开不显示标题的bug
APP.configLoadOK = () => {
  wx.setNavigationBarTitle({
    title: wx.getStorageSync('mallName')
  })
}

//获取应用实例
Page({
    data: {
        inputVal: "", // 搜索框内容
        recomGoods: [], // 推荐商品
        // kanjiaList: [], //砍价商品列表
        // pingtuanList: [], //拼团商品列表
        // loadingHidden: false, // loading
        // selectCurrent: 0,
        categories: [],
        // activeCategoryId: 0,
        goods: [],
        // scrollTop: 0,
        loadingMoreHidden: true,
        goodsGroups: [],
        curPage: 0,
        pageSize: 6,
        waitCount: 0,
        nextTime: 0
    },
    tabClick(e) {
        wx.setStorageSync("_categoryId", e.currentTarget.id)
        wx.switchTab({url: '/pages/category/category'})
    },
    tapBanner(e) {
        const url = e.currentTarget.dataset.url
        if (url) {
            wx.navigateTo({url})
        }
    },
    loadData() {
        let curTime = Date.now()
        if (curTime < this.data.nextTime) {
            return
        }
        this.data.nextTime = curTime+1000*2

        // let tabBarBadge_time = wx.getStorageSync('tabBarBadge_time') || 0
        // tabBarBadge_time = Number(tabBarBadge_time)
        // if (tabBarBadge_time+1000*60*1 < curTime) {
        //     wx.setStorageSync('tabBarBadge_time', curTime)
        //     AUTH.checkHasLogined().then(isLogined => {
        //         if (isLogined) {
        //             this.setData({wxlogin: isLogined})
        //             TOOLS.showTabBarBadge()
        //         }
        //     })
        // }

        let index_time = wx.getStorageSync('index_index_time') || 0
        index_time = Number(index_time)
        if (index_time+1000*60*1 > curTime) {
            return
        }
        let isConnected = wx.getStorageSync('isConnected')
        if (isConnected && AUTH.isLogined()) {
            // console.log("index/index  已登陆");
            this.enter();
        }else{
            // console.log("start/loading  未登陆");
            this.getServerInfo();
        }
    },
    onLoad(options) {
        wx.showShareMenu({
            withShareTicket: true
        })
        if (options) {
            let route = options.route
            if (options && options.scene) {
                const scene = decodeURIComponent(options.scene)
                const tmps = scene.split(',')
                if (tmps[0] == 'qrcode-index') {
                    route = '/pages/start/start' + tmps[1]
                    if (tmps.length > 2) {
                        wx.setStorageSync('referrer', tmps[2])
                    }
                }else if (tmps[0] == 'qrcode-goods') {
                    route = '/pages/goods-details/index' + tmps[1]
                    if (tmps.length > 2) {
                        wx.setStorageSync('referrer', tmps[2])
                    }
                }
            }
            if (options && options.inviter_id) {
                wx.setStorageSync('referrer', options.inviter_id)
                if (options.shareTicket) {
                    // this.processShareTicket(options.inviter_id, options.shareTicket)
                }
            }
        }
        this.loadData()
    },
    checkLogin(){
        if (isNaN(this.data.waitCount)) {
            this.data.waitCount = 0
        }
        this.data.waitCount++
        if (this.data.waitCount>10) {
            AUTH.login();
             // this.login()
        }
        const isLogined = AUTH.isLogined();
        if (!isLogined) {
            // console.log("start/loading  未登陆 waitCount =", this.data.waitCount)
            setTimeout(() => {
                this.checkLogin();
            }, 1000)
            return
        }
        this.enter();
    },
    getServerInfo() {
        const accountInfo = wx.getAccountInfoSync()
        const appid = accountInfo.miniProgram.appId
        wx.setStorageSync('wxAppid', appid)
        const version = CONFIG.version
        // console.log("start/loading getServerInfo===> appid =", appid)
        const _this = this
        WXAPI.request('/subdomain/appid/wxapp', false, 'get', {appid, version}, function(url){
            setTimeout(() => {
                wx.showModal({
                    title: '错误提示',
                    content: '连接不上服务器，错误代码:1',
                    showCancel: false,
                    confirmText: '重试',
                    success(res) {
                        _this.getServerInfo()
                    }
                })
            }, 3000)
        }).then(function(res){
            if (res.code != 0) {
                wx.showModal({
                    title: '错误提示',
                    content: res.msg,
                    showCancel: false,
                    confirmText: '重试',
                    success(res) {
                        _this.getServerInfo()
                    }
                })
            } else {
                // console.log("start/loading getServerInfo===> 3");
                const data = res.data;
                WXAPI.init2(data.host, data.subdomain)
                for (var i = 0; i < data.config.length; i++) {
                    let item = data.config[i];
                    wx.setStorageSync(item.key, item.value);
                }
                wx.setNavigationBarTitle({
                  title: wx.getStorageSync('mallName')
                }) 
                // wx.setNavigationBarTitle({
                //     title: wx.getStorageSync('mallName')
                // });
                AUTH.login();
                _this.checkLogin();
            }
        });
    },
    onShow(e) {
        this.loadData()
        let animation = wx.createAnimation({
            duration: 1000,
            timingFunction: 'linear',
            delay: 0,
            transformOrigin: '50% 50%',
            success: function(res) {
             console.log("res")
            }
        })
        if (this.timerHandle) {
            clearInterval(this.timerHandle)
        }

        let index = 1.3;
        let that = this;
        this.timerHandle = setInterval(function() {
            if (index > 1) {
                index = 1.0
            } else {
                index = 1.3
            }
            animation.scale(index).step()
            that.setData({
                animation: animation.export()
            })
        }, 1000)
    },
    onHide() {
        if (this.timerHandle) {
            clearInterval(this.timerHandle)
            this.timerHandle = null
        }
    },
    enter(){
        wx.setNavigationBarTitle({
            title: wx.getStorageSync('mallName')
        })
        this.requestData()
    },
    onPageScroll(e) {
        // let scrollTop = this.data.scrollTop
        // this.setData({
        //     scrollTop: e.scrollTop
        // })
    },
    onShareAppMessage() {
        return {
            title: '"' + wx.getStorageSync('mallName') + '" ' + CONFIG.shareProfile,
            path: '/pages/start/loading?inviter_id=' + wx.getStorageSync('uid') + '&route=/pages/index/index'
        }
    },
    // onReachBottom() {
    //     if (!this.data.loadingMoreHidden) {
    //         return
    //     }
    //     this.setData({
    //         curPage: this.data.curPage + 1
    //     })
    //     this.getGoodsList(this.data.activeCategoryId, true)
    // },
    // onPullDownRefresh() {
    //     this.setData({curPage: 1})
    //     this.getGoodsList(this.data.activeCategoryId)
    //     wx.stopPullDownRefresh()
    // },
    async wxaMpLiveRooms() {
        const res = await WXAPI.wxaMpLiveRooms()
        if (res.code == 0 && res.data.length > 0) {
            this.setData({
                aliveRooms: res.data
            })
        }
    },
    goGM(){
        wx.navigateTo({url: "/pages/gm/index"})
    },
    async requestData() {
        const res = await WXAPI.pageIndex()
        let resData = res.data
        let categories = resData.categories || []
        let recomGoods = resData.recomGoods || []
        let seckillGoods = resData.seckillGoods || []
        let teamGoods = resData.teamGoods || []
        let banners    = resData.banners || []
        let noticeList = resData.noticeList || []
        let goodsDynamic = resData.goodsDynamic || []
        let gm = resData.gm
        wx.setStorageSync('gm', gm)

        if (resData.cartCount && resData.cartCount > 0) {
            wx.setTabBarBadge({index: 2, text: `${resData.cartCount}`});
        } else {
            wx.removeTabBarBadge({index: 2});
        }

        var curTime = new Date().getTime()
        curTime = Math.floor(curTime/1000)

        for (let i = 0; i < seckillGoods.length; i++) {
            let goods = seckillGoods[i]
            if (goods.startTime && goods.startTime > 0) {
                goods.startInterval = (goods.startTime-curTime)*1000
            } else {
                goods.startInterval = -1
            }
            if (goods.endTime && goods.endTime > 0) {
                goods.endInterval = (goods.endTime-curTime)*1000
            } else {
                goods.endInterval = -1
            }
        }

        // for (let i = 0; i < teamGoods.length; i++) {
        //     let goods = teamGoods[i]
        //     if (goods.startTime && goods.startTime > 0) {
        //         goods.startInterval = (goods.startTime-curTime)*1000
        //     } else {
        //         goods.startInterval = -1
        //     }
        //     if (goods.endTime && goods.endTime > 0) {
        //         goods.endInterval = (goods.endTime-curTime)*1000
        //     } else {
        //         goods.endInterval = -1
        //     }
        // }

        recomGoods.sort(function(a, b){
            return a.order-b.order
        })
        seckillGoods.sort(function(a, b){
            return a.order-b.order
        })
        teamGoods.sort(function(a, b){
            return a.order-b.order
        })
        categories.sort(function(a, b){
            return a.order - b.order
        })

        let tmpGroups = res.data.goodsGroups || [];
        let goodsGroups = [];
        for (var i = 0; i < categories.length; i++) {
            let cItem = categories[i];
            for (var j = 0; j < tmpGroups.length; j++) {
                let group = tmpGroups[j];
                if (group.categoryId == cItem.id) {
                    group.name = cItem.name;
                    group.pic = cItem.pic;
                    group.goodsList = group.goodsList || []
                    if (group.categoryId == 10) {
                        goodsGroups.unshift(group)
                    } else {
                        goodsGroups.push(group)
                    }
                    break
                }
            }
        }

        this.setData({
            categories: categories,
            noticeList: noticeList,
            recomGoods: recomGoods,
            seckillGoods: seckillGoods,
            teamGoods: teamGoods,
            goodsGroups: goodsGroups,
            goodsDynamic: goodsDynamic,
            banners: banners,
            curPage: 0,
            gm: gm
        })

        // let curTime = Date.now()
        wx.setStorageSync('index_index_time', curTime)
    },
    getCache(key){
        let cache = this.data.cache
        // let nowTime = Date.now();
        // let saveTime = cache.saveTime || 0;
        // cache.saveTime = nowTime+1000*60;
        // if (saveTime > nowTime) {
            return cache[key]
        // }
    },
    setCache(key, val){
        // let nowTime = Date.now();
        let cache = this.data.cache
        // cache.saveTime = nowTime+1000*60;
        cache[key] = val
    },
    async getGoodsList(categoryId, append) {
        if (categoryId == 0) {
            categoryId = -1
        }
        let params = {
            categoryId: categoryId,
            nameLike: this.data.inputVal,
            page: this.data.curPage,
            pageSize: this.data.pageSize
        }
        // let reqKey = [params.categoryId, params.page, params.pageSize, params.nameLike].join(":")
        // console.log("getGoodsList=====>> 请求数据")
        // let res = this.getCache(reqKey)
        // if(res) {
        //     console.log("getGoodsList=====>> 无需下载")
        //     wx.hideLoading()
        // } else {
            // console.log("getGoodsList=====>> 下载")
            wx.showLoading({"mask": true})
            let res = await WXAPI.goods(params)
            // this.setCache(reqKey, res)
            wx.hideLoading()
        // }
        this.setData({loadingMoreHidden: res.code != 10000})
        var newDatas = []
        if (res.data && res.data.length>0) {
            for (var i = 0; i < res.data.length; i++) {
                let item = res.data[i]
                // if (item.mark == 0) {
                    newDatas.push(item)
                // }
            }
        }
        if (!append) {
            let newData = {
                goods:newDatas
            }
            this.setData(newData)
            return
        }
        if (newDatas.length == 0) {
            return
        }
        let goodsMap = {}
        let goods = this.data.goods
        for (let i = 0; i < goods.length; i++) {
            let item = goods[i]
            goodsMap[item.id] = item
        }
        for (let i = 0; i < newDatas.length; i++) {
            let item = newDatas[i]
            if (!goodsMap[item.id]) {
                goods.push(item);
            }
        }
        this.setData({goods:goods})
    },
    // getCoupons() {
    //     var that = this;
        // WXAPI.coupons().then(function(res) {
        //     if (res.code == 0) {
        //         that.setData({
        //             coupons: res.data
        //         });
        //     }
        // })
    // },
    goNotice(e) {
        const noticeId = e.currentTarget.dataset.id
        wx.navigateTo({
          url: '/pages/notice/show?id=' + noticeId,
        })
    },
    goCoupons(e) {
        wx.navigateTo({url: "/pages/coupons/index"})
    },
    bindinput(e) {
        this.setData({inputVal: e.detail.value})
    },
    bindconfirm(e) {
        let value = e.detail.value;
        if (!value || value.length == 0) {
            return;
        }
        this.setData({inputVal: value})
        wx.navigateTo({url: '/pages/goods/list?name=' + this.data.inputVal})
    }
})