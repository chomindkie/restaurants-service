test:
	go test ./... -v -coverprofile .coverage.txt
	go tool cover -func .coverage.txt
initialize:
	go mod tidy
	go get -v
	go run main.go
