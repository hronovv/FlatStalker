package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"flat-stalker/internal/api"
	appbot "flat-stalker/internal/bot"
	"flat-stalker/internal/config"
	"flat-stalker/internal/db"
	"flat-stalker/internal/repository"
	"flat-stalker/internal/worker"

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

	listings := repository.NewListings(pool)
	seen := repository.NewSeenAds(pool)

	go tgBot.Start()
	go worker.New(cfg.Worker.Interval, listings, seen, tgBot).Start(ctx)

	gin.SetMode(cfg.App.GinMode)
	router := gin.Default()
	router.Use(api.CORS(cfg.CORS.Origins))
	router.POST("/message", tgBot.Message)

	links := &api.LinksHandler{
		Users:    repository.NewUsers(pool),
		Listings: listings,
	}
	links.Register(router)

	log.Printf("http listening on %s", cfg.Server.Addr)
	if err := router.Run(cfg.Server.Addr); err != nil {
		log.Fatal(err)
	}
}
