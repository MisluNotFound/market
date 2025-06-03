import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, Button, message, Tag, Space } from 'antd';
import styled from 'styled-components';
import ProductService from '../services/product';
import UserService from '../services/user';

const Container = styled.div`
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
`;

const StyledCard = styled(Card)`
  width: 90%;
  max-width: 600px;
  text-align: center;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
`;

const Title = styled.h1`
  font-size: 24px;
  color: #333;
  margin-bottom: 8px;
`;

const Subtitle = styled.p`
  font-size: 14px;
  color: #666;
  margin-bottom: 24px;
`;

const TagsContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
  margin-bottom: 32px;
  padding: 0 24px;
`;

const StyledTag = styled(Tag)`
  padding: 8px 16px;
  font-size: 14px;
  border-radius: 16px;
  cursor: pointer;
  user-select: none;
  margin: 0;
  
  &:hover {
    opacity: 0.8;
  }
`;

const ButtonGroup = styled(Space)`
  margin-top: 24px;
`;

const InterestTags = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [tags, setTags] = useState([]);
    const [selectedTags, setSelectedTags] = useState([]);

    useEffect(() => {
        fetchTags();
    }, []);

    const fetchTags = async () => {
        try {
            const response = await ProductService.getTags();
            setTags(response.data || []);
        } catch (error) {
            message.error('获取标签失败');
        }
    };

    const handleTagClick = (tagId) => {
        setSelectedTags(prev => {
            if (prev.includes(tagId)) {
                return prev.filter(id => id !== tagId);
            }
            return [...prev, tagId];
        });
    };

    const handleSubmit = async () => {
        if (selectedTags.length === 0) {
            message.warning('请至少选择一个兴趣标签');
            return;
        }

        setLoading(true);
        try {
            const userId = localStorage.getItem('userId');
            await UserService.selectInterestTags(userId, selectedTags);
            message.success('设置成功');
            navigate('/');
        } catch (error) {
            message.error('设置失败');
        } finally {
            setLoading(false);
        }
    };

    const handleSkip = () => {
        navigate('/');
    };

    return (
        <Container>
            <StyledCard>
                <Title>选择你感兴趣的内容</Title>
                <Subtitle>帮助我们了解你的偏好，为你推荐更好的内容</Subtitle>

                <TagsContainer>
                    {tags.map(tag => (
                        <StyledTag
                            key={tag.id}
                            color={selectedTags.includes(tag.id) ? '#1890ff' : 'default'}
                            onClick={() => handleTagClick(tag.id)}
                        >
                            {tag.tagName}
                        </StyledTag>
                    ))}
                </TagsContainer>

                <ButtonGroup>
                    <Button onClick={handleSkip}>跳过</Button>
                    <Button
                        type="primary"
                        onClick={handleSubmit}
                        loading={loading}
                        disabled={selectedTags.length === 0}
                    >
                        完成
                    </Button>
                </ButtonGroup>
            </StyledCard>
        </Container>
    );
};

export default InterestTags; 