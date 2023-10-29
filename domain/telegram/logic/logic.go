package logic

import (
	"fmt"
	"github.com/chazari-x/hmtpk_schedule/config"
	"github.com/chazari-x/hmtpk_schedule/domain/telegram/logger"
	. "github.com/chazari-x/hmtpk_schedule/domain/telegram/model"
	"github.com/chazari-x/hmtpk_schedule/redis"
	"github.com/chazari-x/hmtpk_schedule/schedule"
	"github.com/chazari-x/hmtpk_schedule/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type Logic struct {
	cfg      *config.Telegram
	schedule *schedule.Schedule
	bot      *tgbotapi.BotAPI
	logger   *logger.Logger
	storage  *storage.Storage
	redis    *redis.Redis
}

func NewLogic(cfg *config.Telegram, schedule *schedule.Schedule, bot *tgbotapi.BotAPI, logger *logger.Logger, storage *storage.Storage, redis *redis.Redis) *Logic {
	return &Logic{cfg, schedule, bot, logger, storage, redis}
}

func (l *Logic) UpdateMessage(callbackQuery *tgbotapi.CallbackQuery) {
	log.Trace(callbackQuery.Data)
	var day string
	var sch tgbotapi.MessageConfig
	var buttons tgbotapi.InlineKeyboardMarkup
	var nextWeek = time.Now().AddDate(0, 0, 7).Format("02.01.2006")
	_, week := time.Now().AddDate(0, 0, 7).ISOWeek()
	if strings.Contains(callbackQuery.Data, TeacherScheduleNextCode) {
		data := strings.Split(callbackQuery.Data, TeacherScheduleNextCode)
		day = data[0]
		callbackQuery.Message.Text = data[1]
		buttons = l.getInlineKeyboard(TeacherScheduleNextCode, day, callbackQuery.Message.Text)
		sch = l.getTeacherSchedule(callbackQuery.Message, day, nextWeek, week)
	} else if strings.Contains(callbackQuery.Data, TeacherScheduleCode) {
		data := strings.Split(callbackQuery.Data, TeacherScheduleCode)
		day = data[0]
		callbackQuery.Message.Text = data[1]
		buttons = l.getInlineKeyboard(TeacherScheduleCode, day, callbackQuery.Message.Text)
		sch = l.getTeacherSchedule(callbackQuery.Message, day, "", 0)
	} else if strings.Contains(callbackQuery.Data, GroupScheduleNextCode) {
		data := strings.Split(callbackQuery.Data, GroupScheduleNextCode)
		day = data[0]
		callbackQuery.Message.Text = data[1]
		buttons = l.getInlineKeyboard(GroupScheduleNextCode, day, callbackQuery.Message.Text)
		sch = l.getGroupSchedule(callbackQuery.Message, day, nextWeek, week)
	} else if strings.Contains(callbackQuery.Data, GroupScheduleCode) {
		data := strings.Split(callbackQuery.Data, GroupScheduleCode)
		day = data[0]
		callbackQuery.Message.Text = data[1]
		buttons = l.getInlineKeyboard(GroupScheduleCode, day, callbackQuery.Message.Text)
		sch = l.getGroupSchedule(callbackQuery.Message, day, "", 0)
	} else if strings.Contains(callbackQuery.Data, MyScheduleNextCode) {
		day = strings.ReplaceAll(callbackQuery.Data, MyScheduleNextCode, "")
		buttons = l.getInlineKeyboard(MyScheduleNextCode, day, "")
		sch = l.getMySchedule(callbackQuery.Message, day, nextWeek, week)
	} else if strings.Contains(callbackQuery.Data, MyScheduleCode) {
		day = strings.ReplaceAll(callbackQuery.Data, MyScheduleCode, "")
		buttons = l.getInlineKeyboard(MyScheduleCode, day, "")
		sch = l.getMySchedule(callbackQuery.Message, day, "", 0)
	}

	msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, sch.Text)
	msg.ReplyMarkup = &buttons
	msg.ParseMode = "html"
	log.Info(len(msg.Text))
	if _, err := l.bot.Send(msg); err != nil {
		if !strings.Contains(err.Error(), "message is not modified") {
			log.Error(err)
			callbackResponse := tgbotapi.NewCallback(callbackQuery.ID, "Произошла ошибка")
			if _, err := l.bot.AnswerCallbackQuery(callbackResponse); err != nil {
				log.Error(err)
			}
			return
		}
		callbackResponse := tgbotapi.NewCallback(callbackQuery.ID, "Расписание не изменено")
		if _, err := l.bot.AnswerCallbackQuery(callbackResponse); err != nil {
			log.Error(err)
		}
		return
	}

	dayNum, err := strconv.Atoi(day)
	if err != nil {
		log.Error(err)
		return
	}
	callbackResponse := tgbotapi.NewCallback(callbackQuery.ID, fmt.Sprintf("Показано расписание на %s", Weekday(dayNum).String()))
	if _, err := l.bot.AnswerCallbackQuery(callbackResponse); err != nil {
		log.Error(err)
	}
}

