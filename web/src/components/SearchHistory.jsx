import { useEffect, useState } from 'react';
import '../styles/searchHistory.css';

const SearchHistory = ({ isVisible, onSelect, onClose }) => {
    const [history, setHistory] = useState([]);

    useEffect(() => {
        const fetchHistory = async () => {
            try {
                const userId = localStorage.getItem("userId")
                const response = await fetch(`http://localhost:3200/api/search/${userId}/history`);
                const result = await response.json();
                if (result.code === 200) {
                    setHistory(result.data.history);
                }
            } catch (error) {
                console.error('获取搜索历史失败:', error);
            }
        };

        if (isVisible) {
            fetchHistory();
        }
    }, [isVisible]);

    if (!isVisible) return null;

    return (
        <div className="search-history-container">
            <div className="search-history-header">
                <h3>搜索历史</h3>
                <button onClick={onClose} className="close-button">×</button>
            </div>
            <ul className="search-history-list">
                {history.map((item, index) => (
                    <li
                        key={index}
                        onClick={() => onSelect(item.Keyword)}
                        className="search-history-item"
                    >
                        {item.Keyword}
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default SearchHistory; 