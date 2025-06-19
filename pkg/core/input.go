package core

// Parser は入力ファイルをパースするインターフェース
type Parser interface {
	// Parse は入力ファイルをパースして、開始リソースとルールを返す
	Parse(input string) (
		firstResource Node,
		newNode func(resources any) Node,
		edgeRules []*EdgeRule,
		err error,
	)
}
