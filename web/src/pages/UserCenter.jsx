import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  FaUser,
  FaShoppingCart,
  FaStore,
  FaCommentDots,
  FaChevronRight
} from 'react-icons/fa';
import AuthService from '../services/auth';
import OrderService from '../services/order';
import '../styles/user-center.css';

const UserCenter = () => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [orderStatus, setOrderStatus] = useState({
    bought: 0,
    sold: 0,
    beEvaluated: 0
  });
  const navigate = useNavigate();

  const checkAuth = async () => {
    const userData = await AuthService.getCurrentUser();
    if (!userData) {
      alert('请先登录');
      navigate('/login');
      return false;
    }
    return true;
  };

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const userData = await AuthService.getCurrentUser();
        if (!userData) {
          navigate('/login');
          return;
        }
        setUser(userData);
      } finally {
        setLoading(false);
      }
    };

    const fetchData = async () => {
      try {
        const userData = await AuthService.getCurrentUser();
        if (!userData) {
          navigate('/login');
          return;
        }
        setUser(userData);

        // 获取订单状态
        const status = await OrderService.getAllOrderStatus(userData.id);
        console.log(status)
        setOrderStatus({
          bought: status.bought || 0,
          sold: status.sold || 0,
          beEvaluated: status.beEvaluated || 0
        });
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [navigate]);

  if (loading) return <div className="loading">加载中...</div>;

  return (
    <div className="user-center">
      {/* 用户信息头部 */}
      <div className="user-header" onClick={() => navigate('/profile')}>
        <div className="avatar-container">
          <img
            src={user.avatar || '/default-avatar.png'}
            alt="用户头像"
            className="avatar"
          />
        </div>
        <div className="user-info">
          <h3>{user.username}</h3>
          <p>查看个人资料 <FaChevronRight /></p>
        </div>
      </div>

      {/* 订单状态分类 */}
      <div className="order-sections">
        <h2 className="section-title">我的订单</h2>

        <div className="order-cards">
          <div className="order-card" onClick={async () => {
            navigate('/orders/sold');
          }}>
            <div className="icon-container sold">
              <FaStore />
              {orderStatus.sold > 0 && <span className="badge">{orderStatus.sold}</span>}
            </div>
            <span>我卖出的</span>
          </div>

          <div className="order-card" onClick={async () => {
            navigate('/orders/bought');
          }}>
            <div className="icon-container bought">
              <FaShoppingCart />
              {orderStatus.bought > 0 && <span className="badge">{orderStatus.bought}</span>}
            </div>
            <span>我买到的</span>
          </div>

          <div className="order-card" onClick={async () => {
            navigate('/orders/reviews');
          }}>
            <div className="icon-container review">
              <FaCommentDots />
              {orderStatus.beEvaluated > 0 && <span className="badge">{orderStatus.beEvaluated}</span>}
            </div>
            <span>待评价</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default UserCenter;