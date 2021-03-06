const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')

Page({
    data: {},
    onLoad: function(options) {},
    onShow: function() {
        AUTH.checkHasLogined().then(isLogined => {
            this.setData({
                wxlogin: isLogined
            })
            if (isLogined) {
                this.goodsFavList()
            }
        })
    },
    async goodsFavList() {
        // 搜索商品
        wx.showLoading({
            title: '加载中',
        })
        const _data = {
            page: 1,
            pageSize: 200,
        }
        const res = await WXAPI.goodsFavList(_data)
        wx.hideLoading()
        if (res.code == 0) {
            this.setData({
                goods: res.data,
            })
        } else {
            this.setData({
                goods: null
            })
        }
    },
    async removeFav(e) {
        const goodsId = e.currentTarget.dataset.id
        const res = await WXAPI.goodsFavDelete(goodsId)
        if (res.code == 0) {
            wx.showToast({
                title: '取消收藏',
                icon: 'success'
            })
            this.goodsFavList()
        } else {
            wx.showToast({
                title: res.msg,
                icon: 'none'
            })
        }
    },
})