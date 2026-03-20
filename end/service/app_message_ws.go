package service

import (
	"encoding/json"
	"net/http"
	"ohome/global"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	appMessagePushEventType = "message_sync"

	appMessageWSWriteWait  = 10 * time.Second
	appMessageWSPongWait   = 60 * time.Second
	appMessageWSPingPeriod = (appMessageWSPongWait * 9) / 10
	appMessageWSSendBuffer = 16
)

type appMessagePushPayload struct {
	Event       string    `json:"event"`
	Reason      string    `json:"reason,omitempty"`
	UnreadCount int64     `json:"unreadCount"`
	MessageIDs  []uint    `json:"messageIds,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

type appMessageWSClient struct {
	userID uint
	conn   *websocket.Conn
	send   chan []byte
	hub    *appMessageWSHub
	once   sync.Once
}

type appMessageWSHub struct {
	mu sync.RWMutex

	upgrader websocket.Upgrader
	clients  map[uint]map[*appMessageWSClient]struct{}
}

var appMessagePushHub = newAppMessageWSHub()

func newAppMessageWSHub() *appMessageWSHub {
	return &appMessageWSHub{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients: make(map[uint]map[*appMessageWSClient]struct{}),
	}
}

func ServeAppMessageWS(c *gin.Context, userID uint) error {
	return appMessagePushHub.serve(c, userID)
}

func NotifyAppMessageOwners(ownerUserIDs []uint, reason string, messageIDs ...uint) {
	appMessagePushHub.notifyOwners(ownerUserIDs, reason, messageIDs)
}

func (h *appMessageWSHub) serve(c *gin.Context, userID uint) error {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return err
	}

	client := &appMessageWSClient{
		userID: userID,
		conn:   conn,
		send:   make(chan []byte, appMessageWSSendBuffer),
		hub:    h,
	}
	h.addClient(client)

	go client.writePump()
	go client.readPump()
	return nil
}

func (h *appMessageWSHub) notifyOwners(ownerUserIDs []uint, reason string, messageIDs []uint) {
	seen := make(map[uint]struct{}, len(ownerUserIDs))
	service := &AppMessageService{}
	now := time.Now()

	for _, ownerUserID := range ownerUserIDs {
		if ownerUserID == 0 {
			continue
		}
		if _, ok := seen[ownerUserID]; ok {
			continue
		}
		seen[ownerUserID] = struct{}{}

		clients := h.snapshotClients(ownerUserID)
		if len(clients) == 0 {
			continue
		}

		unreadCount, err := service.CountUnread(ownerUserID, "")
		if err != nil {
			if global.Logger != nil {
				global.Logger.Errorf("App Message WS CountUnread Error: user=%d err=%s", ownerUserID, err.Error())
			}
			continue
		}

		payload, err := json.Marshal(appMessagePushPayload{
			Event:       appMessagePushEventType,
			Reason:      reason,
			UnreadCount: unreadCount,
			MessageIDs:  messageIDs,
			Timestamp:   now,
		})
		if err != nil {
			if global.Logger != nil {
				global.Logger.Errorf("App Message WS Marshal Error: user=%d err=%s", ownerUserID, err.Error())
			}
			continue
		}

		for _, client := range clients {
			if client.trySend(payload) {
				continue
			}
			client.close()
		}
	}
}

func (h *appMessageWSHub) addClient(client *appMessageWSClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client.userID]; !ok {
		h.clients[client.userID] = make(map[*appMessageWSClient]struct{})
	}
	h.clients[client.userID][client] = struct{}{}
}

func (h *appMessageWSHub) removeClient(client *appMessageWSClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients, ok := h.clients[client.userID]
	if !ok {
		return
	}
	delete(clients, client)
	if len(clients) == 0 {
		delete(h.clients, client.userID)
	}
}

func (h *appMessageWSHub) snapshotClients(userID uint) []*appMessageWSClient {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.clients[userID]
	if !ok || len(clients) == 0 {
		return nil
	}

	result := make([]*appMessageWSClient, 0, len(clients))
	for client := range clients {
		result = append(result, client)
	}
	return result
}

func (c *appMessageWSClient) trySend(payload []byte) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()

	select {
	case c.send <- payload:
		return true
	default:
		return false
	}
}

func (c *appMessageWSClient) close() {
	c.once.Do(func() {
		c.hub.removeClient(c)
		close(c.send)
		_ = c.conn.Close()
	})
}

func (c *appMessageWSClient) readPump() {
	defer c.close()

	c.conn.SetReadLimit(512)
	_ = c.conn.SetReadDeadline(time.Now().Add(appMessageWSPongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(appMessageWSPongWait))
	})

	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (c *appMessageWSClient) writePump() {
	ticker := time.NewTicker(appMessageWSPingPeriod)
	defer func() {
		ticker.Stop()
		c.close()
	}()

	for {
		select {
		case payload, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(appMessageWSWriteWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(appMessageWSWriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
