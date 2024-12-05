package main

import (
	"os"
	"strings"
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

func (t TemplateFunc) ToSnakeCase(s string) string {
	return strcase.ToSnake(s)
}

func (t TemplateFunc) ToKebab(s string) string {
	return strcase.ToKebab(s)
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

func (t TemplateFunc) ToPlural(s string) string {
	return t.pluarizerClient.Plural(s)
}

func (t TemplateFunc) Replace(old, new, src string) string {
	return strings.Replace(src, old, new, -1)
}

func (t TemplateFunc) Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func initFileTemplate(file string) (*template.Template, error) {
	var err error
	var buf []byte
	buf, err = os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	templateFunc := NewTemplateFunc(pluralize.NewClient())
	tmpl, err := template.New("gen-protoc").Funcs(template.FuncMap{
		"toCamelCase":      templateFunc.ToCamelCase,
		"toKebeb":          templateFunc.ToKebab,
		"toLowerCamelCase": templateFunc.ToLowerCamelCase,
		"toSnakeCase":      templateFunc.ToSnakeCase,
		"toSingular":       templateFunc.ToSingular,
		"toPlural":         templateFunc.ToPlural,
		"replace":          templateFunc.Replace,
		"contains":         templateFunc.Contains,
	}).Parse(string(buf))
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func initOutputPathTemplate(outputPath string) (*template.Template, error) {
	templateFunc := NewTemplateFunc(pluralize.NewClient())
	tmpl, err := template.New("gen-protoc-output-path").Funcs(template.FuncMap{
		"toCamelCase":      templateFunc.ToCamelCase,
		"toLowerCamelCase": templateFunc.ToLowerCamelCase,
		"toSnakeCase":      templateFunc.ToSnakeCase,
		"toSingular":       templateFunc.ToSingular,
		"toPlural":         templateFunc.ToPlural,
		"replace":          strings.Replace,
		"contains":         templateFunc.Contains,
	}).Parse(outputPath)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
