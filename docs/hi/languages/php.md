---
layout: default
title: PHP
parent: भाषा गाइड
nav_order: 13
---

# PHP समर्थन

[English](../../languages/php.md) | [한국어](../ko/languages/php.md) | [日本語](../ja/languages/php.md) | [हिन्दी](php.md) | [Deutsch](../de/languages/php.md)

## समर्थित एक्सटेंशन

- `.php`

## निष्कर्षण लक्ष्य

| तत्व | Kind | उदाहरण |
|------|------|---------|
| फ़ंक्शन | `function` | `function greet($name)` |
| मेथड | `method` | `public function getName()` |
| क्लास | `class` | `class User` |
| इंटरफ़ेस | `interface` | `interface Repository` |
| ट्रेट | `type` | `trait Loggable` |
| एनम | `enum` | `enum Status` |
| कॉन्स्टेंट | `variable` | `const MAX_SIZE = 100;` |
| प्रॉपर्टी | `variable` | `public $name;` |
| PHPDoc टिप्पणी | `doc` | `/** विवरण */` |
| लाइन टिप्पणी | `doc` | `// विवरण` |
| use स्टेटमेंट | `import` | `use App\Services\UserService;` |
| require/include | `import` | `require 'vendor/autoload.php';` |

## उदाहरण

### इनपुट

```php
<?php

namespace App\Services;

use App\Models\User;

/**
 * UserService उपयोगकर्ता-संबंधित कार्यों को संभालता है।
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

### आउटपुट (XML)

```xml
<file path="example.php" language="php">
  <type kind="class" line="10">
    <name>UserService</name>
    <text>class UserService</text>
    <doc>UserService उपयोगकर्ता-संबंधित कार्यों को संभालता है।</doc>
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

## टिप्पणियाँ

### दृश्यता

- सभी घोषणाएं दृश्यता (public, private, protected) की परवाह किए बिना निकाली जाती हैं
- `--include-private` फ्लैग से फ़िल्टर करें (भविष्य की सुविधा)

### बॉडी हटाना

जब `--include-body` फ्लैग का उपयोग नहीं किया जाता है:

- फ़ंक्शन/मेथड: खुलने वाले ब्रेस `{` के बाद बॉडी हटा दी जाती है
- क्लास/इंटरफ़ेस/ट्रेट: क्लास नाम के बाद बॉडी हटा दी जाती है
- कॉन्स्टेंट/प्रॉपर्टी: जैसा है वैसा रखें

### आयात निष्कर्षण

`--include-imports` फ्लैग के साथ:

- `use` स्टेटमेंट: `use App\Services\UserService;`
- `require` स्टेटमेंट: `require 'vendor/autoload.php';`
- `require_once` स्टेटमेंट: `require_once 'config.php';`
- `include` स्टेटमेंट: `include 'file.php';`
- `include_once` स्टेटमेंट: `include_once 'helpers.php';`

### दस्तावेज़ टिप्पणियाँ

- PHPDoc टिप्पणियाँ (`/** ... */`) निकाली जाती हैं
- नियमित बहु-पंक्ति टिप्पणियाँ (`/* ... */`) और लाइन टिप्पणियाँ (`//`, `#`) भी समर्थित हैं
