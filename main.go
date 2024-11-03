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
			b := bytes.NewBuffer([]byte{})
			err := g.fileTemplate.Execute(b, message)
			if err != nil {
				return nil, err
			}
			outputPathBuffer := bytes.NewBuffer([]byte{})
			err = g.outputPathTemplate.Execute(outputPathBuffer, message)
			if err != nil {
				return nil, err
			}

			outputPath := outputPathBuffer.String()
			files = append(files, &plugin.CodeGeneratorResponse_File{
				Name:    &outputPath,
				Content: proto.String(b.String()),
			})
		}
	case "service":
		for _, service := range fileDescriptor.Services {
			b := bytes.NewBuffer([]byte{})
			err := g.fileTemplate.Execute(b, service)
			if err != nil {
				return nil, err
			}
			outputPathBuffer := bytes.NewBuffer([]byte{})
			err = g.outputPathTemplate.Execute(outputPathBuffer, service)
			if err != nil {
				return nil, err
			}

			outputPath := outputPathBuffer.String()
			files = append(files, &plugin.CodeGeneratorResponse_File{
				Name:    &outputPath,
				Content: proto.String(b.String()),
			})
		}
	case "file":
		b := bytes.NewBuffer([]byte{})
		err := g.fileTemplate.Execute(b, fileDescriptor)
		if err != nil {
			return nil, err
		}
		outputPath := g.option.OutputPath
		files = append(files, &plugin.CodeGeneratorResponse_File{
			Name:    &outputPath,
			Content: proto.String(b.String()),
		})
	}
	return files, nil
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
			resp.File = append(resp.File, files...)
		}
	}

	return &resp
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
