# CLI Task Management Program

A minimal CLI for tracking and managing tasks.


## Requirements
- Written in Go
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


## Considerations
- Centralized vs decentralized sync
  - Centralized: App (and DB) live in a server, clients to connect. Everything
    syncs to the server.
  - Decentralized: Each instance of the app is independent, maintains it's own
    database. Databases sync with each other.

## Task Schema
- Simple data serialization format (YAML or TOML)
- Fields:
  - Name or short description
  - Deadline OR Timeslot OR Untimed
  - Tags
  - Priority
  - Notes
  - Subtasks
