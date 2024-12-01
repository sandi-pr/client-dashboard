package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

// NotificationMessage struct untuk pesan notifikasi.
type NotificationMessage struct {
	UserID  string `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Time    string `json:"time"`
	Type    string `json:"type"`
	Meta    U.Dict `json:"meta"`
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
	if strings.Contains(message.Type, "progress") {
		if message.Type == "progress" {
			err := ws.SaveNotification(message)

			if err != nil {
				fmt.Println("Error saving notification:", err)
			}
		} else {
			err := UpdateContentProgressByName(message.Title, message.Content)
			if err != nil {
				fmt.Println("Error saving notification:", err)
			}
		}
	} else {
		if message.Content == "Hexacube has finished analyzing the uploaded documents. You can now proceed to explore the insights!" {
			fmt.Println("ready notification")
		} else {
			err := ws.SaveNotification(message)

			if err != nil {
				fmt.Println("Error saving notification:", err)
			}
		}
	}

	for client := range ws.Clients {
		err := client.WriteJSON(message)
		if err != nil {
			fmt.Println("Error sending notification:", err)
			ws.RemoveClient(client)
		}
	}
}

func (ws *NotificationWebSocket) SaveNotification(message NotificationMessage) error {
	strapi := strapi.NewStrapiWithToken(STRAPI_URL, STRAPI_TOKEN)

	status, _, err := strapi.Create("notification-table", U.Dict{
		"title":   message.Title,
		"content": message.Content,
		"user_id": message.UserID,
		"sent_at": message.Time,
		"type":    message.Type,
		"read":    false,
		"meta":    message.Meta,
	})

	if err != nil {
		fmt.Println("Error saving notification:", err)
	}

	if status != 200 {
		fmt.Println("Error saving notification:", status)
	}

	return nil
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

func SendNotification(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	data, ok := ctx.Params.GetString("data")

	if !ok {
		return ctx.Error(500, "Internal Server Error")
	}

	decodedData, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		fmt.Printf("Error decoding Base64 data: %v\n", err)
		return ctx.Error(400, "Invalid Base64 encoding in 'data' parameter")
	}

	fmt.Printf("Decoded JSON data: %s\n", decodedData)

	var dataJson NotificationMessage
	err = json.Unmarshal([]byte(decodedData), &dataJson)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(400, "invalid parameters data")
	}

	notificationWS.SendNotification(dataJson)

	return 200, U.Dict{}, nil
}
