import { useState } from 'react';
import { FaSearch } from 'react-icons/fa';
import SearchHistory from './SearchHistory';
import '../styles/header.css';

const Header = ({ onSearch }) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [showHistory, setShowHistory] = useState(false);

  const handleSearch = (e) => {
    e.preventDefault();
    onSearch(searchTerm);
    setShowHistory(false);
  };

  const handleHistorySelect = (keyword) => {
    setSearchTerm(keyword);
    setShowHistory(false);
  };

  return (
    <header className="header">
      <form className="search-form" onSubmit={handleSearch}>
        <div className="search-container">
          <input
            type="text"
            placeholder="搜索商品..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            onFocus={() => setShowHistory(true)}
            className="search-input"
          />
          <button type="submit" className="search-button">
            <FaSearch />
          </button>
          <SearchHistory
            isVisible={showHistory}
            onSelect={handleHistorySelect}
            onClose={() => setShowHistory(false)}
          />
        </div>
      </form>
    </header>
  );
};

export default Header;