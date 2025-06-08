package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/yuukiiwai/blindspot/pkg/input"
)

var parserMap = map[string](func() (input.Parser, error)){
	".tengo": input.NewTengoParser,
}

func getSupportedExtensions() []string {
	extensions := make([]string, 0, len(parserMap))
	for ext := range parserMap {
		extensions = append(extensions, ext)
	}
	return extensions
}

func getParser(inputFile string) (input.Parser, error) {
	parser, ok := parserMap[path.Ext(inputFile)]
	if !ok {
		return nil, fmt.Errorf("unsupported file extension: %s\nWe can only parse %v files", path.Ext(inputFile), strings.Join(getSupportedExtensions(), ", "))
	}
	return parser()
}
