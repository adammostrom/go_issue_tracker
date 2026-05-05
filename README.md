Issuetracker

Issuetracker is written entirely in Go and is designed to be simple to use. It can run both via a command line interface via the command `issuetracker <command>` or as a backend service for a web page: `ìssuetracker start`. 

The idea is to use it as a dependency within your project, to give a simple overview and tracking of project-related issues. 
Database: SQLite

The idea is to contain it as follows:

your_project/
└── external/
    └── issuetracker/
        ├── issuetracker.bin          
        └── .issuetracker/
            └── issuetracker_sqlite3.db

## How to run
Make sure Go is installed.
Download the repo into the project or directory of your choice



## Data Storage

issuetracker is distributed as a single binary.

On first run, it creates the following directory next to the executable:

.issuetracker/
└── issuedb.db

This SQLite database stores all project issues.
The directory is created automatically and should not be committed to version control
