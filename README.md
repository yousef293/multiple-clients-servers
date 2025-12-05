# multiple-clients-servers

# Real-Time Chat System

A concurrent TCP chat server built with Go, featuring real-time message broadcasting and user notifications.

## Features

-  Real-time message broadcasting to all connected clients
-  Multi-client support with unique user IDs
- Join/leave notifications
- Thread-safe client management with Mutex
-  Non-blocking message delivery using goroutines and channels
-  No self-echo (clients don't receive their own messages)

## Architecture

- **Server**: Manages client connections, broadcasts messages using channels
- **Client**: Connects to server, sends/receives messages concurrently
- **Concurrency**: Uses goroutines for handling multiple clients and channels for communication

## Getting Started

### Prerequisites
- Go 1.16 or higher

### Running the Server
```bash
cd server
go run main.go
