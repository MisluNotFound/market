import React, { useState } from 'react';
import { Table, Button, Modal, Form, Input, Space, message, Select } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import type { Tag } from '../types/tag';
import type { Category } from '../types/category';

interface TagManagerProps {
    tags: Tag[];
    categories: Category[];
    onAddTag: (tag: Partial<Tag>) => Promise<void>;
    onUpdateTag: (tagId: number, tag: Partial<Tag>) => Promise<void>;
    onDeleteTag: (tagId: number) => Promise<void>;
}

const TagManager: React.FC<TagManagerProps> = ({
    tags,
    categories,
    onAddTag,
    onUpdateTag,
    onDeleteTag,
}) => {
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [editingTag, setEditingTag] = useState<Tag | null>(null);
    const [form] = Form.useForm();

    const showModal = (tag?: Tag) => {
        setEditingTag(tag || null);
        setIsModalVisible(true);
        form.resetFields();
        if (tag) {
            form.setFieldsValue({
                tagName: tag.tagName,
                categoryID: tag.categoryID,
            });
        }
    };

    const handleOk = async () => {
        try {
            const values = await form.validateFields();
            if (editingTag) {
                await onUpdateTag(editingTag.id, values);
            } else {
                await onAddTag(values);
            }
            setIsModalVisible(false);
            message.success('操作成功');
        } catch (error) {
            message.error('操作失败');
        }
    };

    const getCategoryName = (categoryId: number) => {
        const findCategory = (categories: Category[]): string => {
            for (const cat of categories) {
                if (cat.ID === categoryId) return cat.TypeName;
                if (cat.children) {
                    const name = findCategory(cat.children);
                    if (name) return name;
                }
            }
            return '';
        };
        return findCategory(categories);
    };

    const columns = [
        {
            title: '标签名称',
            dataIndex: 'tagName',
            key: 'tagName',
        },
        {
            title: '所属分类',
            dataIndex: 'categoryID',
            key: 'categoryID',
            render: (categoryId: number) => getCategoryName(categoryId),
        },
        {
            title: '操作',
            key: 'action',
            render: (_: any, record: Tag) => (
                <Space>
                    <Button
                        icon={<EditOutlined />}
                        onClick={() => showModal(record)}
                    >
                        编辑
                    </Button>
                    <Button
                        danger
                        icon={<DeleteOutlined />}
                        onClick={() => onDeleteTag(record.id)}
                    >
                        删除
                    </Button>
                </Space>
            ),
        },
    ];

    const getCategoryOptions = (categories: Category[]): { value: number; label: string }[] => {
        let options: { value: number; label: string }[] = [];
        categories.forEach(cat => {
            options.push({ value: cat.ID, label: cat.TypeName });
            if (cat.children) {
                options = options.concat(getCategoryOptions(cat.children));
            }
        });
        return options;
    };

    return (
        <div>
            <div style={{ marginBottom: '16px' }}>
                <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={() => showModal()}
                >
                    新增标签
                </Button>
            </div>
            <Table
                columns={columns}
                dataSource={tags}
                rowKey="id"
            />
            <Modal
                title={editingTag ? '编辑标签' : '新增标签'}
                open={isModalVisible}
                onOk={handleOk}
                onCancel={() => setIsModalVisible(false)}
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        name="tagName"
                        label="标签名称"
                        rules={[{ required: true, message: '请输入标签名称' }]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        name="categoryID"
                        label="所属分类"
                        rules={[{ required: true, message: '请选择所属分类' }]}
                    >
                        <Select
                            options={getCategoryOptions(categories)}
                            placeholder="请选择所属分类"
                        />
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};

export default TagManager; 