# File Server [![GoReport](https://goreportcard.com/badge/github.com/cixtor/fileserver)](https://goreportcard.com/report/github.com/cixtor/fileserver) [![GoDoc](https://godoc.org/github.com/cixtor/fileserver?status.svg)](https://godoc.org/github.com/cixtor/fileserver)

> In computing, a file server (or fileserver) is a computer attached to a network that has the primary purpose of providing a location for shared disk access, i.e. shared storage of computer files (such as documents, sound files, photographs, movies, images, databases, etc.) that can be accessed by the workstations that are attached to the same computer network. The term server highlights the role of the machine in the client–server scheme, where the clients are the workstations using the storage.
>
> A file server is not intended to perform computational tasks, and does not run programs on behalf of its clients. It is designed primarily to enable the storage and retrieval of data while the computation is carried out by the workstations.
> 
> — http://en.wikipedia.org/wiki/File_server

## Installation

```
go get -u github.com/cixtor/fileserver
```

## Usage

```
$ fileserver -p 5690
File Server
Listening on http://localhost:5690
Started at Wed Jan 30 13:36:26 2019
Document root is /Users/cixtor/public_html
Press Ctrl-C to quit.
[::1] - - [30/Jan/2019:21:36:36 +0000] "GET /index.html HTTP/1.1" 301 0 0.0000s
[::1] - - [30/Jan/2019:21:36:36 +0000] "GET / HTTP/1.1" 200 225 0.0004s
[::1] - - [30/Jan/2019:21:36:40 +0000] "GET /logger.go HTTP/1.1" 200 976 0.0090s
[::1] - - [30/Jan/2019:21:36:48 +0000] "GET /README.md HTTP/1.1" 200 1161 0.0003s
[::1] - - [30/Jan/2019:21:36:50 +0000] "GET /.git/ HTTP/1.1" 200 463 0.0005s
[::1] - - [30/Jan/2019:21:36:53 +0000] "GET /.git/HEAD HTTP/1.1" 200 23 0.0002s
^C
Server stopped
```
