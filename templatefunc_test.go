package main_test

import (
	templatefunc "github.com/deresmos/protoc-gen-template"
	"github.com/gertd/go-pluralize"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToSingular(t *testing.T) {
	tf := templatefunc.NewTemplateFunc(pluralize.NewClient())

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "tests -> test",
			input: "tests",
			want:  "test",
		},
		{
			name:  "TestUsers -> TestUser",
			input: "TestUsers",
			want:  "TestUser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tf.ToSingular(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tf := templatefunc.NewTemplateFunc(pluralize.NewClient())

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "camelCase -> camel_case",
			input: "camelCase",
			want:  "camel_case",
		},
		{
			name:  "lowerCamelCase -> lower_camel_case",
			input: "lowerCamelCase",
			want:  "lower_camel_case",
		},
		{
			name:  "snake_case -> snake_case",
			input: "snake_case",
			want:  "snake_case",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tf.ToSnakeCase(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tf := templatefunc.NewTemplateFunc(pluralize.NewClient())

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "CamelCase to CamelCase",
			input: "CamelCase",
			want:  "CamelCase",
		},
		{
			name:  "snake_case -> SnakeCase",
			input: "snake_case",
			want:  "SnakeCase",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tf.ToCamelCase(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToLowerCamelCase(t *testing.T) {
	tf := templatefunc.NewTemplateFunc(pluralize.NewClient())

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "CamelCase -> camelCase",
			input: "CamelCase",
			want:  "camelCase",
		},
		{
			name:  "lowerCamelCase -> lowerCamelCase",
			input: "lowerCamelCase",
			want:  "lowerCamelCase",
		},
		{
			name:  "snake_case -> snakeCase",
			input: "snake_case",
			want:  "snakeCase",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tf.ToLowerCamelCase(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReplace(t *testing.T) {
	tf := templatefunc.NewTemplateFunc(pluralize.NewClient())

	tests := []struct {
		name string
		old  string
		new  string
		src  string
		want string
	}{
		{
			name: "Replace: foo -> bar",
			old:  "foo",
			new:  "bar",
			src:  "foofoo",
			want: "barbar",
		},
		{
			name: "Replace: hello -> world",
			old:  "hello",
			new:  "world",
			src:  "hellohello",
			want: "worldworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tf.Replace(tt.old, tt.new, tt.src)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContains(t *testing.T) {
	tf := templatefunc.NewTemplateFunc(pluralize.NewClient())

	tests := []struct {
		name   string
		s      string
		substr string
		want   bool
	}{
		{
			name:   "String contains substring",
			s:      "hello world",
			substr: "world",
			want:   true,
		},
		{
			name:   "String does not contain substring",
			s:      "hello world",
			substr: "goodbye",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tf.Contains(tt.s, tt.substr)
			assert.Equal(t, tt.want, got)
		})
	}
}
