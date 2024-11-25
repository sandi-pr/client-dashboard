package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"gitlab.hexacode.org/go-libs/microservice/config"
)

// NotificationMessage struct untuk pesan notifikasi.
type NotificationMessage struct {
	UserID  string `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

// NotificationWebSocket struct untuk mengelola WebSocket notifikasi.
type NotificationWebSocket struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan NotificationMessage
}

// NewNotificationWebSocket inisialisasi instance WebSocket notifikasi.
func NewNotificationWebSocket() *NotificationWebSocket {
	return &NotificationWebSocket{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan NotificationMessage),
	}
}

// AddClient menambahkan klien baru ke WebSocket.
func (ws *NotificationWebSocket) AddClient(client *websocket.Conn) {
	ws.Clients[client] = true
}

// RemoveClient menghapus klien dari WebSocket.
func (ws *NotificationWebSocket) RemoveClient(client *websocket.Conn) {
	delete(ws.Clients, client)
	client.Close()
}

// SendNotification mengirimkan pesan notifikasi ke semua klien yang terhubung.
func (ws *NotificationWebSocket) SendNotification(message NotificationMessage) {
	for client := range ws.Clients {
		err := client.WriteJSON(message)
		if err != nil {
			fmt.Println("Error sending notification:", err)
			ws.RemoveClient(client)
		}
	}
}

// HandleMessages menangani pengiriman pesan notifikasi menggunakan channel broadcast.
func (ws *NotificationWebSocket) HandleMessages() {
	for {
		msg := <-ws.Broadcast
		ws.SendNotification(msg)
	}
}

// ServeNotifications menangani koneksi WebSocket.
func (ws *NotificationWebSocket) ServeNotifications(ctx *config.ConfigContext, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection:", err)
		return
	}

	tokenEncoded, ok := ctx.URLParams.GetString("token")
	if !ok {
		fmt.Println("Unauthorized: No token provided")
		return
	}

	token, err := url.QueryUnescape(string(tokenEncoded))
	if err != nil {
		fmt.Println("Error decoding token:", err)
		return
	}

	idUser, err := getId(string(token))
	if err != nil {
		fmt.Println("Unauthorized")
		return
	}

	ws.AddClient(conn)
	defer ws.RemoveClient(conn)

	for {
		var msg NotificationMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		// Set UserID untuk pesan
		msg.UserID = strconv.Itoa(int(idUser))
		msg.Time = time.Now().Format(time.RFC3339)

		// Kirim pesan ke channel broadcast
		ws.Broadcast <- msg

		// Proses pesan jika perlu (misalnya simpan di database atau log).
		fmt.Printf("Notification received: %+v\n", msg)
	}
}
