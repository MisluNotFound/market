import axios from 'axios';

const API_BASE_URL = 'http://localhost:3200/api';

const ProductService = {
  // 创建商品
  createProduct: async (formData, userId) => {
    console.log(formData)
    try {
      const response = await axios.post(`${API_BASE_URL}/product/${userId}`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
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
      const response = await axios.get(`${API_BASE_URL}/${userId}/${productId}`);
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  // 更新商品
  updateProduct: async (userId, productId, formData) => {
    try {
      const response = await axios.put(`${API_BASE_URL}/${userId}/${productId}`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
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
      const response = await axios.put(`${API_BASE_URL}/${userId}/${productId}/off-shelves`);
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  // 上架商品
  onShelves: async (userId, productId) => {
    try {
      const response = await axios.put(`${API_BASE_URL}/${userId}/${productId}/on-shelves`);
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  // 标记为已售出
  soldOut: async (userId, productId) => {
    try {
      const response = await axios.put(`${API_BASE_URL}/${userId}/${productId}/sold`);
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },
  // 获取商品列表(带分页)
  getProductList: async (page = 1, size = 10) => {
    try {
      const response = await axios.get(`${API_BASE_URL}/product/products`, {
        params: {
          page,
          size
        }
      });
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  getProducts: async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/products`);
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  getProductDetail: async (userId, productId) => {
    try {
      const response = await axios.get(`${API_BASE_URL}/product/${userId}/${productId}`);
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },

  searchProducts: async (keyword, page = 1, size = 10) => {
    try {
      const response = await axios.get(`${API_BASE_URL}/product/search`, {
        params: {
          keyword,
          page,
          size
        }
      });
      return response.data;
    } catch (error) {
      throw error.response.data;
    }
  },
}

export default ProductService