import React, { useState, useEffect, useRef } from 'react';
import { useParams } from 'react-router-dom';
import ProductCard from '../components/ProductCard';
import UserService from '../services/user';
import ProductService from '../services/product';
import '../styles/user-profile.css';

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

const UserProfile = () => {
    const { userId } = useParams();
    const [user, setUser] = useState(null);
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [loadingMore, setLoadingMore] = useState(false);
    const [pageInfo, setPageInfo] = useState({
        page: 1,
        size: 10,
        total: 0
    });
    const isFetching = useRef(false);

    // 获取用户信息
    useEffect(() => {
        const fetchUserInfo = async () => {
            try {
                const userData = await UserService.getUserInfo(userId);
                setUser(userData);
            } catch (error) {
                console.error('获取用户信息失败:', error);
            }
        };

        fetchUserInfo();
    }, [userId]);

    // 获取用户商品
    const fetchProducts = async (isLoadMore = false) => {
        if (isFetching.current) return;
        isFetching.current = true;

        try {
            if (isLoadMore) {
                setLoadingMore(true);
            } else {
                setLoading(true);
            }

            const data = await ProductService.getProductList(pageInfo.page, pageInfo.size);
            const newProducts = Array.isArray(data.data?.products) ? data.data.products : [];

            if (newProducts.length > 0) {
                setProducts(prev =>
                    isLoadMore ? [...prev, ...newProducts] : newProducts
                );
            }

            setPageInfo(prev => ({
                ...prev,
                total: data.total || 0
            }));
        } catch (error) {
            console.error('获取商品列表失败:', error);
        } finally {
            setLoadingMore(false);
            setLoading(false);
            isFetching.current = false;
        }
    };

    useEffect(() => {
        if (products.length === 0) {
            fetchProducts();
        }
    }, []);

    const handleLoadMore = () => {
        setPageInfo(prev => ({ ...prev, page: prev.page + 1 }));
        fetchProducts(true);
    };

    const reputationInfo = user?.reputation ? getReputationInfo(user.reputation) : null;

    if (!user) return <div className="loading">加载中...</div>;

    return (
        <div className="user-profile-container">
            {/* 用户信息区域 */}
            <div className="user-info-section">
                <div className="user-info-header">
                    <img
                        src={user.avatar || '/placeholder-user.png'}
                        alt={user.username}
                        className="user-avatar"
                    />
                    <div className="user-info-content">
                        <div className="user-info-main">
                            <h2 className="username">{user.username}</h2>
                            {reputationInfo && (
                                <span className={`reputation-tag ${reputationInfo.className}`}>
                                    {reputationInfo.text}
                                </span>
                            )}
                        </div>
                        {reputationInfo && user.reputation > 0 && (
                            <div className="user-credit-info">
                                <div className="credit-item">
                                    <span className="label">总评价</span>
                                    <span className="value">{user.totalComment}</span>
                                </div>
                                <div className="credit-item">
                                    <span className="label">好评</span>
                                    <span className="value positive">{user.positiveComment}</span>
                                </div>
                                <div className="credit-item">
                                    <span className="label">差评</span>
                                    <span className="value negative">{user.negativeComment}</span>
                                </div>
                                <div className="credit-item">
                                    <span className="label">信誉值</span>
                                    <span className={`value ${reputationInfo.className}`}>
                                        {user.reputation}
                                    </span>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </div>

            {/* 商品列表区域 */}
            <div className="user-products-section">
                <h3 className="section-title">发布的商品</h3>
                {loading ? (
                    <div className="loading">加载中...</div>
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
                        {!loadingMore && products.length > 0 && pageInfo.page * pageInfo.size < pageInfo.total && (
                            <div className="load-more-container">
                                <button onClick={handleLoadMore} className="load-more-btn">
                                    加载更多
                                </button>
                            </div>
                        )}
                    </>
                )}
            </div>
        </div>
    );
};

export default UserProfile; 