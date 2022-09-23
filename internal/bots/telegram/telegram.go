package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ssstatistics/internal/config"
)

var Bot *TgBot

func Initialize() {
	Bot = &TgBot{}
	go Bot.Initialize()
}

type TgBot struct {
	Self *tgbotapi.BotAPI
}

func (t *TgBot) Send(s string) error {
	for _, id := range config.Config.TelegramBot.TrustedChatIDs {
		msg := tgbotapi.NewMessage(int64(id), s)
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		_, err := t.Self.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TgBot) Initialize() {
	bot, err := tgbotapi.NewBotAPI(config.Config.TelegramBot.Token)
	if err != nil {
		println(err)
		return
	}
	bot.Debug = true

	t.Self = bot

}
