const WXAPI = require('apifm-wxapi')

async function checkSession() {
    console.log("auth.checkSession=====>>")
    return new Promise((resolve, reject) => {
        wx.checkSession({
            success() {
                console.log("auth.checkSession=====<<success")
                return resolve(true)
            },
            fail() {
                console.log("auth.checkSession=====<<fail")
                return resolve(false)
            }
        });
    });
}

function onLoginSuccess(data) {
    console.log("auth.onLoginSuccess=====>> data=", data)
    wx.setStorageSync('token', data.token)
    wx.setStorageSync('uid', data.uid)
    if (data.nickName && data.nickName.length > 0) {
        wx.setStorageSync('nickName', data.nickName)
    }

    var curTime = new Date().getTime()
    wx.setStorageSync('loginTime', curTime)
}

function loginOut() {
    console.log("auth.loginOut=====>>")
    wx.removeStorageSync('token')
    wx.removeStorageSync('uid')
}

function isLogined(){
    const token = wx.getStorageSync('token')
    console.log("auth.isLogined=====>> token=", token)
    if (!token) {
        return false
    }
    var loginTime = wx.getStorageSync('loginTime')
    console.log("auth.isLogined=====>> loginTime=", loginTime)
    if (loginTime) {
        var curTime = new Date().getTime()
        if (curTime-loginTime < 1000*60*60*24*30) {
            return true
        }
    }
    return true
    // return false
}

// 检测登录状态，返回 true / false
async function checkHasLogined() {
    const token = wx.getStorageSync('token')
    console.log("auth.checkHasLogined=====>> token =", token)
    if (!token) {
        console.log("auth.checkHasLogined=====>> 本地无token")
        return false
    }
    var loginTime = wx.getStorageSync('loginTime')
    if (loginTime) {
        var curTime = new Date().getTime()
        if (curTime-loginTime<1000*60*60*24) {
            console.log("auth.checkHasLogined=====>>未过10分钟")
            return true
        }
    }
    const loggined = await checkSession()
    console.log("auth.checkSession=====>> loggined =", loggined)
    if (!loggined) {
        console.log("auth.checkHasLogined=====>> 微信session失效")
        return false
    }
    const res = await WXAPI.checkToken(token)
    console.log("auth.checkHasLogined=====>> res =", res)
    if (res.code != 0) {
        console.log("auth.checkHasLogined=====>> token无效")
        return false
    }
    console.log("auth.checkHasLogined=====>> 登录验证成功")
    onLoginSuccess(res.data)
    return true
}

async function wxaCode() {
    return new Promise((resolve, reject) => {
        wx.login({
            success(res) {
                return resolve(res.code)
            },
            fail() {
                wx.showToast({
                    title: '获取code失败',
                    icon: 'none'
                })
                return resolve('获取code失败')
            }
        })
    })
}

async function getUserInfo() {
    return new Promise((resolve, reject) => {
        wx.getUserInfo({
            success: res => {
                return resolve(res)
            },
            fail: err => {
                console.error(err)
                return resolve()
            }
        })
    })
}

async function login(page, noReg) {
    console.log("auth.login=========>> page=", page)
    const _this = this
    wx.login({
        success: function(res) {
            console.log("auth.login=========>>",res)
            WXAPI.login_wx(res.code).then(function(res) {
                if (res.code == 10000) {
                    console.log("auth.login=========>>去注册 ")
                    if (noReg) {
                        return
                    }
                    // 去注册
                    _this.register_simple(page)
                    return
                }
                if (res.code != 0) {
                    // 登录错误
                    wx.showModal({
                        title: '无法登录',
                        content: res.msg,
                        showCancel: false
                    })
                    return
                }
                console.log("auth.login=========>>token登陆成功 ")
                onLoginSuccess(res.data)
                if (page) {
                    page.onShow()
                }
            })
        }
    })
}

async function register_simple(page) {
    let _this = this;
    wx.login({
        success: function(res) {
            let code = res.code; // 微信登录接口返回的 code 参数，下面注册接口需要用到
            let referrer = '' // 推荐人
            let referrer_storge = wx.getStorageSync('referrer');
            if (referrer_storge) {
                referrer = referrer_storge;
            }
            // 下面开始调用注册接口
            WXAPI.register_simple({
                code: code,
                referrer: referrer
            }).then(function(res) {
                console.log("register.register =======>>", res)
                _this.login(page, true);
            })
        }
    })
}

async function register(page) {
    let _this = this;
    wx.login({
        success: function(res) {
            let code = res.code; // 微信登录接口返回的 code 参数，下面注册接口需要用到
            wx.getUserInfo({
                success: function(res) {
                    console.log("register.getUserInfo =======>>", res)
                    let iv = res.iv;
                    let encryptedData = res.encryptedData;
                    let referrer = '' // 推荐人
                    let referrer_storge = wx.getStorageSync('referrer');
                    if (referrer_storge) {
                        referrer = referrer_storge;
                    }
                    // 下面开始调用注册接口
                    WXAPI.register_complex({
                        code: code,
                        encryptedData: encryptedData,
                        iv: iv,
                        referrer: referrer
                    }).then(function(res) {
                        console.log("register.register =======>>", res)
                        _this.login(page, true);
                    })
                }
            })
        }
    })
}

async function checkAndAuthorize(scope) {
    return new Promise((resolve, reject) => {
        wx.getSetting({
            success(res) {
                if (!res.authSetting[scope]) {
                    wx.authorize({
                        scope: scope,
                        success() {
                            resolve() // 无返回参数
                        },
                        fail(e) {
                            console.error(e)
                            // if (e.errMsg.indexof('auth deny') != -1) {
                            //   wx.showToast({
                            //     title: e.errMsg,
                            //     icon: 'none'
                            //   })
                            // }
                            wx.showModal({
                                title: '无权操作',
                                content: '需要获得您的授权',
                                showCancel: false,
                                confirmText: '立即授权',
                                confirmColor: '#e64340',
                                success(res) {
                                    wx.openSetting();
                                },
                                fail(e) {
                                    console.error(e)
                                    reject(e)
                                },
                            })
                        }
                    })
                } else {
                    resolve() // 无返回参数
                }
            },
            fail(e) {
                console.error(e)
                reject(e)
            }
        })
    })
}


module.exports = {
    checkHasLogined: checkHasLogined,
    wxaCode: wxaCode,
    getUserInfo: getUserInfo,
    login: login,
    register: register,
    register_simple: register_simple,
    loginOut: loginOut,
    isLogined: isLogined,
    checkAndAuthorize: checkAndAuthorize
}