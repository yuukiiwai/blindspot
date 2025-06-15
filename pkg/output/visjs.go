package output

import (
	"errors"

	"github.com/yuukiiwai/blindspot/pkg/core"
)

// VisjsFormatter Vis.js形式の出力フォーマッター
type VisjsFormatter struct{}

// NewVisjsFormatter 新しいVisjsFormatterを作成
func NewVisjsFormatter() *VisjsFormatter {
	return &VisjsFormatter{}
}

// Format ステートマシンをVis.js形式で出力
func (f *VisjsFormatter) Format(generator *core.Generator) (string, error) {
	// TODO: Implement
	return "", errors.New("not implemented")
}
