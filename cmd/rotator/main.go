package main

import (
	"log"

	"github.com/FedoseevAlex/banner-rotation/internal/app"
	"github.com/FedoseevAlex/banner-rotation/internal/config"
	"github.com/FedoseevAlex/banner-rotation/internal/server"
)

func main() {
	cfg, err := config.ReadConfig("./configs/config.toml")
	if err != nil {
		log.Fatal(err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	srv, err := server.NewServer(application, cfg.Server)
	if err != nil {
		log.Fatal(err)
	}

	err = srv.Start()
	if err != nil {
		log.Fatal(err)
	}
}
