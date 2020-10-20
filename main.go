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

	app := &cli.App{
		Name:  "task",
		Usage: "manage tasks from the terminal",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all tasks",
				Action: func(c *cli.Context) error {
					for _, task := range tl.Tasks() {
						fmt.Println(task)
					}
					tl.Db.Close()
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
					tl.Db.Close()
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{"r", "rm"},
				Usage:   "remove a task from the list",
				Action: func(c *cli.Context) error {
					tl.RemoveTask(tl.SelectTask())
					tl.Db.Close()
					return nil
				},
			},
			{
				Name:    "complete",
				Aliases: []string{"c"},
				Usage:   "mark a task as 'completed'",
				Action: func(c *cli.Context) error {
					tl.CompleteTask(tl.SelectTask())
					tl.Db.Close()
					return nil
				},
			},
			{
				Name:    "uncomplete",
				Aliases: []string{"C"},
				Usage:   "mark a task as `incomplete`",
				Action: func(c *cli.Context) error {
					tl.UncompleteTask(tl.SelectTask())
					tl.Db.Close()
					return nil
				},
			},
			{
				Name:    "toggle",
				Aliases: []string{"t"},
				Usage:   "toggle the completion state of a task",
				Action: func(c *cli.Context) error {
					tl.ToggleComplete(tl.SelectTask())
					tl.Db.Close()
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
