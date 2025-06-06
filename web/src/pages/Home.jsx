import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import ProductCard from '../components/ProductCard';
import Header from '../components/Header';
import ProductService from '../services/product';
import '../styles/home.css';

const Home = () => {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const navigate = useNavigate();
  const isMounted = useRef(true); // 使用useRef来跟踪组件挂载状态

  // 组件卸载时设置isMounted为false
  useEffect(() => {
    return () => {
      isMounted.current = false;
    };
  }, []);

  useEffect(() => {
    let isActive = true; // 局部变量替代 isMounted

    const fetchProducts = async () => {
      try {
        setLoading(true);
        const response = await ProductService.getProductList(page, 10);

        if (!isActive) return;

        const productList = response?.data?.products || [];

        setProducts(prev => {
          if (page === 1) return productList;
          // 合并所有商品，确保不重复
          const allProducts = [...prev, ...productList];
          const uniqueProducts = Array.from(
            new Map(allProducts.map(item => [item.id, item])).values()
          );
          return uniqueProducts;
        });

        setHasMore(productList.length > 0);
      } catch (error) {
        if (!isActive) return;
        console.error('获取商品列表失败:', error);
        setProducts([]);
      } finally {
        if (isActive) setLoading(false);
      }
    };

    fetchProducts();

    return () => {
      isActive = false; // 仅标记本次 effect 失效
    };
  }, [page]);

  const handleSearch = (keyword) => {
    if (keyword.trim()) {
      navigate(`/search?keyword=${encodeURIComponent(keyword.trim())}`);
    }
  };

  return (
    <div className="home-container">
      <Header onSearch={handleSearch} />
      <div className="product-grid">
        {products.map((product) => (
          <ProductCard key={product.id} product={product} />
        ))}
      </div>
      {loading && <div className="loading">加载中...</div>}
      {!loading && hasMore && (
        <button
          className="load-more"
          onClick={() => setPage(prev => prev + 1)}
        >
          加载更多
        </button>
      )}
      {!loading && products.length === 0 && (
        <div className="no-products">暂无商品数据</div>
      )}
    </div>
  );
};

export default Home;