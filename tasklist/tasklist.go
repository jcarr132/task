/*
Tasklist handles interaction with the JSON database. This package
defines the TaskList struct type, which wraps the database connection, and
the Task struct, which holds data for a single task.
*/
package tasklist

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"

	"encoding/json"
	// "fmt"
)

/*
TaskList simply wraps the database driver and provides methods for acting
on the list as a whole.
*/
type TaskList struct {
	Db *bolt.DB
}

/*
NewTaskList returns a new TaskList struct containing a connection to the
database.
*/
func NewTasklist() TaskList {
	db, err := bolt.Open("data", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("tasks"))
		return err
	})

	return TaskList{
		Db: db,
	}
}

/*
Tasks() queries the database returns a slice containing the tasks stored
within.
*/
func (tl TaskList) Tasks() []Task {
	db, err := bolt.Open("data", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var tasks []Task

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tasks"))
		if err != nil {
			panic(err)
		}

		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var task Task
			err := json.Unmarshal(v, &task)
			if err != nil {
				return err
			}
			tasks = append(tasks, task)
		}
		return nil
	})

	return tasks
}

// TODO docstring
func (tl TaskList) AddTask(task Task) error {
	db, err := bolt.Open("data", 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("tasks"))
		// bucket := tx.Bucket([]byte("tasks"))

		buf, err := json.Marshal(task)
		if err != nil {
			return err
		}

		key, err := task.TaskId.MarshalBinary()
		if err != nil {
			return err
		}

		return bucket.Put(key, buf)
	})

	return err
}

// TODO docstring
// TODO reimplement with boltdb
func (tl TaskList) CompleteTask(task Task) {
	task.Complete = true
	tl.AddTask(task)
}

/*
The Task struct holds data about a task. Each Task is assigned a random UUID
which is used as it's primary identifier.
*/
type Task struct {
	TaskId   uuid.UUID `json:"taskid"`
	Name     string    `json:"name"`
	Complete bool      `json:"complete"`
	Notes    string    `json:"notes"`
	// TODO deadline/timeslot
	// TODO tags
	// TODO priority
	// TODO subtasks
}

// TODO docstring
func NewTask(name string) Task {
	return Task{
		TaskId:   uuid.New(),
		Name:     name,
		Complete: false,
		Notes:    "",
	}
}

func (t Task) String() string {
	return fmt.Sprintf("Task: %s\nComplete: %v\n", t.Name, t.Complete)
}
