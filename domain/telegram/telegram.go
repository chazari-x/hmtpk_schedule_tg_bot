package telegram

import (
	"github.com/chazari-x/hmtpk_schedule/config"
	"github.com/chazari-x/hmtpk_schedule/domain/telegram/logger"
	"github.com/chazari-x/hmtpk_schedule/domain/telegram/logics"
	"github.com/chazari-x/hmtpk_schedule/redis/redis"
	"github.com/chazari-x/hmtpk_schedule/schedule"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func Start(cfg *config.Telegram, _ *redis.Redis, schedule *schedule.Schedule) error {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return err
	}

	//bot.Debug = true
	log.Infof("Авторизован как %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Error(err)
		return err
	}

	l := logger.NewLogger(cfg, bot)
	logic := logics.NewLogic(cfg, schedule)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Tracef("[%s - %d] %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)

		answer, err := logic.GetMessage(update.Message)
		if err != nil {
			log.Error(err)
			l.Error(err)
		}

		_, err = bot.Send(answer)
		if err != nil {
			log.Error(err)
			l.Error(err)
		}
	}

	return nil
}
