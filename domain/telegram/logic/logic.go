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
	if strings.Contains(callbackQuery.Data, TeacherScheduleCode) {
		data := strings.Split(callbackQuery.Data, TeacherScheduleCode)
		teacher := data[1]

		buttons = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(Weekday(1).String(), TeacherSchCode(1).Code(teacher)),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(2).String(), TeacherSchCode(2).Code(teacher)),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(3).String(), TeacherSchCode(3).Code(teacher)),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(4).String(), TeacherSchCode(4).Code(teacher)),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(5).String(), TeacherSchCode(5).Code(teacher)),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(6).String(), TeacherSchCode(6).Code(teacher)),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(0).String(), TeacherSchCode(0).Code(teacher)),
		))

		day = data[0]
		callbackQuery.Message.Text = teacher
		sch = l.getTeacherSchedule(callbackQuery.Message, day)
	} else if strings.Contains(callbackQuery.Data, GroupScheduleCode) {
		data := strings.Split(callbackQuery.Data, GroupScheduleCode)
		group := data[1]

		buttons = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(Weekday(1).String(), GroupSchCode(1).Code(group)),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(2).String(), GroupSchCode(2).Code(group)),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(3).String(), GroupSchCode(3).Code(group)),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(4).String(), GroupSchCode(4).Code(group)),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(5).String(), GroupSchCode(5).Code(group)),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(6).String(), GroupSchCode(6).Code(group)),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(0).String(), GroupSchCode(0).Code(group)),
		))

		day = data[0]
		callbackQuery.Message.Text = group
		sch = l.getGroupSchedule(callbackQuery.Message, day)
	} else if strings.Contains(callbackQuery.Data, MyScheduleCode) {
		buttons = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(Weekday(1).String(), MySchCode(1).Code()),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(2).String(), MySchCode(2).Code()),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(3).String(), MySchCode(3).Code()),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(4).String(), MySchCode(4).Code()),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(5).String(), MySchCode(5).Code()),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(6).String(), MySchCode(6).Code()),
			tgbotapi.NewInlineKeyboardButtonData(Weekday(0).String(), MySchCode(0).Code()),
		))

		day = strings.ReplaceAll(callbackQuery.Data, MyScheduleCode, "")
		sch = l.getMySchedule(callbackQuery.Message, day)
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

	get, e := l.redis.Get(fmt.Sprintf("chat-%d", message.Chat.ID))
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
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.Chat.ID), ""); err != nil {
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
			msg = l.getGroupSchedule(message, "")
		}
	case TeacherSchedule:
		switch message.Text {
		case Home:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.Chat.ID), ""); err != nil {
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
			msg = l.getTeacherSchedule(message, "")
		}
	case ChangeMyGroup:
		switch message.Text {
		case Home:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.Chat.ID), ""); err != nil {
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
					err := l.storage.ChangeGroupID(int(message.Chat.ID), group)
					if err == nil {
						buttons = Home
						if err := l.redis.Set(fmt.Sprintf("chat-%d", message.Chat.ID), ""); err != nil {
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
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.Chat.ID), ""); err != nil {
				log.Errorln(err)
			}
			buttons = MySchCode(7).Code()
			msg = l.getMySchedule(message, "")
		case OtherSchedule:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.Chat.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(OtherSchedule).Value())
		case GroupSchedule:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.Chat.ID), GroupSchedule); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(GroupSchedule).Value())
		case TeacherSchedule:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.Chat.ID), TeacherSchedule); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(TeacherSchedule).Value())
		case Support:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.Chat.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(Button(Support).Value(), l.cfg.Support.Href))
		case Settings:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.Chat.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(Settings).Value())
		case ChangeMyGroup:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.Chat.ID), ChangeMyGroup); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(ChangeMyGroup).Value())
		case OtherButtons:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.Chat.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(OtherButtons).Value())
		case Statistics:
			buttons = message.Text
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.Chat.ID), ""); err != nil {
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
			if err := l.redis.Set(fmt.Sprintf("chat-%d", message.Chat.ID), ""); err != nil {
				log.Errorln(err)
			}
			msg = tgbotapi.NewMessage(message.Chat.ID, Button(Home).Value())
		}
	}

	msg.ReplyMarkup = l.getButtons(buttons, id)
	msg.ParseMode = "html"

	_, err := l.bot.Send(msg)
	if err != nil {
		log.Errorln(err)
	}
}

