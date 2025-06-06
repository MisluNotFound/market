import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { FaHeart } from 'react-icons/fa';
import { message } from 'antd';
import ProductService from '../services/product';
import '../styles/product-card.css';

const getReputationInfo = (reputation) => {
  if (!reputation || reputation === 0) return null;
  if (reputation > 80) {
    return { text: '信誉极好', className: 'excellent' };
  } else if (reputation >= 60) {
    return { text: '信誉良好', className: 'good' };
  } else {
    return { text: '信誉较差', className: 'poor' };
  }
};

const ProductCard = ({ product: productData }) => {
  const navigate = useNavigate();
  const product = productData;
  const [isLiked, setIsLiked] = useState(product.isLiked || false);
  const [isLoading, setIsLoading] = useState(false);

  const handleClick = () => {
    navigate(`/product/${product.product.id}`, {
      state: { product }
    });
  };

  const handleFavoriteClick = async (e) => {
    e.stopPropagation();
    if (isLoading) return;

    try {
      setIsLoading(true);
      const userId = localStorage.getItem('userId');
      if (!userId) {
        message.warning('请先登录');
        return;
      }

      if (isLiked) {
        await ProductService.dislikeProduct(userId, product.product.id);
        setIsLiked(false);
        message.success('已取消收藏');
      } else {
        await ProductService.likeProduct(userId, product.product.id);
        setIsLiked(true);
        message.success('收藏成功');
      }
    } catch (error) {
      message.error(error.message || '操作失败');
    } finally {
      setIsLoading(false);
    }
  };

  const reputationInfo = product.credit ? getReputationInfo(product.credit.reputation) : null;

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
          <div className="user-info">
            <span className="username">{product.user.username}</span>
            {reputationInfo && (
              <span className={`reputation-tag ${reputationInfo.className}`}>
                {reputationInfo.text}
              </span>
            )}
          </div>
        </div>
        <div className="product-price">
          <span className="current-price">¥{product.product.price || '0'}</span>
          {product.product.originalPrice && product.product.originalPrice > 0 && (
            <span className="original-price">¥{product.product.originalPrice}</span>
          )}
        </div>
        <p className="product-desc">{product.product.describe || '暂无描述'}</p>
        <div className="product-meta">
          <span className="location">{product.address}</span>
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
      <div
        className={`favorite-icon ${isLiked ? 'liked' : ''}`}
        onClick={handleFavoriteClick}
        style={{ cursor: isLoading ? 'not-allowed' : 'pointer' }}
      >
        <FaHeart />
      </div>
    </div>
  );
};

export default ProductCard;