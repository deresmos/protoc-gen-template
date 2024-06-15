.PHONY: test
test: go-test integration-test

.PHONY: go-test
go-test:
	go test -v ./...

.PHONY: integration-test
integration-test:
	protoc --template_out='template=test/message.template,lang=typescript,generate_type=message,output_path=./test/output/message/{{toSnakeCase .MessageName}}.txt:.' test/message.proto
	protoc --template_out='template=test/service.template,lang=typescript,generate_type=service,output_path=./test/output/service/{{toSnakeCase .ServiceName}}.txt:.' test/service.proto
	git diff --exit-code --quiet
