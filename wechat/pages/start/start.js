const WXAPI = require('apifm-wxapi')
const CONFIG = require('../../config.js')
//获取应用实例
Page({
    data: {
        banners: [],
        swiperMaxNumber: 0,
        swiperCurrent: 0
    },
    onLoad: function() {
        const _this = this
        // wx.setNavigationBarTitle({
        //     title: wx.getStorageSync('mallName')
        // })
        const version = wx.getStorageSync('app_show_pic_version')
        if (version && version == CONFIG.version) {
            if (CONFIG.shopMod) {
                wx.redirectTo({url: '/pages/shop/select'})
            } else {
                wx.switchTab({url: '/pages/index/index'})
            }
        } else {
            // 展示启动页
            WXAPI.banners({
                types: 'app'
            }).then(function(res) {
                if (res.code == 700) {
                    if (CONFIG.shopMod) {
                        wx.redirectTo({url: '/pages/shop/select'})
                    } else {
                        wx.switchTab({url: '/pages/index/index'})
                    }
                } else {
                    let banners = res.data.app
                    _this.setData({
                        banners: banners,
                        swiperMaxNumber: banners.length
                    })
                }
            }).catch(function(e) {
                if (CONFIG.shopMod) {
                    wx.redirectTo({url: '/pages/shop/select'})
                } else {
                    wx.switchTab({url: '/pages/index/index'})
                }
            })
        }
    },
    onShow: function() {
    },
    swiperchange: function(e) {
        this.setData({
            swiperCurrent: e.detail.current
        })
    },
    goToIndex: function(e) {
        let isConnected = wx.getStorageSync('isConnected')
        if (isConnected) {
            wx.setStorage({key: 'app_show_pic_version', data: CONFIG.version})
            if (CONFIG.shopMod) {
                wx.redirectTo({url: '/pages/shop/select'})
            } else {
                wx.switchTab({url: '/pages/index/index'})
            }
        } else {
            wx.showToast({title: '当前无网络', icon: 'none'})
        }
    },
    imgClick() {
        if (this.data.swiperCurrent + 1 != this.data.swiperMaxNumber) {
            wx.showToast({title: '左滑进入', icon: 'none' })
        }
    }
});