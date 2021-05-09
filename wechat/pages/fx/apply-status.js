const CONFIG = require('../../config.js')
const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')
const ImageUtil = require('../../utils/image')

Page({

    /**
     * 页面的初始数据
     */
    data: {
        wxlogin: true,

        applyStatus: -2, // -1 表示未申请，0 审核中 1 不通过 2 通过
        applyInfo: {},
        canvasHeight: 0,
        currentPages: undefined,
    },

    /**
     * 生命周期函数--监听页面加载
     */
    onLoad: function(options) {
    },

    /**
     * 生命周期函数--监听页面初次渲染完成
     */
    onReady: function() {
    },

    /**
     * 生命周期函数--监听页面显示
     */
    onShow() {
        const _this = this
        AUTH.checkHasLogined().then(isLogined => {
            this.setData({
                wxlogin: isLogined
            })
            if (isLogined) {
                this.doneShow();
            }
        })
    },
    async doneShow() {
        const _this = this
        WXAPI.fxApplyProgress().then(res => {
            if (res.code == 2000) {
                this.setData({
                    wxlogin: false
                })
                return
            }
            if (res.code === 0) {
                _this.setData({
                    status: res.data.status,
                    remark: res.data.remark,
                })
                if (res.data.status == 3) {
                    _this.fetchQrcode()
                }
            }
        })
        this.setData({
            currentPages: getCurrentPages()
        })
    },

    fetchQrcode() {
        wx.showLoading({
            title: '加载中',
            mask: true
        })
        const _this = this
        WXAPI.wxaQrcode({
            scene: 'qf,' + wx.getStorageSync('uid'),
            page: 'pages/start/loading',
            is_hyaline: true,
            expireHours: 1
        }).then(res => {
            wx.hideLoading()
            if (res.code == 0) {
                _this.showCanvas(res.data)
            }else{
                wx.showToast({
                    title: res.msg,
                    icon: 'none'
                })
            }
        })
    },

    showCanvas(qrcode) {
        const _this = this
        let ctx
        wx.getImageInfo({
            src: qrcode,
            success: (res) => {
                const imageSize = ImageUtil.imageUtil(res.width, res.height)
                const qrcodeWidth = imageSize.windowWidth / 2
                _this.setData({
                    canvasHeight: qrcodeWidth
                })
                ctx = wx.createCanvasContext('firstCanvas')
                ctx.setFillStyle('#fff')
                ctx.fillRect(0, 0, imageSize.windowWidth, imageSize.imageHeight + qrcodeWidth)
                ctx.drawImage(res.path, (imageSize.windowWidth - qrcodeWidth) / 2, 0, qrcodeWidth, qrcodeWidth)
                setTimeout(function() {
                    wx.hideLoading()
                    ctx.draw()
                }, 1000)
            }
        })
    },
    onShareAppMessage() {
        return {
            title: '"' + wx.getStorageSync('mallName') + '" ' + CONFIG.shareProfile,
            path: '/pages/start/loading?inviter_id=' + wx.getStorageSync('uid') + '&route=/pages/index/index'
        }
    },
    bindSave: function(e) {
        wx.navigateTo({
            url: "/pages/fx/apply"
        })
    },
    goShop: function(e) {
        wx.switchTab({
            url: '/pages/index/index',
        })
    },
    saveToMobile() { //下载二维码到手机
        wx.canvasToTempFilePath({
            canvasId: 'firstCanvas',
            success: function(res) {
                let tempFilePath = res.tempFilePath
                wx.saveImageToPhotosAlbum({
                    filePath: tempFilePath,
                    success: (res) => {
                        wx.showModal({
                            content: '二维码已保存到手机相册',
                            showCancel: false,
                            confirmText: '知道了',
                            confirmColor: '#333'
                        })
                    },
                    fail: (res) => {
                        wx.showToast({
                            title: res.errMsg,
                            icon: 'none',
                            duration: 2000
                        })
                    }
                })
            }
        })
    },
    goIndex() {
        wx.switchTab({
            url: '/pages/index/index',
        });
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