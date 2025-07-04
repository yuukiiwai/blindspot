package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/yuukiiwai/blindspot/pkg/core"
	"github.com/yuukiiwai/blindspot/pkg/std-impl/stringlist"
)

var parserMap = map[string](func() (core.Parser, error)){
	".json": stringlist.NewRuledJsonParser,
}

func getSupportedExtensions() []string {
	extensions := make([]string, 0, len(parserMap))
	for ext := range parserMap {
		extensions = append(extensions, ext)
	}
	return extensions
}

func getParser(inputFile string) (core.Parser, error) {
	parser, ok := parserMap[path.Ext(inputFile)]
	if !ok {
		return nil, fmt.Errorf("unsupported file extension: %s\nWe can only parse %v files", path.Ext(inputFile), strings.Join(getSupportedExtensions(), ", "))
	}
	return parser()
}

func getCommandDefinition() string {
	return `
	Usage:
		blindspot -input <input_file> [OPTIONS]
		blindspot -help

	Options:
		-output string (mermaid, visjs, dot) default: mermaid
		-log-severity string (debug, info, warn, error) default: warn

	Examples:
		blindspot -input rules.json -output visjs
		blindspot -input rules.json -output dot -log-severity debug
	`
}
