---
layout: default
title: Bash/Shell
parent: 言語ガイド
nav_order: 13
---

# Bash/Shell対応

[English](../../languages/bash.md) | [한국어](../../ko/languages/bash.md) | [日本語](bash.md) | [हिन्दी](../../hi/languages/bash.md) | [Deutsch](../../de/languages/bash.md)

## 対応拡張子

- `.sh`
- `.bash`

## 抽出対象

| 要素 | 種類 | 例 |
|------|------|-----|
| 関数 | `function` | `function greet { ... }` |
| 関数 | `function` | `greet() { ... }` |
| 変数代入 | `variable` | `NAME="value"` |
| 宣言 | `variable` | `declare VERBOSE` |
| ローカル変数 | `variable` | `local count=0` |
| 読み取り専用変数 | `variable` | `readonly VERSION="1.0"` |
| コメント | `doc` | `# 説明` |
| source文 | `import` | `source /path/to/lib.sh` |
| ドット文 | `import` | `. ./config.sh` |

## 例

### 入力

```bash
#!/bin/bash

# 設定
CONFIG_PATH="/etc/myapp"
VERSION="1.0.0"
declare VERBOSE=false

# アプリケーションをデプロイ
function deploy {
    local app_name="$1"
    echo "Deploying $app_name"
}

# プロジェクトをビルド
build() {
    echo "Building..."
}

source ./utils.sh
. ./config.sh
```

### 出力 (XML)

```xml
<file path="deploy.sh" language="bash">
  <variable kind="variable" line="4">
    <name>CONFIG_PATH</name>
    <text>CONFIG_PATH="/etc/myapp"</text>
  </variable>
  <variable kind="variable" line="5">
    <name>VERSION</name>
    <text>VERSION="1.0.0"</text>
  </variable>
  <variable kind="variable" line="6">
    <name>VERBOSE</name>
    <text>declare VERBOSE=false</text>
  </variable>
  <function kind="function" line="9">
    <name>deploy</name>
    <text>function deploy</text>
  </function>
  <function kind="function" line="15">
    <name>build</name>
    <text>build()</text>
  </function>
</file>
```

## 注意事項

### 可視性

- すべての宣言が抽出されます（Bashにはアクセス修飾子がありません）
- 関数内の`local`変数も解析時に宣言されていれば抽出されます

### 関数構文

Bashは2つの関数宣言スタイルをサポートしています：

- `function 名前 { ... }` - `function`キーワードを使用
- `名前() { ... }` - 括弧を使用

どちらも`function`種類として抽出されます。

### 本体の削除

`--include-body`フラグを使用しない場合：

- 関数：開き中括弧`{`以降の本体を削除
- 変数：最初の行のみ保持（複数行代入を処理）

### インポート抽出

- `source`および`.`コマンドが`--include-imports`フラグで抽出されます
- 引用符あり/なしのパスの両方をサポートします

### ドキュメントコメント

- `#`で始まるシェルコメントが抽出されます
- シバン行（`#!/bin/bash`）はコメントとして扱われません

### 制限事項

- ネストされた関数がサポートされています
- 関数本体内のヒアドキュメントも正しく処理されます
- 複雑な変数展開はシグネチャにそのまま保持されます
