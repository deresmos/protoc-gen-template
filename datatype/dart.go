package datatype

import (
	"fmt"
	"regexp"
	"strings"

	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

type DartDataType struct{}

func (t DartDataType) getTypeName(f *descriptor.FieldDescriptorProto) (string, error) {
	switch f.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "String", nil
	case descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_SINT64:
		return "int", nil
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return "double", nil
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "bool", nil
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		switch f.GetTypeName() {
		case ".google.protobuf.Timestamp":
			return "DateTime", nil
		default:
			if strings.Contains(f.GetTypeName(), "__") {
				re := regexp.MustCompile(`\.__(\w+)$`)
				return re.FindStringSubmatch(f.GetTypeName())[1], nil
			}

			return strings.TrimPrefix(f.GetTypeName(), "."), nil
		}
	}

	return "", fmt.Errorf("unknown type: %s", f.GetType())
}

func (t DartDataType) repeatedFormat() string {
	return "List<%s>"
}
