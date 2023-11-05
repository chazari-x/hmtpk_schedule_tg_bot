package time

import "time"

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
		return Now().Weekday().String()
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
		return Now().Weekday().ShortString()
	}
}

type Time time.Time

func Now() Time {
	return Time(time.Now())
}

func (t Time) Weekday() Weekday {
	a := time.Now().Weekday()
	if a == 0 {
		a = 7
	}

	return Weekday(a)
}
