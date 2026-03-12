# Python サポート

🌐 [English](../../languages/python.md) | [한국어](../../ko/languages/python.md) | [日本語](python.md) | [हिन्दी](../../hi/languages/python.md) | [Deutsch](../../de/languages/python.md)

## サポート拡張子

- `.py`

## 抽出対象

| 要素 | Kind | 例 |
|------|------|-----|
| 関数 | `function` | `def greet():` |
| 非同期関数 | `function` | `async def fetch():` |
| クラス | `class` | `class User:` |
| メソッド | `method` | `def __init__(self):` |
| クラスメソッド | `method` | `def method(cls):` |
| モジュールレベル変数 | `variable` | `API_URL = "..."` |
| コメント | `doc` | `# Comment` |

## 例

### 入力

```python
# User model for the application.
class User:
    """Represents a user in the system."""

    def __init__(self, name: str, email: str):
        """Initialize a new user."""
        self.name = name
        self.email = email

    def get_display_name(self) -> str:
        """Return the display name."""
        return f"{self.name} <{self.email}>"


# Create a new user instance.
def create_user(name: str, email: str) -> User:
    """Factory function to create a user."""
    return User(name, email)


async def fetch_user(user_id: int) -> User:
    """Fetch a user from the database."""
    pass
```

### 出力（XML）

```xml
<file path="user.py" language="python">
  <signature kind="class" line="2">
    <name>User</name>
    <text>class User</text>
    <doc>User model for the application.</doc>
  </signature>
  <signature kind="method" line="5">
    <name>__init__</name>
    <text>def __init__(self, name: str, email: str)</text>
  </signature>
  <signature kind="method" line="10">
    <name>get_display_name</name>
    <text>def get_display_name(self) -> str</text>
  </signature>
  <signature kind="function" line="16">
    <name>create_user</name>
    <text>def create_user(name: str, email: str) -> User</text>
    <doc>Create a new user instance.</doc>
  </signature>
  <signature kind="function" line="21">
    <name>fetch_user</name>
    <text>async def fetch_user(user_id: int) -> User</text>
  </signature>
</file>
```

## 注意事項

### エクスポート判定

- Pythonはすべての要素をpublicとして扱う
- `_private`または`__mangled`命名規則も含まれる

### メソッド vs 関数の判定

- 最初のパラメータが`self`または`cls`の場合は`method`に分類
- それ以外は`function`に分類
- `@staticmethod`デコレータが付いたメソッドは`function`に分類（selfなし）

### 非同期処理

- `async def`は`function` kindに統一
- シグネチャテキストに`async`キーワードを含む

### 本体削除

`--include-body`フラグ未使用時：

- 関数/メソッド：シグネチャ末尾のコロン（`:`）以降の本体を削除
- クラス：クラス名と継承情報のみ保持

`--include-private`を使用して非公開/unexportedシンボルを含める。

### 型ヒント内のコロン処理

複雑な型ヒント（例：`Dict[str, int]`）のコロンは関数終端のコロンと区別される：

```python
def func(x: Dict[str, List[int]]) -> str:  # 最後のコロンのみ関数終端
```

### Docstring（将来サポート）

- 現在のバージョン：関数/クラス上の`#`コメントのみdocとしてキャプチャ
- 将来のバージョン：トリプルクォートdocstring（`"""..."""`）サポート予定

### デコレータ

デコレータはシグネチャに含まれない：

```python
@decorator
def func():  # シグネチャ: "def func()"
```

### サポートされていない要素

- Lambda式
- ネストされた関数（外側の関数のみキャプチャ）
