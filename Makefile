build: 
	@go build -o ./bin/go-basketball

run: build
	@ ./bin/go-basketball