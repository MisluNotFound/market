import { useEffect, useState, useRef, useCallback } from 'react';
import { useLocation } from 'react-router-dom';
import AuthService from '../services/auth';
import ProductService from '../services/product';
import '../styles/chat.css';

export default function Chat({ userId }) {
  const location = useLocation();
  const [conversations, setConversations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [activeConversation, setActiveConversation] = useState(null);
  const [conversationMessages, setConversationMessages] = useState(new Map());
  const messages = activeConversation ? conversationMessages.get(activeConversation.id) || [] : [];
  const [input, setInput] = useState('');
  const [currentProduct, setCurrentProduct] = useState(null);
  const imService = useRef(AuthService.getIMService());

  const sortMessages = useCallback((messages) => {
    const sorted = [...messages].sort((a, b) => {
      const timeA = new Date(a.time).getTime();
      const timeB = new Date(b.time).getTime();
      console.log(`排序比较: ${a.id} (${timeA}) vs ${b.id} (${timeB})`);
      return timeA - timeB; // 升序排列，最早的消息在前
    });
    console.log('排序后的消息:', sorted);
    return sorted;
  }, []);

  useEffect(() => {
    if (!imService.current) {
      console.error('IMService not initialized');
      return;
    }
    const initChat = async () => {
      try {
        const currentUserId = localStorage.getItem('userId');
        if (!currentUserId) throw new Error('用户未登录');
        const convList = await imService.current.getConversationList(currentUserId);
        setConversations(convList);

        // 处理从商品详情跳转过来的自动选择
        if (location.state?.fromUserID && location.state?.toUserID) {
          const targetConv = convList.find(c =>
            (c.fromUserID === location.state.fromUserID && c.toUserID === location.state.toUserID) ||
            (c.fromUserID === location.state.toUserID && c.toUserID === location.state.fromUserID)
          );
          if (targetConv) {
            const conv = {
              ...targetConv,
              id: targetConv.fromUserID === userId ? targetConv.toUserID : targetConv.fromUserID
            };
            console.log('准备调用selectConversation:', conv);
            await selectConversation(conv);
            console.log('selectConversation调用完成');
          }
        }
      } catch (error) {
        console.error('初始化聊天失败:', error);
      } finally {
        setLoading(false);
      }
    };

    const handleMessage = (message) => {
      console.log('处理新消息:', message);
      const conversationId = message.from === userId ? message.to : message.from;

      // 自动激活或创建对话
      if (!activeConversation || activeConversation.id !== conversationId) {
        const targetConv = conversations.find(c =>
          c.id === conversationId ||
          c.fromUserID === conversationId ||
          c.toUserID === conversationId
        );

        if (targetConv) {
          console.log('自动激活对话:', targetConv);
          setActiveConversation({
            id: conversationId,
            name: targetConv.username,
            avatar: targetConv.avatar
          });
        } else {
          console.log('创建新对话记录:', conversationId);
          setActiveConversation({
            id: conversationId,
            name: '新对话',
            avatar: '/placeholder-user.png'
          });
        }
      }

      setConversationMessages(prev => {
        const newMap = new Map(prev);
        const currentMessages = newMap.get(conversationId) || [];
        const newMessage = {
          id: message.id || message.tempID || Date.now().toString(),
          from: message.from_user_id || message.from,
          content: message.content,
          time: message.timestamp || new Date().toISOString()
        };
        console.log('格式化后的消息:', newMessage);

        const updatedMessages = sortMessages([...currentMessages, newMessage]);
        newMap.set(conversationId, updatedMessages);
        console.log('更新后的消息列表:', updatedMessages);
        return newMap;
      });
    };

    const removeHandler = imService.current?.addMessageHandler(handleMessage);
    initChat();

    return () => {
      removeHandler();
    };
  }, [userId, location.state]);

  const handleSend = () => {
    if (input.trim() && activeConversation) {
      const newMessage = {
        id: Date.now().toString(),
        from: userId,
        content: input,
        time: new Date().toISOString()
      };

      setConversationMessages(prev => {
        const newMap = new Map(prev);
        const currentMessages = newMap.get(activeConversation.id) || [];
        newMap.set(activeConversation.id, sortMessages([...currentMessages, newMessage]));
        return newMap;
      });

      imService.current?.sendMessage(input, activeConversation.id);
      setInput('');
    }
  };

  const selectConversation = async (conversation) => {
    console.log(conversation.username, conversation.avatar)
    setActiveConversation({
      id: conversation.id,
      name: conversation.username,
      avatar: conversation.avatar,
      productId: conversation.currentProductID
    });
    console.log(conversation)
    // 获取关联商品信息
    if (conversation.currentProductID) {
      try {
        console.log('获取商品详情参数:', {
          userId: conversation.fromUserID,
          productId: conversation.currentProductID
        });
        const product = await ProductService.getProductDetail(
          conversation.fromUserID,
          conversation.currentProductID
        );
        console.log('商品详情响应:', product.data.product);
        setCurrentProduct(product.data.product);
      } catch (error) {
        console.error('获取商品信息失败:', error);
        setCurrentProduct(null);
      }
    } else {
      setCurrentProduct(null);
    }
    try {
      // 获取对话消息时，fromUserID应该是对方用户ID
      console.log(conversation)
      const otherUserId = conversation.fromUserID
      const msgResponse = await imService.current.getConversationMessages(
        conversation.fromUserID,  // fromUserID应该是对方用户
        conversation.toUserID
      );
      const sortedMessages = sortMessages(
        (msgResponse?.data?.Messages || []).map(msg => ({
          id: msg.id,
          from: msg.from_user_id,
          content: msg.content,
          time: msg.timestamp
        }))
      );
      setConversationMessages(prev => {
        const newMap = new Map(prev);
        newMap.set(conversation.id, sortedMessages);
        return newMap;
      });
    } catch (error) {
      console.error('获取消息失败:', error);
    }
  };

  if (loading) return <div className="loading">加载中...</div>;

  return (
    <div className="chat-app">
      <div className="conversation-list">
        {conversations.map(conversation => (
          <div
            key={conversation.id}
            className={`conversation-item ${activeConversation?.id === conversation.id ? 'active' : ''}`}
            onClick={() => selectConversation(conversation)}
          >
            <div
              className="conversation-avatar"
              style={{ backgroundImage: `url(${conversation.avatar || '/placeholder-user.png'})` }}
            ></div>
            <div className="conversation-info">
              <div className="conversation-name">{conversation.username}</div>
              <div className="conversation-preview">{conversation.lastMessageContent}</div>
              <div className="conversation-time">
                {new Date(conversation.lastMessageTime).toLocaleString()}
              </div>
            </div>
            {conversation.unread > 0 && (
              <div className="unread-badge">{conversation.unread}</div>
            )}
          </div>
        ))}
      </div>

      <div className="chat-area">
        {activeConversation ? (
          <>
            {currentProduct && (
              <div className="chat-product-info">
                <img
                  src={currentProduct.avatar.split(',')[0].trim()}
                  alt="商品图片"
                  className="product-image"
                />
                <div className="product-details">
                  <div className="product-name">{currentProduct.title}</div>
                  <div className="product-info-line">
                    <span className="product-price">¥{currentProduct.price}</span>
                    <span className="product-description">{currentProduct.describe}</span>
                  </div>
                </div>
              </div>
            )}
            <div className="chat-header">
              <div className="chat-user">
                <img
                  src={activeConversation.avatar || '/placeholder-user.png'}
                  className="chat-user-avatar"
                  alt="用户头像"
                />
                <div className="chat-title">{activeConversation.name}</div>
              </div>
              {currentProduct && (
                <div className="chat-product-brief">
                  正在交易: {currentProduct.title}
                </div>
              )}
            </div>
            <div className="chat-messages">
              {messages.map((msg, index) => (
                <div key={index} className={`message-container ${msg.from === userId ? 'sent' : 'received'}`}>
                  {msg.from !== userId && (
                    <div
                      className="message-avatar"
                      style={{ backgroundImage: `url(${activeConversation.avatar})` }}
                    />
                  )}
                  <div className={`message ${msg.from === userId ? 'sent' : 'received'}`}>
                    <div className="message-content">{msg.content}</div>
                    <div className="message-time">
                      {new Date(msg.time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                    </div>
                  </div>
                </div>
              ))}
            </div>
            <div className="chat-input">
              <textarea
                value={input}
                onChange={(e) => setInput(e.target.value)}
                placeholder="输入消息..."
                onKeyPress={(e) => e.key === 'Enter' && !e.shiftKey && handleSend()}
              />
              <button onClick={handleSend}>发送</button>
            </div>
          </>
        ) : (
          <div className="chat-placeholder">
            <div className="placeholder-content">
              <h3>选择对话开始聊天</h3>
              <p>从左侧列表中选择一个对话或创建新对话</p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}