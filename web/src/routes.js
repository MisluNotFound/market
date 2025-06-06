import { createBrowserRouter } from 'react-router-dom';
import App from './App';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import UserCenter from './pages/UserCenter';
import Profile from './pages/Profile';
import Orders from './pages/Orders';
import OrderDetail from './pages/OrderDetail';
import CreateProduct from './pages/CreateProduct';
import ProductDetail from './pages/ProductDetail';
import EditProduct from './pages/EditProduct';
import MyProducts from './pages/MyProducts';
import Favorites from './pages/Favorites';
import AddressManagement from './pages/AddressManagement';
import Chat from './components/Chat';
import SearchResults from './pages/SearchResults';
import InterestTags from './pages/InterestTags';
import AuthService from './services/auth';
import UserProfile from './pages/UserProfile';

const protectedLoader = async () => {
  const user = await AuthService.getCurrentUser();
  if (!user) {
    return { redirect: '/login' };
  }
  return null;
};

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      { path: '/', element: <Home /> },
      { path: '/search', element: <SearchResults /> },
      { path: '/login', element: <Login /> },
      { path: '/register', element: <Register /> },
      { path: '/interest-tags', element: <InterestTags /> },
      {
        path: '/user-center',
        element: <UserCenter />,
        loader: protectedLoader
      },
      {
        path: '/profile',
        element: <Profile />,
        loader: protectedLoader
      },
      {
        path: '/orders/:type',
        element: <Orders />,
        loader: protectedLoader
      },
      {
        path: '/order/:id',
        element: <OrderDetail />,
        loader: protectedLoader
      },
      {
        path: '/create-product',
        element: <CreateProduct />,
        loader: protectedLoader
      },
      {
        path: '/product/:id',
        element: <ProductDetail />,
        loader: protectedLoader
      },
      {
        path: '/my-products',
        element: <MyProducts />,
        loader: protectedLoader
      },
      {
        path: '/my-favorites',
        element: <Favorites />,
        loader: protectedLoader
      },
      {
        path: '/edit-product/:id',
        element: <EditProduct />,
        loader: protectedLoader
      },
      {
        path: '/addresses',
        element: <AddressManagement />,
        loader: protectedLoader
      },
      {
        path: '/chat',
        element: <Chat userId="current-user-id" />,
        loader: protectedLoader
      },
      {
        path: '/user/:userId',
        element: <UserProfile />
      },
    ],
  },
]);

export default router;