package output

import "github.com/yuukiiwai/blindspot/pkg/core"

// Formatter 出力フォーマッターのインターフェース
type Formatter interface {
	// Format ステートマシンを指定された形式で出力
	Format(generator *core.Generator) (string, error)
}
