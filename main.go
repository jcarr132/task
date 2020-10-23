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
	tl, err := tasklist.NewTasklist()
	if err != nil {
		log.Fatal(err)
	}
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
						Usage:   "include the priority value for each task",
						Aliases: []string{"p"},
					},
				},
				Action: func(c *cli.Context) error {
					tasks, err := tl.Tasks()
					if err != nil {
						return err
					}
					for i, task := range tasks {
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
					return tl.AddTask(task)
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
					task, err := tl.SelectTask(c.Int("target"))
					if err != nil {
						return err
					}
					return tl.RemoveTask(task)
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
					task, err := tl.SelectTask(c.Int("target"))
					if err != nil {
						return err
					}

					return tl.CompleteTask(task)
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
					task, err := tl.SelectTask(c.Int("target"))
					if err != nil {
						return err
					}
					return tl.UncompleteTask(task)
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
					task, err := tl.SelectTask(c.Int("target"))
					if err != nil {
						return err
					}
					return tl.ToggleComplete(task)
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
					task, err := tl.SelectTask(c.Int("target"))
					if err != nil {
						return err
					}
					task.Priority = newVal
					return tl.UpdateTask(task)
				},
			},
			{
				// Next command
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
