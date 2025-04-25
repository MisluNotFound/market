import React from 'react';
import { useParams, useLocation, useNavigate } from 'react-router-dom';
import ProductService from '../services/product';
import AuthService from '../services/auth';
import OrderService from '../services/order';
import '../styles/product-detail.css';

const ProductDetail = () => {
  const { id } = useParams();
  const location = useLocation();
  const navigate = useNavigate();
  const [product, setProduct] = React.useState(location.state?.product || null);
  const [loading, setLoading] = React.useState(!location.state?.product);
  const [currentUser, setCurrentUser] = React.useState(null);

  React.useEffect(() => {
    const fetchCurrentUser = async () => {
      const user = await AuthService.getCurrentUser();
      setCurrentUser(user);
    };
    fetchCurrentUser();
  }, []);

  // React.useEffect(() => {
  //   if (!location.state?.product) {
  //     const fetchProduct = async () => {
  //       try {
  //         const response = await ProductService.getProductDetail(id);
  //         setProduct(response.data);
  //       } catch (error) {
  //         console.error('获取商品详情失败:', error);
  //       } finally {
  //         setLoading(false);
  //       }
  //     };

  //     fetchProduct();
  //   }
  // }, [id, location.state]);

  if (loading) return <div className="loading">加载中...</div>;
  if (!product) return <div className="error">商品不存在</div>;
  return (
    <div className="product-detail-container">
      <div className="user-section">
        <div className="seller-info">
          <img
            src={product.user?.avatar || '/placeholder-user.png'}
            className="seller-avatar"
          />
          <span className="seller-name">{product.user?.username || '未知用户'}</span>
        </div>
      </div>

      <div className="product-images">
        {product.product.avatar.split(',').map((imgUrl, index) => (
          <img
            key={index}
            src={imgUrl.trim()}
            alt={`商品图片 ${index + 1}`}
          />
        ))}
      </div>

      <div className="product-info">
        <h1 className="product-title">{product.title}</h1>
        <div className="price-section">
          <span className="current-price">¥{product.product.price}</span>
          {product.product.originalPrice && (
            <span className="original-price">¥{product.product.originalPrice}</span>
          )}
          <span
            className={`shipping-price ${product.product.shippingPrice === 0 ? 'free' : ''}`}
            style={{ marginLeft: '10px' }}
          >
            {product.product.shippingPrice === 0 ? '包邮' : `运费: ¥${product.product.shippingPrice}`}
          </span>
        </div>
        <div className="product-meta">
          <span>发布时间: {new Date(product.product.publishAt).toLocaleDateString()}</span>
          {/* <span>浏览: {product.product.viewCount}次</span> */}
        </div>
        <div className="product-description">
          <p>{product.product.describe}</p>
        </div>
      </div>

      {currentUser && currentUser.id !== product.user?.id && (
        <div className="action-bar">
          <button className="contact-btn">联系卖家</button>
          <button
            className="buy-btn"
            onClick={async () => {
              try {
                await OrderService.purchaseProduct(
                  currentUser.id,
                  product.product.id,
                  product.product.price, // totalAmount
                  product.product.shippingPrice // shipAmount
                );
                alert('购买成功！');
                navigate('/user-center');
              } catch (error) {
                alert(`购买失败: ${error.message}`);
              }
            }}
          >
            立即购买
          </button>
        </div>
      )}
    </div>
  );
};

export default ProductDetail;