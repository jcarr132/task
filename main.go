package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"task/tasklist"

	"github.com/urfave/cli/v2"
)

func main() {

	tl := tasklist.NewTasklist()
	fmt.Println(reflect.TypeOf(tl.Tasks()[0]))

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
						fmt.Println(i, task.Name, task.Notes)
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
