package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type tower struct {
	server *Server
}

func registerTower(server *Server) {
	t := &tower{server: server}

	t.server.engine.GET("/servers", t.handleGetAllServers)
}

func (t *tower) handleGetAllServers(c *gin.Context) {
	response(c, http.StatusOK, t.server.service.GetAvailableServerIPs())
}
