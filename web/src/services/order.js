import axios from 'axios';

const API_BASE_URL = 'http://localhost:3200/api/order';

const OrderService = {
  // 购买商品
  purchaseProduct: async (userId, productId, totalAmount, shipAmount = 0) => {

    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.post(`${API_BASE_URL}/${userId}/${productId}`, {
        totalAmount,
        shipAmount
      }, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || '购买商品失败');
    }
  },

  // 获取订单列表
  getOrderList: async (userId, page, isBought) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.get(`${API_BASE_URL}/${userId}/list`, {
        params: {
          page,
          isBought
        },
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return {
        data: response.data.data
      };
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取订单列表失败');
    }
  },

  // 获取订单详情
  getOrder: async (userId, orderId) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.get(`${API_BASE_URL}/${userId}/${orderId}`, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取订单详情失败');
    }
  },

  // 确认发货
  confirmShipped: async (userId, orderId, refund) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.put(`${API_BASE_URL}/shipped/${userId}/${orderId}`, { refund }, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || '确认发货失败');
    }
  },

  // 确认签收
  confirmSigned: async (userId, orderId, refund) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const formData = new URLSearchParams();
      formData.append('refund', refund);

      const response = await axios.put(
        `${API_BASE_URL}/signed/${userId}/${orderId}`,
        formData,
        {
          headers: {
            'Authorization': `Bearer ${accessToken}`,
            'Content-Type': 'application/x-www-form-urlencoded'
          }
        }
      );
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || '确认签收失败');
    }
  },

  // 支付订单
  payOrder: async (orderId) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const userId = localStorage.getItem('userId')
      console.log(userId, orderId)
      const response = await axios.put(`${API_BASE_URL}/pay/${userId}/${orderId}`, {}, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || '支付订单失败');
    }
  },

  // 获取所有订单状态
  getAllOrderStatus: async (userId) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.get(`${API_BASE_URL}/${userId}/status`, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取订单状态失败');
    }
  },

  // 退款
  refundOrder: async (userId, orderId) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.post(`${API_BASE_URL}/refund/${userId}/${orderId}`, {}, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || '退款失败');
    }
  },

  getOrderDetail: async (orderId, userId) => {
    try {
      console.log(orderId, userId)
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.get(`${API_BASE_URL}/${userId}/${orderId}`, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取订单详情失败');
    }
  },

  createOrder: async (orderData) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.post(`${API_BASE_URL}`, orderData, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || '创建订单失败');
    }
  },

  cancelOrder: async (orderId) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const userId = localStorage.getItem('userId')
      const response = await axios.put(`${API_BASE_URL}/cancel/${userId}/${orderId}`, {}, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || '取消订单失败');
    }
  },

  // 获取订单评论
  getOrderComments: async (orderId, page = 1, size = 10) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const userId = localStorage.getItem('userId');
      const response = await axios.get(`${API_BASE_URL}/comment/${orderId}`, {
        params: { page, size },
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取评论失败');
    }
  },

  // 创建订单评论
  createOrderComment: async (orderId, comment, isGood, pics = '') => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const userId = localStorage.getItem('userId');
      const response = await axios.post(`${API_BASE_URL}/comment/${orderId}`, {
        comment,
        isGood,
        pics
      }, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || '提交评论失败');
    }
  },

  // 回复订单评论
  replyOrderComment: async (orderId, commentId, comment) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const userId = localStorage.getItem('userId');
      const response = await axios.post(`${API_BASE_URL}/comment/${orderId}/reply`, {
        commentID: commentId,
        comment
      }, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || '回复评论失败');
    }
  },

  // 获取未评价订单
  getUncommentOrders: async (userId) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.get(`${API_BASE_URL}/${userId}/uncomment`, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return {
        data: response.data.data
      };
    } catch (error) {
      throw new Error(error.response?.data?.message || '获取未评价订单失败');
    }
  },

  // 获取订单状态
  async getOrderStatus(orderId) {
    const response = await axios.get(`${API_BASE_URL}/status/${orderId}`);
    return response.data;
  }
};

export default OrderService;