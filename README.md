# Go GitHub

## Web mode
1. Run

   ```sh
   go run cmd/web/web.go
   ```
2. Open browser: http://localhost:8080/profile/antonybudianto

## GRPC mode
1. Run

    ```sh
    go run cmd/grpc_server/server.go
    ```

2. Try using GRPC client:

    ```sh
    go run cmd/grpc_client/client.go
    ```

3. Misc: Generate proto
   
    ```sh
    protoc --go_out=plugins=grpc:. protos/*.proto
    ```

## CLI mode
1. Run

   ```sh
   go run cmd/cli/cli.go <your-username>
   ```


# License
MIT
