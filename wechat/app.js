
// let __CacheMap = {};

App({
    onLaunch: function() {
        const that = this;
        // 检测新版本
        const updateManager = wx.getUpdateManager();
        updateManager.onUpdateReady(function() {
            wx.showModal({
                title: '更新提示',
                content: '新版本已经准备好，是否重启应用？',
                success(res) {
                    if (res.confirm) {
                        // 新的版本已经下载好，调用 applyUpdate 应用新版本并重启
                        updateManager.applyUpdate()
                    }
                }
            });
        });
        /**
         * 初次加载判断网络情况
         * 无网络状态下根据实际情况进行调整
         */
        wx.getNetworkType({
            success(res) {
                const networkType = res.networkType
                if (networkType === 'none') {
                    wx.removeStorageSync('isConnected');
                    // wx.removeStorageSync('loginTime');
                    wx.showToast({
                        title: '当前无网络',
                        icon: 'loading',
                        duration: 2000
                    })
                }
            }
        });
        /**
         * 监听网络状态变化
         * 可根据业务需求进行调整
         */
        wx.onNetworkStatusChange(function(res) {
            if (!res.isConnected) {
                wx.removeStorageSync('isConnected');
                // wx.removeStorageSync('loginTime');
                wx.showToast({
                    title: '网络已断开',
                    icon: 'loading',
                    duration: 2000
                })
            } else {
                wx.setStorageSync('isConnected', true)
                wx.hideToast()
            }
        });
        wx.setStorageSync('isConnected', true)
        // wx.removeStorageSync('loginTime');
        wx.setStorageSync('gm', false)

        // let menuButtonObject = wx.getMenuButtonBoundingClientRect();
        // console.log("小程序胶囊信息",menuButtonObject)
        // wx.getSystemInfo({
        //   success: res => {
        //     let statusBarHeight = res.statusBarHeight,
        //       navTop = menuButtonObject.top,//胶囊按钮与顶部的距离
        //       navHeight = statusBarHeight + menuButtonObject.height + (menuButtonObject.top - statusBarHeight)*2;//导航高度
        //     this.globalData.navHeight = navHeight;
        //     this.globalData.navTop = navTop;
        //     this.globalData.windowHeight = res.windowHeight;
        //     this.globalData.menuButtonObject = menuButtonObject;
        //     console.log("navHeight",navHeight);
        //   },
        //   fail(err) {
        //     console.log(err);
        //   }
        // })
    },
    // globalData: {
    // }
})


// {
//     "pagePath": "pages/order-list/index",
//     "iconPath": "images/nav/order-off.png",
//     "selectedIconPath": "images/nav/order-on.png",
//     "text": "订单"
// },