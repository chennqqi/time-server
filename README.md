# Time Server

The TCP time server according to [RFC868](https://tools.ietf.org/html/rfc868)

## Installation

Get the package

```bash
go get -v github.com/ilyakaznacheev/time-server
```

Build executables with `make`

```bash
make build
```

or manual

```bash
go build -o time-server
go build -o time-client cmd/client/client.go 
```

## Usage

### Server

Run the server on localhost at the port you want

```bash
./time-server -p 11037
```

### Client

Run the client to call the RFC868 time server via RFC

```bash
./time-client time.nist.gov 37
```
