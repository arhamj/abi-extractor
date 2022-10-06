.PHONY: build
build:
	go build -o abi-extractor cmd/*.go

.PHONY: install
install:
	go build -o ${GOPATH}/bin/abi-extractor cmd/*.go

.PHONY: test
test:
	go test -v -coverpkg=./pkg/... ./pkg/...