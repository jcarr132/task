package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"task/tasklist"

	"github.com/araddon/dateparse"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
)

func main() {
	var tl tasklist.TaskList
	var err error
	var dbpath string

	app := &cli.App{
		Name:    "task",
		Version: "v0.0.1",
		Authors: []*cli.Author{
			{
				Name:  "JK Carr",
				Email: "joshcarr132@gmail.com",
			},
		},
		HelpName: "task",
		Usage:    "manage tasks from the terminal",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dbpath",
				Aliases: []string{"db"},
				Usage:   "set the location of the task database (default $HOME/.taskdb)",
			},
		},
		Before: func(c *cli.Context) error {
			fmt.Printf("Task CLI - %s\n", c.App.Version)

			if c.IsSet("dbpath") {
				dbpath, err = homedir.Expand(c.String("dbpath"))
			} else {
				dbpath, err = homedir.Dir()
				dbpath = path.Join(dbpath, ".taskdb")
			}

			tl, err = tasklist.NewTasklist(dbpath)
			LogFatalIfErr(err)
			fmt.Printf("Opened task database at %s\n\n", dbpath)
			return nil
		},
		After: func(c *cli.Context) error {
			err = tl.Db.Close()
			LogFatalIfErr(err)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list all tasks",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "sort",
						Aliases: []string{"s"},
						Usage:   "sort the task list based on `VAL` = [priority|created|deadline]",
					},
				},
				Action: func(c *cli.Context) error {
					tl.PrintTasks()
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
						Usage:   "set the priority `VALUE` for the task",
					},
					&cli.StringFlag{
						Name:    "deadline",
						Value:   "nil",
						Aliases: []string{"d"},
						Usage:   "assign a `DEADLINE` for the task",
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
					if c.IsSet("deadline") {
						deadline, err := dateparse.ParseAny(c.String("deadline"))
						LogFatalIfErr(err)
						task.Deadline = deadline
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

	err = app.Run(os.Args)
	LogFatalIfErr(err)
}

/* TargetFlag is used by several commands to specify a Task from the TaskList
* to act on. */
var TargetFlag = cli.IntFlag{
	Name:    "target",
	Aliases: []string{"t"},
	Usage:   "the task to act on (use 'task list' to get value)",
}

func LogFatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
