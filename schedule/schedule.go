package schedule

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chazari-x/hmtpk_schedule/config"
	"github.com/chazari-x/hmtpk_schedule/redis"
	"github.com/chazari-x/hmtpk_schedule/schedule/model"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Schedule struct {
	cfg *config.Schedule
	r   *redis.Redis
}

var Groups map[string]int

func NewSchedule(cfg *config.Schedule, r *redis.Redis) *Schedule {
	s := &Schedule{cfg, r}
	Groups = make(map[string]int)
	for _, group := range cfg.Groups {
		Groups[group.Name] = group.ID
	}

	return s
}

func (s *Schedule) GetGroups() []string {
	var groups []string

	for _, g := range s.cfg.Groups {
		groups = append(groups, g.Name)
	}

	return groups
}

func (s *Schedule) GetGroup(groupName string) string {
	if Groups[groupName] != 0 {
		return groupName
	}

	return ""
}

func (s *Schedule) GetTeachers() []string {
	var teachers []string

	for _, g := range s.cfg.Teachers {
		teachers = append(teachers, g.Name)
	}

	return teachers
}

func (s *Schedule) GetScheduleByGroupID(group, date string, week int) ([]model.Schedule, error) {
	var weeklySchedule []model.Schedule

	if week == 0 {
		_, week = time.Now().ISOWeek()
	}

	log.Trace(group)

	if redisWeeklySchedule, err := s.r.Get(strconv.Itoa(week) + ":" + strconv.Itoa(Groups[group])); err == nil && redisWeeklySchedule != "" {
		if json.Unmarshal([]byte(redisWeeklySchedule), &weeklySchedule) == nil {
			log.Trace("Данные получены из Redis")
			return weeklySchedule, nil
		}
	}

	href := fmt.Sprintf("https://hmtpk.ru/ru/students/schedule/?group=%d&date_edu1c=%s&send=Показать#current", Groups[group], date)
	resp, err := http.Post(href, "", nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Ошибка: %s", resp.Status))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	for scheduleElementNum := 2; scheduleElementNum <= 8; scheduleElementNum++ {
		scheduleDateElement := doc.Children().Find(fmt.Sprintf("div.raspcontent.m5 div:nth-child(%d) div.panel-heading.edu_today > h2", scheduleElementNum))
		weeklySchedule = append(weeklySchedule, model.Schedule{
			Date: scheduleDateElement.Text(),
		})

		lessonsElement := doc.Children().Find(fmt.Sprintf("div.raspcontent.m5 div:nth-child(%d) div.panel-body > #mobile-friendly > tbody:nth-child(2)", scheduleElementNum))
		var lessons []model.Lesson
		for lessonNum := 1; lessonNum < 14; lessonNum++ {
			var lesson model.Lesson
			var exists bool
			lessonElement := lessonsElement.Find(fmt.Sprintf("tr:nth-child(%d)", lessonNum))
			for lessonAttributeNum := 1; lessonAttributeNum <= 5; lessonAttributeNum++ {
				lessonElementAttribute := lessonElement.Find(fmt.Sprintf("td:nth-child(%d)", lessonAttributeNum))
				var value string
				value, exists = lessonElementAttribute.Attr("data-title")
				if exists {
					switch value {
					case "Номер урока":
						lesson.Num = lessonElementAttribute.Text()
					case "Время":
						if lesson.Num == "" {
							lesson.Num = lessons[len(lessons)-1].Num
						}
						lesson.Time = lessonElementAttribute.Text()
					case "Название предмета":
						lesson.Name = lessonElementAttribute.Text()
					case "Кабинет":
						lesson.Room = lessonElementAttribute.Text()
					case "Преподаватель":
						lesson.Teacher = lessonElementAttribute.Text()
					}
				} else if lessonAttributeNum == 5 {
					exists = true
				} else if lessonAttributeNum == 1 {
					break
				}
			}

			if exists {
				lessons = append(lessons, lesson)
			} else {
				break
			}
		}

		weeklySchedule[len(weeklySchedule)-1].Lessons = lessons
	}

	//log.Trace(weeklySchedule)

	if marshal, err := json.Marshal(weeklySchedule); err == nil {
		if err := s.r.Set(strconv.Itoa(week)+":"+strconv.Itoa(Groups[group]), string(marshal)); err != nil {
			log.Error(err)
		} else {
			log.Trace("Данные сохранены в Redis")
		}
	}

	return weeklySchedule, nil
}