func (l *Logic) getButtons(list, id string) interface{} {
	var replyMarkup tgbotapi.ReplyKeyboardMarkup
	switch list {
	case Settings:
		replyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(Home),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(ChangeMyGroup),
			),
		)
	case OtherSchedule:
		replyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(Home),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(GroupSchedule),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(TeacherSchedule),
			),
		)
	case MySchCode(0).Code(), MySchCode(1).Code(), MySchCode(2).Code(), MySchCode(3).Code(), MySchCode(4).Code(), MySchCode(5).Code(), MySchCode(6).Code(), MySchCode(7).Code():
		return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(Monday.String(), MySchCode(1).Code()),
			tgbotapi.NewInlineKeyboardButtonData(Tuesday.String(), MySchCode(2).Code()),
			tgbotapi.NewInlineKeyboardButtonData(Wednesday.String(), MySchCode(3).Code()),
			tgbotapi.NewInlineKeyboardButtonData(Thursday.String(), MySchCode(4).Code()),
			tgbotapi.NewInlineKeyboardButtonData(Friday.String(), MySchCode(5).Code()),
			tgbotapi.NewInlineKeyboardButtonData(Saturday.String(), MySchCode(6).Code()),
			tgbotapi.NewInlineKeyboardButtonData(Sunday.String(), MySchCode(0).Code()),
		))
	case GroupSchCode(0).Code(id), GroupSchCode(1).Code(id), GroupSchCode(2).Code(id), GroupSchCode(3).Code(id), GroupSchCode(4).Code(id), GroupSchCode(5).Code(id), GroupSchCode(6).Code(id), GroupSchCode(7).Code(id):
		return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(Monday.String(), GroupSchCode(1).Code(id)),
			tgbotapi.NewInlineKeyboardButtonData(Tuesday.String(), GroupSchCode(2).Code(id)),
			tgbotapi.NewInlineKeyboardButtonData(Wednesday.String(), GroupSchCode(3).Code(id)),
			tgbotapi.NewInlineKeyboardButtonData(Thursday.String(), GroupSchCode(4).Code(id)),
			tgbotapi.NewInlineKeyboardButtonData(Friday.String(), GroupSchCode(5).Code(id)),
			tgbotapi.NewInlineKeyboardButtonData(Saturday.String(), GroupSchCode(6).Code(id)),
			tgbotapi.NewInlineKeyboardButtonData(Sunday.String(), GroupSchCode(0).Code(id)),
		))
	case TeacherSchCode(0).Code(id), TeacherSchCode(1).Code(id), TeacherSchCode(2).Code(id), TeacherSchCode(3).Code(id), TeacherSchCode(4).Code(id), TeacherSchCode(5).Code(id), TeacherSchCode(6).Code(id), TeacherSchCode(7).Code(id):
		return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(Monday.String(), TeacherSchCode(1).Code(id)),
			tgbotapi.NewInlineKeyboardButtonData(Tuesday.String(), TeacherSchCode(2).Code(id)),
			tgbotapi.NewInlineKeyboardButtonData(Wednesday.String(), TeacherSchCode(3).Code(id)),
			tgbotapi.NewInlineKeyboardButtonData(Thursday.String(), TeacherSchCode(4).Code(id)),
			tgbotapi.NewInlineKeyboardButtonData(Friday.String(), TeacherSchCode(5).Code(id)),
			tgbotapi.NewInlineKeyboardButtonData(Saturday.String(), TeacherSchCode(6).Code(id)),
			tgbotapi.NewInlineKeyboardButtonData(Sunday.String(), TeacherSchCode(0).Code(id)),
		))
	case OtherButtons:
		replyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(Home),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(Support),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(Settings),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(Statistics),
			),
		)
	default:
		if strings.Contains(list, ChangeMyGroup) {
			numString := strings.ReplaceAll(list, ChangeMyGroup+" ", "")
			if numString == "" || numString == ChangeMyGroup {
				numString = "1"
			}
			num, err := strconv.Atoi(numString)
			if err != nil {
				log.Error(err)
			}

			var keyboard [][]tgbotapi.KeyboardButton
			keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(Home)))

			if num > 1 {
				keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(ChangeMyGroup+" "+strconv.Itoa(num-1))))
			}

			groups := l.schedule.GetGroups()
			for i := len(groups) - 50*(num-1) - 1; i >= 0 && i >= len(groups)-50*(num-1)-50 && i < len(groups); i-- {
				keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(groups[i])))
			}

			if num*50 < len(groups) {
				keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(ChangeMyGroup+" "+strconv.Itoa(num+1))))
			}

			replyMarkup = tgbotapi.NewReplyKeyboard(keyboard...)
		} else if strings.Contains(list, GroupSchedule) {
			numString := strings.ReplaceAll(list, GroupSchedule+" ", "")
			if numString == "" || numString == GroupSchedule {
				numString = "1"
			}
			num, err := strconv.Atoi(numString)
			if err != nil {
				log.Error(err)
			}

			var keyboard [][]tgbotapi.KeyboardButton
			keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(Home)))

			if num > 1 {
				keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(GroupSchedule+" "+strconv.Itoa(num-1))))
			}

			groups := l.schedule.GetGroups()
			for i := len(groups) - 50*(num-1) - 1; i >= 0 && i >= len(groups)-50*(num-1)-50 && i < len(groups); i-- {
				keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(groups[i])))
			}

			if num*50 < len(groups) {
				keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(GroupSchedule+" "+strconv.Itoa(num+1))))
			}

			replyMarkup = tgbotapi.NewReplyKeyboard(keyboard...)
		} else if strings.Contains(list, TeacherSchedule) {
			numString := strings.ReplaceAll(list, TeacherSchedule+" ", "")
			if numString == "" || numString == TeacherSchedule {
				numString = "1"
			}
			num, err := strconv.Atoi(numString)
			if err != nil {
				log.Error(err)
			}

			var keyboard [][]tgbotapi.KeyboardButton
			keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(Home)))

			if num > 1 {
				keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(TeacherSchedule+" "+strconv.Itoa(num-1))))
			}

			teachers := l.schedule.GetTeachers()
			for i := len(teachers) - 50*(num-1) - 1; i >= 0 && i >= len(teachers)-50*(num-1)-50 && i < len(teachers); i-- {
				keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(teachers[i])))
			}

			if num*50 < len(teachers) {
				keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(TeacherSchedule+" "+strconv.Itoa(num+1))))
			}

			replyMarkup = tgbotapi.NewReplyKeyboard(keyboard...)
		} else {
			replyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(MySchedule),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(OtherSchedule),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(OtherButtons),
				),
			)
		}
	}

	replyMarkup.ResizeKeyboard = true

	return replyMarkup
}

