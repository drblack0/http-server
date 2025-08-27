# HTTP Server (Raw TCP)

A minimal HTTP/1.1 server implementation in Go, built directly on top of raw TCP sockets. This project demonstrates how to parse HTTP requests, manage headers, and handle connections without relying on Go's standard `net/http` package.

## Features

- Parses HTTP request lines and headers from raw TCP streams
- Handles multiple requests and connections
- Custom header parsing logic
- Easily extensible for request bodies and response handling

## Project Structure

```
http-server/
├── cmd/
│   ├── tcpListener/      # TCP server entrypoint
│   │   └── main.go
│   └── udpSender/        # UDP sender utility (optional)
│       └── main.go
├── internal/
│   ├── headers/          # HTTP header parsing logic
│   │   ├── headers.go
│   │   └── headers_test.go
│   └── request/          # HTTP request parsing logic
│       ├── request.go
│       └── request_test.go
├── main.go               # Example usage / entrypoint
├── go.mod
├── go.sum
└── rawget.http           # Example HTTP request
```

## Getting Started

### Prerequisites

- Go 1.18 or newer

### Running the TCP Listener

Start the server:

```bash
go run cmd/tcpListener/main.go
```

The server listens on port `42069` by default.

### Sending Requests

You can use `curl`, Postman, or PowerShell to send HTTP requests:

```bash
curl -v http://localhost:42069/coffee
```

Or with PowerShell:

```powershell
Invoke-WebRequest -Uri "http://localhost:42069/coffee" `
  -Method POST `
  -Headers @{ "Content-Type" = "application/json" } `
  -Body '{"flavor":"dark mode"}'
```

### Example Output

The server will print parsed request lines and headers to the console.

## Testing

Run unit tests for header and request parsing:

```bash
go test ./internal/headers
go test ./internal/request
```

## How It Works

- **TCP Listener:** Accepts raw TCP connections and reads incoming data.
- **Request Parser:** Parses the HTTP request line and headers using a state machine.
- **Header Parser:** Handles header field validation and supports multi-value headers.

## Extending

- Add support for HTTP request bodies (e.g., POST data)
- Implement HTTP response generation
- Add routing and handler logic


**Author:** Kartik Saxena
