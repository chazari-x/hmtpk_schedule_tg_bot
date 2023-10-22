package telegram

import (
	"github.com/chazari-x/hmtpk_schedule/config"
	"github.com/chazari-x/hmtpk_schedule/domain/telegram/logger"
	"github.com/chazari-x/hmtpk_schedule/domain/telegram/logic"
	"github.com/chazari-x/hmtpk_schedule/redis"
	"github.com/chazari-x/hmtpk_schedule/schedule"
	"github.com/chazari-x/hmtpk_schedule/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Start(cfg *config.Telegram, redis *redis.Redis, schedule *schedule.Schedule, storage *storage.Storage) error {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return err
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Error(err)
		return err
	}

	newLogger := logger.NewLogger(cfg, bot)
	newLogic := logic.NewLogic(cfg, schedule, bot, newLogger, storage, redis)

	go func() {
		for update := range updates {
			if update.Message == nil {
				if update.CallbackQuery != nil {
					newLogic.UpdateMessage(update.CallbackQuery)
				}
				continue
			}

			log.Infof("[%s: %d - %d] %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Chat.ID, update.Message.Text)

			if err := storage.InsertChat(int(update.Message.Chat.ID)); err != nil {
				log.Error(err)
			}

			newLogic.SendAnswer(update.Message)
		}
	}()

	log.Infof("Press CTRL-C to exit. Авторизован как %s.", bot.Self.UserName)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return nil
}