const WXAPI = require('apifm-wxapi')
Page({

    /**
     * 页面的初始数据
     */
    data: {
        commissions: []
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
    onShow: function() {
        const _this = this
        WXAPI.fxCommisionLog({}).then(res => {
            if (res.code == 0 && res.data) {
                _this.setData({
                    commissions: res.data
                })
            } else {
                _this.setData({
                    commissions: []
                })
            }
        })
    },

    /**
     * 生命周期函数--监听页面隐藏
     */
    onHide: function() {

    },

    /**
     * 生命周期函数--监听页面卸载
     */
    onUnload: function() {

    },

    /**
     * 页面相关事件处理函数--监听用户下拉动作
     */
    onPullDownRefresh: function() {

    },

    /**
     * 页面上拉触底事件的处理函数
     */
    onReachBottom: function() {

    }
})