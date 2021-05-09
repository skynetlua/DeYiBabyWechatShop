const WXAPI = require('apifm-wxapi')

var sliderWidth = 96;
Page({

    /**
     * 页面的初始数据
     */
    data: {
        members1: [],
        members2: [],
        activeIndex: 0,
        sliderOffset: 0,
        sliderLeft: 0,
        // grids: [0, 1, 2, 3, 4, 5, 6, 7, 8]
    },

    onLoad: function() {
        var that = this;
        wx.getSystemInfo({
            success: function(res) {
                that.setData({
                    sliderLeft: (res.windowWidth / 2 - sliderWidth) / 2,
                    sliderOffset: res.windowWidth / 2 * that.data.activeIndex
                });
            }
        });
    },
    tabClick: function(e) {
        this.setData({
            sliderOffset: e.currentTarget.offsetLeft,
            activeIndex: e.currentTarget.id
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
        const _this = this;
        WXAPI.fxMembers({
            pageSize: 1000
        }).then(res => {
            if (res.code == 700) {
                _this.setData({
                    members1: [],
                    members2: [],
                });
            }
            if (res.code == 0) {
                let members1 = [];
                let members2 = [];
                res.data.forEach(ele => {
                    if (ele.level == 1) {
                        members1.push(ele);
                    }else if (ele.level == 2) {
                        members2.push(ele);
                    }
                });
                _this.setData({
                    members1: members1,
                    members2: members2,
                });
            }
        });
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