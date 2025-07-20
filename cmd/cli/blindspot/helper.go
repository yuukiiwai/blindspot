package main

import (
	"fmt"
	"strings"

	"github.com/yuukiiwai/blindspot/pkg/core"
	"github.com/yuukiiwai/blindspot/pkg/std-impl/cud"
	"github.com/yuukiiwai/blindspot/pkg/std-impl/stringlist"
)

var parserMap = map[string](func() (core.Parser, error)){
	"stringlist": stringlist.NewRuledJsonParser,
	"cud":        cud.NewCudYamlParser,
}

func getSupportedFormats() []string {
	formats := make([]string, 0, len(parserMap))
	for format := range parserMap {
		formats = append(formats, format)
	}
	return formats
}

func getParser(inputFormat string) (core.Parser, error) {
	parser, ok := parserMap[inputFormat]
	if !ok {
		return nil, fmt.Errorf("unsupported input format: %s\nWe can only parse %v formats", inputFormat, strings.Join(getSupportedFormats(), ", "))
	}
	return parser()
}

func getCommandDefinition() string {
	return `
	Usage:
		blindspot <input_file> [OPTIONS]
		blindspot -help

	Required:
		<input_file> string (入力ファイルのパス)

	Options:
		-input string (stringlist, cud) default: stringlist
		-output string (mermaid, visjs, dot) default: mermaid
		-log-severity string (debug, info, warn, error) default: warn
		--limit int64 (反復回数の上限、無限ループ防止) default: 0 (無制限)

	Examples:
		blindspot rules.json -input stringlist -output mermaid
		blindspot rules.json -input cud -output visjs
		blindspot rules.json -input stringlist -output dot -log-severity debug
		blindspot rules.json -input stringlist -output mermaid --limit 1000
	`
}
