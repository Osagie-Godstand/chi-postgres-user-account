build-api:
	@go build -o bin/api ./cmd/api/

run: build-api
	@./bin/api

test: 
	@go test -v ./...

docker:
	@echo "building docker file"
	@docker build -t api -f Dockerfile .
	@echo "running API inside Docker container"
	@docker run -p 9090:9090 api


clean: 
	@rm -rf bin

