# afrochain-faucet

This repository provides a fully functional faucet implementation tailored for Cosmos SDK-based blockchain networks. The faucet allows developers and testers to request tokens from a predefined pool for testing purposes, streamlining the development and testing processes in blockchain environments.

## Features

- **Rate Limiting**: 24-hour cooldown per wallet address to prevent abuse
- **In-Memory Cache**: Built-in rate limiting without external dependencies
- **RESTful API**: Simple HTTP endpoints for token requests
- **Blockchain Agnostic**: Works with any Cosmos SDK-based chain

## Prerequisites

- Go 1.24.6 or higher
- Cosmos SDK blockchain CLI tool (afrochaind)
- Wallet with tokens to distribute

## Installation

```bash
# Clone the repository
git clone https://github.com/afrochainorg/afrochain-faucet
cd afrochain-faucet

# Install dependencies
go mod tidy

# Build the binary
go build -o faucet
```

## Configuration

### Required Parameters

- `--port`: Server port (e.g., 9000)
- `--cli`: Blockchain CLI binary name (e.g., afrochaind)
- `--address`: Faucet wallet address that will send tokens
- `--alias`: Wallet alias/name in keyring
- `--node`: Blockchain RPC endpoint
- `--home`: Blockchain home directory
- `--keyring-backend`: Keyring backend type (usually "test" for testnets)
- `--chain-id`: Blockchain network chain ID

## Usage

### Start the Faucet Server

```bash
./faucet --port 9000 \
         --cli afrochaind \
         --address afro13sllcdsqhjektac5r6h50dvjrthm0yt6m6sfkm \
         --alias alice \
         --node https://127.0.0.1:26657 \
         --home /home/afrochain/.afrochaind-1 \
         --keyring-backend test \
         --chain-id afrochain-1
```

### Request Tokens

```bash
curl -d '{"recipient":"afro13sllcdsqhjektac5r6h50dvjrthm0yt6m6sfkm", "amount":"5000000aafro"}' \
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
  "recipient": "afro13sllcdsqhjektac5r6h50dvjrthm0yt6m6sfkm",
  "amount": "5000000aafro"
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

Each wallet address can only request tokens **once per 24 hours**. The rate limiting is enforced using an in-memory cache with automatic expiration.

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
│   └── ratelimit.go     # In-memory rate limiting
└── go.mod              # Go module dependencies
```
