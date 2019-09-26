# Go GitHub

## Requirements

- Go >= 1.11
- Github Personal Access Token
  - Get one from https://github.com/settings/tokens
  - Set the token as env variable `GH_ACCESS_TOKEN`

## Web mode
1. Run

   ```sh
   go run cmd/web/web.go
   ```
2. Open browser: http://localhost:8080/gh/profile/antonybudianto

## GRPC mode
1. Run

    ```sh
    go run cmd/grpc_server/server.go
    ```

2. Try using GRPC client:

    ```sh
    go run cmd/grpc_client/client.go <github-username>
    ```

3. Misc: Generate proto
   
    ```sh
    protoc --go_out=plugins=grpc:. protos/*.proto
    ```

## CLI mode
1. Run

   ```sh
   go run cmd/cli/cli.go <github-username>
   ```

## Build (for GCP Debian)

```sh
$ env GOOS=linux GOARCH=amd64 go build -o bin/web cmd/web/web.go
```

# License
MIT
