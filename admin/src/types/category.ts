export type DataType = 'STRING' | 'ENUM' | 'NUMBER' | 'BOOLEAN';

export interface CategoryAttribute {
    ID: number;
    Name: string;
    DataType: DataType;
    Required: boolean;
    Options: string[] | null;
    Unit: string;
}

export interface Category {
    ID: number;
    TypeName: string;
    Level: number;
    IsLeaf: boolean;
    children: Category[];
    attributes: CategoryAttribute[] | null;
}

export interface CategoryResponse {
    code: number;
    msg: string;
    data: Category[];
} 