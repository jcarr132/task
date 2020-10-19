/* Tasklist handles interaction with the JSON database. This package
defines the TaskList struct type, which wraps the database connection, and
the Task struct, which holds data for a single task.  */
package tasklist

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"

	"encoding/json"
)

/* TaskList simply wraps the database driver and provides methods for acting
on the list as a whole.  */
type TaskList struct {
	Db *bolt.DB
}

/* NewTaskList returns a new TaskList struct containing a connection to the
database.  */
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

/* Tasks() queries the database returns a slice containing the tasks stored
within.  */
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

/* AddTask accepts a Task struct as an argument and saves it to the database with
its UUID (Task.TaskId) as the database key. */
func (tl TaskList) AddTask(task Task) error {
	db, err := bolt.Open("data", 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("tasks"))

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

/* CompleteTask sets a Task's `complete` field to `true` and re-adds it to the
database, overwriting the previous version.  */
func (tl TaskList) CompleteTask(task Task) {
	task.Complete = true
	tl.AddTask(task)
}

/* UncompleteTask sets a Task's `complete` field to `false` and re-adds it to the
database, overwriting the previous version. */
func (tl TaskList) UncompleteTask(task Task) {
	task.Complete = false
	tl.AddTask(task)
}

/* ToggleComplete changes the `complete` field of a Task from `true` to `false` or
from `true` to `false` as appropriate and re-adds it to the database, overwriting
the previous version. */
func (tl TaskList) ToggleComplete(task Task) {
	if task.Complete == true {
		task.Complete = false
	} else {
		task.Complete = true
	}
	tl.AddTask(task)
}

/* SelectTask prints an enumerated list of tasks to stdout and accepts an integer
input from the user indicating which Task struct to return. Used in conjunction with
another method that accepts a Task struct such as TaskList.Complete(...).

Example:
					tl.CompleteTask(tl.SelectTask())
*/
func (tl TaskList) SelectTask() Task {
	tasks := tl.Tasks()
	for i, task := range tasks {
		fmt.Println(i+1, task)
	}

	fmt.Print("Enter selection: ")
	var selection int
	_, err := fmt.Scanf("%d", &selection)
	if err != nil {
		panic(err)
	}

	return tasks[selection-1]
}

/* RemoveTask deletes a task from the database. */
func (tl TaskList) RemoveTask(task Task) {
	db, err := bolt.Open("data", 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("tasks"))

		key, err := task.TaskId.MarshalBinary()
		if err != nil {
			return err
		}

		return bucket.Delete(key)
	})
}

/* The Task struct holds data about a task. Each Task is assigned a random UUID
which is used as it's primary identifier.  */
type Task struct {
	TaskId   uuid.UUID `json:"taskid"`
	Name     string    `json:"name"`
	Complete bool      `json:"complete"`
	Notes    string    `json:"notes"`
	// TODO implement the rest of the fields
	// deadline/timeslot
	// tags
	// priority
	// subtasks
}

/* NewTask returns a new Task struct with the given name and randomly generated
UUID.  By default, the new Task is incomplete (Task.Complete = false), and
has no notes associated with it.*/
func NewTask(name string) Task {
	return Task{
		TaskId:   uuid.New(),
		Name:     name,
		Complete: false,
		Notes:    "",
	}
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
