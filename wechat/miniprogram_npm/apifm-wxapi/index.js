module.exports = (function() {
var __MODS__ = {};
var __DEFINE__ = function(modId, func, req) { var m = { exports: {}, _tempexports: {} }; __MODS__[modId] = { status: 0, func: func, req: req, m: m }; };
var __REQUIRE__ = function(modId, source) { if(!__MODS__[modId]) return require(source); if(!__MODS__[modId].status) { var m = __MODS__[modId].m; m._exports = m._tempexports; var desp = Object.getOwnPropertyDescriptor(m, "exports"); if (desp && desp.configurable) Object.defineProperty(m, "exports", { set: function (val) { if(typeof val === "object" && val !== m._exports) { m._exports.__proto__ = val.__proto__; Object.keys(val).forEach(function (k) { m._exports[k] = val[k]; }); } m._tempexports = val }, get: function () { return m._tempexports; } }); __MODS__[modId].status = 1; __MODS__[modId].func(__MODS__[modId].req, m, m.exports); } return __MODS__[modId].m.exports; };
var __REQUIRE_WILDCARD__ = function(obj) { if(obj && obj.__esModule) { return obj; } else { var newObj = {}; if(obj != null) { for(var k in obj) { if (Object.prototype.hasOwnProperty.call(obj, k)) newObj[k] = obj[k]; } } newObj.default = obj; return newObj; } };
var __REQUIRE_DEFAULT__ = function(obj) { return obj && obj.__esModule ? obj.default : obj; };
__DEFINE__(1587387818640, function(require, module, exports) {
/* eslint-disable */
// 小程序开发api接口工具包，https://github.com/gooking/wxapi
var API_BASE_URL = 'http://127.0.0.1:444'
var subDomain = 'api'
// function EventEmitter(){this._list = {}}
// EventEmitter.prototype = {
//     on(id, cb) {
//         var obj = this._list[id];if(!obj){this._list[id] = cb;return;}
//         if(typeof obj == 'function') {if(obj != cb) {this._list[id] = [obj, cb];}
//         }else{for (var i = 0; i < obj.length; i++){if(obj[i] == cb){return;}}obj.push(cb);}
//     },
//     emit(id) {
//         var obj = this._list[id];if(!obj) {console.log("emit not obj.id =", id);return;}
//         if (typeof obj == 'function') {obj.apply(null, arguments);}else{for(var i = 0; i < obj.length; i++) {obj[i].apply(null, arguments);}}
//     },
//     remove(id, cb) {
//         var obj = this._list[id];if (!obj) return;if(!cb) {delete this._list[id];return;}
//         if (typeof obj == 'function') {if (obj == cb) delete this._list[id];
//         }else{for (var i = 0; i < obj.length; i++) {if (obj[i] == cb) {obj.splice(i, 1);break;}}}
//     }
// }
// const emitter = new EventEmitter()

const request = (url, needSubDomain, method, data, errCb) => {
    const _url = API_BASE_URL + (needSubDomain ? '/' + subDomain : '') + url
    if (data) {
        data.token = wx.getStorageSync('token')
    }
    console.log("==>>:", _url, method, data)
    return new Promise((resolve, reject) => {
        wx.request({
            url: _url,
            method: method,
            data: data,
            header: {
                'Content-Type': 'application/x-www-form-urlencoded'
            },
            success(request) {
                console.log("<<==:", _url, request.data)
                resolve(request.data)
            },
            fail(error) {
                console.log("request fail:", _url, method, data)
                if (errCb) {
                    errCb(url)
                }
                reject(error)
            },
            complete(aaa) {
                // 加载完成
            }
        })
    })
}

/**
 * 小程序的promise没有finally方法，自己扩展下
 */
// Promise.prototype.finally = function(callback) {
//     var Promise = this.constructor;
//     return this.then(
//         function(value) {
//             Promise.resolve(callback()).then(
//                 function() {
//                     return value;
//                 }
//             );
//         },
//         function(reason) {
//             Promise.resolve(callback()).then(
//                 function() {
//                     throw reason;
//                 }
//             );
//         }
//     );
// }


module.exports = {
    init2: (a, b) => {
        API_BASE_URL = a;
        subDomain = b;
        wx.setStorageSync('remote_host', API_BASE_URL)
        wx.setStorageSync('remote_subdomain', subDomain)
    },
    init: (b) => {
        subDomain = b;
        wx.setStorageSync('remote_subdomain', subDomain)
    },
    request,
    // emitter,
    //page
    pageIndex: () => {
        return request('/page/index', true, 'get', {})
    },
    pageGoodsDetail: (goodsId) => {
        return request('/page/goods/detail', true, 'get', {goodsId})
    },

    fetchSubDomainByWxappAppid: (appid) => {
        return request('/subdomain/appid/wxapp', false, 'get', { appid })
    },
    vipLevel: () => {
        return request('/config/vipLevel', true, 'get')
    },

    //wechat
    login_wx: (code) => {
        return request('/wechat/login', true, 'post', {code,type: 2})
    },
    register_complex: (data) => {
        return request('/wechat/register/complex', true, 'post', data)
    },
    register_simple: (data) => {
        return request('/wechat/register/simple', true, 'post', data)
    },
    wxpay: (data) => {
        return request('/wechat/pay', true, 'post', data)
    },
    wxaQrcode: (data) => {
        return request('/wechat/qrcode', true, 'post', data)
    },
    bindMobileWxa: (encryptedData, iv) => {
        return request('/wechat/bindMobile', true, 'post', {encryptedData,iv})
    },

    //config
    queryConfigValue: (key) => {
        return request('/config/value', true, 'get', { key })
    },
    queryConfigBatch: (keys) => {
        return request('/config/values', true, 'get', { keys })
    },

    //user
    checkToken: (token) => {
        return request('/user/check/token', true, 'get', {token})
    },
    checkReferrer: (referrer) => {
        return request('/user/check/referrer', true, 'get', {referrer})
    },
    userDetail: () => {
        return request('/user/detail', true, 'get', {})
    },
    userAmount: () => {
        return request('/user/amount', true, 'get', {})
    },
    userWxinfo: (token) => {
        return request('/user/wxinfo', true, 'get', {token})
    },
    cashLogs: (data) => {
        return request('/user/cashLog', true, 'post', data)
    },
    payLogs: (data) => {
        return request('/user/payLog', true, 'post', data)
    },

    //address
    addAddress: (data) => {
        return request('/address/add', true, 'post', data)
    },
    updateAddress: (data) => {
        return request('/address/update', true, 'post', data)
    },
    deleteAddress: (id) => {
        return request('/address/delete', true, 'post', {id})
    },
    queryAddress: () => {
        return request('/address/list', true, 'get', {})
    },
    defaultAddress: () => {
        return request('/address/default', true, 'get', {})
    },
    addressDetail: (id) => {
        return request('/address/detail', true, 'get', {id})
    },

    //banner
    banners: (data) => {
        return request('/banner/list', true, 'get', data)
    },

    //goods
    goodsCategory: () => {
        return request('/goods/category/all', true, 'get')
    },
    goodsCategoryDetail: (id) => {
        return request('/goods/category/info', true, 'get', { id })
    },
    goods: (data) => {
        return request('/goods/list', true, 'post', data)
    },
    goodsDetail: (goodsId) => {
        return request('/goods/detail', true, 'get', {goodsId})
    },
    goodsSku: (id) => {
        return request('/goods/sku', true, 'get', {id})
    },
    goodsReputation: (data) => {
        return request('/goods/reputation', true, 'post', data)
    },
    goodsSubtypes: (categoryId) => {
        return request('/goods/category/subtypes', true, 'get', {categoryId})
    },
    goodsSubList: (categoryId, subType) => {
        return request('/goods/category/sublist', true, 'get', {categoryId, subType})
    },

    ///fav
    goodsFavList: (data) => {
        return request('/fav/list', true, 'post', data)
    },
    goodsFavPut: (goodsId) => {
        return request('/fav/add', true, 'post', {goodsId})
    },
    goodsFavCheck: (goodsId) => {
        return request('/fav/check', true, 'get', {goodsId})
    },
    goodsFavDelete: (goodsId) => {
        return request('/fav/delete', true, 'post', {goodsId})
    },

    //subshop
    fetchShops: (data) => {
        return request('/subshop/list', true, 'post', data)
    },
    fetchMyShops: (token) => {
        return request('/subshop/my', true, 'get', { token })
    },
    shopSubdetail: (id) => {
        return request('/subshop/detail', true, 'get', { id })
    },
    shopSubApply: (data) => {
        return request('/subshop/apply', true, 'post', data)
    },

    //discount
    coupons: (data) => {
        return request('/discount/coupon', true, 'get', data)
    },
    myCoupons: (data) => {
        return request('/discount/my', true, 'get', data)
    },
    
    couponDetail: (id) => {
        return request('/discount/detail', true, 'get', {id})
    },
    fetchCoupons: (data) => {
        return request('/discount/fetch', true, 'post', data)
    },
    sendCoupons: (data) => {
        return request('/discount/send', true, 'post', data)
    },
    exchangeCoupons: (number, pwd) => {
        return request('/discount/exchange', true, 'post', {number,pwd})
    },

    //live
    wxaMpLiveRooms: () => {
        return request('/live/rooms', true, 'get')
    },
    wxaMpLiveRoomHisVedios: (roomId) => {
        return request('/live/his', true, 'get', {roomId})
    },

    //cart
    cartInfo: () => {
        return request('/cart/info', true, 'get', {})
    },
    cartList: () => {
        return request('/cart/list', true, 'get', {})
    },
    cartAdd: (goodsId, buyNumber, skuId) => {
        wx.setStorageSync('tabBarBadge_time', 0)
        return request('/cart/add', true, 'post', {goodsId, buyNumber, skuId})
    },
    cartModifyNumber: (id, number) => {
        return request('/cart/modifyNumber', true, 'post', {id,number})
    },
    cartRemove: (id) => {
        wx.setStorageSync('tabBarBadge_time', 0)
        return request('/cart/remove', true, 'post', {id})
    },
    cartEmpty: () => {
        wx.setStorageSync('tabBarBadge_time', 0)
        return request('/cart/empty', true, 'post', {})
    },
    // cartQuick: (goodsId, buyNumber, skuId) => {
    //     return request('/cart/quick', true, 'get', {goodsId, buyNumber, skuId})
    // },

    //order
    orderPrepare: (data) => {
        return request('/order/prepare', true, 'post', data)
    },
    orderCreate: (data) => {
        return request('/order/create', true, 'post', data)
    },
    orderList: () => {
        return request('/order/list', true, 'post', {})
    },
    orderDetail: (orderId, hxNumber = '') => {
        return request('/order/detail', true, 'get', {orderId, hxNumber})
    },
    orderDelivery: (orderId) => {
        return request('/order/delivery', true, 'post', {orderId})
    },
    orderReputation: (data) => {
        return request('/order/reputation', true, 'post', data)
    },
    orderClose: (orderId) => {
        return request('/order/close', true, 'post', {orderId})
    },
    orderDelete: (orderId) => {
        return request('/order/delete', true, 'post', {orderId})
    },
    orderPay: (orderId) => {
        return request('/order/pay', true, 'post', {orderId})
    },
    orderHX: (hxNumber) => {
        return request('/order/hx', true, 'post', {hxNumber})
    },
    orderStatistics: () => {
        return request('/order/statistics', true, 'get', {})
    },
    orderRefunds: (orderId) => {
        return request('/order/refund', true, 'get', {orderId})
    },
    refundApply: (data) => {
        return request('/order/refundApply/apply', true, 'post', data)
    },
    refundApplyDetail: (orderId) => {
        return request('/order/refundApply/info', true, 'get', {orderId})
    },
    refundApplyCancel: (orderId) => {
        return request('/order/refundApply/cancel', true, 'post', {orderId})
    },
 
    //distribute
    fxInfo: () => {
        return request('/distribute/info', true, 'get', {})
    },
    fxApply: (name, mobile) => {
        return request('/distribute/apply', true, 'post', {name, mobile})
    },
    fxApplyProgress: () => {
        return request('/distribute/apply/progress', true, 'get', {})
    },
    fxMembers: (data) => {
        return request('/distribute/members', true, 'post', data)
    },
    fxCommisionLog: (data) => {
        return request('/distribute/log', true, 'post', data)
    },

    //team
    myTeamInfo: () => {
        return request('/team/myteam', true, 'get', {})
    },
    teamCreate: (name, mobile, slogan) => {
        return request('/team/create', true, 'post', {name, mobile, slogan})
    },
    teamJoin: (name, mobile) => {
        return request('/team/join', true, 'post', {name, mobile})
    },
    teamKick: (memberId) => {
        return request('/team/kick', true, 'post', {memberId})
    },
    teamAppoint: (memberId, post) => {
        return request('/team/appoint', true, 'post', {memberId, post})
    },
    teamMembers: () => {
        return request('/team/members', true, 'post', {})
    },

    //gm
    gmOrderList: (status) => {
        return request("/gm/order/list", true, "get", {status})
    },
    gmOrderDo: (orderId, status, playerId) => {
        return request("/gm/order/do", true, "get", {orderId, status, playerId})
    },
    gmOrderDetail: (orderId, playerId) => {
        return request("/gm/order/detail", true, "get", {orderId, playerId})
    },
    gmOrderCoupon: (orderId, playerId, amount, couponId) => {
        return request("/gm/order/coupon", true, "post", {orderId, playerId, amount, couponId})
    },
    gmRefundConfirm: (orderId, playerId) => {
        return request("/gm/refund/confirm", true, "get", {orderId, playerId})
    },
    gmRefundCancel: (orderId, playerId) => {
        return request("/gm/refund/cancel", true, "get", {orderId, playerId})
    },
    gmTeamList: () => {
        return request("/gm/team/list", true, "get", {})
    },
    gmDoTeam: (teamId, status) => {
        return request("/gm/do/team", true, "get", {teamId, status})
    },
    gmGoodsInfo: (goodsId) => {
        return request("/gm/goods/info", true, "get", {goodsId})
    },
    gmGoodsGoodsData: (barCode) => {
        return request("/gm/goods/goodsdata", true, "get", {barCode})
    },
    gmGoodsGoodsDatas: () => {
        return request("/gm/goods/goodsdatas", true, "get", {})
    },
    gmGoodsUpdate: (goodsInfo) => {
        return request("/gm/goods/update", true, "post", goodsInfo)
    },
    gmGoodsUpdateInfo: (goodsInfo) => {
        return request("/gm/goods/update/info", true, "post", goodsInfo)
    },
    gmGoodsRemove: (goodsId) => {
        return request("/gm/goods/remove", true, "get", {goodsId})
    },
    gmGoodsList: (data) => {
        data = data || {};
        return request("/gm/goods/list", true, "get", data)
    },
    gmBarCodeGoodsId: (barCode) => {
        return request("/gm/goods/barcode", true, "get", {barCode})
    },
    gmGoodsCategory: (data) => {
        data = data || {};
        return request("/gm/goods/category", true, "get", data)
    },
    gmUploadGoods: (filePath, categoryId, goodsId, part, idx) => {
        const uploadUrl = API_BASE_URL + '/' + subDomain + '/gm/upload/goods'
        const token = wx.getStorageSync('token')
        return new Promise((resolve, reject) => {
            wx.uploadFile({
                url: uploadUrl,
                filePath: filePath,
                name: 'upfile',
                formData: {
                    token:token,
                    categoryId:categoryId,
                    goodsId:goodsId,
                    part:part,
                    idx:idx
                },
                success(res) {
                    resolve(JSON.parse(res.data));
                },
                fail(error) {
                    reject(error);
                },
                complete(aaa) {
                    // 加载完成
                }
            });
        });
    },
    gmCategoryList: () => {
        return request("/gm/category/list", true, "get", {})
    },
    gmCategoryUpdate: (data) => {
        return request("/gm/category/update", true, "post", data)
    },
    gmCategoryRemove: (categoryId) => {
        return request("/gm/category/remove", true, "get", {categoryId})
    },
    gmUploadCategory: (filePath, categoryId) => {
        const uploadUrl = API_BASE_URL + '/' + subDomain + '/gm/upload/category'
        const token = wx.getStorageSync('token')
        return new Promise((resolve, reject) => {
            wx.uploadFile({
                url: uploadUrl,
                filePath: filePath,
                name: 'upfile',
                formData: {
                    token:token,
                    categoryId:categoryId,
                },
                success(res) {
                    resolve(JSON.parse(res.data));
                },
                fail(error) {
                    reject(error);
                },
                complete(aaa) {
                    // 加载完成
                }
            });
        });
    },

    gmGoodsLoadPicture: () => {
        return request("/gm/goods/load/picture", true, "get", {})
    },

    gmUploadExcel: (filePath, fileName, opt) => {
        const uploadUrl = API_BASE_URL + '/' + subDomain + '/gm/upload/excel'
        const token = wx.getStorageSync('token')
        return new Promise((resolve, reject) => {
            wx.uploadFile({
                url: uploadUrl,
                filePath: filePath,
                name: 'upfile',
                formData: {
                    token:token,
                    fileName:fileName,
                    opt:opt,
                },
                success(res) {
                    resolve(JSON.parse(res.data));
                },
                fail(error) {
                    reject(error);
                },
                complete(aaa) {
                    // 加载完成
                }
            });
        });
    },


    //deposit
    // depositList: (data) => {
    //     return request('/deposit/list', true, 'post', data)
    // },
    //  payDeposit: (data) => {
    //     return request('/deposit/pay', true, 'post', data)
    // },
    // depositInfo: (id) => {
    //     return request('/deposit/info', true, 'get', {id})
    // },
    // depositBackApply: (id) => {
    //     return request('/deposit/back/apply', true, 'post', {id})
    // },

    //withdraw
    withDrawApply: (money) => {
        return request('/withdraw/apply', true, 'post', {money})
    },
    withDrawDetail: (id) => {
        return request('/withdraw/detail', true, 'get', {id})
    },
    withDrawLogs: (data) => {
        return request('/withdraw/list', true, 'post', data)
    },

    //region
    province: () => {
        return request('/region/province', false, 'get', {})
    },
    nextRegion: (pid) => {
        return request('/region/child', false, 'get', {pid})
    },

    uploadFile: (token, sendFilePath, expireHours = '') => {
        const uploadUrl = API_BASE_URL + '/' + subDomain + '/upload/file'
        return new Promise((resolve, reject) => {
            wx.uploadFile({
                url: uploadUrl,
                filePath: sendFilePath,
                name: 'upfile',
                formData: {
                    token,
                    expireHours
                },
                success(res) {
                    resolve(JSON.parse(res.data))
                },
                fail(error) {
                    reject(error)
                },
                complete(aaa) {
                    // 加载完成
                }
            })
        })
    },

    noticeList: (data) => {
        return request('/notice/list', true, 'post', data)
    },
    noticeLastOne: (type) => {
        return request('/notice/lastone', true, 'get', {type})
    },
    noticeDetail: (id) => {
        return request('/notice/detail', true, 'get', {id})
    },


    //
    // queryMobileLocation: (mobile = '') => {
    //     return request('/common/mobile-segment/location', false, 'get', { mobile })
    // },
    // nextMobileSegment: (data) => {
    //     return request('/common/mobile-segment/next', false, 'post', data)
    // },

    // scoreRules: (data) => {
    //     return request('/score/send/rule', true, 'post', data)
    // },
    // scoreSignRules: () => {
    //     return request('/score/sign/rules', true, 'get', {})
    // },
    // scoreSign: (token) => {
    //     return request('/score/sign', true, 'post', {
    //         token
    //     })
    // },
    // scoreSignLogs: (data) => {
    //     return request('/score/sign/logs', true, 'post', data)
    // },
    // scoreTodaySignedInfo: (token) => {
    //     return request('/score/today-signed', true, 'get', {
    //         token
    //     })
    // },
    // scoreExchange: (token, number) => {
    //     return request('/score/exchange', true, 'post', {
    //         number,
    //         token
    //     })
    // },
    // scoreExchangeCash: (token, deductionScore) => {
    //     return request('/score/exchange/cash', true, 'post', {
    //         deductionScore,
    //         token
    //     })
    // },
    // scoreLogs: (data) => {
    //     return request('/score/logs', true, 'post', data)
    // },
    // shareGroupGetScore: (code, referrer, encryptedData, iv) => {
    //     return request('/score/share/wxa/group', true, 'post', {
    //         code,
    //         referrer,
    //         encryptedData,
    //         iv
    //     })
    // },
    // kanjiaSet: (goodsId) => {
    //     return request('/shop/goods/kanjia/set/v2', true, 'get', { goodsId })
    // },
    // kanjiaJoin: (token, kjid) => {
    //     return request('/shop/goods/kanjia/join', true, 'post', {
    //         kjid,
    //         token
    //     })
    // },
    // kanjiaDetail: (kjid, joiner) => {
    //     return request('/shop/goods/kanjia/info', true, 'get', {
    //         kjid,
    //         joiner
    //     })
    // },
    // kanjiaHelp: (token, kjid, joiner, remark = '') => {
    //     return request('/shop/goods/kanjia/help', true, 'post', {
    //         kjid,
    //         joinerUser: joiner,
    //         token,
    //         remark
    //     })
    // },
    // kanjiaClear: (token, kjid) => {
    //     return request('/shop/goods/kanjia/clear', true, 'post', {
    //         kjid,
    //         token
    //     })
    // },
    // kanjiaMyJoinInfo: (token, kjid) => {
    //     return request('/shop/goods/kanjia/my', true, 'get', {
    //         kjid,
    //         token
    //     })
    // },
    // kanjiaHelpDetail: (token, kjid, joiner) => {
    //     return request('/shop/goods/kanjia/myHelp', true, 'get', {
    //         kjid,
    //         joinerUser: joiner,
    //         token
    //     })
    // },

    // addTempleMsgFormid: (token, type, formId) => {
    //     return request('/template-msg/wxa/formId', true, 'post', {
    //         token,
    //         type,
    //         formId
    //     })
    // },
    // sendTempleMsg: (data) => {
    //     return request('/template-msg/put', true, 'post', data)
    // },
    
    // ttpay: (data) => {
    //     return request('/pay/tt/microapp', true, 'post', data)
    // },
    // payQuery: (token, outTradeId) => {
    //     return request('/pay/query', true, 'get', { token, outTradeId })
    // },
    // wxpaySaobei: (data) => {
    //     return request('/pay/lcsw/wxapp', true, 'post', data)
    // },
    // wxpayWepayez: (data) => {
    //     return request('/pay/wepayez/wxapp', true, 'post', data)
    // },
    // alipay: (data) => {
    //     return request('/pay/alipay/semiAutomatic/payurl', true, 'post', data)
    // },
    
    // loginWxaMobile: (code, encryptedData, iv) => {
    //     return request('/user/wxapp/login/mobile', true, 'post', {
    //         code,
    //         encryptedData,
    //         iv
    //     })
    // },
    // login_username: (data) => {
    //     return request('/user/username/login', true, 'post', data)
    // },
    // bindUsername: (token, username, pwd = '') => {
    //     return request('/user/username/bindUsername', true, 'post', {
    //         token,
    //         username,
    //         pwd
    //     })
    // },
    // login_mobile: (mobile, pwd, deviceId = '', deviceName = '') => {
    //     return request('/user/m/login', true, 'post', {
    //         mobile,
    //         pwd,
    //         deviceId,
    //         deviceName
    //     })
    // },
    // resetPwdUseMobileCode: (mobile, pwd, code) => {
    //     return request('/user/m/reset-pwd', true, 'post', {
    //         mobile,
    //         pwd,
    //         code
    //     })
    // },
    // resetPwdUseEmailCode: (email, pwd, code) => {
    //     return request('/user/email/reset-pwd', true, 'post', {
    //         email,
    //         pwd,
    //         code
    //     })
    // },
    
    // register_username: (data) => {
    //     return request('/user/username/register', true, 'post', data)
    // },
    // register_mobile: (data) => {
    //     return request('/user/m/register', true, 'post', data)
    // },
    
    
    // goodsLimitations: (goodsId, priceId = '') => {
    //     return request('/shop/goods/limitation', true, 'get', {
    //         goodsId,
    //         priceId
    //     })
    // },
    // goodsPrice: (goodsId, propertyChildIds) => {
    //     return request('/shop/goods/price', true, 'post', {
    //         goodsId,
    //         propertyChildIds
    //     })
    // },
    // goodsPriceDaily: (goodsId, priceId = '') => {
    //     return request('/shop/goods/price/day', true, 'get', {
    //         goodsId,
    //         priceId
    //     })
    // },
    // goodsPriceFreight: (data) => {
    //     return request('/shop/goods/price/freight', true, 'get', data)
    // },
    // goodsRebate: (token, goodsId) => {
    //     return request('/shop/goods/rebate/v2', true, 'get', {
    //         token,
    //         goodsId
    //     })
    // },
    
   
    
    
    
    // pingtuanSet: (goodsId) => {
    //     return request('/shop/goods/pingtuan/set', true, 'get', {
    //         goodsId
    //     })
    // },
    // pingtuanSets: (goodsIdArray) => {
    //     return request('/shop/goods/pingtuan/sets', true, 'get', {
    //         goodsId: goodsIdArray.join()
    //     })
    // },
    // pingtuanOpen: (token, goodsId) => {
    //     return request('/shop/goods/pingtuan/open', true, 'post', {
    //         goodsId,
    //         token
    //     })
    // },
    // pingtuanList: (data) => {
    //     return request('/shop/goods/pingtuan/list/v2', true, 'post', data)
    // },
    // pingtuanJoinUsers: (tuanId) => {
    //     return request('/shop/goods/pingtuan/joiner', true, 'get', { tuanId })
    // },
    // pingtuanMyJoined: (data) => {
    //     return request('/shop/goods/pingtuan/my-join-list', true, 'post', data)
    // },
    // friendlyPartnerList: (type = '') => {
    //     return request('/friendly-partner/list', true, 'post', {
    //         type
    //     })
    // },
    // friendList: (data) => {
    //     return request('/user/friend/list', true, 'post', data)
    // },
    // addFriend: (token, uid) => {
    //     return request('/user/friend/add', true, 'post', { token, uid })
    // },
    // friendUserDetail: (token, uid) => {
    //     return request('/user/friend/detail', true, 'get', { token, uid })
    // },
    // videoDetail: (videoId) => {
    //     return request('/media/video/detail', true, 'get', {
    //         videoId
    //     })
    // },
    
    // bindMobileSms: (token, mobile, code, pwd = '') => {
    //     return request('/user/m/bind-mobile', true, 'post', {
    //         token,
    //         mobile,
    //         code,
    //         pwd
    //     })
    // },
    
    
    // rechargeSendRules: () => {
    //     return request('/user/recharge/send/rule', true, 'get')
    // },
    // payBillDiscounts: () => {
    //     return request('/payBill/discount', true, 'get')
    // },
    // payBill: (token, money) => {
    //     return request('/payBill/pay', true, 'post', { token, money })
    // },

    
    // uploadFileFromUrl: (remoteFileUrl = '', ext = '') => {
    //     return request('/dfs/upload/url', true, 'post', { remoteFileUrl, ext })
    // },
    // uploadFileList: (path = '') => {
    //     return request('/dfs/upload/list', true, 'post', { path })
    // },
    
    // cmsCategories: () => {
    //     return request('/cms/category/list', true, 'get', {})
    // },
    // cmsCategoryDetail: (id) => {
    //     return request('/cms/category/info', true, 'get', { id })
    // },
    // cmsArticles: (data) => {
    //     return request('/cms/news/list', true, 'post', data)
    // },
    // cmsArticleUsefulLogs: (data) => {
    //     return request('/cms/news/useful/logs', true, 'post', data)
    // },
    // cmsArticleDetail: (id) => {
    //     return request('/cms/news/detail', true, 'get', { id })
    // },
    // cmsArticlePreNext: (id) => {
    //     return request('/cms/news/preNext', true, 'get', { id })
    // },
    // cmsArticleCreate: (data) => {
    //     return request('/cms/news/put', true, 'post', data)
    // },
    // cmsArticleDelete: (token, id) => {
    //     return request('/cms/news/del', true, 'post', { token, id })
    // },
    // cmsArticleUseless: (data) => {
    //     return request('/cms/news/useful', true, 'post', data)
    // },
    // cmsPage: (key) => {
    //     return request('/cms/page/info/v2', true, 'get', { key })
    // },
    // cmsTags: () => {
    //     return request('/cms/tags/list', true, 'get', {})
    // },
    // invoiceList: (data) => {
    //     return request('/invoice/list', true, 'post', data)
    // },
    // invoiceApply: (data) => {
    //     return request('/invoice/apply', true, 'post', data)
    // },
    // invoiceDetail: (token, id) => {
    //     return request('/invoice/info', true, 'get', { token, id })
    // },
    

    // addComment: (data) => {
    //     return request('/comment/add', true, 'post', data)
    // },
    // commentList: (data) => {
    //     return request('/comment/list', true, 'post', data)
    // },
    // modifyUserInfo: (data) => {
    //     return request('/user/modify', true, 'post', data)
    // },
    // uniqueId: (type = '') => {
    //     return request('/uniqueId/get', true, 'get', { type })
    // },
    // queryBarcode: (barcode = '') => {
    //     return request('/barcode/info', true, 'get', { barcode })
    // },
    // luckyInfo: (id) => {
    //     return request('/luckyInfo/info/v2', true, 'get', { id })
    // },
    // luckyInfoJoin: (id, token) => {
    //     return request('/luckyInfo/join', true, 'post', { id, token })
    // },
    // luckyInfoJoinMy: (id, token) => {
    //     return request('/luckyInfo/join/my', true, 'get', { id, token })
    // },
    // luckyInfoJoinLogs: (data) => {
    //     return request('/luckyInfo/join/logs', true, 'post', data)
    // },
    // jsonList: (data) => {
    //     return request('/json/list', true, 'post', data)
    // },
    // jsonSet: (data) => {
    //     return request('/json/set', true, 'post', data)
    // },
    // jsonDelete: (token = '', id) => {
    //     return request('/json/delete', true, 'post', { token, id })
    // },
    // graphValidateCodeUrl: (key = Math.random()) => {
    //     const _url = API_BASE_URL + '/' + subDomain + '/verification/pic/get?key=' + key
    //     return _url
    // },
    // graphValidateCodeCheck: (key = Math.random(), code) => {
    //     return request('/verification/pic/check', true, 'post', { key, code })
    // },
    // shortUrl: (url = '') => {
    //     return request('/common/short-url/shorten', false, 'post', { url })
    // },
    // smsValidateCode: (mobile, key = '', picCode = '') => {
    //     return request('/verification/sms/get', true, 'get', { mobile, key, picCode })
    // },
    // smsValidateCodeCheck: (mobile, code) => {
    //     return request('/verification/sms/check', true, 'post', { mobile, code })
    // },
    // mailValidateCode: (mail) => {
    //     return request('/verification/mail/get', true, 'get', { mail })
    // },
    // mailValidateCodeCheck: (mail, code) => {
    //     return request('/verification/mail/check', true, 'post', { mail, code })
    // },
    // mapDistance: (lat1, lng1, lat2, lng2) => {
    //     return request('/common/map/distance', false, 'get', { lat1, lng1, lat2, lng2 })
    // },
    // mapQQAddress: (location = '', coord_type = '5') => {
    //     return request('/common/map/qq/address', false, 'get', { location, coord_type })
    // },
    // mapQQSearch: (data) => {
    //     return request('/common/map/qq/search', false, 'post', data)
    // },
    // virtualTraderList: (data) => {
    //     return request('/virtualTrader/list', true, 'post', data)
    // },
    // virtualTraderDetail: (token, id) => {
    //     return request('/virtualTrader/info', true, 'get', { token, id })
    // },
    // virtualTraderBuy: (token, id) => {
    //     return request('/virtualTrader/buy', true, 'post', { token, id })
    // },
    // virtualTraderMyBuyLogs: (data) => {
    //     return request('/virtualTrader/buy/logs', true, 'post', data)
    // },
    // queuingTypes: (status = '') => {
    //     return request('/queuing/types', true, 'get', { status })
    // },
    // queuingGet: (token, typeId, mobile = '') => {
    //     return request('/queuing/get', true, 'post', { token, typeId, mobile })
    // },
    // queuingMy: (token, typeId = '', status = '') => {
    //     return request('/queuing/my', true, 'get', { token, typeId, status })
    // },
    // idcardCheck: (token, name, idCardNo) => {
    //     return request('/user/idcard', true, 'post', { token, name, idCardNo })
    // },
    // loginout: (token) => {
    //     return request('/user/loginout', true, 'get', { token })
    // },
    // userLevelList: (data) => {
    //     return request('/user/level/list', true, 'post', data)
    // },
    // userLevelDetail: (levelId) => {
    //     return request('/user/level/info', true, 'get', { id: levelId })
    // },
    // userLevelPrices: (levelId) => {
    //     return request('/user/level/prices', true, 'get', { levelId })
    // },
    // userLevelBuy: (token, priceId, isAutoRenew = false, remark = '') => {
    //     return request('/user/level/buy', true, 'post', {
    //         token,
    //         userLevelPriceId: priceId,
    //         isAutoRenew,
    //         remark
    //     })
    // },
    // userLevelBuyLogs: (data) => {
    //     return request('/user/level/buyLogs', true, 'post', data)
    // },
    // messageList: (data) => {
    //     return request('/user/message/list', true, 'post', data)
    // },
    // messageRead: (token, id) => {
    //     return request('/user/message/read', true, 'post', { token, id })
    // },
    // messageDelete: (token, id) => {
    //     return request('/user/message/del', true, 'post', { token, id })
    // },
    // bindOpenid: (token, code) => {
    //     return request('/user/wxapp/bindOpenid', true, 'post', {
    //         token,
    //         code,
    //         type: 2
    //     })
    // },
    // encryptedData: (code, encryptedData, iv) => {
    //     return request('/user/wxapp/decode/encryptedData', true, 'post', {
    //         code,
    //         encryptedData,
    //         iv
    //     })
    // },
    // scoreDeductionRules: (type = '') => {
    //     return request('/score/deduction/rules', true, 'get', { type })
    // },
    // voteItems: (data) => {
    //     return request('/vote/items', true, 'post', data)
    // },
    // voteItemDetail: (id) => {
    //     return request('/vote/info', true, 'get', { id })
    // },
    // vote: (token, voteId, items, remark) => {
    //     return request('/vote/vote', true, 'post', {
    //         token,
    //         voteId,
    //         items: items.join(),
    //         remark
    //     })
    // },
    // myVote: (token, voteId) => {
    //     return request('/vote/vote/info', true, 'get', {
    //         token,
    //         voteId,
    //     })
    // },
    // voteLogs: (data) => {
    //     return request('/vote/vote/list', true, 'post', data)
    // },
    // yuyueItems: (data) => {
    //     return request('/yuyue/items', true, 'post', data)
    // },
    // yuyueItemDetail: (id) => {
    //     return request('/yuyue/info', true, 'get', { id })
    // },
    // yuyueJoin: (data) => {
    //     return request('/yuyue/join', true, 'post', data)
    // },
    // yuyueJoinPay: (token, joinId) => {
    //     return request('/yuyue/pay', true, 'post', {
    //         token,
    //         joinId
    //     })
    // },
    // yuyueJoinUpdate: (token, joinId, extJsonStr) => {
    //     return request('/yuyue/join/update', true, 'post', {
    //         token,
    //         joinId,
    //         extJsonStr
    //     })
    // },
    // yuyueMyJoinInfo: (token, joinId) => {
    //     return request('/yuyue/join/info', true, 'post', {
    //         token,
    //         joinId
    //     })
    // },
    // yuyueMyJoinLogs: (data) => {
    //     return request('/yuyue/join/list', true, 'post', data)
    // },
    // yuyueTeams: (data) => {
    //     return request('/yuyue/info/teams', true, 'post', data)
    // },
    // yuyueTeamDetail: (teamId) => {
    //     return request('/yuyue/info/team', true, 'get', { teamId })
    // },
    // yuyueTeamMembers: (data) => {
    //     return request('/yuyue/info/team/members', true, 'post', data)
    // },
    // yuyueTeamDeleteMember: (token, joinId) => {
    //     return request('/yuyue/info/team/members/del', true, 'post', data)
    // },
    // register_email: (data) => {
    //     return request('/user/email/register', true, 'post', data)
    // },
    // login_email: (data) => {
    //     return request('/user/email/login', true, 'post', data)
    // },
    // bindEmail: (token, email, code, pwd = '') => {
    //     return request('/user/email/bindUsername', true, 'post', {
    //         token,
    //         email,
    //         code,
    //         pwd
    //     })
    // },
    // siteStatistics: () => {
    //     return request('/site/statistics', true, 'get')
    // },

    // cmsArticleFavPut: (token, newsId) => {
    //     return request('/cms/news/fav/add', true, 'post', { token, newsId })
    // },
    // cmsArticleFavCheck: (token, newsId) => {
    //     return request('/cms/news/fav/check', true, 'get', { token, newsId })
    // },
    // cmsArticleFavList: (data) => {
    //     return request('/cms/news/fav/list', true, 'post', data)
    // },
    // cmsArticleFavDeleteById: (token, id) => {
    //     return request('/cms/news/fav/delete', true, 'post', { token, id })
    // },
    // cmsArticleFavDeleteByNewsId: (token, newsId) => {
    //     return request('/cms/news/fav/delete', true, 'post', { token, newsId })
    // },
    
    // growthLogs: (data) => {
    //     return request('/growth/logs', true, 'post', data)
    // },
    // exchangeScoreToGrowth: (token, deductionScore) => {
    //     return request('/growth/exchange', true, 'post', {
    //         token,
    //         deductionScore
    //     })
    // },
    
}
}, function(modId) {var map = {}; return __REQUIRE__(map[modId], modId); })
return __REQUIRE__(1587387818640);
})()
//# sourceMappingURL=index.js.map