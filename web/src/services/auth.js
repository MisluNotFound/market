import UserService from './user';
import { IMService } from './im';

let imServiceInstance = null;

const AuthService = {
  // 登录方法
  login: async (phone, password) => {
    try {
      const response = await UserService.login({ phone, password });
      localStorage.setItem('accessToken', response.accessToken);
      localStorage.setItem('userId', response.userID);
      localStorage.setItem('refreshToken', response.refreshToken);
      // 初始化IMService
      imServiceInstance = new IMService(response.userID);
      return response;
    } catch (error) {
      throw error;
    }
  },

  // 注册方法
  register: async (userData) => {
    try {
      const response = await UserService.register(userData);
      return response;
    } catch (error) {
      throw error;
    }
  },

  // 获取当前用户信息
  getCurrentUser: async () => {
    const userId = localStorage.getItem('userId');
    if (!userId) return null;

    try {
      const user = await UserService.getUserInfo(userId);
      // 确保IMService已初始化
      if (!imServiceInstance) {
        imServiceInstance = new IMService(userId);
      }
      return user;
    } catch (error) {
      return null;
    }
  },

  // 登出方法
  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('userId');
    // 关闭IMService连接
    if (imServiceInstance) {
      imServiceInstance.close();
      imServiceInstance = null;
    }
  },

  // 检查是否已认证
  isAuthenticated: () => {
    return !!localStorage.getItem('accessToken');
  },

  // 获取认证token
  getToken: () => {
    return localStorage.getItem('accessToken');
  },

  // 更新用户基本信息
  updateBasicInfo: async (userId, userData) => {
    try {
      const response = await UserService.updateBasic(userId, userData);
      return response;
    } catch (error) {
      throw error;
    }
  },

  // 更新用户密码
  updatePassword: async (userId, passwordData) => {
    try {
      const response = await UserService.updatePassword(userId, passwordData);
      return response;
    } catch (error) {
      throw error;
    }
  },

  // 上传用户头像
  uploadAvatar: async (userId, avatarFile) => {
    try {
      const response = await UserService.uploadAvatar(userId, avatarFile);
      return response;
    } catch (error) {
      throw error;
    }
  },

  // 获取IMService实例
  getIMService: () => {
    return imServiceInstance;
  }
};

export default AuthService;