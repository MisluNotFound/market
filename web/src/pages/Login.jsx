import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import AuthService from '../services/auth';
import '../styles/login.css';
import { message } from 'antd';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await AuthService.login(username, password);
      message.success('登录成功');
      if (response.needSelectTags) {
        console.log(response.needSelectTags)
        navigate('/interest-tags');
      } else {
        navigate('/');
      }
    } catch (error) {
      message.error(error.message || '登录失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-container">
      <h2>登录二手市场</h2>
      {error && <div className="error-message">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>手机号</label>
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>

        <div className="form-group">
          <label>密码</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>

        <button type="submit" className="login-btn" disabled={loading}>登录</button>
      </form>

      <div className="register-prompt">
        没有账号？<Link to="/register">前去注册</Link>
      </div>
    </div>
  );
};

export default Login;