import axios from 'axios';
import AuthService from './auth';

const API_BASE = 'http://localhost:3200/api/address';

export default {
    // 创建地址
    async createAddress(addressData) {
        const user = await AuthService.getCurrentUser();
        if (!user) throw new Error('用户未登录');
        const accessToken = localStorage.getItem('accessToken');

        return axios.post(`${API_BASE}/${user.id}`, addressData, {
            headers: {
                Authorization: `Bearer ${accessToken}`
            }
        });
    },

    // 更新地址
    async updateAddress(addressId, addressData) {
        const accessToken = localStorage.getItem('accessToken');
        return axios.put(`${API_BASE}/${addressId}`, addressData, {
            headers: {
                Authorization: `Bearer ${accessToken}`
            }
        });
    },

    // 获取地址列表
    async getAddresses(page = 1, pageSize = 10) {
        const user = await AuthService.getCurrentUser();
        if (!user) throw new Error('用户未登录');
        const accessToken = localStorage.getItem('accessToken');

        return axios.get(`${API_BASE}/${user.id}`, {
            params: { page, pageSize },
            headers: {
                Authorization: `Bearer ${accessToken}`
            }
        });
    },

    // 删除地址
    async deleteAddress(addressId) {
        const accessToken = localStorage.getItem('accessToken');
        return axios.delete(`${API_BASE}/${addressId}`, {
            headers: {
                Authorization: `Bearer ${accessToken}`
            }
        });
    },

    async setDefaultAddress(addressId, isDefault) {
        const accessToken = localStorage.getItem('accessToken');
        return axios.put(`${API_BASE}/default/${addressId}`, { isDefault }, {
            headers: {
                Authorization: `Bearer ${accessToken}`
            }
        });
    }
};