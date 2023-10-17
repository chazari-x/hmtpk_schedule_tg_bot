package logger

import (
	"github.com/chazari-x/hmtpk_schedule/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type Logger struct {
	cfg *config.Telegram
	bot *tgbotapi.BotAPI
}

func NewLogger(cfg *config.Telegram, bot *tgbotapi.BotAPI) *Logger {
	return &Logger{cfg: cfg, bot: bot}
}

func (l *Logger) Message(message string) {
	msg := tgbotapi.NewMessage(l.cfg.Support.ID, message)
	_, err := l.bot.Send(msg)
	log.Error(err)
}

func (l *Logger) Error(error error) {
	msg := tgbotapi.NewMessage(l.cfg.Support.ID, error.Error())
	_, err := l.bot.Send(msg)
	log.Error(err)
}
