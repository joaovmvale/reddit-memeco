# MemeCo API

A Go-based API service that serves memes with rate limiting capabilities. The service implements a sliding window rate limiter to protect against excessive requests.

## Features

- RESTful API endpoints for meme retrieval
- Sliding window rate limiting (3 requests per second per client in the implementation example)
- Support for random meme selection
- Client IP detection with X-Forwarded-For header support
- Comprehensive test suite

## API Endpoints

- `GET /memes/random` - Get a random meme
- `GET /memes/{id}` - Get a specific meme by ID

## Requirements

- Go 1.20 or higher
- gorilla/mux for routing

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/reddit-memeco.git
cd reddit-memeco
```

2. Install dependencies:
```bash
go mod download
```

## Running the Server

Start the server on port 8080:
```bash
go run cmd/server/main.go
```

The server will start and listen on `http://localhost:8080`.

## Testing

### Running Unit Tests

Run all tests:
```bash
go test ./... -v
```

Run specific test packages:
```bash
go test ./pkg/meme -v
go test ./api/handlers -v
```

### Running API Tests

The project includes two test scripts to verify the API functionality:

1. Basic API test:
```bash
./scripts/test_api.sh
```
This script tests:
- Random meme retrieval
- Specific meme retrieval
- Non-existent meme handling
- Basic rate limiting
- Multiple client simulation

2. Rate limit test:
```bash
./scripts/test_rate_limit.sh
```
This script tests:
- Single client burst requests
- Multiple clients making concurrent requests
- Rate limit enforcement

## Project Structure

```
.
├── api/
│   └── handlers/         # HTTP request handlers
├── cmd/
│   └── server/          # Server entry point
├── pkg/
│   ├── meme/            # Meme service implementation
│   └── ratelimiter/     # Rate limiter implementation
├── scripts/             # Test scripts
└── example/             # Example usage
```

## Rate Limiter Implementation

The rate limiter uses a sliding window approach:
- Tracks requests per client IP
- Allows 3 requests per second per client
- Uses a cleanup goroutine to remove old entries
- Thread-safe implementation with mutex locks

## License

MIT License - see LICENSE file for details