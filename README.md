# XMapp

## Overview

This project is a Golang-based API application that leverages Docker to manage dependencies, including MongoDB, Redis, and Kafka. The project also supports integration with Postman for testing and debugging, and includes a linter configuration with Golangci-lint.

## Prerequisites
To run the application, ensure the following tools are installed on your machine:

- [Docker](https://www.docker.com/) (required for containerized execution)
- [Make](https://www.gnu.org/software/make/) (optional for local builds)
- [VS Code](https://code.visualstudio.com/) (optional for running integration tests)
- [Postman](https://www.postman.com/) (optional for testing API endpoints)

## Installation

1. Run with Docker: To start the application along with its dependencies, navigate to the project directory and run the following command:

`docker-compose up -d`

This will automatically spin up all required services (MongoDB, Redis, Kafka) along with the API itself.

2. Run locally: If you prefer to run the project on your local machine, you can build the application using `make`:

`make api`

This will generate the API executable, which can be found under `cmd/api/`.

**Note**: Before running the executable, ensure that Docker's API is stopped, but that the required dependencies (MongoDB, Redis, Kafka) are still running.

## API Usage
You can interact with the API using [Postman](https://www.postman.com/). Ready-to-import collections and environment configurations are available in the `scripts/postman` folder.

### Authentication
To authenticate with the API, use the following credentials from the `.env` file located in `cmd/api/.env`:

Email: `TEST_EMAIL`
Password: `TEST_PASS`

Each API call produces detailed logs visible in the application output. Mutating operations also trigger messages to be published and consumed from Kafka, along with corresponding debug logs.

## Integration Tests

There is currently one integration test available in the `integration_tests` directory. However, this test requires refactoring.

### Running Tests:

The integration tests rely on the `.env` configuration. You can easily run them using the "Run Test" functionality in [VS Code](https://code.visualstudio.com/).

## Code Linting
This project uses [Golangci-lint](https://github.com/golangci/golangci-lint) for static code analysis and linting. You can run the linter with:

`make lint`

This will apply the rules specified in the `golangci.yml` configuration file.