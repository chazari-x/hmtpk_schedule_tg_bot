package logics

import (
	"fmt"
	"github.com/chazari-x/hmtpk_schedule/config"
	"github.com/chazari-x/hmtpk_schedule/schedule"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

type Logic struct {
	cfg      *config.Telegram
	schedule *schedule.Schedule
}

func NewLogic(cfg *config.Telegram, schedule *schedule.Schedule) *Logic {
	return &Logic{cfg, schedule}
}

const (
	Support     = "Служба поддержки"
	GotoSupport = "Связаться со службой поддержки"
	SupportText = `Дорогие пользователи!

Если у вас есть какие-либо вопросы, предложения или вам требуется помощь, наша служба поддержки готова вам помочь. Вы можете связаться с нами по ссылке: <a href="%s">Служба поддержки бота</a>.

Мы всегда готовы ответить на ваши вопросы и рассмотреть ваши запросы. Не стесняйтесь обращаться!

Спасибо за вашу поддержку и использование нашего бота.

С наилучшими пожеланиями,
Команда разработчиков`

	GroupScheduleIcoOne = "👩‍🎓"
	GroupScheduleIcoTwo = "👨‍🎓"
	GroupSchedule       = GroupScheduleIcoOne + " Расписание группы " + GroupScheduleIcoTwo
	GroupScheduleText   = "Пожалуйста, выберите или введите полный номер группы."

	TeacherScheduleIcoOne = "👩‍🏫"
	TeacherScheduleIcoTwo = "👨‍🏫"
	TeacherSchedule       = TeacherScheduleIcoOne + " Расписание преподавателя " + TeacherScheduleIcoTwo
	TeacherScheduleText   = "Пожалуйста, выберите или введите ФИО преподавателя."

	GotoHome = "Перейти в главное меню"
	HomeText = `Дорогие пользователи!

Этот бот в настоящее время находится в стадии активной разработки. Мы работаем над его улучшением и добавлением новых функций, чтобы предоставить вам наилучший опыт.

Пожалуйста, будьте терпеливы и следите за нашими обновлениями. В скором времени бот будет работать в полном объеме и предоставлять вам больше возможностей.

Спасибо за ваше понимание и интерес к нашему проекту!

С наилучшими пожеланиями,
Команда разработчиков :3`
)

func (l *Logic) GetMessage(message *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	switch message.Text {
	case GotoHome:
		msg = tgbotapi.NewMessage(message.Chat.ID, HomeText)
	case GroupSchedule:
		msg = tgbotapi.NewMessage(message.Chat.ID, GroupScheduleText)
	case TeacherSchedule:
		msg = tgbotapi.NewMessage(message.Chat.ID, TeacherScheduleText)
	case Support:
		msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(SupportText, l.cfg.Support.Href))
	default:
		switch {
		case strings.Contains(message.Text, GroupScheduleIcoOne) && strings.Contains(message.Text, GroupScheduleIcoTwo):
			msg = tgbotapi.NewMessage(message.Chat.ID, "Вы выбрали группу "+message.Text)
		case strings.Contains(message.Text, TeacherScheduleIcoOne) && strings.Contains(message.Text, TeacherScheduleIcoTwo):
			msg = tgbotapi.NewMessage(message.Chat.ID, "Вы выбрали преподавателя "+message.Text)
		default:
			msg = tgbotapi.NewMessage(message.Chat.ID, HomeText)
		}
	}

	msg.ReplyMarkup = l.getButtons(message.Text)
	msg.ParseMode = "html"

	return msg, nil
}

func (l *Logic) getButtons(list string) tgbotapi.ReplyKeyboardMarkup {
	var replyMarkup tgbotapi.ReplyKeyboardMarkup
	switch list {
	case GroupSchedule:
		var keyboard [][]tgbotapi.KeyboardButton
		keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(GotoHome)))

		for _, group := range l.schedule.GetGroups() {
			group = fmt.Sprintf("%s %s %s", GroupScheduleIcoOne, group, GroupScheduleIcoTwo)
			keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(group)))
		}

		replyMarkup = tgbotapi.NewReplyKeyboard(keyboard...)
	case TeacherSchedule:
		var keyboard [][]tgbotapi.KeyboardButton
		keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(GotoHome)))

		for _, teacher := range l.schedule.GetTeachers() {
			teacher = fmt.Sprintf("%s %s %s", TeacherScheduleIcoOne, teacher, TeacherScheduleIcoTwo)
			keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(teacher)))
		}

		replyMarkup = tgbotapi.NewReplyKeyboard(keyboard...)
	default:
		replyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(GroupSchedule),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(TeacherSchedule),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(Support),
			),
		)
	}

	replyMarkup.OneTimeKeyboard = true // Скрыть клавиатуру после использования
	replyMarkup.ResizeKeyboard = true

	return replyMarkup
}
