const WXAPI = require('apifm-wxapi')
const CONFIG = require('../../config.js')
const AUTH = require('../../utils/auth')
const IMAGE = require('../../utils/image')
const TOOLS = require('../../utils/tools.js')

Page({
    data: {
        wxlogin: true,

        name:"",
        // icon:"",

        barCode:"",
        desc:"",
        skus:[],
        newSkuIdx: "",
        newSku:"",
        newSkuPrice: "",
        pics:[],
        contents:[],

        originPrice:0.0,
        sellPrice:0.0,
        // minPrice:0.0,
        numberStore:9999,
        numberSell:0,

        windowWidth: 0,
        pressImage:false,

        categorys: [],
        categoryIds:[],
        categoryIndex: 0,

        statusNames: ["录入", "待上架", "上架", "下架"],
        statuss:[0, 1, 2, 3],
        statusIndex: 0,
        isInit: false,

        markNames: ['无', '推荐', '秒杀', '拼团'],

        startDate: "2021-03-12",
        endDate: "2021-03-12",
        startTime: "00:01",
        endTime: "23:59",

        fromDatePattern: "",
        toDatePattern: "2031-03-12",
        fromTimePattern: "00:01",
        toTimePattern: "23:59"

        // pickerIdx:0,
        // pickerValues: []
    },
    requestCategoryList() {
        const that = this;
        WXAPI.gmCategoryList().then(res => {
            if (res.code == 0) {
                let categoryList = res.data.categorys || [];
                categoryList.sort(function(a, b){
                    return a.order - b.order;
                });
                let categorys = ["选择目录"];
                let categoryIds = [0];
                for (var i = 0; i < categoryList.length; i++) {
                    let item = categoryList[i];
                    categorys.push(item.name);
                    categoryIds.push(item.id);
                }
                that.setData({
                    categorys: categorys,
                    categoryIds: categoryIds
                });
            }
        });
    },
    setStatus:function(status){
        let statusIndex = 0;
        let statuss = this.data.statuss
        for (var i = 0; i < statuss.length; i++) {
            if (statuss[i] == status) {
                statusIndex = i;
                break;
            }
        }
        this.setData({
            statusIndex: statusIndex
        });
    },
    statusChange(e) {
        let that = this;
        wx.showActionSheet({
            itemList: this.data.statusNames,
            success (res) {
                that.setData({
                    status: res.tapIndex,
                });
            },
            fail (res) {
                // console.log(res.errMsg)
            }
        })
    },
    getCategoryId(){
        let categoryIndex = this.data.categoryIndex;
        let categoryId = this.data.categoryIds[categoryIndex];
        return categoryId;
    },
    setCategoryId(categoryId){
        for (var i = 0; i < this.data.categoryIds.length; i++) {
            let id = this.data.categoryIds[i];
            if (id == categoryId) {
                this.setData({
                    categoryIndex: i
                });
                break;
            }
        }
    },
    onLoad(e) {
        const that = this
        wx.getSystemInfo({
            success (res) {
                that.setData({
                    windowWidth: res.windowWidth
                });
            }
        });
        let barCode = e.barcode;
        if (barCode) {
            this.setData({barCode: barCode});
        }
        let goodsId = e.id;
        if (goodsId) {
            this.setData({goodsId: goodsId});
        }
    },
    categoryChange(e) {
        this.setData({
            categoryIndex: e.detail.value
        });
    },
    onShow() {
        if (this.data.isInit) {
            return
        }
        this.data.isInit = true
        const that = this
        AUTH.checkHasLogined().then(isLogined => {
            if (isLogined) {
                that.setData({
                    wxlogin: isLogined
                });
                that.requestCategoryList();
                if (that.data.barCode) {
                    that.requestGoodsData(that.data.barCode);
                } else {
                    that.requestGoodsInfo(that.data.goodsId);
                }
            }
        });
    },
    splitItems(str){
        if (!str || str == "") {
            return [];
        }
        let tmps = str.split(";");
        let items = [];
        for (var i = 0; i < tmps.length; i++) {
            let tmp = tmps[i];
            if (tmp && tmp.length > 0) {
                items.push(tmp);
            }
        }
        return items;
    },
    requestGoodsInfo(goodsId){
        let that = this;
        WXAPI.gmGoodsInfo(goodsId || 0).then(function(res) {
            if (res.code == 0) {
                const resData = res.data.goods;

                var skuGroups = [];
                var skuList = [];
                if (resData.skuJson && resData.skuJson.length > 0) {
                    skuList = JSON.parse(resData.skuJson);
                    // skuList.push({"label":"规格2","id":100,"name":"2NB码96片","price":0});
                    // skuList.push({"label":"规格2","id":200,"name":"2NB码96片","price":0});
                    // skuList.push({"label":"规格3","id":10000,"name":"3NB码96片","price":0});
                    // skuList.push({"label":"规格3","id":20000,"name":"3NB码96片","price":0});
                    let skuGroupMap = {};
                    for (var i = 0; i < skuList.length; i++) {
                        let sku = skuList[i];
                        let skuGroup = skuGroupMap[sku.label];
                        let level = TOOLS.getLevel(sku.id)
                        if (!skuGroup) {
                            skuGroup = {
                                label:sku.label,
                                level:level,
                                skuList:[],
                                skuMap: {}
                            };
                            skuGroupMap[sku.label] = skuGroup;
                            skuGroups.push(skuGroup);
                        } else {
                            if (skuGroup.level != level) {
                                console.log("shippingCartInfo skuGroup.level =", skuGroup.level, "level =", level);
                            }
                        }
                        skuGroup.skuList.push(sku)
                        skuGroup.skuMap[sku.name] = sku
                    }
                }
                resData.skuList = skuList
                resData.skuGroups = skuGroups
                
                resData.pics           = resData.pics || [];
                resData.contents       = resData.contents || [];
                resData.publicPics     = resData.publicPics || [];
                resData.publicContents = resData.publicContents || [];

                if (!resData.id) {
                    that.setData({
                        order:resData.order,
                        numberStore:resData.numberStore,
                        numberSell:resData.numberSell,
                        mainType    :"",
                        subType     :"",
                        pics        :resData.pics,
                        contents    :resData.contents,

                        skuList     :resData.skuList,
                        skuGroups   :resData.skuGroups,

                        publicPics   :resData.publicPics,
                        publicContents:resData.publicContents
                    });
                    if (that.data.barCode && that.data.barCode.length > 0) {
                        wx.showToast({title: "后台有商品，但未上架小商城", icon: 'none'})
                        return
                    }
                    wx.showToast({title: "小商城无该商品", icon: 'none'})
                    return
                }

                let data = {
                    goodsId     :resData.goodsId,
                    name        :resData.name,
                    status      :resData.status,
                    mark        :resData.mark,
                    barCode     :resData.barCode,

                    promote     :resData.promote || "",
                    order       :resData.order,
                    mainType    :resData.mainType,
                    subType     :resData.subType,

                    originPrice :resData.originPrice,
                    sellPrice   :resData.sellPrice,
                    numberStore :resData.numberStore,
                    numberSell  :resData.numberSell,

                    skuList     :resData.skuList,
                    skuGroups   :resData.skuGroups,

                    pics        :resData.pics,
                    contents    :resData.contents,

                    publicPics  :resData.publicPics,
                    publicContents:resData.publicContents,
                };
                that.setData(data);
                that.setCategoryId(resData.categoryId);
            } else {
                wx.showToast({title: res.msg, icon: 'none'});
            }
        });
    },
    getSkuGroup(level) {
        let skuGroups = this.data.skuGroups;
        for (var i = 0; i < skuGroups.length; i++) {
            let skuGroup = skuGroups[i];
            if (skuGroup.level == level) {
                return skuGroup
            }
        }
        wx.showToast({title: "规格组不存在", icon: 'none'});
        this.updateSkuGroups();
    },
    updateSkuGroups() {
        this.setData({
            skuGroups: this.data.skuGroups
        })
    },
    onRemoveSkuGroupTap(e) {
        let level = e.currentTarget.dataset.level;
        let skuGroup = this.getSkuGroup(level);
        if (!skuGroup) {
            return;
        }
        let that = this
        wx.showModal({
            title: '确定要删除规格组['+skuGroup.label+']',
            content: '',
            success: function(res) {
                if (res.confirm) {
                    let skuGroups = that.data.skuGroups;
                    for (var i = 0; i < skuGroups.length; i++) {
                        let skuGroup = skuGroups[i];
                        if (skuGroup.level == level) {
                            skuGroups.splice(i, 1);
                            break;
                        }
                    }
                    that.updateSkuGroups();                
                }
            }
        });
    },
    onAddSkuGroupTap(e) {
        let skuGroups = this.data.skuGroups;
        let newGroupLevel = Number(this.data.newGroupLevel);
        if (isNaN(newGroupLevel) || (newGroupLevel != 1 && newGroupLevel != 2)) {
            wx.showToast({title: "规格组级别只支持1或者2", icon: 'none'});
            return;
        }
        let maxLevel = 0;
        for (var i = 0; i < skuGroups.length; i++) {
            let skuGroup = skuGroups[i];
            if (skuGroup.level > maxLevel) {
                maxLevel = skuGroup.level;
            }
            if (skuGroup.level == newGroupLevel) {
                wx.showToast({title: "规格组已存在级别"+newGroupLevel, icon: 'none'});
                return
            }
        }
        maxLevel++;
        if (newGroupLevel > maxLevel) {
            newGroupLevel = maxLevel;
        }

        let that = this
        wx.showModal({
            title: '确定要增加规格组'+newGroupLevel,
            content: '',
            success: function(res) {
                if (res.confirm) {
                    let skuGroup = {
                        label: "规格"+newGroupLevel,
                        level: newGroupLevel,
                        skuList: [],
                        skuMap: {}
                    };
                    skuGroups.push(skuGroup);
                    that.updateSkuGroups();               
                }
            }
        });
    },
    onAddSkuGroupLevelInput(e) {
        let value = e.detail.value;
        this.setData({
            newGroupLevel: value
        })
    },
    onSkuLabelInput(e) {
        let level = e.currentTarget.dataset.level;
        let value = e.detail.value;
        let skuGroup = this.getSkuGroup(level);
        if (!skuGroup) {
            return;
        }
        skuGroup.label = value;
        this.updateSkuGroups();
    },
    onAddSkuTap(e){
        let level = e.currentTarget.dataset.level;
        let skuGroup = this.getSkuGroup(level);
        if (!skuGroup) {
            return;
        }
        let name = skuGroup.newSku;
        if (!name || name == "") {
            wx.showToast({title: "请先输入规格名称", icon: 'none'});
            return;
        }
        let price = 0;
        if (skuGroup.newSkuPrice && skuGroup.newSkuPrice.length > 0) {
            price = Number(skuGroup.newSkuPrice);
            if (isNaN(price)) {
                price = 0;
            }
        }
        let sku = skuGroup.skuMap[name];
        if (!sku) {
            skuGroup.skuMap[name] = {
                label: skuGroup.label,
                id: 0,
                name: name,
                price: price
            };
            sku = skuGroup.skuMap[name];
            skuGroup.skuList.push(sku)
        } else {
            sku.price = price
        }
        skuGroup.newSku = "";
        skuGroup.newSkuPrice = "";
        this.updateSkuGroups();
    },
    onNewSkuInput(e){
        let level = e.currentTarget.dataset.level;
        let value = e.detail.value;
        let skuGroup = this.getSkuGroup(level);
        if (!skuGroup) {
            return;
        }
        skuGroup.newSku = value;
        this.updateSkuGroups();
    },
    onNewSkuPriceInput(e){
        let level = e.currentTarget.dataset.level;
        let value = e.detail.value;
        let skuGroup = this.getSkuGroup(level);
        if (!skuGroup) {
            return;
        }
        skuGroup.newSkuPrice = value;
        this.updateSkuGroups();
    },
    onRemoveSkuTap(e){
        let level = e.currentTarget.dataset.level;
        let skuId = e.currentTarget.dataset.skuid;
        let skuGroup = this.getSkuGroup(level);
        if (!skuGroup) {
            return;
        }
        let selectSku = null;
        for (var i = 0; i < skuGroup.skuList.length; i++) {
            let sku = skuGroup.skuList[i];
            if (sku.id == skuId) {
                selectSku = sku;
                break;
            }
        }
        if (!selectSku) {
            this.updateSkuGroups();
            return;
        }
        let that = this;
        wx.showModal({
            title: '确定删除规格组['+skuGroup.label+']中的['+selectSku.name+']？',
            content: '',
            success: function(res) {
                if (res.confirm) {
                    for (var i = 0; i < skuGroup.skuList.length; i++) {
                        let sku = skuGroup.skuList[i];
                        if (sku.id == skuId) {
                            delete skuGroup.skuMap[sku.name];
                            skuGroup.skuList.splice(i, 1);
                            that.updateSkuGroups();
                            break;
                        }
                    }
                }
            }
        });
    },
    processSkuGroups() {
        let skuGroups = this.data.skuGroups;
        let skuList = [];
        for (var i = 0; i < skuGroups.length; i++) {
            let skuGroup = skuGroups[i]
            let base = 1;
            let level = skuGroup.level
            while (level > 1) {
                base = base*100;
                level--;
            }
            for (var j = 0; j < skuGroup.skuList.length; ++j) {
                let sku = skuGroup.skuList[j];
                sku.label = skuGroup.label;
                sku.id = base*(j+1);
                skuList.push(sku);
            }
        }
        // console.log('processSkuGroups===>> skuList =', skuList);
        return skuList;
    },
    requestGoodsData(barCode) {
        let that = this;
        WXAPI.gmGoodsGoodsData(barCode).then(function(res) {
            var data = res.data
            if (data.barCode && data.barCode.length > 0) {
                that.setData({
                    goodsId: data.goodsId,
                    barCode: data.barCode,
                    name: data.name,
                    sellPrice: data.sellPrice,
                    mainType:data.mainType,
                    subType:data.subType,
                    numberStore: 9999
                });
                that.requestGoodsInfo(data.goodsId);
            } else {
                wx.showToast({
                    title: "后台无该商品数据",
                    icon: 'none'
                });
            }
        })
    },
    scanOrderCode() {
        const that = this;
        wx.scanCode({
            onlyFromCamera: false,
            success(res) {
                console.log("wx.scanCode success res =", res)
                if (res.result) {
                    that.setData({
                        barCode: res.result
                    });
                    that.requestGoodsData(res.result)
                } else {
                    that.setData({
                        barCode: ""
                    });
                }
            },
            fail(err) {
                console.log("wx.scanCode fail err =", err);
                wx.showToast({
                    title: "扫码取消",
                    icon: 'none'
                });
            }
        })
    },
    removeTap(e) {
        const url = e.currentTarget.dataset.url;
        const type = e.currentTarget.dataset.type;
        const that = this;
        wx.showModal({
            title: '确定删除该照片？',
            content: '',
            success: function(res) {
                if (res.confirm) {
                    if (type == "pic") {
                        let pics = that.data.pics;
                        let idx = -1;
                        for (let i = 0; i < pics.length; i++) {
                            let pic = pics[i];
                            if (url == pic) {
                                idx = i;
                                break;
                            }
                        }
                        let publicPics = that.data.publicPics
                        if (idx >= 0) {
                            pics.splice(idx, 1);
                            publicPics.splice(idx, 1);
                        }
                        that.setData({
                            pics: pics,
                            publicPics:publicPics
                        });
                    } else if (type == "content") {
                        let contents = that.data.contents;
                        let idx = -1;
                        for (let i = 0; i < contents.length; i++) {
                            let pic = contents[i];
                            if (url == pic) {
                                idx = i;
                                break;
                            }
                        }
                        let publicContents = that.data.publicContents;
                        if (idx >= 0) {
                            contents.splice(idx, 1);
                            publicContents.splice(idx, 1);
                        }
                        that.setData({
                            contents: contents,
                            publicContents:publicContents,
                        });
                    } else {
                        that.setData({
                            publicPic: null,
                            icon: ""
                        });
                    }
                }
            }
        });
    },
    chooseImage(e) {
        const that = this;
        const type = e.currentTarget.dataset.type;
        wx.chooseImage({
            sizeType: ['original', 'compressed'],
            sourceType: ['album', 'camera'],
            success: function(res) {
                var canvasId = "pressCanvas";
                var drawWidth = that.data.windowWidth;
                for (var i = 0; i < res.tempFilePaths.length; i++) {
                    var imagePath = res.tempFilePaths[i];
                    that.setData({
                        pressImage:true,
                    })
                    IMAGE.getLessLimitSizeImage(canvasId, imagePath, 300, drawWidth, function(_imagePath){
                        that.setData({
                            pressImage:false,
                        })
                        if (type == "pic") {
                            let publicPics = that.data.publicPics;
                            publicPics.push(_imagePath);
                            that.data.pics.push(_imagePath);
                            that.setData({
                                pics: that.data.pics,
                                publicPics:publicPics,
                            });
                        } else if (type == "content") {
                            let publicContents = that.data.publicContents;
                            publicContents.push(_imagePath);
                            that.data.contents.push(_imagePath);
                            that.setData({
                                contents: that.data.contents,
                                publicContents:publicContents,
                            });
                        } else {
                            that.setData({
                                icon: _imagePath
                            });
                        }
                    });
                }
            }
        })
    },
    watchTextInput(e) {
        const field = e.currentTarget.dataset.field;
        let fieldValue = e.detail.value;
        this.data[field] = fieldValue;
        let data = {};
        data[field] = fieldValue;
        this.setData(data);
    },
    numJianTap(e) {
        const field = e.currentTarget.dataset.field;
        let fieldValue = this.data[field];
        if (fieldValue <= 0) {
            fieldValue = 0;
        }else{
            fieldValue--;
        }
        let data = {};
        data[field] = fieldValue;
        this.setData(data);
    },
    numJiaTap(e) {
        const field = e.currentTarget.dataset.field;
        let fieldValue = this.data[field];
        fieldValue++;
        let data = {};
        data[field] = fieldValue;
        this.setData(data);
    },
    watchInput(e) {
        const field = e.currentTarget.dataset.field;
        let fieldValue = Number(e.detail.value);
        if (!fieldValue && fieldValue != 0) {
            fieldValue = this.data[field];
        }
        this.data[field] = fieldValue;
        let data = {};
        data[field] = fieldValue;
        this.setData(data);
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
    requestUpdate(){
        let data = this.data;
        if (!data.name || data.name == "") {
            wx.showToast({title: "请设置商品名字",icon: 'none'});
            return;
        }
        // if (!data.desc || data.desc == "") {
        //     wx.showToast({title: "请设置商品个性描述",icon: 'none'});
        //     return;
        // }
        // if (!data.enterPrice) {
        //     wx.showToast({title: "请设置商品进货价格",icon: 'none'});
        //     return;
        // }
        // if (!data.minPrice) {
        //     wx.showToast({title: "请设置商品保底价格",icon: 'none'});
        //     return;
        // }
        

        let skuList = this.processSkuGroups()
        let skuJson = JSON.stringify(skuList);

        let categoryId = this.getCategoryId()
        let params = {
            goodsId     :data.goodsId,
            name        :data.name,
            status      :data.status,
            barCode     :data.barCode || '',
            // desc        :data.desc,
            // skus        :data.skus.join(";"),
            // skuPrices   :data.skuPrices.join(";"),
            
            // enterPrice  :data.enterPrice,
            // minPrice    :data.minPrice,
            sellPrice   :data.sellPrice,
            originPrice :data.originPrice,
            numberStore :data.numberStore,
            numberSell  :data.numberSell,
            categoryId  :categoryId,

            promote     :data.promote || '',
            order       :data.order,
            // express     :data.express,
            mainType    :data.mainType,
            subType     :data.subType,
            skuJson     :skuJson,
            // icon        :data.icon,
            pics        :data.pics.join(";"),
            contents    :data.contents.join(";")
        };
        let goodsId = data.goodsId;
        wx.showLoading({title: '上传中'})
        let that = this;
        WXAPI.gmGoodsUpdate(params).then(function(res) {
            wx.hideLoading();
            console.log("requestUpdate res =", res);
            if (res.code == 0) {
                that.requestGoodsInfo(goodsId);
                wx.showToast({title:"商品更新成功", icon: 'none'});
                wx.navigateBack();
                return;
            }
            if (res.msg) {
                wx.showToast({title:res.msg, icon: 'none'})
            }
        });
    },
    uploadFile(task) {
        let goodsId = task.goodsId;
        let filePath = task.filePath;
        let part = task.part;
        let idx = task.idx;
        let categoryId = this.getCategoryId()
        let that = this;
        WXAPI.gmUploadGoods(filePath, categoryId, goodsId, part, idx).then(function(res) {
            console.log("uploadFile res =", res);
            if (res.code == 0) {
                let resData = res.data;
                if (resData.goodsId != goodsId
                    ||resData.categoryId != categoryId
                    ||resData.part != part
                    ||resData.idx != idx) {
                    console.error("uploadFile======>>error");
                    wx.showToast({title: "上传图片，返回参数出错",icon: 'none'});
                    return
                }
                idx = idx-1
                let url = resData.url;
                if (part == "icon") {
                    that.data.icon = url;
                }else if (part == "pic") {
                    that.data.pics[idx] = url;
                }else if (part == "content") {
                    that.data.contents[idx] = url;
                }else{
                    console.error("uploadFile======>>");
                }
                let tasks = that.data.tasks;
                for (var i = 0; i < tasks.length; i++) {
                    if (filePath == tasks[i].filePath) {
                        tasks.splice(i, 1);
                        break;
                    }
                }
                if (tasks.length <= 0) {
                    console.log("uploadFile ===>>完成上传");
                    that.saveTap();
                }
            }
        });
    },
    checkUploadFile(filePath, part, idx){
        let findIdx = filePath.indexOf("/goods/");
        if (findIdx >= 0) {
            return false;
        }
        findIdx = filePath.indexOf("/auto/");
        if (findIdx >= 0) {
            return false;
        }
        findIdx = filePath.indexOf("/picture/");
        if (findIdx >= 0) {
            return false;
        }
        return true;
    },
    saveTap(){
        // let skuGroups = this.data.skuGroups;
        // console.log("skuGroups =", skuGroups);
        // return;
        let catogoryId = this.getCategoryId()
        if (catogoryId <= 0) {
            wx.showToast({title: "请指定主目录",icon: 'none'});
            return
        }

        let tasks = [];
        const goodsId = this.data.goodsId;
        // let icon = this.data.icon;
        // if (icon && icon.length > 0) {
        //     if(this.checkUploadFile(icon)){
        //         let task = {
        //             goodsId :goodsId,
        //             filePath:icon,
        //             part    :"icon",
        //             idx     :0
        //         };
        //         tasks.push(task);
        //     }
        // }
        let pics = this.data.pics;
        if (pics && pics.length > 0) {
            for (var i = 0; i < pics.length; i++) {
                let pic = pics[i];
                if (this.checkUploadFile(pic)) {
                    let task = {
                        goodsId :goodsId,
                        filePath:pic,
                        part    :"pic",
                        idx     :i+1
                    };
                    tasks.push(task);
                }
            }
        }
        let contents = this.data.contents;
        if (contents && contents.length > 0) {
            for (var i = 0; i < contents.length; i++) {
                let content = contents[i];
                if (this.checkUploadFile(content)) {
                    let task = {
                        goodsId :goodsId,
                        filePath:content,
                        part    :"content",
                        idx     :i+1
                    };
                    tasks.push(task);
                }
            }
        }
        // if (!icon || icon == "") {
        //     wx.showToast({title: "请设置商品图标",icon: 'none'});
        //     return;
        // }
        if (!pics || pics.length == 0) {
            wx.showToast({title: "请设置商品图片",icon: 'none'});
            return;
        }
        if (!contents || contents.length == 0) {
            wx.showToast({title: "请设置商品详情图片",icon: 'none'});
            return;
        }
        this.data.tasks = tasks;
        if (tasks.length>0) {
            for (var i = 0; i < tasks.length; i++) {
                this.uploadFile(tasks[i]);
            }
            wx.showToast({title: "上传图片中", icon: 'none'});
            return;
        }
        // wx.showToast({title: "上传商品资料", icon: 'none'});
        this.requestUpdate();
    }
    // previewImage(e) {
    //     const url = e.currentTarget.dataset.url
    //     if(typeof url == "string"){
    //         wx.previewImage({current: url,urls: [url]})
    //     }else{
    //         wx.previewImage({current: url[0],urls: url})
    //     }
    // },
   
})
