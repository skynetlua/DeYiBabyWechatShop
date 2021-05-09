const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')
const CONFIG = require('../../config.js')

Page({
    data: {
        // initSuccess: false
        waitCount: 0
    },
    onLoad: function(options) {
        wx.setNavigationBarTitle({title: "Q-Baby母婴生活馆"}) 
        if (options) {
            this.data.route = options.route
        }
        this.waitCount = 0
        const accountInfo = wx.getAccountInfoSync();
        let appId = accountInfo.miniProgram.appId;
        wx.setStorageSync('wxAppid', appId);
        // 处理参数
        // {
        //     scene: 'qrcode-goods,?id=' + basicInfo.id + ',' + wx.getStorageSync('uid'),
        //     page: 'pages/start/loading',
        //     is_hyaline: true,
        //     autoColor: true,
        //     expireHours: 1
        // }
        if (options && options.scene) {
            const scene = decodeURIComponent(options.scene) // 处理扫码进来的业务逻辑,格式为：   qrcode-*,queryString,referrer
            const _const = scene.split(',')
            if (_const[0] == 'qrcode-index') {
                this.data.route = '/pages/start/start' + _const[1]
                if (_const.length > 2) {
                    wx.setStorageSync('referrer', _const[2])
                }
            }
            if (_const[0] == 'qrcode-goods') {
                this.data.route = '/pages/goods-details/index' + _const[1]
                if (_const.length > 2) {
                    wx.setStorageSync('referrer', _const[2])
                }
            }
        }
        console.log("start/loading onLoad===>");
        if (options && options.inviter_id) {
            wx.setStorageSync('referrer', options.inviter_id)
            if (options.shareTicket) {
                this.processShareTicket(options.inviter_id, options.shareTicket)
                return
            // } else {
                // this.routePage()
            }
        // } else {
            // this.routePage()
        }
        this.getServerInfo(appId)
    },
    getServerInfo(appid) {
        console.log("start/loading getServerInfo===> appid =", appid)
        const _this = this
        const version = CONFIG.version
        WXAPI.request('/subdomain/appid/wxapp', false, 'get', {appid, version}, function(url){
            console.log("url =", url)
            wx.showModal({
                title: '错误提示',
                content: '连接不上服务器，错误代码:1',
                showCancel: false,
                confirmText: '重试',
                success(res) {
                    _this.getServerInfo(appid)
                }
            })
        }).then(function(res){
            console.log("start/loading getServerInfo===> 1");
            if (res.code != 0) {
                console.log("start/loading getServerInfo===> 2");
                wx.showModal({
                    title: '错误提示',
                    content: res.msg,
                    showCancel: false,
                    confirmText: '重试',
                    success(res) {
                        _this.getServerInfo(appid)
                    }
                })
            } else {
                console.log("start/loading getServerInfo===> 3");
                const data = res.data;
                WXAPI.init2(data.host, data.subdomain);
                for (var i = 0; i < data.config.length; i++) {
                    let item = data.config[i];
                    wx.setStorageSync(item.key, item.value);
                }
                wx.setNavigationBarTitle({
                    title: wx.getStorageSync('mallName')
                })

                _this.login()
                _this.routePage()
                // wx.setStorageSync('subDomain', data.subdomain)
            }
        });
    },
    async login() {
        this.waitCount = 0
        // 自动登录
        const isLogined = await AUTH.checkHasLogined()
        if (!isLogined) {
            console.log("start/loading  未登陆")
            AUTH.login();
        } else {
            console.log("start/loading  已登陆")
        }
    },
    async routePage() {
        this.waitCount++
        if (this.waitCount>30) {
             this.login()
        }
        const isLogined = AUTH.isLogined()
        if (!isLogined) {
            console.log("start/loading  未登陆 waitCount =", this.waitCount)
            setTimeout(() => {
                this.routePage()
            }, 500)
            return
        }
        console.log("start/loading  已登陆")
        // if (!this.data.initSuccess) {
        //     console.log("start/loading  routePage===>no pass")
        //     setTimeout(() => {
        //         this.routePage()
        //     }, 500);
        //     return
        // }
        // console.log("start/loading  routePage======>> pass")
        // 页面跳转,  参数请用 url 编码
        // let pageUrl = this.data.route ? this.data.route : '/pages/start/start'
        // wx.reLaunch({
        //     url: decodeURIComponent(pageUrl)
        // })
        let pageUrl = this.data.route
        if (pageUrl && pageUrl.length > 2) {
            wx.reLaunch({
                url: decodeURIComponent(pageUrl)
            })
            return;
        }
        const version = wx.getStorageSync('app_show_pic_version')
        if (version && version == CONFIG.version) {
            if (CONFIG.shopMod) {
                wx.reLaunch({url: '/pages/shop/select'})
            } else {
                wx.reLaunch({url: '/pages/index/index'})
            }
        }else{
            wx.reLaunch({url: '/pages/start/start'})
        }
    },
    async processShareTicket(inviter_id, shareTicket) {
        const isLogined = AUTH.isLogined()
        if (!AUTH.isLogined()) {
            setTimeout(() => {
                this.processShareTicket(inviter_id, shareTicket)
            }, 500)
            return
        }
        // 处理分享到群奖励
        const code = AUTH.wxaCode()
        wx.getShareInfo({
            shareTicket: shareTicket,
            success: res => {
                WXAPI.shareGroupGetScore(
                    code,
                    inviter_id,
                    res.encryptedData,
                    res.iv
                ).then(_res => {
                    console.log(_res)
                    this.routePage()
                }).catch(err => {
                    console.error(err)
                    this.routePage()
                })
            },
            fail: () => {
                this.routePage()
            }
        })
    }
})