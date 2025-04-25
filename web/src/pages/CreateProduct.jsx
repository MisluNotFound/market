import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import ProductService from '../services/product';
import AuthService from '../services/auth';
import '../styles/create-product.css';

const CreateProduct = () => {
  const [formData, setFormData] = useState({
    originalPrice: '',
    price: '',
    describe: '',
    shipMethod: 'included',
    shipPrice: '',
    canSelfPickup: false
  });
  const [pics, setPics] = useState([]);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleFileChange = (e) => {
    const newFiles = Array.from(e.target.files);
    const totalFiles = pics.length + newFiles.length;

    if (totalFiles > 5) {
      setError('最多只能上传5张图片');
      return;
    }

    setPics(prev => [...prev, ...newFiles]);
    setError('');
  };

  const removeImage = (index) => {
    const newPics = [...pics];
    newPics.splice(index, 1);
    setPics(newPics);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    try {
      const data = new FormData();
      data.append('originalPrice', formData.originalPrice);
      data.append('price', formData.price);
      data.append('describe', formData.describe);
      data.append('shipMethod', formData.shipMethod);
      data.append('shipPrice', formData.shipPrice);
      data.append('canSelfPickup', formData.canSelfPickup);

      pics.forEach((file) => {
        data.append('pics', file);
      });

      const user = await AuthService.getCurrentUser();
      if (!user) {
        setError('请先登录');
        return;
      }
      await ProductService.createProduct(data, user.id);
      navigate('/');
    } catch (err) {
      setError(err.message || '发布商品失败');
    }
  };

  return (
    <div className="create-product-container">
      <h2>发布商品</h2>
      {error && <div className="error-message">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div className="form-row">
          <label className="form-label">商品图片</label>
          <div className="file-upload-container">
            <label className="file-upload-label">
              {pics.length > 0 ? `${pics.length}/5张图片已选择` : '点击上传图片(最多5张)'}
              <input
                type="file"
                className="file-input"
                multiple
                accept="image/*"
                onChange={handleFileChange}
                required
              />
            </label>
            <div className="image-preview-container">
              {pics.map((file, index) => (
                <div key={index} className="image-preview-item">
                  <img
                    src={URL.createObjectURL(file)}
                    alt={`预览 ${index + 1}`}
                    className="preview-image"
                  />
                  <button
                    type="button"
                    className="remove-image-btn"
                    onClick={() => removeImage(index)}
                  >
                    ×
                  </button>
                </div>
              ))}
              {pics.length < 5 && (
                <div className="image-preview-item empty-slot">
                  <span>+</span>
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="form-row">
          <label className="form-label">原价</label>
          <input
            type="number"
            className="form-control"
            name="originalPrice"
            value={formData.originalPrice}
            onChange={handleChange}
          />
        </div>

        <div className="form-row">
          <label className="form-label">现价</label>
          <input
            type="number"
            className="form-control"
            name="price"
            value={formData.price}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-row">
          <label className="form-label">商品描述</label>
          <textarea
            className="form-control description-input"
            name="describe"
            value={formData.describe}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-row">
          <label className="form-label">配送方式</label>
          <select
            className="form-control"
            name="shipMethod"
            value={formData.shipMethod}
            onChange={handleChange}
            required
          >
            <option value="included">包邮</option>
            <option value="fixed">固定运费</option>
          </select>
        </div>

        {formData.shipMethod === 'fixed' && (
          <div className="form-row">
            <label className="form-label">运费</label>
            <input
              type="number"
              className="form-control"
              name="shipPrice"
              value={formData.shipPrice}
              onChange={handleChange}
              required={formData.shipMethod === 'fixed'}
            />
          </div>
        )}

        <div className="checkbox-container">
          <input
            type="checkbox"
            name="canSelfPickup"
            checked={formData.canSelfPickup}
            onChange={handleChange}
            id="canSelfPickup"
          />
          <label className="checkbox-label" htmlFor="canSelfPickup">支持自提</label>
        </div>

        <button type="submit" className="submit-btn">发布商品</button>
      </form>
    </div>
  );
};

export default CreateProduct;