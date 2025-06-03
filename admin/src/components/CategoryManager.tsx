import React, { useState } from 'react';
import { Tree, Button, Modal, Form, Input, Select, Space, Table, message, Tag } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import type { DataNode } from 'antd/es/tree';
import type { Category, CategoryAttribute } from '../types/category';

interface CategoryManagerProps {
    categories: Category[];
    onAddCategory: (parentId: number | null, category: Partial<Category>) => Promise<void>;
    onUpdateCategory: (categoryId: number, category: Partial<Category>) => Promise<void>;
    onDeleteCategory: (categoryId: number) => Promise<void>;
    onAddAttribute: (categoryId: number, attribute: Partial<CategoryAttribute>) => Promise<void>;
    onUpdateAttribute: (categoryId: number, attributeId: number, attribute: Partial<CategoryAttribute>) => Promise<void>;
    onDeleteAttribute: (categoryId: number, attributeId: number) => Promise<void>;
}

const CategoryManager: React.FC<CategoryManagerProps> = ({
    categories,
    onAddCategory,
    onUpdateCategory,
    onDeleteCategory,
    onAddAttribute,
    onUpdateAttribute,
    onDeleteAttribute,
}) => {
    const [selectedCategory, setSelectedCategory] = useState<Category | null>(null);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [modalType, setModalType] = useState<'category' | 'attribute'>('category');
    const [modalAction, setModalAction] = useState<'add' | 'edit'>('add');
    const [form] = Form.useForm();

    const convertToTreeData = (categories: Category[]): DataNode[] => {
        return categories.map(cat => ({
            key: cat.ID.toString(),
            title: cat.TypeName,
            children: cat.children ? convertToTreeData(cat.children) : undefined,
        }));
    };

    const handleSelect = (selectedKeys: React.Key[]) => {
        const findCategory = (categories: Category[], id: number): Category | null => {
            for (const cat of categories) {
                if (cat.ID === id) return cat;
                if (cat.children) {
                    const found = findCategory(cat.children, id);
                    if (found) return found;
                }
            }
            return null;
        };

        if (selectedKeys.length > 0) {
            const category = findCategory(categories, Number(selectedKeys[0]));
            setSelectedCategory(category);
        }
    };

    const showModal = (type: 'category' | 'attribute', action: 'add' | 'edit', record?: CategoryAttribute) => {
        setModalType(type);
        setModalAction(action);
        setIsModalVisible(true);
        form.resetFields();

        if (action === 'edit') {
            if (type === 'category' && selectedCategory) {
                form.setFieldsValue({
                    typeName: selectedCategory.TypeName,
                });
            } else if (type === 'attribute' && record) {
                form.setFieldsValue({
                    attributeId: record.ID,
                    name: record.Name,
                    dataType: record.DataType,
                    required: record.Required,
                    options: record.Options || undefined,
                    unit: record.Unit,
                });
            }
        }
    };

    const handleOk = async () => {
        try {
            const values = await form.validateFields();
            if (modalType === 'category') {
                if (modalAction === 'add') {
                    const level = selectedCategory ? selectedCategory.Level + 1 : 1;
                    await onAddCategory(selectedCategory?.ID || null, {
                        TypeName: values.typeName,
                        Level: level,
                        IsLeaf: level === 3,
                    });
                } else {
                    if (selectedCategory) {
                        await onUpdateCategory(selectedCategory.ID, {
                            TypeName: values.typeName,
                        });
                    }
                }
            } else {
                if (modalAction === 'add' && selectedCategory) {
                    await onAddAttribute(selectedCategory.ID, {
                        Name: values.name,
                        DataType: values.dataType,
                        Required: values.required,
                        Options: values.dataType === 'ENUM' ? values.options : null,
                        Unit: values.unit || '',
                    });
                } else if (selectedCategory && values.attributeId) {
                    await onUpdateAttribute(selectedCategory.ID, values.attributeId, {
                        Name: values.name,
                        DataType: values.dataType,
                        Required: values.required,
                        Options: values.dataType === 'ENUM' ? values.options : null,
                        Unit: values.unit || '',
                    });
                }
            }
            setIsModalVisible(false);
            setSelectedCategory(null);
            message.success('操作成功');
        } catch (error) {
            message.error('操作失败');
        }
    };

    const handleDelete = async (record: CategoryAttribute) => {
        try {
            if (selectedCategory) {
                await onDeleteAttribute(selectedCategory.ID, record.ID);
                setSelectedCategory(null);
                message.success('删除属性成功');
            }
        } catch (error) {
            message.error('删除属性失败');
        }
    };

    const handleDeleteCategory = async () => {
        try {
            if (selectedCategory) {
                await onDeleteCategory(selectedCategory.ID);
                setSelectedCategory(null);
                message.success('删除分类成功');
            }
        } catch (error) {
            message.error('删除分类失败');
        }
    };

    const attributeColumns = [
        {
            title: '属性名称',
            dataIndex: 'Name',
            key: 'name',
        },
        {
            title: '数据类型',
            dataIndex: 'DataType',
            key: 'dataType',
        },
        {
            title: '是否必填',
            dataIndex: 'Required',
            key: 'required',
            render: (required: boolean) => required ? '是' : '否',
        },
        {
            title: '选项',
            dataIndex: 'Options',
            key: 'options',
            render: (options: string[] | null) => {
                if (!options) return '-';
                return (
                    <Space wrap>
                        {options.map((opt, index) => (
                            <Tag key={index}>{opt}</Tag>
                        ))}
                    </Space>
                );
            },
        },
        {
            title: '单位',
            dataIndex: 'Unit',
            key: 'unit',
            render: (unit: string) => unit || '-',
        },
        {
            title: '操作',
            key: 'action',
            render: (_: any, record: CategoryAttribute) => (
                <Space>
                    <Button
                        icon={<EditOutlined />}
                        onClick={() => showModal('attribute', 'edit', record)}
                    >
                        编辑
                    </Button>
                    <Button
                        danger
                        icon={<DeleteOutlined />}
                        onClick={() => handleDelete(record)}
                    >
                        删除
                    </Button>
                </Space>
            ),
        },
    ];

    return (
        <div style={{ display: 'flex', gap: '24px' }}>
            <div style={{ width: '300px' }}>
                <div style={{ marginBottom: '16px' }}>
                    <Space wrap>
                        <Button
                            type="primary"
                            icon={<PlusOutlined />}
                            onClick={() => showModal('category', 'add')}
                        >
                            {selectedCategory ? '新增子分类' : '新增一级分类'}
                        </Button>
                        <Button
                            icon={<EditOutlined />}
                            disabled={!selectedCategory}
                            onClick={() => showModal('category', 'edit')}
                        >
                            编辑分类
                        </Button>
                        <Button
                            danger
                            icon={<DeleteOutlined />}
                            disabled={!selectedCategory}
                            onClick={handleDeleteCategory}
                        >
                            删除分类
                        </Button>
                    </Space>
                </div>
                <Tree
                    treeData={convertToTreeData(categories)}
                    onSelect={handleSelect}
                    selectedKeys={selectedCategory ? [selectedCategory.ID.toString()] : []}
                />
            </div>
            <div style={{ flex: 1 }}>
                <div style={{ marginBottom: '16px' }}>
                    <Button
                        type="primary"
                        icon={<PlusOutlined />}
                        disabled={!selectedCategory || selectedCategory.Level !== 3}
                        onClick={() => showModal('attribute', 'add')}
                    >
                        新增属性
                    </Button>
                </div>
                {selectedCategory?.Level === 3 && (
                    <Table
                        columns={attributeColumns}
                        dataSource={selectedCategory?.attributes || []}
                        rowKey="ID"
                    />
                )}
                {selectedCategory && selectedCategory.Level !== 3 && (
                    <div style={{ textAlign: 'center', color: '#999', padding: '20px' }}>
                        请选择三级分类以管理属性
                    </div>
                )}
                {!selectedCategory && (
                    <div style={{ textAlign: 'center', color: '#999', padding: '20px' }}>
                        请选择一个分类
                    </div>
                )}
            </div>
            <Modal
                title={`${modalAction === 'add' ? '新增' : '编辑'}${modalType === 'category' ? '分类' : '属性'}`}
                open={isModalVisible}
                onOk={handleOk}
                onCancel={() => setIsModalVisible(false)}
            >
                <Form form={form} layout="vertical">
                    {modalType === 'category' ? (
                        <Form.Item
                            name="typeName"
                            label="分类名称"
                            rules={[{ required: true, message: '请输入分类名称' }]}
                        >
                            <Input />
                        </Form.Item>
                    ) : (
                        <>
                            <Form.Item name="attributeId" hidden>
                                <Input />
                            </Form.Item>
                            <Form.Item
                                name="name"
                                label="属性名称"
                                rules={[{ required: true, message: '请输入属性名称' }]}
                            >
                                <Input />
                            </Form.Item>
                            <Form.Item
                                name="dataType"
                                label="数据类型"
                                rules={[{ required: true, message: '请选择数据类型' }]}
                            >
                                <Select>
                                    <Select.Option value="STRING">文本</Select.Option>
                                    <Select.Option value="NUMBER">数字</Select.Option>
                                    <Select.Option value="BOOLEAN">布尔值</Select.Option>
                                    <Select.Option value="ENUM">枚举</Select.Option>
                                    <Select.Option value="DATE">日期</Select.Option>
                                </Select>
                            </Form.Item>
                            <Form.Item
                                name="required"
                                label="是否必填"
                                initialValue={false}
                            >
                                <Select>
                                    <Select.Option value={true}>是</Select.Option>
                                    <Select.Option value={false}>否</Select.Option>
                                </Select>
                            </Form.Item>
                            <Form.Item
                                noStyle
                                shouldUpdate={(prevValues, currentValues) =>
                                    prevValues?.dataType !== currentValues?.dataType
                                }
                            >
                                {({ getFieldValue }) =>
                                    getFieldValue('dataType') === 'ENUM' && (
                                        <Form.Item
                                            name="options"
                                            label="选项列表"
                                            rules={[{ required: true, message: '请输入选项列表' }]}
                                        >
                                            <Select
                                                mode="tags"
                                                placeholder="请输入选项，按回车键确认"
                                                style={{ width: '100%' }}
                                            />
                                        </Form.Item>
                                    )
                                }
                            </Form.Item>
                            <Form.Item name="unit" label="单位">
                                <Input />
                            </Form.Item>
                        </>
                    )}
                </Form>
            </Modal>
        </div>
    );
};

export default CategoryManager; 