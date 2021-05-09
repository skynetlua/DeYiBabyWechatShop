const WXAPI = require('apifm-wxapi')
const IMAGE = require('../../utils/image')


Page({
    data: {
        orderId: 1,
        amount: 0.00,
        type: 0,
        typeItems: [
            {name: '我要退款(无需退货)', value: '0', checked: true},
            {name: '我要退货退款', value: '1'},
        ],
        logisticsStatus: 0,
        logisticsStatusItems: [
            { name: '未收到货', value: '0', checked: true },
            { name: '已收到货', value: '1' }
        ],

        reasons: [
            "不喜欢/不想要",
            "空包裹",
            "未按约定时间发货",
            "快递/物流一直未送达",
            "货物破损已拒签",
            "退运费",
            "规格尺寸与商品页面描述不符",
            "功能/效果不符",
            "质量问题",
            "少件/漏发",
            "包装/商品破损",
            "发票问题",
        ],
        reasonIndex: 0,

        files: [],
        pics: [],

        windowWidth: 0,
    },
    onLoad: function(e) {
        this.setData({
            orderId: parseInt(e.id)
            // amount: Number(e.amount)
        });
        const _this = this
        wx.getSystemInfo({
            success (res) {
                _this.setData({
                    windowWidth: res.windowWidth
                });
            }
        })
    },
    onShow() {
        const _this = this
        WXAPI.refundApplyDetail(this.data.orderId).then(res => {
            if (res.code == 0) {
                let refund = res.data.refund;
                _this.setData({
                    refund: refund,
                    status: refund.refundStatus,
                    amount: refund.amount,
                });
                // let data = _this.data;
                // if (data.amount !== data.orderInfo.amountReal) {
                //     console.log("refundApplyDetail data.amount =", data.amount);
                //     console.log("refundApplyDetail orderInfo.amountReal =", data.orderInfo.amountReal);
                //     console.error("refundApplyDetail 金额不相等");
                // }
            }
        })
    },
    refundApplyCancel() {
        WXAPI.refundApplyCancel(this.data.orderId).then(res => {
            if (res.code == 0) {
                // wx.switchTab({url: "/pages/order-list/index"})
                wx.navigateBack({})
            }
        })
    },
    typeItemsChange: function(e) {
        const typeItems = this.data.typeItems;
        const curValue = e.detail.value;
        for (var i = 0, len = typeItems.length; i < len; ++i) {
            typeItems[i].checked = typeItems[i].value == curValue;
        }
        this.setData({
            typeItems: typeItems,
            type: curValue
        });
    },
    logisticsStatusItemsChange: function(e) {
        const items = this.data.logisticsStatusItems;
        const curValue = e.detail.value;
        for (var i = 0, len = items.length; i < len; ++i) {
            items[i].checked = items[i].value == curValue;
        }
        this.setData({
            logisticsStatusItems: items,
            logisticsStatus: curValue
        });
    },
    reasonChange: function(e) {
        this.setData({
            reasonIndex: e.detail.value
        })
    },
    chooseImage: function(e) {
        const that = this;
        wx.chooseImage({
            sizeType: ['original', 'compressed'],
            sourceType: ['album', 'camera'],
            success: function(res) {
                var canvasId = "pressCanvas";
                var drawWidth = that.data.windowWidth;
                for (var i = 0; i < res.tempFilePaths.length; i++) {
                    var imagePath = res.tempFilePaths[i];
                    IMAGE.getLessLimitSizeImage(canvasId, imagePath, 300, drawWidth, function(_imagePath){
                        that.data.files.push(_imagePath)
                        that.setData({
                            files: that.data.files
                        });
                    });
                }
                // that.setData({
                //     files: that.data.files.concat(res.tempFilePaths)
                // });
            }
        })
    },
    previewImage: function(e) {
        wx.previewImage({
            current: e.currentTarget.id,
            urls: this.data.files
        })
    },
    async uploadPics() {
        for (let i = 0; i < this.data.files.length; i++) {
            const res = await WXAPI.uploadFile(wx.getStorageSync('token'), this.data.files[i])
            if (res.code == 0) {
                this.data.pics.push(res.data.url)
            }
        }
    },
    async bindSave(e) {
        // let amount = e.detail.value.amount;
        let remark = e.detail.value.remark || "";
        if (this.data.files.length>0) {
            await this.uploadPics();
        }
        let data = this.data;
        WXAPI.refundApply({
            orderId:      data.orderId,
            refundType:   data.type,
            logistics:    data.logisticsStatus,
            reasonId:    data.reasonIndex,
            // amount:             amount,
            remark:       remark,
            pics:          data.pics.join(",")
        }).then(res => {
            if (res.code == 0) {
                wx.showModal({
                    title: '成功',
                    content: '提交成功，请耐心等待我们处理！',
                    showCancel: false,
                    confirmText: '我知道了',
                    success(res) {
                        // wx.switchTab({url: "/pages/order-list/index"})
                        wx.navigateBack({})
                    }
                })
            } else {
                wx.showModal({
                    title: '失败',
                    content: res.msg,
                    showCancel: false,
                    confirmText: '我知道了',
                    success(res) {
                        // wx.switchTab({url: "/pages/order-list/index"})
                        wx.navigateBack({})
                    }
                })
            }
        })
    }
});