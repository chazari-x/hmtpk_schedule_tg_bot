package model

import "time"

type MySchCode int

func (c MySchCode) Code() string {
	switch c {
	case 0:
		return "0" + MyScheduleCode
	case 1:
		return "1" + MyScheduleCode
	case 2:
		return "2" + MyScheduleCode
	case 3:
		return "3" + MyScheduleCode
	case 4:
		return "4" + MyScheduleCode
	case 5:
		return "5" + MyScheduleCode
	case 6:
		return "6" + MyScheduleCode
	default:
		return "7" + MyScheduleCode
	}
}

type GroupSchCode int

func (c GroupSchCode) Code(group string) string {
	switch c {
	case 0:
		return "0" + GroupScheduleCode + group
	case 1:
		return "1" + GroupScheduleCode + group
	case 2:
		return "2" + GroupScheduleCode + group
	case 3:
		return "3" + GroupScheduleCode + group
	case 4:
		return "4" + GroupScheduleCode + group
	case 5:
		return "5" + GroupScheduleCode + group
	case 6:
		return "6" + GroupScheduleCode + group
	default:
		return "7" + GroupScheduleCode + group
	}
}

type TeacherSchCode int

func (c TeacherSchCode) Code(teacher string) string {
	switch c {
	case 0:
		return "0" + TeacherScheduleCode + teacher
	case 1:
		return "1" + TeacherScheduleCode + teacher
	case 2:
		return "2" + TeacherScheduleCode + teacher
	case 3:
		return "3" + TeacherScheduleCode + teacher
	case 4:
		return "4" + TeacherScheduleCode + teacher
	case 5:
		return "5" + TeacherScheduleCode + teacher
	case 6:
		return "6" + TeacherScheduleCode + teacher
	default:
		return "7" + TeacherScheduleCode + teacher
	}
}

const (
	MyScheduleCode      = "M"
	GroupScheduleCode   = "G"
	TeacherScheduleCode = "T"
)

type Button string

const (
	Home            string = "Перейти в главное меню"
	Support         string = "Служба поддержки"
	Settings        string = "Настройки"
	MySchedule      string = "Мое расписание"
	ChangeMyGroup   string = "Изменить мою группу"
	OtherSchedule   string = "Другое расписание"
	GroupSchedule   string = /*"👩‍🎓" + */ "Группы"        /* + "👨‍🎓"*/
	TeacherSchedule string = /*"👩‍🏫" + */ "Преподаватели" /*+ "👨‍🏫"*/
)

func (b Button) Value() string {
	switch b {
	case Button(Home):
		return `Дорогие пользователи!

Этот бот в настоящее время находится в стадии активной разработки. Мы работаем над его улучшением и добавлением новых функций, чтобы предоставить вам наилучший опыт.

Пожалуйста, будьте терпеливы и следите за нашими обновлениями. В скором времени бот будет работать в полном объеме и предоставлять вам больше возможностей.

Спасибо за ваше понимание и интерес к нашему проекту!

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
		return ""
	case Button(Settings):
		return "Ваши настройки."
	case Button(ChangeMyGroup):
		return "Пожалуйста, выберите или введите полный номер группы."
	case Button(OtherSchedule):
		return "Выберите расписание"
	case Button(GroupSchedule):
		return "Пожалуйста, выберите или введите полный номер группы."
	case Button(TeacherSchedule):
		return "Пожалуйста, выберите или введите ФИО преподавателя."
	default:
		return "-"
	}
}

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

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
