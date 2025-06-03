import axios from 'axios';

const AMAP_KEY = '391c64d0c9a28e082671ddd2cad3240b';

export default {
    // 初始化地图服务
    initMapService(containerId, center) {
        return new Promise((resolve) => {
            if (window.AMap) {
                const map = new AMap.Map(containerId, {
                    zoom: 15,
                    center: [center.longitude, center.latitude],
                    viewMode: '2D'
                });
                // 添加标记点
                new AMap.Marker({
                    position: [center.longitude, center.latitude],
                    map: map
                });
                resolve(map);
                return;
            }

            const script = document.createElement('script');
            script.src = `https://webapi.amap.com/maps?v=2.0&key=${AMAP_KEY}`;
            script.onload = () => {
                const map = new AMap.Map(containerId, {
                    zoom: 15,
                    center: [center.longitude, center.latitude],
                    viewMode: '2D'
                });
                // 添加标记点
                new AMap.Marker({
                    position: [center.longitude, center.latitude],
                    map: map
                });
                resolve(map);
            };
            document.head.appendChild(script);
        });
    },

    // 获取当前位置
    getCurrentPosition() {
        return new Promise((resolve, reject) => {
            if (!navigator.geolocation) {
                reject(new Error('浏览器不支持地理位置功能'));
                return;
            }

            navigator.geolocation.getCurrentPosition(
                (position) => {
                    this.reverseGeocode(position.coords.latitude, position.coords.longitude)
                        .then(resolve)
                        .catch(reject);
                },
                (error) => reject(error)
            );
        });
    },

    // 逆地理编码
    reverseGeocode(lat, lng) {
        return axios.get(`https://restapi.amap.com/v3/geocode/regeo?key=${AMAP_KEY}&location=${lng},${lat}`)
            .then(response => {
                const result = response.data.regeocode;
                return {
                    address: result.formatted_address,
                    province: result.addressComponent.province,
                    city: result.addressComponent.city,
                    district: result.addressComponent.district,
                    street: result.addressComponent.streetNumber.street,
                    streetNumber: result.addressComponent.streetNumber.number,
                    latitude: lat,
                    longitude: lng
                };
            });
    },

    // 地理编码
    geocode(address) {
        return axios.get(`https://restapi.amap.com/v3/geocode/geo?key=${AMAP_KEY}&address=${encodeURIComponent(address)}`)
            .then(response => {
                const location = response.data.geocodes[0]?.location;
                if (!location) {
                    throw new Error('地址解析失败');
                }
                const [lng, lat] = location.split(',');
                return {
                    latitude: parseFloat(lat),
                    longitude: parseFloat(lng),
                    formattedAddress: response.data.geocodes[0].formatted_address
                };
            });
    },

    // 获取行政区划数据
    async getDistrict(level = 'province', keywords = '') {
        try {
            const response = await axios.get(
                `https://restapi.amap.com/v3/config/district?key=${AMAP_KEY}&keywords=${encodeURIComponent(keywords)}&subdistrict=1&extensions=base`
            );
            const districts = response.data.districts[0]?.districts || [];
            return districts.map(item => ({
                name: item.name,
                level: item.level
            }));
        } catch (error) {
            console.error('获取行政区划失败:', error);
            return [];
        }
    }
};