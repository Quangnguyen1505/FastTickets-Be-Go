// WebSocket Hub quản lý nhiều kết nối
package websocket

import "github.com/gorilla/websocket"

type Client struct {
	Conn      *websocket.Conn
	SessionID string
	Send      chan []byte
}

type Hub struct {
	clients    map[string][]*Client // sessionId -> list of clients
	register   chan *Client
	unregister chan *Client
	Broadcast  chan BroadcastMsg
}

type BroadcastMsg struct {
	SessionID string
	Message   []byte
}

var ChatHub = Hub{
	clients:    make(map[string][]*Client),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	Broadcast:  make(chan BroadcastMsg),
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.SessionID] = append(h.clients[client.SessionID], client)

		case client := <-h.unregister:
			conns := h.clients[client.SessionID]
			for i, c := range conns {
				if c == client {
					h.clients[client.SessionID] = append(conns[:i], conns[i+1:]...)
					break
				}
			}

		case msg := <-h.Broadcast:
			for _, client := range h.clients[msg.SessionID] {
				select {
				case client.Send <- msg.Message:
				default:
					// đóng connection nếu client bị lỗi
				}
			}
		}
	}
}

func (c *Client) writePump() {
	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.Conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