func getDay(day time.Weekday) string {
	return strings.ToLower(Weekday(day).String())
}

func (l *Logic) getMySchedule(message *tgbotapi.Message, date string) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	date = strings.ToLower(date)
	group, err := l.storage.SelectGroupID(int(message.From.ID))
	if err != nil {
		return tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при поиске вашей группы")
	} else if group == "0" || group == "" {
		return tgbotapi.NewMessage(message.Chat.ID, "У вас не выбрана группа")
	}
	if date == "" || date == "7" {
		date = getDay(time.Now().Weekday())
	} else {
		atoi, err := strconv.Atoi(date)
		if err != nil {
			log.Error(err)
			return tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении даты")
		}
		date = getDay(time.Weekday(atoi))
	}
	schs, err := l.schedule.GetScheduleByGroupID(group, time.Now().Format("02.01.2006"))
	if err != nil || len(schs) == 0 {
		if err != nil {
			log.Error(err)
		}
		return tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении расписания")
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, "")

	for _, sch := range schs {
		if strings.Contains(sch.Date, date) {
			msg.Text += fmt.Sprintf("Группа: %s\n%s", group, sch.Date)

			if len(sch.Lessons) == 0 {
				msg.Text += "\nРасписания нет (пар нет)"
				return msg
			}

			for i, lesson := range sch.Lessons {
				if i == 0 {
					msg.Text += fmt.Sprintf("\n\nУрок: <code>%s</code>", lesson.Num)
					msg.Text += fmt.Sprintf(" [<code>%s</code>]", lesson.Time)
					msg.Text += fmt.Sprintf("\n<code>%s</code>", lesson.Name)
					if lesson.Room != "СРС" {
						msg.Text += fmt.Sprintf("\n[<code>%s</code>]", lesson.Room)
					}
					if lesson.Teacher != "<>" {
						msg.Text += fmt.Sprintf(" <code>%s</code>", lesson.Teacher)
					}
				} else {
					if lesson.Num == sch.Lessons[i-1].Num {
						if lesson.Name != sch.Lessons[i-1].Name {
							msg.Text += fmt.Sprintf("\n<code>%s</code>", lesson.Name)
						}
					} else {
						msg.Text += fmt.Sprintf("\n\nУрок: <code>%s</code>", lesson.Num)
						msg.Text += fmt.Sprintf(" [<code>%s</code>]", lesson.Time)
						msg.Text += fmt.Sprintf("\n<code>%s</code>", lesson.Name)
					}
					if lesson.Room != "СРС" {
						msg.Text += fmt.Sprintf("\n[<code>%s</code>]", lesson.Room)
					}
					if lesson.Teacher != "<>" {
						msg.Text += fmt.Sprintf(" <code>%s</code>", lesson.Teacher)
					}
				}
			}
		}
	}

	if msg.Text == "" {
		msg.Text += "Произошла какая-то ошибка"
	}

	return msg
}

