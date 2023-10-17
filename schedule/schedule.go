package schedule

import (
	"encoding/json"
	"fmt"
	"github.com/chazari-x/hmtpk_schedule/config"
	"github.com/chazari-x/hmtpk_schedule/redis/redis"
	sel "github.com/chazari-x/hmtpk_schedule/selenium"
	"github.com/tebeka/selenium"
	"io"
	"net/http"
	"time"
)

type Schedule struct {
	cfg      *config.Schedule
	r        *redis.Redis
	selenium *sel.Selenium
}

func NewSchedule(cfg *config.Schedule, r *redis.Redis, selenium *sel.Selenium) *Schedule {
	return &Schedule{cfg, r, selenium}
}

func (s *Schedule) GetGroups() []string {
	var groups []string

	if s.r != nil {
		if err := s.selenium.OpenURL("https://hmtpk.ru/ru/students/schedule/"); err != nil {
			return nil
		}

		for i := 2; ; i++ {
			value, err := s.selenium.GetElementValue(selenium.ByCSSSelector, fmt.Sprintf("#group > option:nth-child(%d)", i))
			if err != nil {
				break
			}

			name, err := s.selenium.GetElementText(selenium.ByCSSSelector, fmt.Sprintf("#group > option:nth-child(%d)", i))
			if err != nil {
				break
			}

			fmt.Printf("- id: %s\n  name:\"%s\"\n", value, name)
			groups = append(groups, value+" - "+name)
		}

		return groups
	}

	data, err := s.r.GetData("groups")
	if err != nil {
		for _, g := range s.cfg.Groups {
			groups = append(groups, g.Name)
		}

		return groups
	}

	if err = json.Unmarshal([]byte(data), &groups); err != nil {
		for _, g := range s.cfg.Groups {
			groups = append(groups, g.Name)
		}

		return groups
	}

	return groups
}

func (s *Schedule) GetTeachers() []string {
	var teachers []string

	if s.r == nil {
		for _, g := range s.cfg.Teachers {
			teachers = append(teachers, g.Name)
		}

		return teachers
	}

	data, err := s.r.GetData("teachers")
	if err != nil {
		for _, g := range s.cfg.Groups {
			teachers = append(teachers, g.Name)
		}

		return teachers
	}

	if err = json.Unmarshal([]byte(data), &teachers); err != nil {
		for _, g := range s.cfg.Groups {
			teachers = append(teachers, g.Name)
		}

		return teachers
	}

	return teachers
}

func (s *Schedule) GetScheduleByGroupID(group int, date time.Time) ([]byte, error) {
	href := fmt.Sprintf("https://hmtpk.ru/ru/students/schedule/?group=%d&date_edu1c=%s&send=Показать#current", group, date.String())
	r, err := http.Get(href)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = r.Body.Close()
	}()

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", r.StatusCode)
	}

	// todo чтение хтмл и конвертация в структуру

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//
//func (c *Schedule) GetScheduleByTeacher(teacher string, date time.Time) ([]byte, error) {
//	href := fmt.Sprintf("https://hmtpk.ru/ru/students/schedule/?group=%s&date_edu1c=%s&send=Показать#current", teacher, date.String())
//	r, err := http.Get(href)
//	if err != nil {
//		return nil, err
//	}
//	defer func() {
//		_ = r.Body.Close()
//	}()
//
//	if r.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("%d", r.StatusCode)
//	}
//
//	data, err := io.ReadAll(r.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	return data, nil
//}
