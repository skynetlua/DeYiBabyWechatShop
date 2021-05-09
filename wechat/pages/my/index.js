const CONFIG = require('../../config.js')
const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')
const TOOLS = require('../../utils/tools.js')

Page({
    data: {
        wxlogin: true,
        balance: 0.00,
        score: 0,
        growth: 0,
        vipLevel: 1,
        version:CONFIG.version,
        isInit: false
        // score_sign_continuous: 0,
        // rechargeOpen: false // 是否开启充值[预存]功能
    },
    loadData() {
        let my_time = wx.getStorageSync('my_index_time') || 0
        my_time = Number(my_time)
        let curTime = Date.now()
        if (this.data.isInit && my_time+1000*60*3 > curTime) {
            return
        }
        wx.setStorageSync('my_index_time', curTime)
        this.data.isInit = true

        const that = this
        this.setData({
            gm: wx.getStorageSync('gm')
        })
        AUTH.checkHasLogined().then(isLogined => {
            var nickName = wx.getStorageSync('nickName');
            isLogined = isLogined && nickName != '';
            this.setData({
                wxlogin: isLogined
            })
            if (isLogined) {
                that.getUserInfo();
                // TOOLS.showTabBarBadge();
            }
        })
    },
    onLoad() {
        this.loadData()
    },
    onShow() {
        this.loadData()
    },
    aboutUs() {
        wx.showModal({
            title: '关于我们',
            content: 'Q-Baby母婴生活馆（敏捷店）\n微信号：18926627466\n客服电话：18926627466\n地址：清远市清城区东城街道敏捷水岸花园1号楼',
            showCancel: false
        })
    },
    loginOut() {
        AUTH.loginOut()
        wx.reLaunch({
            url: '/pages/my/index'
        })
    },
    getPhoneNumber(e) {
        if (!e.detail.errMsg || e.detail.errMsg != "getPhoneNumber:ok") {
            wx.showModal({
                title: '提示',
                content: e.detail.errMsg,
                showCancel: false
            })
            return;
        }
        WXAPI.bindMobileWxa(e.detail.encryptedData, e.detail.iv).then(res => {
            if (res.code === 10002) {
                this.setData({
                    wxlogin: false
                })
                return
            }
            if (res.code == 0) {
                wx.showToast({
                    title: '绑定成功',
                    icon: 'success',
                    duration: 2000
                })
                this.getUserInfo();
            } else {
                wx.showModal({
                    title: '提示',
                    content: res.msg,
                    showCancel: false
                })
            }
        })
    },
    getUserInfo() {
        var that = this;
        WXAPI.userDetail().then(function(res) {
            if (res.code == 0) {
                let userInfo = res.data.userInfo
                if (!userInfo.nickName || userInfo.nickName.length == 0) {
                    that.setData({
                        wxlogin: false
                    })
                    return
                }
                let _data = {
                    userInfo: userInfo,
                    isSeller: userInfo.isSeller,
                }
                // if (userInfo.mobile) {
                //     _data.userMobile = data.base.mobile
                // }
                // if (that.data.order_hx_uids && that.data.order_hx_uids.indexOf(data.base.id) != -1) {
                //     _data.canHX = true // 具有扫码核销的权限
                // }
                that.setData(_data);
            }
        })
    },
    getUserAmount() {
        var that = this;
        WXAPI.userAmount().then(function(res) {
            if (res.code == 0) {
                let data = res.data
                that.setData({
                    balance: data.balance.toFixed(2),
                    score: data.score,
                    growth: data.growth
                });
            }
        })
    },
    goAsset() {
        wx.navigateTo({url: "/pages/asset/index"})
    },
    goScore() {
        wx.navigateTo({
            url: "/pages/score/index"
        })
    },
    goOrder(e) {
        wx.navigateTo({
            url: "/pages/order-list/index?type=" + e.currentTarget.dataset.type
        })
    },
    cancelLogin() {
        this.setData({wxlogin: true});
    },
    goLogin() {
        this.setData({wxlogin: false});
    },
    processLogin(e) {
        if (!e.detail.userInfo) {
            wx.showToast({title: '已取消',icon: 'none'});
            return;
        }
        AUTH.register(this);
    },
    scanOrderCode() {
        wx.scanCode({
            onlyFromCamera: true,
            success(res) {
                wx.navigateTo({
                    url: '/pages/order-details/scan-result?hxNumber=' + res.result,
                });
            },
            fail(err) {
                console.error(err);
                wx.showToast({
                    title: err.errMsg,
                    icon: 'none'
                });
            }
        })
    },
    clearStorage() {
        wx.clearStorageSync();
        wx.showToast({title: '已清除',icon: 'success'});
    },
})