# Stock Exchange Simulator

A Go-based application designed to simulate a stock exchange's core functionality, with an order matching engine, trade generator, and real-time tick/candle streaming.

## Project Modules

This system is built around three core services. It's a monolith for now, will probably break TickBus into a new service if it gets too heavy with streaming. Don't mind the names, they're random identifiers an LLM generated for me.

### 1. Nexus (Order receiver for now, will potentially become the RMS)
- **Purpose:** Handles incoming buy and sell orders from users. This is the client facing module of the system.
- **Key Features:**
    - Order book management (DB)
    - Order cancellation and modification (TODO)
    - Currently only supports limit orders (TODO: market orders, gtt, etc.).
    - Order validation before placement (RMS, TODO)

### 2. TradeMaster (Order Matcher & Tradebook Generator)
- **Purpose:** Maintains a ledger and generates trades for every matched pair of orders.
- **Key Features:**
    - Order matching engine (with redis sorted sets)
    - Price discovery
    - Partial order fulfillment
    - Trade book management (DB)

### 3. TickBus (The Tick Stream)
- **Purpose:** Manages and streams real-time tick data.
- **Key Features:**
    - Real-time Last Traded Price (LTP) updates per stock.

---

## System design brief

Current wip: Order matching and trade generation (Nexus + TradeMaster).

### Order Matching Workflow

1.  **Order Placement:** A user submits an order (e.g., BUY 10 shares of IDEA at 10 bucks) to the `Nexus` API.
2.  **Add to Redis queue:** `Nexus` validates the order, then enqueues a 'PlaceOrder' task into an Asynq queue and immediately confirms receipt to the user (exchange order id).
3.  **Queue item processed:** A `TradeMaster` worker, running in the background, picks up the 'PlaceOrder' task. The queue is used because if the order matching process (next step) fails, the order will be missed and must wait for another order to close a trade.
4.  **Order Book Matching:** The worker queries the **Order Book** in Redis. It is implemented with a sorted set. For a BUY order, it looks for the best-priced SELL order (lowest price) in the corresponding Redis Sorted Set.
5.  **Trade Execution:**
    *   **If a match is found:** `TradeMaster` generates a `Trade`, saves it to PostgreSQL, and removes the fulfilled order from the Redis sorted set.
    *   **If no match is found:** The new order is added to the Order Book (the appropriate Redis Sorted Set) to await a future match.

## Setup

### Prerequisites

- Go (1.18+)
- PostgreSQL
- Redis

### Installation

Clone the repository and install dependencies:
```sh
git clone <repository-url>
cd stock-exchange-simulator
go mod tidy
```

### Running the Application

1.  Ensure your PostgreSQL and Redis servers are running.
2.  Update the database connection details in `pkg/db/db.go` if they are not the default.
3.  Run the main application:
```sh
go run cmd/stock-exchange-simulator/main.go
```
The application will start, connecting to the database, initializing services, and starting the API server on port `:8080`.

