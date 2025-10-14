# 🌐 PulseHTTP — Minimal HTTP/1.1 Server on Raw TCP
> *A handcrafted HTTP server written in Go, built from the ground up using raw TCP sockets. No frameworks. No `net/http`. Just protocols, sockets, and packets.*

**Author:** Kartik Saxena

---

## ✨ Why I Built This

Modern developers rely on abstractions like `net/http`, but rarely see what happens underneath.
**PulseHTTP** was my attempt to strip away the layers and **implement a functional HTTP/1.1 server** myself — to understand:

* How HTTP requests are framed over TCP
* How headers are parsed and structured
* How to handle persistent connections and multiple requests

This project is fully functional, unit tested, and ready to extend.

---

## ⚡ Features

* 🧾 Parses HTTP request lines and headers from raw TCP streams
* 🔁 Handles multiple requests and persistent connections
* 🧠 Custom header parsing logic (multi-value support)
* 🧪 Unit-tested request and header parsing modules
* 🧭 Modular internal packages for extensibility

---

## 🧰 Tech Stack

* **Language:** Go (1.18+)
* **Protocol:** HTTP/1.1
* **Transport:** Raw TCP sockets
* **Testing:** Go’s built-in test framework

---

## 🏗️ Architecture

```
+-------------+
|   Client    |
| (curl, etc) |
+------+------+
       |
       v
+-------------+       +----------------+
| TCP Listener|-----> | Request Parser |
+-------------+       +----------------+
                            |
                            v
                     +--------------+
                     | Header Logic |
                     +--------------+
                            |
                            v
                     +---------------+
                     | Future Router |
                     +---------------+
```

* `cmd/tcpListener`: TCP server entry point
* `internal/request`: Request parsing logic
* `internal/headers`: Header parsing logic
* `cmd/udpSender`: Optional UDP tool (for experiments)

---

## 🧪 Running the Server

```bash
go run cmd/tcpListener/main.go
```

**Default Port:** `42069`

Then send a request:

```bash
curl -v http://localhost:42069/coffee
```

Or:

```powershell
Invoke-WebRequest -Uri "http://localhost:42069/coffee"
```

The server prints parsed request lines and headers to the console.

---

## 🧠 Example Output

```
> Accepting connection on 127.0.0.1:42069
> Parsed Request:
  Method: GET
  Path: /coffee
  Version: HTTP/1.1
  Headers:
    Host: localhost:42069
    User-Agent: curl/8.0.1
    Accept: */*
```

---

## 🧪 Tests

```bash
go test ./internal/headers
go test ./internal/request
```

Tests cover:

* Header parsing edge cases
* Multi-value support
* Request line validation

---

## 🧭 Roadmap

* [ ] Add support for request bodies
* [ ] Implement proper HTTP response generation
* [ ] Support chunked transfer encoding
* [ ] Add simple routing
* [ ] Benchmark against Go’s standard `net/http`

---

## 🌱 What I Learned

* How TCP streams are read and buffered
* How HTTP/1.1 parsing works under the hood
* How headers and CRLF terminations are structured
* Why abstractions like `net/http` are so powerful

---

## 🪄 Potential Extensions

This project can be evolved into:

* A lightweight, minimal HTTP framework
* A custom reverse proxy or router
* An educational tool for teaching protocols
* A base for experimenting with HTTP/2 upgrades

---
