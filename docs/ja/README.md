<p align="center">
  <img src="../../assets/logo.png" alt="Brf.it logo" width="128">
</p>

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

<br/>

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
<function>
  export async function fetchUser(
    id: string
  ): Promise<User>
</function>
```

</td>
</tr>
</table>

---

<br/>

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

<br/>

## 実際の使用例

**[SAMPLE.md](SAMPLE.md)** | **[SAMPLE.xml](SAMPLE.xml)**

brfit自体でパッケージングしたこのプロジェクトです。コミットごとに自動生成されます。

---

<br/>

## 機能

| 機能 | 説明 |
|------|------|
| Tree-sitterベース | 正確なASTパースで言語構造を分析 |
| 複数フォーマット | XML、Markdown出力をサポート |
| トークンカウント | 出力トークン数を自動計算 |
| gitignore対応 | 不要なファイルを自動除外 |
| クロスプラットフォーム | Linux、macOS、Windowsをサポート |

---

<br/>

## サポート言語

| 言語 | 拡張子 | ドキュメント |
|------|--------|--------------|
| Go | `.go` | [Goガイド](languages/go.md) |
| TypeScript | `.ts`、`.tsx` | [TypeScriptガイド](languages/typescript.md) |
| JavaScript | `.js`、`.jsx` | [TypeScriptガイド](languages/typescript.md) |
| Python | `.py` | [Pythonガイド](languages/python.md) |
| C | `.c`、`.h` | [Cガイド](languages/c.md) |
| Java | `.java` | [Javaガイド](languages/java.md) |
| Rust | `.rs` | [Rustガイド](languages/rust.md) |
| Swift | `.swift` | [Swiftガイド](languages/swift.md) |
| Kotlin | `.kt`, `.kts` | [Kotlinガイド](languages/kotlin.md) |
| C# | `.cs` | [C#ガイド](languages/csharp.md) |
| Lua | `.lua` | [Luaガイド](languages/lua.md) |
| PHP | `.php` | [PHPガイド](languages/php.md) |
| Ruby | `.rb` | [Rubyガイド](languages/ruby.md) |
| Scala | `.scala`, `.sc` | [Scalaガイド](languages/scala.md) |
| Elixir | `.ex`, `.exs` | [Elixirガイド](languages/elixir.md) |
| SQL | `.sql` | [SQLガイド](languages/sql.md) |
| YAML | `.yaml`、`.yml` | [YAMLガイド](languages/yaml.md) |
| TOML | `.toml` | [TOMLガイド](languages/toml.md) |

---

<br/>

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
| `--include-imports` | | import文を含める | `false` |
| `--include-private` | | 非公開/unexportedシンボルを含める | `false` |
| `--ignore` | `-i` | ignoreファイルパス（複数回指定可能） | `.gitignore` |
| `--include` | | 含めるglobパターン（複数指定可能） | |
| `--exclude` | | 除外するglobパターン（複数指定可能） | |
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

# importを含める（そのまま）
brfit . --include-imports
```

---

<br/>

## ライセンス

MITライセンス — 個人・商用プロジェクトで自由に使用できます。

詳細は[LICENSE](LICENSE)をご覧ください。
