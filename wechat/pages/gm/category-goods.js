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

        startDate: "2021-03-12",
        endDate: "2021-03-12",
        startTime: "00:01",
        endTime: "23:59",

        fromDatePattern: "",
        toDatePattern: "2031-03-12",
        fromTimePattern: "00:01",
        toTimePattern: "23:59",
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
        this.setData({
            fromDatePattern: TOOLS.time2dateformat()
            // toDatePattern: TOOLS.time2dateformat()
        });
    },
    setSelectCategoryId(selectCategoryId){
        this.setData({
            selectCategoryId: selectCategoryId
        });
        this.data.goodsList = [];
        this.getGoodsList();
        wx.setStorageSync('gm_categoryId', selectCategoryId)
    },
    bindStartDateChange: function(e) {
        this.setData({
            startDate: e.detail.value
        })
    },
    bindStartTimeChange: function(e) {
        this.setData({
            startTime: e.detail.value
        })
    },
    bindEndDateChange: function(e) {
        this.setData({
            endDate: e.detail.value
        })
    },
    bindEndTimeChange: function(e) {
        this.setData({
            endTime: e.detail.value
        })
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
        let goodsId = e.currentTarget.dataset.id;
        wx.navigateTo({
            url: "/pages/goods-details/index?id=" + goodsId
        });
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
        // wx.navigateTo({url: "/pages/gm/goods-edit?id="+goodsId});
        for (var i = 0; i < this.data.goodsList.length; i++) {
            let goods = this.data.goodsList[i];
            if (goods.goodsId == goodsId) {
                let selectGoods = {};
                for (var key in goods) {
                    selectGoods[key] = goods[key]
                }
                this.data._selectGoods = goods;
                this.setCategoryId(goods.categoryId)

                this.setData({
                    selectGoods: selectGoods,
                    startDate: TOOLS.time2dateformat(goods.startTime),
                    startTime: TOOLS.time2timeformat(goods.startTime),
                    endDate: TOOLS.time2dateformat(goods.endTime),
                    endTime: TOOLS.time2timeformat(goods.endTime)
                });
                break;
            }
        }
    },
    closePopup() {
        let selectGoods = this.data.selectGoods;
        selectGoods.categoryId = this.getCategoryId()
        let keys = Object.keys(selectGoods);
        for (var i = 0; i < keys.length; i++) {
            let key = keys[i];
            if (selectGoods[key] != this.data._selectGoods[key]) {
                const that = this;
                wx.showModal({
                    title: '确定取消编辑？',
                    content: '',
                    success: function(res) {
                        if (res.confirm) {
                            that.setData({
                                selectGoods: null
                            });
                        }
                    }
                });
                return;
            }
        }
        this.setData({
            selectGoods: null
        });
    },
    numMinusTap(e) {
        let selectGoods = this.data.selectGoods;
        const field = e.currentTarget.dataset.field;
        let fieldValue = selectGoods[field];
        if (fieldValue <= 0) {
            fieldValue = 0;
        }else{
            fieldValue--;
        }
        selectGoods[field] = fieldValue;
        this.setData({
            selectGoods: selectGoods,
        });
    },
    numPlusTap(e) {
        let selectGoods = this.data.selectGoods;
        const field = e.currentTarget.dataset.field;
        let fieldValue = selectGoods[field];
        fieldValue++;
        selectGoods[field] = fieldValue;
        this.setData({
            selectGoods: selectGoods,
        });
    },
    watchInput(e) {
        let selectGoods = this.data.selectGoods;
        const field = e.currentTarget.dataset.field;
        let fieldValue = Number(e.detail.value);
        if (!fieldValue && fieldValue != 0) {
            fieldValue = selectGoods[field];
        }
        selectGoods[field] = fieldValue;
        this.setData({
            selectGoods: selectGoods,
        });
    },
    // watchInputStartDate(e) {
    //     this.data.startDate = e.detail.value
    // },
    // watchInputEndDate(e) {
    //     this.data.endDate = e.detail.value
    // },
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
    statusChange(e) {
        let that = this;
        wx.showActionSheet({
            itemList: this.data.statusNames,
            success (res) {
                let selectGoods = that.data.selectGoods;
                selectGoods.status = res.tapIndex;
                that.setData({
                    selectGoods: selectGoods,
                });
            },
            fail (res) {
                // console.log(res.errMsg)
            }
        })
    },

    markChange(){
        let that = this;
        wx.showActionSheet({
            itemList: this.data.markNames,
            success (res) {
                let selectGoods = that.data.selectGoods;
                selectGoods.mark = res.tapIndex;
                that.setData({
                    selectGoods: selectGoods,
                });
            },
            fail (res) {
                // console.log(res.errMsg)
            }
        })
    },
    requestUpdate(){
        let selectGoods = this.data.selectGoods;
        this.setData({
            selectGoods: null
        });
        if (!selectGoods) {
            wx.showToast({title: '未选择商品', icon: 'none'})
            return;
        }
        selectGoods.categoryId = this.getCategoryId()
        if (selectGoods.mark == 2) {
            let startTime = Math.floor(Date.parse(this.data.startDate + ' ' + this.data.startTime)/1000);
            if (startTime == 0 || isNaN(startTime)) {
                wx.showToast({title: '限时秒杀开始时间错误', icon: 'none'})
                return
            }
            let endTime = Math.floor(Date.parse(this.data.endDate + ' ' + this.data.endTime)/1000);
            if (endTime == 0  || isNaN(endTime)) {
                wx.showToast({title: '限时秒杀结束时间错误', icon: 'none'})
                return
            }
            if (startTime > endTime) {
                wx.showToast({title: '开始时间比结束时间大', icon: 'none'})
                return
            }
            selectGoods.startTime = startTime
            selectGoods.endTime = endTime
        } else {
            selectGoods.startTime = 0
            selectGoods.endTime = 0
        }
        this.clearCache();
        let that = this;
        WXAPI.gmGoodsUpdateInfo(selectGoods).then(function(res) {
            if (res.code == 0) {
                wx.showToast({title: '更新成功', icon: 'none'});
                that.requestCategoryList(that.data.selectCategoryId);
            }else{
                wx.showToast({title: res.msg, icon: 'none'});
            }
        });
    },
    requestRemove(){
        if (this.data._selectGoods.status == 2) {
            wx.showToast({title: '上架商品不能删除', icon: 'none'})
            return
        }
        let selectGoods = this.data.selectGoods;
        this.setData({
            selectGoods: null
        });
        if (!selectGoods) {
            wx.showToast({title: '未选择商品', icon: 'none'})
            return;
        }
        this.clearCache();
        let that = this;
        WXAPI.gmGoodsRemove(selectGoods.goodsId).then(function(res) {
            if (res.code == 0) {
                wx.showToast({title: '移除成功', icon: 'none'});
                that.requestCategoryList(that.data.selectCategoryId);
            }else{
                wx.showToast({title: res.msg, icon: 'none'});
            }
        });
    },
    bindinput(e) {
        this.setData({
            inputVal: e.detail
        })
    },
    bindconfirm(e) {
        let searchName = e.detail;
        this.setData({
            inputVal: searchName
        });
        if (!searchName) {return}

        let dataList = this.data.originGoodsList;
        let goodsList = [];
        for (var i = 0; i < dataList.length; i++) {
            let data = dataList[i];
            if (data.name.indexOf(searchName) > 0) {
                goodsList.push(data)
            }
        }
        this.setData({
            goodsList: goodsList
        });
    },
    requestCategoryList(selectCategoryId) {
        const that = this;
        WXAPI.gmCategoryList().then(res => {
            if (res.code == 0) {
                let categoryList = res.data.categorys || [];
                categoryList.sort(function(a, b){
                    return a.order - b.order;
                });
                that.setData({
                    categories: categoryList,
                });

                let categoryNames = [];
                let categoryIds = [];
                for (var i = 0; i < categoryList.length; i++) {
                    let item = categoryList[i];
                    categoryNames.push(item.name);
                    categoryIds.push(item.id);
                }
                that.setData({
                    categoryNames: categoryNames,
                    categoryIds:categoryIds,
                });

                let category = categoryList[0];
                if (category) {
                    that.setSelectCategoryId(selectCategoryId || category.id);
                }
            }
        });
    },
    onShow() {
        this.data.groupGoodsList = {}
        let that = this;
        const categoryId = wx.getStorageSync('gm_categoryId')
        AUTH.checkHasLogined().then(isLogined => {
            if (isLogined) {
                this.setData({
                    wxlogin: isLogined
                })
            }
            that.requestCategoryList(categoryId);
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