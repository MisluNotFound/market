import { useState, useEffect, useRef } from 'react';
import { Cascader } from 'antd';

import { Button, List, message, Modal, Switch } from 'antd';
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
    const mapInstanceRef = useRef(null); // 保存地图实例

    useEffect(() => {
        loadUser();
        loadAddresses();
    }, []);

    useEffect(() => {
        // 在 Modal 打开时初始化地图
        if (mapVisible && currentLocation) {
            initMap();
        }
    }, [mapVisible, currentLocation]);

    const initMap = async () => {
        try {
            const map = await LocationService.initMapService('map', currentLocation);
            mapInstanceRef.current = map; // 保存地图实例
            initMapClickHandler(map); // 初始化地图点击事件
        } catch (error) {
            message.error('地图服务初始化失败');
        }
    };

    // 初始化地图点击事件
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
            } catch (error) {
                message.error('获取位置信息失败');
            }
        });
    };

    // 更新位置信息
    const updateLocation = (field, value) => {
        setCurrentLocation(prev => {
            const updated = { ...prev, [field]: value };

            // 当省市区变化时，自动生成地址
            if (['province', 'city', 'district'].includes(field)) {
                updated.address = `${updated.province || ''}${updated.city || ''}${updated.district || ''}`;

                // 如果地图已初始化，更新地图中心
                if (mapInstanceRef.current && updated.province && updated.city && updated.district) {
                    LocationService.geocode(updated.address).then(loc => {
                        mapInstanceRef.current.setCenter([loc.longitude, loc.latitude]);
                    }).catch(() => {
                        message.error('地址解析失败');
                    });
                }
            }
            return updated;
        });
    };

    // 加载行政区划数据
    const loadDistricts = async (selectedOptions) => {
        const targetOption = selectedOptions[selectedOptions.length - 1];
        targetOption.loading = true;

        try {
            const districts = await LocationService.getDistrict(
                selectedOptions.length === 0 ? 'province' :
                    selectedOptions.length === 1 ? 'city' : 'district',
                targetOption.value || ''
            );

            targetOption.loading = false;
            targetOption.children = districts.map(item => ({
                label: item.name,
                value: item.name,
                isLeaf: selectedOptions.length >= 2
            }));

            setCurrentLocation(prev => ({
                ...prev,
                province: selectedOptions[0]?.value || '',
                city: selectedOptions[1]?.value || '',
                district: selectedOptions[2]?.value || ''
            }));
        } catch (error) {
            message.error('获取行政区划失败');
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

    const getCurrentLocation = async () => {
        try {
            setLoading(true);
            const location = await LocationService.getCurrentPosition();
            console.log('获取到的位置信息:', location);
            setCurrentLocation(location);
            setMapVisible(true);
        } catch (error) {
            message.error('获取位置失败: ' + error.message);
        } finally {
            setLoading(false);
        }
    };

    const loadUser = async () => {
        const user = await AuthService.getCurrentUser();
        setUser(user);
    };

    const handleSaveAddress = async (address) => {
        try {
            console.log("saving", address)
            setLoading(true);
            // const receiver = document.getElementById('receiver').value;
            // const phone = document.getElementById('phone').value;
            // const detail = document.getElementById('address-detail').value;

            // if (!receiver) {
            //     message.error('请输入收件人姓名');
            //     return;
            // }
            // if (!phone || !/^1[3-9]\d{9}$/.test(phone)) {
            //     message.error('请输入正确的手机号');
            //     return;
            // }
            // if (!detail) {
            //     message.error('请输入详细地址');
            //     return;
            // }

            const addressData = {
                address: address.address,
                city: address.city,
                district: address.district,
                province: address.province,
                isDefault: false,
                detail: "六号楼",
                phone: "19862503536",
                receiver: "于泽灏",
                latitude: address.latitude,
                longitude: address.longitude,
                street: address.street,
                streetNumber: address.streetNumber
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
            loadAddresses();
        } catch (error) {
            message.error('保存地址失败');
        } finally {
            setLoading(false);
        }
    };

    const handleEditAddress = async (address) => {
        try {
            setLoading(true);
            setEditingAddressId(address.addressID);
            console.log(parseFloat(address.latitude), address.latitude);
            console.log(parseFloat(address.longitude), address.longitude);
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
        } finally {
            setLoading(false);
        }
    };

    const handleSetDefault = async (addressId, isDefault) => {
        try {
            setLoading(true);
            await AddressService.updateAddress(addressId, { isDefault });
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
                }}
                onOk={() => handleSaveAddress(currentLocation)}
                okText="保存"
                cancelText="取消"
                width={800}
            >
                {currentLocation && (
                    <div className="map-container">
                        <div className="map-form-row">
                            <div id="map" style={{ height: '400px', width: '100%' }}></div>
                            <div className="address-form">
                                <div className="form-group">
                                    <label>省市区:</label>
                                    <Cascader
                                        options={[]}
                                        loadData={loadDistricts}
                                        onChange={(value) => {
                                            updateLocation('province', value[0] || '');
                                            updateLocation('city', value[1] || '');
                                            updateLocation('district', value[2] || '');
                                        }}
                                        changeOnSelect
                                        placeholder="请选择省市区"
                                        style={{ width: '100%' }}
                                        fieldNames={{ label: 'label', value: 'value', children: 'children' }}
                                        onFocus={() => loadDistricts([])}
                                    />
                                </div>
                                <div className="form-group">
                                    <label>详细地址:</label>
                                    <input
                                        type="text"
                                        value={currentLocation.detail || ''}
                                        onChange={(e) => updateLocation('detail', e.target.value)}
                                        placeholder="例如: xx栋xx单元xx号"
                                    />
                                </div>
                                <div className="form-group">
                                    <label>收件人:</label>
                                    <input
                                        type="text"
                                        value={currentLocation.receiver || user?.username || ''}
                                        onChange={(e) => updateLocation('receiver', e.target.value)}
                                    />
                                </div>
                                <div className="form-group">
                                    <label>手机号:</label>
                                    <input
                                        type="tel"
                                        value={currentLocation.phone || user?.phone || ''}
                                        onChange={(e) => updateLocation('phone', e.target.value)}
                                    />
                                </div>
                            </div>
                        </div>
                        <div className="action-buttons">
                            <Button onClick={getCurrentLocation}>获取当前位置</Button>
                            <Button onClick={() => setMapVisible(false)}>取消</Button>
                        </div>
                    </div>
                )}
            </Modal>
        </div>
    );
};

export default AddressManagement;