package model

type Schedule struct {
	Date    string   `json:"date"`
	Lessons []Lesson `json:"lesson"`
}

type Lesson struct {
	Num     string `json:"num"`
	Time    string `json:"time"`
	Name    string `json:"name"`
	Room    string `json:"room"`
	Group   string `json:"group"`
	Teacher string `json:"teacher"`
}
