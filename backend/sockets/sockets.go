package sockets

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins for now. In production, you should restrict this.
    },
}

type Client struct {
    conn *websocket.Conn
    send chan []byte
}

type TimerState struct {
    Duration  int64     `json:"duration"`
    StartTime time.Time `json:"startTime"`
}

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    timerState TimerState
    mutex      sync.Mutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}

func (h *Hub) startTimer(duration int64) {
    h.mutex.Lock()
    defer h.mutex.Unlock()

    h.timerState = TimerState{
        Duration:  duration,
        StartTime: time.Now(),
    }

    message, _ := json.Marshal(map[string]interface{}{
        "type":      "timer_update",
        "duration":  duration,
        "startTime": h.timerState.StartTime,
    })
    h.broadcast <- message
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    client := &Client{conn: conn, send: make(chan []byte, 256)}
    hub.register <- client

    go client.writePump(hub)
    go client.readPump(hub)
}

func (c *Client) readPump(hub *Hub) {
    defer func() {
        hub.unregister <- c
        c.conn.Close()
    }()

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }

        var data map[string]interface{}
        if err := json.Unmarshal(message, &data); err != nil {
            log.Println("Error unmarshalling message:", err)
            continue
        }

        if data["type"] == "start_timer" {
            duration := int64(data["duration"].(float64))
            hub.startTimer(duration)
        }
    }
}

func (c *Client) writePump(hub *Hub) {
    defer func() {
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            if err := w.Close(); err != nil {
                return
            }
        }
    }
}

// func main() {
//     hub := newHub()
//     go hub.run()

//     http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
//         serveWs(hub, w, r)
//     })

//     log.Println("Server starting on :8080")
//     err := http.ListenAndServe(":8080", nil)
//     if err != nil {
//         log.Fatal("ListenAndServe: ", err)
//     }
// }