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

	// reinitialize the database (testing)
	os.Remove("data")
	tl.AddTask(tasklist.NewTask("buy groceries"))
	tl.AddTask(tasklist.NewTask("work on cli task manager program"))
	tl.AddTask(tasklist.NewTask("pay bills"))

	// test completing a task
	t := tl.Tasks()[0]
	tl.CompleteTask(t)

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

					for _, task := range tl.Tasks() {
						fmt.Println(task)
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
