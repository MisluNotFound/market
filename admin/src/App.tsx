import { useState, useEffect } from 'react';
import { Layout, theme, Tabs, message } from 'antd';
import CategoryManager from './components/CategoryManager';
import TagManager from './components/TagManager';
import type { Category, CategoryAttribute } from './types/category';
import type { Tag } from './types/tag';
import {
  getCategories,
  createCategory,
  updateCategory,
  deleteCategory,
  createAttribute,
  updateAttribute,
  deleteAttribute,
  getTags,
  createTag,
  updateTag,
  deleteTag,
} from './services/api';

const { Header, Content } = Layout;

function App() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  // 加载初始数据
  useEffect(() => {
    fetchCategories();
    fetchTags();
  }, []);

  const fetchCategories = async () => {
    try {
      const response = await getCategories();
      setCategories(response.data.data);
    } catch (error) {
      message.error('获取分类列表失败');
    }
  };

  const fetchTags = async () => {
    try {
      const response = await getTags();
      setTags(response.data.data);
    } catch (error) {
      message.error('获取标签列表失败');
    }
  };

  // 分类管理相关函数
  const handleAddCategory = async (parentId: number | null, category: Partial<Category>) => {
    try {
      await createCategory({
        categoryName: category.TypeName!,
        parentID: parentId || undefined,
        level: category.Level as 1 | 2 | 3 | undefined,
        attributes: category.attributes?.map(attr => ({
          name: attr.Name,
          dataType: attr.DataType,
          required: attr.Required,
          options: attr.Options || undefined,
          unit: attr.Unit,
        })),
      });
      message.success('添加分类成功');
      await fetchCategories();
    } catch (error) {
      message.error('添加分类失败');
    }
  };

  const handleUpdateCategory = async (categoryId: number, category: Partial<Category>) => {
    try {
      await updateCategory({
        categoryID: categoryId,
        categoryName: category.TypeName!,
      });
      message.success('更新分类成功');
      await fetchCategories();
    } catch (error) {
      message.error('更新分类失败');
    }
  };

  const handleDeleteCategory = async (categoryId: number) => {
    try {
      await deleteCategory({ categoryID: categoryId });
      message.success('删除分类成功');
      await fetchCategories();
    } catch (error) {
      message.error('删除分类失败');
    }
  };

  const handleAddAttribute = async (categoryId: number, attribute: Partial<CategoryAttribute>) => {
    try {
      await createAttribute({
        categoryID: categoryId,
        name: attribute.Name!,
        dataType: attribute.DataType!,
        required: attribute.Required!,
        options: attribute.Options || undefined,
        unit: attribute.Unit,
      });
      message.success('添加属性成功');
      await fetchCategories();
    } catch (error) {
      message.error('添加属性失败');
    }
  };

  const handleUpdateAttribute = async (categoryId: number, attributeId: number, attribute: Partial<CategoryAttribute>) => {
    try {
      await updateAttribute({
        attributeID: attributeId,
        name: attribute.Name!,
        dataType: attribute.DataType!,
        required: attribute.Required!,
        options: attribute.Options || undefined,
        unit: attribute.Unit,
      });
      message.success('更新属性成功');
      await fetchCategories();
    } catch (error) {
      message.error('更新属性失败');
    }
  };

  const handleDeleteAttribute = async (categoryId: number, attributeId: number) => {
    try {
      await deleteAttribute({ attributeID: attributeId });
      message.success('删除属性成功');
      await fetchCategories();
    } catch (error) {
      message.error('删除属性失败');
    }
  };

  // 标签管理相关函数
  const handleAddTag = async (tag: Partial<Tag>) => {
    try {
      await createTag({
        tagName: tag.tagName!,
        categoryID: tag.categoryID!,
      });
      message.success('添加标签成功');
      fetchTags();
    } catch (error) {
      message.error('添加标签失败');
    }
  };

  const handleUpdateTag = async (tagId: number, tag: Partial<Tag>) => {
    try {
      await updateTag({
        tagID: tagId,
        tagName: tag.tagName!,
        categoryID: tag.categoryID!,
      });
      message.success('更新标签成功');
      fetchTags();
    } catch (error) {
      message.error('更新标签失败');
    }
  };

  const handleDeleteTag = async (tagId: number) => {
    try {
      await deleteTag({ tagID: tagId });
      message.success('删除标签成功');
      fetchTags();
    } catch (error) {
      message.error('删除标签失败');
    }
  };

  const items = [
    {
      key: 'category',
      label: '分类管理',
      children: (
        <CategoryManager
          categories={categories}
          onAddCategory={handleAddCategory}
          onUpdateCategory={handleUpdateCategory}
          onDeleteCategory={handleDeleteCategory}
          onAddAttribute={handleAddAttribute}
          onUpdateAttribute={handleUpdateAttribute}
          onDeleteAttribute={handleDeleteAttribute}
        />
      ),
    },
    {
      key: 'tag',
      label: '兴趣标签',
      children: (
        <TagManager
          tags={tags}
          categories={categories}
          onAddTag={handleAddTag}
          onUpdateTag={handleUpdateTag}
          onDeleteTag={handleDeleteTag}
        />
      ),
    },
  ];

  return (
    <Layout>
      <Header style={{ display: 'flex', alignItems: 'center' }}>
        <h1 style={{ color: '#fff', margin: 0 }}>商品管理系统</h1>
      </Header>
      <Content style={{ padding: '24px' }}>
        <div
          style={{
            background: colorBgContainer,
            padding: 24,
            borderRadius: borderRadiusLG,
            minHeight: 'calc(100vh - 152px)',
          }}
        >
          <Tabs items={items} />
        </div>
      </Content>
    </Layout>
  );
}

export default App;
