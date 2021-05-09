const wxpay = require('../../utils/pay.js')
const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')

const StatusType = {
    StatusPay          :0,
    StatusSend         :1,
    StatusReceive      :2,
    StatusRepute       :3,
    StatusFinish       :4,
    StatusRefundApply  :5,
    StatusRefundRefuse :6,
    StatusClose        :-1,
    StatusRefundFinish :-2,
}

var sliderWidth = 96;

Page({
    data: {
        tabs: ["待支付", "待发货", "待收货"],
        tabValues: [StatusType.StatusPay, StatusType.StatusSend, StatusType.StatusReceive],
        status: -1,
        activeIndex: 0,
        sliderOffset: 0,
        sliderLeft: 0,
        typeOrders: {},
        isAddCoupon: true,
        typeGoodsMap: {}
    },
    setStatus: function(status) {
        const orderList = this.data.typeOrders[status] || null;
        const goodsMap = this.data.typeGoodsMap[status] || null;
        this.setData({
            status :status,
            goodsMap :goodsMap,
            orderList :orderList,
        });
        if (!orderList) {
            this.requestData(status);
        }
    },

    onLoad: function(options) {
        const that = this;
        wx.getSystemInfo({
            success: function(res) {
                that.setData({
                    sliderLeft: (res.windowWidth / that.data.tabs.length - sliderWidth) / 2,
                    sliderOffset: res.windowWidth / that.data.tabs.length * that.data.activeIndex
                });
            }
        });
    },

    tabClick: function(e) {
        this.setData({
            sliderOffset: e.currentTarget.offsetLeft,
            activeIndex: e.currentTarget.id
        });
        let status = this.data.tabValues[e.currentTarget.id];
        this.setStatus(status);
    },

    requestData: function(status) {
        let that = this;
        WXAPI.gmOrderList(status).then(function(res) {
            if (res.code == 0) {
                let orderList = res.data.orderList || null;
                let goodsMap = res.data.goodsMap || {};
                that.data.typeOrders[status] = orderList;
                that.data.typeGoodsMap[status] = goodsMap;
                if (orderList) {
                    orderList.sort(function(a, b) {
                        return b.timeStamp-a.timeStamp;
                    })
                }
                that.setData({
                    status: status,
                    goodsMap: goodsMap,
                    orderList: orderList,
                });
            } else {
                wx.showToast({title: res.msg, icon: 'none'});
                that.setData({
                    goodsMap:null,
                    orderList: null,
                });
            }
        })
    },
    closePopup:function() {
        this.data.curOrder.amountCoupon = this.amountCoupon
        this.setData({
            curOrder: null
        });
    },
    couponOrderTap:function(e) {
        const orderId = e.currentTarget.dataset.id
        let selectOrder;
        for (var i = 0; i < this.data.orderList.length; i++) {
            let order = this.data.orderList[i]
            if (order.orderId == orderId) {
                selectOrder = order
                break
            }
        }
        if (!selectOrder) {return}

        this.amountCoupon = selectOrder.amountCoupon
        this.setData({
            curOrder: selectOrder
        });
    },
    numMinusTap(e) {
        let curOrder = this.data.curOrder;
        const field = e.currentTarget.dataset.field;
        let fieldValue = curOrder[field];
        if (fieldValue <= 0) {
            fieldValue = 0;
        }else{
            fieldValue--;
        }
        curOrder[field] = fieldValue;
        this.setData({
            curOrder: curOrder,
        });
    },

    numPlusTap(e) {
        let curOrder = this.data.curOrder;
        const field = e.currentTarget.dataset.field;
        let fieldValue = curOrder[field];
        fieldValue++;
        curOrder[field] = fieldValue;
        this.setData({
            curOrder: curOrder,
        });
    },

    watchInput(e) {
        let curOrder = this.data.curOrder;
        const field = e.currentTarget.dataset.field;
        let fieldValue = Number(e.detail.value);
        if (!fieldValue && fieldValue != 0) {
            fieldValue = curOrder[field];
        }
        curOrder[field] = fieldValue;
        this.setData({
            curOrder: curOrder,
        });
    },
    requestCoupon:function() {
        if (!this.data.curOrder) {
            return
        }
        let curOrder = this.data.curOrder
        let amountCoupon = curOrder.amountCoupon
        curOrder.amountCoupon = this.amountCoupon
        if (amountCoupon > curOrder.amountGoods) {
            amountCoupon = curOrder.amountGoods
        }
        this.setData({
            curOrder: null
        });
        let that = this
        wx.showModal({
            title: '确定要给订单'+curOrder.orderNumber+'使用['+(amountCoupon/100)+'元]优惠券',
            content: '',
            success: function(res) {
                if (res.confirm) {
                    WXAPI.gmOrderCoupon(curOrder.orderId, curOrder.playerId, amountCoupon, curOrder.couponId).then(function(res) {
                        if (res.code == 0) {
                            that.onShow(true);
                        } else {
                            wx.showToast({title: res.msg, icon: 'none'});
                        }
                    });
                }
            }
        });
    },
    refundConfirmTap: function(e) {
        let that = this
        const orderId = e.currentTarget.dataset.id
        const playerId = e.currentTarget.dataset.playerid
        wx.showModal({
            title: '请确保买家已退货，再操作退款',
            content: '',
            success: function(res) {
                if (res.confirm) {
                    WXAPI.gmRefundConfirm(orderId, playerId).then(function(res) {
                        if (res.code == 0) {
                            that.onShow(true);
                        } else {
                            wx.showToast({title: res.msg, icon: 'none'});
                        }
                    })
                }
            }
        });
    },
    refundCancelTap: function(e) {
        let that = this
        const orderId = e.currentTarget.dataset.id
        const playerId = e.currentTarget.dataset.playerid
        wx.showModal({
            title: '请联系买家，是否要确定取消退款？',
            content: '',
            success: function(res) {
                if (res.confirm) {
                    WXAPI.gmRefundCancel(orderId, playerId).then(function(res) {
                        if (res.code == 0) {
                            that.onShow(true);
                        } else {
                            wx.showToast({title: res.msg, icon: 'none'});
                        }
                    })
                }
            }
        });
    },
    sendOrderTap: function(e) {
        const that = this;
        const orderId = e.currentTarget.dataset.id;
        const playerId = e.currentTarget.dataset.playerid
        wx.showModal({
            title: '请确保发货',
            content: '',
            success: function(res) {
                if (res.confirm) {
                    let status = StatusType.StatusReceive;
                    WXAPI.gmOrderDo(orderId, status, playerId).then(function(res) {
                        if (res.code == 0) {
                            that.onShow(true);
                        } else {
                            wx.showToast({title: res.msg, icon: 'none'});
                        }
                    })
                }
            }
        });
    },
    
    onReady: function() {
    },
    onShow: function(isForce) {
        let statusType = this.data.status
        if (!isForce && statusType != -1) {
            return
        }
        if (statusType == -1) {
            statusType = StatusType.StatusPay
        }
        this.data.typeOrders = {};
        this.data.typeGoodsMap = {};
        AUTH.checkHasLogined().then(isLogined => {
            if (isLogined) {
                this.requestData(statusType);
            } else {
                wx.showModal({
                    title: '提示',
                    content: '本次操作需要您的登录授权',
                    cancelText: '暂不登录',
                    confirmText: '前往登录',
                    success(res) {
                        if (res.confirm) {
                            wx.switchTab({
                                url: "/pages/my/index"
                            })
                        } else {
                            wx.navigateBack()
                        }
                    }
                })
            }
        })
    },
    onHide: function() {
        // 生命周期函数--监听页面隐藏
    },
    onUnload: function() {
        // 生命周期函数--监听页面卸载
    },
    onPullDownRefresh: function() {
        // 页面相关事件处理函数--监听用户下拉动作
    },
    onReachBottom: function() {
        // 页面上拉触底事件的处理函数
    }
})