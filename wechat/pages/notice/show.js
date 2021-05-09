const app = getApp();
const WXAPI = require('apifm-wxapi')
Page({
    data: {
        notice: {
            title: '',
            content: ''
        },
        noticeId: 0
    },
    onLoad: function(options) {
        this.setData({
            noticeId: options.id
        });
        var that = this;
        WXAPI.noticeDetail(options.id).then(function(res) {
            if (res.code == 0) {
                that.setData({
                    notice: res.data
                });
            }
        })
    },
    onShareAppMessage() {
        let notice = this.data.notice;
        let noticeId = this.data.noticeId;
        let uid = wx.getStorageSync('uid');
        let path = '/pages/start/loading?inviter_id=' + uid + '&route=/pages/notice/show%3fid%3d' + noticeId;
        return {
            title: notice.title,
            path: path
        };
    }
})

