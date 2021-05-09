const WXAPI = require('apifm-wxapi')
const CONFIG = require('../../config.js')
const AUTH = require('../../utils/auth')
const TOOLS = require('../../utils/tools.js')

import Poster from 'wxa-plugin-canvas/poster/poster'

Page({
    data: {
        wxlogin: true,

        cartNum: 0,
        hideShopPopup: true,

        buyNumber: 0,
        buyNumMin: 1,
        buyNumMax: 0,
        goods:{sellPrice:0, originPrice:0},
        selectSkuIds: [],
        teamBuy: 0,
        shopType: "addShopCar", //购物类型，加入购物车或立即购买，默认为加入购物车,
        countDownTime: 0,
        countDownData: {},
        teamBuyList: []
    },
    onCountDownChange(e) {
        let data = e.detail;
        if (data.hours == 0 && data.minutes == 0 && data.seconds == 0) {
            console.log("onCountDownChange===>>");
            this.onShow();
        }
        if (data.hours < 10) {
            data.hours = '0'+data.hours;
        }
        if (data.minutes < 10) {
            data.minutes = '0'+data.minutes;
        }
        if (data.seconds < 10) {
            data.seconds = '0'+data.seconds;
        }
        this.setData({
            countDownData: data
        });
    },
    async onLoad(e) {
        // console.log("goods-details e =", e)
        this.setData({
            goodsId: e.id,
            mode: e.mode?e.mode:null
        })
    },
    callPhone() {
        wx.setClipboardData({
            data: '18926627466',
            success: function (res) {
                wx.showModal({
                    title: '联系客服',
                    content: '微信号已复制成功，是否拨打电话给客服？',
                    confirmText: "拨打电话",
                    cancelText: "取消",
                    success: function(res) {
                        if (res.confirm) {
                            wx.makePhoneCall({
                                phoneNumber:'18926627466',
                                success:function() {
                                    wx.showToast({title: '拨打电话成功！',icon: 'none'})
                                },
                                fail:function() {
                                    wx.showToast({title: '拨打电话失败！',icon: 'none'})
                                }
                            })
                        }
                    }
                });
            },
            fail:function() {
                wx.showToast({title: '微信号复制失败！',icon: 'none'})
            }
        })
    },
    onShow() {
        AUTH.checkHasLogined().then(isLogined => {
            if (isLogined) {
                this.requestGoodsDetail();
            }
            this.setData({
                wxlogin: isLogined
            })
        })
    },
    onHide() {
        this.stopTimer();
    },
    startTimer() {
        this.stopTimer();
        let that = this;
        this.timerHandle = setInterval(function() {
            that.updateTimer();
        }, 300);
    },
    stopTimer() {
        if (this.timerHandle) {
            clearInterval(this.timerHandle);
            this.timerHandle = null;
        }
    },
    updateTimer() {
        let teamBuyList = this.data.teamBuyList;
        for (var i = 0; i < teamBuyList.length; i++) {
            let order = teamBuyList[i];
            if (!order.EndTime) {
                continue;
            }
            let countTime = order.EndTime*1000-Date.now();
            order.countDownTime = countTime;
            if (countTime <= 0) {
                continue;
            }
            countTime = Math.floor(countTime/100);
            order.microseconds = countTime%10;

            countTime = Math.floor(countTime/10);
            order.seconds = countTime%60;

            countTime = Math.floor(countTime/60);
            order.minutes = countTime%60;

            countTime = Math.floor(countTime/60);
            order.hours = countTime;
            
            if (order.hours < 10) {
                order.hours = '0'+order.hours;
            }
            if (order.minutes < 10) {
                order.minutes = '0'+order.minutes;
            }
            if (order.seconds < 10) {
                order.seconds = '0'+order.seconds;
            }
        }
        this.setData({
            teamBuyList: teamBuyList
        });
    },
    processTeamOrders(teamBuyList) {
        teamBuyList = teamBuyList.reverse()
        let teamOrderMap = {};
        for (var i = 0; i < teamBuyList.length; i++) {
            let teamOrder = teamBuyList[i];
            teamOrderMap[teamOrder.OrderId] = teamOrder;
        }
        let teamingOrders = [];
        let teamedOrders = [];
        for (var i = 0; i < teamBuyList.length; i++) {
            let teamOrder = teamBuyList[i];
            let pTeamOrder = teamOrderMap[teamOrder.TeamBuy]
            if (pTeamOrder && teamOrder.TeamBuy != 2) {
                if (!pTeamOrder.childs) {
                    pTeamOrder.childs = [];
                }
                pTeamOrder.childs.push(teamOrder);
            } else {
                if (teamOrder.Status == 1) {
                    teamingOrders.push(teamOrder);
                } else {
                    teamedOrders.push(teamOrder);
                }
            }
        }
        this.data.teamingCount = teamingOrders.length;
        if (teamingOrders.length < 6) {
            for (var i = 0; i < teamedOrders.length; i++) {
                teamingOrders.push(teamedOrders[i]);
                if (teamingOrders.length >= 6) {
                    break;
                }
            }
        }
        return teamingOrders;
    },
    requestGoodsDetail() {
        let that = this;
        let goodsId = this.data.goodsId;
        WXAPI.pageGoodsDetail(goodsId).then(function(res) {
            if (res.code == 0) {
                var data = res.data
                var goods = data.goods
                // if (data.shopId) {
                //     that.shopSubdetail(data.shopId)
                // }
                if (goods.videoId) {
                    that.getVideoSrc(goods.videoId);
                }
                goods.isClose = false
                if (goods.mark == 2) {
                    var curTime = new Date().getTime();
                    curTime = Math.floor(curTime/1000)
                    var endInterval = goods.endTime-curTime;
                    if (endInterval <= 0) {
                        goods.isClose = true;
                    }
                    var startInterval = goods.startTime-curTime;
                    if (startInterval > 0) {
                        goods.isClose = true;
                    }
                    if (startInterval <= 0 && endInterval > 0) {
                        that.setData({
                            countDownTime:endInterval*1000
                        })
                    }
                }
                var skuGroups = [];
                var skuList = [];
                if (goods.skuJson && goods.skuJson.length > 0) {
                    skuList = JSON.parse(goods.skuJson);
                    // skuList.push({"label":"规格2","id":100,"name":"2NB码96片","price":0});
                    // skuList.push({"label":"规格2","id":200,"name":"2NB码96片","price":0});
                    // skuList.push({"label":"规格3","id":10000,"name":"3NB码96片","price":0});
                    // skuList.push({"label":"规格3","id":20000,"name":"3NB码96片","price":0});
                    let skuGroupMap = {};
                    for (var i = 0; i < skuList.length; i++) {
                        let sku = skuList[i];
                        let skuGroup = skuGroupMap[sku.label];
                        let level = TOOLS.getLevel(sku.id)
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
                                console.log("shippingCartInfo skuGroup.level =", skuGroup.level, "level =", level);
                            }
                        }
                        skuGroup.skuList.push(sku)
                    }
                }
                // for (var i = 0; i < data.teamBuyList.length; i++) {
                //     let teamBuy = data.teamBuyList[i];
                //     that.parseOrder(teamBuy);
                // }
                let teamBuyList = data.teamBuyList;
                if (teamBuyList && teamBuyList.length > 0) {
                    teamBuyList = that.processTeamOrders(teamBuyList);
                    that.startTimer();
                }
                that.setData({
                    goods: goods,
                    skuList: skuList,
                    skuGroups: skuGroups,
                    buyNumMax: goods.numberStore,
                    buyNumber: goods.numberStore > 0 ? 1 : 0,
                    faved: data.isFavorite == 1,
                    cartNum: data.cartNum,
                    buyCount: data.buyCount,
                    buyLimit: data.buyLimit,
                    playerId: data.playerId,
                    teamBuyList: teamBuyList,
                    teamingCount: that.data.teamingCount || 0
                });
            }else{
                wx.showToast({title: res.msg, icon: 'none'});
            }
        });
    },
    reputation: function(goodsId) {
        var that = this;
        WXAPI.goodsReputation({
            goodsId: goodsId
        }).then(function(res) {
            if (res.code == 0) {
                that.setData({
                    reputations: res.data
                });
            }
        })
    },
    async addFav() {
        AUTH.checkHasLogined().then(isLogined => {
            this.setData({
                wxlogin: isLogined
            })
            if (isLogined) {
                let goodsId = this.data.goodsId;
                if (this.data.faved) {
                    WXAPI.goodsFavDelete(goodsId).then(res => {
                        if (res.code == 0) {
                            this.setData({faved: false})
                        }
                    })
                } else {
                    WXAPI.goodsFavPut(goodsId).then(res => {
                        if (res.code == 0) {
                            this.setData({faved: true})
                        }
                    })
                }
            }
        })
    },
    async shopSubdetail(shopId) {
        const res = await WXAPI.shopSubdetail(shopId)
        if (res.code == 0) {
            this.setData({
                shopInfo: res.data
            })
        }
    },
    goShopCar: function() {
        wx.reLaunch({
            url: "/pages/shop-cart/index"
        });
    },
    isCanBuy() {
        if (this.data.buyLimit && this.data.buyLimit > 0) {
            if (this.data.buyCount && this.data.buyCount > 0) {
                if (this.data.buyCount >= this.data.buyLimit) {
                    wx.showToast({title: '超过限购，限购'+this.data.buyLimit+'次数',icon: 'none'})
                    return false
                }
            }
        }
        return true
    },
    onAddShopCartClick: function() {
        if (!this.isCanBuy()) {
            return
        }
        this.setData({
            shopType: "addShopCar"
        })
        this.showSkuPanel();
    },
    onToBuyClick: function() {
        if (!this.isCanBuy()) {
            return
        }
        this.setData({
            shopType: "tobuy"
        });
        this.showSkuPanel();
    },
    onSingleBuy: function() {
        if (!this.isCanBuy()) {
            return
        }
        this.setData({
            teamBuy: 1,
            shopType: "tobuy"
        });
        this.showSkuPanel();
    },
    onTeamBuy: function(e) {
        const orderId = e.currentTarget.dataset.orderid;
        let teamBuy = 2;
        if (orderId) {
            teamBuy = Number(orderId);
            for (var i = 0; i < this.data.teamBuyList.length; i++) {
                let teamBuyOrder = this.data.teamBuyList[i];
                if (teamBuyOrder.OrderId == orderId) {
                    if (teamBuyOrder.PlayerId == this.data.playerId) {
                        wx.showToast({
                            title: '不能和自己拼单',
                            icon: 'none'
                        })
                        return;
                    }
                    break;
                }
            }
        }
        if (!this.isCanBuy()) {
            return
        }
        this.setData({
            teamBuy: teamBuy,
            shopType: "tobuy"
        });
        this.showSkuPanel();
    },
    showSkuPanel: function() {
        this.updateSkuInfos()
        this.setData({
            hideShopPopup: false
        })
    },
    onCloseSkuPanel: function() {
        this.setData({
            hideShopPopup: true
        })
    },
    updateSkuInfos: function() {
        let skuGroups = this.data.skuGroups;
        let selectSkuIds = this.data.selectSkuIds;
        const goods = this.data.goods;
        let selectSkuPrice = goods.sellPrice;
        let selectSkuIcon = goods.icon;
        let teamBuy = this.data.teamBuy;
        if (teamBuy && teamBuy > 0) {
            if (teamBuy == 1) {
                selectSkuPrice = goods.originPrice;
            } else {
                selectSkuPrice = goods.sellPrice;
            }
        }
        for (var i = 0; i < skuGroups.length; i++) {
            let selectSkuId = selectSkuIds[i];
            let skuGroup = skuGroups[i];
            skuGroup.skuList.forEach(sku => {
                if (sku.id == selectSkuId) {
                    sku.active = true;
                    if (selectSkuId < 100) {
                        if (goods.skuPics && goods.skuPics[selectSkuId]) {
                            selectSkuIcon = goods.skuPics[selectSkuId]
                        }
                    }
                    if (sku.price && sku.price > 0) {
                        selectSkuPrice = sku.price
                    }
                } else {
                    sku.active = false;
                }
            })
        }
        this.setData({
            skuGroups: skuGroups,
            selectSkuIcon: selectSkuIcon,
            selectSkuPrice: selectSkuPrice,
        });
    },
    numJianTap: function() {
        if (this.data.buyNumber > this.data.buyNumMin) {
            var currentNum = this.data.buyNumber;
            currentNum--;
            this.setData({
                buyNumber: currentNum
            })
        }else{
            wx.showToast({
                title: '不能小于1',
                icon: 'none'
            })
        }
    },
    numJiaTap: function() {
        if (this.data.buyNumber < this.data.buyNumMax) {
            var currentNum = this.data.buyNumber;
            currentNum++;
            this.setData({
                buyNumber: currentNum
            })
        }else{
            wx.showToast({
                title: '库存有限',
                icon: 'none'
            })
        }
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
    makeOrderSkuId() {
        if (this.data.buyNumber < 1) {
            wx.showModal({
                title: '提示',
                content: '亲，购买数量不能为0喔',
                showCancel: false
            })
            return null;
        }
        let skuGroups = this.data.skuGroups;
        let count = skuGroups.length;
        if (count == 0) {
            return 0;
        }
        let selectSkuIds = this.data.selectSkuIds;
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
                this.showSkuPanel();
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
    },
    /**
     * 加入购物车
     */
    async addShopCar() {
        const isLogined = await AUTH.checkHasLogined();
        if (!isLogined) {
            this.setData({wxlogin: false});
            return;
        }
        let orderSkuId = this.makeOrderSkuId()
        if (orderSkuId == null) {
            return;
        }
        let data = this.data;
        let goodsId = data.goods.id;
        let buyNumber = data.buyNumber;
        // const selectSkuId = this.data.selectSkuId || 0;
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
        this.requestGoodsDetail();
    },
    /**
     * 立即购买
     */
    async buyNow() {
        const isLogined = await AUTH.checkHasLogined()
        if (!isLogined) {
            this.setData({wxlogin: false})
            return
        }
        let orderSkuId = this.makeOrderSkuId()
        if (orderSkuId == null) {
            return;
        }
        const data = this.data;
        this.onCloseSkuPanel();
        let teamBuy = this.data.teamBuy;
        wx.navigateTo({
            url: "/pages/to-pay-order/index?id="+data.goods.id+"&skuid="+orderSkuId+"&num="+data.buyNumber+"&team="+teamBuy
        })
    },
    onShareAppMessage() {
        let goods = this.data.goods;
        let uid = wx.getStorageSync('uid');
        let path = '/pages/start/loading?inviter_id=' + uid + '&route=/pages/goods-details/index%3fid%3d' + goods.id;
        return {
            title: goods.name,
            path: path
        };
    },
    getVideoSrc: function(videoId) {
        var that = this;
        WXAPI.videoDetail(videoId).then(function(res) {
            if (res.code == 0) {
                that.setData({
                    videoMp4Src: res.data.fdMp4
                });
            }
        })
    },
    goIndex() {
        wx.switchTab({
            url: '/pages/index/index',
        });
    },
    cancelLogin() {
        this.setData({
            wxlogin: true
        })
    },
    processLogin(e) {
        if (!e.detail.userInfo) {
            wx.showToast({
                title: '登陆失败',
                icon: 'none',
            })
            return;
        }
        AUTH.register(this);
    },
    previewImage(e) {
        const url = e.currentTarget.dataset.url;
        if(typeof url == "string"){
            wx.previewImage({current: url, urls: [url]});
        }else{
            wx.previewImage({current: url[0], urls: url});
        }
    },
    // viewImage(e) {
    //     const index = e.currentTarget.dataset.index;
    //     const goods = this.data.goods;
    //     wx.previewImage({current: goods.contents[index], urls: goods.contents});
    // },
    async drawSharePic() {
        const goods = this.data.goods;
        var params = {
            scene: 'qg,id=' + goods.id + ',' + wx.getStorageSync('uid'),
            page: 'pages/start/loading',
            is_hyaline: true,
            autoColor: true,
            expireHours: 1
        };
        const qrcodeRes = await WXAPI.wxaQrcode(params);
        if (qrcodeRes.code != 0) {
            wx.showToast({
                title: qrcodeRes.msg,
                icon: 'none'
            });
            return;
        }
        const qrcode = qrcodeRes.data;
        const _this = this;
        wx.getImageInfo({
            src: goods.icon,
            success(res) {
                const height = 490 * res.height / res.width;
                _this.drawSharePicDone(height, qrcode);
            },
            fail(e) {
                console.error(e);
            }
        });
    },
    drawSharePicDone(picHeight, qrcode) {
        var goods = this.data.goods;
        const _baseHeight = 74 + (picHeight + 120);
        let posterConfig = {
            width: 750,
            height: picHeight + 660,
            backgroundColor: '#fff',
            debug: false,
            blocks: [
                {x: 76, y: 74, width: 604, height: picHeight + 120, borderWidth: 2, borderColor: '#c2aa85', borderRadius: 8}
            ],
            images: [
                {x: 133, y : 133, width: 490, height: picHeight, url: goods.icon},
                {x: 76, y: _baseHeight + 199, width: 222, height: 222, url: qrcode}
            ],
            texts: [
                {x: 375, y: _baseHeight + 80, width: 650, lineNum: 2, text: goods.name, textAlign: 'center', fontSize: 40, color: '#333'},
                {x: 375, y: _baseHeight + 180, text: '￥' + goods.sellPrice/100, textAlign: 'center', fontSize: 50, color: '#e64340'},
                {x: 352, y: _baseHeight + 320, text: '长按识别小程序码', fontSize: 28, color: '#999'}
            ],
        };
        this.setData({
            posterConfig:posterConfig
        }, () => {
            Poster.create();
        });
    },
    onPosterSuccess(e) {
        console.log('success:', e);
        this.setData({
            posterImg: e.detail,
            showposterImg: true
        })
    },
    onPosterFail(e) {
        console.error('fail:', e);
    },
    closePoster() {
        this.setData({
            showposterImg: false
        });
    },
    savePosterPic() {
        const _this = this;
        wx.saveImageToPhotosAlbum({
            filePath: this.data.posterImg,
            success: (res) => {
                wx.showModal({
                    content: '已保存到手机相册',
                    showCancel: false,
                    confirmText: '知道了',
                    confirmColor: '#333'
                })
            },
            complete: () => {
                _this.setData({showposterImg: false});
            },
            fail: (res) => {
                wx.showToast({
                    title: res.errMsg,
                    icon: 'none',
                    duration: 2000
                })
            }
        })
    },
})
