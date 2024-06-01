package main

import (
	"context"
	deepstate "ds/internal/commands/deepState"
	"ds/internal/config"
	"ds/internal/database"

	"github.com/go-telegram/bot"
)

type HistoryRecord struct {
	Id            int    `json:"id"`
	DescriptionEn string `json:"descriptionEn"`
	Description   string `json:"description"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
	Datetime      string `json:"datetime"`
	Status        bool   `json:"status"`
}

func main() {
	loadErr := config.Load()

	if loadErr != nil {
		panic(loadErr)
	}
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	database.Migrate(db)

	opts := []bot.Option{
		bot.WithDefaultHandler(deepstate.DateRequestHandler(db)),
	}
	b, err := bot.New(config.Config(config.BOT_TOKEN_KEY), opts...)
	if err != nil {
		panic(err)
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/state", bot.MatchTypeExact, deepstate.ActualStateHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/parse", bot.MatchTypeExact, deepstate.HistoryHandler(db))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/map", bot.MatchTypeExact, deepstate.GeoHistoryHandler(db))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/about", bot.MatchTypeExact, deepstate.HelpHandler)

	b.Start(context.Background())

}
