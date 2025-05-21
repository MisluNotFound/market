import { useState, useEffect } from 'react';

// 递归查找分类
const findCategoryById = (categories, id) => {
  for (const category of categories) {
    if (category.ID === id) return category;
    if (category.children?.length > 0) {
      const found = findCategoryById(category.children, id);
      if (found) return found;
    }
  }
  return null;
};

// 递归渲染分类选项
const renderCategoryOptions = (categories, level = 0) => {
  return categories.map(category => (
    <>
      <option
        key={category.ID}
        value={category.ID}
        disabled={!category.IsLeaf}
        style={{
          color: level === 0 ? '#333' : level === 1 ? '#555' : '#777',
          fontWeight: level === 0 ? 'bold' : level === 1 ? '500' : 'normal',
          paddingLeft: `${level * 16}px`,
          fontSize: level === 0 ? '14px' : level === 1 ? '13px' : '12px'
        }}
      >
        {category.TypeName}
      </option>
      {category.children?.length > 0 && renderCategoryOptions(category.children, level + 1)}
    </>
  ));
};

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
    canSelfPickup: false,
    categoryId: '',
    attributes: {},
    condition: 'new',
    usedTime: ''
  });
  const [pics, setPics] = useState([]);
  const [error, setError] = useState('');
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const response = await ProductService.getCategories();
        setCategories(response.data);
      } catch (err) {
        setError('获取分类失败');
      }
    };
    fetchCategories();
  }, []);

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleAttributeChange = (attrId, value) => {
    setFormData(prev => ({
      ...prev,
      attributes: {
        ...prev.attributes,
        [attrId]: value
      }
    }));
  };

  // 渲染属性输入框
  const renderAttributeInput = (attr) => {
    switch (attr.DataType) {
      case 'STRING':
        return (
          <input
            type="text"
            className="form-control"
            onChange={(e) => handleAttributeChange(attr.ID, e.target.value)}
            required={attr.Required}
          />
        );
      case 'NUMBER':
        return (
          <input
            type="number"
            className="form-control"
            onChange={(e) => handleAttributeChange(attr.ID, e.target.value)}
            required={attr.Required}
          />
        );
      case 'ENUM':
        return (
          <select
            className="form-control"
            onChange={(e) => handleAttributeChange(attr.ID, e.target.value)}
            required={attr.Required}
          >
            <option value="">请选择</option>
            {attr.Options.map(option => (
              <option key={option} value={option}>{option}</option>
            ))}
          </select>
        );
      default:
        return null;
    }
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

    // 检查是否登录
    const userId = localStorage.getItem('userId');
    if (!userId) {
      setError('请先登录');
      navigate('/login');
      return;
    }

    // 检查图片上传
    if (pics.length === 0) {
      setError('请上传至少一张商品图片');
      return;
    }

    try {
      const data = new FormData();
      data.append('originalPrice', formData.originalPrice);
      data.append('price', formData.price);
      data.append('describe', formData.describe);
      data.append('shipMethod', formData.shipMethod);
      data.append('shipPrice', formData.shipPrice);
      data.append('canSelfPickup', formData.canSelfPickup);
      data.append('categories', formData.categoryId);

      // 转换属性格式为map[uint]string
      const attributesMap = {};
      Object.keys(formData.attributes).forEach(key => {
        attributesMap[parseInt(key)] = formData.attributes[key];
      });
      data.append('attributes', JSON.stringify(attributesMap));
      data.append('condition', formData.condition);
      data.append('usedTime', formData.usedTime);

      pics.forEach((file) => {
        data.append('pics', file);
      });

      const response = await ProductService.createProduct(data, userId);
      if (response.code === 200) {
        navigate('/', { state: { message: '商品发布成功' } });
      } else {
        setError(response.msg || '发布商品失败');
      }
    } catch (err) {
      console.error('发布商品错误:', err);
      setError(err.message || '发布商品失败');
    }
  };

  return (
    <div className="create-product-container">
      <h2>发布商品</h2>
      {error && <div className="error-message">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div className="form-row">
          <label className="form-label">商品分类</label>
          <select
            className="form-control"
            name="categoryId"
            value={formData.categoryId}
            onChange={(e) => {
              const categoryId = e.target.value;
              const category = findCategoryById(categories, parseInt(categoryId));
              setSelectedCategory(category);
              setFormData(prev => ({
                ...prev,
                categoryId,
                attributes: {}
              }));
            }}
            required
          >
            <option value="">请选择分类</option>
            {renderCategoryOptions(categories)}
          </select>
        </div>

        {selectedCategory?.attributes?.length > 0 && (
          <div className="form-row">
            <label className="form-label">商品属性</label>
            <div className="attributes-container">
              {selectedCategory.attributes.map(attr => (
                <div key={attr.ID} className="attribute-item" style={{ display: 'flex', alignItems: 'center', marginBottom: '10px' }}>
                  <label style={{ width: '120px', textAlign: 'left' }}>
                    {attr.Name}{attr.Required && <span className="required">*</span>}
                  </label>
                  <div style={{ flex: 1 }}>
                    {renderAttributeInput(attr)}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

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
          <label className="form-label">商品状态</label>
          <select
            className="form-control"
            name="condition"
            value={formData.condition}
            onChange={handleChange}
          >
            <option value="new">全新</option>
            <option value="excellent">九成新</option>
            <option value="good">八成新</option>
            <option value="used">使用过</option>
          </select>
        </div>

        {formData.condition === 'used' && (
          <div className="form-row">
            <label className="form-label">使用时间</label>
            <input
              type="text"
              className="form-control"
              name="usedTime"
              value={formData.usedTime}
              onChange={handleChange}
              required={formData.condition === 'used'}
            />
          </div>
        )}

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