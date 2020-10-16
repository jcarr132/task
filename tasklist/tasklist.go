/*
Tasklist handles interaction with the JSON database (`simdb`). This package
defines the TaskList struct type, which wraps the database connection, and
the Task struct, which holds data for a single task.
*/
package tasklist

import (
	"github.com/google/uuid"
	"github.com/sonyarouje/simdb/db"
)

/*
TaskList simply wraps the database driver and provides methods for acting
on the list as a whole.
*/
type TaskList struct {
	Db db.Driver
}

/*
NewTaskList returns a new TaskList struct containing a connection to the
`simdb` database.
*/
func NewTasklist() TaskList {
	driver, err := db.New("data")
	if err != nil {
		panic(err)
	}

	return TaskList{
		Db: *driver,
	}
}

/*
Tasks() queries the database returns a slice containing the tasks stored
within.
*/
func (tl TaskList) Tasks() []Task {
	var tasks []Task
	err := tl.Db.Open(Task{}).Get().AsEntity(&tasks)
	if err != nil {
		panic(err)
	}

	return tasks
}

// TODO docstring
func (tl TaskList) AddTask(task Task) {
	err := tl.Db.Insert(task)
	if err != nil {
		panic(err)
	}
}

// TODO docstring
// FIXME Db.Update doesn't correctly match on UUID
func (tl TaskList) CompleteTask(task Task) {
	task.Complete = true
	err := tl.Db.Update(task)
	if err != nil {
		panic(err)
	}
}

/*
The Task struct holds data about a task. Each Task is assigned a random UUID
which is used as it's primary identifier.
*/
type Task struct {
	TaskID   uuid.UUID `json:"taskid"`
	Name     string    `json:"name"`
	Complete bool      `json:"complete"`
	Notes    string    `json:"notes"`
	// TODO deadline/timeslot
	// TODO tags
	// TODO priority
	// TODO subtasks
}

/*
Task implements ID to conform to the database library (`simdb`) requirements.
*/
func (t Task) ID() (jsonField string, value interface{}) {
	value = t.TaskID
	jsonField = "taskid"
	return
}

// TODO docstring
func NewTask(name string) Task {
	return Task{
		TaskID:   uuid.New(),
		Name:     name,
		Complete: false,
		Notes:    "",
	}
}
