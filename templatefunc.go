package main

import (
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
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

func (t TemplateFunc) HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func (t TemplateFunc) HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
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
		"toKebab":          templateFunc.ToKebab,
		"toLowerCamelCase": templateFunc.ToLowerCamelCase,
		"toSnakeCase":      templateFunc.ToSnakeCase,
		"toSingular":       templateFunc.ToSingular,
		"toPlural":         templateFunc.ToPlural,
	}).Funcs(sprig.FuncMap()).Parse(string(buf))
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
		"hasPrefix":        templateFunc.HasPrefix,
		"hasSuffix":        templateFunc.HasSuffix,
	}).Funcs(sprig.FuncMap()).Parse(outputPath)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
