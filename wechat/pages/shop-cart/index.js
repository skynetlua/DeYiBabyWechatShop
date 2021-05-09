const WXAPI = require('apifm-wxapi')
const TOOLS = require('../../utils/tools.js')
const AUTH = require('../../utils/auth')

Page({
    data: {
        wxlogin: true,
        saveHidden: true,
        allSelect: true,
        noSelect: false,
        amountGoods: 0,
        delBtnWidth: 120,
        isInit: false
    },
    //获取元素自适应后的实际宽度
    getEleWidth(w) {
        var real = 0;
        try {
            var windowWidth = wx.getSystemInfoSync().windowWidth
            var scale = 750/w
            real = Math.floor(windowWidth / scale);
            return real;
        } catch (e) {
            return w;
        }
    },
    initEleWidth() {
        var delBtnWidth = this.getEleWidth(this.data.delBtnWidth);
        this.setData({
            delBtnWidth: delBtnWidth
        });
    },
    onLoad() {
        this.initEleWidth();
    },
    onShow() {
        AUTH.checkHasLogined().then(isLogined => {
            this.setData({
                wxlogin: isLogined
            })
            if (isLogined) {
                this.requestCartList()
            }
        })
    },
    async requestCartList() {
        const res = await WXAPI.cartList()
        if (res.code == 0) {
            const items = res.data.items
            let amountGoods = res.data.amountGoods || 0
            if (items && items.length > 0) {
                wx.setTabBarBadge({index: 2, text: `${items.length}`});
            }else{
                wx.removeTabBarBadge({index: 2});
            }
            this.setData({
                items: items,
                amountGoods: amountGoods
            })
        } else {
            this.setData({
                items: null,
                amountGoods: 0
            })
            wx.removeTabBarBadge({index: 2});
        }
    },
    toIndexPage() {
        wx.switchTab({
            url: "/pages/index/index"
        });
    },
    touchS(e) {
        if (e.touches.length == 1) {
            this.setData({
                startX: e.touches[0].clientX
            });
        }
    },
    touchM(e) {
        const index = e.currentTarget.dataset.index;
        if (e.touches.length == 1) {
            var moveX = e.touches[0].clientX;
            var disX = this.data.startX - moveX;
            var delBtnWidth = this.data.delBtnWidth;
            var left = "";
            if (disX == 0 || disX < 0) { //如果移动距离小于等于0，container位置不变
                left = "margin-left:0px";
            } else if (disX > 0) { //移动距离大于0，container left值等于手指移动距离
                left = "margin-left:-" + disX + "px";
                if (disX >= delBtnWidth) {
                    left = "left:-" + delBtnWidth + "px";
                }
            }
            this.data.items[index].left = left
            this.setData({
                items: this.data.items
            })
        }
    },
    touchE(e) {
        var index = e.currentTarget.dataset.index;
        if (e.changedTouches.length == 1) {
            var endX = e.changedTouches[0].clientX;
            var disX = this.data.startX - endX;
            var delBtnWidth = this.data.delBtnWidth;
            var left = disX > delBtnWidth / 2 ? "margin-left:-" + delBtnWidth + "px" : "margin-left:0px";
            this.data.items[index].left = left
            this.setData({
                items: this.data.items
            })
        }
    },
    toDetailsTap(e) {
        wx.navigateTo({url: "/pages/goods-details/index?id=" + e.currentTarget.dataset.id})
    },
    async delItem(e) {
        const id = e.currentTarget.dataset.id
        wx.showModal({
            content: '确定要删除该商品吗？',
            success: (res) => {
                if (res.confirm) {
                    this.delItemDone(id)
                }
            }
        })
    },
    async delItemDone(id) {
        const res = await WXAPI.cartRemove(id)
        if (res.code != 0) {
            wx.showToast({
                title: res.msg,
                icon: 'none'
            })
            setTimeout(() => {
                this.requestCartList()
            }, 1000)
        } else {
            this.requestCartList()
        }
    },
    setItemNum(itemId, itemNum){
        var self = this;
        WXAPI.cartModifyNumber(itemId, itemNum).then(res => {
            if (res.code != 0) {
                wx.showToast({
                    title: res.msg,
                    icon: 'none'
                })
            } else {
                this.requestCartList()
            }
        })
    },
    async jiaBtnTap(e) {
        const index = e.currentTarget.dataset.index;
        const item = this.data.items[index];
        const num = item.numberBuy + 1;
        this.setItemNum(item.id, num);
    },
    async jianBtnTap(e) {
        const index = e.currentTarget.dataset.index;
        const item = this.data.items[index]
        const num = item.numberBuy - 1
        if (num <= 0) {
            // 弹出删除确认
            wx.showModal({
                content: '确定要删除该商品吗？',
                success: (res) => {
                    if (res.confirm) {
                        this.delItemDone(item.id)
                    }
                }
            })
            return
        }
        this.setItemNum(item.id, num);
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
    }
})