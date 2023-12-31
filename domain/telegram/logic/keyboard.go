package logic

import (
	"fmt"
	. "github.com/chazari-x/hmtpk_schedule/domain/telegram/model"
	time2 "github.com/chazari-x/hmtpk_schedule/time"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func (l *Logic) getInlineKeyboard(list, dayStr, value string) tgbotapi.InlineKeyboardMarkup {
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		log.Error(err)
		day = int(time2.Now().Weekday())
	}
	t := int(time2.Now().Weekday())
	var keyboard tgbotapi.InlineKeyboardMarkup
	switch list {
	case TeacherScheduleCode:
		buttons := [][]tgbotapi.InlineKeyboardButton{{}, {}, {}, {}}
		for i := 1; i <= 7 && i > 0; i++ {
			if i != day {
				buttonDay := fmt.Sprintf("%s: %s", time2.Weekday(i).ShortString(), time.Now().AddDate(0, 0, i-t).Format("02.01"))
				if len(buttons[0]) < 3 {
					buttons[0] = append(buttons[0], tgbotapi.NewInlineKeyboardButtonData(buttonDay, TeacherSchCode(i).Code(value)))
				} else {
					buttons[1] = append(buttons[1], tgbotapi.NewInlineKeyboardButtonData(buttonDay, TeacherSchCode(i).Code(value)))
				}
			}
		}

		buttonUpdate := fmt.Sprintf("Обновить - %s: %s", time2.Weekday(day).ShortString(), time.Now().AddDate(0, 0, day-t).Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonUpdate, TeacherSchCode(day).Code(value)))
		//if day != t {
		buttonToday := fmt.Sprintf("Сегодня - %s: %s", time2.Weekday(t).ShortString(), time.Now().Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonToday, TeacherSchCode(t).Code(value)))
		//}

		buttonNext := fmt.Sprintf("Следующая неделя: %s - %s", time.Now().AddDate(0, 0, 7+1-t).Format("02.01"), time.Now().AddDate(0, 0, 7+7-t).Format("02.01"))
		buttons[3] = append(buttons[3], tgbotapi.NewInlineKeyboardButtonData(buttonNext, TeacherSchNextCode(1).Code(value)))

		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(buttons[0]...),
			tgbotapi.NewInlineKeyboardRow(buttons[1]...),
			tgbotapi.NewInlineKeyboardRow(buttons[2]...),
			tgbotapi.NewInlineKeyboardRow(buttons[3]...),
		)
	case TeacherScheduleNextCode:
		buttons := [][]tgbotapi.InlineKeyboardButton{{}, {}, {}, {}}
		for i := 1; i <= 7 && i > 0; i++ {
			if i != day {
				buttonDay := fmt.Sprintf("%s: %s", time2.Weekday(i).ShortString(), time.Now().AddDate(0, 0, 7+i-t).Format("02.01"))
				if len(buttons[0]) < 3 {
					buttons[0] = append(buttons[0], tgbotapi.NewInlineKeyboardButtonData(buttonDay, TeacherSchNextCode(i).Code(value)))
				} else {
					buttons[1] = append(buttons[1], tgbotapi.NewInlineKeyboardButtonData(buttonDay, TeacherSchNextCode(i).Code(value)))
				}
			}
		}

		buttonUpdate := fmt.Sprintf("Обновить - %s: %s", time2.Weekday(day).ShortString(), time.Now().AddDate(0, 0, 7+day-t).Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonUpdate, TeacherSchNextCode(day).Code(value)))
		buttonToday := fmt.Sprintf("Сегодня - %s: %s", time2.Weekday(t).ShortString(), time.Now().Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonToday, TeacherSchCode(t).Code(value)))

		buttonPast := fmt.Sprintf("Текущая неделя: %s - %s", time.Now().AddDate(0, 0, 1-t).Format("02.01"), time.Now().AddDate(0, 0, 7-t).Format("02.01"))
		buttons[3] = append(buttons[3], tgbotapi.NewInlineKeyboardButtonData(buttonPast, TeacherSchCode(t).Code(value)))

		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(buttons[0]...),
			tgbotapi.NewInlineKeyboardRow(buttons[1]...),
			tgbotapi.NewInlineKeyboardRow(buttons[2]...),
			tgbotapi.NewInlineKeyboardRow(buttons[3]...),
		)
	case GroupScheduleCode:
		buttons := [][]tgbotapi.InlineKeyboardButton{{}, {}, {}, {}}
		for i := 1; i <= 7 && i > 0; i++ {
			if i != day {
				buttonDay := fmt.Sprintf("%s: %s", time2.Weekday(i).ShortString(), time.Now().AddDate(0, 0, i-t).Format("02.01"))
				if len(buttons[0]) < 3 {
					buttons[0] = append(buttons[0], tgbotapi.NewInlineKeyboardButtonData(buttonDay, GroupSchCode(i).Code(value)))
				} else {
					buttons[1] = append(buttons[1], tgbotapi.NewInlineKeyboardButtonData(buttonDay, GroupSchCode(i).Code(value)))
				}
			}
		}

		buttonUpdate := fmt.Sprintf("Обновить - %s: %s", time2.Weekday(day).ShortString(), time.Now().AddDate(0, 0, day-t).Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonUpdate, GroupSchCode(day).Code(value)))
		//if day != t {
		buttonToday := fmt.Sprintf("Сегодня - %s: %s", time2.Weekday(t).ShortString(), time.Now().Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonToday, GroupSchCode(t).Code(value)))
		//}

		buttonNext := fmt.Sprintf("Следующая неделя: %s - %s", time.Now().AddDate(0, 0, 7+1-t).Format("02.01"), time.Now().AddDate(0, 0, 7+7-t).Format("02.01"))
		buttons[3] = append(buttons[3], tgbotapi.NewInlineKeyboardButtonData(buttonNext, GroupSchNextCode(1).Code(value)))

		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(buttons[0]...),
			tgbotapi.NewInlineKeyboardRow(buttons[1]...),
			tgbotapi.NewInlineKeyboardRow(buttons[2]...),
			tgbotapi.NewInlineKeyboardRow(buttons[3]...),
		)
	case GroupScheduleNextCode:
		buttons := [][]tgbotapi.InlineKeyboardButton{{}, {}, {}, {}}
		for i := 1; i <= 7 && i > 0; i++ {
			if i != day {
				buttonDay := fmt.Sprintf("%s: %s", time2.Weekday(i).ShortString(), time.Now().AddDate(0, 0, 7+i-t).Format("02.01"))
				if len(buttons[0]) < 3 {
					buttons[0] = append(buttons[0], tgbotapi.NewInlineKeyboardButtonData(buttonDay, GroupSchNextCode(i).Code(value)))
				} else {
					buttons[1] = append(buttons[1], tgbotapi.NewInlineKeyboardButtonData(buttonDay, GroupSchNextCode(i).Code(value)))
				}
			}
		}

		buttonUpdate := fmt.Sprintf("Обновить - %s: %s", time2.Weekday(day).ShortString(), time.Now().AddDate(0, 0, 7+day-t).Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonUpdate, GroupSchNextCode(day).Code(value)))
		buttonToday := fmt.Sprintf("Сегодня - %s: %s", time2.Weekday(t).ShortString(), time.Now().Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonToday, GroupSchCode(t).Code(value)))

		buttonPast := fmt.Sprintf("Текущая неделя: %s - %s", time.Now().AddDate(0, 0, 1-t).Format("02.01"), time.Now().AddDate(0, 0, 7-t).Format("02.01"))
		buttons[3] = append(buttons[3], tgbotapi.NewInlineKeyboardButtonData(buttonPast, GroupSchCode(t).Code(value)))

		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(buttons[0]...),
			tgbotapi.NewInlineKeyboardRow(buttons[1]...),
			tgbotapi.NewInlineKeyboardRow(buttons[2]...),
			tgbotapi.NewInlineKeyboardRow(buttons[3]...),
		)
	case MyScheduleCode:
		buttons := [][]tgbotapi.InlineKeyboardButton{{}, {}, {}, {}}
		for i := 1; i <= 7 && i > 0; i++ {
			if i != day {
				buttonDay := fmt.Sprintf("%s: %s", time2.Weekday(i).ShortString(), time.Now().AddDate(0, 0, i-t).Format("02.01"))
				if len(buttons[0]) < 3 {
					buttons[0] = append(buttons[0], tgbotapi.NewInlineKeyboardButtonData(buttonDay, MySchCode(i).Code()))
				} else {
					buttons[1] = append(buttons[1], tgbotapi.NewInlineKeyboardButtonData(buttonDay, MySchCode(i).Code()))
				}
			}
		}

		buttonUpdate := fmt.Sprintf("Обновить - %s: %s", time2.Weekday(day).ShortString(), time.Now().AddDate(0, 0, day-t).Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonUpdate, MySchCode(day).Code()))
		//if day != t {
		buttonToday := fmt.Sprintf("Сегодня - %s: %s", time2.Weekday(t).ShortString(), time.Now().Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonToday, MySchCode(t).Code()))
		//}

		buttonNext := fmt.Sprintf("Следующая неделя: %s - %s", time.Now().AddDate(0, 0, 7+1-t).Format("02.01"), time.Now().AddDate(0, 0, 7+7-t).Format("02.01"))
		buttons[3] = append(buttons[3], tgbotapi.NewInlineKeyboardButtonData(buttonNext, MySchNextCode(1).Code()))

		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(buttons[0]...),
			tgbotapi.NewInlineKeyboardRow(buttons[1]...),
			tgbotapi.NewInlineKeyboardRow(buttons[2]...),
			tgbotapi.NewInlineKeyboardRow(buttons[3]...),
		)
	case MyScheduleNextCode:
		buttons := [][]tgbotapi.InlineKeyboardButton{{}, {}, {}, {}}
		for i := 1; i <= 7 && i > 0; i++ {
			if i != day {
				buttonDay := fmt.Sprintf("%s: %s", time2.Weekday(i).ShortString(), time.Now().AddDate(0, 0, 7+i-t).Format("02.01"))
				if len(buttons[0]) < 3 {
					buttons[0] = append(buttons[0], tgbotapi.NewInlineKeyboardButtonData(buttonDay, MySchNextCode(i).Code()))
				} else {
					buttons[1] = append(buttons[1], tgbotapi.NewInlineKeyboardButtonData(buttonDay, MySchNextCode(i).Code()))
				}
			}
		}

		buttonUpdate := fmt.Sprintf("Обновить - %s: %s", time2.Weekday(day).ShortString(), time.Now().AddDate(0, 0, 7+day-t).Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonUpdate, MySchNextCode(day).Code()))
		buttonToday := fmt.Sprintf("Сегодня - %s: %s", time2.Weekday(t).ShortString(), time.Now().Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonToday, MySchCode(t).Code()))

		buttonPast := fmt.Sprintf("Текущая неделя: %s - %s", time.Now().AddDate(0, 0, 1-t).Format("02.01"), time.Now().AddDate(0, 0, 7-t).Format("02.01"))
		buttons[3] = append(buttons[3], tgbotapi.NewInlineKeyboardButtonData(buttonPast, MySchCode(t).Code()))

		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(buttons[0]...),
			tgbotapi.NewInlineKeyboardRow(buttons[1]...),
			tgbotapi.NewInlineKeyboardRow(buttons[2]...),
			tgbotapi.NewInlineKeyboardRow(buttons[3]...),
		)
	}

	return keyboard
}

