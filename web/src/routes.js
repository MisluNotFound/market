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
import Chat from './components/Chat';
import AuthService from './services/auth';

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
      { path: '/login', element: <Login /> },
      { path: '/register', element: <Register /> },
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
        path: '/edit-product/:id',
        element: <EditProduct />,
        loader: protectedLoader
      },
      {
        path: '/chat',
        element: <Chat userId="current-user-id" />,
        loader: protectedLoader
      },
    ],
  },
]);

export default router;