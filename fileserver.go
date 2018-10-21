package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"
)

var directory string
var serverPort string

func init() {
	flag.StringVar(&directory, "d", "./", "Directory path where the server will run")
	flag.StringVar(&serverPort, "p", "8080", "Port number where the server will run")

	flag.Usage = func() {
		fmt.Println("File Server")
		fmt.Println("https://cixtor.com/")
		fmt.Println("https://github.com/cixtor/fileserver")
		fmt.Println("https://en.wikipedia.org/wiki/File_server")
		fmt.Println()
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	var err error
	var finfo os.FileInfo
	var abspath string

	if finfo, err = os.Stat(directory); err != nil {
		fmt.Printf("cannot access directory `%s`\n", directory)
		os.Exit(1)
	}

	if !finfo.IsDir() {
		fmt.Println("file cannot be used as a directory")
		os.Exit(1)
	}

	if _, err = strconv.Atoi(serverPort); err != nil {
		fmt.Printf("cannot start server on port `:%s`\n", serverPort)
		os.Exit(1)
	}

	if abspath, err = filepath.Abs(directory); err != nil {
		fmt.Println("filepath.Abs", err)
		os.Exit(1)
	}

	fmt.Printf("File Server\n")
	fmt.Printf("Listening on http://0.0.0.0:%s\n", serverPort)
	fmt.Printf("Started at %s\n", time.Now().Format(time.ANSIC))
	fmt.Printf("Document root is %s\n", abspath)
	fmt.Printf("Press Ctrl-C to quit.\n")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Printf("\nServer stopped\n")
			os.Exit(0)
		}
	}()

	mux := http.DefaultServeMux
	mog := NewLogger(mux, os.Stderr)
	mux.Handle("/", http.FileServer(http.Dir(directory)))
	server := &http.Server{Addr: ":" + serverPort, Handler: mog}
	log.Fatal(server.ListenAndServe())
}
