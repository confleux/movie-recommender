package main

import (
	"context"
	"log"
	"tg-bot/internal/handler"
	"tg-bot/internal/repository"
	"tg-bot/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"tg-bot/internal/config"
)

func main() {
	cfg := config.MustLoad()

	pool, err := pgxpool.New(context.Background(), cfg.Postgres.Uri)
	if err != nil {
		log.Fatalf("postges connection: %v", err)
	}
	defer pool.Close()

	movieRepository := repository.NewMovieRepository(pool)

	movieService := service.NewMovieService(movieRepository)

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatalf("tg bot api initialization: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updateHandler := handler.NewUpdateHandler(bot, movieService)

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		go func() {
			err := updateHandler.ProcessUpdate(&update)
			if err != nil {
				log.Printf("process update: %v", err)
			}
		}()
	}
}
