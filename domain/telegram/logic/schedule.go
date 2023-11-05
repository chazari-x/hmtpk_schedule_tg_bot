package logic

import (
	"fmt"
	. "github.com/chazari-x/hmtpk_schedule/domain/telegram/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func getDay(day time.Weekday) string {
	return strings.ToLower(Weekday(day).String())
}

func (l *Logic) getMySchedule(message *tgbotapi.Message, dayName, date, group string, week int) tgbotapi.MessageConfig {
	if week == 0 {
		_, week = time.Now().ISOWeek()
	}
	if date == "" {
		date = time.Now().Format("02.01.2006")
	}
	var msg tgbotapi.MessageConfig
	dayName = strings.ToLower(dayName)
	if dayName == "" || dayName == "0" {
		dayName = getDay(time.Now().Weekday())
	} else {
		day, err := strconv.Atoi(dayName)
		if err != nil {
			log.Error(err)
			return tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении даты")
		}
		dayName = getDay(time.Weekday(day))
	}
	schs, err := l.schedule.GetScheduleByGroupID(group, date, week)
	if err != nil || len(schs) == 0 {
		if err != nil {
			log.Error(err)
		}
		return tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении расписания")
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, "")

	for _, sch := range schs {
		if strings.Contains(sch.Date, dayName) {
			msg.Text += fmt.Sprintf("Группа: %s\n%s", group, sch.Date)

			if len(sch.Lessons) == 0 {
				msg.Text += "\nРасписания нет (пар нет)"
				return msg
			}

			var lessonNum string
			var lessonText string
			for i, lesson := range sch.Lessons {
				if num, err := strconv.Atoi(lesson.Num); num%2 == 0 && err == nil {
					if lessonNum != lesson.Num {
						lessonText = ""
						lessonNum = lesson.Num
					} else {
						lessonNum = lesson.Num
					}

					if lesson.Num == sch.Lessons[i-1].Num {
						if lesson.Name[:len(lesson.Name)-4] != sch.Lessons[i-1].Name[:len(sch.Lessons[i-1].Name)-4] {
							lessonText += fmt.Sprintf("\n%s", lesson.Name)
						}
					} else {
						if lesson.Name[len(lesson.Name)-3:] == "(1)" || lesson.Name[len(lesson.Name)-3:] == "(2)" {
							lessonText += fmt.Sprintf("\n%s", lesson.Name[:len(lesson.Name)-3])
						} else {
							lessonText += fmt.Sprintf("\n%s", lesson.Name)
						}
					}
					if lesson.Room != "СРС" {
						lessonText += fmt.Sprintf("\n[%s]", lesson.Room)
					}
					if lesson.Teacher != "<>" {
						if lesson.Name[len(lesson.Name)-3:] == "(1)" {
							lessonText += fmt.Sprintf(" %s - подгр. 1", lesson.Teacher)
						} else if lesson.Name[len(lesson.Name)-3:] == "(2)" {
							lessonText += fmt.Sprintf(" %s - подгр. 2", lesson.Teacher)
						} else {
							lessonText += fmt.Sprintf(" %s", lesson.Teacher)
						}
					}

					if i < len(sch.Lessons)-1 {
						if lesson.Num == sch.Lessons[i+1].Num {
							continue
						}
					}

					if strings.HasSuffix(msg.Text, lessonText) {
						msg.Text = msg.Text[:len(msg.Text)-len(lessonText)] + fmt.Sprintf(" и %s [%s]", lesson.Num, lesson.Time) + msg.Text[len(msg.Text)-len(lessonText):]
					} else {
						msg.Text += fmt.Sprintf("</code>\n\n<code>Урок: %s [%s]", lesson.Num, lesson.Time) + lessonText
					}
				} else {
					if i == 0 {
						msg.Text += fmt.Sprintf("\n\n<code>Урок: %s", lesson.Num)
						msg.Text += fmt.Sprintf(" [%s]", lesson.Time)
						if lesson.Name[len(lesson.Name)-3:] == "(1)" || lesson.Name[len(lesson.Name)-3:] == "(2)" {
							msg.Text += fmt.Sprintf("\n%s", lesson.Name[:len(lesson.Name)-3])
						} else {
							msg.Text += fmt.Sprintf("\n%s", lesson.Name)
						}
						if lesson.Room != "СРС" {
							msg.Text += fmt.Sprintf("\n[%s]", lesson.Room)
						}
						if lesson.Teacher != "<>" {
							if lesson.Name[len(lesson.Name)-3:] == "(1)" {
								msg.Text += fmt.Sprintf(" %s - подгр. 1", lesson.Teacher)
							} else if lesson.Name[len(lesson.Name)-3:] == "(2)" {
								msg.Text += fmt.Sprintf(" %s - подгр. 2", lesson.Teacher)
							} else {
								msg.Text += fmt.Sprintf(" %s", lesson.Teacher)
							}
						}
					} else {
						if lesson.Num == sch.Lessons[i-1].Num {
							if lesson.Name[:len(lesson.Name)-4] != sch.Lessons[i-1].Name[:len(sch.Lessons[i-1].Name)-4] {
								msg.Text += fmt.Sprintf("\n%s", lesson.Name)
							}
						} else {
							msg.Text += fmt.Sprintf("</code>\n\n<code>Урок: %s", lesson.Num)
							msg.Text += fmt.Sprintf(" [%s]", lesson.Time)
							if lesson.Name[len(lesson.Name)-3:] == "(1)" || lesson.Name[len(lesson.Name)-3:] == "(2)" {
								msg.Text += fmt.Sprintf("\n%s", lesson.Name[:len(lesson.Name)-3])
							} else {
								msg.Text += fmt.Sprintf("\n%s", lesson.Name)
							}
						}
						if lesson.Room != "СРС" {
							msg.Text += fmt.Sprintf("\n[%s]", lesson.Room)
						}
						if lesson.Teacher != "<>" {
							if lesson.Name[len(lesson.Name)-3:] == "(1)" {
								msg.Text += fmt.Sprintf(" %s - подгр. 1", lesson.Teacher)
							} else if lesson.Name[len(lesson.Name)-3:] == "(2)" {
								msg.Text += fmt.Sprintf(" %s - подгр. 2", lesson.Teacher)
							} else {
								msg.Text += fmt.Sprintf(" %s", lesson.Teacher)
							}
						}
					}
				}
			}
		}
	}

	if msg.Text == "" {
		msg.Text += "Произошла какая-то ошибка"
	} else {
		msg.Text += "</code>"
	}

	return msg
}

