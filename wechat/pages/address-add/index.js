const WXAPI = require('apifm-wxapi')
const AUTH = require('../../utils/auth')
Page({
    data: {
        provinces: undefined, // 省份数据数组
        pIndex: 0, //选择的省下标
        cities: undefined, // 城市数据数组
        cIndex: 0, //选择的市下标
        areas: undefined, // 区县数数组
        aIndex: 0, //选择的区下标
    },
    async provinces(provinceId, cityId, areaId) {
        const res = await WXAPI.province();
        if (res.code == 0) {
            for (var i = 0; i < res.data.length; i++) {
                let item = res.data[i]
                if (item.name == '广东省'){
                    res.data.splice(i, 1);
                    res.data.unshift(item);
                    break;
                }
            }
            const provinces = [{
                id: 0,
                name: '请选择'
            }].concat(res.data);
            let pIndex = 0;
            if (provinceId) {
                pIndex = provinces.findIndex(ele => {
                    return ele.id == provinceId;
                });
            }
            this.setData({
                pIndex,
                provinces: provinces
            });
            if (provinceId) {
                const e = { detail: { value: pIndex } };
                this.provinceChange(e, cityId, areaId);
            }
        }
    },
    async provinceChange(e, cityId, areaId) {
        const index = e.detail.value;
        this.setData({
            pIndex: index
        });
        const pid = this.data.provinces[index].id;
        if (pid == 0) {
            this.setData({
                cities: null,
                cIndex: 0,
                areas: null,
                aIndex: 0
            });
            return;
        }
        const res = await WXAPI.nextRegion(pid);
        if (res.code == 0) {
            for (var i = 0; i < res.data.length; i++) {
                let item = res.data[i]
                if (item.name == '清远市') {
                    res.data.splice(i, 1);
                    res.data.unshift(item);
                    break;
                }
            }
            const cities = [{
                id: 0,
                name: '请选择'
            }].concat(res.data);
            let cIndex = 0;
            if (cityId) {
                cIndex = cities.findIndex(ele => {
                    return ele.id == cityId;
                });
            }
            this.setData({
                cIndex,
                cities: cities
            });
            if (cityId) {
                const e = { detail: { value: cIndex } };
                this.cityChange(e, areaId);
            }
        }
    },
    async cityChange(e, areaId) {
        const index = e.detail.value;
        this.setData({
            cIndex: index
        });
        const pid = this.data.cities[index].id;
        if (pid == 0) {
            this.setData({
                areas: null,
                aIndex: 0
            });
            return;
        }
        const res = await WXAPI.nextRegion(pid);
        if (res.code == 0) {
            for (var i = 0; i < res.data.length; i++) {
                let item = res.data[i]
                if (item.name == '清城区') {
                    res.data.splice(i, 1);
                    res.data.unshift(item);
                    break;
                }
            }
            const areas = [{
                id: 0,
                name: '请选择'
            }].concat(res.data);
            let aIndex = 0;
            if (areaId) {
                aIndex = areas.findIndex(ele => {
                    return ele.id == areaId;
                });
            }
            this.setData({
                aIndex,
                areas: areas
            });
            if (areaId) {
                const e = { detail: { value: aIndex } };
                this.areaChange(e);
            }
        }
    },
    async areaChange(e) {
        const index = e.detail.value
        this.setData({
            aIndex: index
        })
    },
    async bindSave(e) {
        var data = this.data;
        if (data.pIndex == 0) {
            wx.showToast({
                title: '请选择省份',
                icon: 'none'
            })
            return
        }
        if (data.cIndex == 0) {
            wx.showToast({
                title: '请选择城市',
                icon: 'none'
            })
            return
        }
        const linkMan = e.detail.value.linkMan;
        const address = e.detail.value.address;
        const mobile = e.detail.value.mobile;
        const code = '322000';
        if (linkMan == "") {
            wx.showToast({
                title: '请填写联系人姓名',
                icon: 'none'
            })
            return;
        }
        if (mobile == "") {
            wx.showToast({
                title: '请填写手机号码',
                icon: 'none'
            })
            return;
        }
        if (address == "") {
            wx.showToast({
                title: '请填写详细地址',
                icon: 'none'
            })
            return;
        }
        const postData = {
            linkMan: linkMan,
            address: address,
            mobile: mobile,
            code: code,
        }
        if (data.pIndex > 0) {
            postData.provinceId = data.provinces[data.pIndex].id
        }
        if (data.cIndex > 0) {
            postData.cityId = data.cities[data.cIndex].id
        }
        if (data.aIndex > 0) {
            postData.areaId = data.areas[data.aIndex].id
        }
        let res
        if (data.id) {
            postData.id = data.id
            res = await WXAPI.updateAddress(postData)
        } else {
            res = await WXAPI.addAddress(postData)
        }
        if (res.code != 0) {
            wx.hideLoading();
            wx.showToast({
                title: res.msg,
                icon: 'none'
            });
            return;
        } else {
            wx.navigateBack();
        }
    },
    async onLoad(e) {
        if (e.id) { // 修改初始化数据库数据
            const res = await WXAPI.addressDetail(e.id);
            if (res.code == 0) {
                var data = res.data;
                this.setData({
                    id: e.id,
                    addressData: data.info
                });
                this.provinces(data.info.provinceId, data.info.cityId, data.info.areaId);
            } else {
                wx.showModal({
                    title: '错误',
                    content: '无法获取快递地址数据',
                    showCancel: false
                });
            }
        } else {
            this.provinces();
        }
    },
    deleteAddress: function(e) {
        const id = e.currentTarget.dataset.id;
        wx.showModal({
            title: '提示',
            content: '确定要删除该收货地址吗？',
            success: function(res) {
                if (res.confirm) {
                    WXAPI.deleteAddress(id).then(function(res) {
                        wx.navigateBack({})
                    })
                } else {
                    console.log('用户点击取消')
                }
            }
        })
    },
    async readFromWx() {
        let that = this;
        wx.chooseAddress({
            success: function(res) {
                console.log(res)
                const provinceName = res.provinceName;
                const cityName = res.cityName;
                const areaName = res.countyName;
                // 读取省
                const pIndex = that.data.provinces.findIndex(ele => {
                    return ele.name == provinceName
                })
                if (pIndex != -1) {
                    const e = {
                        detail: {
                            value: pIndex
                        }
                    }
                    that.provinceChange(e, 0, 0).then(() => {
                        // 读取市
                        const cIndex = that.data.cities.findIndex(ele => {
                            return ele.name == cityName
                        })
                        if (cIndex != -1) {
                            const e = {
                                detail: {
                                    value: cIndex
                                }
                            }
                            that.cityChange(e, 0).then(() => {
                                // 读取区县
                                const aIndex = that.data.areas.findIndex(ele => {
                                    return ele.name == areaName;
                                });
                                if (aIndex != -1) {
                                    const e = {
                                        detail: {
                                            value: aIndex
                                        }
                                    };
                                    that.areaChange(e);
                                }
                            })
                        }
                    })
                }
                const addressData = {};
                addressData.linkMan = res.userName;
                addressData.mobile = res.telNumber;
                addressData.address = res.detailInfo;
                that.setData({addressData});
            }
        })
    },
})

