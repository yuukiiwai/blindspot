# blindspot - Generate Finite Automaton
ルールに則って状態遷移図を自動生成することで、盲点となっている状態を発見し、遷移ルール作成を手助けする。  
Automatic generation of state transition diagrams according to the rules helps to find blind spot states and create transition rules.
## 目次
- [日本語パート - Japanese](#日本語パート)
- [English - 英語パート](#english)

# 日本語パート
## インストール方法
```sh
$ go install github.com/yuukiiwai/blindspot/cmd/cli/blindspot@latest
```

## 使用方法
1. 状態遷移jsonを作成する
2. コマンドを実行する
    ```sh
    $ blindspot -input data.json -output mermaid
    ```

## コントリビューター向け
### 解決する課題
![初期案](./1st-design.jpg)

# English
## Installation
```sh
$ go install github.com/yuukiiwai/blindspot/cmd/cli/blindspot@latest
```

## Usage
1. make state json
2. execute under command
    ```sh
    $ blindspot -input data.json -output mermaid
    ```

## For Contributers
## What I want to resolve
![1st disign](./1st-design.jpg)