func (l *Logic) getGroupSchedule(message *tgbotapi.Message, dayName, date string, week int) tgbotapi.MessageConfig {
	if date == "" {
		date = time.Now().Format("02.01.2006")
	}
	if week == 0 {
		_, week = time.Now().ISOWeek()
	}
	var msg tgbotapi.MessageConfig
	dayName = strings.ToLower(dayName)
	//group := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(message.Text, GroupScheduleIcoOne, ""), GroupScheduleIcoTwo, ""), " ", "")
	group := message.Text
	if dayName == "" || dayName == "0" {
		dayName = getDay(time.Now().Weekday())
	} else {
		day, err := strconv.Atoi(dayName)
		if err != nil {
			log.Error(err)
			return tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении даты")
		}
		dayName = getDay(time.Weekday(day))
	}
	schs, err := l.schedule.GetScheduleByGroupID(group, date, week)
	if err != nil || len(schs) == 0 {
		if err != nil {
			log.Error(err)
		}
		return tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении расписания")
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, "")

	for _, sch := range schs {
		if strings.Contains(sch.Date, dayName) {
			msg.Text += fmt.Sprintf("Группа: %s\n%s", group, sch.Date)

			if len(sch.Lessons) == 0 {
				msg.Text += "\nРасписания нет (пар нет)"
				return msg
			}

			var lessonNum string
			var lessonText string
			for i, lesson := range sch.Lessons {
				if num, err := strconv.Atoi(lesson.Num); num%2 == 0 && err == nil {
					if lessonNum != lesson.Num {
						lessonText = ""
						lessonNum = lesson.Num
					} else {
						lessonNum = lesson.Num
					}

					if lesson.Num == sch.Lessons[i-1].Num {
						if lesson.Name[:len(lesson.Name)-4] != sch.Lessons[i-1].Name[:len(sch.Lessons[i-1].Name)-4] {
							lessonText += fmt.Sprintf("\n%s", lesson.Name)
						}
					} else {
						if lesson.Name[len(lesson.Name)-3:] == "(1)" || lesson.Name[len(lesson.Name)-3:] == "(2)" {
							lessonText += fmt.Sprintf("\n%s", lesson.Name[:len(lesson.Name)-3])
						} else {
							lessonText += fmt.Sprintf("\n%s", lesson.Name)
						}
					}
					if lesson.Room != "СРС" {
						lessonText += fmt.Sprintf("\n[%s]", lesson.Room)
					}
					if lesson.Teacher != "<>" {
						if lesson.Name[len(lesson.Name)-3:] == "(1)" {
							lessonText += fmt.Sprintf(" %s - подгр. 1", lesson.Teacher)
						} else if lesson.Name[len(lesson.Name)-3:] == "(2)" {
							lessonText += fmt.Sprintf(" %s - подгр. 2", lesson.Teacher)
						} else {
							lessonText += fmt.Sprintf(" %s", lesson.Teacher)
						}
					}

					if i < len(sch.Lessons)-1 {
						if lesson.Num == sch.Lessons[i+1].Num {
							continue
						}
					}

					if strings.HasSuffix(msg.Text, lessonText) {
						msg.Text = msg.Text[:len(msg.Text)-len(lessonText)] + fmt.Sprintf(" и %s [%s]", lesson.Num, lesson.Time) + msg.Text[len(msg.Text)-len(lessonText):]
					} else {
						msg.Text += fmt.Sprintf("</code>\n\n<code>Урок: %s [%s]", lesson.Num, lesson.Time) + lessonText
					}
				} else {
					if i == 0 {
						msg.Text += fmt.Sprintf("\n\n<code>Урок: %s", lesson.Num)
						msg.Text += fmt.Sprintf(" [%s]", lesson.Time)
						if lesson.Name[len(lesson.Name)-3:] == "(1)" || lesson.Name[len(lesson.Name)-3:] == "(2)" {
							msg.Text += fmt.Sprintf("\n%s", lesson.Name[:len(lesson.Name)-3])
						} else {
							msg.Text += fmt.Sprintf("\n%s", lesson.Name)
						}
						if lesson.Room != "СРС" {
							msg.Text += fmt.Sprintf("\n[%s]", lesson.Room)
						}
						if lesson.Teacher != "<>" {
							if lesson.Name[len(lesson.Name)-3:] == "(1)" {
								msg.Text += fmt.Sprintf(" %s - подгр. 1", lesson.Teacher)
							} else if lesson.Name[len(lesson.Name)-3:] == "(2)" {
								msg.Text += fmt.Sprintf(" %s - подгр. 2", lesson.Teacher)
							} else {
								msg.Text += fmt.Sprintf(" %s", lesson.Teacher)
							}
						}
					} else {
						if lesson.Num == sch.Lessons[i-1].Num {
							if lesson.Name[:len(lesson.Name)-4] != sch.Lessons[i-1].Name[:len(sch.Lessons[i-1].Name)-4] {
								msg.Text += fmt.Sprintf("\n%s", lesson.Name)
							}
						} else {
							msg.Text += fmt.Sprintf("</code>\n\n<code>Урок: %s", lesson.Num)
							msg.Text += fmt.Sprintf(" [%s]", lesson.Time)
							if lesson.Name[len(lesson.Name)-3:] == "(1)" || lesson.Name[len(lesson.Name)-3:] == "(2)" {
								msg.Text += fmt.Sprintf("\n%s", lesson.Name[:len(lesson.Name)-3])
							} else {
								msg.Text += fmt.Sprintf("\n%s", lesson.Name)
							}
						}
						if lesson.Room != "СРС" {
							msg.Text += fmt.Sprintf("\n[%s]", lesson.Room)
						}
						if lesson.Teacher != "<>" {
							if lesson.Name[len(lesson.Name)-3:] == "(1)" {
								msg.Text += fmt.Sprintf(" %s - подгр. 1", lesson.Teacher)
							} else if lesson.Name[len(lesson.Name)-3:] == "(2)" {
								msg.Text += fmt.Sprintf(" %s - подгр. 2", lesson.Teacher)
							} else {
								msg.Text += fmt.Sprintf(" %s", lesson.Teacher)
							}
						}
					}
				}
			}
		}
	}

	if msg.Text == "" {
		msg.Text += "Произошла какая-то ошибка"
	} else {
		msg.Text += "</code>"
	}

	return msg
}