func (l *Logic) getGroupSchedule(message *tgbotapi.Message, date string) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	date = strings.ToLower(date)
	//group := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(message.Text, GroupScheduleIcoOne, ""), GroupScheduleIcoTwo, ""), " ", "")
	group := message.Text
	if date == "" || date == "7" {
		date = getDay(time.Now().Weekday())
	} else {
		atoi, err := strconv.Atoi(date)
		if err != nil {
			log.Error(err)
			return tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении даты")
		}
		date = getDay(time.Weekday(atoi))
	}
	schs, err := l.schedule.GetScheduleByGroupID(group, time.Now().Format("02.01.2006"))
	if err != nil || len(schs) == 0 {
		if err != nil {
			log.Error(err)
		}
		return tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении расписания")
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, "")

	for _, sch := range schs {
		if strings.Contains(sch.Date, date) {
			msg.Text += fmt.Sprintf("Группа: %s\n%s", group, sch.Date)

			if len(sch.Lessons) == 0 {
				msg.Text += "\nРасписания нет (пар нет)"
				return msg
			}

			for i, lesson := range sch.Lessons {
				if i == 0 {
					msg.Text += fmt.Sprintf("\n\nУрок: <code>%s</code>", lesson.Num)
					msg.Text += fmt.Sprintf(" [<code>%s</code>]", lesson.Time)
					msg.Text += fmt.Sprintf("\n<code>%s</code>", lesson.Name)
					if lesson.Room != "СРС" {
						msg.Text += fmt.Sprintf("\n[<code>%s</code>]", lesson.Room)
					}
					if lesson.Teacher != "<>" {
						msg.Text += fmt.Sprintf(" <code>%s</code>", lesson.Teacher)
					}
				} else {
					if lesson.Num == sch.Lessons[i-1].Num {
						if lesson.Name != sch.Lessons[i-1].Name {
							msg.Text += fmt.Sprintf("\n<code>%s</code>", lesson.Name)
						}
					} else {
						msg.Text += fmt.Sprintf("\n\nУрок: <code>%s</code>", lesson.Num)
						msg.Text += fmt.Sprintf(" [<code>%s</code>]", lesson.Time)
						msg.Text += fmt.Sprintf("\n<code>%s</code>", lesson.Name)
					}
					if lesson.Room != "СРС" {
						msg.Text += fmt.Sprintf("\n[<code>%s</code>]", lesson.Room)
					}
					if lesson.Teacher != "<>" {
						msg.Text += fmt.Sprintf(" <code>%s</code>", lesson.Teacher)
					}
				}
			}
		}
	}

	if msg.Text == "" {
		msg.Text += "Произошла какая-то ошибка"
	}

	return msg
}

func (l *Logic) getTeacherSchedule(message *tgbotapi.Message, date string) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	date = strings.ToLower(date)
	//teacher := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(message.Text, TeacherScheduleIcoOne+" ", ""), " "+TeacherScheduleIcoTwo, ""), " ", "+")
	teacher := strings.ReplaceAll(message.Text, " ", "+")
	if date == "" || date == "7" {
		date = getDay(time.Now().Weekday())
	} else {
		atoi, err := strconv.Atoi(date)
		if err != nil {
			log.Error(err)
			return tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении даты")
		}
		date = getDay(time.Weekday(atoi))
	}
	schs, err := l.schedule.GetScheduleByTeacher(teacher, time.Now().Format("02.01.2006"))
	if err != nil || len(schs) == 0 {
		if err != nil {
			log.Error(err)
		}
		return tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении расписания")
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, "")

	for _, sch := range schs {
		if strings.Contains(sch.Date, date) {
			msg.Text += fmt.Sprintf("Преподаватель: %s\n%s", teacher, sch.Date)

			if len(sch.Lessons) == 0 {
				msg.Text += "\nРасписания нет (пар нет)"
				return msg
			}

			for _, lesson := range sch.Lessons {
				msg.Text += fmt.Sprintf(`

Урок: <code>%s</code> [<code>%s</code>]
<code>%s</code>
[<code>%s</code>] <code>%s</code>`,
					lesson.Num,
					lesson.Time,
					lesson.Name,
					lesson.Room,
					lesson.Group)
			}
		}
	}

	if msg.Text == "" {
		msg.Text += "Произошла какая-то ошибка"
	}

	return msg
}
