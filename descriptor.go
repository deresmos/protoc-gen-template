package main

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/deresmos/protoc-gen-template/datatype"
)

type FileDescriptor struct {
	PackageName string
	Messages    []MessageDescriptor
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

func (g *FileDescriptorGenerator) Run(messageTypes []*descriptor.DescriptorProto) (*FileDescriptor, error) {
	types, err := g.generateMessageDescriptor(messageTypes, nil)
	if err != nil {
		return nil, err
	}

	return &FileDescriptor{
		PackageName: g.packageName,
		Messages:    types,
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
