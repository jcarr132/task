package main

import (
	// "errors"
	"fmt"
	"log"
	"os"
	"strconv"

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
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "priority",
						Aliases: []string{"p"},
					},
				},
				Action: func(c *cli.Context) error {
					for i, task := range tl.Tasks() {
						if c.Bool("priority") {
							fmt.Println(i+1, task, task.Priority)
						} else {
							fmt.Println(i+1, task)
						}
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
						return err
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
				Flags: []cli.Flag{
					&TargetFlag,
				},
				Action: func(c *cli.Context) error {
					tl.RemoveTask(tl.SelectTask(c.Int("target")))
					return nil
				},
			},
			{
				Name:    "complete",
				Aliases: []string{"c"},
				Usage:   "mark a task as 'completed'",
				Flags: []cli.Flag{
					&TargetFlag,
				},
				Action: func(c *cli.Context) error {
					tl.CompleteTask(tl.SelectTask(c.Int("target")))
					return nil
				},
			},
			{
				Name:    "uncomplete",
				Aliases: []string{"C"},
				Usage:   "mark a task as 'incomplete'",
				Flags: []cli.Flag{
					&TargetFlag,
				},
				Action: func(c *cli.Context) error {
					tl.UncompleteTask(tl.SelectTask(c.Int("target")))
					return nil
				},
			},
			{
				Name:    "toggle",
				Aliases: []string{"tog"},
				Usage:   "toggle the completion state of a task",
				Flags: []cli.Flag{
					&TargetFlag,
				},
				Action: func(c *cli.Context) error {
					tl.ToggleComplete(tl.SelectTask(c.Int("target")))
					return nil
				},
			},
			{
				Name:    "priority",
				Aliases: []string{"p"},
				Usage:   "set the priority level of a task",
				Flags: []cli.Flag{
					&TargetFlag,
				},
				Action: func(c *cli.Context) error {
					newVal, err := strconv.Atoi(c.Args().First())
					if err != nil {
						return err
					}
					task := tl.SelectTask(c.Int("target"))
					task.Priority = newVal
					tl.UpdateTask(task)
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

/* TargetFlag is used by several commands to specify a Task from the TaskList
* to act on. */
var TargetFlag = cli.IntFlag{
	Name:    "target",
	Aliases: []string{"t"},
	Usage:   "the task to act on (use 'task list' to get value)",
}
