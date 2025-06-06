import { useState, useEffect, useRef } from 'react';
import { Cascader, Form, Input, Button, List, message, Modal, Switch, Radio } from 'antd';
import { EnvironmentOutlined, SearchOutlined } from '@ant-design/icons';
import LocationService from '../services/location';
import AddressService from '../services/address';
import AuthService from '../services/auth';
import '../styles/address-management.css';

const AddressManagement = () => {
    const [addresses, setAddresses] = useState([]);
    const [loading, setLoading] = useState(true);
    const [user, setUser] = useState(null);
    const [mapVisible, setMapVisible] = useState(false);
    const [currentLocation, setCurrentLocation] = useState(null);
    const [editingAddressId, setEditingAddressId] = useState(null);
    const [addressForm] = Form.useForm();
    const mapInstanceRef = useRef(null);
    const [locationMode, setLocationMode] = useState('auto'); // 'auto' 或 'manual'
    const [searchKeyword, setSearchKeyword] = useState('');

    useEffect(() => {
        loadUser();
        loadAddresses();
    }, []);

    useEffect(() => {
        if (mapVisible && currentLocation) {
            initMap();
        }
    }, [mapVisible, currentLocation]);

    const initMap = async () => {
        try {
            const map = await LocationService.initMapService('map', currentLocation);
            mapInstanceRef.current = map;
            initMapClickHandler(map);
        } catch (error) {
            message.error('地图服务初始化失败');
        }
    };

    const initMapClickHandler = (map) => {
        map.on('click', async (e) => {
            try {
                const location = await LocationService.reverseGeocode(e.lnglat.lat, e.lnglat.lng);
                setCurrentLocation(prev => ({
                    ...prev,
                    ...location,
                    latitude: e.lnglat.lat,
                    longitude: e.lnglat.lng
                }));
                addressForm.setFieldsValue({
                    province: location.province,
                    city: location.city,
                    district: location.district,
                    street: location.street,
                    streetNumber: location.streetNumber
                });
            } catch (error) {
                message.error('获取位置信息失败');
            }
        });
    };

    const handleSearch = async () => {
        if (!searchKeyword.trim()) {
            message.warning('请输入搜索关键词');
            return;
        }

        try {
            const location = await LocationService.geocode(searchKeyword);
            if (mapInstanceRef.current) {
                mapInstanceRef.current.setCenter([location.longitude, location.latitude]);

                // 从formattedAddress中解析地址信息
                const addressInfo = await LocationService.reverseGeocode(location.latitude, location.longitude);

                setCurrentLocation(prev => ({
                    ...prev,
                    ...location,
                    ...addressInfo
                }));

                // 确保所有地址字段都被正确填充
                addressForm.setFieldsValue({
                    province: addressInfo.province || '',
                    city: addressInfo.city || '',
                    district: addressInfo.district || '',
                    street: addressInfo.street || '',
                    streetNumber: addressInfo.streetNumber || '',
                    detail: location.formattedAddress || ''
                });

                // 如果地图已初始化，更新地图中心
                if (location.latitude && location.longitude) {
                    mapInstanceRef.current.setCenter([location.longitude, location.latitude]);
                }
            }
        } catch (error) {
            message.error('地址搜索失败');
        }
    };

    const getCurrentLocation = async () => {
        try {
            setLoading(true);
            const location = await LocationService.getCurrentPosition();
            setCurrentLocation(location);
            setMapVisible(true);
            setLocationMode('auto');
            addressForm.setFieldsValue({
                province: location.province,
                city: location.city,
                district: location.district,
                street: location.street,
                streetNumber: location.streetNumber,
            });
        } catch (error) {
            message.error('获取位置失败: ' + error.message);
        } finally {
            setLoading(false);
        }
    };

    const handleSaveAddress = async () => {
        try {
            const values = await addressForm.validateFields();
            setLoading(true);

            // 构建完整地址字符串
            const fullAddress = `${values.province}${values.city}${values.district}${values.street}${values.streetNumber}`;

            const addressData = {
                ...values,
                address: fullAddress, // 添加完整地址
                latitude: currentLocation.latitude,
                longitude: currentLocation.longitude,
                isDefault: false
            };

            if (editingAddressId) {
                await AddressService.updateAddress(editingAddressId, addressData);
            } else {
                await AddressService.createAddress(addressData);
            }

            message.success(editingAddressId ? '地址更新成功' : '地址保存成功');
            setMapVisible(false);
            setEditingAddressId(null);
            setCurrentLocation(null);
            addressForm.resetFields();
            loadAddresses();
        } catch (error) {
            if (error.errorFields) {
                message.error('请填写完整的地址信息');
            } else {
                message.error('保存地址失败');
            }
        } finally {
            setLoading(false);
        }
    };

    const handleEditAddress = async (address) => {
        try {
            setLoading(true);
            setEditingAddressId(address.addressID);
            setCurrentLocation({
                address: address.address,
                city: address.city,
                district: address.district,
                province: address.province,
                street: address.street,
                streetNumber: address.streetNumber,
                latitude: parseFloat(address.latitude),
                longitude: parseFloat(address.longitude)
            });
            setMapVisible(true);
            setLocationMode('manual');
            addressForm.setFieldsValue({
                receiver: address.receiver,
                phone: address.phone,
                province: address.province,
                city: address.city,
                district: address.district,
                street: address.street,
                streetNumber: address.streetNumber,
                detail: address.detail
            });
        } finally {
            setLoading(false);
        }
    };

    const handleSetDefault = async (addressId, isDefault) => {
        try {
            setLoading(true);
            await AddressService.setDefaultAddress(addressId, isDefault);
            message.success('默认地址设置成功');
            loadAddresses();
        } catch (error) {
            message.error('设置默认地址失败');
        } finally {
            setLoading(false);
        }
    };

    const handleDeleteAddress = async (addressId) => {
        try {
            setLoading(true);
            await AddressService.deleteAddress(addressId);
            message.success('地址删除成功');
            loadAddresses();
        } catch (error) {
            message.error('删除地址失败');
        } finally {
            setLoading(false);
        }
    };

    const loadUser = async () => {
        try {
            const user = await AuthService.getCurrentUser();
            setUser(user);
        } catch (error) {
            message.error('获取用户信息失败');
        }
    };

    const loadAddresses = async () => {
        try {
            setLoading(true);
            const response = await AddressService.getAddresses();
            if (response.data.code === 200) {
                setAddresses(response.data.data.addresses || []);
            } else {
                message.error(response.data.msg || '获取地址列表失败');
            }
        } catch (error) {
            message.error('获取地址列表失败');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="address-management">
            <div className="header">
                <h2>收货地址管理</h2>
                <Button type="primary" onClick={getCurrentLocation} loading={loading}>
                    新增收货地址
                </Button>
            </div>

            <List
                dataSource={addresses}
                renderItem={(item) => (
                    <List.Item
                        actions={[
                            <Switch
                                key="default"
                                checkedChildren="默认"
                                unCheckedChildren="设置"
                                checked={item.isDefault}
                                onChange={(checked) => handleSetDefault(item.addressID, checked)}
                            />,
                            <a key="edit" onClick={() => handleEditAddress(item)}>
                                编辑
                            </a>,
                            <a key="delete" onClick={() => handleDeleteAddress(item.addressID)}>
                                删除
                            </a>,
                        ]}
                    >
                        <List.Item.Meta
                            title={`${item.receiver} ${item.phone}`}
                            description={`${item.province}${item.city}${item.district} ${item.street}${item.streetNumber} ${item.detail}`}
                        />
                    </List.Item>
                )}
            />

            <Modal
                title="选择收货地址"
                open={mapVisible}
                onCancel={() => {
                    setMapVisible(false);
                    setCurrentLocation(null);
                    setEditingAddressId(null);
                    addressForm.resetFields();
                }}
                onOk={handleSaveAddress}
                okText="保存"
                cancelText="取消"
                width={800}
            >
                {currentLocation && (
                    <div className="map-container">
                        <div className="location-mode">
                            <Radio.Group value={locationMode} onChange={(e) => setLocationMode(e.target.value)}>
                                <Radio.Button value="auto">
                                    <EnvironmentOutlined /> 自动定位
                                </Radio.Button>
                                <Radio.Button value="manual">
                                    <SearchOutlined /> 手动输入
                                </Radio.Button>
                            </Radio.Group>
                        </div>

                        <div className="map-form-row">
                            <div id="map" style={{ height: '400px', width: '100%' }}></div>
                            <div className="address-form">
                                <Form
                                    form={addressForm}
                                    layout="vertical"
                                    initialValues={{
                                        receiver: user?.username || '',
                                        phone: user?.phone || ''
                                    }}
                                >
                                    <Form.Item
                                        name="receiver"
                                        label="收货人"
                                        rules={[{ required: true, message: '请输入收货人姓名' }]}
                                    >
                                        <Input placeholder="请输入收货人姓名" />
                                    </Form.Item>

                                    <Form.Item
                                        name="phone"
                                        label="手机号码"
                                        rules={[
                                            { required: true, message: '请输入手机号码' },
                                            { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码' }
                                        ]}
                                    >
                                        <Input placeholder="请输入手机号码" />
                                    </Form.Item>

                                    {locationMode === 'manual' && (
                                        <div className="search-box">
                                            <Input
                                                placeholder="搜索地址"
                                                value={searchKeyword}
                                                onChange={(e) => setSearchKeyword(e.target.value)}
                                                onPressEnter={handleSearch}
                                                suffix={
                                                    <SearchOutlined onClick={handleSearch} style={{ cursor: 'pointer' }} />
                                                }
                                            />
                                        </div>
                                    )}

                                    <Form.Item
                                        name="province"
                                        label="省份"
                                        rules={[{ required: true, message: '请选择省份' }]}
                                    >
                                        <Input disabled />
                                    </Form.Item>

                                    <Form.Item
                                        name="city"
                                        label="城市"
                                        rules={[{ required: true, message: '请选择城市' }]}
                                    >
                                        <Input disabled />
                                    </Form.Item>

                                    <Form.Item
                                        name="district"
                                        label="区县"
                                        rules={[{ required: true, message: '请选择区县' }]}
                                    >
                                        <Input disabled />
                                    </Form.Item>

                                    <Form.Item
                                        name="street"
                                        label="街道"
                                        rules={[{ required: true, message: '请输入街道' }]}
                                    >
                                        <Input placeholder="请输入街道" />
                                    </Form.Item>

                                    <Form.Item
                                        name="streetNumber"
                                        label="门牌号"
                                        rules={[{ required: true, message: '请输入门牌号' }]}
                                    >
                                        <Input placeholder="请输入门牌号" />
                                    </Form.Item>

                                    <Form.Item
                                        name="detail"
                                        label="详细地址"
                                        rules={[{ required: true, message: '请输入详细地址' }]}
                                    >
                                        <Input.TextArea
                                            placeholder="请输入详细地址，如：xx栋xx单元xx号"
                                            rows={3}
                                        />
                                    </Form.Item>
                                </Form>
                            </div>
                        </div>
                    </div>
                )}
            </Modal>
        </div>
    );
};

export default AddressManagement;