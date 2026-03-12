# C++ サポート

🌐 [English](../../languages/cpp.md) | [한국어](../../ko/languages/cpp.md) | [日本語](cpp.md) | [हिन्दी](../../hi/languages/cpp.md) | [Deutsch](../../de/languages/cpp.md)

## サポートされる拡張子

- `.cpp`
- `.hpp`
- `.h`

## 抽出対象

| 要素 | Kind | 例 |
|------|------|-----|
| クラス | `class` | `class User { ... }` |
| 構造体 | `struct` | `struct Point { int x, y; }` |
| メソッド | `method` | `void User::getName()` |
| コンストラクタ | `constructor` | `User(string name)` |
| デストラクタ | `destructor` | `~User()` |
| 関数 | `function` | `int add(int a, int b)` |
| 名前空間 | `namespace` | `namespace utils { }` |
| テンプレート | `template` | `template<typename T> class Box` |
| 列挙型 | `enum` | `enum Color { RED, GREEN }` |
| Typedef | `typedef` | `typedef unsigned int uint` |
| マクロ | `macro` | `#define MAX_SIZE 100` |
| Include | (import) | `#include <iostream>` |
| コメント | `doc` | `// Comment` |

## 例

### 入力

```cpp
#include <iostream>
#include <string>

// ユーザーデータを管理するUserクラス
class User {
public:
    User(const std::string& name);
    ~User();

    std::string getName() const;
    void setName(const std::string& name);

private:
    std::string name_;
};

namespace utils {
    // ヘルパー関数
    int calculateHash(const std::string& input);
}

template<typename T>
class Box {
    T value;
public:
    T getValue() const;
};

#define MAX_USERS 100
```

### 出力（XML）

```xml
<file path="example.hpp" language="cpp">
  <signature kind="class" line="5">
    <name>User</name>
    <text>class User</text>
    <doc>ユーザーデータを管理するUserクラス</doc>
  </signature>
  <signature kind="method" line="11">
    <name>getName</name>
    <text>std::string getName() const;</text>
  </signature>
  <signature kind="method" line="12">
    <name>setName</name>
    <text>void setName(const std::string& name);</text>
  </signature>
  <signature kind="namespace" line="18">
    <name>utils</name>
    <text>namespace utils</text>
  </signature>
  <signature kind="function" line="20">
    <name>calculateHash</name>
    <text>int calculateHash(const std::string& input);</text>
    <doc>ヘルパー関数</doc>
  </signature>
  <signature kind="template" line="23">
    <name>Box</name>
    <text>template&lt;typename T&gt; class Box</text>
  </signature>
  <signature kind="macro" line="30">
    <name>MAX_USERS</name>
    <text>#define MAX_USERS 100</text>
  </signature>
</file>
```

## 注意事項

### アクセス制御

- すべてのアクセスレベル（public、private、protected）が抽出される
- 可視性修飾子によるフィルタリングなし
- AIが完全なクラス構造を理解するのに有用

### 本体削除

`--include-body`フラグ未使用時：

- 関数/メソッド：開き波括弧`{`以降の本体を削除
- クラス/構造体/名前空間：本体を削除、宣言のみ保持
- テンプレート：基礎となる宣言本体を削除
- Enum/Typedef/Macro：全文保持

### テンプレートサポート

基本的なテンプレートサポートが含まれる：

```cpp
template<typename T>
class Box { ... };         // キャプチャされる

template<typename T>
T getMax(T a, T b) { ... } // キャプチャされる
```

### 名前空間サポート

単純およびネストされた名前空間の両方がキャプチャされる：

```cpp
namespace outer {
    namespace inner {
        void helper();     // 3つすべてキャプチャ
    }
}
```

### Include文

`--include-imports`を使用して`#include`ディレクティブを抽出。
`--include-private`を使用して非公開/unexportedシンボルを含める。

```cpp
#include <iostream>        // システムinclude
#include "myheader.h"      // ローカルinclude
```

## サポートされていない要素（v1）

| 要素 | 理由 |
|------|------|
| 演算子オーバーロード | `operator+`、`operator<<` - 特殊ケース、稀 |
| Friend宣言 | `friend class Bar` - アクセス制御例外 |
| Using宣言 | `using namespace std` - 単純なエイリアス |
| Lambda式 | `[](int x) { ... }` - インライン定義 |
| テンプレート特殊化 | `template<> class Box<int>` - 複雑なパース |
| 可変引数テンプレート | `template<typename... Args>` - 高度なパターン |
| C++20 Concepts | `template<Integral T>` - 限定的なコンパイラサポート |
| C++20 Modules | `import std;` - 限定的なコンパイラサポート |
| グローバル変数 | 将来のバージョンで追加される可能性あり |
