package network

import (
	"chat_server_go/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type api struct {
	server *Server
}

func registerServer(server *Server) {
	a := &api{server: server}

	server.engine.GET("/rooms", a.roomList)
	server.engine.POST("/room", a.makeRoom)
	server.engine.GET("/room", a.room)
	server.engine.GET("/enter-room", a.enterRoom)

	// r := NewRoom()
	// go r.Run()

	// server.engine.GET("/room", r.ServeHttp)
}

func (a *api) roomList(c *gin.Context) {
	if res, err := a.server.service.RoomList(); err != nil {
		response(c, http.StatusInternalServerError, err.Error())
	} else {
		response(c, http.StatusOK, res)
	}
}

func (a *api) makeRoom(c *gin.Context) {
	var req types.BodyRoomReq

	if err := c.ShouldBindJSON(&req); err != nil {
		response(c, http.StatusUnprocessableEntity, err.Error())
	} else if err := a.server.service.MakeRoom(req.Name); err != nil {
		response(c, http.StatusInternalServerError, err.Error())
	} else {
		response(c, http.StatusOK, "Success")
	}
}

func (a *api) room(c *gin.Context) {
	var req types.FormRoomReq

	if err := c.ShouldBindQuery(&req); err != nil {
		response(c, http.StatusUnprocessableEntity, err.Error())
	} else if res, err := a.server.service.Room(req.Name); err != nil {
		response(c, http.StatusInternalServerError, err.Error())
	} else {
		response(c, http.StatusOK, res)
	}
}

func (a *api) enterRoom(c *gin.Context) {
	var req types.FormRoomReq

	if err := c.ShouldBindQuery(&req); err != nil {
		response(c, http.StatusUnprocessableEntity, err.Error())
	} else if res, err := a.server.service.EnterRoom(req.Name); err != nil {
		response(c, http.StatusInternalServerError, err.Error())
	} else {
		response(c, http.StatusOK, res)
	}
}