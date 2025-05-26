package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ChatSocketHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("sessionId")
	if sessionId == "" {
		http.Error(w, "Missing sessionId", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil) // từ gorilla/websocket
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	client := &Client{
		Conn:      conn,
		SessionID: sessionId,
		Send:      make(chan []byte, 256),
	}
	ChatHub.register <- client

	// go client.readPump()  // nếu cần nhận message từ client
	go client.writePump() // ghi message từ server tới client
}
