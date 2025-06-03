export interface Tag {
    id: number;
    tagName: string;
    categoryID: number;
}

export interface TagResponse {
    code: number;
    msg: string;
    data: Tag[];
} 