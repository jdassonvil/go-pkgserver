Go powered http file server

## Installation
The server is writing to /resources (this might be made configurable).

## Usage
Usage is as simple as it could be.

Use an `HTTP GET` to retrieve a file, if missing you will end up with an `HTTP 404` error code.
If the target path is a directory the server will answer you a JSON listing of the contents.

```
$ curl -X GET http://pkgserver.go/path/filename
```

Use an `HTTP POST` to store a file at the selected path. You can't override an existing file.

```
$ curl -X POST --data-binary "@myfile" http://pkgserver.go/path/filename
```
