const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')
const TOOLS = require('../../utils/tools.js')

Page({
    data: {
        listType: 1, // 1为1个商品一行，2为2个商品一行    
        name: '', // 搜索关键词
        orderBy: '', // 排序规则
        loadingMoreHidden: true,

        showSkuPanel: false,
        selectSkuIds: []
    },
    onLoad(options) {
        let listType = wx.getStorageSync("listType")
        if (!listType) {
            listType = 2
        }
        this.setData({
			name: options.name,
            listType: listType,
            // categoryId: options.categoryId
        })
        if (options.id) {
            this.setData({
                categoryId: options.id
            })
        }
        this.search()
    },
    onReady() {
    },
    onShow() {
    },
    async search() {
        wx.showLoading({title: '加载中'});
        let data = {
            page: 0,
            pageSize: 1000,
            orderBy: this.data.orderBy,
        }
        if (this.data.name) {
            data.nameLike = this.data.name
        }
		// if (this.data.categoryId) {
  //           data.categoryId = this.data.categoryId;
  //       }
        const res = await WXAPI.goods(data)
        wx.hideLoading();
        this.setData({loadingMoreHidden: res.code != 10000})
        if (res.code == 0 || res.code == 10000) {
            this.setData({
                goodsList: res.data,
            });
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
    // bindinput(e) {
    //     this.setData({
    //         name: e.detail.value
    //     })
    // },
    // bindconfirm(e) {
    //     this.setData({
    //         name: e.detail.value
    //     })
    //     this.search();
    // },
    onSearch(e) {
        let value = e.detail;
        if (!value || value.length == 0) {
            return;
        }
        this.setData({
            name: value
        })
        this.search();
    },
    filter(e) {
        this.setData({
            orderBy: e.currentTarget.dataset.val
        });
        this.updateDatas(this.data.goodsList);
        // this.search()
    },
    updateDatas(items){
        if (!items) {
            return
        }
        var orderBy = this.data.orderBy
        if (orderBy) {
            items.sort(function(a, b) {
                if (orderBy == "addedDown") {
                    return a.updateTime-b.updateTime
                } else if (orderBy == "priceUp") {
                    return a.sellPrice-b.sellPrice
                } else {
                    return a.numberOrder-b.numberOrder
                }
            })
        }
        this.setData({goodsList: items})
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