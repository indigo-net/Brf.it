# TypeScript サポート

🌐 [English](typescript.md) | [한국어](typescript.ko.md) | [日本語](typescript.ja.md) | [हिन्दी](typescript.hi.md) | [Deutsch](typescript.de.md)

## サポート拡張子

- `.ts`
- `.tsx`
- `.js`（JavaScript）
- `.jsx`（JSX）

## 抽出対象

| 要素 | Kind | 例 |
|------|------|-----|
| 関数宣言 | `function` | `function greet()` |
| アロー関数 | `arrow` | `const greet = () => {}` |
| メソッド | `method` | `class A { method() {} }` |
| クラス | `class` | `class User {}` |
| インターフェース | `interface` | `interface Props {}` |
| 型エイリアス | `type` | `type ID = string` |
| コメント | `doc` | `// Comment` |

## 例

### 入力

```typescript
/**
 * User interface representing a user entity.
 */
interface User {
    id: number;
    name: string;
}

/**
 * Creates a new user with the given name.
 */
function createUser(name: string): User {
    return { id: Date.now(), name };
}

/**
 * User service for managing users.
 */
class UserService {
    private users: User[] = [];

    /**
     * Adds a user to the service.
     */
    addUser(user: User): void {
        this.users.push(user);
    }
}

// Arrow function example
const formatName = (user: User): string => {
    return user.name.toUpperCase();
};
```

### 出力（XML）

```xml
<file path="user.ts" language="typescript">
  <signature kind="interface" line="4">
    <name>User</name>
    <text>interface User</text>
    <doc>User interface representing a user entity.</doc>
  </signature>
  <signature kind="function" line="12">
    <name>createUser</name>
    <text>function createUser(name: string): User</text>
    <doc>Creates a new user with the given name.</doc>
  </signature>
  <signature kind="class" line="18">
    <name>UserService</name>
    <text>class UserService</text>
    <doc>User service for managing users.</doc>
  </signature>
  <signature kind="method" line="24">
    <name>addUser</name>
    <text>addUser(user: User): void</text>
    <doc>Adds a user to the service.</doc>
  </signature>
  <signature kind="arrow" line="30">
    <name>formatName</name>
    <text>const formatName = (user: User): string => </text>
    <doc>Arrow function example</doc>
  </signature>
</file>
```

## 注意事項

### エクスポート判定

- `export`キーワードがある要素のみ抽出（デフォルト）
- JavaScriptファイルの場合はすべての要素を抽出

### アロー関数

- `const`/`let`/`var`で宣言されたアロー関数をキャプチャ
- 変数名を関数名として使用
- `export const`形式もサポート

### 本体削除

`--include-body`フラグ未使用時：

- 関数/メソッド：中括弧`{`以降の本体を削除
- アロー関数：`=>`以降の本体を削除
- クラス/インターフェース：中括弧`{`以降の内容を削除

### JSDocサポート

- `/** ... */`形式のJSDocコメントを自動リンク
- 関数/クラス直前のコメントがdocとしてキャプチャされる

### JavaScript互換性

- `.js`、`.jsx`ファイルはTypeScriptパーサーで処理
- 型情報がなくても関数/クラスを抽出可能
