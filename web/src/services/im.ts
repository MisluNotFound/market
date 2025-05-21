import axios from 'axios';

declare const crypto: Crypto;

type MessageHandler = (message: Message) => void;

interface Message {
  tempID?: string;
  from: string;
  to: string;
  content: string;
  mediaType: string;
  type: number;
}

interface Conversation {
  FromUserID: string;
  ToUserID: string;
  LastMessageContent: string;
  LastMessageTime: string;
  LastReadMessageID: string;
  MarkDeleted: boolean;
  CurrentProductID: string;
  unreadCount: number;
  id: string;
  username: string;
  school: string;
  phone: string;
  gender: string;
  avatar: string;
  sellerCredit: number;
  purchaseCredit: number;
  isCertificated: boolean;
  address: string;
}

interface MessageListResponse {
  Messages: Array<{
    id: string;
    conversation_id: string;
    from_user_id: string;
    to_user_id: string;
    content: string;
    media_type: string;
    timestamp: string;
  }>;
}

interface PageReq {
  page: number;
  size: number;
}

class IMService {
  private static instance: IMService | null = null;
  private userId: string;
  private ws: WebSocket | null;
  private pendingMessages: Map<string, Message>;
  private messageHandlers: Set<MessageHandler>;

  private constructor(userId: string) {
    this.userId = userId;
    this.ws = null;
    this.pendingMessages = new Map();
    this.messageHandlers = new Set();
    this.connect();
  }

  private connect(): void {
    if (!this.userId) {
      throw new Error('UserID is required for WebSocket connection');
    }

    const wsUrl = `ws://localhost:3300/api/im/ws?userID=${this.userId}`;
    console.log('Connecting to WebSocket:', wsUrl);

    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      console.log('WebSocket connected successfully');
    };

    this.ws.onmessage = (event) => {
      console.log('Received WebSocket message:', event.data);
      try {
        const message = JSON.parse(event.data) as Message;
        console.log('Parsed message:', message);
        this.handleMessage(message);
      } catch (error) {
        console.error('Error parsing message:', error, 'Raw data:', event.data);
      }
    };

    this.ws.onclose = () => {
      console.log('WebSocket disconnected');
    };

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  }

  sendMessage(content: string, to: string, mediaType = 'text'): string {
    const tempID = crypto.randomUUID();
    const message: Message = {
      tempID,
      from: this.userId,
      to,
      content,
      mediaType,
      type: 1
    };

    this.pendingMessages.set(tempID, message);
    this.ws?.send(JSON.stringify(message));
    return tempID;
  }

  private handleMessage(message: Message): void {
    switch (message.type) {
      case 2: // ACK
        this.pendingMessages.delete(message.tempID!);
        break;
      case 3: // FAIL
        const pendingMsg = this.pendingMessages.get(message.tempID!);
        if (pendingMsg) {
          this.ws?.send(JSON.stringify(pendingMsg));
        }
        break;
      default:
        this.notifyHandlers(message);
    }
  }

  addMessageHandler(handler: MessageHandler): () => void {
    this.messageHandlers.add(handler);
    return () => this.messageHandlers.delete(handler);
  }

  private notifyHandlers(message: Message): void {
    console.log('Notifying handlers for message:', message);
    console.log('Current handler count:', this.messageHandlers.size);
    this.messageHandlers.forEach(handler => {
      console.log('Calling handler:', handler);
      handler(message);
    });
  }

  async createConversation(fromUserID: string, toUserID: string, productID: string): Promise<any> {
    const response = await axios.post('http://localhost:3200/api/conversation/create', {
      fromUserID,
      toUserID,
      productID
    });
    return response.data;
  }

  async getConversationList(userID: string): Promise<Conversation[]> {
    console.log(userID)
    const response = await axios.get(`http://localhost:3200/api/conversation/${userID}/list`);
    return response.data.data;
  }

  async getConversationMessages(
    fromUserID: string,
    toUserID: string,
    pageReq: PageReq = { page: 1, size: 20 }
  ): Promise<MessageListResponse> {
    const response = await axios.get('http://localhost:3200/api/conversation/messages', {
      params: {
        fromUserID,
        toUserID,
        ...pageReq
      }
    });
    return response.data;
  }

  public close(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  public static getInstance(userId: string): IMService {
    if (!IMService.instance) {
      IMService.instance = new IMService(userId);
    }
    return IMService.instance;
  }
}

export { IMService };