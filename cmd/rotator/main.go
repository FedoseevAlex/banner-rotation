package main

import (
	"flag"
	"log"

	"github.com/FedoseevAlex/banner-rotation/internal/app"
	"github.com/FedoseevAlex/banner-rotation/internal/config"
	"github.com/FedoseevAlex/banner-rotation/internal/server"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "Path to config file.")
}

func main() {
	flag.Parse()

	cfg, err := config.ReadConfig(configPath)
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
