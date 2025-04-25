import axios from 'axios';

const API_BASE_URL = 'http://localhost:3200/api/user';

const handleResponse = (response) => {
  if (response.data.code !== 200) {
    throw new Error(response.data.msg);
  }
  return response.data.data;
};

const UserService = {
  // 注册用户
  register: async (userData) => {
    try {
      const response = await axios.post(`${API_BASE_URL}/register`, userData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      });
      return handleResponse(response);
    } catch (error) {
      throw new Error(error.response?.data?.msg || '注册失败');
    }
  },

  // 用户登录
  login: async (credentials) => {
    try {
      const response = await axios.post(`${API_BASE_URL}/login`, credentials);
      return handleResponse(response);
    } catch (error) {
      throw new Error(error.response?.data?.msg || '登录失败');
    }
  },

  // 获取用户信息
  getUserInfo: async (userId) => {
    try {
      const response = await axios.get(`${API_BASE_URL}/${userId}`);
      return handleResponse(response);
    } catch (error) {
      throw new Error(error.response?.data?.msg || '获取用户信息失败');
    }
  },

  // 上传头像
  uploadAvatar: async (userId, avatarFile) => {
    try {
      const formData = new FormData();
      formData.append('avatar', avatarFile);
      
      const response = await axios.put(`${API_BASE_URL}/${userId}/avatar`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      });
      return handleResponse(response);
    } catch (error) {
      throw new Error(error.response?.data?.msg || '上传头像失败');
    }
  },

  // 更新基本信息
  updateBasic: async (userId, userData) => {
    try {
      const response = await axios.put(`${API_BASE_URL}/${userId}/basic`, userData);
      return handleResponse(response);
    } catch (error) {
      throw new Error(error.response?.data?.msg || '更新基本信息失败');
    }
  },

  // 更新密码
  updatePassword: async (userId, passwordData) => {
    try {
      const response = await axios.put(`${API_BASE_URL}/${userId}/password`, passwordData);
      return handleResponse(response);
    } catch (error) {
      throw new Error(error.response?.data?.msg || '更新密码失败');
    }
  }
};

export default UserService;