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

var dir_path = flag.String("path", "./", "Set the directory path where the server will run")
var server_port = flag.String("port", "8080", "Set the port number where the server will run")

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

	finfo, err := os.Stat(*dir_path)
	if err != nil {
		flag.Usage()
		fmt.Printf("\nDirectory does not exists: %s\n", *dir_path)
		os.Exit(1)
	}

	if !finfo.IsDir() {
		flag.Usage()
		fmt.Printf("\nServing individual files is not allowed\n")
		os.Exit(1)
	}

	port_re := regexp.MustCompile(`^[0-9]{2,4}$`)
	var port_match []string = port_re.FindStringSubmatch(*server_port)

	if port_match == nil {
		flag.Usage()
		fmt.Printf("\nError. Invalid port number\n")
		os.Exit(1)
	}

	fmt.Printf("File Server\n")
	fmt.Printf("Document root: %s\n", *dir_path)
	fmt.Printf("Listening on.: http://localhost:%s/\n", *server_port)
	fmt.Printf("Started at...: %s\n", time.Now().Format(time.RFC850))
	fmt.Printf("Press Ctrl-C to quit\n")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Printf("\nServer stopped\n")
			os.Exit(0)
		}
	}()

	http.Handle("/", http.FileServer(http.Dir(*dir_path)))
	err = http.ListenAndServe(":"+*server_port, nil)

	if err != nil {
		flag.Usage()
		fmt.Println()
		log.Fatal(err)
	}
}
