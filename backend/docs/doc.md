# Backend Architecture

## Overview

  Sepolia (blockchain)
        ↓ WebSocket (events)
    Listener
        ↓ channel (AuctionCreatedEvent)
    Worker Pool
        ↓
    Service
        ↓
    Repository
        ↓
    PostgreSQL
        ↑
    HTTP Handlers
        ↑
    Web App

## Goroutines

| Goroutine               | Role                                        |
| ----------------------- | ------------------------------------------- |
| main                    | Waits for OS signal, orchestrates shutdown  |
| HTTP Server             | Handles incoming REST requests              |
| Blockchain Listener     | Subscribes to contract events via WebSocket |
| Worker Pool (N workers) | Consumes events and persists to database    |

## Shutdown Flow

CTRL+C
→ cancel() propagates to all goroutines via context
→ listener stops WebSocket subscription
→ workers drain remaining events and exit
→ HTTP server finishes in-flight requests
→ process exits cleanly

## Communication

- Goroutines communicate via channels, never shared memory
- `listener.Events` — buffered channel (size 100), listener → workers
- All database calls propagate context for cancellation

## Layers

handler → receives HTTP request, calls service
service → business logic, calls repository
repository → database queries (database/sql, no ORM)

## Key Packages

| Package               | Role                                                 |
| --------------------- | ---------------------------------------------------- |
| `internal/auction`    | Models, repository, service, handlers                |
| `internal/blockchain` | Event listener, ABI parsing                          |
| `internal/database`   | Connection pool setup, migrations                    |
| `cmd/api`             | Entry point, dependency injection, graceful shutdown |
