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
    IsRunning bool      `json:"isRunning"`
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
	log.Printf("startTimer duration %v", duration)
    h.mutex.Lock()
    defer h.mutex.Unlock()

    h.timerState = TimerState{
        Duration:  duration,
        StartTime: time.Now(),
    }

    message, _ := json.Marshal(map[string]interface{}{
        "type":      "timer_update",
        "time":  duration,
        "startTime": h.timerState.StartTime,
    })
    h.broadcast <- message
}

func (h *Hub) sendTimerState(client *Client) {
    h.mutex.Lock()
    defer h.mutex.Unlock()

    message, _ := json.Marshal(map[string]interface{}{
        "type":       "timer_update",
        "timerState": h.timerState,
    })
    client.send <- message
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    client := &Client{conn: conn, send: make(chan []byte, 256)}
    hub.register <- client

	hub.sendTimerState(client)

    go client.writePump(hub)
    go client.readPump(hub)
}


func (h *Hub) handleTimerEvent(eventType string, duration int64) {
    h.mutex.Lock()
    defer h.mutex.Unlock()

    switch eventType {
    case "start":
        h.timerState = TimerState{
            Duration:  duration,
            StartTime: time.Now(),
            IsRunning: true,
        }
    case "pause":
        h.timerState.IsRunning = false
        h.timerState.Duration = duration
        h.timerState.StartTime = time.Now()
    case "stop":
        h.timerState = TimerState{} // Reset timer state
    }

    message, _ := json.Marshal(map[string]interface{}{
        "type":       "timer_update",
        "timerState": h.timerState,
    })
    h.broadcast <- message
}

func (c *Client) readPump(hub *Hub) {
	log.Println("readPump")
    defer func() {
        hub.unregister <- c
        c.conn.Close()
    }()

    for {
        _, message, err := c.conn.ReadMessage()
		log.Printf("readPump message %v, error %v", message, err)
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

        if eventType, ok := data["type"].(string); ok {
			duration := int64(data["duration"].(float64))
			hub.handleTimerEvent(eventType, duration)
		}
    }
}

func (c *Client) writePump(hub *Hub) {
	log.Println("writePump")
    defer func() {
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
			log.Printf("writePump message %v, ok %v", message, ok)
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