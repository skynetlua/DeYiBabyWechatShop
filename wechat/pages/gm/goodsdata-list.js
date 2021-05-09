const CONFIG = require('../../config.js')
const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')
const TOOLS = require('../../utils/tools.js')

Page({
    /**
     * 页面的初始数据
     */
    data: {
        categories: [],
        selectCategoryId:0,

        goodsList: [],
        onLoadStatus: true,
        scrolltop: 0,
        groupGoodsList: {},
        selectGoods: null,

        categorys: [],
        categoryIds:[],
        categoryIndex: 0,

        statusNames: ['录入', '待上架', '上架', '下架'],
        markNames: ['普通', '推荐', '秒杀', '拼团'],

    },

    getCategoryId(){
        let categoryIndex = this.data.categoryIndex;
        let catogoryId = this.data.categoryIds[categoryIndex];
        return catogoryId;
    },
    setCategoryId(catogoryId){
        for (var i = 0; i < this.data.categoryIds.length; i++) {
            let id = this.data.categoryIds[i];
            if (id == catogoryId) {
                this.setData({
                    categoryIndex: i
                });
                break;
            }
        }
    },
    categoryChange(e) {
        this.setData({
            categoryIndex: e.detail.value
        });
    },
    onLoad(options) {
    },
    setSelectCategoryId(selectCategoryId){
        this.setData({
            selectCategoryId: selectCategoryId
        });
        // this.data.goodsList = [];
        // this.getGoodsList();
    },
    
    clearCache() {
        let categoryId = this.data.selectCategoryId;
        this.data.groupGoodsList[categoryId] = null;
    },
    async getGoodsList() {
        let categoryId = this.data.selectCategoryId;
        let res = this.data.groupGoodsList[categoryId]
        // let reqKey = 'gm-category-goods'+ categoryId
        // let res = APP.getCache(reqKey);
        if(!res) {
            wx.showLoading({title: '加载中'});
            res = await WXAPI.gmGoodsCategory({
                categoryId: categoryId,
            });
            this.data.groupGoodsList[categoryId] = res;
            // APP.setCache(reqKey, res);
            wx.hideLoading();
        }
        let goodsList = res.data
        goodsList.sort(function(a, b){
            return b.updateTime-a.updateTime
        })
        this.setData({
            goodsList: goodsList,
            originGoodsList: goodsList
        });
    },
    toDetailsTap(e) {
        // let goodsId = e.currentTarget.dataset.id;
        // wx.navigateTo({
        //     url: "/pages/goods-details/index?id=" + goodsId
        // });
    },
    onCategoryClick(e) {
        var categoryId = e.target.dataset.id;
        if (categoryId === this.data.selectCategoryId) {
            this.setData({
                scrolltop: 0,
            });
        } else {
            this.setSelectCategoryId(categoryId);
        }
    },
    selectGoodsTap(e){
        let goodsId = e.currentTarget.dataset.id;
        
    },
    goodsEditTap(){
        let selectGoods = this.data.selectGoods;
        this.setData({
            selectGoods: null
        });
        if (!selectGoods) {
            wx.showToast({title: '未选择商品', icon: 'none'})
            return;
        }
        wx.navigateTo({url: "/pages/gm/goods-edit?id="+selectGoods.goodsId});
    },
    onShow() {
        this.data.groupGoodsList = {}
        let that = this;
        AUTH.checkHasLogined().then(isLogined => {
            if (isLogined) {
                this.setData({
                    wxlogin: isLogined
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