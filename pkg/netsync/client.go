package netsync

import (
	"context"
	"fmt"
	"log"
	"main/pkg/api"
	"main/pkg/serial"

	"nhooyr.io/websocket"
)

// TODO: need to make game data concurrency safe

const maxMessageSize = 2048

type UpdateHandler func(*serial.EventWrapper) error

type Client struct {
	addr         string
	port         string
	Username     string
	HandleUpdate UpdateHandler
	conn         *websocket.Conn
	send         chan []byte
}

type ApiPositionUpdate struct {
	Username string         `json:"name"`
	Position api.ApiVector2 `json:"position"`
}

func NewClient(addr, port, username string) (c *Client) {
	return &Client{
		addr:     addr,
		port:     port,
		Username: username,
		send:     make(chan []byte, 256),
	}
}

func (c *Client) WriteMessage(event serial.EventWrapper) {
	b, err := event.Serialize()
	if err != nil {
		log.Println("serialization error")
	}
	c.send <- b
}

func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()
	c.conn.SetReadLimit(maxMessageSize)

	for {
		_, message, err := c.conn.Read(ctx)
		if err != nil {
			if true {
				log.Println("unexpected-close")
			}
			break
		}
		wrapper, err := serial.Deserialize(message)
		if err != nil {
			log.Println("deserialization error")
			continue
		}
		err = c.HandleUpdate(wrapper)
		if err != nil {
			log.Println("update handling error")
			continue
		}
	}
}

func (c *Client) writePump(ctx context.Context) {
	defer func() {
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				log.Println("send channel closed")
				return
			}

			w, err := c.conn.Writer(ctx, websocket.MessageText)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func (c *Client) Start() {
	ctx := context.Background()

	urlFormat := "ws://%s:%s/ws?username=%s"
	url := fmt.Sprintf(urlFormat, c.addr, c.port, c.Username)
	// url := "ws://localhost:8087/ws?username=ted"
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	// defer conn.Close(websocket.StatusInternalError, "connection closed")

	c.conn = conn

	go c.writePump(ctx)
	go c.readPump(ctx)
}
