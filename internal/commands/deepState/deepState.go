package deepstate

import (
	"context"
	api "ds/internal/api/deepState"
	m "ds/internal/database/models"
	deepstate "ds/internal/services/deepState"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Goldziher/go-utils/sliceutils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/gorm"
)

const dateFormat = "2006-01-02"

var dict = map[api.AreaStatusType]string{
	api.Liberated:         "Звільнено",
	api.Occupied:          "Всього тимчасово окуповано",
	api.Occupied_after:    "Окуповано після 24 лютого 2022 ",
	api.Occupied_to:       "Окуповано до 24 лютого 2022",
	api.Unspecified:       "Статус невідомий",
	api.Other_territories: "Інші",
}

func ActualStateHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	res, err := deepstate.GetActualState()
	if err != nil {
		println(err.Error())
	}

	text := ""

	var totalOccupiedArea float32
	var totalOccupiedPercent float32
	occupied_before := sliceutils.Find(res, func(a api.Area, _ int, _ []api.Area) bool {
		return a.Type == api.Occupied_to
	})
	occupied_after := sliceutils.Find(res, func(a api.Area, _ int, _ []api.Area) bool {
		return a.Type == api.Occupied_after
	})
	if occupied_before != nil && occupied_after != nil {
		totalOccupiedArea = occupied_before.Area + occupied_after.Area
		pBefore, _ := strconv.ParseFloat(occupied_before.Percent, 32)
		pAfter, _ := strconv.ParseFloat(occupied_before.Percent, 32)
		totalOccupiedPercent = float32(pBefore + pAfter)
	}

	for _, areas := range res {
		status := dict[areas.Type]
		text = fmt.Sprintf("%s\n%s: %f (%s%%)", text, status, areas.Area, areas.Percent)
		if areas.Type == api.Occupied_after {

			text = fmt.Sprintf("%s\n%s: %f (%f%%)", text, dict[api.Occupied], totalOccupiedArea, totalOccupiedPercent)
		}
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      fmt.Sprintf("*Контроль території:*\n%s", bot.EscapeMarkdown(string(text))),
		ParseMode: models.ParseModeMarkdown,
	})
}

func HistoryHandler(db *gorm.DB) func(context.Context, *bot.Bot, *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		res, err := deepstate.GetHistory(db)
		if err != nil {
			println(err.Error())
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      fmt.Sprintf("Збережено нових записів: %d", res),
			ParseMode: models.ParseModeMarkdown,
		})
	}
}

func DateRequestHandler(db *gorm.DB) func(context.Context, *bot.Bot, *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {

		regex := *regexp.MustCompile(`^(\d{2}).(\d{2}).(\d{4})$`)
		parsed := regex.FindAllStringSubmatch(update.Message.Text, -1)
		if len(parsed) == 0 {
			return
		}
		r := parsed[0]
		year, _ := strconv.ParseInt(r[3], 10, 64)
		month, _ := strconv.ParseInt(r[2], 10, 64)
		day, _ := strconv.ParseInt(r[1], 10, 64)

		t := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, &time.Location{})

		var text string

		res := deepstate.GetRecordsByDate(db, t)
		if len(res) == 0 {
			text = "Нічого не знайдено"
		} else {
			text = strings.Join(sliceutils.Map(res, func(r m.HistoryRecord, _ int, _ []m.HistoryRecord) string {
				return fmt.Sprintf("%s: %s", r.CreatedAtDS.Format(dateFormat), r.Description)
			}), "\n")
		}

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      text,
			ParseMode: models.ParseModeHTML,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func HelpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	about := fmt.Sprintf(
		"%s\n[Підтримати наших захисників](https://t.me/help_deepstateua)",
		bot.EscapeMarkdown(`Бот, який надасть вам актуальні відомості про стан справ на фронті.
Джерело даних: https://deepstatemap.live`),
	)
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      about,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
