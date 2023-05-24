# GitHub API Demo
This is a simple API demo that can authorize user with either a GitHub Token, or through harrison-minibucks GH App.

## Setup
1. Install [golang](https://go.dev/dl/)
2. Install [Protobuf](https://github.com/protocolbuffers/protobuf/releases)
3. Optionally, install [Make](https://www.gnu.org/software/make/) cli, or `sudo apt install make`
4. Run `make init` to install the dependency tools (Or run the command lines listed in the Makefile)
5. Generate a new GitHub Personal Access Token ([PAT](https://docs.github.com/en/rest/overview/authenticating-to-the-rest-api?apiVersion=2022-11-28#basic-authentication))

### Optionally
1. If you want to use GitHub App, please [setup](https://docs.github.com/en/apps/creating-github-apps/setting-up-a-github-app/about-creating-github-apps) one and add the `client_id` and `client_secret` to `./configs/secret.yaml`, sample is shown on `./secret-sample.yaml`.
2. Setup the GitHub APP callback url to `http://localhost:8000/github/callback`
3. Otherwise, you may use GitHub PAT for the `Authorization` header to test the API

## Development
1. Set `go env -w GO111MODULE=on` (to accomodaate Kratos)
2. Run `go mod tidy`
3. Run `make db`
4. Run `make run`

## Docker Setup
1. Install [Docker](https://docs.docker.com/engine/install/)
2. Run `docker compose up`

## Unit Test
Run `make test` to run unit tests

## HTTP API Testing
1. Please refer to the file `./rests/todo-api.http` and `./rests/gh-api.http` (Install [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) in VSCode for easier testing).
2. Or, you can head to [Swagger](https://editor.swagger.io/) and paste the contents in `openapi.yaml`
3. You can head to http://localhost:8000/github/login to login into your GitHub App, and a `Session` will be returned
4. You could either use `Authorization` with GitHub PAT, or use `Session` with the session ID returned as your header

**/github/avatar and /github/logout is not implemented yet**
*GRPC server is not started.*

## Others
You may refer to the Makefile.