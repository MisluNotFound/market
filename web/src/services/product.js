import axios from 'axios';

const API_BASE_URL = 'http://localhost:3200/api';

const ProductService = {
  // 创建商品
  createProduct: async (formData, userId) => {
    console.log(formData)
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.post(`${API_BASE_URL}/product/${userId}`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  // 获取单个商品
  getProduct: async (userId, productId) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.get(`${API_BASE_URL}/${userId}/${productId}`, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  // 更新商品
  updateProduct: async (userId, productId, formData) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.put(`${API_BASE_URL}/${userId}/${productId}`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  // 下架商品
  offShelves: async (userId, productId) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.put(`${API_BASE_URL}/product/${userId}/${productId}/off-shelves`, {}, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  // 上架商品
  onShelves: async (userId, productId) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.put(`${API_BASE_URL}/product/${userId}/${productId}/on-shelves`, {}, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  // 标记为已售出
  soldOut: async (userId, productId) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.put(`${API_BASE_URL}/product/${userId}/${productId}/sold`, {}, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  selling: async (userId, productId) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.put(`${API_BASE_URL}/product/${userId}/${productId}/selling`, {}, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },
  // 获取商品列表(带分页)
  getProductList: async (page = 1, size = 10) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.get(`${API_BASE_URL}/product/products`, {
        params: {
          page,
          size
        },
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  getProducts: async () => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.get(`${API_BASE_URL}/products`, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  getProductDetail: async (userId, productId) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.get(`${API_BASE_URL}/product/${userId}/${productId}`, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  searchProducts: async (searchParams) => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.post(`${API_BASE_URL}/search/products`, {
        keyword: searchParams.keyword,
        categories: searchParams.categories || [],
        attributes: [],
        sort: {
          field: searchParams.sort?.field || 'createTime',
          decs: searchParams.sort?.decs || true
        }
      }, {
        params: {
          page: searchParams.page || 1,
          size: searchParams.size || 20
        },
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || { message: '搜索失败' };
    }
  },

  // 获取商品分类
  getCategories: async () => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.get(`${API_BASE_URL}/product/category`, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  // 获取用户商品列表
  getUserProducts: async (page = 1, size = 10) => {
    try {
      const userId = localStorage.getItem("userId");
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.get(`${API_BASE_URL}/product/${userId}`, {
        params: { page, size },
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  // 获取用户收藏商品
  getUserFavorites: async (page = 1, size = 10) => {
    try {
      const userId = localStorage.getItem("userId");
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.get(`${API_BASE_URL}/product/${userId}/favorites`, {
        params: { page, size },
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  // 获取标签列表
  getTags: async () => {
    try {
      const accessToken = localStorage.getItem('accessToken');
      const response = await axios.get(`${API_BASE_URL}/product/tags`, {
        headers: {
          Authorization: `Bearer ${accessToken}`
        }
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || { message: '获取标签失败' };
    }
  },
}

export default ProductService