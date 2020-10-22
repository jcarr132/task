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
		Name:    "task",
		Version: "v0.0.1",
		Authors: []*cli.Author{
			{
				Name:  "J Carr",
				Email: "joshcarr132@gmail.com",
			},
		},
		HelpName: "task",
		Usage:    "manage tasks from the terminal",
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
				Name:      "add",
				Aliases:   []string{"a"},
				Usage:     "add a task to the tasklist",
				ArgsUsage: "[name: (string)]",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "priority",
						Aliases: []string{"p"},
						Usage:   "set the priority value for the task",
					},
				},
				Action: func(c *cli.Context) error {
					name := c.Args().Get(0)
					task, err := tasklist.NewTask(name)
					if err != nil {
						cli.ShowCommandHelp(c, "add")
						log.Fatal(err)
					}
					if c.Int("priority") != 0 {
						task.Priority = c.Int("priority")
					}
					tl.AddTask(task)
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
				Usage:   "mark a task as 'incomplete'",
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
