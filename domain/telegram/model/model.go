package model

import (
	"fmt"
	"strconv"
	"time"
)

type MySchCode int

func (c MySchCode) Code() string {
	if c > 7 || c < 1 {
		return "" + MyScheduleCode
	}
	return strconv.Itoa(int(c)) + MyScheduleCode
}

type MySchNextCode int

func (c MySchNextCode) Code() string {
	if c > 7 || c < 1 {
		return "" + MyScheduleNextCode
	}
	return strconv.Itoa(int(c)) + MyScheduleNextCode
}

type GroupSchCode int

func (c GroupSchCode) Code(group string) string {
	if c > 7 || c < 1 {
		return "" + MyScheduleCode + group
	}
	return strconv.Itoa(int(c)) + GroupScheduleCode + group
}

type GroupSchNextCode int

func (c GroupSchNextCode) Code(group string) string {
	if c > 7 || c < 1 {
		return "" + GroupScheduleNextCode + group
	}
	return strconv.Itoa(int(c)) + GroupScheduleNextCode + group
}

type TeacherSchCode int

func (c TeacherSchCode) Code(teacher string) string {
	if c > 7 || c < 1 {
		return "" + TeacherScheduleCode + teacher
	}
	return strconv.Itoa(int(c)) + TeacherScheduleCode + teacher
}

type TeacherSchNextCode int

