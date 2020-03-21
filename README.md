## protoc-gen-template

A protoc plugin to output file with text/template literal

## Install

```bash
$ go get github.com/tzmfreedom/protoc-gen-template
```

## Usage

generate class files with following command
```bash
$ protoc -I. --template_out=template=go.template:. target.proto
```

You should write template file with [text/template](https://golang.org/pkg/text/template/)
```
type {{ .Prefix }}{{ .Type.Name }} struct {
    {{ range .Type.Field }}{{ getName .Name }} {{ propertyType . $.PackageName }}
    {{ end }}
}
```
