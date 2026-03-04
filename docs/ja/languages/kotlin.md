# Kotlin サポート

[English](../../languages/kotlin.md) | [한국어](../../ko/languages/kotlin.md) | [日本語](kotlin.md) | [हिन्दी](../../hi/languages/kotlin.md) | [Deutsch](../../de/languages/kotlin.md)

## 対応拡張子

- `.kt`
- `.kts`

## 抽出対象

| 要素 | Kind | 例 |
|------|------|------|
| 関数 | `function` | `fun add(a: Int, b: Int): Int` |
| Suspend 関数 | `function` | `suspend fun fetchData(url: String): String` |
| 拡張関数 | `function` | `fun String.isEmail(): Boolean` |
| クラス | `class` | `class User(val name: String)` |
| データクラス | `class` | `data class Point(val x: Double, val y: Double)` |
| Sealed クラス | `class` | `sealed class Result<out T>` |
| Enum クラス | `enum` | `enum class Color { RED, GREEN, BLUE }` |
| インターフェース | `interface` | `interface Repository<T>` |
| オブジェクト | `class` | `object AppConfig` |
| コンパニオンオブジェクト | `class` | `companion object Factory` |
| プロパティ (val/var) | `variable` | `val MAX_SIZE = 100` |
| 型エイリアス | `type` | `typealias Handler<T> = (T) -> Unit` |
| Enum エントリ | `variable` | `RED("#FF0000")` |
| セカンダリコンストラクタ | `constructor` | `constructor(name: String)` |
| ドキュメントコメント | `doc` | `/** ドキュメント */` |

## 例

### 入力

```kotlin
/** APIレスポンス用のユーザーデータクラス。 */
data class User(
    val id: Long,
    val name: String,
    val email: String
) {
    fun isValid(): Boolean = email.contains("@")
}

/** ユーザー操作用のリポジトリインターフェース。 */
interface UserRepository {
    suspend fun getUser(id: Long): User?
    fun save(user: User): Boolean
}

val DEFAULT_TIMEOUT: Long = 5000L
```

### 出力 (XML)

```xml
<file path="user.kt" language="kotlin">
  <type kind="class" line="2">
    <name>User</name>
    <text>data class User(
    val id: Long,
    val name: String,
    val email: String
)</text>
    <doc>APIレスポンス用のユーザーデータクラス。</doc>
  </type>
  <function kind="function" line="7">
    <name>isValid</name>
    <text>fun isValid(): Boolean = email.contains("@")</text>
  </function>
  <type kind="interface" line="11">
    <name>UserRepository</name>
    <text>interface UserRepository</text>
    <doc>ユーザー操作用のリポジトリインターフェース。</doc>
  </type>
  <function kind="function" line="12">
    <name>getUser</name>
    <text>suspend fun getUser(id: Long): User?</text>
  </function>
  <function kind="function" line="13">
    <name>save</name>
    <text>fun save(user: User): Boolean</text>
  </function>
  <variable kind="variable" line="16">
    <name>DEFAULT_TIMEOUT</name>
    <text>val DEFAULT_TIMEOUT: Long = 5000L</text>
  </variable>
</file>
```

## 注意事項

### 可視性

- すべての宣言が抽出されます（Kotlinはデフォルトで `public`）
- アクセス修飾子（`public`、`internal`、`private`、`protected`）はシグネチャにそのまま保持されます

### 関数修飾子

- `suspend`、`inline`、`infix`、`operator`、`tailrec` 関数はすべて kind `function` として分類されます
- 修飾子はシグネチャテキストに保持されます
- 単一式関数（`fun double(x: Int) = x * 2`）は完全に保持されます

### ジェネリクス

- ジェネリック型パラメータ（`<T>`、`<T : Comparable<T>>`）が完全に保持されます
- `where` 句と分散アノテーション（`in`、`out`）がシグネチャに含まれます
- `reified` 型パラメータも保持されます

### クラス

- `data class`、`sealed class`、`abstract class`、`open class`、`inner class`、`annotation class`、`value class` はすべて kind `class` として分類されます
- `enum class` は kind `enum` として分類されます
- `interface` と `sealed interface` は kind `interface` として分類されます

### オブジェクト

- `object` 宣言（シングルトン）は kind `class` として分類されます
- `companion object` ブロックが抽出されます。名前のないコンパニオンには合成名「Companion」が付与されます

### 本体の削除

`--include-body` フラグを使用しない場合:

- 関数/メソッド: 開き括弧 `{` 以降の本体を削除
- 単一式関数: 完全に保持（式自体がシグネチャの一部）
- クラス/インターフェース/列挙型: 開き括弧 `{` 以降の本体を削除
- プロパティ (val/var): 値式は保持
- 型エイリアス: 完全に保持

### ドキュメントコメント

- `/** ... */`（KDoc）と `//` 行コメントの両方が抽出されます
- KDocコメントは後続の宣言に関連付けられます
