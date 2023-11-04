package model

import (
	"strconv"
	"time"
)

type MySchCode int

func (c MySchCode) Code() string {
	if c > 6 || c < 0 {
		return "" + MyScheduleCode
	}
	return strconv.Itoa(int(c)) + MyScheduleCode
}

type MySchNextCode int

func (c MySchNextCode) Code() string {
	if c > 6 || c < 0 {
		return "" + MyScheduleNextCode
	}
	return strconv.Itoa(int(c)) + MyScheduleNextCode
}

type GroupSchCode int

func (c GroupSchCode) Code(group string) string {
	if c > 6 || c < 0 {
		return "" + MyScheduleCode + group
	}
	return strconv.Itoa(int(c)) + GroupScheduleCode + group
}

type GroupSchNextCode int

func (c GroupSchNextCode) Code(group string) string {
	if c > 6 || c < 0 {
		return "" + GroupScheduleNextCode + group
	}
	return strconv.Itoa(int(c)) + GroupScheduleNextCode + group
}

type TeacherSchCode int

func (c TeacherSchCode) Code(teacher string) string {
	if c > 6 || c < 0 {
		return "" + TeacherScheduleCode + teacher
	}
	return strconv.Itoa(int(c)) + TeacherScheduleCode + teacher
}

type TeacherSchNextCode int

func (c TeacherSchNextCode) Code(teacher string) string {
	if c > 6 || c < 0 {
		return "" + TeacherScheduleNextCode + teacher
	}
	return strconv.Itoa(int(c)) + TeacherScheduleNextCode + teacher
}

const (
	MyScheduleCode          = "M"
	GroupScheduleCode       = "G"
	TeacherScheduleCode     = "T"
	MyScheduleNextCode      = "MN"
	GroupScheduleNextCode   = "GN"
	TeacherScheduleNextCode = "TN"
)

type Button string

const (
	Start           string = "/start"
	Home            string = "Перейти в главное меню"
	MySchedule      string = "Мое расписание"
	OtherSchedule   string = "Другое расписание"
	GroupSchedule   string = /*"👩‍🎓" + */ "Группы"        /* + "👨‍🎓"*/
	TeacherSchedule string = /*"👩‍🏫" + */ "Преподаватели" /*+ "👨‍🏫"*/
	OtherButtons    string = "Другие кнопки"
	Support         string = "Служба поддержки"
	Settings        string = "Настройки"
	ChangeMyGroup   string = "Изменить мою группу"
	Statistics      string = "Статистика"
)

func (b Button) Value() string {
	switch b {
	case Button(Start):
		return `Дорогие пользователи,

Мы рады представить вам нашего нового Telegram бота! На данный момент, бот находится в стадии активной разработки, и мы работаем над расширением его функциональности.

Несмотря на то, что бот все еще находится в процессе разработки, вы уже можете начать использовать его и получать пользу от доступных функций. Мы стараемся делать его лучше с каждым обновлением.

`
	case Button(Home):
		return `Для получения расписания, у вас есть два варианта:

1. Для вашего собственного расписания нажмите кнопку "Мое расписание".

2. Чтобы получить расписание для другой группы или преподавателя, нажмите кнопку "Другое расписание".

Спасибо за использование нашего бота! Не стесняйтесь задавать вопросы, если у вас есть какие-либо. Вам всегда готовы помочь! ("Настройки" -> "Служба поддержки")

С наилучшими пожеланиями,
Команда разработчиков.`
	case Button(Support):
		return `Дорогие пользователи!

Если у вас есть какие-либо вопросы, предложения или вам требуется помощь, наша служба поддержки готова вам помочь. Вы можете связаться с нами по ссылке: <a href="%s">Служба поддержки бота</a>.

Мы всегда готовы ответить на ваши вопросы и рассмотреть ваши запросы. Не стесняйтесь обращаться!

Спасибо за вашу поддержку и использование нашего бота.

С наилучшими пожеланиями,
Команда разработчиков.`
	case Button(MySchedule):
		return "-"
	case Button(Settings):
		return "Ваши настройки."
	case Button(ChangeMyGroup):
		return "Пожалуйста, выберите или введите полный номер группы."
	case Button(OtherSchedule):
		return "Выберите расписание."
	case Button(GroupSchedule):
		return "Пожалуйста, выберите или введите полный номер группы."
	case Button(TeacherSchedule):
		return "Пожалуйста, выберите или введите ФИО преподавателя."
	case Button(OtherButtons):
		return "Показаны остальные кнопки."
	case Button(Statistics):
		return "Статистика использования бота:\n\nза день: %d пользователя\n\nза месяц: %d пользователя"
	default:
		return "-"
	}
}

type Weekday int

func (d Weekday) String() string {
	switch d {
	case 1:
		return "Понедельник"
	case 2:
		return "Вторник"
	case 3:
		return "Среда"
	case 4:
		return "Четверг"
	case 5:
		return "Пятница"
	case 6:
		return "Суббота"
	case 0:
		return "Воскресенье"
	default:
		return Weekday(time.Now().Weekday()).String()
	}
}

func (d Weekday) ShortString() string {
	switch d {
	case 1:
		return "ПН"
	case 2:
		return "ВТ"
	case 3:
		return "СР"
	case 4:
		return "ЧТ"
	case 5:
		return "ПТ"
	case 6:
		return "СБ"
	case 0:
		return "ВС"
	default:
		return Weekday(time.Now().Weekday()).ShortString()
	}
}
