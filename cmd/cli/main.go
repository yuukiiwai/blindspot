package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yuukiiwai/blindspot/pkg/core"
	"github.com/yuukiiwai/blindspot/pkg/output"
)

func main() {
	// コマンドライン引数の定義
	inputFile := flag.String("input", "", "入力ファイルのパス")
	outputFormat := flag.String("format", "mermaid", "出力形式 (mermaid, visjs, dot)")
	flag.Parse()

	// 入力ファイルの確認
	if *inputFile == "" {
		log.Fatal("入力ファイルを指定してください (-input)")
	}

	// 入力ファイルの読み込み
	ruleFile, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("入力ファイルの読み込みに失敗: %v", err)
	}

	// パーサーの作成
	parser, err := getParser(*inputFile)
	if err != nil {
		log.Fatalf("パーサーの作成に失敗: %v", err)
	}

	// ルールのパース
	firstResources, edgeRules, err := parser.Parse(string(ruleFile))
	if err != nil {
		log.Fatalf("ルールのパースに失敗: %v", err)
	}

	// ジェネレーターの作成
	generator := core.NewGenerator(firstResources, edgeRules)

	// ステートマシンの生成
	if err := generator.Generate(); err != nil {
		log.Fatalf("ステートマシンの生成に失敗: %v", err)
	}

	// フォーマッターの選択
	var formatter output.Formatter
	switch *outputFormat {
	case "mermaid":
		formatter = output.NewMermaidFormatter()
	case "visjs":
		formatter = output.NewVisjsFormatter()
	case "dot":
		formatter = output.NewDotFormatter()
	default:
		log.Fatalf("未対応の出力形式: %s", *outputFormat)
	}

	// 出力の生成
	result, err := formatter.Format(generator)
	if err != nil {
		log.Fatalf("出力の生成に失敗: %v", err)
	}

	// 結果の出力
	fmt.Println(result)
}
