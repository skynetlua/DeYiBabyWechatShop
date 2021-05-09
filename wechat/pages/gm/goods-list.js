const CONFIG = require('../../config.js')
const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')
const TOOLS = require('../../utils/tools.js') // TOOLS.showTabBarBadge();

Page({
    /**
     * 页面的初始数据
     */
    data: {
        categories: [
            {id:0,name:"录入"},
            {id:1,name:"待上架"},
            {id:2,name:"上架"},
            {id:3,name:"下架"},
        ],
        selectCategoryId:0,
        goodsList: [],
        onLoadStatus: true,
        scrolltop: 0,
        selectGoods:null,
        cache:{},
        categorys: [],
        categoryIds:[],
        categoryIndex: 0,
        startDate: "2021/3/12 16:27:00",
        endDate: "2021/3/12 16:27:00"
    },
    requestCategoryList() {
        const that = this;
        WXAPI.gmCategoryList().then(res => {
            if (res.code == 0) {
                let categoryList = res.data.categorys || [];
                categoryList.sort(function(a, b){
                    return a.order - b.order;
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
            }
        });
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
        let startDate = TOOLS.time2datestr()
        let endDate = TOOLS.time2datestr()
        this.setData({
            startDate: startDate,
            endDate: endDate
        });
        this.requestCategoryList()
    },
    setSelectId:function(selectCategoryId){
        this.setData({
            selectCategoryId: selectCategoryId
        });
        this.getGoodsList();
    },
    getCache(key){
        let cache = this.data.cache;
        if (!cache) {
            this.data.cache = {};
            cache = this.data.cache;
        }
        // let nowTime = Date.now();
        // let saveTime = cache.saveTime || 0;
        // cache.saveTime = nowTime+1000*60;
        // if (saveTime > nowTime) {
            return cache[key];
        // }
    },
    setCache(key, val){
        // let nowTime = Date.now();
        let cache = this.data.cache;
        if (!cache) {
            this.data.cache = {};
            cache = this.data.cache;
        }
        // cache.saveTime = nowTime+1000*60;
        cache[key] = val;
    },
    async getGoodsList(page, pagesize) {
        let status = this.data.selectCategoryId;
        let params = {
            status: status,
            page: page || 0,
            pageSize: pagesize || 500
        }
        let reqKey = [params.status, params.page, params.pageSize].join(":");
        let res = this.getCache(reqKey, 10);
        if(res){
            // console.log("getGoodsList=====>> 无需下载");
            // wx.hideLoading()
        }else{
            // console.log("getGoodsList=====>> 下载");
            wx.showLoading({title: '加载中',});
            res = await WXAPI.gmGoodsList(params);
            this.setCache(reqKey, res);
            wx.hideLoading();
        }
        this.setData({
            goodsList: res.data
        });
    },
    toDetailsTap: function(e) {
        let goodsId = e.currentTarget.dataset.id;
        wx.navigateTo({
            url: "/pages/goods-details/index?id=" + goodsId
        })
    },
    onCategoryClick: function(e) {
        var id = e.target.dataset.id;
        if (id === this.data.selectCategoryId) {
            this.setData({
                scrolltop: 0,
            })
        } else {
            this.setSelectId(id);
        }
    },
    selectGoodsTap:function(e){
        let goodsId = e.currentTarget.dataset.id;
        for (var i = 0; i < this.data.goodsList.length; i++) {
            let goods = this.data.goodsList[i];
            if (goods.goodsId == goodsId) {
                let selectGoods = {
                    "goodsId"       :goods.goodsId,
                    "status"        :goods.status,
                    "numberStore"   :goods.numberStore,
                    "name"          :goods.name,
                    "pic"           :goods.pic,
                    "sellPrice"     :goods.sellPrice,
                    "mark"          :goods.mark,
                    "order"         :goods.order,
                    "categoryId"    :goods.categoryId,
                };
                this.data._selectGoods = goods;
                this.setCategoryId(goods.categoryId)

                let startDate = TOOLS.time2datestr(goods.startTime)
                let endDate = TOOLS.time2datestr(goods.endTime)
                this.setData({
                    selectGoods: selectGoods,
                    startDate: startDate,
                    endDate: endDate
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

    watchInputStartDate(e) {
        this.data.startDate = e.detail.value
    },

    watchInputEndDate(e) {
        this.data.endDate = e.detail.value
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
    statusChange(e) {
        let selectGoods = this.data.selectGoods;
        if (selectGoods.status == 0) {
            wx.showToast({title: '此状态不能修改', icon: 'none'})
            return;
        } else if (selectGoods.status == 1) {
            selectGoods.status = 2;
        } else if (selectGoods.status == 2) {
            if (this.data._selectGoods.status == 1) {
                selectGoods.status = 1;
            }else{
                selectGoods.status = 3;
            }
        } else if (selectGoods.status == 3) {
            selectGoods.status = 2;
        }
        this.setData({
            selectGoods: selectGoods,
        });
    },
    markChange(){
        let selectGoods = this.data.selectGoods;
        if (selectGoods.mark == 0) {
            selectGoods.mark = 1;
        } else if (selectGoods.mark == 1) {
            selectGoods.mark = 2;
        } else if (selectGoods.mark == 2) {
            selectGoods.mark = 0;
        }
        this.setData({
            selectGoods: selectGoods,
        });
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
            let startTime = 0
            try{
                startTime = TOOLS.datestr2time(this.data.startDate)
            } catch(err) {
                wx.showModal({title: '温馨提示', content: '限时秒杀开始时间错误', showCancel: false})
                return;
            }
            if (startTime == 0 || isNaN(startTime)) {
                wx.showModal({title: '温馨提示', content: '限时秒杀开始时间错误', showCancel: false})
                return
            }

            let endTime = 0
            try{
                endTime = TOOLS.datestr2time(this.data.endDate)
            } catch(err) {
                wx.showModal({title: '温馨提示', content: '限时秒杀结束时间错误', showCancel: false})
                return;
            }
            if (endTime == 0  || isNaN(endTime)) {
                wx.showModal({title: '温馨提示', content: '限时秒杀结束时间错误', showCancel: false})
                return
            }
            if (startTime > endTime) {
                wx.showModal({title: '温馨提示',content: '开始时间比结束时间大',showCancel: false})
                return
            }
            selectGoods.startTime = startTime
            selectGoods.endTime = endTime
        }
        this.data.cache = {};
        let that = this;
        WXAPI.gmGoodsUpdateInfo(selectGoods).then(function(res) {
            if (res.code == 0) {
                wx.showToast({title: '更新成功', icon: 'none'});
            }else{
                wx.showToast({title: res.msg, icon: 'none'});
            }
            that.onShow(that.data.selectCategoryId);
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
        this.data.cache = {};
        let that = this;
        WXAPI.gmGoodsRemove(selectGoods.goodsId).then(function(res) {
            if (res.code == 0) {
                wx.showToast({title: '移除成功', icon: 'none'});
            }else{
                wx.showToast({title: res.msg, icon: 'none'});
            }
            that.onShow(that.data.selectCategoryId);
        });
    },

    bindinput(e) {
        this.setData({
            inputVal: e.detail.value
        })
    },
    bindconfirm(e) {
        this.setData({
            inputVal: e.detail.value
        })
        wx.navigateTo({
            url: '/pages/goods/list?name=' + this.data.inputVal,
        })
    },
    onShow(selectId) {
        this.data.cache = {};
        let that = this;
        AUTH.checkHasLogined().then(isLogined => {
            if (isLogined) {
                this.setData({
                    wxlogin: isLogined
                })
            }
            that.setSelectId(selectId || 0);
        })
    },
    
})