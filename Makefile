## build: Compile the app into the executable binery
build:
	@echo ">  Building the app..."
	@go build -o time-server
	@go build -o time-client cmd/client/client.go 
	@echo ">  Done"

.PHONY: build