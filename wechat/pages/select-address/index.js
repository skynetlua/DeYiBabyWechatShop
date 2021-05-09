const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')

Page({
    data: {
        addressList: []
    },

    selectTap: function(e) {
        var id = e.currentTarget.dataset.id;
        WXAPI.updateAddress({
            id: id,
            isDefault: 1
        }).then(function(res) {
            wx.navigateBack({})
        })
    },

    addAddess: function() {
        wx.navigateTo({
            url: "/pages/address-add/index"
        })
    },

    editAddess: function(e) {
        wx.navigateTo({
            url: "/pages/address-add/index?id=" + e.currentTarget.dataset.id
        })
    },

    onLoad: function() {},
    onShow: function() {
        AUTH.checkHasLogined().then(isLogined => {
            if (isLogined) {
                this.initShippingAddress();
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
    initShippingAddress: function() {
        var that = this;
        WXAPI.queryAddress().then(function(res) {
            if (res.code == 0) {
                that.setData({
                    addressList: res.data.addresss
                });
            } else if (res.code == 700) {
                that.setData({
                    addressList: null
                });
            }
        })
    }

})