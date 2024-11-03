package main

import (
	"slices"
	"strings"

	"github.com/deresmos/protoc-gen-template/datatype"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

type FileDescriptor struct {
	PackageName string
	Messages    []MessageDescriptor
	Services    []ServiceDescriptor
}

func (f *FileDescriptor) Append(fileDescriptor *FileDescriptor) *FileDescriptor {
	if f == nil {
		return fileDescriptor
	}

	return &FileDescriptor{
		PackageName: f.PackageName,
		Messages:    append(f.Messages, fileDescriptor.Messages...),
		Services:    append(f.Services, fileDescriptor.Services...),
	}
}

type MessageDescriptor struct {
	MessageName  string
	Fields       MessageFieldDescriptorList
	Parents      MessageDescriptorList
	ItemMessages MessageDescriptorList
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

type MessageFieldDescriptor struct {
	FieldName     string
	DataTypeName  string
	IsOptional    bool
	IsRequired    bool
	IsTimestamp   bool
	IsRepeated    bool
	IsMessageType bool
}

type MessageFieldDescriptorList []MessageFieldDescriptor

func (m MessageFieldDescriptorList) HasTimestamp() bool {
	return slices.ContainsFunc(m, func(field MessageFieldDescriptor) bool {
		return field.IsTimestamp
	})
}

type ServiceDescriptor struct {
	ServiceName string
	Methods     []ServiceMethodDescriptor
	Messages    MessageDescriptorList
}

type ServiceMethodDescriptor struct {
	MethodName    string
	InputMessage  *MessageDescriptor
	OutputMessage *MessageDescriptor
	Dependencies  []MessageDescriptor
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
		nestedTypes, err := g.generateMessageDescriptor(messageType.NestedType, append(parents, newMessageType))
		if err != nil {
			return nil, err
		}

		var itemMessages MessageDescriptorList
		for _, nestedType := range nestedTypes {
			if strings.HasPrefix(nestedType.MessageName, "__") {
				itemMessages = append(itemMessages, nestedType)
				continue
			}
		}

		var filteredNestedMessages []MessageDescriptor
		for _, nestedType := range nestedTypes {
			if strings.HasPrefix(nestedType.MessageName, "__") {
				continue
			}
			filteredNestedMessages = append(filteredNestedMessages, nestedType)
		}
		newMessageType.ItemMessages = itemMessages
		types = append(types, newMessageType)
		types = append(types, filteredNestedMessages...)
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
			FieldName:     field.GetName(),
			DataTypeName:  typeName,
			IsOptional:    field.GetProto3Optional(),
			IsRequired:    !field.GetProto3Optional(),
			IsTimestamp:   field.GetTypeName() == ".google.protobuf.Timestamp",
			IsRepeated:    field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED,
			IsMessageType: field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE,
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
			MethodName:    method.GetName(),
			InputMessage:  messages.GetByMessageName(strings.TrimPrefix(method.GetInputType(), ".")),
			OutputMessage: messages.GetByMessageName(strings.TrimPrefix(method.GetOutputType(), ".")),
		}
		params = append(params, param)
	}

	return params
}
