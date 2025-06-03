import React, { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Select, Space, Empty, Spin, message } from 'antd';
import styled from 'styled-components';
import ProductService from '../services/product';
import ProductCard from '../components/ProductCard';
import Header from '../components/Header';

const Container = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
`;

const FilterSection = styled.div`
  margin-bottom: 24px;
  padding: 16px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
`;

const ProductGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin-top: 20px;
`;

const LoadMore = styled.div`
  text-align: center;
  margin-top: 20px;
  padding: 10px 0;
  cursor: pointer;
  color: #1890ff;
  &:hover {
    opacity: 0.8;
  }
`;

const SearchResults = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [products, setProducts] = useState([]);
    const [sortOption, setSortOption] = useState({ field: 'create_at', decs: true });
    const [page, setPage] = useState(1);
    const [hasMore, setHasMore] = useState(true);

    // 从 URL 获取搜索关键词
    const searchParams = new URLSearchParams(location.search);
    const keyword = searchParams.get('keyword') || '';

    // 搜索商品
    useEffect(() => {
        if (keyword) {
            searchProducts();
        }
    }, [keyword, sortOption]);

    const searchProducts = async (isLoadMore = false) => {
        if (!isLoadMore) {
            setPage(1);
        }
        setLoading(true);

        try {
            const searchParams = {
                keyword,
                sort: sortOption,
                page: isLoadMore ? page : 1,
                size: 20
            };

            const response = await ProductService.searchProducts(searchParams);
            const newProducts = response.data?.products || [];
            if (isLoadMore) {
                setProducts(prev => [...prev, ...newProducts]);
            } else {
                setProducts(newProducts);
            }

            setHasMore(response.data?.hasMore || false);
            if (isLoadMore) {
                setPage(prev => prev + 1);
            }
        } catch (error) {
            message.error('搜索失败');
        } finally {
            setLoading(false);
        }
    };

    const sortOptions = [
        { label: '最新发布', value: 'create_at,true' },
        { label: '价格从低到高', value: 'price,false' },
        { label: '价格从高到低', value: 'price,true' }
    ];

    const handleSortChange = (value) => {
        const [field, decs] = value.split(',');
        setSortOption({ field, decs: decs === 'true' });
    };

    const handleSearch = (keyword) => {
        if (keyword.trim()) {
            navigate(`/search?keyword=${encodeURIComponent(keyword.trim())}`);
        }
    };

    return (
        <div className="home-container">
            <Header onSearch={handleSearch} />
            <Container>
                <FilterSection>
                    <Space direction="vertical" style={{ width: '100%' }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                            <span>搜索结果：{keyword}</span>
                            <Select
                                style={{ width: 150 }}
                                placeholder="排序方式"
                                options={sortOptions}
                                onChange={handleSortChange}
                                defaultValue="create_at,true"
                            />
                        </div>
                    </Space>
                </FilterSection>

                <Spin spinning={loading}>
                    {products.length > 0 ? (
                        <ProductGrid>
                            {products.map((item, index) => (
                                <ProductCard
                                    key={`${item.product.id}-${index}`}
                                    product={item}
                                />
                            ))}
                        </ProductGrid>
                    ) : (
                        <Empty description="暂无相关商品" />
                    )}
                </Spin>

                {hasMore && !loading && (
                    <LoadMore onClick={() => searchProducts(true)}>
                        加载更多
                    </LoadMore>
                )}
            </Container>
        </div>
    );
};

export default SearchResults; 