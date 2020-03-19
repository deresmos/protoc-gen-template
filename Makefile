.PHONY: run
run: install
	protoc -I. \
	  -I${GOPATH}/src \
		--entity_out=file=./js.template:. hello.proto

.PHONY: instlal
install: format
	go install

.PHONY: format
format:
	gofmt -w .
