import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import AuthService from '../services/auth';
import UserProfile from '../components/UserProfile';
import '../styles/profile.css';

const Profile = () => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const userData = await AuthService.getCurrentUser();
        if (!userData) {
          navigate('/login');
          return;
        }
        setUser(userData);
      } catch (err) {
        setError('获取用户信息失败');
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, [navigate]);

  const handleBack = () => {
    navigate('/user-center');
  };

  const handleUpdate = async (updateData) => {
    try {
      await AuthService.updateBasicInfo(user.id, updateData);
      const updatedUser = await AuthService.getCurrentUser();
      setUser(updatedUser);
    } catch (err) {
      setError(err.message);
    }
  };

  const handleAvatarChange = async (file) => {
    try {
      await AuthService.uploadAvatar(user.id, file);
      const updatedUser = await AuthService.getCurrentUser();
      setUser(updatedUser);
    } catch (err) {
      setError(err.message);
    }
  };

  if (loading) return <div className="loading">加载中...</div>;
  if (!user) return <div>未登录</div>;

  return (
    <div className="profile-page">
      {/* <button className="back-btn" onClick={handleBack}>返回用户中心</button> */}
      {/* <h2>个人资料</h2> */}
      {error && <div className="error-message">{error}</div>}

      <UserProfile
        user={user}
        editable={true}
        onUpdate={handleUpdate}
        onAvatarChange={handleAvatarChange}
      />
    </div>
  );
};

export default Profile;