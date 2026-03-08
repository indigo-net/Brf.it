---
layout: default
title: PHP
parent: 言語ガイド
nav_order: 13
---

# PHP サポート

[English](../../languages/php.md) | [한국어](../ko/languages/php.md) | [日本語](php.md) | [हिन्दी](../hi/languages/php.md) | [Deutsch](../de/languages/php.md)

## サポート拡張子

- `.php`

## 抽出対象

| 要素 | Kind | 例 |
|------|------|------|
| 関数 | `function` | `function greet($name)` |
| メソッド | `method` | `public function getName()` |
| クラス | `class` | `class User` |
| インターフェース | `interface` | `interface Repository` |
| トレイト | `type` | `trait Loggable` |
| 列挙型 | `enum` | `enum Status` |
| 定数 | `variable` | `const MAX_SIZE = 100;` |
| プロパティ | `variable` | `public $name;` |
| PHPDoc コメント | `doc` | `/** 説明 */` |
| 行コメント | `doc` | `// 説明` |
| use 文 | `import` | `use App\Services\UserService;` |
| require/include | `import` | `require 'vendor/autoload.php';` |

## 例

### 入力

```php
<?php

namespace App\Services;

use App\Models\User;

/**
 * UserService はユーザー関連の操作を処理します。
 */
class UserService {
    private $repository;

    public function __construct($repo) {
        $this->repository = $repo;
    }

    public function findUser($id): ?User {
        return $this->repository->find($id);
    }
}

interface RepositoryInterface {
    public function find($id);
    public function save($entity);
}

trait Loggable {
    public function log($message) {
        echo $message;
    }
}

function helper($data) {
    return $data;
}

const MAX_ITEMS = 100;
const APP_NAME = "Brf.it";
```

### 出力 (XML)

```xml
<file path="example.php" language="php">
  <type kind="class" line="10">
    <name>UserService</name>
    <text>class UserService</text>
    <doc>UserService はユーザー関連の操作を処理します。</doc>
  </type>
  <function kind="method" line="15">
    <name>__construct</name>
    <text>public function __construct($repo)</text>
  </function>
  <function kind="method" line="19">
    <name>findUser</name>
    <text>public function findUser($id): ?User</text>
  </function>
  <type kind="interface" line="25">
    <name>RepositoryInterface</name>
    <text>interface RepositoryInterface</text>
  </type>
  <function kind="method" line="26">
    <name>find</name>
    <text>public function find($id)</text>
  </function>
  <function kind="method" line="27">
    <name>save</name>
    <text>public function save($entity)</text>
  </function>
  <type kind="type" line="30">
    <name>Loggable</name>
    <text>trait Loggable</text>
  </type>
  <function kind="method" line="31">
    <name>log</name>
    <text>public function log($message)</text>
  </function>
  <function kind="function" line="36">
    <name>helper</name>
    <text>function helper($data)</text>
  </function>
  <variable kind="variable" line="40">
    <name>MAX_ITEMS</name>
    <text>const MAX_ITEMS = 100;</text>
  </variable>
  <variable kind="variable" line="41">
    <name>APP_NAME</name>
    <text>const APP_NAME = "Brf.it";</text>
  </variable>
</file>
```

## 注意事項

### 可視性

- すべての宣言が可視性（public、private、protected）に関係なく抽出されます
- `--include-private` フラグでフィルタリング（将来の機能）

### 本体の削除

`--include-body` フラグを使用しない場合：

- 関数/メソッド：開き括弧 `{` 以降の本体が削除されます
- クラス/インターフェース/トレイト：クラス名以降の本体が削除されます
- 定数/プロパティ：そのまま維持されます

### インポートの抽出

`--include-imports` フラグを使用する場合：

- `use` 文：`use App\Services\UserService;`
- `require` 文：`require 'vendor/autoload.php';`
- `require_once` 文：`require_once 'config.php';`
- `include` 文：`include 'file.php';`
- `include_once` 文：`include_once 'helpers.php';`

### ドキュメントコメント

- PHPDoc コメント（`/** ... */`）が抽出されます
- 通常の複数行コメント（`/* ... */`）と行コメント（`//`、`#`）もサポートされています
