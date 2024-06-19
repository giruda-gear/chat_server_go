package network

import (
	"chat_server_go/service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Room struct {
	Forward chan *message // send messages to Clients
	Join    chan *Client  // when socket is connected
	Leave   chan *Client  // when socket is disconnected
	Clients map[*Client]bool
	service *service.Service
}

type message struct {
	Name      string `json:"name"`
	Message   string `json:"message"`
	Room      string `json:"room"`
	CreatedAt int64  `json:"createdAt"`
}

type Client struct {
	Socket *websocket.Conn
	Send   chan *message
	Room   *Room
	Name   string `json:"name"`
}

func NewRoom(service *service.Service) *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
		service: service,
	}
}

func (c *Client) Read() {
	defer c.Socket.Close()
	for {
		var msg *message
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				break
			} else {
				panic(err)
			}
		} else {
			log.Println("READ:", msg, "CLIENT:", c.Name)
			msg.CreatedAt = time.Now().Unix()
			msg.Name = c.Name

			c.Room.Forward <- msg
		}
	}
}

func (c *Client) Write() {
	defer c.Socket.Close()
	for msg := range c.Send {
		log.Println("WRITE:", msg, "CLIENT:", c.Name)
		err := c.Socket.WriteJSON(msg)
		if err != nil {
			panic(err)
		}
	}
}

func (r *Room) Run() {
	for {
		select {
		case Client := <-r.Join:
			r.Clients[Client] = true
		case Client := <-r.Leave:
			r.Clients[Client] = false
			close(Client.Send)
			delete(r.Clients, Client)
		case msg := <-r.Forward:
			go r.service.InsertChatting(msg.Name, msg.Message, msg.Room)

			for Client := range r.Clients {
				Client.Send <- msg
			}
		}
	}
}

func (r *Room) ServeHttp(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatalln("---- serveHttp:", err)
	}

	authCookie, err := c.Request.Cookie("auth")
	if err != nil {
		panic(err)
	}

	Client := &Client{
		Socket: socket,
		Send:   make(chan *message, messageBufferSize),
		Room:   r,
		Name:   authCookie.Value,
	}

	r.Join <- Client

	defer func() { r.Leave <- Client }()

	go Client.Write()

	Client.Read()
}
