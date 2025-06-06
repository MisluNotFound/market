import { useEffect, useState, useRef, useCallback } from 'react';
import { useLocation } from 'react-router-dom';
import AuthService from '../services/auth';
import ProductService from '../services/product';
import '../styles/chat.css';

export default function Chat() {
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

  const userId = localStorage.getItem('userId');


  const parseProductMessage = (content) => {
    const parts = content.split(',');
    const description = parts[0].replace('describe ', '');
    const avatar = parts[1].replace('avatar ', '');
    const price = parts[2].replace('price ', '');
    return { description, avatar, price };
  };

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
      console.log(message);

      // 处理消息确认
      if (message.type === 2) {
        const conversationId = message.to;
        setConversationMessages(prev => {
          const newMap = new Map(prev);
          console.log("map", newMap);
          const currentMessages = newMap.get(conversationId) || [];
          const updatedMessages = currentMessages.map(msg => {
            if (msg.tempID === message.tempID) {
              console.log('找到匹配的临时消息:', msg);
              return { ...msg, id: message.id };
            }
            return msg;
          });
          console.log("updated", updatedMessages);
          newMap.set(conversationId, updatedMessages);
          return newMap;
        });
        return;
      }

      // 处理撤回消息
      if (message.type === 4) {
        const conversationId = message.from;
        setConversationMessages(prev => {
          const newMap = new Map(prev);
          const currentMessages = newMap.get(conversationId) || [];
          const updatedMessages = currentMessages.filter(msg => msg.id !== message.id);
          console.log(currentMessages);
          newMap.set(conversationId, updatedMessages);
          return newMap;
        });
        return;
      }

      // 对于其他类型的消息，使用from字段来确定会话ID
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
        let newMessage;

        if (message.media_type === 'link') {
          const productInfo = parseProductMessage(message.content);
          newMessage = {
            id: message.id || message.tempID || Date.now().toString(),
            tempID: message.tempID,
            from: message.from_user_id,
            content: productInfo,
            mediaType: 'link',
            time: message.timestamp || new Date().toISOString(),
            type: message.type || 1
          };
        } else {
          newMessage = {
            id: message.id || message.tempID || Date.now().toString(),
            tempID: message.tempID,
            from: message.from_user_id,
            content: message.content,
            time: message.timestamp || new Date().toISOString(),
            type: message.type || 1
          };
        }

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
      // 使用 IMService 发送消息并获取 tempID
      const tempID = imService.current?.sendMessage(input, activeConversation.id);

      if (tempID) {
        const newMessage = {
          tempID: tempID,
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

        setInput('');
      }
    }
  };

  const selectConversation = async (conversation) => {
    console.log("select conversation", conversation)
    setActiveConversation({
      id: conversation.id,
      name: conversation.username,
      avatar: conversation.avatar,
    });

    try {
      const msgResponse = await imService.current.getConversationMessages(
        conversation.fromUserID,
        conversation.toUserID
      );

      console.log('获取消息响应:', msgResponse.data.Messages);
      const sortedMessages = sortMessages(
        (msgResponse?.data?.Messages || []).map(msg => ({
          id: msg.id,
          tempID: msg.tempID,
          from: msg.from_user_id,
          content: msg.content,
          time: msg.timestamp,
          mediaType: msg.media_type,
          type: msg.type || 1,
          to: msg.to_user_id
        }))
      );
      setConversationMessages(prev => {
        const newMap = new Map(prev);
        newMap.set(conversation.id, sortedMessages);
        console.log(newMap)
        return newMap;
      });
    } catch (error) {
      console.error('获取消息失败:', error);
    }
  };

  const handleRecall = (msg) => {
    console.log(msg)
    if (activeConversation) {
      const messageId = msg.id
      // 发送撤回消息，type为4，content为要撤回的消息ID
      imService.current?.sendMessage(messageId, activeConversation.id, 'text', 4);

      // 从本地消息列表中移除被撤回的消息
      setConversationMessages(prev => {
        const newMap = new Map(prev);
        const currentMessages = newMap.get(activeConversation.id) || [];
        const updatedMessages = currentMessages.filter(msg => msg.id !== messageId);
        newMap.set(activeConversation.id, updatedMessages);
        return newMap;
      });
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
            {conversation.unreadCount > 0 && (
              <div className="unread-badge">{conversation.unreadCount}</div>
            )}
          </div>
        ))}
      </div>

      <div className="chat-area">
        {activeConversation ? (
          <>
            <div className="chat-header">
              <div className="chat-user">
                <img
                  src={activeConversation.avatar || '/placeholder-user.png'}
                  className="chat-user-avatar"
                  alt="用户头像"
                />
                <div className="chat-title">{activeConversation.name}</div>
              </div>
            </div>
            <div className="chat-messages">
              {messages.map((msg, index) => (
                <div
                  key={index}
                  className={`message-container ${msg.from === userId ? 'sent' : 'received'}`}
                  onContextMenu={(e) => {
                    e.preventDefault();
                    const menu = document.getElementById(`message-menu-${index}`);
                    if (menu) {
                      menu.style.display = 'block';
                      menu.style.left = `${e.clientX}px`;
                      menu.style.top = `${e.clientY}px`;
                    }
                  }}
                  onMouseEnter={() => {
                    const recall = document.getElementById(`recall-${index}`);
                    if (recall) recall.style.display = 'block';
                  }}
                  onMouseLeave={() => {
                    const recall = document.getElementById(`recall-${index}`);
                    if (recall) recall.style.display = 'none';
                  }}
                >
                  {msg.from === userId && (
                    <div
                      id={`recall-${index}`}
                      className="recall-btn"
                      style={{ display: 'none' }}
                      onClick={() => handleRecall(msg)}
                    >
                      撤回
                    </div>
                  )}
                  {msg.from !== userId && (
                    <div
                      className="message-avatar"
                      style={{ backgroundImage: `url(${activeConversation.avatar})` }}
                    />
                  )}
                  <div className={`message ${msg.from === userId ? 'sent' : 'received'}`}>
                    {msg.mediaType === 'link' ? (
                      <div className="product-message">
                        <img src="http://127.0.0.1:3200/api/assert/1/7bed6ab8-378f-11f0-939b-0242ac150003/888cb4aa35be5d3c.jpg" alt="商品图片" className="product-image" />
                        <div className="product-info">
                          <div className="product-description">
                            {msg.content.length > 20 ? `${msg.content.substring(0, 20)}...` : " Xbox 冰雪白游戏手柄"}
                          </div>
                          <div className="product-price">¥299</div>
                        </div>
                      </div>
                    ) : (
                      <div className="message-content">{msg.content}</div>
                    )}
                    <div className="message-time">
                      {new Date(msg.time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                    </div>
                  </div>
                  <div id={`message-menu-${index}`} className="message-menu" style={{ display: 'none' }}>
                    <div className="menu-item">撤回消息</div>
                  </div>
                </div>
              ))}
            </div>

            <div className="chat-input">
              {/* <div className="chat-actions">
                <button className="chat-action-btn">发送商品</button>
                <button className="chat-action-btn">发送订单</button>
              </div> */}
              <div className="input-container">
                <textarea
                  value={input}
                  onChange={(e) => setInput(e.target.value)}
                  placeholder="输入消息..."
                  onKeyPress={(e) => e.key === 'Enter' && !e.shiftKey && handleSend()}
                />
                <button className="send-btn" onClick={handleSend}>发送</button>
              </div>
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