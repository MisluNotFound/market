import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Select } from 'antd';
import ProductService from '../services/product';
import AuthService from '../services/auth';
import AddressService from '../services/address';
import '../styles/create-product.css';

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

const EditProduct = () => {
  const { id: productId } = useParams();
  const [formData, setFormData] = useState({
    originalPrice: '',
    price: '',
    describe: '',
    shipMethod: 'included',
    shipPrice: '',
    canSelfPickup: false,
    categoryId: '',
    attributes: {},
    addressId: ''
  });
  const [originalAttributes, setOriginalAttributes] = useState({});
  const [pics, setPics] = useState([]);
  const [deletedPics, setDeletedPics] = useState([]);
  const [error, setError] = useState('');
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState(null);
  const [existingPics, setExistingPics] = useState([]);
  const [addresses, setAddresses] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchData = async () => {
      try {
        if (!productId) {
          throw new Error('缺少商品ID参数');
        }

        const userId = localStorage.getItem('userId');
        if (!userId) {
          navigate('/login');
          return;
        }

        const categoriesRes = await ProductService.getCategories();
        const productRes = await ProductService.getProductDetail(userId, productId);
        const product = productRes.data;

        // 获取地址列表
        const addressResponse = await AddressService.getAddresses();
        if (addressResponse?.data?.code === 200) {
          const addressList = addressResponse.data.data.addresses || [];
          setAddresses(addressList);

          // 根据商品的location信息匹配地址
          if (product.product.location) {
            const matchedAddress = addressList.find(addr =>
              addr.addressID == product.product.location
            );

            if (matchedAddress) {
              setFormData(prev => ({
                ...prev,
                addressId: matchedAddress.addressID
              }));
            }
          }
        }

        setCategories(categoriesRes.data);
        setFormData(prev => ({
          ...prev,
          originalPrice: product.product.originalPrice || '',
          price: product.product.price,
          describe: product.product.describe,
          shipMethod: product.product.shipMethod,
          shipPrice: product.product.shipPrice,
          canSelfPickup: product.product.canSelfPickup,
          categoryId: product.categories[0],
          attributes: product.attributes || {}
        }));
        setOriginalAttributes(product.attributes || {});
        const category = findCategoryById(categoriesRes.data, product.categories[0]);
        setSelectedCategory(category);

        setExistingPics(product.product.avatar.split(',').filter(url => url.trim()));
      } catch (err) {
        console.log(err);
        setError('获取商品信息失败');
      }
    };

    fetchData();
  }, [productId, navigate]);

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

  const renderAttributeInput = (attr) => {
    const value = formData.attributes[attr.ID] ?? originalAttributes[attr.ID] ?? '';
    switch (attr.DataType) {
      case 'STRING':
        return (
          <input
            type="text"
            className="form-control"
            value={value}
            onChange={(e) => handleAttributeChange(attr.ID, e.target.value)}
            required={attr.Required}
          />
        );
      case 'NUMBER':
        return (
          <input
            type="number"
            className="form-control"
            value={value}
            onChange={(e) => handleAttributeChange(attr.ID, e.target.value)}
            required={attr.Required}
          />
        );
      case 'ENUM':
        return (
          <select
            className="form-control"
            value={value}
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
    const totalFiles = pics.length + existingPics.length - deletedPics.length + newFiles.length;

    if (totalFiles > 5) {
      setError('最多只能上传5张图片');
      return;
    }

    setPics(prev => [...prev, ...newFiles]);
    setError('');
  };

  const removeImage = (index, isExisting) => {
    if (isExisting) {
      setDeletedPics(prev => [...prev, existingPics[index]]);
      setExistingPics(prev => prev.filter((_, i) => i !== index));
    } else {
      setPics(prev => prev.filter((_, i) => i !== index));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    const userId = localStorage.getItem('userId');
    if (!userId) {
      setError('请先登录');
      navigate('/login');
      return;
    }

    if (!formData.addressId) {
      setError('请选择收货地址');
      return;
    }

    try {
      const data = new FormData();
      data.append('originalPrice', formData.originalPrice);
      data.append('price', formData.price);
      data.append('describe', formData.describe);
      data.append('shipMethod', formData.shipMethod || 'included');
      data.append('shipPrice', formData.shipMethod === 'fixed' ? formData.shipPrice : '0');
      data.append('canSelfPickup', formData.canSelfPickup);
      data.append('categories', parseInt(formData.categoryId));
      data.append('addressId', formData.addressId);

      const attributesMap = {};
      const hasAttributeChanges = Object.keys(formData.attributes).length > 0;
      const finalAttributes = hasAttributeChanges
        ? formData.attributes
        : originalAttributes;

      Object.keys(finalAttributes).forEach(key => {
        attributesMap[parseInt(key)] = finalAttributes[key];
      });
      data.append('attributes', JSON.stringify(attributesMap));

      // 添加删除的图片
      deletedPics.forEach(pic => {
        data.append('deletedPics', pic);
      });

      // 添加新上传的图片
      pics.forEach((file) => {
        data.append('addedPics', file);
      });

      const response = await ProductService.updateProduct(data, userId, productId);
      if (response.code === 200) {
        navigate(`/product/${productId}`, {
          state: {
            product: {
              user: {
                id: userId
              }
            }
          }
        });
      } else {
        setError(response.msg || '修改商品失败');
      }
    } catch (err) {
      console.error('修改商品错误:', err);
      setError(err.message || '修改商品失败');
    }
  };

  return (
    <div className="create-product-container">
      <h2>修改商品</h2>
      {error && <div className="error-message">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div className="form-row">
          <label className="form-label">收货地址</label>
          <Select
            className="form-control"
            value={formData.addressId}
            onChange={(value) => setFormData(prev => ({ ...prev, addressId: value }))}
            placeholder="请选择收货地址"
            required
          >
            {addresses.map(address => (
              <Select.Option key={address.addressID} value={address.addressID}>
                {`${address.receiver} ${address.phone} - ${address.province}${address.city}${address.district} ${address.street}${address.streetNumber} ${address.detail}`}
              </Select.Option>
            ))}
          </Select>
        </div>

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
              {existingPics.length + pics.length - deletedPics.length > 0
                ? `${existingPics.length + pics.length - deletedPics.length}/5张图片已选择`
                : '点击上传图片(最多5张)'}
              <input
                type="file"
                className="file-input"
                multiple
                accept="image/*"
                onChange={handleFileChange}
              />
            </label>
            <div className="image-preview-container">
              {existingPics.map((url, index) => (
                <div key={`existing-${index}`} className="image-preview-item">
                  <img
                    src={url}
                    alt={`预览 ${index + 1}`}
                    className="preview-image"
                  />
                  <button
                    type="button"
                    className="remove-image-btn"
                    onClick={() => removeImage(index, true)}
                  >
                    ×
                  </button>
                </div>
              ))}
              {pics.map((file, index) => (
                <div key={`new-${index}`} className="image-preview-item">
                  <img
                    src={URL.createObjectURL(file)}
                    alt={`预览 ${existingPics.length + index + 1}`}
                    className="preview-image"
                  />
                  <button
                    type="button"
                    className="remove-image-btn"
                    onClick={() => removeImage(index, false)}
                  >
                    ×
                  </button>
                </div>
              ))}
              {existingPics.length + pics.length - deletedPics.length < 5 && (
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

        <button type="submit" className="submit-btn">保存修改</button>
      </form>
    </div>
  );
};

export default EditProduct;