func (l *Logic) SendAnswer(message *tgbotapi.Message) {
	var msg tgbotapi.MessageConfig

	get, e := l.redis.Get(fmt.Sprintf("chat-%d", message.From.ID))
	if e != nil {
		if !strings.Contains(e.Error(), "redis: nil") {
			log.Errorln(e)
		}
	}

	if get != "" {
		log.Traceln(get)
	}

	var buttons = get
	var id string
	switch get {
	case GroupSchedule:
		switch message.Text {
		case Home:
			buttons = message.Text
			get = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(Home).Value())
		default:
			if strings.Contains(message.Text, GroupSchedule) {
				buttons = message.Text
				msg = tgbotapi.NewMessage(message.Chat.ID, Button(GroupSchedule).Value())
				break
			}
			id = message.Text
			buttons = GroupSchCode(7).Code(id)
			msg = l.getGroupSchedule(message, "", "", 0)
		}
	case TeacherSchedule:
		switch message.Text {
		case Home:
			buttons = message.Text
			get = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(Home).Value())
		default:
			if strings.Contains(message.Text, TeacherSchedule) {
				buttons = message.Text
				msg = tgbotapi.NewMessage(message.Chat.ID, Button(TeacherSchedule).Value())
				break
			}
			id = message.Text
			buttons = TeacherSchCode(7).Code(id)
			msg = l.getTeacherSchedule(message, "", "", 0)
		}
	case ChangeMyGroup:
		switch message.Text {
		case Home:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(Home).Value())
		default:
			if message.Text != "" {
				if strings.Contains(message.Text, ChangeMyGroup) {
					buttons = message.Text
					msg = tgbotapi.NewMessage(message.Chat.ID, Button(ChangeMyGroup).Value())
					break
				}

				group := l.schedule.GetGroup(message.Text)
				if group != "" {
					err := l.storage.ChangeGroupID(int(message.From.ID), group)
					if err == nil {
						buttons = Home
						if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
							log.Errorln(err)
						}
						msg = tgbotapi.NewMessage(message.Chat.ID, "Вы изменили свою группу.")
						break
					}

					log.Errorln(err)
				}
			}

			msg = tgbotapi.NewMessage(message.Chat.ID, "Введена неверная группа.")
		}
	default:
		switch message.Text {
		case MySchedule:
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			buttons = MySchCode(7).Code()
			msg = l.getMySchedule(message, "", "", 0)
		case OtherSchedule:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(OtherSchedule).Value())
		case GroupSchedule:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), GroupSchedule); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(GroupSchedule).Value())
		case TeacherSchedule:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), TeacherSchedule); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(TeacherSchedule).Value())
		case Support:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(Button(Support).Value(), l.cfg.Support.Href))
		case Settings:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(Settings).Value())
		case ChangeMyGroup:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ChangeMyGroup); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(ChangeMyGroup).Value())
		case OtherButtons:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(OtherButtons).Value())
		case Statistics:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			day, month, err := l.storage.GetActiveChats()
			if err != nil {
				log.Errorln(err)
				msg = tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении статистики.")
				break
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(Button(Statistics).Value(), day, month))
		default:
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(Home).Value())
		}
	}

	_, _ = l.bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
		ChatID:    message.Chat.ID,
		MessageID: message.MessageID,
	})

	msg.ReplyMarkup = l.getKeyboard(buttons, id)
	msg.ParseMode = "html"

	_, err := l.bot.Send(msg)
	if err != nil {
		log.Errorln(err)
	}

	//if get != TeacherSchedule && get != GroupSchedule && message.Text != MySchedule {
	//	defer func(m tgbotapi.Message) {
	//		time.Sleep(time.Second * 30)
	//		_, _ = l.bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
	//			ChatID:    m.Chat.ID,
	//			MessageID: m.MessageID,
	//		})
	//	}(m)
	//}
}
