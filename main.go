package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/deresmos/protoc-gen-template/datatype"
	"google.golang.org/protobuf/proto"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
	plugin "google.golang.org/protobuf/types/pluginpb"
)

func parseReq(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var req plugin.CodeGeneratorRequest
	if err = proto.Unmarshal(buf, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

type fileGenerator struct {
	packageName             string // TODO:
	option                  *ProtoOption
	fileDescriptorGenerator *FileDescriptorGenerator
	fileTemplate            *template.Template
	outputPathTemplate      *template.Template
}

func (g *fileGenerator) run(fileDescriptor *FileDescriptor) ([]*plugin.CodeGeneratorResponse_File, error) {
	var files []*plugin.CodeGeneratorResponse_File
	switch g.option.GenerateType {
	case "message":
		for _, message := range fileDescriptor.Messages {
			responseFile, err := g.generateResponseFile(message)
			if err != nil {
				return nil, err
			}
			files = append(files, responseFile)
		}
	case "service":
		for _, service := range fileDescriptor.Services {
			responseFile, err := g.generateResponseFile(service)
			if err != nil {
				return nil, err
			}
			files = append(files, responseFile)
		}
	case "method":
		for _, service := range fileDescriptor.Services {
			for _, method := range service.Methods {
				responseFile, err := g.generateResponseFile(method)
				if err != nil {
					return nil, err
				}
				files = append(files, responseFile)
			}
		}
	case "file":
		responseFile, err := g.generateResponseFile(fileDescriptor)
		if err != nil {
			return nil, err
		}
		files = append(files, responseFile)
	}

	files = filterResponseFiles(files, func(file *plugin.CodeGeneratorResponse_File) bool {
		return file.GetName() != ""
	})

	return files, nil
}

func filterResponseFiles(files []*plugin.CodeGeneratorResponse_File, filter func(*plugin.CodeGeneratorResponse_File) bool) []*plugin.CodeGeneratorResponse_File {
	var newFiles []*plugin.CodeGeneratorResponse_File
	for _, file := range files {
		if filter(file) {
			newFiles = append(newFiles, file)
		}
	}
	return newFiles
}

func (g *fileGenerator) generateResponseFile(data any) (*plugin.CodeGeneratorResponse_File, error) {
	b := bytes.NewBuffer([]byte{})
	err := g.fileTemplate.Execute(b, data)
	if err != nil {
		return nil, err
	}
	outputPathBuffer := bytes.NewBuffer([]byte{})
	err = g.outputPathTemplate.Execute(outputPathBuffer, data)
	if err != nil {
		return nil, err
	}

	outputPath := outputPathBuffer.String()
	return &plugin.CodeGeneratorResponse_File{
		Name:    &outputPath,
		Content: proto.String(b.String()),
	}, nil
}

func processReq(req *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	protoOption, err := NewProtoOptionFromString(req.GetParameter())
	if err != nil {
		panic(err)
	}
	dataType, err := datatype.NewDataType(protoOption.Language)
	if err != nil {
		panic(err)
	}

	fileDescriptorGenerator := NewFileDescriptorGenerator(req.GetParameter(), dataType)
	fileTmpl, err := initFileTemplate(protoOption.TemplatePath)
	if err != nil {
		panic(err)
	}
	outputTmpl, err := initOutputPathTemplate(protoOption.OutputPath)
	if err != nil {
		panic(err)
	}

	fileGenerator := &fileGenerator{
		packageName:             req.GetParameter(),
		option:                  protoOption,
		fileDescriptorGenerator: fileDescriptorGenerator,
		fileTemplate:            fileTmpl,
		outputPathTemplate:      outputTmpl,
	}

	files := make(map[string]*descriptor.FileDescriptorProto)
	for _, f := range req.ProtoFile {
		files[f.GetName()] = f
	}
	var resp plugin.CodeGeneratorResponse
	features := uint64(plugin.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	resp.SupportedFeatures = &features

	if protoOption.AllowMerge {
		var megeredFileDescriptor *FileDescriptor
		for _, fname := range req.FileToGenerate {
			f := files[fname]
			fileDescriptor, err := fileGenerator.fileDescriptorGenerator.Run(f)
			if err != nil {
				panic(err)
			}

			megeredFileDescriptor = megeredFileDescriptor.Append(fileDescriptor)
		}

		files, err := fileGenerator.run(megeredFileDescriptor)
		if err != nil {
			panic(err)
		}
		if !protoOption.Overwrite {
			files = filterFirstTimeOutputFiles(files)
		}
		resp.File = append(resp.File, files...)
	} else {
		for _, fname := range req.FileToGenerate {
			f := files[fname]
			fileDescriptor, err := fileGenerator.fileDescriptorGenerator.Run(f)
			if err != nil {
				panic(err)
			}

			files, err := fileGenerator.run(fileDescriptor)
			if err != nil {
				panic(err)
			}
			if !protoOption.Overwrite {
				files = filterFirstTimeOutputFiles(files)
			}
			resp.File = append(resp.File, files...)
		}
	}

	return &resp
}

func filterFirstTimeOutputFiles(files []*plugin.CodeGeneratorResponse_File) []*plugin.CodeGeneratorResponse_File {
	var newFiles []*plugin.CodeGeneratorResponse_File
	for _, file := range files {
		// パスが存在するかチェック
		_, err := os.Stat(file.GetName())
		if !os.IsNotExist(err) {
			log.Printf("Skip generate %s. Bacause overwrite option is false.", file.GetName())
			continue
		}

		newFiles = append(newFiles, file)
	}

	return newFiles
}

func emitResp(resp *plugin.CodeGeneratorResponse) error {
	buf, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(buf)
	return err
}

func Run(r io.Reader) error {
	req, err := parseReq(r)
	if err != nil {
		return err
	}

	resp := processReq(req)

	return emitResp(resp)
}

func main() {
	if err := Run(os.Stdin); err != nil {
		log.Fatalln(err)
	}
}
