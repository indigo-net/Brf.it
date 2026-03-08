---
layout: default
title: PHP
parent: 언어 가이드
nav_order: 13
---

# PHP 지원

[English](../../languages/php.md) | [한국어](php.md) | [日本語](../ja/languages/php.md) | [हिन्दी](../hi/languages/php.md) | [Deutsch](../de/languages/php.md)

## 지원 확장자

- `.php`

## 추출 대상

| 요소 | Kind | 예시 |
|------|------|------|
| 함수 | `function` | `function greet($name)` |
| 메서드 | `method` | `public function getName()` |
| 클래스 | `class` | `class User` |
| 인터페이스 | `interface` | `interface Repository` |
| 트레이트 | `type` | `trait Loggable` |
| 열거형 | `enum` | `enum Status` |
| 상수 | `variable` | `const MAX_SIZE = 100;` |
| 프로퍼티 | `variable` | `public $name;` |
| PHPDoc 주석 | `doc` | `/** 설명 */` |
| 한 줄 주석 | `doc` | `// 설명` |
| use 문 | `import` | `use App\Services\UserService;` |
| require/include | `import` | `require 'vendor/autoload.php';` |

## 예시

### 입력

```php
<?php

namespace App\Services;

use App\Models\User;

/**
 * UserService는 사용자 관련 작업을 처리합니다.
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

### 출력 (XML)

```xml
<file path="example.php" language="php">
  <type kind="class" line="10">
    <name>UserService</name>
    <text>class UserService</text>
    <doc>UserService는 사용자 관련 작업을 처리합니다.</doc>
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

## 참고사항

### 가시성

- 모든 선언이 가시성(public, private, protected)에 관계없이 추출됩니다
- `--include-private` 플래그로 필터링 (향후 기능)

### 본문 제거

`--include-body` 플래그를 사용하지 않을 때:

- 함수/메서드: 여는 중괄호 `{` 이후의 본문이 제거됩니다
- 클래스/인터페이스/트레이트: 클래스 이름 이후의 본문이 제거됩니다
- 상수/프로퍼티: 그대로 유지됩니다

### Import 추출

`--include-imports` 플래그 사용 시:

- `use` 문: `use App\Services\UserService;`
- `require` 문: `require 'vendor/autoload.php';`
- `require_once` 문: `require_once 'config.php';`
- `include` 문: `include 'file.php';`
- `include_once` 문: `include_once 'helpers.php';`

### 문서 주석

- PHPDoc 주석(`/** ... */`)이 추출됩니다
- 일반 여러 줄 주석(`/* ... */`)과 한 줄 주석(`//`, `#`)도 지원됩니다
