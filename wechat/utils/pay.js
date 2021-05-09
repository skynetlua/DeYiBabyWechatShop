const WXAPI = require('apifm-wxapi')

/**
 * type: order 支付订单 recharge 充值 paybill 优惠买单
 * data: 扩展数据对象，用于保存参数
 */
function wxpay(orderId, money, cb) {
    console.log("wxpay ====>> orderId =", orderId, "money =", money)
    const postData = {
        orderId: orderId,
        money: money
    }
    wx.showLoading({title: '请求支付中', mask: true})
    WXAPI.wxpay(postData).then(function(res) {
        console.log("WXAPI.wxpay ====>> orderId =", orderId, "money =", money)
        wx.setStorageSync('successOrderId', orderId)
        wx.hideLoading()
        if (res.code == 0) {
            // 发起支付
            const resData = res.data
            wx.requestPayment({
                timeStamp: resData.timeStamp,
                signType: resData.signType,
                nonceStr: resData.nonceStr,
                package: resData.package,
                paySign: resData.paySign,
                fail: function(err) {
                    console.log("wxpay error", err)
                    wx.showToast({
                        title: '订单支付失败'
                    });
                    if (cb) {
                        cb()
                    }
                    // wx.redirectTo({
                    //     url: redirectUrl
                    // });
                },
                success: function() {
                    console.log("wx.requestPayment ====>> orderId =", orderId, "money =", money)
                    wx.showToast({
                        title: '支付成功'
                    });
                    if (cb) {
                        cb(1)
                    }
                    // wx.redirectTo({
                    //     url: redirectUrl
                    // });
                }
            })
        } else if (res.code == 1) {
            wx.showToast({
                title: '订单已支付成功'
            });
            if (cb) {
                cb(1)
            }
            // wx.redirectTo({
            //     url: redirectUrl
            // });
        } else {
            wx.showModal({
                title: '出错了',
                content: res.msg?res.msg:JSON.stringify(res),
                showCancel: false
            })
        }
    })
}

// function _requestPay(orderId, amount, cb){
//     if (amount <= 0) {
//         WXAPI.orderPay(orderId).then(function(res) {
//             if (cb) {cb()}
//         })
//     } else {
//         wxpay(orderId, amount, "/pages/order-list/index");
//     }
// }

function requestPay(orderId, amountReal, cb) {
    wx.showModal({
        title: '请确认订单',
        content: '订单金额: ' + (amountReal/100) + ' 元',
        confirmText: "确认",
        cancelText: "取消",
        success: function(res) {
            if (res.confirm) {
                // _requestPay(orderId, amountReal, cb)
                wxpay(orderId, amountReal, cb);
            } else {
                console.log('用户点击取消支付')
                if (cb) {
                    cb()
                }
            }
        }
    });

    // wxpay(orderId, amountReal, "/pages/order-list/index");
    // WXAPI.userAmount().then(function(res) {
    //     if (res.code == 0) {
    //         let data = res.data;
    //         let balance = data.balance?data.balance.toFixed(2):0;
    //         let amountLeft = amountReal - balance;
    //         let _msg = '订单金额: ' + amountReal + ' 元';
    //         if (balance > 0) {
    //             _msg += ',可用余额为 ' + balance + ' 元'
    //             if (amountLeft > 0) {
    //                 _msg += ',仍需微信支付 ' + amountLeft + ' 元'
    //             }
    //         }
    //         wx.showModal({
    //             title: '请确认支付',
    //             content: _msg,
    //             confirmText: "确认支付",
    //             cancelText: "取消支付",
    //             success: function(res) {
    //                 if (res.confirm) {
    //                     _requestPay(orderId, amountLeft, cb)
    //                 } else {
    //                     console.log('用户点击取消支付')
    //                 }
    //             }
    //         });
    //     } else {
    //         wx.showModal({
    //             title: '错误',
    //             content: '无法获取用户信息',
    //             showCancel: false
    //         })
    //     }
    // })
}

module.exports = {
    wxpay: wxpay,
    requestPay: requestPay
}

