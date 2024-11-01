# Buy Better Notification System

## Overview
The Buy Better Notification System utilizes the Server-Sent Events (SSE) protocol, allowing clients to connect 
and receive continuous notifications from the Buy Better System. It employs NAT Jetstream, SSE, channels, 
and mutexes to facilitate the broadcasting of client connections to the system.

Note: Currently, there are no publishers integrated into the system anywhere else yet.

## Project Structure
```
├── Makefile
├── cmd
│   └── api
│       └── main.go      # API and Consumer logic
├── go.mod
├── go.sum
├── internal
│   ├── category
│   │   └── category.go  # SSE logic for category notification
│   └── state
│       └── state.go     # Server state for client connections
├── publisher.go         # Publisher mock up
└── readme.md
```

## Tech stack
- NATs jetstream for Pub/Sub
- Gin
- Goroutine (Mutex, Channel)
- Server Sent Event

