PID = $(shell cat tmp/legsc.pid)

setup:
	go get -u golang.org/x/lint/golint
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/oxequa/realize
	go get -u github.com/motemen/gore
	go get -u github.com/golang/mock/gomock
	go get -u github.com/mitchellh/gox
	go install github.com/golang/mock/mockgen
lint:
	go vet ./...
	golint -set_exit_status ./...
fmt:lint
	goimports -w .
vet:
	go tool vet ./
build: fmt lint
	go build -ldflags="-s -w" -o dist/legsc main.go
start:
	go run main.go start
run:
	go run main.go start -f -c ./config.toml
stop:
	go run main.go stop
restart:
	go run main.go restart
test:
	env ENV=test go test -cover -race ./...
