package datatype

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type TypeScriptDataType struct{}

func (t TypeScriptDataType) getTypeName(f *descriptor.FieldDescriptorProto) (string, error) {
	switch f.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "string", nil
	case descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_SINT64:
		return "number", nil
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return "bigint", nil
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "boolean", nil
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		switch f.GetTypeName() {
		case ".google.protobuf.Timestamp":
			return "Date", nil
		default:
			return strings.TrimPrefix(f.GetTypeName(), "."), nil
		}
	}

	return "", fmt.Errorf("unknown type: %s", f.GetType())
}

func (t TypeScriptDataType) repeatedFormat() string {
	return "%s[]"
}
