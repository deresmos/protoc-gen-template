.PHONY: test
test: go-install go-test integration-test

.PHONY: go-install
go-install:
	go install .

.PHONY: go-test
go-test:
	go test -v ./...

.PHONY: integration-test
integration-test:
	rm -rf ./test/output/*
	protoc --template_out='template=test/message.template,lang=typescript,generate_type=message,output_path=./test/output/message/{{toSnakeCase .MessageName}}.txt:.' test/message.proto
	protoc --template_out='template=test/service.template,lang=typescript,generate_type=service,output_path=./test/output/service/{{toSnakeCase .ServiceName}}.txt:.' test/service.proto
	protoc --template_out='template=test/allow-merge/merge.template,lang=typescript,generate_type=file,allow_merge=true,output_path=./test/output/allow-merge/merge.txt:.' test/allow-merge/*.proto
	protoc --template_out='template=test/option/overwrite/overwrite.template,lang=typescript,generate_type=file,overwrite=false,output_path=./test/option/overwrite/output/overwrite.txt:.' test/option/overwrite/main.proto
	git diff --exit-code --quiet ./test/output
