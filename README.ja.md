# Brf.it

🌐 [English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md) | [हिन्दी](README.hi.md) | [Deutsch](README.de.md)

[![Release](https://img.shields.io/github/v/release/indigo-net/Brf.it)](https://github.com/indigo-net/Brf.it/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/indigo-net/Brf.it)](https://goreportcard.com/report/github.com/indigo-net/Brf.it)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**AIアシスタントのためのコードブリーフィングツール**

Brf.itはコードベースから関数シグネチャを抽出し、実装の詳細を削除してAIが必要とする重要な情報のみを残します。トークン使用量を大幅に削減できます。

---

## 主な機能

AIアシスタントに生のコードを渡す代わりに：

<table>
<tr>
<th>Before（50+トークン）</th>
<th>After with brfit（8トークン）</th>
</tr>
<tr>
<td>

```typescript
export async function fetchUser(
  id: string
): Promise<User> {
  const response = await fetch(
    `${API_URL}/users/${id}`
  );
  if (!response.ok) {
    throw new Error('User not found');
  }
  const data = await response.json();
  return {
    id: data.id,
    name: data.name,
    email: data.email,
    createdAt: new Date(data.created_at)
  };
}
```

</td>
<td>

```xml
<signature>
  export async function fetchUser(
    id: string
  ): Promise<User>
</signature>
```

</td>
</tr>
</table>

---

## インストール

### macOS（Homebrew）

```bash
brew install indigo-net/tap/brfit
```

### リリースからダウンロード

[Releases](https://github.com/indigo-net/Brf.it/releases)から最新のバイナリをダウンロードしてください。

### ソースからビルド

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o brfit ./cmd/brfit
```

---

## 使い方

```bash
brfit [パス] [オプション]
```

### クイック例

```bash
# 現在のディレクトリからシグネチャを抽出
brfit .

# Markdown形式で出力
brfit . -f md

# ファイルに保存
brfit . -o output.xml

# 関数本体を含める（完全なコード）
brfit . --include-body

# ディレクトリツリーをスキップ
brfit . --no-tree
```

### CLIオプション

| オプション | 短縮形 | 説明 | デフォルト |
|------------|--------|------|------------|
| `--format` | `-f` | 出力形式（`xml`、`md`） | `xml` |
| `--output` | `-o` | 出力ファイルパス | stdout |
| `--include-body` | | 関数本体を含める | `false` |
| `--ignore` | `-i` | ignoreファイルパス | `.gitignore` |
| `--include-hidden` | | 隠しファイルを含める | `false` |
| `--no-tree` | | ディレクトリツリーをスキップ | `false` |
| `--no-tokens` | | トークンカウントを無効化 | `false` |
| `--max-size` | | 最大ファイルサイズ（バイト） | `512000` |
| `--version` | `-v` | バージョンを表示 | |

---

## サポート言語

| 言語 | 拡張子 | ドキュメント |
|------|--------|--------------|
| Go | `.go` | [Goガイド](docs/languages/go.ja.md) |
| TypeScript | `.ts`、`.tsx` | [TypeScriptガイド](docs/languages/typescript.ja.md) |
| JavaScript | `.js`、`.jsx` | [TypeScriptガイド](docs/languages/typescript.ja.md) |
| Python | `.py` | [Pythonガイド](docs/languages/python.ja.md) |

---

## 出力例

### XML（デフォルト）

```xml
<?xml version="1.0" encoding="UTF-8"?>
<brfit>
  <metadata>
    <tree>pkg/
└── scanner/
    └── scanner.go</tree>
  </metadata>
  <files>
    <file path="pkg/scanner/scanner.go" language="go">
      <signature>func Scan(root string) (*Result, error)</signature>
      <doc>Scan recursively scans the directory.</doc>
    </file>
  </files>
</brfit>
```

### Markdown

```markdown
# Brf.it Output

## Directory Tree

pkg/
└── scanner/
    └── scanner.go

## Files

### pkg/scanner/scanner.go

\`\`\`go
func Scan(root string) (*Result, error)
\`\`\`

> Scan recursively scans the directory.
```

---

## ライセンス

MIT
