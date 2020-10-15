package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/sonyarouje/simdb/db"
	"github.com/urfave/cli/v2"
)

func main() {
	driver, err := db.New("data")
	if err != nil {
		panic(err)
	}

	task := Task{
		TaskID: uuid.New(),
		Name:   "example task",
	}

	err = driver.Insert(task)
	if err != nil {
		panic(err)
	}

	app := &cli.App{
		Name:  "task",
		Usage: "manage tasks from the terminal",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all tasks",
				Action: func(c *cli.Context) error {
					fmt.Println("running 'task list'")
					return nil
				},
			},
			{
				// Next subcommand
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type Task struct {
	TaskID uuid.UUID
	Name   string `json:"name"`
}

func (t Task) ID() (jsonField string, value interface{}) {
	value = t.TaskID
	jsonField = "id"
	return
}
