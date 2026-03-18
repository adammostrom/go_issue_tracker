# TODO

[/] - Implement json convertion functions in handler
    [ ] - function: backend struct to json response
    [ ] - 

[ ] - Implement parsing/check for the create issue handler/request
[X] - Fix column internal_id "no such column" for a GET/id request
    [/] - Remove the panic if server cant find a specific ID (app shouldnt crash)
            UPDATE: Was SQL error, cant prevent crash

[ ] - Implement function for deleting a value
    [ ] - In router
    [ ] - In server
    [ ] - In database_layer

# LOG

## 2026-03-16
Have the response JSON transfer the basic struct information, keep the logs out and have them be sent via another GET request like:

/issues/{id}/logs


## 2026-03-17

Downloaded and switched to sqlite for now.

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