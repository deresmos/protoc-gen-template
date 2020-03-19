.PHONY: run
run: install
	protoc -I. \
	  -I${GOPATH}/src \
	  --entity_out=. hello.proto

.PHONY: instlal
install: format
	go install

.PHONY: format
format:
	gofmt -w .
