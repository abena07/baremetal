
# baremetal

a concurrent tcp server with a custom wire protocol

## why
the sole aim of this project is to understand what powers http servers at the transport layer

## blog
read the write-up: [baremetal: building a tcp server](https://www.bennett-eghan.com/blog/baremetal)

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
