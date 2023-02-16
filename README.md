# host2ip

A tiny HTTP API server that resolves hostname IP addresses

# API

List of available API calls:

- `/` - shows available junctions
- `/lookup/:name` - fetches a list of IP addresses for the specified FQDN or hostname

# Usage

To start the service on `[::1]:7029` run:

    go run main.go

The server is now accessible at `http://localhost:7029`.

Example query:

    curl -s http://localhost:7029/lookup/google.bg | jq

```json
{
  "addresses": [
    "172.217.16.131",
    "2a00:1450:4001:808::2003"
  ]
}
```

# License 

[MIT](LICENSE)
