import { Link } from 'react-router-dom';
import { FaHome, FaPlusCircle, FaUser } from 'react-icons/fa';
import '../styles/navbar.css';

const Navbar = () => {
  return (
    <nav className="navbar">
      <Link to="/" className="nav-item">
        <FaHome className="nav-icon" />
        <span>首页</span>
      </Link>
      <Link to="/create-product" className="nav-item add-button">
        <FaPlusCircle className="nav-icon" />
      </Link>
      <Link to="/user-center" className="nav-item">
        <FaUser className="nav-icon" />
        <span>我的</span>
      </Link>
    </nav>
  );
};

export default Navbar;