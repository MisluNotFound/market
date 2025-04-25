import React, { useState, useEffect, useRef } from 'react';
import ProductCard from '../components/ProductCard';
import Header from '../components/Header';
import ProductService from '../services/product';
import '../styles/home.css';

const Home = () => {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [searchKeyword, setSearchKeyword] = useState('');
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
        let response;
        if (searchKeyword) {
          response = await ProductService.searchProducts(searchKeyword, page, 10);
        } else {
          response = await ProductService.getProductList(page, 10);
        }

        // 检查是否仍需更新
        if (!isActive) return;

        const productList = response?.data?.products || [];
        setProducts(prev => (page === 1 ? productList : [...prev, ...productList]));
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
  }, [page, searchKeyword]);

  const handleSearch = (keyword) => {
    setSearchKeyword(keyword);
    setPage(1); // 搜索时重置页码
  };

  return (
    <div className="home-container">
      <Header onSearch={handleSearch} />
      <div className="product-grid">
        {products.map((product, index) => (
          console.log(product),
          <ProductCard key={`${product.id}-${index}`} product={product} />
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