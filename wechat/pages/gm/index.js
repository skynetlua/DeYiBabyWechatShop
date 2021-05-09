const CONFIG = require('../../config.js')
const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')
const TOOLS = require('../../utils/tools.js')

Page({
    data: {
        wxlogin: true,
        showInput: false,
        barCode: ''
    },
    onLoad() {},
    onShow() {
        const _this = this
        AUTH.checkHasLogined().then(isLogined => {
            this.setData({
                wxlogin: isLogined
            })
        })
    },
    aboutUs: function() {
    },
    // loginOut() {
    //     AUTH.loginOut()
    //     wx.reLaunch({url: '/pages/my/index'})
    // },
    categoryList: function(){
        wx.navigateTo({url: "/pages/gm/category-list"});
    },
    goodsCategory: function(){
        wx.navigateTo({url: "/pages/gm/category-goods"});
    },
    orderList: function() {
        wx.navigateTo({url: "/pages/gm/order-list"});
    },
    refund: function() {
        wx.navigateTo({url: "/pages/gm/refund"});
    },
    teamList: function(e) {
        wx.navigateTo({url: "/pages/gm/team-list"});
    },
    goodsEdit: function(e) {
        wx.navigateTo({url: "/pages/gm/goods-edit"});
    },
    goGoodsEdit:function() {
        let barCode = this.data.barCode;
        wx.showLoading({title: '查询中'})
        WXAPI.gmBarCodeGoodsId(barCode).then(function(res) {
            wx.hideLoading();
            if (res.code == 0) {
                wx.navigateTo({url: "/pages/gm/goods-edit?barcode="+barCode});
            }else if (res.code == 1) {
                wx.showToast({title:"很抱歉，后台无该商品数据", icon: 'none'})
            }else{
                wx.showToast({title:res.msg, icon: 'none'})
            }
        });
    },
    onGoodsQR: function(e){
        const that = this;
        wx.scanCode({
            onlyFromCamera: false,
            success(res) {
                if (res.result) {
                    that.setData({
                        barCode: res.result
                    });
                    that.goGoodsEdit();
                }
            },
            fail(err) {
                wx.showToast({
                    title: err.errMsg,
                    icon: 'none'
                });
            }
        })
    },
    onGoodsCode: function() {
        this.setData({
            showInput: true
        });
    },
    closePopup: function() {
        this.setData({
            showInput: false
        });
    },
    watchInput(e) {
        this.setData({
            barCode: e.detail.value
        });
    },
    goodsList: function(e) {
        wx.navigateTo({url: "/pages/gm/goods-list"});
    },
    controlPanel(){
        wx.navigateTo({url: "/pages/gm/control-panel"});
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
    }
})