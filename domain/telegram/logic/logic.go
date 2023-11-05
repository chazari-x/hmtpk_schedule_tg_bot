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
		if group, err := l.storage.SelectGroupID(callbackQuery.Message.From.ID); err != nil {
			if db, err := l.storage.Ping(); err != nil {
				log.Errorln(err)
			} else {
				l.storage.DB = db
			}

			sch = tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Произошла ошибка при поиске вашей группы")
		} else if group == "0" || group == "" {
			sch = tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "У вас не выбрана группа")
		} else {
			sch = l.getMySchedule(callbackQuery.Message, day, nextWeek, group, week)
		}
	} else if strings.Contains(callbackQuery.Data, MyScheduleCode) {
		day = strings.ReplaceAll(callbackQuery.Data, MyScheduleCode, "")
		buttons = l.getInlineKeyboard(MyScheduleCode, day, "")
		if group, err := l.storage.SelectGroupID(callbackQuery.Message.From.ID); err != nil {
			if db, err := l.storage.Ping(); err != nil {
				log.Errorln(err)
			} else {
				l.storage.DB = db
			}

			sch = tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Произошла ошибка при поиске вашей группы")
		} else if group == "0" || group == "" {
			sch = tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "У вас не выбрана группа")
		} else {
			sch = l.getMySchedule(callbackQuery.Message, day, "", group, 0)
		}
	}

	msg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, sch.Text)
	msg.ReplyMarkup = &buttons
	msg.ParseMode = "html"
	log.Info(len(msg.Text))
	if _, err := l.bot.Send(msg); err != nil {
		if !strings.Contains(err.Error(), "message is not modified") {
			log.Error(fmt.Errorf("%s: %s", err, msg.Text))
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

	if strings.HasPrefix(message.Text, "/") {
		get = ""
	}

	if get != "" {
		log.Traceln(get)
	}

	var buttons = get
	var id string
	switch get {
	case GroupSchedule.String(), GroupSchedule.Cmd():
		switch message.Text {
		case Home.String(), Start.String(), Home.Cmd(), Start.Cmd():
			buttons = Home.String()
			get = Home.String()
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Home.Value())
		default:
			if strings.Contains(message.Text, GroupSchedule.String()) {
				buttons = message.Text
				msg = tgbotapi.NewMessage(message.Chat.ID, GroupSchedule.Value())
				break
			}
			id = message.Text
			buttons = GroupSchCode(0).Code(id)
			msg = l.getGroupSchedule(message, "", "", 0)
		}
	case TeacherSchedule.String(), TeacherSchedule.Cmd():
		switch message.Text {
		case Home.String(), Start.String(), Home.Cmd(), Start.Cmd():
			buttons = Home.String()
			get = Home.String()
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Home.Value())
		default:
			if strings.Contains(message.Text, TeacherSchedule.String()) {
				buttons = message.Text
				msg = tgbotapi.NewMessage(message.Chat.ID, TeacherSchedule.Value())
				break
			}
			id = message.Text
			buttons = TeacherSchCode(0).Code(id)
			msg = l.getTeacherSchedule(message, "", "", 0)
		}
	case ChangeMyGroup.String(), ChangeMyGroup.Cmd():
		switch message.Text {
		case Home.String(), Start.String(), Home.Cmd(), Start.Cmd():
			buttons = Home.String()
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Home.Value())
		default:
			if message.Text != "" {
				if strings.Contains(message.Text, ChangeMyGroup.String()) {
					buttons = message.Text
					msg = tgbotapi.NewMessage(message.Chat.ID, ChangeMyGroup.Value())
					break
				}

				group := l.schedule.GetGroup(message.Text)
				if group != "" {
					err := l.storage.ChangeGroupID(message.From.ID, group)
					if err == nil {
						buttons = Home.String()
						if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
							log.Errorln(err)
						}
						msg = tgbotapi.NewMessage(message.Chat.ID, "Вы изменили свою группу.")
						break
					}

					log.Errorln(err)

					if db, err := l.storage.Ping(); err != nil {
						log.Errorln(err)
					} else {
						l.storage.DB = db
					}
				}
			}

			msg = tgbotapi.NewMessage(message.Chat.ID, "Введена неверная группа.")
		}
	default:
		switch message.Text {
		case Start.String():
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("%s%s", Start.Value(), Home.Value()))
		case MySchedule.String(), MySchedule.Cmd():
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			buttons = MySchCode(0).Code()
			group, err := l.storage.SelectGroupID(message.From.ID)
			if err != nil {
				if db, err := l.storage.Ping(); err != nil {
					log.Errorln(err)
				} else {
					l.storage.DB = db
				}

				msg = tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при поиске вашей группы")
			} else if group == "0" || group == "" {
				buttons = ChangeMyGroup.String()
				if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ChangeMyGroup.String()); err != nil {
					log.Errorln(err)
				}
				msg = tgbotapi.NewMessage(message.Chat.ID, ChangeMyGroup.Value())
			} else {
				msg = l.getMySchedule(message, "", "", group, 0)
			}
		case OtherSchedule.String(), OtherSchedule.Cmd():
			buttons = OtherSchedule.String()
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, OtherSchedule.Value())
		case GroupSchedule.String(), GroupSchedule.Cmd():
			buttons = GroupSchedule.String()
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), GroupSchedule.String()); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, GroupSchedule.Value())
		case TeacherSchedule.String(), TeacherSchedule.Cmd():
			buttons = TeacherSchedule.String()
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), TeacherSchedule.String()); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, TeacherSchedule.Value())
		case Support.String(), Support.Cmd():
			buttons = Support.String()
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(Support.Value(), l.cfg.Support.Href))
		case Settings.String(), Settings.Cmd():
			buttons = Settings.String()
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Settings.Value())
		case ChangeMyGroup.String(), ChangeMyGroup.Cmd():
			buttons = ChangeMyGroup.String()
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ChangeMyGroup.String()); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, ChangeMyGroup.Value())
		case OtherButtons.String(), OtherButtons.Cmd():
			buttons = OtherButtons.String()
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, OtherButtons.Value())
		case Statistics.String(), Statistics.Cmd():
			buttons = Statistics.String()
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			day, month, err := l.storage.GetActiveChats()
			if err != nil {
				log.Errorln(err)
				msg = tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении статистики.")

				if db, err := l.storage.Ping(); err != nil {
					log.Errorln(err)
				} else {
					l.storage.DB = db
				}

				break
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(Statistics.Value(), day, month))
		case Polity.String(), Polity.Cmd():
			buttons = Polity.String()
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(Polity.Value(), l.cfg.Support.Href))
		default:
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.From.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Home.Value())
		}
	}

	//_, _ = l.bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
	//	ChatID:    message.Chat.ID,
	//	MessageID: message.MessageID,
	//})

	msg.ReplyMarkup = l.getKeyboard(buttons, id)
	msg.ParseMode = "html"

	_, err := l.bot.Send(msg)
	if err != nil {
		log.Errorln(err)
	}
}
