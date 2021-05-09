const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')

Page({
    /**
     * 页面的初始数据
     */
    data: {
    },

    onLoad: function(options) {
        
    },

    /**
     * 生命周期函数--监听页面初次渲染完成
     */
    onReady: function() {

    },

    uploadGoodsExcel(){
        let that = this;
        wx.chooseMessageFile({
            count: 1,
            type: 'file',
            success(res) {
                const tempFile = res.tempFiles[0];
                if (!tempFile) {
                    wx.showToast({title: "请选择Excel文件", icon: 'none'});
                    return;
                }
                if (tempFile.name.lastIndexOf(".xlsx")<=0) {
                    wx.showToast({title: "请选择格式为.xlsx的Excel文件", icon: 'none'});
                    return;
                }
                that.uploadFile(tempFile.path, tempFile.name);
            }
        });
    },

    uploadFile(filePath, fileName) {
        let opt = "goods";
        let that = this;
        WXAPI.gmUploadExcel(filePath, fileName, opt).then(function(res) {
            if (res.code == 0) {
                wx.showToast({title: "商品导入成功", icon: 'none'});
                if (res.data && res.data.length > 0) {
                    wx.showModal({
                        title: '提示',
                        content: '下列商品导入失败：\n'+res.data.join("\n"),
                    });
                }
            }else{
                wx.showToast({title: res.msg, icon: 'none'});
            }
        });
    },

    goodsLoadPicture(){
        WXAPI.gmGoodsLoadPicture().then(function(res) {
            if (res.code == 0) {
                wx.showToast({title: "商品导入成功", icon: 'none'});
                wx.showModal({
                    title: '提示',
                    content: '下列商品导入失败：\n'+res.data.join("\n"),
                });
            }else{
                wx.showToast({title: res.msg, icon: 'none'});
            }
        });
    },

    /**
     * 生命周期函数--监听页面显示
     */
    onShow: function() {
        const that = this;
        AUTH.checkHasLogined().then(isLogined => {
            if (isLogined) {
                that.setData({
                    wxlogin: isLogined
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

    },

    goIndex() {
        wx.switchTab({
            url: '/pages/index/index',
        });
    },
    cancelLogin() {
        this.setData({
            wxlogin: true
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