func (s *Schedule) GetScheduleByTeacher(teacher, date string, week int) ([]model.Schedule, error) {
	var weeklySchedule []model.Schedule

	if week == 0 {
		_, week = time.Now().ISOWeek()
	}

	log.Trace(teacher)

	if redisWeeklySchedule, err := s.r.Get(strconv.Itoa(week) + ":" + teacher); err == nil && redisWeeklySchedule != "" {
		if json.Unmarshal([]byte(redisWeeklySchedule), &weeklySchedule) == nil {
			log.Trace("Данные получены из Redis")
			return weeklySchedule, nil
		}
	}

	href := fmt.Sprintf("https://hmtpk.ru/ru/teachers/schedule/?teacher=%s&date_edu1c=%s&send=Показать#current", teacher, date)
	resp, err := http.Post(href, "", nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Ошибка: %s", resp.Status))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	for scheduleElementNum := 1; scheduleElementNum <= 7; scheduleElementNum++ {
		scheduleDateElement := doc.Children().Find(fmt.Sprintf("div.raspcontent.m5 div:nth-child(%d) div.panel-heading.edu_today > h2", scheduleElementNum))
		weeklySchedule = append(weeklySchedule, model.Schedule{
			Date: scheduleDateElement.Text(),
		})

		lessonsElement := doc.Children().Find(fmt.Sprintf("div.raspcontent.m5 div:nth-child(%d) div.panel-body > table.table > tbody:nth-child(2)", scheduleElementNum))
		var lessons []model.Lesson
		for lessonNum := 1; lessonNum < 14; lessonNum++ {
			var lesson model.Lesson
			lessonElement := lessonsElement.Find(fmt.Sprintf("tr:nth-child(%d)", lessonNum))
			for lessonAttributeNum := 1; lessonAttributeNum <= 5; lessonAttributeNum++ {
				lessonElementAttribute := lessonElement.Find(fmt.Sprintf("td:nth-child(%d)", lessonAttributeNum))
				value := lessonElementAttribute.Text()
				if value == "" {
					break
				}

				value = strings.ReplaceAll(value, "\n", "")
				value = strings.TrimSpace(value)
				switch lessonAttributeNum {
				case 1:
					lesson.Num = value
				case 2:
					lesson.Time = value
				case 3:
					lesson.Name = value
				case 4:
					lesson.Group = value
				case 5:
					room := strings.TrimSpace(regexp.MustCompile("\\W-[0-9]{1,3}$").FindString(value))
					if room == "" {
						lesson.Room = fmt.Sprintf("%s", strings.TrimSpace(strings.TrimRight(value, room)))
					} else {
						lesson.Room = fmt.Sprintf("%s (%s)", room, strings.TrimSpace(strings.TrimRight(value, room)))
					}
				}
			}

			if lesson.Num != "" {
				lessons = append(lessons, lesson)
			}
		}

		weeklySchedule[len(weeklySchedule)-1].Lessons = lessons
	}

	//log.Trace(weeklySchedule)

	if marshal, err := json.Marshal(weeklySchedule); err == nil {
		if err := s.r.Set(strconv.Itoa(week)+":"+teacher, string(marshal)); err != nil {
			log.Error(err)
		} else {
			log.Trace("Данные сохранены в Redis")
		}
	}

	return weeklySchedule, nil
}
