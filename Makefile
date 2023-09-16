build-app:
	@go build -o bin/app ./cmd/api/

run: build-app
	@./bin/app


docker:
	@echo "building docker file"
	@docker build -t app -f Dockerfile .
	@echo "running API inside Docker container"
	@docker run -p 9090:9090 app


clean: 
	@rm -rf bin

