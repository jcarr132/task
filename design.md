# CLI Task Management Program

A minimal CLI for tracking and managing tasks.


## Requirements
- Written in Go
- Command (CLI) and Interactive (TUI) modes
- Text file database
- Import external calendars (e.g., Google)
- Multiple independent lists of tasks
- Support recurrence
  - Advanced recurrence specification (e.g., last Friday of every month)
- Multiple high-level views (calendar, agenda, task list view)
- Archive for completed items
- Subtasks


## Optional Features
- Local notification daemon
- Push notifications across devices


## Task Schema
- Simple data serialization format (YAML or TOML)
- Fields:
  - ID
  - Name or short description
  - Deadline OR Timeslot OR Untimed
  - Tags
  - Priority
  - Notes
  - Subtasks


## First Prototype
The initial prototype should establish the basic functionality of managing
tasks from the command line. TUI features are extraneous at this point. The
emphasis is on developing a Task List API which will serve both the CLI and
TUI interfaces.
