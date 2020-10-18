/*
Tasklist handles interaction with the JSON database. This package
defines the TaskList struct type, which wraps the database connection, and
the Task struct, which holds data for a single task.
*/
package tasklist

import (
	"github.com/boltdb/bolt"
	"github.com/google/uuid"

	"encoding/binary"
	"encoding/json"
	"fmt"
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

	return TaskList{
		Db: db,
	}
}

/*
Tasks() queries the database returns a slice containing the tasks stored
within.
*/
// FIXME
func (tl TaskList) Tasks() []Task {
	var tasks []Task

	tl.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		c := b.Cursor()

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

	fmt.Println("len(tasks) = ", len(tasks))
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
		b, err := tx.CreateBucketIfNotExists([]byte("tasks"))

		buf, err := json.Marshal(task)
		if err != nil {
			return err
		}

		key, err := task.TaskId.MarshalBinary()
		if err != nil {
			return err
		}

		return b.Put(key, buf)
	})

	return err
}

// TODO docstring
// FIXME Db.Update doesn't correctly match on UUID
// func (tl TaskList) CompleteTask(task Task) {
// task.Complete = true
// err := tl.Db.Update(task)
// if err != nil {
// 	panic(err)
// }
// }

// func LoadFromDB(db bolt.DB) []Task {
// 	db.View(func(tx *bolt.Tx) error {

// 		return nil
// 	})

// 	return
// }

// func SaveToDB

/*
The Task struct holds data about a task. Each Task is assigned a random UUID
which is used as it's primary identifier.
*/
type Task struct {
	TaskId uuid.UUID `json:"taskid"`
	// TaskId   int       `json:"taskid"`
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
	Notes    string `json:"notes"`
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

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
