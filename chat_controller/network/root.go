package network

import (
	"chat_controller_go/service"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine

	service *service.Service

	port string
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

	registerTower(s)

	return s
}

func (s *Server) Start() error {
	log.Println("Start Controller Server!!")
	return s.engine.Run(s.port)
}
