# GitHub API Demo
This is a simple API demo that can authorize user with either a GitHub Token, or through harrison-minibucks GH App.

## Setup
1. Install [golang](https://go.dev/dl/)
2. Install [Protobuf](https://github.com/protocolbuffers/protobuf/releases)
3. Optionally, install [Make](https://www.gnu.org/software/make/) cli, or `sudo apt install make`
4. Run `make init` to install the dependency tools (Or run the command lines listed in the Makefile)

## Development
1. Run `go mod tidy`
2. Run `make db`
3. Run `make run`

## Docker Setup
1. Install [Docker](https://docs.docker.com/engine/install/)
2. Run `docker compose up`

## Testing
Refer to the file `rests/todo-api.http` (Install [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) in VSCode for easier testing)

*GRPC server is not started.*

## Others
You may refer to the Makefile.