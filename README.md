## Installation

In order to run the application, [Docker](https://www.docker.com/) is required.

Navigate to the project's directory and run the command `docker-compose up -d`.
This will spin up every dependency (MongoDB, Redis, Kafka) along with the api project.

If you want to run the project in your machine, you can use [Make](https://www.gnu.org/software/make/).

The command `make api` will build the project and produce an executable found at **cmd/api**.

Before running the executable, make sure that the Docker api is stoped and all of the dependencies are still running.

## Using the api

You can call the api endpoints using [Postman](https://www.postman.com/).

In the folder **scripts/postman** you will find ready to import collections and environment.

In order to login, you can use the `TEST_EMAIL` and `TEST_PASS` credentials found in **cmd/api/.env**.

Each request will produce logs which you can view in the programs output and each mutating operation will publish and consume a message from Kafka along with a debug log of the consumed message.

## Integration tests

Currently there is a single integration test implemented under the **integration_tests** folder, which needs refactoring.

Because the test uses the .env file in order to run, you can use the 'Run test' functionality of [VS Code](https://code.visualstudio.com/).

## Linter

The program uses [Golangci-lint](https://github.com/golangci/golangci-lint) for linting.
The command `make lint` will run the linter using the **golangci.yml** configuration file.