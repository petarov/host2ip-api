# host2ip-api

A tiny HTTP API server that resolves hostname IP addresses

# Usage

To start the service on `[::1]:7029` run:

    go run main.go

The server is now accessible at `http://localhost:7029`.

# API

List of available API calls:

- `/` - shows available junctions
- `/lookup/:name` - fetches a list of IP addresses for the specified FQDN or hostname

# License 

[MIT](LICENSE)