func (l *Logic) getTeacherSchedule(message *tgbotapi.Message, dayName, date string, week int) tgbotapi.MessageConfig {
	if date == "" {
		date = time.Now().Format("02.01.2006")
	}
	if week == 0 {
		_, week = time.Now().ISOWeek()
	}
	var msg tgbotapi.MessageConfig
	dayName = strings.ToLower(dayName)
	//teacher := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(message.Text, TeacherScheduleIcoOne+" ", ""), " "+TeacherScheduleIcoTwo, ""), " ", "+")
	teacher := strings.ReplaceAll(message.Text, " ", "+")
	if dayName == "" || dayName == "0" {
		dayName = getDay(time.Now().Weekday())
	} else {
		atoi, err := strconv.Atoi(dayName)
		if err != nil {
			log.Error(err)
			return tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении даты")
		}
		dayName = getDay(time.Weekday(atoi))
	}
	schs, err := l.schedule.GetScheduleByTeacher(teacher, date, week)
	if err != nil || len(schs) == 0 {
		if err != nil {
			log.Error(err)
		}
		return tgbotapi.NewMessage(message.Chat.ID, "Произошла ошибка при получении расписания")
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, "")

	for _, sch := range schs {
		if strings.Contains(sch.Date, dayName) {
			msg.Text += fmt.Sprintf("Преподаватель: %s\n%s", strings.ReplaceAll(teacher, "+", " "), sch.Date)

			if len(sch.Lessons) == 0 {
				msg.Text += "\nРасписания нет (пар нет)"
				return msg
			}

			for i, lesson := range sch.Lessons {
				if num, err := strconv.Atoi(lesson.Num); num%2 == 0 && err == nil {
					var lessonText string

					if lesson.Name[len(lesson.Name)-3:] == "(1)" || lesson.Name[len(lesson.Name)-3:] == "(2)" {
						lessonText += fmt.Sprintf("\n%s", lesson.Name[:len(lesson.Name)-3])
					} else {
						lessonText += fmt.Sprintf("\n%s", lesson.Name)
					}
					if lesson.Name[len(lesson.Name)-3:] == "(1)" {
						lessonText += fmt.Sprintf("\n[%s] %s - подгр. 1", lesson.Room, lesson.Group)
					} else if lesson.Name[len(lesson.Name)-3:] == "(2)" {
						lessonText += fmt.Sprintf("\n[%s] %s - подгр. 2", lesson.Room, lesson.Group)
					} else {
						lessonText += fmt.Sprintf("\n[%s] %s", lesson.Room, lesson.Group)
					}

					if strings.HasSuffix(msg.Text, lessonText) {
						msg.Text = msg.Text[:len(msg.Text)-len(lessonText)] + fmt.Sprintf(" и %s [%s]", lesson.Num, lesson.Time) + msg.Text[len(msg.Text)-len(lessonText):]
					} else {
						msg.Text += fmt.Sprintf("</code>\n\n<code>Урок: %s [%s]", lesson.Num, lesson.Time) + lessonText
					}
				} else {
					if i == 0 {
						msg.Text += fmt.Sprintf("\n\n<code>Урок: %s", lesson.Num)
					} else {
						msg.Text += fmt.Sprintf("</code>\n\n<code>Урок: %s", lesson.Num)
					}
					msg.Text += fmt.Sprintf(" [%s]", lesson.Time)
					if lesson.Name[len(lesson.Name)-3:] == "(1)" || lesson.Name[len(lesson.Name)-3:] == "(2)" {
						msg.Text += fmt.Sprintf("\n%s", lesson.Name[:len(lesson.Name)-3])
					} else {
						msg.Text += fmt.Sprintf("\n%s", lesson.Name)
					}
					if lesson.Name[len(lesson.Name)-3:] == "(1)" {
						msg.Text += fmt.Sprintf("\n[%s] %s - подгр. 1", lesson.Room, lesson.Group)
					} else if lesson.Name[len(lesson.Name)-3:] == "(2)" {
						msg.Text += fmt.Sprintf("\n[%s] %s - подгр. 2", lesson.Room, lesson.Group)
					} else {
						msg.Text += fmt.Sprintf("\n[%s] %s", lesson.Room, lesson.Group)
					}
				}
			}
		}
	}

	if msg.Text == "" {
		msg.Text += "Произошла какая-то ошибка"
	} else {
		msg.Text += "</code>"
	}

	return msg
}
