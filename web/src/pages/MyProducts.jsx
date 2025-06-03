import { useState, useEffect, useRef } from 'react';
import { useParams } from 'react-router-dom';
import ProductCard from '../components/ProductCard';
import ProductService from '../services/product';
import '../styles/home.css';

const MyProducts = () => {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [pageInfo, setPageInfo] = useState({
    page: 1,
    size: 10,
    total: 0
  });
  const isFetching = useRef(false);

  const fetchProducts = async (isLoadMore = false) => {
    if (isFetching.current) return;
    isFetching.current = true;
    let isActive = true;
    try {
      if (isLoadMore) {
        setLoadingMore(true);
      } else {
        setLoading(true);
      }

      const data = await ProductService.getUserProducts(
        pageInfo.page,
        pageInfo.size
      );
      if (!isActive) return;

      const newProducts = Array.isArray(data.data?.products) ? data.data.products : [];
      console.log(newProducts)
      if (newProducts.length > 0) {
        setProducts(prev =>
          isLoadMore
            ? [...prev, ...newProducts]
            : newProducts
        );
      }
      setPageInfo(prev => ({
        ...prev,
        total: data.total || 0
      }));
    } catch (error) {
      if (!isActive) return;
      console.error('获取商品列表失败:', error);
    } finally {
      if (isActive) {
        if (isLoadMore) {
          setLoadingMore(false);
        } else {
          setLoading(false);
        }
      }
      isFetching.current = false;
    }
  };

  useEffect(() => {
    let isActive = true;
    // 只在初始加载且没有数据时请求
    if (products.length === 0) {
      fetchProducts();
    }
    return () => {
      isActive = false;
    };
  }, []);

  const handleLoadMore = () => {
    setPageInfo(prev => ({ ...prev, page: prev.page + 1 }));
    fetchProducts(true);
  };

  return (
    <div className="home-container">
      <h2>我的商品</h2>
      {loading ? (
        <div>加载中...</div>
      ) : (
        <>
          <div className="product-grid">
            {products.length > 0 ? (
              products.map(product => (
                <ProductCard key={product.id} product={product} />
              ))
            ) : (
              <div className="empty-message">暂无商品</div>
            )}
          </div>
          {loadingMore && <div className="loading-more">加载中...</div>}
          {!loadingMore &&
            products.length > 0 &&
            (
              <div style={{ marginBottom: '80px', textAlign: 'center' }}>
                <button onClick={handleLoadMore} className="load-more-btn">
                  加载更多
                </button>
              </div>
            )}
        </>
      )}
    </div>
  );
};

export default MyProducts;