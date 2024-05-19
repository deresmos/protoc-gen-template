package main

import (
	"fmt"
	"strings"
)

type ProtoOption struct {
	TemplatePath string
	Language     string
	FileNameExt  string
	OutputPath   string
	GenerateType string
}

func NewProtoOptionFromString(protoOption string) (*ProtoOption, error) {
	templatePath, err := parseProtoOption(protoOption, "template")
	if err != nil {
		return nil, err
	}
	language, err := parseProtoOption(protoOption, "lang")
	if err != nil {
		return nil, err
	}
	outputDirectory, err := parseProtoOption(protoOption, "output")
	if err != nil {
		return nil, err
	}
	generateType, err := parseProtoOption(protoOption, "generate_type")
	if err != nil {
		return nil, err
	}

	return &ProtoOption{
		TemplatePath: templatePath,
		Language:     language,
		FileNameExt:  getFileNameExtFromLanguage(language),
		OutputPath:   outputDirectory,
		GenerateType: generateType,
	}, nil
}

func parseProtoOption(optionString string, fieldName string) (string, error) {
	spec := strings.Split(optionString, ",")
	for _, p := range spec {
		if strings.Contains(p, fieldName) {
			return strings.Split(p, "=")[1], nil
		}
	}

	return "", fmt.Errorf("option `%s` not found", fieldName)
}

func getFileNameExtFromLanguage(language string) string {
	switch language {
	case "dart":
		return "dart"
	case "typescript":
		return "ts"
	default:
		return "txt"
	}
}
