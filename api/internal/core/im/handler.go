package im

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mislu/market-api/internal/types/request"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn      *websocket.Conn
	send      chan request.Message
	userID    string
	pendingMu sync.Mutex                 // 队列锁
	pending   map[string]*PendingMessage // 消息ID->消息
}

// 待确认消息结构
type PendingMessage struct {
	request.Message
	timer     *time.Timer
	expiresAt time.Time
	retries   uint
}

var (
	clients   = make(map[string]*Client) // userID -> Client
	clientsMu sync.Mutex
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 在升级连接前验证 userID
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		w.Write([]byte("missing user id"))
		return
	}

	// 升级 HTTP 连接为 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket 升级错误: %v", err)
		return
	}

	// 设置连接关闭处理程序
	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("用户 %s 的 WebSocket 连接关闭: %d %s", userID, code, text)
		removeClient(userID)
		return nil
	})

	// 创建客户端
	client := &Client{
		conn:    conn,
		send:    make(chan request.Message, 256),
		userID:  userID,
		pending: make(map[string]*PendingMessage, 0),
	}

	// 添加客户端到管理
	addClient(userID, client)

	// 启动读写协程
	go client.writePump()
	go client.readPump()
}

func addClient(userID string, client *Client) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	clients[userID] = client
}

func removeClient(userID string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	delete(clients, userID)
}

func getClient(userID string) (*Client, bool) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	client, ok := clients[userID]
	return client, ok
}

func (c *Client) readPump() {
	defer c.conn.Close()

	for {
		messageType, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("用户 %s 的 WebSocket 读取错误: %v", c.userID, err)
			}
			break
		}

		log.Printf("from 用户 %s 收到原始消息 (类型 %d): %s", c.userID, messageType, string(data))

		var msg request.Message
		if messageType == websocket.TextMessage {
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Printf("用户 %s 的 JSON 解析错误: %v, 原始数据: %s", c.userID, err, string(data))
				continue
			}

			c.handleMessage(&msg)
		} else {
			log.Printf("用户 %s 收到非文本消息 (类型 %d)，暂不支持", c.userID, messageType)
			continue // 跳过非文本消息
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.send:
			log.Printf("用户%s 收到消息\n", c.userID)
			if !ok {
				// 通道关闭
				log.Println("通道关闭")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(msg); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
			log.Println("写入用户ws连接")

		case <-ticker.C:
			// 定期发送 ping 保持连接
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("ping error ", err)
				return
			}
		}
	}
}

func Init() {
	http.HandleFunc("/api/im/ws", HandleWebSocket)

	server := &http.Server{
		Addr:    ":3300",
		Handler: nil,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("服务器运行错误: ", err)
		}
	}()
}
