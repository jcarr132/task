package tasklist

import (
	"github.com/google/uuid"
	"github.com/sonyarouje/simdb/db"
)

type TaskList struct {
	Db db.Driver
}

func NewTasklist() TaskList {
	driver, err := db.New("data")
	if err != nil {
		panic(err)
	}

	return TaskList{
		Db: *driver,
	}
}

func (tl TaskList) Tasks() []Task {
	var tasks []Task
	err := tl.Db.Open(Task{}).Get().AsEntity(&tasks)
	if err != nil {
		panic(err)
	}

	return tasks
}

type Task struct {
	TaskID uuid.UUID
	Name   string `json:"name"`
	Notes  string `json:"notes"`
}

// func TaskFromJSON()

func (t Task) ID() (jsonField string, value interface{}) {
	value = t.TaskID
	jsonField = "id"
	return
}
