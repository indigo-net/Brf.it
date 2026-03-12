# Java サポート

🌐 [English](../../languages/java.md) | [한국어](../../ko/languages/java.md) | [日本語](java.md) | [हिन्दी](../../hi/languages/java.md) | [Deutsch](../../de/languages/java.md)

## サポート拡張子

- `.java`

## 抽出対象

| 要素 | Kind | 例 |
|------|------|-----|
| クラス | `class` | `public class User { ... }` |
| インターフェース | `interface` | `public interface Repository<T> { ... }` |
| メソッド | `method` | `public String getName() { ... }` |
| コンストラクタ | `constructor` | `public User(String name) { ... }` |
| Enum | `enum` | `public enum Status { ... }` |
| アノテーション | `annotation` | `public @interface Inject { ... }` |
| Record (Java 14+) | `record` | `public record Point(int x, int y) { ... }` |
| フィールド | `field` | `public static final String API = "..."` |
| コメント | `doc` | `// Comment` または `/* Block */` |

## 例

### 入力

```java
package com.example;

/**
 * User class represents a user in the system.
 */
public class User {
    private String name;

    public User(String name) {
        this.name = name;
    }

    public String getName() {
        return name;
    }

    private void internalMethod() {
        // Private method
    }
}

public interface Repository<T> {
    T findById(String id);
    void save(T entity);
}

public enum Status {
    PENDING, ACTIVE, COMPLETED
}

public record Point(int x, int y) {}
```

### 出力（XML）

```xml
<file path="User.java" language="java">
  <signature kind="class" line="6">
    <name>User</name>
    <text>public class User</text>
  </signature>
  <signature kind="constructor" line="9">
    <name>User</name>
    <text>public User(String name)</text>
  </signature>
  <signature kind="method" line="13">
    <name>getName</name>
    <text>public String getName()</text>
  </signature>
  <signature kind="interface" line="22">
    <name>Repository</name>
    <text>public interface Repository&lt;T&gt;</text>
  </signature>
  <signature kind="method" line="23">
    <name>findById</name>
    <text>T findById(String id);</text>
  </signature>
  <signature kind="method" line="24">
    <name>save</name>
    <text>void save(T entity);</text>
  </signature>
  <signature kind="enum" line="27">
    <name>Status</name>
    <text>public enum Status</text>
  </signature>
  <signature kind="record" line="31">
    <name>Point</name>
    <text>public record Point(int x, int y)</text>
  </signature>
</file>
```

## 注意事項

### 可視性フィルタリング

- `public`、`protected`、package-private（デフォルト）：デフォルトで抽出
- `private`：`--include-body`使用時のみ含まれる

### ジェネリクス処理

ジェネリック型パラメータはシグネチャに含まれる：

```java
public class Box<T extends Comparable<T>>  // 完全にキャプチャ
public <U> U transform(Function<T, U> fn)  // メソッド型パラメータを含む
```

### アノテーション出力

メソッドとクラスのアノテーションはシグネチャテキストに含まれる：

```java
@Override
public String toString()  // シグネチャに@Overrideを含む
```

### Record サポート（Java 14+）

Recordはコンポーネントパラメータと共に抽出される：

```java
public record User(String name, int age)  // コンポーネント保持
```

### 内部/ネストクラス

すべてのネストクラスは個別のシグネチャとして抽出される：

```java
public class Outer {
    public static class Nested { ... }  // 個別に抽出
    public class Inner { ... }          // 同様に抽出
}
```

### 抽象メソッド

インターフェースの抽象メソッドは`;`で終わる（本体なし）：

```java
interface Foo {
    void bar();  // そのままキャプチャ
}
```

### 本体削除

`--include-body`フラグ未使用時：

- メソッド/コンストラクタ：開き波括弧`{`以降の本体を削除
- クラス/インターフェース/Enum：開き波括弧`{`以降の本体を削除
- 抽象メソッド：そのまま保持（`;`で終わる）

`--include-private`を使用して非公開/unexportedシンボルを含める。

### Javadoc（将来サポート）

- 現在のバージョン：宣言上の`//`および`/* */`コメントがdocとしてキャプチャ
- 将来のバージョン：`/** */` Javadoc解析サポート予定

### サポートされていない要素

- staticイニシャライザブロック
- 匿名クラス
- Lambda式