func (l *Logic) getKeyboard(list, value string) interface{} {
	var keyboard tgbotapi.ReplyKeyboardMarkup
	switch list {
	case Settings.String():
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(Home.String()),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(ChangeMyGroup.String()),
			),
			//tgbotapi.NewKeyboardButtonRow(
			//	tgbotapi.NewKeyboardButton(Support),
			//),
			//tgbotapi.NewKeyboardButtonRow(
			//	tgbotapi.NewKeyboardButton(Statistics),
			//),
			//tgbotapi.NewKeyboardButtonRow(
			//	tgbotapi.NewKeyboardButton(Polity),
			//),
		)
	case OtherSchedule.String():
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(Home.String()),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(GroupSchedule.String()),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(TeacherSchedule.String()),
			),
		)
	case MySchCode(0).Code(), MySchCode(1).Code(), MySchCode(2).Code(), MySchCode(3).Code(),
		MySchCode(4).Code(), MySchCode(5).Code(), MySchCode(6).Code(), MySchCode(7).Code():
		buttons := [][]tgbotapi.InlineKeyboardButton{{}, {}, {}, {}}
		day := int(time2.Now().Weekday())
		for i := 1; i <= 7 && i > 0; i++ {
			if i != day {
				buttonDay := fmt.Sprintf("%s: %s", time2.Weekday(i).ShortString(), time.Now().AddDate(0, 0, i-day).Format("02.01"))
				if len(buttons[0]) < 3 {
					buttons[0] = append(buttons[0], tgbotapi.NewInlineKeyboardButtonData(buttonDay, MySchCode(i).Code()))
				} else {
					buttons[1] = append(buttons[1], tgbotapi.NewInlineKeyboardButtonData(buttonDay, MySchCode(i).Code()))
				}
			}
		}

		buttonUpdate := fmt.Sprintf("Обновить - %s: %s", time2.Weekday(day).ShortString(), time.Now().Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonUpdate, MySchCode(day).Code()))
		//if day != int(day) {
		buttonToday := fmt.Sprintf("Сегодня - %s: %s", time2.Weekday(day).ShortString(), time.Now().Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonToday, MySchCode(day).Code()))
		//}

		buttonNext := fmt.Sprintf("Следующая неделя: %s - %s", time.Now().AddDate(0, 0, 7+1-day).Format("02.01"), time.Now().AddDate(0, 0, 7+7-day).Format("02.01"))
		buttons[3] = append(buttons[3], tgbotapi.NewInlineKeyboardButtonData(buttonNext, MySchNextCode(1).Code()))

		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(buttons[0]...),
			tgbotapi.NewInlineKeyboardRow(buttons[1]...),
			tgbotapi.NewInlineKeyboardRow(buttons[2]...),
			tgbotapi.NewInlineKeyboardRow(buttons[3]...),
		)
	case GroupSchCode(0).Code(value), GroupSchCode(1).Code(value), GroupSchCode(2).Code(value), GroupSchCode(3).Code(value),
		GroupSchCode(4).Code(value), GroupSchCode(5).Code(value), GroupSchCode(6).Code(value), GroupSchCode(7).Code(value):
		buttons := [][]tgbotapi.InlineKeyboardButton{{}, {}, {}, {}}
		day := int(time2.Now().Weekday())
		for i := 1; i <= 7 && i > 0; i++ {
			if i != day {
				buttonDay := fmt.Sprintf("%s: %s", time2.Weekday(i).ShortString(), time.Now().AddDate(0, 0, i-day).Format("02.01"))
				if len(buttons[0]) < 3 {
					buttons[0] = append(buttons[0], tgbotapi.NewInlineKeyboardButtonData(buttonDay, GroupSchCode(i).Code(value)))
				} else {
					buttons[1] = append(buttons[1], tgbotapi.NewInlineKeyboardButtonData(buttonDay, GroupSchCode(i).Code(value)))
				}
			}
		}

		buttonUpdate := fmt.Sprintf("Обновить - %s: %s", time2.Weekday(day).ShortString(), time.Now().Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonUpdate, GroupSchCode(day).Code(value)))
		//if day != t {
		buttonToday := fmt.Sprintf("Сегодня - %s: %s", time2.Weekday(day).ShortString(), time.Now().Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonToday, GroupSchCode(day).Code(value)))
		//}

		buttonNext := fmt.Sprintf("Следующая неделя: %s - %s", time.Now().AddDate(0, 0, 7+1-day).Format("02.01"), time.Now().AddDate(0, 0, 7+7-day).Format("02.01"))
		buttons[3] = append(buttons[3], tgbotapi.NewInlineKeyboardButtonData(buttonNext, GroupSchNextCode(1).Code(value)))

		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(buttons[0]...),
			tgbotapi.NewInlineKeyboardRow(buttons[1]...),
			tgbotapi.NewInlineKeyboardRow(buttons[2]...),
			tgbotapi.NewInlineKeyboardRow(buttons[3]...),
		)
	case TeacherSchCode(0).Code(value), TeacherSchCode(1).Code(value), TeacherSchCode(2).Code(value), TeacherSchCode(3).Code(value),
		TeacherSchCode(4).Code(value), TeacherSchCode(5).Code(value), TeacherSchCode(6).Code(value), TeacherSchCode(7).Code(value):
		buttons := [][]tgbotapi.InlineKeyboardButton{{}, {}, {}, {}}
		day := int(time2.Now().Weekday())
		for i := 1; i <= 7 && i > 0; i++ {
			if i != day {
				buttonDay := fmt.Sprintf("%s: %s", time2.Weekday(i).ShortString(), time.Now().AddDate(0, 0, i-day).Format("02.01"))
				if len(buttons[0]) < 3 {
					buttons[0] = append(buttons[0], tgbotapi.NewInlineKeyboardButtonData(buttonDay, TeacherSchCode(i).Code(value)))
				} else {
					buttons[1] = append(buttons[1], tgbotapi.NewInlineKeyboardButtonData(buttonDay, TeacherSchCode(i).Code(value)))
				}
			}
		}

		buttonUpdate := fmt.Sprintf("Обновить - %s: %s", time2.Weekday(day).ShortString(), time.Now().Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonUpdate, TeacherSchCode(day).Code(value)))
		//if day != t {
		buttonToday := fmt.Sprintf("Сегодня - %s: %s", time2.Weekday(day).ShortString(), time.Now().Format("02.01"))
		buttons[2] = append(buttons[2], tgbotapi.NewInlineKeyboardButtonData(buttonToday, TeacherSchCode(day).Code(value)))
		//}

		buttonNext := fmt.Sprintf("Следующая неделя: %s - %s", time.Now().AddDate(0, 0, 7+1-day).Format("02.01"), time.Now().AddDate(0, 0, 7+7-day).Format("02.01"))
		buttons[3] = append(buttons[3], tgbotapi.NewInlineKeyboardButtonData(buttonNext, TeacherSchNextCode(1).Code(value)))

		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(buttons[0]...),
			tgbotapi.NewInlineKeyboardRow(buttons[1]...),
			tgbotapi.NewInlineKeyboardRow(buttons[2]...),
			tgbotapi.NewInlineKeyboardRow(buttons[3]...),
		)
	case OtherButtons.String():
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(Home.String()),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(Polity.String()),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(Support.String()),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(Settings.String()),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(Statistics.String()),
			),
		)
	default:
		if strings.Contains(list, ChangeMyGroup.String()) {
			numString := strings.ReplaceAll(list, ChangeMyGroup.String(), "")
			numString = strings.TrimSpace(numString)
			if numString == "" || numString == ChangeMyGroup.String() {
				numString = "1"
			}
			num, err := strconv.Atoi(numString)
			if err != nil {
				log.Error(err)
			}

			var buttons [][]tgbotapi.KeyboardButton
			buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(Home.String())))

			if num > 1 {
				buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(ChangeMyGroup.String()+" "+strconv.Itoa(num-1))))
			}

			groups := l.schedule.GetGroups()
			for i := len(groups) - 50*(num-1) - 1; i >= 0 && i >= len(groups)-50*(num-1)-50 && i < len(groups); i-- {
				buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(groups[i])))
			}

			if num*50 < len(groups) {
				buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(ChangeMyGroup.String()+" "+strconv.Itoa(num+1))))
			}

			keyboard = tgbotapi.NewReplyKeyboard(buttons...)
		} else if strings.Contains(list, GroupSchedule.String()) {
			numString := strings.ReplaceAll(list, GroupSchedule.String(), "")
			numString = strings.TrimSpace(numString)
			if numString == "" || numString == GroupSchedule.String() {
				numString = "1"
			}
			num, err := strconv.Atoi(numString)
			if err != nil {
				log.Error(err)
			}

			var buttons [][]tgbotapi.KeyboardButton
			buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(Home.String())))

			if num > 1 {
				buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(GroupSchedule.String()+" "+strconv.Itoa(num-1))))
			}

			groups := l.schedule.GetGroups()
			for i := len(groups) - 50*(num-1) - 1; i >= 0 && i >= len(groups)-50*(num-1)-50 && i < len(groups); i-- {
				buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(groups[i])))
			}

			if num*50 < len(groups) {
				buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(GroupSchedule.String()+" "+strconv.Itoa(num+1))))
			}

			keyboard = tgbotapi.NewReplyKeyboard(buttons...)
		} else if strings.Contains(list, TeacherSchedule.String()) {
			numString := strings.ReplaceAll(list, TeacherSchedule.String(), "")
			numString = strings.TrimSpace(numString)
			if numString == "" || numString == TeacherSchedule.String() {
				numString = "1"
			}
			num, err := strconv.Atoi(numString)
			if err != nil {
				log.Error(err)
			}

			var buttons [][]tgbotapi.KeyboardButton
			buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(Home.String())))

			if num > 1 {
				buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(TeacherSchedule.String()+" "+strconv.Itoa(num-1))))
			}

			teachers := l.schedule.GetTeachers()
			for i := len(teachers) - 50*(num-1) - 1; i >= 0 && i >= len(teachers)-50*(num-1)-50 && i < len(teachers); i-- {
				buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(teachers[i])))
			}

			if num*50 < len(teachers) {
				buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(TeacherSchedule.String()+" "+strconv.Itoa(num+1))))
			}

			keyboard = tgbotapi.NewReplyKeyboard(buttons...)
		} else {
			keyboard = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(MySchedule.String()),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(OtherSchedule.String()),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(OtherButtons.String()),
				),
			)
		}
	}

	keyboard.ResizeKeyboard = true

	return keyboard
}
