import { Outlet } from 'react-router-dom';
import Navbar from './components/Navbar';
import './App.css';

function App() {
  return (
    <div className="app">
      <div className="content-container">
        <Outlet />
      </div>
      <Navbar />
    </div>
  );
}

export default App;