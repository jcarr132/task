package main

import (
	"fmt"
	"log"
	"os"
	// "reflect"

	"task/tasklist"

	"github.com/urfave/cli/v2"
)

func main() {

	tl := tasklist.NewTasklist()

	tl.AddTask(tasklist.NewTask("function test"))
	tl.AddTask(tasklist.NewTask("function test"))
	tl.AddTask(tasklist.NewTask("function test"))

	// t := tl.Tasks()[0]
	// fmt.Println(reflect.TypeOf(t))
	// fmt.Println(t.TaskID)
	// fmt.Println(t)
	// tl.CompleteTask(t)

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

					for i, task := range tl.Tasks() {
						fmt.Println(i, task.TaskId, task.Name, task.Notes, task.Complete)
					}

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
