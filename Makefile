SHELL=/bin/sh

.PHONY: ci
ci: test

.PHONY: gen
gen:
	GO111MODULE=off go get github.com/matryer/moq
	moq -out controller/mock_service.go -pkg controller ./service Service

.PHONY: fe-build
fe-build:
	cd ./client && npm run build && cd ..

.PHONY: be-build
be-build:
	CGO_ENABLED=0 go build -ldflags "-w" -a -o water-jug-riddle-service .

.PHONY: all-build
all-build:
	rice append -i . --exec water-jug-riddle-service

.PHONY: build
build: fe-build be-build all-build

.PHONY: test
test:
	go test -cover ./... -count=1

.PHONY: cover
cover:
	GO111MODULE=off go get github.com/axw/gocov/gocov
	GO111MODULE=off go get -u gopkg.in/matm/v1/gocov-html
	${GOPATH}/bin/gocov test ./... | ${GOPATH}/bin/gocov-html > coverage.html
	open coverage.html

.PHONY: release
release: build test cover

.PHONY: run
run:
	go run main.go