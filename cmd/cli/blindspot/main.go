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
	// FlagSetを使用して混合引数を処理
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var help bool
	fs.BoolVar(&help, "help", false, "ヘルプを表示")
	inputFormat := fs.String("input", "stringlist", "入力形式 (stringlist, cud)")
	outputFormat := fs.String("output", "mermaid", "出力形式 (mermaid, visjs, dot)")
	logSeverity := fs.String("log-severity", "warn", "ログの重大度 (debug, info, warn, error)")
	limitFlag := fs.Int64("limit", -1, "反復回数の上限")

	// 引数が足りない場合はヘルプを表示
	if len(os.Args) < 2 {
		fmt.Println(getCommandDefinition())
		os.Exit(1)
	}

	// ヘルプが最初の引数の場合
	if os.Args[1] == "-help" || os.Args[1] == "--help" {
		fmt.Println(getCommandDefinition())
		os.Exit(0)
	}

	// 最初の引数を入力ファイルとして取得
	inputFile := os.Args[1]

	// 残りの引数をフラグとして解析
	fs.Parse(os.Args[2:])

	// limitが指定されていない場合はnilポインタを使用
	var limit *int64
	if *limitFlag != -1 {
		limit = limitFlag
	}

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

	// 入力ファイルの読み込み
	ruleFile, err := os.ReadFile(inputFile)
	if err != nil {
		slog.Error("入力ファイルの読み込みに失敗", "error", err)
		os.Exit(1)
	}

	// パーサーの作成
	parser, err := getParser(*inputFormat)
	if err != nil {
		slog.Error("パーサーの作成に失敗", "error", err)
		os.Exit(1)
	}

	// ルールのパース
	firstResources, newNode, edgeRules, err := parser.Parse(string(ruleFile))
	if err != nil {
		slog.Error("ルールのパースに失敗", "error", err)
		os.Exit(1)
	}

	if limit != nil {
		fmt.Printf("反復回数の上限は%dです.よろしいですか？(y/n)", *limit)
		var input string
		fmt.Scanln(&input)
		if input != "y" {
			os.Exit(0)
		}
	} else {
		fmt.Println("反復回数の上限は設定されていません.論理的に終了条件が存在しない場合,コンピューターに不具合が発生する可能性があります.また,生成後のステートマシンの出力中に停止する可能性があります.よろしいですか？(y/n)")
		var input string
		fmt.Scanln(&input)
		if input != "y" {
			os.Exit(0)
		}
	}

	// ジェネレーターの作成
	generator := core.NewGenerator(newNode, firstResources, edgeRules, limit)

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
