package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

// FormatPattern defines the structure of the system logs.
const FormatPattern = "%s - - [%s] \"%s %d %d\" %f\n"

var dirPath = flag.String("d", "./", "Directory path where the server will run")
var serverPort = flag.String("p", "8080", "Port number where the server will run")

func main() {
	flag.Usage = func() {
		fmt.Println("File Server")
		fmt.Println("https://cixtor.com/")
		fmt.Println("https://github.com/cixtor/fileserver")
		fmt.Println("https://en.wikipedia.org/wiki/File_server")
		fmt.Println()
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}

	flag.Parse()

	finfo, err := os.Stat(*dirPath)

	if err != nil {
		fmt.Println("Directory does not exists:", *dirPath)
		os.Exit(1)
	}

	if !finfo.IsDir() {
		fmt.Println("Serving individual files is not possible")
		os.Exit(1)
	}

	if _, err := strconv.Atoi(*serverPort); err != nil {
		fmt.Println("Invalid port number, use one over 1024")
		os.Exit(1)
	}

	fmt.Printf("File Server\n")
	fmt.Printf("Document root: %s\n", *dirPath)
	fmt.Printf("Listening on.: http://0.0.0.0:%s/\n", *serverPort)
	fmt.Printf("Started at...: %s\n", time.Now().Format(time.RFC850))
	fmt.Printf("Press Ctrl-C to quit\n")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Printf("\nServer stopped\n")
			os.Exit(0)
		}
	}()

	mux := http.DefaultServeMux

	mux.Handle("/", http.FileServer(http.Dir(*dirPath)))

	logging := NewLoggingHandler(mux, os.Stderr)
	server := &http.Server{Addr: ":" + *serverPort, Handler: logging}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
