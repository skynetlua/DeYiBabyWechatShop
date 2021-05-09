const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')
const wxpay = require('../../utils/pay.js')

Page({
    data: {
        wxlogin: true,

        goodsList: [],

        // hasNoCoupons: true,
        amountGoods: 0,
        amountLogistics: 0,
        amountCoupon: 0,
        amountReal: 0,

        // coupons: [],
        // curCoupon: null, // 当前选择使用的优惠券
        // curCouponShowText: '请选择使用优惠券', // 当前选择使用的优惠券

        remark: '',
        sendType: 0, // 配送方式 0,1 分别表示到店自取, 快递
        isInitData: false,
        steps0: [{
            text: '支付开团'
        },{
            text: '拼团中',
        },{
            text: '拼团成功',
        }],
        steps1: [{
            text: '支付拼团'
        },{
            text: '拼团成功',
        }]
    },
    onLoad(e) {
        if (e.id && e.num) {
            this.setData({
                goodsId: e.id,
                orderSkuId: e.skuid || 0,
                buyNumber: e.num || 1,
                teamBuy: e.team || 0,
                quickBuy: 1
            });
        } else {
            this.setData({
                teamBuy: 0,
                quickBuy: 0
            });
        }
    },
    onShow() {
        if (!this.isInitData) {
            AUTH.checkHasLogined().then(isLogined => {
                this.setData({
                    wxlogin: isLogined
                })
                if (isLogined) {
                    this.isInitData = true
                    this.prepareOrder()
                }
            })
        } else {
            if (!this.data.curAddressData) {
                this.initShippingAddress();
            }
        }
    },
    async prepareOrder() {
        let data = this.data;
        let params = {
            sendType: data.sendType,
            goodsId: data.goodsId,
            orderSkuId: data.orderSkuId,
            buyNumber: data.buyNumber,
            quickBuy: data.quickBuy,
            teamBuy: data.teamBuy
        }
        let res = await WXAPI.orderPrepare(params);
        if (res.code != 0) {
            wx.showToast({
                title: res.msg,
                icon: 'none'
            });
            setTimeout(() => {
                wx.navigateBack();
            }, 3000)
            return
        }
        let resData = res.data;
        resData.firstGoods = resData.goodsList[0] || {}
        this.setData(resData);
        this.initShippingAddress();
    },
    async initShippingAddress() {
        const res = await WXAPI.defaultAddress();
        if (res.code == 0) {
            this.setData({
                curAddressData: res.data
            });
        } else {
            wx.showToast({
                title: res.msg,
                icon: 'none'
            });
        }
    },
    remarkChange(e) {
        this.data.remark = e.detail.value;
    },
    // onClearCartChange(e) {
    //     let value = e.detail.value;
    //     console.log("onClearCartChange value =", value);
    // },
    // getGoodsInfos() {
    //     var goodsList = this.data._goodsList
    //     if (goodsList.length == 0) {
    //         return ""
    //     }
    //     let goodsInfos = []
    //     for (let i = 0; i < goodsList.length; i++) {
    //         let item = goodsList[i];
    //         let goodsInfo = {
    //             goodsId: item.goodsId,
    //             skuId: item.skuId,
    //             number: item.numberBuy,
    //         }
    //         goodsInfos.push(goodsInfo)
    //     }
    //     return JSON.stringify(goodsInfos)
    // },
    goCreateOrder() {
        var data = this.data
        if (data.orderId && data.amountReal) {
            this.processAfterCreateOrder(resData.orderId, resData.amountReal);
            return
        }
        let inviterId = wx.getStorageSync('referrer');
        let params = {
            sendType:   data.sendType,
            goodsId:    data.goodsId,
            orderSkuId: data.orderSkuId,
            buyNumber:  data.buyNumber,
            quickBuy:   data.quickBuy,
            teamBuy:    data.teamBuy,
            remark:     data.remark,
            inviterId:  inviterId ? inviterId:0
        }
        if (params.sendType > 0) {
            var address = data.curAddressData
            if (!address) {
                wx.showToast({title: '请设置收货地址', icon: 'none'})
                return;
            }
            params.addressId = address.id;
        }
        // if (data.curCoupon) {
        //     postData.couponId = data.curCoupon.id;
        // }
        var that = this;
        WXAPI.orderCreate(params).then(function(res) {
            if (res.code != 0) {
                wx.showModal({
                    title: '错误',
                    content: res.msg,
                    showCancel: false
                })
                return;
            }
            let resData = res.data;
            resData.firstGoods = resData.goodsList[0] || {}
            that.setData(resData);
            that.processAfterCreateOrder(resData.orderId, resData.amountReal)
        })
    },
    async processAfterCreateOrder(orderId, amountReal) {
        wxpay.requestPay(orderId, amountReal, function(){
            wx.redirectTo({url: "/pages/order-list/index"});
            // wx.switchTab({url: "/pages/order-list/index"})
        })
    },
    addAddress() {
        wx.navigateTo({
            url: "/pages/address-add/index"
        })
    },
    selectAddress() {
        wx.navigateTo({
            url: "/pages/select-address/index"
        })
    },
    // async getMyCoupons() {
    //     const res = await WXAPI.myCoupons({status: 0})
    //     if (res.code == 0) {
    //         var coupons = res.data.filter(entity => {
    //             return entity.moneyHreshold <= this.data.allGoodsAndYunPrice;
    //         })
    //         if (coupons.length > 0) {
    //             coupons.forEach(ele => {
    //                 ele.nameExt = ele.name + ' [满' + ele.moneyHreshold + '元可减' + ele.money + '元]'
    //             })
    //             this.setData({
    //                 hasNoCoupons: false,
    //                 coupons: coupons
    //             });
    //         }
    //     }
    // },
    // bindChangeCoupon(e) {
    //     const selIndex = e.detail.value;
    //     let coupon = this.data.coupons[selIndex];
    //     this.setData({
    //         curCoupon: coupon,
    //         amountCoupon: coupon.money,
    //         curCouponShowText: coupon.nameExt
    //     });
    // },
    radioChange(e) {
        let sendType = parseInt(e.detail.value);
        this.setData({
            sendType: sendType
        })
        // this.processYunfei()
    },
    cancelLogin() {
        wx.navigateBack()
    },
    processLogin(e) {
        if (!e.detail.userInfo) {
            wx.showToast({
                title: '已取消',
                icon: 'none',
            })
            return;
        }
        AUTH.register(this);
    },
})