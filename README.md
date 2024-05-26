## protoc-gen-template

A protoc plugin to output file with text/template literal

## Install

```bash
go install github.com/deresmos/protoc-gen-template
```

## Usage

generate class files with following command
```bash
protoc --template_out='template=ts.template,lang=typescript,generate_type=message,output_path=./{{toSnakeCase .MessageName}}.ts:.' schema.proto
```

You should write template file with [text/template](https://golang.org/pkg/text/template/)

### TypeScript

```
export interface {{toSingular .MessageName}}Doc {
    {{ range .Fields }}{{ toLowerCamelCase .Name }}: {{ .DataTypeName }},
    {{ end }}
}
```

### Dart

```
@freezed
class {{toSingular .MessageName}} with _${{toSingular .MessageName}} {
  const {{toSingular .MessageName}}._();
  const factory {{toSingular .MessageName}}({
    {{ range .Fields }}{{if .IsRequired}}required {{end}}{{ .DataTypeName }} {{ toLowerCamelCase .Name }},
    {{ end }}
  }) = _{{toSingular .MessageName}};

  factory {{toSingular .MessageName}}.fromJson(Map<String, dynamic> json) => _${{toSingular .MessageName}}FromJson(json);
}
```
