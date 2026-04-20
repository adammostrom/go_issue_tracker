# TODO

[/] - Implement json convertion functions in handler
    [ ] - function: backend struct to json response 

[/] - Implement parsing/check for the create issue handler/request
    2026-03-31: Make basic validation for string length (min/max)
[X] - Fix column internal_id "no such column" for a GET/id request
    [X] - Remove the panic if server cant find a specific ID (app shouldnt crash)
            UPDATE: Was SQL error, cant prevent crash

[ ] - Implement sanitation. In that case, where? for what? All fields?

[X] - Implement function for deleting a value
    [X] - In router
    [X] - In server
    [X] - In database_layer
[X] - Implement GET for logs (issues/id/logs)
    [X] - In router
    [X] - In server
    [X] - In database_layer
[/] - Implement UPDATE for an issue (update fields)
    [/] - In router
    [/] - In server
    [/] - In database_layer
[/] - Implement "add log entry" for an issue
    [X] - In router
    [X] - In server
    [X] - In database_layer

[ ] - Add proper validation for getting an issue that doesnt exist (handle it). 
[ ] - Add a VIEWS in SQL
    [ ] - for GET all
    [ ] - for GET one
    [ ] - for GET logs

Implement LOGs handling for:
[ ] - Handler
    [ ] - CREATE
    [ ] - DELETE ALL
    [ ] - DELETE ONE (id? index? latest? first?)
    [ ] - GET ALL
    [ ] - GET latest/first (filtering)
[ ] - CLI
    [ ] - CREATE
    [ ] - DELETE ALL
    [ ] - DELETE ONE (id? index? latest? first?)
    [ ] - GET ALL
    [ ] - GET latest/first (filtering)

[ ] - Split models into their respective files (like Issue.go)
    [/] - add model verification (ex: validIssue or similar, check: invariants, example: id < 0 return false)

[/] - Make it so that external_ref turns the string into all lowercase before parsing, but all capital for viewing
    [/] - Make restriction on external_ref to certain length, no symbols (only numbers and digits) etc, include in model verification.
    
[X] - change external_ref to string/text instead of int
[X] - Apply logging entries for whenever a PATCH is made. (ex: 2026-03-25: active changed to false)

[/] - Add validation and checking for:
    [X] - Title (min 2, max 30)? -> Max 20
    [X] - Description: Also a max -> Max 30
    [/] - External Ref, min, max, no symbols, only integers and characters

[ ] - Make a "clear log" for a specific issue.

[/] - Implement logic to filter issues based on "active", "time created" etc.
    [X] - Implement filtering for active
    [X] - Implement filtering for deacrive
    [ ] - Implement filtering for time created (age) 

[ ] - Implement sorting
    [ ] - Implement sorting by title (ascending/descending)
    [ ] - Implement sorting by external ref (asc/desc)
    [ ] - Implement sorting by time created (date generally)
[ ] - Implement functionality for erasing log entries 

[X] - Refactor CLI handling and cli command functions -> Command structure /build function in cli.go / 2026-04-14
[ ] - Refactor printout in CLI for issues found

[ ] - Add some type of "started/not started" or find a way how to make the active/deactive reflect 3 states: not started, started, finished 

[ ] - Update README on how to build (using go install) to create into user path, so it can be run globally (binary is placed in local/usr/bin)
[ ] - Add "INIT" for initiating the issuetracker where it creates a folder ".issuetracker/" and the db file inside it. 
    [ ] - Implement checks if it exists, if user tries to do "issuetracker issues list" and its not been initiated, an error should reflect this
    [ ] - Create "ux" for users failing a command and having to run "issuetracker init"
    [ ] - BASICALLY: HAVE ONE DB INIT FUNCTION, AND ONE DB RUN FUNCTION. Have check for existing DB in the run, and if non existing, return fatal err.


# LOG

## 2026-04-21
Added setting inactive and active in CLI (for router this is done with patch request). Added a getAllIssues test, small fixes and tweaks
Also added "create log entry" in router/handler

## 2026-04-15

Maybe add subtasks as checkboxes like how I write it here? Typ som att en issue har "parts" eller "components" som bara har en titel, och en "finished/unfinished", typ:

title: Implement UI
ExternalREF: IMPL_UI
Description: Implement the UI for all functions
Active: True
Components: {
    [] Implement handling
    [] Implement error correction
}


## 2026-04-14 
Implemented the CLI commands tree/struct in cli.go, chose to go the most simple route and keep it slightly more verbose for readability over practically.

## 2026-04-03
Implement a CLI endpoint as an alternative to the routing/handling HTTP endpoint. Do not use any HTTP here:
- Simple CLI
- Have it implement a IssueServer interface (issueServer)
- Handle its own errors

## 2026-04-01

When trying to create issue with external ref that already exists, no information/error print is specifying that. It just silently ignores the query. FIX.

Also, fetching logs return empty strings

## 2026-03-25
Start writing tests

## 2026-03-24
Updated the service to implement a database interface instead, purpose -> swappable with mockDB for testing

## 2026-03-21

Basically finish this:
POST   /issues
GET    /issues
GET    /issues/{id}
PATCH  /issues/{id}/resolve
POST   /issues/{id}/logs
GET    /issues/{id}/logs

Also: change external_ref to string/text instead of int
## 2026-03-18

Implemented a switch case in the router to catch all requests, check alternative imrpovement:
parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

switch len(parts) 
case 1:
    // /issues → list or create
case 2:
    // /issues/2 → single issue
case 3:
    // /issues/2/log → logs


So, fix the panic (so the program doesnt crash if id doesnt exist) and then fix the column internal_id issue

## 2026-03-17

Downloaded and switched to sqlite for now.


## 2026-03-16
Have the response JSON transfer the basic struct information, keep the logs out and have them be sent via another GET request like:

/issues/{id}/logs




