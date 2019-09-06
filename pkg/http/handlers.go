package http

import (
	"encoding/json"
	"lazer-twitter/persistence"
	"net/http"
	"strings"
	"sync"

	"github.com/fid-dev/go-pflog/log"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

type MessageHandler interface {
	CanHandle(message rawMessage) bool
	Handle(message rawMessage) ([]byte, bool, error)
}

func NewWebSocketHandler(db persistence.Database) *WebSocketHandler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	socketSlice := make([]*websocket.Conn, 0)

	return &WebSocketHandler{
		Database:          db,
		sockets:           socketSlice,
		websocketUpgrader: upgrader,
		messageHandlers: []MessageHandler{
			NewMessageHandler(db),
			NewJoinHandler(db),
			NewLikeHandler(db),
		},
	}
}

type rawMessage struct {
	Typ string `json:"typ"`
	Msg json.RawMessage
}

type WebSocketHandler struct {
	Database          persistence.Database
	sockets           []*websocket.Conn
	socketsMutex      sync.Mutex
	websocketUpgrader websocket.Upgrader
	messageHandlers   []MessageHandler
}

func (h *WebSocketHandler) SocketHandler(w http.ResponseWriter, r *http.Request) {

	h.websocketUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	serverSocket, err := h.websocketUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Errorf("could not establish socket connection %v", err)
		return
	}

	h.socketsMutex.Lock()
	h.sockets = append(h.sockets, serverSocket)
	h.socketsMutex.Unlock()
	if serverSocket != nil {
		h.handleData(serverSocket, w)
	}
}

func (h *WebSocketHandler) handleData(socket *websocket.Conn, w http.ResponseWriter) {
	for {
		_, messageByte, err := socket.ReadMessage()

		if err != nil {
			if strings.Contains(err.Error(), "close 1001") {
				h.removeSocket(socket)
				return
			}
			log.Errorf("could not read message %v", err)
			continue
		}
		raw := rawMessage{}
		err = json.Unmarshal(messageByte, &raw)
		raw.Msg = messageByte
		if err != nil {
			log.Errorf("could not convert to struct %v", err)
			continue
		}
		if len(raw.Typ) <= 0 {
			err := ErrorMessage{
				Typ:     "error",
				Message: "Empty type, please check the tweet requirements",
			}
			byteErr, _ := json.Marshal(err)
			_ = socket.WriteMessage(1, byteErr)
		}

		for _, val := range h.messageHandlers {
			if val.CanHandle(raw) {
				msg, multi, err := val.Handle(raw)

				if err != nil {
					log.Errorf("could not handle %v", err)
					continue
				}

				if multi != true {
					err = socket.WriteMessage(1, msg)
					if err != nil {
						log.Errorf("could not send tweet %v", err)
						continue
					}
				} else {
					h.socketsMutex.Lock()
					for _, val := range h.sockets {
						err = val.WriteMessage(1, msg)
						if err != nil {
							log.Errorf("could not send tweet to client %v", err)
							continue
						}
					}
					h.socketsMutex.Unlock()
				}
			}
		}
	}
}

func (h *WebSocketHandler) removeSocket(socket *websocket.Conn) {
	h.socketsMutex.Lock()
	for i, val := range h.sockets {
		if val == socket {
			h.sockets = append(h.sockets[:i], h.sockets[i+1:]...)
		}
	}
	h.socketsMutex.Unlock()
}
