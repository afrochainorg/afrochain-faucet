# cosmos-rest-faucet

This repository provides a fully functional faucet implementation tailored for Cosmos SDK-based blockchain networks. The faucet allows developers and testers to request tokens from a predefined pool for testing purposes, streamlining the development and testing processes in blockchain environments.

## Features

- **Rate Limiting**: 24-hour cooldown per wallet address to prevent abuse
- **Redis Integration**: Uses Upstash Redis for distributed rate limiting
- **RESTful API**: Simple HTTP endpoints for token requests
- **Blockchain Agnostic**: Works with any Cosmos SDK-based chain

## Prerequisites

- Go 1.22.8 or higher
- Redis instance (Upstash recommended)
- Cosmos SDK blockchain CLI tool (gaiad, sourdoughd, etc.)
- Wallet with tokens to distribute

## Installation

```bash
# Clone the repository
git clone https://github.com/s16rv/cosmos-rest-faucet
cd cosmos-rest-faucet

# Install dependencies
go mod tidy

# Build the binary
go build -o faucet
```

## Configuration

### Required Parameters

- `--port`: Server port (e.g., 9000)
- `--cli`: Blockchain CLI binary name (e.g., gaiad, sourdoughd)
- `--address`: Faucet wallet address that will send tokens
- `--alias`: Wallet alias/name in keyring
- `--node`: Blockchain RPC endpoint
- `--home`: Blockchain home directory
- `--keyring-backend`: Keyring backend type (usually "test" for testnets)
- `--chain-id`: Blockchain network chain ID
- `--redis-url`: Redis connection URL (Upstash format)

### Redis URL Format

For Upstash Redis, the URL format is:

```
rediss://default:PASSWORD@HOST:PORT
```

## Usage

### Start the Faucet Server

```bash
./faucet --port 9000 \
         --cli sourdoughd \
         --address sourdough1abcd2efg4h5ijklmno6pqr7stuvwxyz89abc0d \
         --alias alice \
         --node https://YOUR_RPC_ENDPOINT \
         --home ~/.sourdoughd \
         --keyring-backend test \
         --chain-id sourdough-1 \
         --redis-url "rediss://default:YOUR_PASSWORD@YOUR_HOST:6379"
```

### Request Tokens

```bash
curl -d '{"recipient":"sourdough1abcd2efg4h5ijklmno6pqr7stuvwxyz89abc0d", "amount":"5000000usrdh"}' \
     -H "Content-Type: application/json" \
     -X POST http://localhost:9000/request
```

## API Endpoints

### `GET /ping`

Health check endpoint.

**Response:**

```json
{
  "message": "pong"
}
```

### `POST /request`

Request tokens from the faucet.

**Request Body:**

```json
{
  "recipient": "sourdough1abcd2efg4h5ijklmno6pqr7stuvwxyz89abc0d",
  "amount": "5000000usrdh"
}
```

**Success Response (200):**

```json
{
  "message": "success",
  "txHash": "ABC123DEF456..."
}
```

**Rate Limited Response (429):**

```json
{
  "error": "Rate limit exceeded",
  "message": "You can only request tokens once per 24 hours",
  "timeRemaining": "23h 45m",
  "nextRequestTime": "2025-11-26T15:30:00Z"
}
```

## Rate Limiting

Each wallet address can only request tokens **once per 24 hours**. The rate limiting is enforced using Redis with automatic expiration.

## Development

### Run in Development Mode

```bash
go run main.go [flags...]
```

### Project Structure

```
├── main.go              # Application entry point
├── config/
│   └── config.go        # Configuration management
├── handler/
│   └── handler.go       # HTTP request handlers
├── command/
│   └── command.go       # Blockchain command execution
├── ratelimit/
│   └── ratelimit.go     # Redis-based rate limiting
└── go.mod              # Go module dependencies
```
