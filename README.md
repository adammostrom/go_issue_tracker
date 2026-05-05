# Issuetracker

Issuetracker is a small and simple issue-tracking tool written in Go.
It can be run either as a command-line interface:

    issuetracker <command>

or as a backend service for a web interface:

    issuetracker start

Issuetracker is intended to be used as a **project-local dependency** for your project.

It uses SQLite for storage and requires no external services.

---

## Project Layout

Issuetracker is distributed as a single binary.
When placed inside a project, it is expected to live alongside its own data directory:

```
your_project/
└── external/
    └── issuetracker/
        ├── issuetracker     
        └── .issuetracker/
            └── issuedb.db
```

## Installation & Running

1. Download the appropriate binary for your platform.
2. Place it anywhere inside the project you want to track issues for.
3. Make it executable if necessary: 
`chmod +x issuetracker`
4. Run it: 
`./issuetracker`

## Data Storage

On first run, Issuetracker automatically creates a hidden data directory
next to the executable:
```
.issuetracker/
└── issuedb.db
```

## Web Interface

Run the start command with optional flags:
```
issuetracker start
issuetracker start --port 9090
issuetracker start --bind 0.0.0.0 --port 8080
issuetracker start --base-url https://issues.example.com
```

## Command line interface

Run issuetracker with various commands, a few examples:

`./issuetracker list`

Returns

```
ID   ST   CREATED           EXT REF      TITLE
──── ──── ───────────────── ──────────── ─────────────────────────────
 1    [ ]  2026-05-05 14:41  VAL22         Validate authentication handler
```


`./issuetracker getref VAL22`

```
Issue #1
──────────────────────────────────────────
Title:        Validate authentication handler
Description:  Validate the authentication handler and make sure to write some tests for it.
External Ref: VAL22

Active:       true
Progress:     idle

Logs:
  • 2026-05-05 14:41  Issue created
```

`./issuetracker set progress 1 started `

```
ID   ST   CREATED           EXT REF      TITLE
──── ──── ───────────────── ──────────── ─────────────────────────────
 1    [/]  2026-05-05 14:41  VAL22         Validate authentication handler


Issue #1
──────────────────────────────────────────
Title:        Validate authentication handler
Description:  Validate the authentication handler and make sure to write some tests for it.
External Ref: VAL22

Active:       true
Progress:     idle

Logs:
  • 2026-05-05 14:41  Issue created
  • 2026-05-05-15:30  Progress changed to: started

```


