# singleflight
This is a sample code to demonastrate how we can use singleflight to implement request coalescing in golang. Here we coalescing the requests to sqlite DB and compare the performance.
##Server
The server uses sqlite as backend database. The code creates a user tables and populates it with dummy data. There are two apis v1 and v2. The v1 api is implemented strain=ghtforward and makes a database query for each api call. The v2 api uses golang singleflight package and it makes sure that only one execution is in-flight for a given id at a time, reducing the number of database queries.

```
go run server.go
```
##Client
This is a simple http client that makes 10 requets to both v1 and v2 endpoints and records the execution time. It uses go routines to make concurrent requests.

```go run client.go
```
