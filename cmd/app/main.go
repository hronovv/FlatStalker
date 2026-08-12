package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"flat-stalker/internal/api"
	appbot "flat-stalker/internal/bot"
	"flat-stalker/internal/config"
	"flat-stalker/internal/db"
	"flat-stalker/internal/plan"
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

	listings := repository.NewListings(pool)
	seen := repository.NewSeenAds(pool)
	users := repository.NewUsers(pool)

	tgBot, err := appbot.New(ctx, cfg, users, listings)
	if err != nil {
		panic(err)
	}

	intervals := plan.Intervals{
		Free: cfg.Worker.Free,
		Plus: cfg.Worker.Plus,
		Pro:  cfg.Worker.Pro,
	}

	go tgBot.Start()
	go worker.New(intervals, listings, seen, tgBot).Start(ctx)

	gin.SetMode(cfg.App.GinMode)
	router := gin.Default()
	router.Use(api.CORS(cfg.CORS.Origins))

	apiGroup := router.Group("/api")
	apiGroup.Use(api.TelegramAuth(cfg.Telegram.BotToken, 24*time.Hour))

	links := &api.LinksHandler{
		Users:    users,
		Listings: listings,
	}
	links.Register(apiGroup)

	me := &api.MeHandler{
		Users:     users,
		Intervals: intervals,
	}
	me.Register(apiGroup)

	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	log.Printf("http listening on %s", cfg.Server.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
