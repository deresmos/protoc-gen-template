package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/iancoleman/strcase"
)

type templateBind struct {
	Type        *descriptor.DescriptorProto
	Name        string
	Fields      []*descriptor.FieldDescriptorProto
	Prefix      string
	PackageName string
	Extends     string
}

type method struct {
	Name       string
	HttpMethod string
	Path       string
	InputType  string
	OutputType string
}

type clientBind struct {
	Name         string
	Methods      []*method
	EndpointBase string
}

var (
	messageTemplate *template.Template
)

func parseReq(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var req plugin.CodeGeneratorRequest
	if err = proto.Unmarshal(buf, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func processMessageTypes(packageName, prefix string, messageTypes []*descriptor.DescriptorProto) string {
	var content string
	for _, m := range messageTypes {
		b := bytes.NewBuffer([]byte{})
		err := messageTemplate.Execute(b, templateBind{
			Name:        m.GetName(),
			Fields:      m.GetField(),
			Type:        m,
			Prefix:      prefix,
			PackageName: packageName,
		})
		if err != nil {
			panic(err)
		}
		content += b.String() + "\n"
		content += processMessageTypes(packageName, prefix+m.GetName(), m.NestedType)
	}
	return content
}

func processReq(req *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	template := ""
	suffix := ".out"
	for _, p := range strings.Split(req.GetParameter(), ",") {
		spec := strings.SplitN(p, "=", 2)
		if len(spec) == 1 {
			continue
		}
		name, value := spec[0], spec[1]
		switch name {
		case "template":
			template = value
		case "suffix":
			suffix = value
		}
	}
	initTemplate(template)

	files := make(map[string]*descriptor.FileDescriptorProto)
	for _, f := range req.ProtoFile {
		files[f.GetName()] = f
	}
	var resp plugin.CodeGeneratorResponse
	for _, fname := range req.FileToGenerate {
		f := files[fname]
		outFile := fname + suffix
		content := processMessageTypes(f.GetPackage(), "", f.MessageType)
		resp.File = append(resp.File, &plugin.CodeGeneratorResponse_File{
			Name:    &outFile,
			Content: proto.String(content),
		})
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

func isPrimitive(f *descriptor.FieldDescriptorProto) bool {
	switch f.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_STRING,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_SINT64,
		descriptor.FieldDescriptorProto_TYPE_BOOL,
		descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return true
	}
	return false
}

func getType(f *descriptor.FieldDescriptorProto, packageName string) string {
	switch f.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "string"
	case descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT32:
		return "int32"
	case descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_SINT64:
		return "int64"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "bool"
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return "double"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		t := strings.Replace(f.GetTypeName(), "."+packageName+".", "", -1)
		t = strings.Replace(t, ".", "", -1)
		return t
	}
	return "unknown"
}

func initTemplate(file string) {
	var err error
	var buf []byte
	buf, err = ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	messageTemplate, err = template.New("apex").Funcs(template.FuncMap{
		"isPrimitive": isPrimitive,
		"getName":     strcase.ToCamel,
		"toLower":     strcase.ToLowerCamel,
		"isRepeated": func(f *descriptor.FieldDescriptorProto) bool {
			return f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED
		},
		"getSingleType": getType,
		"propertyType": func(f *descriptor.FieldDescriptorProto, packageName string) string {
			format := "%s"
			if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
				format = "[]%s"
			}
			return fmt.Sprintf(format, getType(f, packageName))
		},
	}).Parse(string(buf))
	if err != nil {
		panic(err)
	}
}

func run() error {
	req, err := parseReq(os.Stdin)
	if err != nil {
		return err
	}

	resp := processReq(req)

	return emitResp(resp)
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
