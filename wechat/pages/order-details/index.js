const CONFIG = require('../../config.js')
const WXAPI = require('apifm-wxapi')
const wxpay = require('../../utils/pay.js')
import wxbarcode from 'wxbarcode'

Page({
    data: {
        orderId: 0,
        goodsList: [],
        orderInfo: {
            status: -11
        },
        logistics: null,
        goodsCoupons: null,
        logisticsTraces: null,

        active: 0,
        steps: [{
            text: '支付开团'
        },{
            text: '拼团中',
        },{
            text: '拼团成功',
        }]
    },
    onLoad(e) {
        this.setData({
            orderId: e.id
        });
        if (e.playerid) {
            this.setData({
                playerid: e.playerid
            });
        }
        if (e.mode) {
            this.setData({
                mode: e.mode
            });
        }
    },
    onShareAppMessage() {
        let playerId = this.data.orderInfo.playerId;
        let orderId = this.data.orderInfo.orderId;
        let uid = wx.getStorageSync('uid');
        let path = '/pages/start/loading?inviter_id=' + uid + '&route=/pages/order-details/index%3fid%3d' + orderId+'%26playerid%3d'+playerId;
        return {
            title: '订单号：'+this.data.orderInfo.orderNumber,
            path: path
        };
    },
    onShow() {
        var that = this;
        var gm = wx.getStorageSync('gm')
        // this.setData({
        //     mode: gm || 0
        // });
        if (gm == 1 && this.data.playerid) {
            WXAPI.gmOrderDetail(this.data.orderId, this.data.playerid).then(function(res) {
                if (res.code != 0) {
                    wx.showModal({
                        title: '错误',
                        content: res.msg,
                        showCancel: false
                    })
                    return;
                }
                let retData = res.data;
                let orderInfo = retData.orderInfo || {};
                // 绘制核销码
                // if (orderInfo.hxNumber && orderInfo.status > 0) {
                //     wxbarcode.qrcode('qrcode', orderInfo.hxNumber, 650, 650);
                // }
                that.setData({
                    orderInfo: orderInfo,
                    goodsList: retData.goodsList,
                    firstGoods: retData.goodsList[0] || {},
                    logistics: retData.logistics,
                    // goodsCoupons: retData.goodsCoupons,
                    // logisticsTraces: retData.logisticsTraces,
                });
            })
            return
        }

        WXAPI.orderDetail(this.data.orderId).then(function(res) {
            if (res.code != 0) {
                wx.showModal({
                    title: '错误',
                    content: res.msg,
                    showCancel: false
                })
                return;
            }
            let retData = res.data;
            let orderInfo = retData.orderInfo || {};
            // 绘制核销码
            // if (orderInfo.hxNumber && orderInfo.status > 0) {
            //     wxbarcode.qrcode('qrcode', orderInfo.hxNumber, 650, 650);
            // }
            that.setData({
                orderInfo: orderInfo,
                goodsList: retData.goodsList,
                logistics: retData.logistics,
                // goodsCoupons: retData.goodsCoupons,
                // logisticsTraces: retData.logisticsTraces,
            });
        })
    },
    toPayTap(e) {
        // 防止连续点击--开始
        if (this.data.payButtonClicked) {
            wx.showToast({
                title: '休息一下~',
                icon: 'none'
            })
            return
        }
        this.data.payButtonClicked = true
        setTimeout(() => {
            this.data.payButtonClicked = false
        }, 2000)

        const that = this;
        let orderInfo = this.data.orderInfo;
        wxpay.requestPay(orderInfo.orderId, orderInfo.amountReal, function() {
            that.onShow();
        })
    },
    toCancelTap() {
        let that = this;
        let orderId = this.data.orderId;
        wx.showModal({
            title: '确定要撤销订单？',
            content: '',
            success(res) {
                if (res.confirm) {
                    WXAPI.orderClose(orderId).then(function(res) {
                        if (res.code == 0) {
                            that.onShow()
                        }
                    })
                }
            }
        })
    },
    toHideTap() {
        let that = this;
        let orderId = this.data.orderId;
        wx.showModal({
            title: '确定要删除订单？',
            content: '',
            success(res) {
                if (res.confirm) {
                    WXAPI.orderClose(orderId).then(function(res) {
                        if (res.code == 0) {
                            wx.navigateBack({})
                        }
                    })
                }
            }
        })
    },
    toDeleteTap() {
        let that = this;
        let orderId = this.data.orderId;
        wx.showModal({
            title: '确定要删除订单？',
            content: '',
            success(res) {
                if (res.confirm) {
                    WXAPI.orderDelete(orderId).then(function(res) {
                        if (res.code == 0) {
                            wx.navigateBack({})
                        }
                    })
                }
            }
        })
    },
    wuliuDetailsTap(e) {
        wx.navigateTo({
            url: "/pages/wuliu/index?id=" + orderId
        })
    },
    confirmBtnTap(e) {
        let that = this;
        let orderId = this.data.orderId;
        wx.showModal({
            title: '确认您已收到商品？',
            content: '',
            success(res) {
                if (res.confirm) {
                    WXAPI.orderDelivery(orderId).then(function(res) {
                        if (res.code == 0) {
                            that.onShow();
                        }
                    })
                }
            }
        })
    },
    submitReputation(e) {
        let reputations = [];
        let value = e.detail.value;
        for (var i = 0; i < 100; i++) {
            let key = "goodsId" + i;
            let goodsId = value[key];
            if (goodsId) {
                key = "repute" + i
                let repute = value[key];
                key = "remark" + i
                let remark = value[key];
                let reputation = {
                    goodsId: parseInt(goodsId),
                    repute: parseInt(repute),
                    remark: remark,
                };
                reputations.push(reputation);
            }else{
                break;
            }
        }
        let that = this;
        WXAPI.orderReputation({
            orderId: this.data.orderId,
            reputations: JSON.stringify(reputations)
        }).then(function(res) {
            if (res.code == 0) {
                that.onShow();
                wx.showToast({
                    title: '评论成功~',
                    icon: 'none'
                })
            }
        });
    },
    refundApply(e) {
        let orderInfo = this.data.orderInfo;
        wx.navigateTo({
            url: "/pages/order/refundApply?id=" + orderInfo.orderId + "&amount=" + orderInfo.amountReal
        })
    }

})