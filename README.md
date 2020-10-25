# Task CLI

Task is a minimal terminal-based program for task/time management written in
Go. It is in very early stages of development but has some basic functionality
implemented.


## Dependencies
- Go 1.15 or later is required to build the program


## Installation
```
go get gitlab.com/_jcarr132/task
```

## Usage
```
> task
Task CLI - v0.0.1

Opened task database at /home/jc/.taskdb

NAME:
   task - manage tasks from the terminal

USAGE:
   task [global options] command [command options] [arguments...]

VERSION:
   v0.0.1

AUTHOR:
   J Carr <joshcarr132@gmail.com>

COMMANDS:
   list, l        list all tasks
   add, a         add a task to the tasklist
   remove, r, rm  remove a task from the list
   complete, c    mark a task as 'completed'
   uncomplete, C  mark a task as 'incomplete'
   toggle, tog    toggle the completion state of a task
   priority, p    set the priority level of a task

   help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --dbpath value, --db value  set the location of the task database (default $HOME/.taskdb)
   --help, -h                  show help (default: false)
   --version, -v               print the version (default: false)
```
