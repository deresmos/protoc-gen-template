package main

import (
	"fmt"
	"strings"
)

type ProtoOption struct {
	TemplatePath string
	Language     string
	OutputPath   string
	GenerateType string
	AllowMerge   bool
	Overwrite    bool
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
	outputDirectory, err := parseProtoOption(protoOption, "output_path")
	if err != nil {
		return nil, err
	}
	generateType, err := parseProtoOption(protoOption, "generate_type")
	if err != nil {
		return nil, err
	}
	allowMerge := parseOptionalOption(protoOption, "allow_merge")
	overwite := parseOptionalOption(protoOption, "overwite")

	return &ProtoOption{
		TemplatePath: templatePath,
		Language:     language,
		OutputPath:   outputDirectory,
		GenerateType: generateType,
		AllowMerge:   allowMerge == "true",
		Overwrite:    overwite != "false",
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

func parseOptionalOption(optionString string, fieldName string) string {
	spec := strings.Split(optionString, ",")
	for _, p := range spec {
		if strings.HasPrefix(p, fieldName) {
			return strings.Split(p, "=")[1]
		}
	}

	return ""
}
