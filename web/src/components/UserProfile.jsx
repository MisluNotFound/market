import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import AuthService from '../services/auth';

const UserProfile = ({ user, editable = false, onUpdate, onAvatarChange }) => {
  const navigate = useNavigate();

  const handleLogout = () => {
    AuthService.logout();
    navigate('/login');
  };
  const [editMode, setEditMode] = useState(false);
  const [formData, setFormData] = useState({
    username: user.username || '',
    gender: user.gender || 'male'
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onUpdate(formData);
    setEditMode(false);
  };

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      onAvatarChange(file);
    }
  };

  return (
    <div className="user-profile">
      <div className="avatar-section">
        <img 
          src={user.avatar || '/default-avatar.png'} 
          alt="用户头像"
          className="avatar"
        />
        {editable && (
          <div className="avatar-upload">
            <label htmlFor="avatar-upload">更换头像</label>
            <input 
              id="avatar-upload"
              type="file" 
              accept="image/*"
              onChange={handleFileChange}
            />
          </div>
        )}
      </div>

      <div className="user-info">
        {editMode ? (
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label>用户名</label>
              <input
                type="text"
                name="username"
                value={formData.username}
                onChange={handleChange}
              />
            </div>
            <div className="form-group">
              <label>性别</label>
              <select 
                name="gender" 
                value={formData.gender}
                onChange={handleChange}
              >
                <option value="male">男</option>
                <option value="female">女</option>
              </select>
            </div>
            <button type="submit">保存</button>
            <button type="button" onClick={() => setEditMode(false)}>取消</button>
          </form>
        ) : (
          <>
            <div className="info-item">
              <span className="label">用户名:</span>
              <span className="value">{user.username}</span>
            </div>
            <div className="info-item">
              <span className="label">手机号:</span>
              <span className="value">{user.phone}</span>
            </div>
            <div className="info-item">
              <span className="label">性别:</span>
              <span className="value">{user.gender === 'male' ? '男' : '女'}</span>
            </div>
            {editable && (
              <div className="profile-actions">
                <button onClick={() => setEditMode(true)}>编辑信息</button>
                <button
                  onClick={handleLogout}
                  className="logout-btn"
                >
                  退出登录
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