SHELL := /bin/zsh

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}{printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Clean proto files
	@rm -rf proto/gen

.PHONY: proto
proto: ## Generate proto files
	@protoc -I=proto --go_out=. --go-grpc_out=. proto/*.proto
	@echo "Proto files generated successfully."

.PHONY: server
server: ## Run the gRPC server
	@echo "Starting server..."
	@go run ./cmd/grpcapi/server.go
	@echo "Server started successfully."