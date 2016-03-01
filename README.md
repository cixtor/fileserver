### File Server

> In computing, a file server (or fileserver) is a computer attached to a network that has the primary purpose of providing a location for shared disk access, i.e. shared storage of computer files (such as documents, sound files, photographs, movies, images, databases, etc.) that can be accessed by the workstations that are attached to the same computer network. The term server highlights the role of the machine in the client–server scheme, where the clients are the workstations using the storage.
>
> A file server is not intended to perform computational tasks, and does not run programs on behalf of its clients. It is designed primarily to enable the storage and retrieval of data while the computation is carried out by the workstations.
> 
> — http://en.wikipedia.org/wiki/File_server

### Installation

I do not distribute binaries for security reasons.

```
go get -u github.com/cixtor/fileserver
```

### Usage

```
$ fileserver -p 5690
  File Server
  Document root: /Users/cixtor/public_html
  Listening on.: http://0.0.0.0:5690/
  Started at...: Tuesday, 10-Oct-13 17:33:25 EST
  Press Ctrl-C to quit
^C
Server stopped
```
