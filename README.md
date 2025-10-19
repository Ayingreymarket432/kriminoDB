# kriminoDB — A Simple Distributed Key-Value Store (Redis-Inspired)

kriminoDB is a minimal, in-memory key-value store built in Go, inspired by Redis.

> **Current Status**: Single-node TCP server supporting `SET` and `GET` commands.
> **Next**: Multi-node clustering, replication, and Raft consensus.

---

## Features (current)

- In-memory key-value storage
- Concurrent-safe access (using `sync.RWMutex`)
- TCP-based command interface
- Supports values with spaces (e.g., `SET name "Ayoub Krimi"`)

---

## Quick Start

### Prerequisites

- Go 1.20+

### Run the Server

```bash
make run
```

By default, the server listens on `localhost:3000`.

### Test with `telnet`

Open a new terminal and connect:

```bash
telnet localhost 3000
```

Then try these commands:

```text
SET name Ayoub
GET name
SET full_name "Ayoub Krimi"
GET full_name
```

Expected responses:

```text
OK
name=Ayoub
OK
full_name=Ayoub Krimi
```

> **Note**: Commands are case-insensitive (`set`, `SET`, `Set` all work).

---

## Command Reference

| Command | Usage               | Description                  |
| ------- | ------------------- | ---------------------------- |
| `SET`   | `SET <key> <value>` | Store a key-value pair       |
| `GET`   | `GET <key>`         | Retrieve the value for a key |

- Keys and values are treated as strings.
- Values **can contain spaces** (everything after the key is treated as the value).
- Invalid commands return an error message.

---

## Built With

- **Go** (goroutines, `net`, `sync`, `bufio`)
- **TCP sockets** for client communication
- **Mutexes** for thread-safe data access

---

## Author

Ayoub Krimi — Final-year Computer Science Student @ ISI Ariana, Tunisia
Open to collaboration, feedback, and distributed systems challenges!
