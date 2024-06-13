package main

import (
	"chat_server_go/config"
	"chat_server_go/network"
	"chat_server_go/repository"
	"chat_server_go/service"
	"flag"
)

var pathFlag = flag.String("config", "./config.toml", "config set")
var port = flag.String("port", ":1010", "port set")

func main() {
	flag.Parse()
	c := config.NewConfig(*pathFlag)

	if rep, err := repository.NewRepository(c); err != nil {
		panic(err)
	} else {
		s := network.NewServer(service.Newservice(rep), rep, *port)
		s.StartServer()
	}
}
