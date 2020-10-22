package main

import (
	"fmt"
	"log"
	"os"

	"task/tasklist"

	"github.com/urfave/cli/v2"
)

func main() {
	tl := tasklist.NewTasklist()
	defer tl.Db.Close()

	app := &cli.App{
		Name:  "task",
		Usage: "manage tasks from the terminal",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all tasks",
				Action: func(c *cli.Context) error {
					for i, task := range tl.Tasks() {
						fmt.Println(i+1, task)
					}
					return nil
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a task to the tasklist",
				Action: func(c *cli.Context) error {
					name := c.Args().Get(0)
					tl.AddTask(tasklist.NewTask(name))
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"r", "rm"},
				Usage:   "remove a task from the list",
				Action: func(c *cli.Context) error {
					// TODO optional argument to select task
					tl.RemoveTask(tl.SelectTask())
					return nil
				},
			},
			{
				Name:    "complete",
				Aliases: []string{"c"},
				Usage:   "mark a task as 'completed'",
				Action: func(c *cli.Context) error {
					tl.CompleteTask(tl.SelectTask())
					return nil
				},
			},
			{
				Name:    "uncomplete",
				Aliases: []string{"C"},
				Usage:   "mark a task as `incomplete`",
				Action: func(c *cli.Context) error {
					tl.UncompleteTask(tl.SelectTask())
					return nil
				},
			},
			{
				Name:    "toggle",
				Aliases: []string{"t"},
				Usage:   "toggle the completion state of a task",
				Action: func(c *cli.Context) error {
					tl.ToggleComplete(tl.SelectTask())
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
