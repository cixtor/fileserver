package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"time"
)

var dirPath = flag.String("d", "./", "Directory path where the server will run")
var serverPort = flag.String("p", "8080", "Port number where the server will run")

func main() {
	flag.Usage = func() {
		fmt.Println("File Server")
		fmt.Println("  http://cixtor.com/")
		fmt.Println("  https://github.com/cixtor/fileserver")
		fmt.Println("  http://en.wikipedia.org/wiki/File_server")
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}

	flag.Parse()

	finfo, err := os.Stat(*dirPath)

	if err != nil {
		flag.Usage()
		fmt.Printf("\nDirectory does not exists: %s\n", *dirPath)
		os.Exit(1)
	}

	if !finfo.IsDir() {
		flag.Usage()
		fmt.Printf("\nServing individual files is not allowed\n")
		os.Exit(1)
	}

	re := regexp.MustCompile(`^[0-9]{2,4}$`)
	match := re.FindStringSubmatch(*serverPort)

	if match == nil {
		flag.Usage()
		fmt.Printf("\nError. Invalid port number\n")
		os.Exit(1)
	}

	fmt.Printf("File Server\n")
	fmt.Printf("Document root: %s\n", *dirPath)
	fmt.Printf("Listening on.: http://localhost:%s/\n", *serverPort)
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

	http.Handle("/", http.FileServer(http.Dir(*dirPath)))

	if err := http.ListenAndServe(":"+*serverPort, nil); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
