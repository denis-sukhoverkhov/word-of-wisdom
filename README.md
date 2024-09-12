# Word-of-wisdom TCP-Sever

## Project Overview

This project contains a TCP server that implements Proof of Work (PoW) to protect against DDoS attacks, serving quotes from a “Word of Wisdom” collection. The project also includes a client that connects to the server, solves the PoW challenge, and receives quotes.

## Choosing a PoW Algorithm

We will use Hashcash as the PoW algorithm. It is simple and widely used in scenarios where the server needs to limit client requests through computational effort.

This approach is computationally challenging for the client, making it harder to flood the server with requests, thus mitigating DDoS attacks.

## Simple and fast project running using docker-compose

 ```bash
   make up
   ```

### Environment Variables

For local running this project relies on environment variables to configure certain parameters. Make sure to create an `.env` filein the root directory with appropriate values (`.env_example`), such as:

```bash
# The difficulty level for the Proof-of-Work (PoW) task issued to clients
APP_DIFFICULTY=4 
APP_ADDR="127.0.0.1:8080"
#  The number of worker processes that handle client connections concurrently.
APP_WORKERCOUNT=10
# The time (in seconds) the server will wait for active connections to finish before shutting down
APP_SHUTDOWNTIMEOUT=5

# Client-specific variables
APP_SERVERADDR="127.0.0.1:8080"
# The number of requests per second (RPS) that the client will send to the server. 
APP_RPS=5
# The total number of requests the client will send to the server.
APP_TOTALREQUESTS=100
```

**Important: the `.env` file redeclare base env variables from `configs` directory.**

**Important: youdon't need configure `.env` file if you run application through `docker-compose`**

### Makefile Overview

The provided Makefile simplifies the process of building, running, testing, and linting the project, as well as building Docker images for both the server and the client.

### Prerequisites

- **Go**: Ensure that you have Go installed. You can install it from [here](https://golang.org/doc/install).
- **Docker**: Make sure Docker is installed to build and run the Docker containers. Install Docker from [here](https://docs.docker.com/get-docker/).

### Available Make Commands

The Makefile supports various commands to manage your development workflow.

#### Linting and Testing

1. **Install Lint Tool**:
   Installs `golangci-lint` for running lint checks.

   ```bash
   make install-lint
   ```

2. **Run Lint Checks**:
   Runs `golangci-lint` to check the Go code for style, bugs, and best practices.

   ```bash
   make lint
   ```

3. **Go Vet**:
   Runs `go vet` to report any suspicious constructs in the code.

   ```bash
   make vet
   ```

4. **Run Tests**:
   Runs all the tests in the project using `go test`.

   ```bash
   make test
   ```

#### Building the Project

1. **Build the Server**:
   Builds the server binary and places it in the `bin/` directory.

   ```bash
   make build-server
   ```

2. **Build the Client**:
   Builds the client binary and places it in the `bin/` directory.

   ```bash
   make build-client
   ```

#### Running the Project

1. **Run the Server**:
   Runs the server locally using environment variables specified in the `.env` file.

   ```bash
   make run-server
   ```

2. **Run the Client**:
   Runs the client locally using environment variables specified in the `.env` file.

   ```bash
   make run-client
   ```

#### Docker Commands

1. **Build Docker Image for Server**:
   Builds a Docker image for the server using the Dockerfile in `Dockerfile.server`.

   ```bash
   make docker-build-server
   ```

2. **Build Docker Image for Client**:
   Builds a Docker image for the client using the Dockerfile in `Dockerfile.client`.

   ```bash
   make docker-build-client
   ```

3. **Build Docker Images for Both Server and Client**:
   Builds Docker images for both the server and the client.

   ```bash
   make docker-build-all
   ```
