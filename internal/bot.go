package internal

import (
	"context"
	deepstate "ds/internal/commands/deepState"
	"ds/internal/config"
	"ds/internal/database"

	"github.com/go-telegram/bot"
	"gorm.io/gorm"
)

func RunBot() {
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

	registerHandlers(b, db)

	b.Start(context.Background())
	defer b.Close(context.Background())
}

func registerHandlers(b *bot.Bot, db *gorm.DB) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/state", bot.MatchTypeExact, deepstate.ActualStateHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/parse", bot.MatchTypeExact, deepstate.HistoryHandler(db))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/map", bot.MatchTypeExact, deepstate.GeoHistoryHandler(db))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/about", bot.MatchTypeExact, deepstate.HelpHandler)
}
