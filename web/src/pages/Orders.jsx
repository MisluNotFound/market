import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import { useNavigate, useParams } from 'react-router-dom';
import { Table, Button, message, Tag } from 'antd';
import AuthService from '../services/auth';
import OrderService from '../services/order';

const Container = styled.div`
  max-width: 1200px;
  margin: 50px auto;
  padding: 20px;
`;

const Title = styled.h2`
  margin-bottom: 24px;
`;

const Orders = () => {
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [orders, setOrders] = useState([]);
  const [hasMore, setHasMore] = useState(true);
  const [page, setPage] = useState(1);
  const navigate = useNavigate();
  const { type } = useParams();

  useEffect(() => {
    if (!AuthService.isAuthenticated()) {
      message.warning('请先登录');
      navigate('/login');
      return;
    }
    fetchOrders();
  }, [navigate]);

  const fetchOrders = async (isLoadMore = false) => {
    if (isLoadMore) {
      setLoadingMore(true);
    } else {
      setLoading(true);
      setPage(1);
    }

    try {
      const user = await AuthService.getCurrentUser();
      if (!user) {
        message.warning('请先登录');
        navigate('/login');
        return;
      }
      const response = await OrderService.getOrderList(user.id, page, type === 'bought');
      const newOrders = Array.isArray(response?.data?.orders) ? response.data.orders : [];
      console.log(newOrders)
      if (isLoadMore) {
        setOrders(prev => [...prev, ...newOrders]);
      } else {
        setOrders(newOrders);
      }

      setHasMore(response?.data?.hasMore || false);
      setPage(prev => prev + 1);
    } catch (error) {
      message.error(error.message);
      setOrders([]);
    } finally {
      setLoading(false);
    }
  };

  const loadMore = () => {
    if (hasMore && !loadingMore) {
      fetchOrders(true);
    }
  };

  const handleViewDetail = (order) => {
    if (!order?.id) {
      message.warning('订单ID无效');
      return;
    }
    navigate(`/order/${order.id}`, {
      state: { order }
    });
  };

  const columns = [
    {
      title: '订单编号',
      dataIndex: ['order', 'id'],
      key: 'id',
    },
    {
      title: '商品名称',
      dataIndex: ['product', 'describe'],
      key: 'productName',
    },
    {
      title: '价格',
      dataIndex: ['order', 'totalAmount'],
      key: 'price',
      render: (price) => price ? `¥${price.toFixed(2)}` : '¥0.00'
    },
    {
      title: '卖家',
      dataIndex: ['user', 'username'],
      key: 'seller',
      render: (username, record) => username || record.user?.username || '未知'
    },
    {
      title: '状态',
      dataIndex: ['order', 'status'],
      key: 'status',
      render: (status) => {
        const statusMap = {
          1: { text: '待付款', color: 'orange' },
          2: { text: '已付款', color: 'blue' },
          3: { text: '已发货', color: 'geekblue' },
          4: { text: '已完成', color: 'green' },
          5: { text: '已退款', color: 'red' },
          6: { text: '已退货', color: 'volcano' },
          7: { text: '已关闭', color: 'default' },
          8: { text: '已取消', color: 'gray' }
        };
        const statusInfo = statusMap[status] || { text: '未知状态', color: 'default' };
        return <Tag color={statusInfo.color}>{statusInfo.text}</Tag>;
      }
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Button
          type="link"
          onClick={() => handleViewDetail(record.order)}
          disabled={!record.order?.id}
        >
          查看详情
        </Button>
      ),
    },
  ];

  return (
    <Container>
      <Title>{type === 'bought' ? '我购买的订单' : '我售出的订单'}</Title>
      <Table
        columns={columns}
        dataSource={orders}
        loading={loading}
        rowKey="id"
        pagination={false}
      />
      {hasMore && (
        <div style={{ textAlign: 'center', marginTop: 16 }}>
          <Button
            onClick={loadMore}
            loading={loadingMore}
            disabled={loadingMore}
          >
            {loadingMore ? '加载中...' : '加载更多'}
          </Button>
        </div>
      )}
    </Container>
  );
};

export default Orders;