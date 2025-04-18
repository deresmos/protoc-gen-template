package datatype

import (
	"fmt"

	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

type dataType interface {
	repeatedFormat() string
	getTypeName(f *descriptor.FieldDescriptorProto) (string, error)
}

type DataType struct {
	dataType dataType
}

func (d DataType) GetName(f *descriptor.FieldDescriptorProto) (string, error) {
	format := "%s"
	if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
		format = d.dataType.repeatedFormat()
	}

	typeName, err := d.dataType.getTypeName(f)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(format, typeName), nil
}

func factoryDataType(lang string) (dataType, error) {
	switch lang {
	case "typescript":
		return &TypeScriptDataType{}, nil
	case "dart":
		return &DartDataType{}, nil
	case "go":
		return &GoDataType{}, nil
	}

	return nil, fmt.Errorf("unknown language: %s", lang)
}

func NewDataType(lang string) (*DataType, error) {
	dataType, err := factoryDataType(lang)
	if err != nil {
		return nil, err
	}

	return &DataType{
		dataType: dataType,
	}, nil
}
