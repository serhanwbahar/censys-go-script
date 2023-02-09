# Censys Go Script

This CLI takes a domain name and returns a JSON with several intelligence data fetched from Censys.

## Installation

- To download required packages, use `go mod download`.
- You only need to Go to build & run the CLI. 
- Makefile is used to provide quick aliases for Go commands, though it is not necessary.
- A [Dockerfile](./Dockerfile) is provided to build an image.

## Usage

You need to have the following environment variables in your environment, or can provide them via `.env` file:

```sh
CENSYS_API_ID=<id>
CENSYS_API_SECRET=<secret>
```

The output will be written as a JSON file under the `out` folder, as `<domain-name>.json`. Makefile provides shorthand commands:

- `make run` to run the CLI via `go run`.
- `make build` to build a binary under the `bin` folder.
- `make test` to run tests.
- `make tidy` as a shorthand for `go mod tidy`.
- `make clean` to delete the binary.
- `make create-image` will use the Dockerfile to build an image for the CLI.

## Methodology

We obtain the IPs of a given domain via:

- `/api/v2/hosts/search?q={{DOMAIN}}` but note that we dont do `q=name: {{DOMAIN}}` although that is how it is described in the [documentation](https://search.censys.io/search/language?resource=hosts).

Then, the following endpoints are used for each IP:

- `/api/v2/hosts/{{IP}}` fetches the entire host entity by IP address and returns the most recent Censys view of the host and its services.
- `/api/v2/hosts/{{IP}}/names` fetches a list of host names for the specified IP address.
- `/api/v2/experimental/hosts/{{IP}}/events?per_page=50` fetches a list of events for the host with the specified IP address.
- **Tags** are not used, as they are restricted.
- **Comments** are not used, as they are restricted.

Note that the results will be saved every few IPs, as specified as a [constant](./internal/constants/constants.go) value.

## Styling

File structure and testing is done as per the [Effective Go](https://go.dev/doc/effective_go) guide.
