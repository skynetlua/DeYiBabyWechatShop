const wxpay = require('../../utils/pay.js')
const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')

Page({
    data: {
        menuItems: [
            {menuId: 8888, label: '全部'},
            {menuId: 0, label: '待支付'},
            {menuId: 1, label: '待发货'},
            {menuId: 2, label: '待收货'},
            {menuId: 3, label: '已完成'},
            {menuId: 4, label: '已取消'},
        ],
        curMenuId: 8888,
        hasRefund: false,
        menuOrders:{},
        badges: [0, 0, 0, 0, 0],
        menuIds: [8888, 0, 1, 2, 3]
    },
    menuTap(e) {
        const curMenuId = e.currentTarget.dataset.menuid;
        this.setMenuId(curMenuId);
    },
    setMenuId(curMenuId) {
        const menuIds = this.data.menuIds;
        const badges = this.data.badges;
        menuIds.some((id, index) => {
            if (id == curMenuId) {
                badges[index] = 0;
                return true;
            }
        });
        const orderList = this.data.menuOrders[curMenuId] || [];
        this.setData({
            badges: badges,
            curMenuId: curMenuId,
            orderList: orderList
        });
    },
    onCountDownChange(e) {
        const orderId = e.currentTarget.dataset.orderid;
        let curMenuId = this.data.curMenuId;
        const orderList = this.data.menuOrders[curMenuId];
        let selectOrder = null;
        for (var i = 0; i < orderList.length; i++) {
            let order = orderList[i];
            if (order.id == orderId) {
                selectOrder = order;
                break;
            }
        }
        if (!selectOrder) {
            return;
        }
        let data = e.detail;
        if (data.hours < 10) {
            data.hours = '0'+data.hours;
        }
        if (data.minutes < 10) {
            data.minutes = '0'+data.minutes;
        }
        if (data.seconds < 10) {
            data.seconds = '0'+data.seconds;
        }
        selectOrder.countDown = data
        this.setData({
            orderList: orderList
        })
    },
    parseOrder(order) {
        if ((order.status == 0 && order.tag == 2) 
            || (order.status == 1 && order.tag == 3)) {
            let endTime = order.endTime;
            if (!endTime) {
                return;
            }
            order.countDownTime = endTime*1000-Date.now();
        }
        // if (countDownTime >= 0) {
        //     let data = {};
        //     data.seconds = countDownTime%60;
        //     countDownTime = Math.floor(countDownTime/60);
        //     data.minutes = countDownTime%60;
        //     countDownTime = Math.floor(countDownTime/60);
        //     data.hours = countDownTime%24;
        //     countDownTime = Math.floor(countDownTime/24);
        //     if (countDownTime > 0) {
        //         data.days = countDownTime;
        //     }
        //     if (data.hours < 10) {
        //         data.hours = '0'+data.hours;
        //     }
        //     if (data.minutes < 10) {
        //         data.minutes = '0'+data.minutes;
        //     }
        //     if (data.seconds < 10) {
        //         data.seconds = '0'+data.seconds;
        //     }
        //     order.countDown = data;
        // }
    },
    requestData(){
        let that = this;
        WXAPI.orderList().then(function(res) {
            if (res.code == 0) {
                let resData = res.data;
                let orderList = resData.orderList;
                var menuOrders = {}
                menuOrders[8888] = orderList;
                for (var i = 0; i < orderList.length; i++) {
                    var order = orderList[i];
                    let menuId = 0;
                    switch(order.status){
                        case 0:
                            menuId = 0;
                            break;
                        case 1:
                            menuId = 1;
                            break;
                        case 2:
                            menuId = 2;
                            break;
                        case 3:
                        case 4:
                            menuId = 3;
                            break;
                        default:
                            menuId = 4;
                    }
                    var menuOrder = menuOrders[menuId];
                    if (!menuOrder) {
                        menuOrders[menuId] = [];
                        menuOrder = menuOrders[menuId];
                    }
                    menuOrder.push(order);
                    that.parseOrder(order);
                }

                const menuIds = [8888, 0, 1, 2, 3];
                const badges = that.data.badges || [];
                let curMenuId = that.data.curMenuId;
                menuIds.forEach((menuId) => {
                    var menuOrder = menuOrders[menuId];
                    badges[i] = 0;
                    if (menuOrder) {
                        menuOrder.sort(function(a, b) {
                            return b.timeStamp-a.timeStamp;
                        })
                        if (curMenuId != menuId) {
                            menuOrder.forEach((order) => {
                                if (order.isTips) {
                                    badges[i]++;
                                }
                            })
                        }
                    }
                })
                menuOrder = menuOrders[curMenuId] || [];

                that.data.menuOrders = menuOrders;
                that.setData({
                    goodsMap: resData.goodsMap,
                    orderList: menuOrder,
                    badges: badges
                });
            } else {
                that.data.menuOrders = {};
                that.setData({
                    goodsMap:{},
                    orderList: null,
                    badges: [0, 0, 0, 0, 0]
                });
            }
        })
    },
    hideOrderTap(e) {
        let that = this;
        const orderId = e.currentTarget.dataset.id;
        wx.showModal({
            title: '确认您已收到商品？',
            content: '',
            success(res) {
                if (res.confirm) {
                    WXAPI.orderClose(orderId).then(function(res) {
                        if (res.code == 0) {
                            that.onShow();
                        }
                    })
                }
            }
        })
    },
    receiveOrderTap(e){
        let that = this;
        const orderId = e.currentTarget.dataset.id;
        wx.showModal({
            title: '确认要删除订单？',
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
    refundApply(e) {
        // 申请售后
        const orderId = e.currentTarget.dataset.id;
        const amount = e.currentTarget.dataset.amount;
        wx.navigateTo({
            url: "/pages/order/refundApply?id=" + orderId + "&amount=" + amount
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
        }, 3000) // 可自行修改时间间隔（目前是3秒内只能点击一次支付按钮）
        // 防止连续点击--结束
        const self = this;
        const orderId = e.currentTarget.dataset.id;
        const amountReal = e.currentTarget.dataset.money;
        wxpay.requestPay(orderId, amountReal, function(){
            self.onShow();
            self.setMenuId(1);
        })
    },
    onLoad(options) {
        if (options && options.type) {
            var status = Number(options.type)
            if (options.type == 99) {
                this.setData({
                    hasRefund: true
                });
            } else {
                this.setData({
                    hasRefund: false,
                    curMenuId: status
                });
            }
        }
        this.setData({
            gm: wx.getStorageSync('gm')
        });
    },
    onReady() {
        // 生命周期函数--监听页面初次渲染完成
    },
    onShow() {
        AUTH.checkHasLogined().then(isLogined => {
            if (isLogined) {
                this.requestData();
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
    // goGM(){
    //     wx.navigateTo({
    //         url: "/pages/gm/index"
    //     })
    // },
    onHide() {
        // 生命周期函数--监听页面隐藏
    },
    onUnload() {
        // 生命周期函数--监听页面卸载
    },
    onPullDownRefresh() {
        // 页面相关事件处理函数--监听用户下拉动作
    },
    onReachBottom() {
        // 页面上拉触底事件的处理函数
    }
})