build:
	sh ./set_version.sh
	go mod tidy
	go build -o ./bin/tplr ./cmd/tplr

build_image:
	docker build . --tag sebasmannem/tplr

debug:
	go build -gcflags "all=-N -l" -o ./bin/tplr ./cmd/tplr
	~/go/bin/dlv --headless --listen=:2345 --api-version=2 --accept-multiclient exec ./bin/tplr -- -c config/tplr.yaml -d

run:
	./bin/tplr -d

fmt:
	gofmt -w .
	goimports -w .
	gci write .

compose:
	./docker-compose-tests.sh

test: gotest sec lint

sec:
	gosec ./...

lint:
	golangci-lint run

gotest:
	go test -v ./...
