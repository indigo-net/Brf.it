# Brf.it

🌐 [English](../../README.md) | [한국어](../ko/README.md) | [日本語](README.md) | [हिन्दी](../hi/README.md) | [Deutsch](../de/README.md)

[![Release](https://img.shields.io/github/v/release/indigo-net/Brf.it)](https://github.com/indigo-net/Brf.it/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/indigo-net/Brf.it)](https://goreportcard.com/report/github.com/indigo-net/Brf.it)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **コードベースをAIが理解しやすい形式にパッケージング**
>
> `50トークン → 8トークン` — 同じ情報、より少ないトークン。

[インストール](#インストール) · [クイックスタート](#クイックスタート) · [サポート言語](#サポート言語)

---

## 動作原理

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

## クイックスタート

### インストール

**macOS（Homebrew）**

```bash
brew install indigo-net/tap/brfit
```

**Linux / macOS（スクリプト）**

```bash
curl -fsSL https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.sh | sh
```

**Windows（PowerShell）**

```powershell
irm https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.ps1 | iex
```

**ソースからビルド**

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o brfit ./cmd/brfit
```

### 初回実行

```bash
brfit .                    # 現在のディレクトリを分析
brfit . -f md              # Markdown形式で出力
brfit . -o briefing.xml    # ファイルに保存
```

---

## 実際の使用例

**[SAMPLE.md](SAMPLE.md)** | **[SAMPLE.xml](SAMPLE.xml)**

brfit自体でパッケージングしたこのプロジェクトです。コミットごとに自動生成されます。

---

## 機能

| 機能 | 説明 |
|------|------|
| Tree-sitterベース | 正確なASTパースで言語構造を分析 |
| 複数フォーマット | XML、Markdown出力をサポート |
| トークンカウント | 出力トークン数を自動計算 |
| gitignore対応 | 不要なファイルを自動除外 |
| クロスプラットフォーム | Linux、macOS、Windowsをサポート |

---

## サポート言語

| 言語 | 拡張子 | ドキュメント |
|------|--------|--------------|
| Go | `.go` | [Goガイド](languages/go.md) |
| TypeScript | `.ts`、`.tsx` | [TypeScriptガイド](languages/typescript.md) |
| JavaScript | `.js`、`.jsx` | [TypeScriptガイド](languages/typescript.md) |
| Python | `.py` | [Pythonガイド](languages/python.md) |
| C | `.c`、`.h` | [Cガイド](languages/c.md) |
| Java | `.java` | [Javaガイド](languages/java.md) |

---

## CLIリファレンス

```bash
brfit [パス] [オプション]
```

### オプション

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

### 例

```bash
# AIアシスタントに渡す（クリップボードにコピー）
brfit . | pbcopy              # macOS
brfit . | xclip               # Linux
brfit . | clip                # Windows

# プロジェクトを分析してファイルに保存
brfit ./my-project -o briefing.xml

# 関数本体を含める（完全なコード）
brfit . --include-body

# ディレクトリツリー出力をスキップ
brfit . --no-tree
```

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

MITライセンス — 個人・商用プロジェクトで自由に使用できます。

詳細は[LICENSE](LICENSE)をご覧ください。
