const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')
const IMAGE = require('../../utils/image')

var sliderWidth = 96;
Page({

    /**
     * 页面的初始数据
     */
    data: {
        categoryList: [],
        selectItem: null,
        pressImage:false,
        pickerIdx:0,
        pickerValues: []
    },

    onPickerChange: function(e) {
        let pickerIdx = e.detail.value
        this.setData({
            pickerIdx: pickerIdx
        })
        let selectItem = this.data.selectItem;
        if (selectItem) {
            let icon = this.data.pickerValues[pickerIdx]
            if (icon && icon.length > 0) {
                selectItem.icon = icon;
            }
        }
    },

    onLoad(options) {
        const _this = this;
        wx.getSystemInfo({
            success (res) {
                _this.setData({
                    windowWidth: res.windowWidth
                });
            }
        });
    },

    createCategoryTap(){
        let item = {
            name:"新条目",
            level:1,
            order:1,
            status:1,
            icon:"",
        };
        const that = this;
        WXAPI.gmCategoryUpdate(item).then(res => {
            if (res.code == 0) {
                that.onShow();
            }
        });
    },

    watchTextInput(e) {
        let selectItem = this.data.selectItem;
        const field = e.currentTarget.dataset.field;
        let fieldValue = e.detail.value;
        selectItem[field] = fieldValue;
        this.setData({selectItem:selectItem});
    },

    selectCategoryTap(e){
        let id = e.currentTarget.dataset.id;
        id = Number(id);
        let selectItem;
        let categoryList = this.data.categoryList;
        for (var i = 0; i < categoryList.length; i++) {
            let category = categoryList[i];
            if (category.id == id) {
                selectItem = category;
                break;
            }
        }
        if (!selectItem) {
            wx.showToast({title: '未选中目录', icon: 'none'});
            return;
        }
        let findIdx = selectItem.icon.indexOf("/api/category/")
        let pickerIdx = 0
        if (findIdx >= 0) {
            let value = selectItem.icon.substr(findIdx)
            for (var i = 0; i < this.data.pickerValues.length; i++) {
                let pValue = this.data.pickerValues[i]
                if (pValue == value) {
                    pickerIdx = i
                    break
                }
            }
        } else {
            selectItem.icon = this.data.pickerValues[pickerIdx] || selectItem.icon
        }

        let category = {
            "id"     :selectItem.id,
            "name"   :selectItem.name,
            "status" :selectItem.status,
            "icon"   :selectItem.icon,
            "level"  :selectItem.level,
            "order"  :selectItem.order,
        }
        category.icon = selectItem.publicIcon?selectItem.publicIcon:selectItem.icon;

        this.data._selectItem = category;
        this.setData({
            pickerIdx: pickerIdx,
            selectItem: category
        });
    },

    closePopup() {
        let selectItem = this.data.selectItem;
        let keys = Object.keys(selectItem);
        for (var i = 0; i < keys.length; i++) {
            let key = keys[i];
            if (selectItem[key] != this.data._selectItem[key]) {
                const that = this;
                wx.showModal({
                    title: '确定取消编辑？',
                    content: '',
                    success: function(res) {
                        if (res.confirm) {
                            that.setData({
                                selectItem: null
                            });
                        }
                    }
                });
                return;
            }
        }

        this.setData({
            selectItem: null
        });
    },

    numMinusTap(e) {
        let selectItem = this.data.selectItem;
        const field = e.currentTarget.dataset.field;
        let fieldValue = selectItem[field];
        if (fieldValue <= 0) {
            fieldValue = 0;
        }else{
            fieldValue--;
        }
        selectItem[field] = fieldValue;
        this.setData({
            selectItem: selectItem,
        });
    },

    numPlusTap(e) {
        let selectItem = this.data.selectItem;
        const field = e.currentTarget.dataset.field;
        let fieldValue = selectItem[field];
        fieldValue++;
        selectItem[field] = fieldValue;
        this.setData({
            selectItem: selectItem,
        });
    },
    statusChangeTap(){
        let selectItem = this.data.selectItem;
        if (!selectItem.status) {
            selectItem.status = 1;
        }else{
            selectItem.status = 0;
        }
        this.setData({
            selectItem: selectItem,
        });
    },

    /**
     * 生命周期函数--监听页面初次渲染完成
     */
    onReady: function() {
    },

    requestCategoryList() {
        const that = this;
        WXAPI.gmCategoryList().then(res => {
            if (res.code == 0) {
                let categoryList = res.data.categorys || [];
                let pickerValues = res.data.resIds || [];
                categoryList.sort(function(a, b){
                    return a.order - b.order;
                });
                that.setData({
                    categoryList: categoryList,
                    pickerValues: pickerValues
                });
            }
        });
    },

    /**
     * 生命周期函数--监听页面显示
     */
    onShow() {
        const that = this;
        AUTH.checkHasLogined().then(isLogined => {
            if (isLogined) {
                that.setData({
                    wxlogin: isLogined
                });
                that.requestCategoryList();
            }
        });
    },
    removeIconTap(){
        let selectItem = this.data.selectItem;
        selectItem.icon = "";
        this.setData({
            selectItem: selectItem
        });
    },
    chooseImage(e) {
        const that = this;
        let selectItem = this.data.selectItem;
        wx.chooseImage({
            sizeType: ['original', 'compressed'],
            sourceType: ['album', 'camera'],
            success: function(res) {
                var canvasId = "pressCanvas";
                var drawWidth = that.data.windowWidth;
                drawWidth = 360;
                for (var i = 0; i < res.tempFilePaths.length; i++) {
                    var imagePath = res.tempFilePaths[i];
                    that.setData({pressImage:true});
                    IMAGE.getLessLimitSizeImage(canvasId, imagePath, 300, drawWidth, function(_imagePath){
                        that.setData({pressImage:false});
                        selectItem.icon = _imagePath;
                        that.setData({
                            selectItem: selectItem
                        });
                    });
                }
            }
        });
    },
    requestRemove(){
        let selectItem = this.data.selectItem;
        this.setData({
            selectItem: null
        });
        if (!selectItem) {
            wx.showToast({title: '未选择目录', icon: 'none'})
            return;
        }
        let that = this;
        WXAPI.gmCategoryRemove(selectItem.id).then(function(res) {
            if (res.code == 0) {
                wx.showToast({title: '移除成功', icon: 'none'});
            }else{
                wx.showToast({title: res.msg, icon: 'none'});
            }
            that.onShow();
        });
    },

    requestUpdate(){
        let selectItem = this.data.selectItem;
        this.setData({
            selectItem: null
        });
        if (!selectItem) {
            wx.showToast({title: '未选择目录', icon: 'none'})
            return;
        }
        let findIdx = selectItem.icon.indexOf("/api/category/")
        if (findIdx < 0) {
            delete selectItem.icon
        }
        let that = this
        WXAPI.gmCategoryUpdate(selectItem).then(function(res) {
            if (res.code == 0) {
                wx.showToast({title: '更新成功', icon: 'none'});
            }else{
                wx.showToast({title: res.msg, icon: 'none'});
            }
            that.onShow();
        });
    },

    uploadFile(filePath, categoryId) {
        let that = this;
        const res = WXAPI.gmUploadCategory(filePath, categoryId).then(function(res) {
            console.log("uploadFile res =", res);
            if (res.code == 0) {
                let resData = res.data;
                if (resData.categoryId != categoryId) {
                    console.error("uploadFile======>>error");
                }
                let url = resData.url;
                that.data.selectItem.icon = url;
                console.log("uploadFile ===>>完成上传");
                that.saveTap();
            }
        });
    },

    checkUploadFile(filePath){
        let findIdx = filePath.indexOf("/api/");
        if (findIdx < 0) {
            return true;
        }
    },

    saveTap(){
        let selectItem = this.data.selectItem;
        let icon = selectItem.icon;
        if (!icon || icon == "") {
            wx.showToast({title: "请设置目录图标",icon: 'none'});
            return;
        }
        if(this.checkUploadFile(icon)){
            this.uploadFile(icon, selectItem.id);
            wx.showToast({title: "上传图片中",icon: 'none'});
            return;
        }
        wx.showToast({title: "上传商品资料",icon: 'none'});
        this.requestUpdate();
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