func (c TeacherSchNextCode) Code(teacher string) string {
	if c > 7 || c < 1 {
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
	Start           Button = "/start"
	Home            Button = "Перейти в главное меню"
	MySchedule      Button = "Мое расписание"
	OtherSchedule   Button = "Другое расписание"
	GroupSchedule   Button = "Группы"
	TeacherSchedule Button = "Преподаватели"
	OtherButtons    Button = "Другие кнопки"
	Support         Button = "Служба поддержки"
	Settings        Button = "Настройки"
	ChangeMyGroup   Button = "Изменить мою группу"
	Statistics      Button = "Статистика"
	Polity          Button = "Политика использования"
)

func (b Button) String() string {
	return string(b)
}

func (b Button) Value() string {
	switch b {
	case Start:
		return fmt.Sprintf(`Дорогие пользователи,

Мы рады представить вам нашего нового Telegram бота! На данный момент, бот находится в стадии активной разработки, и мы работаем над расширением его функциональности.

Несмотря на то, что бот все еще находится в процессе разработки, вы уже можете начать использовать его и получать пользу от доступных функций. Мы стараемся делать его лучше с каждым обновлением.

Используя бот "ХМТПК - РАСПИСАНИЕ", пользователи соглашаются с настоящей политикой использования (%s). Разработчик оставляет за собой право внесения изменений в политику использования без предварительного уведомления.

`, Polity.Cmd())
	case Home:
		return fmt.Sprintf(`Для получения расписания, у вас есть два варианта:

1. Для вашего собственного расписания нажмите кнопку "%s" (%s).

2. Чтобы получить расписание для другой группы (%s) или преподавателя (%s), нажмите кнопку "%s" (%s).

Спасибо за использование нашего бота! Не стесняйтесь задавать вопросы, если у вас есть какие-либо. Вам всегда готовы помочь! (%s)

С наилучшими пожеланиями,
Команда разработчиков.`, MySchedule, MySchedule.Cmd(), GroupSchedule.Cmd(), TeacherSchedule.Cmd(), OtherSchedule, OtherSchedule.Cmd(), Support.Cmd())
	case Support:
		return `Дорогие пользователи!

Если у вас есть какие-либо вопросы, предложения или вам требуется помощь, наша служба поддержки готова вам помочь. Вы можете связаться с нами по ссылке: <a href="%s">Служба поддержки бота</a>.

Мы всегда готовы ответить на ваши вопросы и рассмотреть ваши запросы. Не стесняйтесь обращаться!

Спасибо за вашу поддержку и использование нашего бота.

С наилучшими пожеланиями,
Команда разработчиков.`
	case MySchedule:
		return "-"
	case Settings:
		return fmt.Sprintf(`Изменение настроек:

1. Изменить мою группу (%s);

Для возврата на главную страницу нажмите кнопку "Перейти в главное меню" (%s).`, ChangeMyGroup.Cmd(), Home.Cmd())
	case ChangeMyGroup:
		return fmt.Sprintf(`Пожалуйста, выберите или введите полный номер группы.

Для возврата на главную страницу нажмите кнопку "Перейти в главное меню" (%s).`, Home.Cmd())
	case OtherSchedule:
		return fmt.Sprintf(`Выберите расписание для группы (%s) или для преподавателя (%s).

Для возврата на главную страницу нажмите кнопку "Перейти в главное меню" (%s).`, GroupSchedule.Cmd(), TeacherSchedule.Cmd(), Home.Cmd())
	case GroupSchedule:
		return fmt.Sprintf(`Пожалуйста, выберите или введите полный номер группы.

Для возврата на главную страницу нажмите кнопку "Перейти в главное меню" (%s).`, Home.Cmd())
	case TeacherSchedule:
		return fmt.Sprintf(`Пожалуйста, выберите или введите ФИО преподавателя.

Для возврата на главную страницу нажмите кнопку "Перейти в главное меню" (%s).`, Home.Cmd())
	case OtherButtons:
		return fmt.Sprintf(`Показаны остальные кнопки:

1. Политика использования (%s);

2. Служба поддержки (%s);

3. Настройки (%s);

4. Статистика (%s).

Для возврата на главную страницу нажмите кнопку "Перейти в главное меню" (%s).`, Polity.Cmd(), Support.Cmd(), Settings.Cmd(), Statistics.Cmd(), Home.Cmd())
	case Statistics:
		return "Статистика использования бота:\n\nза день: %d пользователя\n\nза месяц: %d пользователя"
	case Polity:
		return `<b>Политика использования бота "ХМТПК - РАСПИСАНИЕ"</b>

<b>Последнее обновление:</b> 04.11.2023

1. <b>Цель бота:</b> Бот "ХМТПК - РАСПИСАНИЕ" создан с целью предоставления расписания для удобства пользователей. Бот предоставляет доступ к информации о расписании, но не имеет непосредственного отношения к составлению расписания. Разработчик бота не несет ответственности за точность, актуальность или полноту предоставляемой информации.

2. <b>Источник информации:</b> Информация о расписании предоставляется на основе доступных данных из открытых источников, и она может быть подвержена изменениям без предварительного уведомления. Разработчик бота не контролирует и не влияет на источники информации о расписании.

3. <b>Авторские права:</b> Вся информация, предоставляемая ботом "ХМТПК - РАСПИСАНИЕ", охраняется авторскими правами и/или правами интеллектуальной собственности соответствующих организаций или лиц. Использование этой информации вне целей ознакомления может потребовать разрешения правообладателей.

4. <b>Точность информации:</b> Разработчик бота не гарантирует точность, актуальность или полноту предоставляемой информации о расписании. Пользователи обязаны проверять информацию о расписании у официальных источников или организаций, ответственных за составление расписания.

5. <b>Политика конфиденциальности:</b> Бот "ХМТПК - РАСПИСАНИЕ" может собирать данные об использовании бота, такие как действия пользователя и информацию о номере группы пользователя для целей улучшения функциональности и предоставления более релевантной информации. Разработчик бота обязуется соблюдать конфиденциальность данных и не передавать их третьим лицам.

6. <b>Обратная связь:</b> Пользователи могут обращаться к <a href="%s">разработчику бота</a> для предоставления обратной связи, сообщения об ошибках и предложений по улучшению функциональности. Ваши замечания и предложения всегда приветствуются.

7. <b>Отказ от ответственности:</b> Разработчик бота "ХМТПК - РАСПИСАНИЕ" отказывается от какой-либо ответственности за потерю, ущерб или неудовлетворение, связанные с использованием бота или предоставляемой им информацией о расписании.

Используя бот "ХМТПК - РАСПИСАНИЕ", пользователи соглашаются с настоящей политикой использования. Разработчик оставляет за собой право внесения изменений в политику использования без предварительного уведомления.`
	default:
		return "-"
	}
}

func (b Button) Cmd() string {
	switch b {
	case Home:
		return "/home"
	case Support:
		return "/support"
	case MySchedule:
		return "/myschedule"
	case Settings:
		return "/settings"
	case ChangeMyGroup:
		return "/changemygroup"
	case OtherSchedule:
		return "/otherschedule"
	case GroupSchedule:
		return "/group"
	case TeacherSchedule:
		return "/teacher"
	case OtherButtons:
		return "/other"
	case Statistics:
		return "/status"
	case Polity:
		return "/polity"
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
	case 7:
		return "Воскресенье"
	default:
		return Weekday(NowWeekday()).String()
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
	case 7:
		return "ВС"
	default:
		return Weekday(NowWeekday()).ShortString()
	}
}

func NowWeekday() int {
	a := time.Now().Weekday()
	if a == 0 {
		a = 7
	}

	return int(a)
}
