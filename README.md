# NAMNADA LINK - Telegram Bot for Saving Articles

**NAMNADA LINK** is a minimalist Telegram bot that helps you save
interesting articles, read them later, and manage your personal reading list.
Just send any link - the bot will save it automatically.

---

## Features

-   Save articles by simply sending a link
-   Get a random unread article
-   Mark articles as read
-   Delete articles
-   View all saved articles
-   Clean, fast, minimalistic functionality - nothing extra

---

## Commands

    /start  — welcome message  
    /random — get a random unread article  
    /read   — mark an article as read  
    /remove — delete an article  
    /list   — list all saved articles  
    /help   — show help message  

You can also send any link directly - the bot will save it automatically.

---

## Technical Overview

The bot is written in **Go** from scratch without any third-party Telegram libraries.
It uses a clean implementation of the **Telegram Bot API** and Go's standard library.

---

## Running the Bot

Before starting, set:

-   Telegram bot scheme (http/https)
-   Telegram API host
-   Your bot token
-   `BATCH_SIZE` (how many updates to process at once)

``` bash
BATCH_SIZE=100 go run cmd/main.go -tg-bot-scheme 'https' -tg-bot-host 'api.telegram.org' -tg-bot-token 'your_bot_token'
```

---

## Data Storage

Currently the bot uses **in-memory storage**.

Planned improvements:

-   Database support (PostgreSQL or SQLite)
-   Optional file-based storage

The storage interface already allows easy extensions.

---

## Future Plans

-   Convert articles to forwarded Telegram messages (for channel posts)
-   Add language selection (English / Russian)
-   Add persistent database storage
-   Create a full Docker deployment

---

## Project Structure

    .
    ├── cmd/
    │   └── main.go                    # Application entry point
    │
    ├── pkg/
    │   ├── clients/
    │   │   └── telegram/              # Pure Telegram Bot API client
    │   │       ├── telegram.go        # GET updates, send messages
    │   │       └──types.go            # DTOs for Telegram API
    │   │
    │   ├── consumer/                  # Event processing and concurrency logic
    │   │   ├── consumer.go
    │   │   └── event-consumer/        # Concurrent event consumer with error threshold
    │   │       └── event-consumer.go
    │   │
    │   ├── events/
    │   │   └── telegram/              # Parsing incoming messages and command handling
    │   │       ├── commands.go        # /random, /read, /remove, etc.
    │   │       ├── messages.go        # Bot message templates
    │   │       └──telegram.go         # Event transformation to internal types
    │   │
    │   └── storage/
    │       ├── storage.go             # Storage interface
    │       └── memory/                # In-memory implementation
    │           └── memory.go
    │
    ├── go.mod
    └── README.md

### Components

#### **Telegram Client**

Low-level HTTP client for calling Telegram Bot API - no external
libraries.

#### **Event Processor**

Handles events, transforms raw Telegram updates into internal commands.

#### **Event Consumer**

Runs handlers concurrently, provides batching, error counting, and
controlled shutdown.

#### **Storage Layer**

Abstract interface with in-memory implementation.
Easily extendable to PostgreSQL, MongoDB, file storage, Redis, etc.

#### **Tests**

Every module has unit tests using `httptest`, table tests, mocks, and
error scenarios.

---

## Feedback

Have suggestions or found an issue?
Feel free to open an issue - the project is open for improvements and contributions.
