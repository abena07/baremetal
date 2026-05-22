
# baremetal

a concurrent tcp server with a custom wire protocol, built from scratch in go

## why
the sole aim of this project is to understand what powers http servers at the transport layer

## how to run

```bash
go run .
````

server will start on:

```
localhost:8080
```

you can test it using netcat:

```bash
nc localhost 8080
```