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
  const [loading, setLoading] = React.useState(true);
  const [currentUser, setCurrentUser] = React.useState(null);
  const [error, setError] = React.useState(null);

  React.useEffect(() => {
    const fetchData = async () => {
      try {
        const user = await AuthService.getCurrentUser();
        setCurrentUser(user);
        console.log(location.state.product)
        // // 如果location.state中有商品数据，直接使用
        // if (location.state?.product) {
        //   setProduct(location.state.product);
        //   setLoading(false);
        //   return;
        // }

        // 否则调用接口获取商品详情
        const response = await ProductService.getProductDetail(location.state.product.user.id, id);
        setProduct(response.data);
      } catch (error) {
        console.error('获取商品详情失败:', error);
        setError(error.message || '获取商品详情失败');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id, location.state]);

  if (loading) return <div className="loading">加载中...</div>;
  if (error) return <div className="error">{error}</div>;
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
          <span>状态: {product.product.isSelling ? '在售' : '已下架'}{product.product.isSold ? ' (已售罄)' : ''}</span>
        </div>
        <div className="product-description">
          <p>{product.product.describe}</p>
        </div>
      </div>

      <div className="action-bar">
        {currentUser && currentUser.id !== product.user?.id && (
          <>
            <button
              className="contact-btn"
              onClick={async () => {
                try {
                  const response = await fetch('http://localhost:3200/api/conversation/create', {
                    method: 'POST',
                    headers: {
                      'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: new URLSearchParams({
                      fromUserID: currentUser.id,
                      toUserID: product.user?.id,
                      productID: product.product.id
                    })
                  });
                  const result = await response.json();
                  if (result.code !== 200) {
                    throw new Error(result.msg || '创建会话失败');
                  }
                  navigate('/chat', {
                    state: {
                      fromUserID: currentUser.id,
                      toUserID: product.user?.id,
                      productID: product.product.id
                    }
                  });
                } catch (error) {
                  console.error('创建对话失败:', error);
                  alert('创建对话失败，请稍后重试');
                }
              }}
            >
              联系卖家
            </button>
            <button
              className="buy-btn"
              onClick={async () => {
                try {
                  await OrderService.purchaseProduct(
                    currentUser.id,
                    product.product.id,
                    product.product.price,
                    product.product.shippingPrice
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
          </>
        )}
        {currentUser && currentUser.id === product.user?.id && (
          <div className="manage-container">
            <button
              className="manage-btn"
              onClick={(e) => {
                e.stopPropagation();
                document.querySelector('.manage-dropdown').classList.toggle('show');
              }}
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <circle cx="12" cy="12" r="1"></circle>
                <circle cx="12" cy="5" r="1"></circle>
                <circle cx="12" cy="19" r="1"></circle>
              </svg>
              管理
            </button>
            <div className="manage-dropdown">
              <button
                className="dropdown-item"
                onClick={() => {
                  console.log('正在导航到编辑页面，商品ID:', product.product.id);
                  navigate(`/edit-product/${product.product.id}`);
                }}
              >
                编辑商品
              </button>
              {product.product.isSelling ? (
                <button
                  className="dropdown-item"
                  onClick={async () => {
                    try {
                      await ProductService.offShelves(currentUser.id, product.product.id);
                      setProduct(prev => ({
                        ...prev,
                        product: {
                          ...prev.product,
                          isSelling: false
                        }
                      }));
                    } catch (error) {
                      console.error('下架商品失败:', error);
                    }
                  }}
                >
                  下架商品
                </button>
              ) : (
                <button
                  className="dropdown-item"
                  onClick={async () => {
                    try {
                      await ProductService.onShelves(currentUser.id, product.product.id);
                      setProduct(prev => ({
                        ...prev,
                        product: {
                          ...prev.product,
                          isSelling: true
                        }
                      }));
                    } catch (error) {
                      console.error('上架商品失败:', error);
                    }
                  }}
                >
                  上架商品
                </button>
              )}
              {product.product.isSold ? (
                <button
                  className="dropdown-item"
                  onClick={async () => {
                    try {
                      await ProductService.selling(currentUser.id, product.product.id);
                      setProduct(prev => ({
                        ...prev,
                        product: {
                          ...prev.product,
                          isSold: false,
                          isSelling: true
                        }
                      }));
                    } catch (error) {
                      console.error('恢复有货失败:', error);
                    }
                  }}
                >
                  恢复有货
                </button>
              ) : (
                <button
                  className="dropdown-item danger"
                  onClick={async () => {
                    try {
                      await ProductService.soldOut(currentUser.id, product.product.id);
                      setProduct(prev => ({
                        ...prev,
                        product: {
                          ...prev.product,
                          isSold: true,
                          isSelling: false
                        }
                      }));
                    } catch (error) {
                      console.error('标记售罄失败:', error);
                    }
                  }}
                >
                  标记售罄
                </button>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ProductDetail;