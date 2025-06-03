import axios from 'axios';
import { API_BASE_URL, API_URLS } from '../config/api';
import type {
    CreateCategoryReq,
    UpdateCategoryReq,
    DeleteCategoryReq,
    CreateAttributeReq,
    UpdateAttributeReq,
    DeleteAttributeReq,
    CreateInterestTagReq,
    UpdateInterestTagReq,
    DeleteInterestTagReq,
} from '../types/request';
import type { CategoryResponse } from '../types/category';
import type { TagResponse } from '../types/tag';

const api = axios.create({
    baseURL: API_BASE_URL,
    timeout: 10000,
});

// 分类相关 API
export const getCategories = () => {
    return api.get<CategoryResponse>(API_URLS.GET_CATEGORIES);
};

export const createCategory = (data: CreateCategoryReq) => {
    return api.post(API_URLS.CREATE_CATEGORY, data);
};

export const updateCategory = (data: UpdateCategoryReq) => {
    return api.put(API_URLS.UPDATE_CATEGORY, data);
};

export const deleteCategory = (data: DeleteCategoryReq) => {
    return api.delete(API_URLS.DELETE_CATEGORY, { data });
};

// 属性相关 API
export const createAttribute = (data: CreateAttributeReq) => {
    return api.post(API_URLS.CREATE_ATTRIBUTE, data);
};

export const updateAttribute = (data: UpdateAttributeReq) => {
    return api.put(API_URLS.UPDATE_ATTRIBUTE, data);
};

export const deleteAttribute = (data: DeleteAttributeReq) => {
    return api.delete(API_URLS.DELETE_ATTRIBUTE, { data });
};

// 标签相关 API
export const getTags = () => {
    return api.get<TagResponse>(API_URLS.GET_TAGS);
};

export const createTag = (data: CreateInterestTagReq) => {
    return api.post(API_URLS.CREATE_TAG, data);
};

export const updateTag = (data: UpdateInterestTagReq) => {
    return api.put(API_URLS.UPDATE_TAG, data);
};

export const deleteTag = (data: DeleteInterestTagReq) => {
    return api.delete(API_URLS.DELETE_TAG, { data });
}; 