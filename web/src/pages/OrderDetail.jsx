import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import { useParams, useNavigate } from 'react-router-dom';
import { Descriptions, Button, message, Tag, Input, Modal, Spin } from 'antd';
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

const PaymentStatusModal = styled(Modal)`
  .ant-modal-body {
    text-align: center;
    padding: 24px;
  }
`;

const SpinContainer = styled.div`
  margin: 20px 0;
  text-align: center;
`;

const OrderDetail = () => {
  const [loading, setLoading] = useState(false);
  const [order, setOrder] = useState(null);
  const [currentUser, setCurrentUser] = useState(null);
  const { id: orderId } = useParams();
  const navigate = useNavigate();
  const [paymentPolling, setPaymentPolling] = useState(false);
  const [pollingCount, setPollingCount] = useState(0);
  const MAX_POLLING_COUNT = 60; // 最大轮询次数，即60次
  const POLLING_INTERVAL = 3000; // 轮询间隔，3秒

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

  // 轮询订单状态
  const pollOrderStatus = async (orderId) => {
    setPaymentPolling(true);
    setPollingCount(0);

    const poll = async () => {
      if (pollingCount >= MAX_POLLING_COUNT) {
        message.error('支付超时，请稍后在订单中心查看支付状态');
        setPaymentPolling(false);
        return;
      }

      try {
        const response = await OrderService.getOrderStatus(orderId);
        if (response.code === 200) {
          if (response.data.status === 2) {
            message.success('支付成功！');
            setPaymentPolling(false);
            fetchOrder(); // 刷新订单信息
            return;
          }
        }

        setPollingCount(prev => prev + 1);
        setTimeout(poll, POLLING_INTERVAL);
      } catch (error) {
        message.error('获取支付状态失败');
        setPaymentPolling(false);
      }
    };

    poll();
  };

  const handlePay = async () => {
    setLoading(true);
    try {
      const response = await OrderService.payOrder(orderId);
      if (response.code === 200 && response.data?.payURL) {
        // 在新标签页打开支付宝支付页面
        window.open(response.data.payURL, '_blank');
        // 开始轮询支付状态
        pollOrderStatus(orderId);
      } else {
        throw new Error(response.msg || '获取支付链接失败');
      }
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

  // 评论相关状态
  const [comments, setComments] = useState([]);
  const [commentContent, setCommentContent] = useState('');
  const [replyingTo, setReplyingTo] = useState(null);
  const [replyContent, setReplyContent] = useState('');
  const [rating, setRating] = useState(5); // 默认5星好评

  // 渲染评论项
  const renderComment = (comment) => {
    return (
      <div key={comment.id} style={{
        padding: '16px',
        borderBottom: '1px solid #f0f0f0',
        backgroundColor: 'transparent',
        margin: '8px 0',
        borderRadius: '4px',
        border: '1px solid #eee'
      }}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '12px' }}>
          <img
            src={comment.avatar || '/default-avatar.png'}
            alt={comment.username}
            style={{ width: '32px', height: '32px', borderRadius: '50%', marginRight: '8px' }}
          />
          <span style={{ fontWeight: 'bold' }}>{comment.username}</span>
          {comment.replyTo && (
            <span style={{ color: '#666', marginLeft: '8px' }}>
              回复 @{comment.replyTo}
            </span>
          )}
          <span style={{ color: '#999', marginLeft: 'auto' }}>
            {new Date(comment.createdAt).toLocaleString()}
          </span>
        </div>
        <div style={{ marginLeft: '40px', textAlign: 'left' }}>
          <p style={{ margin: '0 0 8px 0' }}>{comment.comment}</p>
          {comment.pics && comment.pics.split(',').map(pic => (
            <img
              key={pic}
              src={pic}
              style={{ width: '100px', height: '100px', objectFit: 'cover', marginRight: '8px' }}
            />
          ))}
          {currentUser && (
            <Button
              type="text"
              size="small"
              onClick={() => setReplyingTo(comment)}
              style={{ padding: '4px 0' }}
            >
              回复
            </Button>
          )}
        </div>

        {replyingTo?.id === comment.id && (
          <div style={{
            marginTop: '12px',
            marginLeft: '40px',
            padding: '12px',
            backgroundColor: '#f9f9f9',
            borderRadius: '4px'
          }}>
            <Input.TextArea
              rows={2}
              value={replyContent}
              onChange={(e) => setReplyContent(e.target.value)}
              placeholder={`回复 ${comment.username}`}
            />
            <div style={{ marginTop: '8px', textAlign: 'right' }}>
              <Button size="small" onClick={() => setReplyingTo(null)}>
                取消
              </Button>
              <Button
                type="primary"
                size="small"
                onClick={() => handleReply(comment.id)}
                style={{ marginLeft: '8px' }}
              >
                发送
              </Button>
            </div>
          </div>
        )}

        {comment.replies?.map(reply => renderComment(reply))}
      </div>
    );
  };

  // 加载评论
  useEffect(() => {
    if (orderId) {
      const loadComments = async () => {
        try {
          const res = await OrderService.getOrderComments(orderId);
          setComments(res.data.comments || []);
        } catch (error) {
          message.error(error.message);
        }
      };
      loadComments();
    }
  }, [orderId]);

  const handleCommentSubmit = async () => {
    if (!commentContent.trim()) {
      message.warning('请输入评论内容');
      return;
    }
    try {
      const isGood = rating >= 3; // 3星以上算好评
      await OrderService.createOrderComment(orderId, commentContent, isGood);
      const res = await OrderService.getOrderComments(orderId);
      setComments(res.data.comments || []);
      message.success('评论提交成功');
      setCommentContent('');
      setRating(5); // 重置评分
    } catch (error) {
      message.error(error.message);
    }
  };

  const handleReply = async (commentId) => {
    if (!replyContent.trim()) {
      message.warning('请输入回复内容');
      return;
    }
    try {
      await OrderService.replyOrderComment(orderId, commentId, replyContent);
      const res = await OrderService.getOrderComments(orderId);
      setComments(res.data.comments || []);
      message.success('回复提交成功');
      setReplyingTo(null);
      setReplyContent('');
    } catch (error) {
      message.error(error.message);
    }
  };

  // 渲染评分星星
  const renderRatingStars = () => {
    return (
      <div style={{ marginBottom: '12px' }}>
        <span style={{ marginRight: '8px' }}>评分:</span>
        {[1, 2, 3, 4, 5].map((star) => (
          <span
            key={star}
            style={{
              cursor: 'pointer',
              color: star <= rating ? '#ffc107' : '#e4e5e9',
              fontSize: '24px'
            }}
            onClick={() => setRating(star)}
          >
            ★
          </span>
        ))}
        <span style={{ marginLeft: '8px', color: '#666' }}>
          {rating >= 3.5 ? '好评' : '差评'}
        </span>
      </div>
    );
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

      <div style={{ marginTop: 24, marginBottom: 24, textAlign: 'center' }}>
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

      {/* 评论区域 */}
      {order.status === 4 && (
        <div style={{ marginTop: '32px', borderTop: '1px solid #f0f0f0', paddingTop: '24px' }}>
          <h3 style={{ marginBottom: '16px', textAlign: 'left' }}>商品评价</h3>

          {/* 评论列表 */}
          <div style={{ marginBottom: '24px' }}>
            {comments.length > 0 ? (
              comments.map(comment => renderComment(comment))
            ) : (
              <div style={{ color: '#999', padding: '24px 0', textAlign: 'left' }}>
                暂无评价
              </div>
            )}
          </div>

          {/* 评论输入框 - 只对未评价的买家显示 */}
          {currentUser &&
            currentUser.id === order.userID &&
            !order.isEvaluated && (
              <div style={{ textAlign: 'left' }}>
                {renderRatingStars()}
                <Input.TextArea
                  rows={4}
                  value={commentContent}
                  onChange={(e) => setCommentContent(e.target.value)}
                  placeholder="写下您的评价..."
                  style={{ marginBottom: '12px' }}
                />
                <div style={{ textAlign: 'right' }}>
                  <Button
                    type="primary"
                    onClick={handleCommentSubmit}
                    disabled={!commentContent.trim()}
                  >
                    提交评价
                  </Button>
                </div>
              </div>
            )}
        </div>
      )}

      {/* 支付状态轮询提示 */}
      <PaymentStatusModal
        title="支付确认中"
        open={paymentPolling}
        footer={null}
        closable={false}
        maskClosable={false}
      >
        <SpinContainer>
          <Spin size="large" />
          <p style={{ marginTop: 16 }}>正在确认支付状态，请稍候...</p>
          <p style={{ color: '#999' }}>剩余等待时间：{Math.max(0, MAX_POLLING_COUNT - pollingCount)} 秒</p>
        </SpinContainer>
      </PaymentStatusModal>
    </Container>
  );
};

export default OrderDetail;