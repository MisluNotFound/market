import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import { useParams, useNavigate } from 'react-router-dom';
import { Descriptions, Button, message, Tag } from 'antd';
import AuthService from '../services/auth';
import OrderService from '../services/order';

const Container = styled.div`
  max-width: 800px;
  margin: 50px auto;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
`;

const Title = styled.h2`
  margin-bottom: 24px;
`;

const OrderDetail = () => {
  const [loading, setLoading] = useState(false);
  const [order, setOrder] = useState(null);
  const [currentUser, setCurrentUser] = useState(null);
  const { id: orderId } = useParams();
  const navigate = useNavigate();

  useEffect(() => {
    if (!AuthService.isAuthenticated()) {
      message.warning('请先登录');
      navigate('/login');
      return;
    }

    const fetchCurrentUser = async () => {
      const user = await AuthService.getCurrentUser();
      setCurrentUser(user);
    };

    fetchOrder();
    fetchCurrentUser();
  }, [orderId, navigate]);

  const fetchOrder = async () => {
    setLoading(true);
    try {
      const user = await AuthService.getCurrentUser();
      const data = await OrderService.getOrderDetail(orderId, user.id);
      const orderData = data.data.order
      const productData = data.data.product
      orderData.product = productData
      console.log(orderData)
      setOrder(orderData);
    } catch (error) {
      message.error(error.message);
    } finally {
      setLoading(false);
    }
  };

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

  const handleCancel = async () => {
    setLoading(true);
    try {
      await OrderService.cancelOrder(orderId);
      message.success('订单已取消');
      fetchOrder();
    } catch (error) {
      message.error(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handlePay = async () => {
    setLoading(true);
    try {
      await OrderService.payOrder(orderId);
      message.success('订单支付成功');
      fetchOrder();
    } catch (error) {
      message.error(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleShip = async (refund) => {
    setLoading(true);
    try {
      await OrderService.confirmShipped(currentUser.id, orderId, refund);
      message.success('订单已发货');
      fetchOrder();
    } catch (error) {
      message.error(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleComplete = async (isRefund = false) => {
    setLoading(true);
    try {
      await OrderService.confirmSigned(currentUser.id, orderId, isRefund);
      message.success(isRefund ? '退货确认成功' : '订单已完成');
      fetchOrder();
    } catch (error) {
      message.error(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleRefund = async () => {
    setLoading(true);
    try {
      await OrderService.refundOrder(currentUser.id, orderId);
      message.success('退款处理成功');
      fetchOrder();
    } catch (error) {
      message.error(error.message);
    } finally {
      setLoading(false);
    }
  };

  if (!order) return null;

  const isBuyer = currentUser && currentUser.id === order.userID;
  const statusInfo = statusMap[order.status] || { text: '未知状态', color: 'default' };

  return (
    <Container>
      <Title>订单详情</Title>
      <div style={{ display: 'flex', gap: '24px', marginBottom: '24px' }}>
        <img
          src={order.product?.avatar?.split(',')[0] || '/placeholder-product.png'}
          alt={order.describe}
          style={{ width: '200px', height: '200px', objectFit: 'cover', borderRadius: '8px' }}
        />
        <div style={{ flex: 1 }}>
          <p style={{ color: '#666', marginBottom: '16px' }}>{order.product?.describe || '暂无商品描述'}</p>
          <div style={{ fontSize: '24px', color: '#ff4d4f', marginBottom: '16px' }}>
            ¥{(order.totalAmount || 0).toFixed(2)}
          </div>
        </div>
      </div>
      <Descriptions bordered column={2}>
        <Descriptions.Item label="订单编号">{order.id}</Descriptions.Item>
        <Descriptions.Item label="状态">
          <Tag color={statusInfo.color}>{statusInfo.text}</Tag>
        </Descriptions.Item>
        <Descriptions.Item label="运费">
          {order.shipAmount === 0 ? (
            <Tag color="green">包邮</Tag>
          ) : (
            `¥${order.shipAmount}`
          )}
        </Descriptions.Item>
        <Descriptions.Item label="下单时间">{order.payTime ? order.payTime.replace('T', ' ').substring(0, 16) : ''}</Descriptions.Item>
        <Descriptions.Item label="完成时间">{order.finishTime ? order.finishTime.replace('T', ' ').substring(0, 16) : ''}</Descriptions.Item>
      </Descriptions>

      <div style={{ marginTop: 24, textAlign: 'center' }}>
        {/* 待付款状态(1): 买家可支付或取消 */}
        {order.status === 1 && isBuyer && (
          <>
            <Button
              type="primary"
              onClick={handlePay}
              loading={loading}
              style={{ marginRight: 16 }}
            >
              立即支付
            </Button>
            <Button
              danger
              onClick={() => handleCancel(currentUser.id, orderId)}
              loading={loading}
            >
              取消订单
            </Button>
          </>
        )}

        {/* 已付款状态(2): 卖家可发货 */}
        {order.status === 2 && !isBuyer && (
          <Button
            type="primary"
            onClick={() => handleShip(false)}
            loading={loading}
          >
            发货
          </Button>
        )}

        {/* 已发货状态(3): 买家可确认收货或退款 */}
        {order.status === 3 && isBuyer && (
          <>
            <Button
              type="primary"
              onClick={() => handleComplete(false)}
              loading={loading}
              style={{ marginRight: 16 }}
            >
              确认收货
            </Button>
            <Button
              danger
              onClick={handleRefund}
              loading={loading}
            >
              申请退款
            </Button>
          </>
        )}

        {/* 已完成状态(4): 买家可申请退款 */}
        {order.status === 4 && isBuyer && (
          <Button
            danger
            onClick={handleRefund}
            loading={loading}
          >
            申请退款
          </Button>
        )}

        {/* 已退款状态(5): 无操作按钮 */}

        {/* 已退货状态(6): 卖家可确认收货 */}
        {order.status === 6 && !isBuyer && (
          <Button
            type="primary"
            onClick={() => handleComplete(true)}
            loading={loading}
          >
            确认退货
          </Button>
        )}
      </div>
    </Container>
  );
};

export default OrderDetail;