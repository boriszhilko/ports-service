# Run the application with docker-compose
run:
	docker-compose up --build

# Run tests with docker-compose using a separate compose file
test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

# Install golangci-lint and run linting on all go files in the current directory and subdirectories
lint:
	docker build --target=lint -f ./build/Dockerfile .