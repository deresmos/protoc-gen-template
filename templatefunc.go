package main

import (
	"os"
	"text/template"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

type TemplateFunc struct {
	pluarizerClient *pluralize.Client
}

func NewTemplateFunc(pluarizerClient *pluralize.Client) TemplateFunc {
	return TemplateFunc{
		pluarizerClient: pluarizerClient,
	}
}

func (t TemplateFunc) ToSingularLowerCamelCase(s string) string {
	return t.pluarizerClient.Singular(strcase.ToLowerCamel(s))
}

func (t TemplateFunc) ToSnakeCase(s string) string {
	return strcase.ToSnake(s)
}

func (t TemplateFunc) ToCamelCase(s string) string {
	return strcase.ToCamel(s)
}

func (t TemplateFunc) ToLowerCamelCase(s string) string {
	return strcase.ToLowerCamel(s)
}

func (t TemplateFunc) ToSingular(s string) string {
	return t.pluarizerClient.Singular(s)
}

func initTemplate(file string) (*template.Template, error) {
	var err error
	var buf []byte
	buf, err = os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	templateFunc := NewTemplateFunc(pluralize.NewClient())
	tmpl, err := template.New("gen-protoc").Funcs(template.FuncMap{
		"toCamelCase":              templateFunc.ToCamelCase,
		"toLowerCamelCase":         templateFunc.ToLowerCamelCase,
		"toSnakeCase":              templateFunc.ToSnakeCase,
		"toSingularLowerCamelCase": templateFunc.ToSingularLowerCamelCase,
		"toSingular":               templateFunc.ToSingular,
	}).Parse(string(buf))
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
