Most developers use Go’s standard `net/http` — but *building an HTTP/1.1 server from scratch on raw TCP* shows a deep understanding of **networking, protocols, and lower-level 
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

## 🧑‍💻 Portfolio Angle

This project demonstrates:

* Deep understanding of **HTTP and TCP**
* Ability to **build infra from scratch**
* Clean code & modular Go design
* Curiosity and low-level problem solving

> 🏆 This is a great repo to showcase in interviews or on your personal site / GitHub profile.

---

## 📜 License

MIT © 2025 Kartik Saxena

---

## 💼 Bonus Tip: Portfolio Presentation

When adding this to your **portfolio site or GitHub profile**, present it like:

> **PulseHTTP** — A handcrafted HTTP server built on raw TCP in Go.
> Implemented request parsing, header handling, and persistent connections without any frameworks.
> ✅ Deep protocol understanding · 🧠 Systems thinking · 🧰 Go networking.

👉 You can also add a **short blog post or LinkedIn article** walking through how you parsed the HTTP request line. That’s a big differentiator in interviews.

---

Would you like me to draft a **1-paragraph portfolio description** (something short, for your personal website or LinkedIn project section)?
