import React from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/product-card.css';

const ProductCard = ({ product: productData }) => {
  const navigate = useNavigate();
  const product = productData;

  const handleClick = () => {
    navigate(`/product/${product.product.id}`, {
      state: { product }
    });
  };

  return (
    <div className="product-card" onClick={handleClick} style={{ cursor: 'pointer' }}>
      <div className="product-image-container">
        <img
          src={product.product.avatar?.split(',')[0] || '/placeholder-product.png'}
          alt={product.product.describe || '商品图片'}
          className="product-image"
        />
      </div>
      <div className="product-info">
        <div className="product-header">
          <img
            src={product.user?.avatar || '/placeholder-user.png'}
            alt={product.user?.username || '用户'}
            className="user-avatar"
          />
          <span className="username">{product.user.username}</span>
        </div>
        <h3 className="product-price">¥{product.product.price || '0'}</h3>
        {product.product.originalPrice && product.product.originalPrice > 0 && (
          <p className="original-price">¥{product.product.originalPrice}</p>
        )}
        <p className="product-desc">{product.product.describe || '暂无描述'}</p>
        <div className="product-meta">
          <span className="location">{product.product.location}</span>
          {product.product.canSelfPickup && (
            <span className="pickup-tag">可自提</span>
          )}
        </div>
      </div>
    </div>
  );
};

export default ProductCard;