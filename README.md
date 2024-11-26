# Ethereum Transaction Parser

This project is an Ethereum blockchain transaction parser that allows querying transactions for subscribed addresses. It uses the Ethereum JSON-RPC API to interact with the blockchain and can be configured to store data in memory or in a LevelDB database.

## Features

- Query the current block of the Ethereum blockchain.
- Subscribe addresses for transaction monitoring.
- Get inbound and outbound transactions for subscribed addresses.
- Store data in memory or LevelDB.
- Expose an HTTP API to interact with the parser.

## Project Structure

- `cmd/`: Contains the main code to start the server.
- `internal/handlers/`: Contains the HTTP handlers.
- `pkg/config/`: Contains the project configuration.
- `pkg/logger/`: Contains the project logger.
- `pkg/parser/`: Contains the parser interface and implementations for memory and LevelDB.
- `pkg/ethereum/`: Contains the Ethereum client to interact with the JSON-RPC API.

## Installation

1. Clone the repository:

```sh
git clone https://github.com/jmsilvadev/tx-parser.git
cd tx-parser
```

2. Install the dependencies:

```sh
go mod tidy
```

## Usage

### Start the Server

To start the server, run the following command:

```sh
go run cmd/main.go
```

### Docker

Build and start the Docker containers:

```sh
make up-build
```

This command will:
- Build the Docker images specified in the `docker-compose.yml` file.
- Start the containers using Docker Compose.


### API Endpoints

- `GET /health`: Check the server health.
- `GET /v1/get-current-block`: Return the current block of the Ethereum blockchain.
- `POST /v1/subscribe?address={address}`: Subscribe an address for transaction monitoring.
- `GET /v1/get-transactions?address={address}`: Return inbound and outbound transactions for a subscribed address.


#### Request Examples

##### Get Current Block

```sh
curl -X GET http://localhost:5000/v1/get-current-block
```

##### Subscribe Address

```sh
curl -X POST http://localhost:5000/v1/subscribe?address=0x123
```

##### Get Transactions

```sh
curl -X GET http://localhost:5000/v1/get-transactions?address=0x123
```

### Tests
To run the tests, use the following command:

```sh
make tests
```

### Commands
```
$ make

build-image                    Build docker image in daemon mode
build-server                   Build server component
clean                          Clean all builts
clean-tests                    Clean tests
down                           Stop docker container
logs                           Watch docker log files
tests-cover                    Run tests with coverage
tests                          Run unit tests
up-build                       Start docker container and rebuild the image
up                             Start docker container
```