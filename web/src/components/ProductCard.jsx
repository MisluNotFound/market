import React from 'react';
import { useNavigate } from 'react-router-dom';
import { FaHeart } from 'react-icons/fa';
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
        {(product.product.isSold || !product.product.isSelling) && (
          <div className="product-status-overlay">
            <span className={`status-text ${product.product.isSold ? 'sold' : 'off-shelf'}`}>
              {product.product.isSold ? '已售罄' : '已下架'}
            </span>
          </div>
        )}
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
        <div className="product-price">
          <span className="current-price">¥{product.product.price || '0'}</span>
          {product.product.originalPrice && product.product.originalPrice > 0 && (
            <span className="original-price">¥{product.product.originalPrice}</span>
          )}
        </div>
        <p className="product-desc">{product.product.describe || '暂无描述'}</p>
        <div className="product-meta">
          <span className="location">{product.product.location}</span>
          {product.product.canSelfPickup && (
            <span className="pickup-tag">可自提</span>
          )}
          <span className={`product-condition ${product.product.condition === '使用过' ? 'used' : 'new'}`}>
            {product.product.condition === '使用过'
              ? `已使用${product.product.usedTime || '未知时长'}`
              : product.product.condition}
          </span>
        </div>
      </div>
      <div className={`favorite-icon ${product.isLiked ? 'liked' : ''}`}>
        <FaHeart />
      </div>
    </div>
  );
};

export default ProductCard;