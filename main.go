package main

import (
	"database/sql"
	"flag"
	"fmt"
	"issuetracker/internal/cli"
	database "issuetracker/internal/database"
	router "issuetracker/internal/router"
	services "issuetracker/internal/services"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Only does initialization: database connection, routing, and starting the HTTP server.
func main() {

	var db *sql.DB
	var err error

	db, err = database.Open()
	if err != nil {
		fmt.Println(err)
		return

	}
	defer db.Close()

	// Initate and load the database schema
	// Initiate the connection to the database (read/write) -> layer between service and db
	db_connection := database.NewDatabaseConnection(db)

	// Create a service struct that delegates to the database for read/write
	issueService := services.NewIssueService(db_connection)

	cli := cli.NewCLI(issueService)

	cmds := cli.BuildCommands()

	if len(os.Args) < 2 {

		cli.PrintCommandUsage(cmds)
		return
	}

	switch os.Args[1] {
	case "start":
		startCmd(os.Args[2:], issueService)
	default:
		cli.Run(cmds, os.Args[1:])

	}

}

func printStart() {
	fmt.Println("")
	fmt.Println("Example usage of start flags:")
	fmt.Println("issuetracker start")
	fmt.Println("issuetracker start --port 9090")
	fmt.Println("issuetracker start --bind 0.0.0.0 --port 8080")
	fmt.Println("issuetracker start --base-url https://issues.example.com")
	fmt.Println("")

}

// Start the HTTP server
func startCmd(args []string, issueService *services.IssueService) {

	printStart()

	fs := flag.NewFlagSet("start", flag.ExitOnError)

	bindAddr := fs.String("bind", "127.0.0.1", "Address to bind HTTP server to")
	port := fs.Int("port", 8080, "Port to run HTTP server on")
	baseURL := fs.String("base-url", "", "public base URL (optional)")
	// Create a router to delegate requests to the server

	fs.Parse(args)

	// Build listen address
	addr := fmt.Sprintf("%s:%d", *bindAddr, *port)

	//If base URL not set, derive it

	effectiveBaseURL := *baseURL
	if effectiveBaseURL == "" {
		effectiveBaseURL = fmt.Sprintf("http://%s", addr)
	}

	r := router.NewRouter(issueService)

	// set up the HTTP server
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    addr, // Add configuration if want to change ports etc.
		Handler: mux,
	}

	// Delegates all HTTP requests to /issues* to the IssuesHandler
	mux.HandleFunc("/issues", r.AllRouting)  // for /issues exact (list, create)
	mux.HandleFunc("/issues/", r.AllRouting) // for /issues/{id} (single issue)

	// Serve the HTML frontend
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("Issuetracker server starting")
	log.Printf("Listening on %s", addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
