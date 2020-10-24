/* Tasklist handles interaction with the JSON database. This package
defines the TaskList struct type, which wraps the database connection, and
the Task struct, which holds data for a single task. */
package tasklist

// TODO: Implement priorities and priority sorting. When viewing the list,
// add a visual indication of priority tasks.

// TODO: Implement deadlines/timeslots could be an array with two values,
// second value nil if it's a deadline, otherwise its a timeslot

// TODO: ordering of items. should keep a list of the order of tasks and when
// retrieving from the database, present them in the same order. Order should
// be manually changeable.

// TODO: better printing. use a table formatting package or something.

import (
	"errors"
	"fmt"

	"github.com/boltdb/bolt"

	"encoding/binary"
	"encoding/json"
)

/* TaskList simply wraps the database driver and provides methods for acting
on the list as a whole. */
type TaskList struct {
	Db *bolt.DB
}

/* NewTaskList returns a new TaskList struct containing a connection to the
database. */
func NewTasklist(dbpath string) (TaskList, error) {
	fmt.Println(dbpath)
	db, err := bolt.Open(dbpath, 0600, nil)
	if err != nil {
		return TaskList{}, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("tasks"))
		return err
	})
	if err != nil {
		return TaskList{}, err
	}

	tl := TaskList{
		Db: db,
	}

	return tl, nil
}

/* Tasks() queries the database returns a slice containing the tasks stored
within. */
func (tl *TaskList) Tasks() ([]Task, error) {
	var tasks []Task

	err := tl.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tasks"))
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
	if err != nil {
		return tasks, err
	}

	return tasks, nil
}

/* AddTask accepts a Task struct as an argument and saves it to the database with
its UUID (Task.TaskId) as the database key. */
func (tl TaskList) AddTask(task Task) error {
	// TODO add multiple tasks in one call
	return tl.Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tasks"))

		id, _ := bucket.NextSequence()
		task.TaskId = int(id)
		buf, err := json.Marshal(task)
		if err != nil {
			return err
		}

		return bucket.Put(itob(task.TaskId), buf)
	})
}

/* UpdateTask has the same functionality as AddTask except that it saves the item
to the databse with the same key rather than assigning a new one. */
func (tl TaskList) UpdateTask(task Task) error {
	return tl.Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tasks"))

		buf, err := json.Marshal(task)
		if err != nil {
			return err
		}
		return bucket.Put(itob(task.TaskId), buf)
	})
}

/* CompleteTask sets a Task's `complete` field to `true` and re-adds it to the
database, overwriting the previous version. */
func (tl TaskList) CompleteTask(task Task) error {
	task.Complete = true
	return tl.UpdateTask(task)
}

/* UncompleteTask sets a Task's `complete` field to `false` and re-adds it to the
database, overwriting the previous version. */
func (tl TaskList) UncompleteTask(task Task) error {
	task.Complete = false
	return tl.UpdateTask(task)
}

/* ToggleComplete changes the `complete` field of a Task from `true` to `false` or
from `true` to `false` as appropriate and re-adds it to the database, overwriting
the previous version. */
func (tl TaskList) ToggleComplete(task Task) error {
	if task.Complete == true {
		task.Complete = false
	} else {
		task.Complete = true
	}
	return tl.UpdateTask(task)
}

// TODO docstring
func (tl TaskList) SetPriority(task Task, p int) error {
	task.Priority = p
	return tl.UpdateTask(task)
}

/* SelectTask prints an enumerated list of tasks to stdout and accepts an integer
input from the user indicating which Task struct to return. Used in conjunction with
another method that accepts a Task struct such as TaskList.Complete(...).

Example:
					tl.CompleteTask(tl.SelectTask())
*/
func (tl TaskList) SelectTask(selection int) (Task, error) {

	tasks, err := tl.Tasks()
	if err != nil {
		return Task{}, err
	}

	if selection != 0 {
		return tasks[selection-1], nil
	}

	for i, task := range tasks {
		fmt.Println(i+1, task)
	}

	fmt.Print("Enter selection: ")
	_, err = fmt.Scanf("%d", &selection)
	if err != nil {
		return Task{}, err
	}
	if selection < 1 || selection > len(tasks) {
		return Task{}, errors.New("invalid selection")
	}

	return tasks[selection-1], nil
}

/* RemoveTask deletes a task from the database. */
func (tl TaskList) RemoveTask(task Task) error {
	return tl.Db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("tasks"))

		key := task.TaskId
		if err != nil {
			return err
		}

		return bucket.Delete(itob(key))
	})
}

/* The Task struct holds data about a task. Each Task is assigned a random UUID
which is used as it's primary identifier. */
type Task struct {
	TaskId   int
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
	Notes    string `json:"notes"`
	Priority int    `json:"priority"`
	// TODO implement the rest of the fields
	// deadline/timeslot
	// tags
	// priority
	// subtasks
}

/* NewTask returns a new Task struct with the given name By default, the new
* Task is incomplete, has zero priority, and has no notes associated with it.
* */
func NewTask(name string) (Task, error) {
	if len(name) < 1 {
		return Task{}, errors.New("cannot create a task without a name")
	}

	task := Task{
		TaskId:   0,
		Name:     name,
		Complete: false,
		Notes:    "",
	}

	return task, nil
}

/* String describes how the string representation of a Task struct and enables
printing with fmt.Println(). */
func (t Task) String() string {
	var checkbox string

	if t.Complete == true {
		checkbox = "[x]"
	} else {
		checkbox = "[ ]"
	}

	return fmt.Sprintf("%s - %s", checkbox, t.Name)
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
