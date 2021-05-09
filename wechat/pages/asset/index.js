const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')

var sliderWidth = 96; // 需要设置slider的宽度，用于计算中间位置

Page({

    /**
     * 页面的初始数据
     */
    data: {
        wxlogin: true,
        balance: 0.00,
        freeze: 0,
        score: 0,
        cashlogs: undefined,

        tabs: ["资金明细", "提现记录"],
        activeIndex: 0,
        sliderOffset: 0,
        sliderLeft: 0,

        withdrawLogs: undefined,
        rechargeOpen: false,
        datas:{},
    },

    /**
     * 生命周期函数--监听页面加载
     */
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

    /**
     * 生命周期函数--监听页面初次渲染完成
     */
    onReady: function() {
    },

    /**
     * 生命周期函数--监听页面显示
     */
    onShow: function() {
        AUTH.checkHasLogined().then(isLogined => {
            this.setData({
                wxlogin: isLogined
            })
            if (isLogined) {
                this.doneShow();
            }
        })
    },
    doneShow: function() {
        const _this = this
        const token = wx.getStorageSync('token')
        if (!token) {
            this.setData({
                wxlogin: false
            })
            return
        }
        WXAPI.userAmount(token).then(function(res) {
            if (res.code == 700) {
                wx.showToast({
                    title: '当前账户存在异常',
                    icon: 'none'
                })
                return
            }
            if (res.code == 2000) {
                this.setData({
                    wxlogin: false
                })
                return
            }
            if (res.code == 0) {
                let data = res.data;
                _this.setData({
                    balance: data.balance.toFixed(2),
                    // amountFreeze: data.freeze.toFixed(2),
                    amountCost: data.amountCost.toFixed(2),
                    // score: res.data.score
                });
            }
        })
        this.fetchTabData(this.data.activeIndex)
    },
    fetchTabData(activeIndex) {
        if (activeIndex == 0) {
            this.cashLogs()
        }else if (activeIndex == 1) {
            this.withdrawLogs()
        // }else if (activeIndex == 2) {
            // this.depositlogs()
        }
    },
    cashLogs() {
        let params = {
            page: 0,
            pageSize: 50
        }
        let key = ["cashLogs", params.page, params.pageSize].join(",");
        let data = this.data.datas[key];
        if (data) {
            this.setData({
                cashlogs: data,
            })
            return;
        }
        const self = this
        WXAPI.cashLogs(params).then(res => {
            if (res.code == 0) {
                self.data.datas[key] = res.data;
                self.setData({
                    cashlogs: res.data
                })
            }
        });
    },
    withdrawLogs() {
        let params = {
            page: 0,
            pageSize: 50
        }
        let key = ["withdrawLogs", params.page, params.pageSize].join(",");
        let data = this.data.datas[key];
        if (data) {
            this.setData({
                withdrawLogs: data,
            })
            return;
        }
        const self = this
        WXAPI.withDrawLogs(params).then(res => {
            if (res.code == 0) {
                self.data.datas[key] = res.data;
                self.setData({
                    withdrawLogs: res.data
                })
            }
        })
    },

    // depositlogs() {
    //     const _this = this
    //     WXAPI.depositList({
    //         page: 1,
    //         pageSize: 50
    //     }).then(res => {
    //         if (res.code == 0) {
    //             _this.setData({
    //                 depositlogs: res.data.result
    //             })
    //         }
    //     })
    // },

    recharge: function(e) {
        wx.navigateTo({
            url: "/pages/recharge/index"
        })
    },
    withdraw: function(e) {
        wx.navigateTo({
            url: "/pages/withdraw/index"
        })
    },
    
    // payDeposit: function(e) {
    //     wx.navigateTo({
    //         url: "/pages/deposit/pay"
    //     })
    // },
    tabClick: function(e) {
        this.setData({
            sliderOffset: e.currentTarget.offsetLeft,
            activeIndex: e.currentTarget.id
        });
        this.fetchTabData(e.currentTarget.id)
    },
    cancelLogin() {
        wx.switchTab({
            url: '/pages/my/index'
        })
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