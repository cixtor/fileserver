package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

var directory string
var serverPort int

func init() {
	flag.StringVar(&directory, "d", "./", "Directory path where the server will run")
	flag.IntVar(&serverPort, "p", 0, "Port number where the server will run")

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

	finfo, err := os.Stat(directory)

	if err != nil {
		fmt.Printf("cannot access directory `%s`\n", directory)
		os.Exit(1)
	}

	if !finfo.IsDir() {
		fmt.Println("file cannot be used as a directory")
		os.Exit(1)
	}

	abspath, err := filepath.Abs(directory)

	if err != nil {
		fmt.Println("filepath.Abs", err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			fmt.Printf("\nServer stopped\n")
			os.Exit(0)
		}
	}()

	mux := http.DefaultServeMux

	mux.Handle("/", http.FileServer(http.Dir(directory)))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: NewLogger(mux, os.Stderr),
	}

	l, err := net.Listen("tcp", server.Addr)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer l.Close()

	addr := strings.Replace(l.Addr().String(), "[::]:", "localhost:", 1)

	fmt.Printf("File Server\n")
	fmt.Printf("Listening on http://%s\n", addr)
	fmt.Printf("Started at %s\n", time.Now().Format(time.ANSIC))
	fmt.Printf("Document root is %s\n", abspath)
	fmt.Printf("Press Ctrl-C to quit.\n")

	if err := server.Serve(l); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
