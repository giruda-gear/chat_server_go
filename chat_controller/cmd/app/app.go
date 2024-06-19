package app

import (
	"chat_controller_go/config"
	"chat_controller_go/network"
	"chat_controller_go/repository"
	"chat_controller_go/service"
)

type App struct {
	cfg *config.Config

	repository *repository.Repository
	service    *service.Service
	network    *network.Server
}

func NewApp(cfg *config.Config) *App {
	a := &App{cfg: cfg}

	var err error
	a.repository, err = repository.NewRepository(cfg)
	if err != nil {
		panic(err)
	}

	a.service = service.NewService(a.repository)
	a.network = network.NewNetwork(a.service, cfg.Info.Port)

	return a
}

func (a *App) Start() error {
	return a.network.Start()
}
