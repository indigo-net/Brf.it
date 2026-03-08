---
layout: default
title: PHP
parent: Language Guides
nav_order: 13
---

# PHP Support

[English](php.md) | [한국어](../ko/languages/php.md) | [日本語](../ja/languages/php.md) | [हिन्दी](../hi/languages/php.md) | [Deutsch](../de/languages/php.md)

## Supported Extensions

- `.php`

## Extraction Targets

| Element | Kind | Example |
|---------|------|---------|
| Function | `function` | `function greet($name)` |
| Method | `method` | `public function getName()` |
| Class | `class` | `class User` |
| Interface | `interface` | `interface Repository` |
| Trait | `type` | `trait Loggable` |
| Enum | `enum` | `enum Status` |
| Constant | `variable` | `const MAX_SIZE = 100;` |
| Property | `variable` | `public $name;` |
| PHPDoc Comment | `doc` | `/** Description */` |
| Line Comment | `doc` | `// Description` |
| use Statement | `import` | `use App\Services\UserService;` |
| require/include | `import` | `require 'vendor/autoload.php';` |

## Example

### Input

```php
<?php

namespace App\Services;

use App\Models\User;

/**
 * UserService handles user-related operations.
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

### Output (XML)

```xml
<file path="example.php" language="php">
  <type kind="class" line="10">
    <name>UserService</name>
    <text>class UserService</text>
    <doc>UserService handles user-related operations.</doc>
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

## Notes

### Visibility

- All declarations are extracted regardless of visibility (public, private, protected)
- Use `--include-private` flag to filter (future feature)

### Body Removal

When `--include-body` flag is not used:

- Functions/Methods: body removed after opening brace `{`
- Classes/Interfaces/Traits: body removed after class name
- Constants/Properties: preserved as-is

### Import Extraction

With `--include-imports` flag:

- `use` statements: `use App\Services\UserService;`
- `require` statements: `require 'vendor/autoload.php';`
- `require_once` statements: `require_once 'config.php';`
- `include` statements: `include 'file.php';`
- `include_once` statements: `include_once 'helpers.php';`

### Doc Comments

- PHPDoc comments (`/** ... */`) are extracted
- Regular multi-line comments (`/* ... */`) and line comments (`//`, `#`) are also supported
