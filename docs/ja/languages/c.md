# C サポート

🌐 [English](../../languages/c.md) | [한국어](../../ko/languages/c.md) | [日本語](c.md) | [हिन्दी](../../hi/languages/c.md) | [Deutsch](../../de/languages/c.md)

## サポートされる拡張子

- `.c`
- `.h`

## 抽出対象

| 要素 | Kind | 例 |
|------|------|-----|
| 関数定義 | `function` | `int add(int a, int b) { ... }` |
| 関数宣言 | `function` | `int add(int a, int b);` |
| 構造体 | `struct` | `struct User { ... };` |
| 列挙型 | `enum` | `enum Color { RED, GREEN, BLUE };` |
| Typedef | `typedef` | `typedef struct { ... } User;` |
| グローバル変数 | `variable` | `int global_count = 0;` |
| オブジェクトマクロ | `macro` | `#define MAX_SIZE 100` |
| 関数マクロ | `macro` | `#define MIN(a, b) ((a) < (b) ? (a) : (b))` |
| コメント | `doc` | `// Comment` |

## 例

### 入力

```c
// User構造体
typedef struct {
    int id;
    char name[64];
} User;

// 新しいユーザーを作成
User* create_user(const char* name);

// 内部ヘルパー
static void init_user(User* u);

#define MAX_USERS 100
#define INIT_USER(u) memset(u, 0, sizeof(User))
```

### 出力 (XML)

```xml
<file path="example.h" language="c">
  <signature kind="typedef" line="2">
    <name>User</name>
    <text>typedef struct { int id; char name[64]; } User;</text>
    <doc>User構造体</doc>
  </signature>
  <signature kind="function" line="8" exported="true">
    <name>create_user</name>
    <text>User* create_user(const char* name);</text>
    <doc>新しいユーザーを作成</doc>
  </signature>
  <signature kind="function" line="11" exported="true">
    <name>init_user</name>
    <text>static void init_user(User* u);</text>
    <doc>内部ヘルパー</doc>
  </signature>
  <signature kind="macro" line="13">
    <name>MAX_USERS</name>
    <text>#define MAX_USERS 100</text>
  </signature>
  <signature kind="macro" line="14">
    <name>INIT_USER</name>
    <text>#define INIT_USER(u) memset(u, 0, sizeof(User))</text>
  </signature>
</file>
```

## 注意事項

### Export検出

- すべてのC関数はデフォルトでexportedとして扱われます
- `static`関数も含まれます（将来：`exported: false`としてマークされる可能性があります）

### 本文削除

`--include-body`フラグを使用しない場合：

- 関数：開き中括弧`{`以降の本文を削除
- Struct/Enum/Typedef/Macro：全文を保持

`--include-private`を使用して非公開/unexportedシンボルを含める。

### ポインタ戻り型

直接戻り型とポインタ戻り型の両方をサポート：

```c
int get_value();        // 直接戻り型
User* create_user();    // ポインタ戻り型
```

### サポートされない要素

- 関数ポインタ（変数として）
- ネストされた構造体（トップレベルのみ抽出）
