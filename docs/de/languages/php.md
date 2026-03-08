---
layout: default
title: PHP
parent: Sprachleitfäden
nav_order: 13
---

# PHP-Unterstützung

[English](../../languages/php.md) | [한국어](../ko/languages/php.md) | [日本語](../ja/languages/php.md) | [हिन्दी](../hi/languages/php.md) | [Deutsch](php.md)

## Unterstützte Erweiterungen

- `.php`

## Extraktionsziele

| Element | Kind | Beispiel |
|---------|------|----------|
| Funktion | `function` | `function greet($name)` |
| Methode | `method` | `public function getName()` |
| Klasse | `class` | `class User` |
| Interface | `interface` | `interface Repository` |
| Trait | `type` | `trait Loggable` |
| Enum | `enum` | `enum Status` |
| Konstante | `variable` | `const MAX_SIZE = 100;` |
| Eigenschaft | `variable` | `public $name;` |
| PHPDoc-Kommentar | `doc` | `/** Beschreibung */` |
| Zeilenkommentar | `doc` | `// Beschreibung` |
| use-Anweisung | `import` | `use App\Services\UserService;` |
| require/include | `import` | `require 'vendor/autoload.php';` |

## Beispiel

### Eingabe

```php
<?php

namespace App\Services;

use App\Models\User;

/**
 * UserService verarbeitet benutzerbezogene Operationen.
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

### Ausgabe (XML)

```xml
<file path="example.php" language="php">
  <type kind="class" line="10">
    <name>UserService</name>
    <text>class UserService</text>
    <doc>UserService verarbeitet benutzerbezogene Operationen.</doc>
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

## Hinweise

### Sichtbarkeit

- Alle Deklarationen werden unabhängig von der Sichtbarkeit (public, private, protected) extrahiert
- Filtern mit `--include-private` Flag (zukünftige Funktion)

### Body-Entfernung

Wenn das `--include-body` Flag nicht verwendet wird:

- Funktionen/Methoden: Body nach der öffnenden Klammer `{` entfernt
- Klassen/Interfaces/Traits: Body nach dem Klassennamen entfernt
- Konstanten/Eigenschaften: bleiben erhalten

### Import-Extraktion

Mit dem `--include-imports` Flag:

- `use`-Anweisungen: `use App\Services\UserService;`
- `require`-Anweisungen: `require 'vendor/autoload.php';`
- `require_once`-Anweisungen: `require_once 'config.php';`
- `include`-Anweisungen: `include 'file.php';`
- `include_once`-Anweisungen: `include_once 'helpers.php';`

### Dokumentationskommentare

- PHPDoc-Kommentare (`/** ... */`) werden extrahiert
- Reguläre mehrzeilige Kommentare (`/* ... */`) und Zeilenkommentare (`//`, `#`) werden ebenfalls unterstützt
