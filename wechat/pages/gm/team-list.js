const WXAPI = require('apifm-wxapi')

var sliderWidth = 96;
Page({

    /**
     * 页面的初始数据
     */
    data: {
        teamList: [],
        typeTeams:{},

        tabs: ["代销列表", "申请审核", "条件不足"],
        activeIndex: 0,
        sliderOffset: 0,
        sliderLeft: 0,
    },

    onLoad: function(options) {
        const that = this;
        wx.getSystemInfo({
            success: function(res) {
                that.setData({
                    sliderLeft: (res.windowWidth / that.data.tabs.length - sliderWidth) / 2,
                    sliderOffset: res.windowWidth / that.data.tabs.length * that.data.activeIndex
                });
            }
        });
    },

    tabClick: function(e) {
        this.setData({
            sliderOffset: e.currentTarget.offsetLeft,
            activeIndex: e.currentTarget.id
        });
        this.fetchTabData(e.currentTarget.id)
    },

    fetchTabData(activeIndex) {
        let teamList = this.data.typeTeams[activeIndex];
        this.setData({
            teamList: teamList || [],
        });
    },

    /**
     * 生命周期函数--监听页面初次渲染完成
     */
    onReady: function() {

    },

    doTeamTap:function(e){
        let status = e.currentTarget.dataset.status
        let teamId = e.currentTarget.dataset.teamid
        const that = this;
        WXAPI.gmDoTeam(teamId, status).then(res => {
            if (res.code == 0) {
                that.onShow()
            }else{
                wx.showToast({
                    title: res.msg,
                    icon: 'none'
                })
            }
        });
    },

    /**
     * 生命周期函数--监听页面显示
     */
    onShow: function() {
        this.data.typeTeams = {}
        const that = this;
        WXAPI.gmTeamList().then(res => {
            if (res.code == 0) {
                let teamList = res.data.teamList || [];
                let typeTeams = that.data.typeTeams;
                for (var i = 0; i < teamList.length; i++) {
                    let team = teamList[i];
                    let idx = 0;
                    if (team.status == 1) {
                        team.statusStr = "申请审核";
                        idx = 1;
                    }else if (team.status == 2) {
                        team.statusStr = "条件不符";
                        idx = 2;
                    }else if (team.status == 3) {
                        team.statusStr = "代销";
                    }else if (team.status == 4) {
                        team.statusStr = "被取消";
                        idx = 2;
                    }
                    let typeTeam = typeTeams[idx];
                    if (!typeTeam) {
                        typeTeams[idx] = [];
                        typeTeam = typeTeams[idx];
                    }
                    typeTeam.push(team);
                }
                teamList = typeTeams[0];
                that.setData({
                    teamList: teamList || [],
                });
            }
        });
    },
    makePhoneCall: function(e) {
        let mobile = e.currentTarget.dataset.mobile
        wx.makePhoneCall({
            phoneNumber: mobile,
            success: function() {
                console.log('拨打成功')
            },
            fail: function() {
                console.log('拨打失败')
            }
        })
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

    }
})