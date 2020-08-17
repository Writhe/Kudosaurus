package main

import (
	"log"
	"os"

	"github.com/writhe/kudosaurus/internal/config"
	"github.com/writhe/kudosaurus/internal/server"
	"github.com/writhe/kudosaurus/internal/sqlitestore"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ltime|log.Ldate|log.LUTC)

	cfg, err := config.GetConfig("config.yml")

	if err != nil {
		panic(err)
	}

	store := sqlitestore.NewSource(cfg.Database.Path, logger)

	serverConfig := server.CommandServerConfig{
		MaxKudos:          cfg.Settings.MaxKudos,
		Token:             cfg.Slack.Token,
		SigningSecret:     cfg.Slack.SigningSecret,
		VerificationToken: cfg.Slack.VerificationToken,
		Port:              cfg.Server.Port,
	}

	server := server.New(store, serverConfig, logger)

	server.Start()
}
