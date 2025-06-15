package output

import (
	"errors"

	"github.com/yuukiiwai/blindspot/pkg/core"
)

// DotFormatter Graphviz（DOT）形式の出力フォーマッター
type DotFormatter struct{}

// NewDotFormatter 新しいDotFormatterを作成
func NewDotFormatter() *DotFormatter {
	return &DotFormatter{}
}

// Format ステートマシンをDOT形式で出力
func (f *DotFormatter) Format(generator *core.Generator) (string, error) {
	// TODO: Implement
	return "", errors.New("not implemented")
}
