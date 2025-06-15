package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/yuukiiwai/blindspot/pkg/core"
	"github.com/yuukiiwai/blindspot/pkg/std-impl/output"
)

func main() {
	// コマンドライン引数の定義
	var help bool
	flag.BoolVar(&help, "help", false, "ヘルプを表示")
	inputFile := flag.String("input", "", "入力ファイルのパス")
	outputFormat := flag.String("output", "mermaid", "出力形式 (mermaid, visjs, dot)")
	logSeverity := flag.String("log-severity", "warn", "ログの重大度 (debug, info, warn, error)")
	flag.Parse()

	if help {
		fmt.Println(getCommandDefinition())
		os.Exit(0)
	}

	// ログの重大度の設定
	var level slog.Level
	switch *logSeverity {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelWarn
	}
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// 入力ファイルの確認
	if *inputFile == "" {
		slog.Error("入力ファイルを指定してください (-input)")
		os.Exit(1)
	}

	// 入力ファイルの読み込み
	ruleFile, err := os.ReadFile(*inputFile)
	if err != nil {
		slog.Error("入力ファイルの読み込みに失敗", "error", err)
		os.Exit(1)
	}

	// パーサーの作成
	parser, err := getParser(*inputFile)
	if err != nil {
		slog.Error("パーサーの作成に失敗", "error", err)
		os.Exit(1)
	}

	// ルールのパース
	firstResources, edgeRules, err := parser.Parse(string(ruleFile))
	if err != nil {
		slog.Error("ルールのパースに失敗", "error", err)
		os.Exit(1)
	}

	// ジェネレーターの作成
	generator := core.NewGenerator(firstResources, edgeRules)

	// ステートマシンの生成
	if err := generator.Generate(); err != nil {
		slog.Error("ステートマシンの生成に失敗", "error", err)
		os.Exit(1)
	}

	// フォーマッターの選択
	var formatter core.Formatter
	switch *outputFormat {
	case "mermaid":
		formatter = output.NewMermaidFormatter()
	case "visjs":
		formatter = output.NewVisjsFormatter()
	case "dot":
		formatter = output.NewDotFormatter()
	default:
		slog.Error("未対応の出力形式", "format", *outputFormat)
		os.Exit(1)
	}

	// 出力の生成
	result, err := formatter.Format(generator)
	if err != nil {
		slog.Error("出力の生成に失敗", "error", err)
		os.Exit(1)
	}

	// 結果の出力
	fmt.Println(result)
}
