package core

// Formatter 出力フォーマッターのインターフェース
type Formatter interface {
	// Format ステートマシンを指定された形式で出力
	Format(generator *Generator) (string, error)
}
