# Public Sector Backend

## Requisities

1. Go 1.27 or higher
2. [GolangCI-Lint](https://golangci-lint.run/docs/welcome/install/local/). In Windows use Chocolatey.

## Development

1. Format the code
   ```sh
   go fmt ./...
   ```

2. Run linting
   ```sh
   golangci-lint run
   ```

3. Run the application
   ```sh
   go run ./cmd/api/main.go
   ```
