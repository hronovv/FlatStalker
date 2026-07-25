package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	appbot "flat-stalker/internal/bot"
	"flat-stalker/internal/config"
	"flat-stalker/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	pool := db.MustNewPool(ctx, cfg.Database)
	defer pool.Close()
	log.Println("database connected")

	tgBot, err := appbot.New(ctx, cfg, pool)
	if err != nil {
		panic(err)
	}

	gin.SetMode(cfg.App.GinMode)
	router := gin.Default()
	router.POST("/message", tgBot.Message)

	go tgBot.Start()

	log.Printf("http listening on %s", cfg.Server.Addr)
	if err := router.Run(cfg.Server.Addr); err != nil {
		log.Fatal(err)
	}
}
