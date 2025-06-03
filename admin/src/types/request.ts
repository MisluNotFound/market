export interface AttributeReq {
    name: string;
    dataType: string;
    required: boolean;
    options?: string[];
    unit?: string;
}

export interface CreateCategoryReq {
    categoryName: string;
    parentID?: number;
    level?: 1 | 2 | 3;
    attributes?: AttributeReq[];
}

export interface UpdateCategoryReq {
    categoryID: number;
    categoryName: string;
}

export interface DeleteCategoryReq {
    categoryID: number;
}

export interface CreateAttributeReq extends AttributeReq {
    categoryID: number;
}

export interface DeleteAttributeReq {
    attributeID: number;
}

export interface UpdateAttributeReq extends AttributeReq {
    attributeID: number;
}

export interface CreateInterestTagReq {
    tagName: string;
    categoryID: number;
}

export interface UpdateInterestTagReq {
    tagID: number;
    tagName: string;
    categoryID: number;
}

export interface DeleteInterestTagReq {
    tagID: number;
} 