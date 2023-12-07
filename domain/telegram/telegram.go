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
		return err
	}

	newLogger := logger.NewLogger(cfg, bot)
	newLogic := logic.NewLogic(cfg, schedule, bot, newLogger, storage, redis)

	go func() {
		for update := range updates {
			if update.Message == nil {
				if update.CallbackQuery != nil {

					log.Infof("[%s: %d - %d] %s", update.CallbackQuery.From.UserName, update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)

					if err := storage.InsertChat(update.CallbackQuery.From.ID); err != nil {
						log.Error(err)

						if db, err := storage.Ping(); err != nil {
							log.Errorln(err)
						} else {
							storage.DB = db
						}
					}

					update.CallbackQuery.Message.From.ID = update.CallbackQuery.From.ID

					newLogic.UpdateMessage(update.CallbackQuery)
				}
				continue
			}

			if update.Message.Text != "" {
				log.Infof("[%s: %d - %d] %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Chat.ID, update.Message.Text)

				if err := storage.InsertChat(update.Message.From.ID); err != nil {
					log.Error(err)

					if db, err := storage.Ping(); err != nil {
						log.Errorln(err)
					} else {
						storage.DB = db
					}
				}

				newLogic.SendAnswer(update.Message)
			}
		}
	}()

	log.Infof("Press CTRL-C to exit. Авторизован как %s.", bot.Self.UserName)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return nil
}
