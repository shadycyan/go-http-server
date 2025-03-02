# HTTP server from scratch 

a HTTP server, written in Go, built from the ground up using TCP, that's capable of handling simple GET/POST requests, serving files, handling multiple concurrent connections, and supports `gzip` compression.

## endpoints 

**`/echo/<string>`** 

echoes back the string provided in the url. if the client supports gzip (via accept-encoding: gzip), it'll compress the response for that extra edge.
 

**`/user-agent`** 

returns your user-agent header.


**`/files/`**  
  - `POST /files/<filename>` writes request body to a file in the specified directory.
  - `GET /files/<filename>` retrieves file data. 


## usage 
 
1. **build:** 

run `go build ./cmd/server` in the project directory.
 
2. **run:**

launch the server with a directory for file ops. for example, if you wanna serve files from a folder called static, create that folder and run:

```bash
./server ./static
```
the server listens on `0.0.0.0:42069`.
 
3. **test endpoints:**  
  - **echo:** 
```bash
curl -v -H "Accept-Encoding: gzip" http://localhost:42069/echo/hello
```
  - **user-agent:** 
```bash
curl -v http://localhost:42069/user-agent
```
  - **files:** 
```bash
curl -v -X POST --data "file content" http://localhost:42069/files/test.txt
curl -v http://localhost:42069/files/test.txt
```

## notes 
- it's a simple server, so it's not production-ready.
- error handling is basic—perfect for learning and experimentation.

enjoy tinkering, and feel free to contribute or drop some feedback in issues!
