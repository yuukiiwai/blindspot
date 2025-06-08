package input

import "github.com/yuukiiwai/blindspot/pkg/core"

// Parser は入力ファイルをパースするインターフェース
type Parser interface {
	// Parse は入力ファイルをパースして、開始リソースとルールを返す
	Parse(input string) ([]string, []*core.EdgeRule, error)
}
