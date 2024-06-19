package network

import (
	"chat_server_go/service"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine

	service *service.Service

	port string
	ip   string
}

func NewNetwork(service *service.Service, port string) *Server {
	s := &Server{engine: gin.New(), service: service, port: port}

	s.engine.Use(gin.Logger())
	s.engine.Use(gin.Recovery())
	s.engine.Use(cors.New(cors.Config{
		AllowWebSockets:  true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"ORIGIN", "Cotent-Length", "Content-Type", "Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Authorization", "X-Requested-With", "expires"},
		ExposeHeaders:    []string{"ORIGIN", "Cotent-Length", "Content-Type", "Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Authorization", "X-Requested-With", "expires"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	}))

	registerServer(s)

	return s
}

func (s *Server) setServerInfo() {
	if addrs, err := net.InterfaceAddrs(); err != nil {
		panic(err.Error())
	} else {
		var ip net.IP

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					ip = ipnet.IP
					break
				}
			}
		}

		if ip == nil {
			panic("No IP address found")
		} else {
			if err = s.service.ServerSet(ip.String()+s.port, true); err != nil {
				panic(err)
			} else {
				s.ip = ip.String()
			}

			s.service.PublishServerStatusEvent(s.ip+s.port, true)
		}
	}
}

func (s *Server) StartServer() error {
	s.setServerInfo()

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGINT)

	// if server goes down
	go func() {
		<-channel

		if err := s.service.ServerSet(s.ip+s.port, false); err != nil {
			// TODO handle a fail case
			log.Println("Failed to set server info", "err:", err)
		}

		s.service.PublishServerStatusEvent(s.ip+s.port, false)

		os.Exit(1)
	}()

	log.Println("Starting Server! port:", s.port)
	return s.engine.Run(s.port)
}
