package main

import (
	"strings"

	"github.com/deresmos/protoc-gen-template/datatype"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type FileDescriptor struct {
	PackageName string
	Messages    []MessageDescriptor
	Services    []ServiceDescriptor
}

type MessageDescriptorList []MessageDescriptor

func (m MessageDescriptorList) GetByMessageName(name string) *MessageDescriptor {
	for _, message := range m {
		if message.MessageName == name {
			return &message
		}
	}

	return nil
}

type MessageDescriptor struct {
	MessageName string
	Fields      []MessageFieldDescriptor
	Parents     []MessageDescriptor
}

type MessageFieldDescriptor struct {
	Name         string
	DataTypeName string
	IsRequired   bool
	IsOptional   bool
	IsTimestamp  bool
	IsRepeated   bool
}

type ServiceDescriptor struct {
	ServiceName string
	Methods     []ServiceMethodDescriptor
	Messages    MessageDescriptorList
}

type ServiceMethodDescriptor struct {
	Name         string
	Input        *MessageDescriptor
	Output       *MessageDescriptor
	Dependencies []MessageDescriptor
}

type FileDescriptorGenerator struct {
	packageName string
	dataType    *datatype.DataType
}

func NewFileDescriptorGenerator(packageName string, dataType *datatype.DataType) *FileDescriptorGenerator {
	return &FileDescriptorGenerator{
		packageName: packageName,
		dataType:    dataType,
	}
}

func (g *FileDescriptorGenerator) Run(f *descriptor.FileDescriptorProto) (*FileDescriptor, error) {
	types, err := g.generateMessageDescriptor(f.MessageType, nil)
	if err != nil {
		return nil, err
	}

	services := g.generateServiceDescriptor(f.Service, MessageDescriptorList(types))

	return &FileDescriptor{
		PackageName: g.packageName,
		Messages:    types,
		Services:    services,
	}, nil
}

func (g *FileDescriptorGenerator) generateMessageDescriptor(messageTypes []*descriptor.DescriptorProto, parents []MessageDescriptor) ([]MessageDescriptor, error) {
	var types []MessageDescriptor

	for _, messageType := range messageTypes {
		fields, err := g.generateMessageFieldDescriptors(messageType.GetField())
		if err != nil {
			return nil, err
		}
		newMessageType := MessageDescriptor{
			MessageName: messageType.GetName(),
			Fields:      fields,
			Parents:     parents,
		}
		types = append(types, newMessageType)
		nestedTypes, err := g.generateMessageDescriptor(messageType.NestedType, append(parents, newMessageType))
		if err != nil {
			return nil, err
		}
		types = append(types, nestedTypes...)
	}

	return types, nil
}

func (g *FileDescriptorGenerator) generateMessageFieldDescriptors(fields []*descriptor.FieldDescriptorProto) ([]MessageFieldDescriptor, error) {
	var params []MessageFieldDescriptor
	for _, field := range fields {
		typeName, err := g.dataType.GetName(field)
		if err != nil {
			return nil, err
		}

		param := MessageFieldDescriptor{
			Name:         field.GetName(),
			DataTypeName: typeName,
			IsRequired:   !field.GetProto3Optional(),
			IsOptional:   field.GetProto3Optional(),
			IsTimestamp:  field.GetTypeName() == ".google.protobuf.Timestamp",
			IsRepeated:   field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED,
		}
		params = append(params, param)
	}

	return params, nil
}

func (g *FileDescriptorGenerator) generateServiceDescriptor(services []*descriptor.ServiceDescriptorProto, messages MessageDescriptorList) []ServiceDescriptor {
	var types []ServiceDescriptor

	for _, service := range services {
		methods := g.generateServiceMethodDescriptors(service.Method, messages)
		newService := ServiceDescriptor{
			ServiceName: strings.TrimSuffix(service.GetName(), "Service"),
			Methods:     methods,
			Messages:    messages,
		}
		types = append(types, newService)
	}

	return types
}

func (g *FileDescriptorGenerator) generateServiceMethodDescriptors(methods []*descriptor.MethodDescriptorProto, messages MessageDescriptorList) []ServiceMethodDescriptor {
	var params []ServiceMethodDescriptor
	for _, method := range methods {
		param := ServiceMethodDescriptor{
			Name:   method.GetName(),
			Input:  messages.GetByMessageName(strings.TrimPrefix(method.GetInputType(), ".")),
			Output: messages.GetByMessageName(strings.TrimPrefix(method.GetOutputType(), ".")),
		}
		params = append(params, param)
	}

	return params
}
